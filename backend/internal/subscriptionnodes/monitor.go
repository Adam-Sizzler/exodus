package subscriptionnodes

import (
	"bytes"
	"context"
	"crypto/tls"
	"crypto/x509"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math/rand/v2"
	"net/http"
	"net/http/httptest"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"

	"exodus/internal/config"
	dbmanager "exodus/internal/db/manager"
	subscriptionapi "exodus/internal/httpapi/subscription"
	systemapi "exodus/internal/httpapi/system"
	"exodus/internal/proto"
	srscore "exodus/internal/srslists"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

const exodusRealIPHeader = "x-exodus-real-ip"
const subGRPCTokenHeader = "x-exodus-grpc-token"

const (
	bridgeOperationSubscriptionInfo    = "subscription_info"
	bridgeOperationSubscriptionContent = "subscription_content"
	bridgeOperationSubpageByShortUUID  = "subpage_config_for_short"
	bridgeOperationSubpageByUUID       = "subpage_config_by_uuid"
)

const (
	subNodeRuntimeStatVersion  = "sub_node_version"
	subNodeRuntimeStatUptime   = "sub_node_uptime"
	subNodeRuntimeStatCPUCount = "cpu_count"
	subNodeRuntimeStatCPUModel = "cpu_model"
	subNodeRuntimeStatTotalRAM = "total_ram"
)

var subNodeVersionPattern = regexp.MustCompile(`^[vV]?\d+\.\d+\.\d+(?:[-+][0-9A-Za-z\.-]+)?$`)

const (
	subNodeStreamInterval      = 20 * time.Second
	subNodeStreamIdleTimeout   = 75 * time.Second
	subNodeStreamWatchInterval = 5 * time.Second
)

type SubNodeMonitor struct {
	manager *dbmanager.DatabaseManager
	cfg     *config.BackendConfig

	nodes     map[string]*subNodeState
	nodesLock sync.RWMutex

	runtimeByNodeName map[string]SubNodeRuntimeSnapshot
	runtimeMu         sync.RWMutex

	globalCtx    context.Context
	globalCancel context.CancelFunc

	syncNow        chan struct{}
	deployNow      chan []string
	subpagePushNow chan subpageConfigPushCommand
	srsSyncNow     chan []string
}

type subpageConfigPushCommand struct {
	uuid        string
	config      []byte
	targetUUIDs []string
}

type subNodeState struct {
	nodeUUID          string
	nodeName          string
	address           string
	port              int
	apiSchema         string
	apiPath           string
	grpcAuthToken     string
	subpageConfigUUID string
	ctx               context.Context
	cancel            context.CancelFunc
	conn              *grpc.ClientConn
	client            proto.NodeServiceClient
	stream            proto.NodeService_StreamNodeDataClient
	streamCancel      context.CancelFunc
	streamGeneration  uint64
	lastResponseAt    time.Time
	sendMu            sync.Mutex

	isConnected  bool
	isConnecting bool
	lastError    string
	mutex        sync.RWMutex
}

type dbSubNode struct {
	UUID              string
	Name              string
	Address           string
	Port              int
	APISchema         string
	APIPath           string
	GRPCAuthToken     string
	SubpageConfigUUID string
}

type SubNodeRuntimeSnapshot struct {
	SingboxVersion *string
	NodeVersion    *string
	SingboxUptime  string
	CPUCount       *int
	CPUModel       *string
	TotalRAM       *string
}

func NewSubNodeMonitor(manager *dbmanager.DatabaseManager, cfg *config.BackendConfig) *SubNodeMonitor {
	return &SubNodeMonitor{
		manager:           manager,
		cfg:               cfg,
		nodes:             make(map[string]*subNodeState),
		runtimeByNodeName: make(map[string]SubNodeRuntimeSnapshot),
		syncNow:           make(chan struct{}, 1),
		deployNow:         make(chan []string, 1),
		subpagePushNow:    make(chan subpageConfigPushCommand, 1),
		srsSyncNow:        make(chan []string, 1),
	}
}

func (sm *SubNodeMonitor) Start(ctx context.Context, wg *sync.WaitGroup) {
	defer wg.Done()

	sm.globalCtx = ctx
	sm.globalCancel = func() {
		sm.nodesLock.RLock()
		defer sm.nodesLock.RUnlock()
		for _, state := range sm.nodes {
			state.cancel()
		}
	}

	sm.cfg.Logger.Trace("Subscription monitor initial sync")
	sm.syncNodes()

	syncTicker := time.NewTicker(30 * time.Second)
	defer syncTicker.Stop()

	for {
		select {
		case <-ctx.Done():
			sm.cfg.Logger.Info("Subscription monitor stopping")
			sm.stopAll()
			return
		case <-sm.syncNow:
			sm.cfg.Logger.Debug("Subscription monitor manual sync requested")
			sm.syncNodes()
		case targets := <-sm.deployNow:
			sm.deployToConnectedNodes(targets)
		case push := <-sm.subpagePushNow:
			sm.pushSubpageConfigToConnectedNodes(push)
		case targets := <-sm.srsSyncNow:
			sm.cfg.Logger.Info("Subscription node SRS sync requested", "node_targets", len(targets))
			sm.syncSRSListsToConnectedNodes(targets)
		case <-syncTicker.C:
			sm.syncNodes()
		}
	}
}

func (sm *SubNodeMonitor) syncNodes() {
	dbNodes, err := sm.loadActiveNodes()
	if err != nil {
		sm.cfg.Logger.Warn("Failed to load subscription nodes from DB", "error", err)
		return
	}

	desired := make(map[string]dbSubNode, len(dbNodes))
	for _, n := range dbNodes {
		desired[n.Name] = n
	}

	sm.nodesLock.Lock()
	defer sm.nodesLock.Unlock()

	toStart := make(map[string]dbSubNode)
	for name, state := range sm.nodes {
		desiredNode, exists := desired[name]
		if !exists {
			sm.cfg.Logger.Debug("Subscription node removed from DB, stopping monitor", "node", name)
			state.cancel()
			if state.conn != nil {
				_ = state.conn.Close()
			}
			delete(sm.nodes, name)
			sm.deleteRuntimeSnapshot(name)
			continue
		}
		if subNodeConfigChanged(state, desiredNode) {
			sm.cfg.Logger.Info(
				"Subscription node config changed, restarting monitor",
				"node", name,
				"old_address", state.address,
				"new_address", desiredNode.Address,
				"old_port", state.port,
				"new_port", desiredNode.Port,
				"old_schema", state.apiSchema,
				"new_schema", desiredNode.APISchema,
				"old_path", state.apiPath,
				"new_path", desiredNode.APIPath,
				"old_subpage_config_uuid", state.subpageConfigUUID,
				"new_subpage_config_uuid", desiredNode.SubpageConfigUUID,
			)
			state.cancel()
			if state.conn != nil {
				_ = state.conn.Close()
			}
			delete(sm.nodes, name)
			toStart[name] = desiredNode
		}
	}

	for name, node := range desired {
		if _, exists := sm.nodes[name]; !exists {
			toStart[name] = node
		}
	}

	for _, node := range toStart {
		sm.startNode(node)
	}
}

func (sm *SubNodeMonitor) loadActiveNodes() ([]dbSubNode, error) {
	nodes := make([]dbSubNode, 0)
	globalGRPCToken, err := sm.loadControlPlaneGRPCToken()
	if err != nil {
		return nil, err
	}
	err = sm.manager.ExecuteHighPriority(func(db dbmanager.DBExecutor) error {
		rows, err := db.Query(`
			SELECT n.uuid, n.name, n.address, n.port, n.api_schema, n.api_path, n.grpc_auth_token,
			       sns.subpage_config_uuid
			FROM sub_nodes n
			LEFT JOIN sub_nodes_to_subscription_page_config sns ON sns.node_uuid = n.uuid
			WHERE n.is_disabled = false
			ORDER BY n.view_position ASC, n.name ASC
		`)
		if err != nil {
			return fmt.Errorf("query sub_nodes: %w", err)
		}
		defer rows.Close()

		for rows.Next() {
			var (
				n             dbSubNode
				port          sql.NullInt64
				grpcAuthToken sql.NullString
				subpageConfig sql.NullString
			)
			if err := rows.Scan(&n.UUID, &n.Name, &n.Address, &port, &n.APISchema, &n.APIPath, &grpcAuthToken, &subpageConfig); err != nil {
				return fmt.Errorf("scan sub_node: %w", err)
			}
			if port.Valid {
				n.Port = int(port.Int64)
			} else {
				n.Port = 2222
			}
			n.APISchema = normalizeSubSchema(n.APISchema)
			if n.APISchema == "tls" {
				n.GRPCAuthToken = strings.TrimSpace(grpcAuthToken.String)
				if n.GRPCAuthToken == "" {
					n.GRPCAuthToken = globalGRPCToken
				}
			} else {
				n.GRPCAuthToken = ""
			}
			if subpageConfig.Valid {
				n.SubpageConfigUUID = normalizeAssignedSubpageConfigUUID(subpageConfig.String)
			}
			nodes = append(nodes, n)
		}
		return rows.Err()
	})
	return nodes, err
}

func (sm *SubNodeMonitor) loadControlPlaneGRPCToken() (string, error) {
	var token sql.NullString
	err := sm.manager.ExecuteHighPriority(func(db dbmanager.DBExecutor) error {
		return db.QueryRow(`
			SELECT grpc_auth_token
			FROM keygen
			ORDER BY created_at ASC
			LIMIT 1
		`).Scan(&token)
	})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", nil
		}
		return "", fmt.Errorf("query keygen grpc_auth_token: %w", err)
	}
	return strings.TrimSpace(token.String), nil
}

func (sm *SubNodeMonitor) startNode(dbNode dbSubNode) {
	ctx, cancel := context.WithCancel(sm.globalCtx)
	state := &subNodeState{
		nodeUUID:          dbNode.UUID,
		nodeName:          dbNode.Name,
		address:           dbNode.Address,
		port:              dbNode.Port,
		apiSchema:         dbNode.APISchema,
		apiPath:           dbNode.APIPath,
		grpcAuthToken:     dbNode.GRPCAuthToken,
		subpageConfigUUID: dbNode.SubpageConfigUUID,
		ctx:               ctx,
		cancel:            cancel,
	}

	sm.nodes[dbNode.Name] = state
	sm.updateConnectionStatus(dbNode.Name, false, true, "Connecting...")
	go sm.monitorNode(state)

	sm.cfg.Logger.Info(
		"Started monitoring subscription node",
		"node", dbNode.Name,
		"address", dbNode.Address,
		"port", dbNode.Port,
		"schema", dbNode.APISchema,
		"path", dbNode.APIPath,
		"subpage_config_uuid", dbNode.SubpageConfigUUID,
	)
}

func (sm *SubNodeMonitor) monitorNode(state *subNodeState) {
	const (
		minBackoff = 2 * time.Second
		maxBackoff = 60 * time.Second
	)
	backoff := minBackoff

	for {
		if state.ctx.Err() != nil {
			sm.cfg.Logger.Debug("Subscription node monitor stopped", "node", state.nodeName)
			return
		}

		sm.connectAndStream(state)
		if state.ctx.Err() != nil {
			sm.cfg.Logger.Debug("Subscription node monitor stopped", "node", state.nodeName)
			return
		}

		wait := subWithJitter(backoff, 0.2)
		timer := time.NewTimer(wait)
		select {
		case <-state.ctx.Done():
			timer.Stop()
			return
		case <-timer.C:
		}

		if backoff < maxBackoff {
			backoff = subMinDuration(maxBackoff, backoff*2)
		}
	}
}

func (sm *SubNodeMonitor) connectAndStream(state *subNodeState) {
	state.mutex.Lock()
	if state.isConnecting {
		state.mutex.Unlock()
		return
	}
	state.isConnecting = true
	state.mutex.Unlock()

	urlTarget := fmt.Sprintf("%s:%d", state.address, state.port)
	opts := make([]grpc.DialOption, 0, 3)

	apiSchema := strings.ToLower(strings.TrimSpace(state.apiSchema))
	switch apiSchema {
	case "tls":
		if strings.TrimSpace(state.grpcAuthToken) == "" {
			sm.cfg.Logger.Warn("Failed to connect to subscription node: missing global gRPC token", "node", state.nodeName)
			sm.updateConnectionStatus(state.nodeName, false, false, "Missing gRPC token in keygen")
			state.mutex.Lock()
			state.isConnecting = false
			state.mutex.Unlock()
			return
		}
		tlsCfg := &tls.Config{
			MinVersion: tls.VersionTLS12,
		}
		if sm.cfg != nil && sm.cfg.Panel.AllowInsecureHTTP {
			tlsCfg.InsecureSkipVerify = true
			sm.cfg.Logger.Warn("Subscription TLS verification is disabled by EXODUS_ALLOW_INSECURE_HTTP", "node", state.nodeName)
		}
		opts = append(opts, grpc.WithTransportCredentials(credentials.NewTLS(tlsCfg)))
	default:
		tlsCfg, tlsErr := sm.loadMTLSConfig(state.ctx)
		if tlsErr != nil {
			sm.cfg.Logger.Warn("Failed to build mTLS config for subscription node", "node", state.nodeName, "error", tlsErr)
			sm.updateConnectionStatus(state.nodeName, false, false, fmt.Sprintf("mTLS config failed: %v", tlsErr))
			state.mutex.Lock()
			state.isConnecting = false
			state.mutex.Unlock()
			return
		}
		opts = append(opts, grpc.WithTransportCredentials(credentials.NewTLS(tlsCfg)))
	}

	cleanPath := normalizeSubPath(state.apiPath)
	opts = append(opts, grpc.WithUnaryInterceptor(subPathPrefixUnaryInterceptor(cleanPath, state.grpcAuthToken)))
	opts = append(opts, grpc.WithStreamInterceptor(subPathPrefixStreamInterceptor(cleanPath, state.grpcAuthToken)))

	conn, err := grpc.NewClient(urlTarget, opts...)
	if err != nil {
		sm.cfg.Logger.Warn("Failed to connect to subscription node", "node", state.nodeName, "address", urlTarget, "error", err)
		sm.updateConnectionStatus(state.nodeName, false, false, fmt.Sprintf("Connection failed: %v", err))
		state.mutex.Lock()
		state.isConnecting = false
		state.mutex.Unlock()
		return
	}

	streamCtx, streamCancel := context.WithCancel(state.ctx)

	client := proto.NewNodeServiceClient(conn)
	stream, err := client.StreamNodeData(streamCtx)
	if err != nil {
		streamCancel()
		sm.cfg.Logger.Warn("Failed to create subscription stream", "node", state.nodeName, "error", err)
		_ = conn.Close()
		sm.updateConnectionStatus(state.nodeName, false, false, fmt.Sprintf("Stream failed: %v", err))
		state.mutex.Lock()
		state.isConnecting = false
		state.mutex.Unlock()
		return
	}

	state.mutex.Lock()
	state.conn = conn
	state.client = client
	state.stream = stream
	state.streamCancel = streamCancel
	state.streamGeneration++
	generation := state.streamGeneration
	state.lastResponseAt = time.Now()
	state.isConnected = true
	state.isConnecting = false
	state.lastError = ""
	state.mutex.Unlock()

	if err := sm.sendNodeRequest(state, &proto.NodeDataRequest{
		Request: &proto.NodeDataRequest_Config{
			Config: &proto.StreamConfig{IntervalSeconds: int32(subNodeStreamInterval / time.Second)},
		},
	}); err != nil {
		sm.cfg.Logger.Warn("Failed to send subscription stream config", "node", state.nodeName, "error", err)
		sm.handleDisconnect(state, fmt.Sprintf("Config failed: %v", err))
		return
	}

	go sm.watchStreamHeartbeat(state, generation, subNodeStreamIdleTimeout, subNodeStreamWatchInterval)

	sm.pushAssignedSubpageConfig(state)

	sm.updateConnectionStatus(state.nodeName, true, false, "Connected")
	sm.cfg.Logger.Info("Subscription node connected", "node", state.nodeName)

	sm.receiveStream(state)
}

func (sm *SubNodeMonitor) receiveStream(state *subNodeState) {
	for {
		state.mutex.RLock()
		stream := state.stream
		state.mutex.RUnlock()
		if stream == nil {
			sm.handleDisconnect(state, "Stream unavailable")
			return
		}

		resp, err := stream.Recv()
		if err == io.EOF {
			sm.handleDisconnect(state, "Stream closed")
			return
		}
		if err != nil {
			if st, ok := status.FromError(err); ok {
				if st.Code() == codes.Canceled && state.ctx.Err() != nil {
					return
				}
				if st.Code() == codes.Unavailable {
					sm.handleDisconnect(state, "Node unavailable")
					return
				}
			}
			sm.handleDisconnect(state, fmt.Sprintf("Stream error: %v", err))
			return
		}

		sm.markStreamActivity(state)

		if err := sm.processResponse(state, resp); err != nil {
			sm.handleDisconnect(state, err.Error())
			return
		}
	}
}

func (sm *SubNodeMonitor) processResponse(state *subNodeState, resp *proto.NodeDataResponse) error {
	if resp == nil {
		return nil
	}

	switch payload := resp.Response.(type) {
	case *proto.NodeDataResponse_Stats:
		sm.updateRuntimeFromStats(state.nodeName, payload.Stats.GetStats())
		return nil
	case *proto.NodeDataResponse_SubscriptionRequest:
		response := sm.handleSubscriptionBridgeRequest(state, payload.SubscriptionRequest)
		return sm.sendNodeRequest(state, &proto.NodeDataRequest{
			Request: &proto.NodeDataRequest_SubscriptionResponse{SubscriptionResponse: response},
		})
	default:
		return nil
	}
}

func (sm *SubNodeMonitor) handleSubscriptionBridgeRequest(state *subNodeState, req *proto.SubscriptionBridgeRequest) *proto.SubscriptionBridgeResponse {
	if req == nil {
		return &proto.SubscriptionBridgeResponse{
			StatusCode: http.StatusBadRequest,
			Error:      "empty bridge request",
		}
	}

	requestID := strings.TrimSpace(req.GetRequestId())
	if requestID == "" {
		requestID = fmt.Sprintf("panel-%d", time.Now().UnixNano())
	}

	ctx, cancel := context.WithTimeout(state.ctx, 20*time.Second)
	defer cancel()

	op := strings.TrimSpace(req.GetOperation())
	switch op {
	case bridgeOperationSubscriptionInfo:
		shortUUID := strings.TrimSpace(req.GetShortUuid())
		if shortUUID == "" {
			return &proto.SubscriptionBridgeResponse{RequestId: requestID, StatusCode: http.StatusBadRequest, Error: "short_uuid is required"}
		}
		path := "/api/sub/" + url.PathEscape(shortUUID) + "/info"
		statusCode, payload, headers := sm.performInternalAPIRequest(ctx, http.MethodGet, path, nil, sm.bridgeHeadersWithClientIP(req))
		return buildBridgeResponse(requestID, statusCode, payload, headers)

	case bridgeOperationSubscriptionContent:
		shortUUID := strings.TrimSpace(req.GetShortUuid())
		if shortUUID == "" {
			return &proto.SubscriptionBridgeResponse{RequestId: requestID, StatusCode: http.StatusBadRequest, Error: "short_uuid is required"}
		}
		path := "/api/sub/" + url.PathEscape(shortUUID)
		clientType := strings.TrimSpace(req.GetClientType())
		if clientType != "" {
			path += "/" + url.PathEscape(clientType)
		}
		statusCode, payload, headers := sm.performInternalAPIRequest(ctx, http.MethodGet, path, nil, sm.bridgeHeadersWithClientIP(req))
		return buildBridgeResponse(requestID, statusCode, payload, headers)

	case bridgeOperationSubpageByShortUUID:
		shortUUID := strings.TrimSpace(req.GetShortUuid())
		if shortUUID == "" {
			return &proto.SubscriptionBridgeResponse{RequestId: requestID, StatusCode: http.StatusBadRequest, Error: "short_uuid is required"}
		}

		payload, marshalErr := json.Marshal(map[string]any{
			"requestHeaders": sm.flattenRequestHeaders(req.GetHeaders()),
		})
		if marshalErr != nil {
			return &proto.SubscriptionBridgeResponse{RequestId: requestID, StatusCode: http.StatusInternalServerError, Error: marshalErr.Error()}
		}

		path := "/api/subscriptions/subpage-config/" + url.PathEscape(shortUUID)
		statusCode, body, headers := sm.performInternalAPIRequest(ctx, http.MethodPost, path, payload, nil)
		return buildBridgeResponse(requestID, statusCode, body, headers)

	case bridgeOperationSubpageByUUID:
		uuidValue := strings.TrimSpace(req.GetSubpageConfigUuid())
		if uuidValue == "" {
			return &proto.SubscriptionBridgeResponse{RequestId: requestID, StatusCode: http.StatusBadRequest, Error: "subpage_config_uuid is required"}
		}
		configPayload, err := sm.fetchSubpageConfigRaw(ctx, uuidValue)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return &proto.SubscriptionBridgeResponse{RequestId: requestID, StatusCode: http.StatusNotFound, Error: "subpage config not found"}
			}
			return &proto.SubscriptionBridgeResponse{RequestId: requestID, StatusCode: http.StatusInternalServerError, Error: err.Error()}
		}
		return &proto.SubscriptionBridgeResponse{
			RequestId:  requestID,
			StatusCode: http.StatusOK,
			Payload:    configPayload,
		}

	default:
		return &proto.SubscriptionBridgeResponse{
			RequestId:  requestID,
			StatusCode: http.StatusBadRequest,
			Error:      "unknown bridge operation",
		}
	}
}

func (sm *SubNodeMonitor) bridgeHeadersWithClientIP(req *proto.SubscriptionBridgeRequest) []*proto.Header {
	headers := append([]*proto.Header{}, req.GetHeaders()...)
	clientIP := strings.TrimSpace(req.GetClientIp())
	if clientIP != "" {
		headers = append(headers, &proto.Header{Key: exodusRealIPHeader, Value: clientIP})
	}
	return headers
}

func (sm *SubNodeMonitor) flattenRequestHeaders(headers []*proto.Header) map[string]string {
	result := make(map[string]string)
	for _, header := range headers {
		if header == nil {
			continue
		}
		key := strings.TrimSpace(header.GetKey())
		if key == "" {
			continue
		}
		if _, exists := result[key]; exists {
			continue
		}
		result[key] = header.GetValue()
	}
	return result
}

func (sm *SubNodeMonitor) performInternalAPIRequest(
	ctx context.Context,
	method,
	path string,
	body []byte,
	headers []*proto.Header,
) (int, []byte, http.Header) {
	handler := sm.resolveInternalHandler(path)
	if handler == nil {
		return http.StatusNotFound, []byte(`{"error":"not found"}`), http.Header{}
	}

	request, err := http.NewRequestWithContext(ctx, method, path, bytes.NewReader(body))
	if err != nil {
		return http.StatusInternalServerError, []byte(err.Error()), http.Header{}
	}
	if len(body) > 0 {
		request.Header.Set("Content-Type", "application/json")
	}
	for _, header := range headers {
		if header == nil {
			continue
		}
		key := strings.TrimSpace(header.GetKey())
		if key == "" {
			continue
		}
		request.Header.Add(key, header.GetValue())
	}

	recorder := httptest.NewRecorder()
	handler.ServeHTTP(recorder, request)

	result := recorder.Result()
	defer result.Body.Close()

	responseBody, readErr := io.ReadAll(result.Body)
	if readErr != nil {
		return http.StatusInternalServerError, []byte(readErr.Error()), result.Header.Clone()
	}

	return result.StatusCode, responseBody, result.Header.Clone()
}

func (sm *SubNodeMonitor) resolveInternalHandler(path string) http.HandlerFunc {
	switch {
	case path == "/api/system/metadata":
		return systemapi.MetadataHandler(sm.cfg)
	case strings.HasPrefix(path, "/api/sub/") || path == "/api/sub":
		return subscriptionapi.SubscriptionPublicHandler(sm.manager, sm.cfg)
	case strings.HasPrefix(path, "/api/subscriptions/") || path == "/api/subscriptions":
		return subscriptionapi.SubscriptionsByPathHandler(sm.manager, sm.cfg)
	default:
		return nil
	}
}

func (sm *SubNodeMonitor) fetchSubpageConfigRaw(ctx context.Context, uuidValue string) ([]byte, error) {
	var payload string
	err := sm.manager.ExecuteHighPriority(func(db dbmanager.DBExecutor) error {
		return db.QueryRowContext(ctx, `
			SELECT config
			FROM subscription_page_config
			WHERE uuid = ?
			LIMIT 1
		`, uuidValue).Scan(&payload)
	})
	if err != nil {
		return nil, err
	}

	raw := []byte(strings.TrimSpace(payload))
	if len(raw) == 0 || !json.Valid(raw) {
		return nil, fmt.Errorf("invalid subpage config payload")
	}

	return raw, nil
}

func buildBridgeResponse(requestID string, statusCode int, body []byte, headers http.Header) *proto.SubscriptionBridgeResponse {
	response := &proto.SubscriptionBridgeResponse{
		RequestId:  requestID,
		StatusCode: int32(statusCode),
		Payload:    body,
		Headers:    toProtoHeaders(headers),
	}
	if statusCode < 200 || statusCode >= 300 {
		trimmed := strings.TrimSpace(string(body))
		if trimmed == "" {
			trimmed = http.StatusText(statusCode)
		}
		response.Error = trimmed
	}
	return response
}

func toProtoHeaders(headers http.Header) []*proto.Header {
	if len(headers) == 0 {
		return nil
	}
	result := make([]*proto.Header, 0, len(headers))
	for key, values := range headers {
		for _, value := range values {
			result = append(result, &proto.Header{Key: key, Value: value})
		}
	}
	return result
}

func (sm *SubNodeMonitor) sendNodeRequest(state *subNodeState, req *proto.NodeDataRequest) error {
	if req == nil {
		return fmt.Errorf("empty node request")
	}

	state.sendMu.Lock()
	defer state.sendMu.Unlock()

	state.mutex.RLock()
	stream := state.stream
	connected := state.isConnected
	state.mutex.RUnlock()
	if !connected || stream == nil {
		return fmt.Errorf("stream is not connected")
	}

	return stream.Send(req)
}

func (sm *SubNodeMonitor) watchStreamHeartbeat(
	state *subNodeState,
	generation uint64,
	idleTimeout time.Duration,
	pollInterval time.Duration,
) {
	ticker := time.NewTicker(pollInterval)
	defer ticker.Stop()

	for {
		select {
		case <-state.ctx.Done():
			return
		case <-ticker.C:
		}

		state.mutex.RLock()
		currentGeneration := state.streamGeneration
		lastResponseAt := state.lastResponseAt
		streamCancel := state.streamCancel
		nodeName := state.nodeName
		isConnected := state.isConnected
		state.mutex.RUnlock()

		if currentGeneration != generation || !isConnected || streamCancel == nil {
			return
		}
		if lastResponseAt.IsZero() || time.Since(lastResponseAt) <= idleTimeout {
			continue
		}

		if sm.cfg != nil {
			sm.cfg.Logger.Warn(
				"Subscription node stream heartbeat timed out",
				"node", nodeName,
				"idle_for", time.Since(lastResponseAt).Round(time.Second),
			)
		}
		streamCancel()
		return
	}
}

func (sm *SubNodeMonitor) markStreamActivity(state *subNodeState) {
	state.mutex.Lock()
	if state.stream != nil {
		state.lastResponseAt = time.Now()
	}
	state.mutex.Unlock()
}

func (sm *SubNodeMonitor) updateRuntimeFromStats(nodeName string, stats []*proto.Stat) {
	cleanNodeName := strings.TrimSpace(nodeName)
	if cleanNodeName == "" || len(stats) == 0 {
		return
	}

	values := make(map[string]string, len(stats))
	for _, stat := range stats {
		if stat == nil {
			continue
		}
		key := strings.ToLower(strings.TrimSpace(stat.GetName()))
		if key == "" {
			continue
		}
		values[key] = strings.TrimSpace(stat.GetValue())
	}
	if len(values) == 0 {
		return
	}

	sm.runtimeMu.Lock()
	runtime := sm.runtimeByNodeName[cleanNodeName]
	if runtime.SingboxUptime == "" {
		runtime.SingboxUptime = "0"
	}

	if version, ok := normalizeSubNodeRuntimeVersion(firstNonEmptyString(
		values[subNodeRuntimeStatVersion],
		values["node_version"],
	)); ok {
		runtime.NodeVersion = stringPtr(version)
	}

	if singboxVersion, ok := normalizeSubNodeRuntimeVersion(values["singbox_version"]); ok {
		runtime.SingboxVersion = stringPtr(singboxVersion)
	}

	if uptime, ok := normalizeSubNodeRuntimeUptime(firstNonEmptyString(
		values[subNodeRuntimeStatUptime],
		values["singbox_uptime"],
	)); ok {
		runtime.SingboxUptime = uptime
	}

	if cpuCount, ok := parseOptionalIntValue(values[subNodeRuntimeStatCPUCount]); ok {
		runtime.CPUCount = intPtr(cpuCount)
	}

	if cpuModel, ok := parseOptionalStringValue(values[subNodeRuntimeStatCPUModel]); ok {
		runtime.CPUModel = stringPtr(cpuModel)
	}

	if totalRAM, ok := parseOptionalStringValue(values[subNodeRuntimeStatTotalRAM]); ok {
		runtime.TotalRAM = stringPtr(totalRAM)
	}

	sm.runtimeByNodeName[cleanNodeName] = runtime
	sm.runtimeMu.Unlock()
}

func (sm *SubNodeMonitor) handleDisconnect(state *subNodeState, reason string) {
	state.mutex.Lock()
	wasConnected := state.isConnected
	streamCancel := state.streamCancel
	conn := state.conn
	state.isConnected = false
	state.isConnecting = false
	state.lastError = reason
	state.lastResponseAt = time.Time{}
	state.stream = nil
	state.streamCancel = nil
	state.conn = nil
	state.mutex.Unlock()

	if streamCancel != nil {
		streamCancel()
	}
	if conn != nil {
		_ = conn.Close()
	}

	if wasConnected {
		sm.updateConnectionStatus(state.nodeName, false, false, reason)
		sm.cfg.Logger.Info("Subscription node disconnected", "node", state.nodeName, "reason", reason)
	}
}

func (sm *SubNodeMonitor) updateConnectionStatus(nodeName string, isConnected, isConnecting bool, message string) {
	err := sm.manager.ExecuteHighPriority(func(db dbmanager.DBExecutor) error {
		var (
			currentConnected  bool
			currentConnecting bool
			currentMessage    sql.NullString
		)
		err := db.QueryRow(`SELECT is_connected, is_connecting, last_status_message FROM sub_nodes WHERE name = ?`, nodeName).
			Scan(&currentConnected, &currentConnecting, &currentMessage)
		if err != nil {
			if err == sql.ErrNoRows {
				return nil
			}
			return fmt.Errorf("query sub node status: %w", err)
		}

		msgStr := ""
		if currentMessage.Valid {
			msgStr = currentMessage.String
		}
		if currentConnected == isConnected && currentConnecting == isConnecting && msgStr == message {
			return nil
		}

		_, err = db.Exec(`
			UPDATE sub_nodes
			SET is_connected = ?,
			    is_connecting = ?,
			    last_status_message = ?,
			    last_status_change = CURRENT_TIMESTAMP,
			    updated_at = CURRENT_TIMESTAMP
			WHERE name = ?
		`, isConnected, isConnecting, message, nodeName)
		if err != nil {
			return fmt.Errorf("update sub node status: %w", err)
		}
		return nil
	})
	if err != nil {
		sm.cfg.Logger.Warn("Failed to update subscription node status", "node", nodeName, "error", err)
	}
}

func (sm *SubNodeMonitor) syncSRSListsToConnectedNodes(requestedNodeUUIDs []string) {
	if sm == nil {
		return
	}

	srsLists, err := srscore.LoadNodeSyncItems(context.Background(), sm.manager)
	if err != nil {
		sm.cfg.Logger.Warn("Failed to load SRS lists for subscription node sync", "error", err)
		return
	}

	payload, err := json.Marshal(map[string]any{
		"srs_lists": srsLists,
	})
	if err != nil {
		sm.cfg.Logger.Warn("Failed to marshal SRS sync payload for subscription nodes", "error", err)
		return
	}

	targetFilter := make(map[string]struct{}, len(requestedNodeUUIDs))
	for _, item := range requestedNodeUUIDs {
		uuid := strings.TrimSpace(item)
		if uuid == "" {
			continue
		}
		targetFilter[uuid] = struct{}{}
	}

	sm.nodesLock.RLock()
	states := make([]*subNodeState, 0, len(sm.nodes))
	for _, state := range sm.nodes {
		if state != nil {
			states = append(states, state)
		}
	}
	sm.nodesLock.RUnlock()

	matchedTargets := 0
	readyTargets := 0
	sentTargets := 0

	for _, state := range states {
		state.mutex.RLock()
		nodeUUID := state.nodeUUID
		nodeName := state.nodeName
		ready := state.isConnected && state.client != nil
		client := state.client
		state.mutex.RUnlock()
		if len(targetFilter) > 0 {
			if _, ok := targetFilter[nodeUUID]; !ok {
				continue
			}
		}
		matchedTargets++
		if !ready {
			continue
		}
		readyTargets++

		ctxBase := sm.globalCtx
		if ctxBase == nil {
			ctxBase = context.Background()
		}
		ctx, cancel := context.WithTimeout(ctxBase, 30*time.Second)
		resp, submitErr := client.SubmitTask(ctx, &proto.NodeTask{
			TaskId:    fmt.Sprintf("sync-srs-%d", time.Now().UnixNano()),
			Operation: "sync_srs_lists",
			Payload:   payload,
		})
		cancel()

		if submitErr != nil {
			sm.cfg.Logger.Warn("SRS sync task failed on subscription node", "node", nodeName, "error", submitErr)
			continue
		}
		if resp == nil || resp.Code != int32(codes.OK) {
			if resp == nil {
				sm.cfg.Logger.Warn("SRS sync returned nil status from subscription node", "node", nodeName)
			} else {
				sm.cfg.Logger.Warn("SRS sync rejected by subscription node", "node", nodeName, "code", resp.Code, "message", resp.Message)
			}
			continue
		}

		sentTargets++
		sm.cfg.Logger.Info("SRS lists synced to subscription node", "node", nodeName, "lists", len(srsLists), "message", resp.Message)
	}

	sm.cfg.Logger.Debug(
		"SRS subscription node sync processed",
		"target_filter_count", len(targetFilter),
		"matched_targets", matchedTargets,
		"ready_targets", readyTargets,
		"sent_targets", sentTargets,
		"lists", len(srsLists),
	)
}

func (sm *SubNodeMonitor) deployToConnectedNodes(requestedNodeUUIDs []string) {
	targetFilter := make(map[string]struct{}, len(requestedNodeUUIDs))
	for _, item := range requestedNodeUUIDs {
		uuid := strings.TrimSpace(item)
		if uuid == "" {
			continue
		}
		targetFilter[uuid] = struct{}{}
	}

	sm.nodesLock.RLock()
	states := make([]*subNodeState, 0, len(sm.nodes))
	for _, state := range sm.nodes {
		if state != nil {
			states = append(states, state)
		}
	}
	sm.nodesLock.RUnlock()

	for _, state := range states {
		state.mutex.RLock()
		nodeUUID := state.nodeUUID
		nodeName := state.nodeName
		ready := state.isConnected && state.stream != nil
		state.mutex.RUnlock()
		if !ready {
			continue
		}
		if len(targetFilter) > 0 {
			if _, ok := targetFilter[nodeUUID]; !ok {
				continue
			}
		}

		err := sm.sendNodeRequest(state, &proto.NodeDataRequest{
			Request: &proto.NodeDataRequest_Config{Config: &proto.StreamConfig{IntervalSeconds: 20}},
		})
		if err != nil {
			sm.cfg.Logger.Warn("Failed to push subscription config over stream", "node", nodeName, "error", err)
			sm.handleDisconnect(state, fmt.Sprintf("Config push failed: %v", err))
			continue
		}
		sm.pushAssignedSubpageConfig(state)
		sm.cfg.Logger.Info("Subscription config push sent", "node", nodeName)
	}
}

func (sm *SubNodeMonitor) pushAssignedSubpageConfig(state *subNodeState) {
	if state == nil {
		return
	}

	state.mutex.RLock()
	nodeName := state.nodeName
	nodeUUID := state.nodeUUID
	subpageConfigUUID := strings.TrimSpace(state.subpageConfigUUID)
	ready := state.isConnected && state.stream != nil
	state.mutex.RUnlock()

	if !ready || subpageConfigUUID == "" {
		return
	}

	configPayload, err := sm.fetchSubpageConfigRaw(state.ctx, subpageConfigUUID)
	if err != nil {
		sm.cfg.Logger.Warn(
			"Failed to fetch assigned subpage config for subscription node",
			"node", nodeName,
			"node_uuid", nodeUUID,
			"subpage_config_uuid", subpageConfigUUID,
			"error", err,
		)
		return
	}

	err = sm.sendNodeRequest(state, &proto.NodeDataRequest{
		Request: &proto.NodeDataRequest_SubpageConfigUpdate{SubpageConfigUpdate: &proto.SubpageConfigUpdate{
			Uuid:   subpageConfigUUID,
			Config: configPayload,
		}},
	})
	if err != nil {
		sm.cfg.Logger.Warn(
			"Failed to push assigned subpage config to subscription node",
			"node", nodeName,
			"node_uuid", nodeUUID,
			"subpage_config_uuid", subpageConfigUUID,
			"error", err,
		)
		sm.handleDisconnect(state, fmt.Sprintf("Subpage config push failed: %v", err))
		return
	}

	sm.cfg.Logger.Info(
		"Assigned subpage config push sent",
		"node", nodeName,
		"node_uuid", nodeUUID,
		"subpage_config_uuid", subpageConfigUUID,
	)
}

func (sm *SubNodeMonitor) pushSubpageConfigToConnectedNodes(command subpageConfigPushCommand) {
	targetFilter := make(map[string]struct{}, len(command.targetUUIDs))
	for _, item := range command.targetUUIDs {
		uuid := strings.TrimSpace(item)
		if uuid == "" {
			continue
		}
		targetFilter[uuid] = struct{}{}
	}

	sm.nodesLock.RLock()
	states := make([]*subNodeState, 0, len(sm.nodes))
	for _, state := range sm.nodes {
		if state != nil {
			states = append(states, state)
		}
	}
	sm.nodesLock.RUnlock()

	matchedTargets := 0
	readyTargets := 0
	sentTargets := 0

	for _, state := range states {
		state.mutex.RLock()
		nodeUUID := state.nodeUUID
		nodeName := state.nodeName
		ready := state.isConnected && state.stream != nil
		state.mutex.RUnlock()
		if len(targetFilter) > 0 {
			if _, ok := targetFilter[nodeUUID]; !ok {
				continue
			}
		}
		matchedTargets++
		if !ready {
			continue
		}
		readyTargets++

		err := sm.sendNodeRequest(state, &proto.NodeDataRequest{
			Request: &proto.NodeDataRequest_SubpageConfigUpdate{SubpageConfigUpdate: &proto.SubpageConfigUpdate{
				Uuid:   command.uuid,
				Config: command.config,
			}},
		})
		if err != nil {
			sm.cfg.Logger.Warn("Failed to push subpage config update", "node", nodeName, "error", err)
			sm.handleDisconnect(state, fmt.Sprintf("Subpage config push failed: %v", err))
			continue
		}
		sentTargets++
		sm.cfg.Logger.Info("Subpage config push sent", "node", nodeName, "uuid", command.uuid)
	}

	sm.cfg.Logger.Debug(
		"Subpage config push processed",
		"uuid", command.uuid,
		"target_filter_count", len(targetFilter),
		"matched_targets", matchedTargets,
		"ready_targets", readyTargets,
		"sent_targets", sentTargets,
		"payload_bytes", len(command.config),
	)
}

func (sm *SubNodeMonitor) stopAll() {
	sm.nodesLock.Lock()
	defer sm.nodesLock.Unlock()

	for name, state := range sm.nodes {
		state.cancel()
		if state.conn != nil {
			_ = state.conn.Close()
		}
		delete(sm.nodes, name)
		sm.deleteRuntimeSnapshot(name)
	}
}

func (sm *SubNodeMonitor) RuntimeSnapshot(nodeName string) (SubNodeRuntimeSnapshot, bool) {
	if sm == nil {
		return SubNodeRuntimeSnapshot{}, false
	}

	cleanNodeName := strings.TrimSpace(nodeName)
	if cleanNodeName == "" {
		return SubNodeRuntimeSnapshot{}, false
	}

	sm.runtimeMu.RLock()
	snapshot, ok := sm.runtimeByNodeName[cleanNodeName]
	sm.runtimeMu.RUnlock()
	if !ok {
		return SubNodeRuntimeSnapshot{}, false
	}

	return cloneRuntimeSnapshot(snapshot), true
}

func (sm *SubNodeMonitor) deleteRuntimeSnapshot(nodeName string) {
	if sm == nil {
		return
	}

	cleanNodeName := strings.TrimSpace(nodeName)
	if cleanNodeName == "" {
		return
	}

	sm.runtimeMu.Lock()
	delete(sm.runtimeByNodeName, cleanNodeName)
	sm.runtimeMu.Unlock()
}

func (sm *SubNodeMonitor) Stop() {
	if sm != nil && sm.globalCancel != nil {
		sm.globalCancel()
	}
}

func (sm *SubNodeMonitor) RequestSync() {
	if sm == nil {
		return
	}
	select {
	case sm.syncNow <- struct{}{}:
	default:
	}
}

func (sm *SubNodeMonitor) RequestDeploy(nodeUUIDs ...string) {
	if sm == nil {
		return
	}
	normalized := normalizeSubNodeUUIDTargets(nodeUUIDs)
	select {
	case sm.deployNow <- normalized:
	default:
		select {
		case <-sm.deployNow:
		default:
		}
		sm.deployNow <- normalized
	}
}

func (sm *SubNodeMonitor) RequestSRSDeploy(nodeUUIDs ...string) {
	if sm == nil {
		return
	}
	normalized := normalizeSubNodeUUIDTargets(nodeUUIDs)
	select {
	case sm.srsSyncNow <- normalized:
	default:
		select {
		case <-sm.srsSyncNow:
		default:
		}
		sm.srsSyncNow <- normalized
	}
}

func (sm *SubNodeMonitor) RequestSubpageConfigPush(uuid string, config []byte, nodeUUIDs ...string) {
	if sm == nil {
		return
	}

	payload := make([]byte, len(config))
	copy(payload, config)
	command := subpageConfigPushCommand{
		uuid:        strings.TrimSpace(uuid),
		config:      payload,
		targetUUIDs: normalizeSubNodeUUIDTargets(nodeUUIDs),
	}
	sm.cfg.Logger.Debug(
		"Subpage config push queued",
		"uuid", command.uuid,
		"target_nodes_count", len(command.targetUUIDs),
		"payload_bytes", len(command.config),
	)

	select {
	case sm.subpagePushNow <- command:
	default:
		select {
		case <-sm.subpagePushNow:
		default:
		}
		sm.subpagePushNow <- command
	}
}

func (sm *SubNodeMonitor) loadMTLSConfig(ctx context.Context) (*tls.Config, error) {
	var (
		caCertPEM     string
		clientCertPEM string
		clientKeyPEM  string
	)

	err := sm.manager.ExecuteHighPriority(func(db dbmanager.DBExecutor) error {
		return db.QueryRowContext(ctx, `
			SELECT ca_cert, client_cert, client_key
			FROM keygen
			ORDER BY created_at ASC
			LIMIT 1
		`).Scan(&caCertPEM, &clientCertPEM, &clientKeyPEM)
	})
	if err != nil {
		return nil, fmt.Errorf("load keygen mTLS material: %w", err)
	}

	clientCert, err := tls.X509KeyPair([]byte(clientCertPEM), []byte(clientKeyPEM))
	if err != nil {
		return nil, fmt.Errorf("parse keygen client certificate/key: %w", err)
	}

	pool := x509.NewCertPool()
	if !pool.AppendCertsFromPEM([]byte(caCertPEM)) {
		return nil, fmt.Errorf("parse keygen CA certificate")
	}

	return &tls.Config{
		MinVersion:   tls.VersionTLS12,
		RootCAs:      pool,
		Certificates: []tls.Certificate{clientCert},
		ServerName:   "internal.exodus.local",
	}, nil
}

func subPathPrefixUnaryInterceptor(prefix, authToken string) grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req, reply any, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		if token := strings.TrimSpace(authToken); token != "" {
			ctx = metadata.AppendToOutgoingContext(ctx, subGRPCTokenHeader, token)
		}
		return invoker(ctx, prefix+method, req, reply, cc, opts...)
	}
}

func subPathPrefixStreamInterceptor(prefix, authToken string) grpc.StreamClientInterceptor {
	return func(ctx context.Context, desc *grpc.StreamDesc, cc *grpc.ClientConn, method string, streamer grpc.Streamer, opts ...grpc.CallOption) (grpc.ClientStream, error) {
		if token := strings.TrimSpace(authToken); token != "" {
			ctx = metadata.AppendToOutgoingContext(ctx, subGRPCTokenHeader, token)
		}
		return streamer(ctx, desc, cc, prefix+method, opts...)
	}
}

func normalizeSubPath(value string) string {
	trimmed := strings.TrimSpace(value)
	if trimmed == "" || trimmed == "/" {
		return ""
	}
	return "/" + strings.Trim(trimmed, "/")
}

func normalizeSubSchema(value string) string {
	normalized := strings.ToLower(strings.TrimSpace(value))
	switch normalized {
	case "tls":
		return "tls"
	case "mtls":
		return "mtls"
	default:
		return "mtls"
	}
}

func normalizeAssignedSubpageConfigUUID(value string) string {
	trimmed := strings.TrimSpace(value)
	if trimmed == "" || trimmed == "00000000-0000-0000-0000-000000000000" {
		return ""
	}
	return trimmed
}

func subNodeConfigChanged(state *subNodeState, desired dbSubNode) bool {
	if state.address != desired.Address {
		return true
	}
	if state.port != desired.Port {
		return true
	}
	if normalizeSubSchema(state.apiSchema) != normalizeSubSchema(desired.APISchema) {
		return true
	}
	if strings.Trim(normalizeSubPath(state.apiPath), "/") != strings.Trim(normalizeSubPath(desired.APIPath), "/") {
		return true
	}
	if strings.TrimSpace(state.grpcAuthToken) != strings.TrimSpace(desired.GRPCAuthToken) {
		return true
	}
	if normalizeAssignedSubpageConfigUUID(state.subpageConfigUUID) != normalizeAssignedSubpageConfigUUID(desired.SubpageConfigUUID) {
		return true
	}
	return false
}

func normalizeSubNodeUUIDTargets(raw []string) []string {
	if len(raw) == 0 {
		return nil
	}
	seen := make(map[string]struct{}, len(raw))
	result := make([]string, 0, len(raw))
	for _, item := range raw {
		uuid := strings.TrimSpace(item)
		if uuid == "" {
			continue
		}
		if _, exists := seen[uuid]; exists {
			continue
		}
		seen[uuid] = struct{}{}
		result = append(result, uuid)
	}
	return result
}

func subWithJitter(base time.Duration, factor float64) time.Duration {
	if base <= 0 || factor <= 0 {
		return base
	}
	delta := int64(float64(base) * factor)
	if delta <= 0 {
		return base
	}
	offset := rand.Int64N(2*delta+1) - delta
	result := base + time.Duration(offset)
	if result <= 0 {
		return base
	}
	return result
}

func subMinDuration(a, b time.Duration) time.Duration {
	if a <= b {
		return a
	}
	return b
}

func normalizeSubNodeRuntimeVersion(value string) (string, bool) {
	trimmed := strings.TrimSpace(value)
	if trimmed == "" {
		return "", false
	}

	lower := strings.ToLower(trimmed)
	switch lower {
	case "unknown", "latest", "(devel)", "sub":
		return "", false
	}

	if subNodeVersionPattern.MatchString(trimmed) {
		if strings.HasPrefix(lower, "v") {
			return "v" + strings.TrimSpace(trimmed[1:]), true
		}
		return trimmed, true
	}

	return trimmed, true
}

func normalizeSubNodeRuntimeUptime(value string) (string, bool) {
	trimmed := strings.TrimSpace(value)
	if trimmed == "" {
		return "", false
	}

	parsed, err := strconv.ParseInt(trimmed, 10, 64)
	if err != nil || parsed < 0 {
		return "", false
	}

	return strconv.FormatInt(parsed, 10), true
}

func parseOptionalIntValue(value string) (int, bool) {
	trimmed := strings.TrimSpace(value)
	if trimmed == "" {
		return 0, false
	}

	parsed, err := strconv.Atoi(trimmed)
	if err != nil {
		return 0, false
	}

	return parsed, true
}

func parseOptionalStringValue(value string) (string, bool) {
	trimmed := strings.TrimSpace(value)
	if trimmed == "" {
		return "", false
	}
	return trimmed, true
}

func firstNonEmptyString(values ...string) string {
	for _, value := range values {
		trimmed := strings.TrimSpace(value)
		if trimmed != "" {
			return trimmed
		}
	}
	return ""
}

func cloneRuntimeSnapshot(snapshot SubNodeRuntimeSnapshot) SubNodeRuntimeSnapshot {
	cloned := SubNodeRuntimeSnapshot{
		SingboxUptime: snapshot.SingboxUptime,
	}

	if snapshot.SingboxVersion != nil {
		cloned.SingboxVersion = stringPtr(*snapshot.SingboxVersion)
	}
	if snapshot.NodeVersion != nil {
		cloned.NodeVersion = stringPtr(*snapshot.NodeVersion)
	}
	if snapshot.CPUCount != nil {
		cloned.CPUCount = intPtr(*snapshot.CPUCount)
	}
	if snapshot.CPUModel != nil {
		cloned.CPUModel = stringPtr(*snapshot.CPUModel)
	}
	if snapshot.TotalRAM != nil {
		cloned.TotalRAM = stringPtr(*snapshot.TotalRAM)
	}
	if cloned.SingboxUptime == "" {
		cloned.SingboxUptime = "0"
	}

	return cloned
}

func stringPtr(value string) *string {
	v := value
	return &v
}

func intPtr(value int) *int {
	v := value
	return &v
}
