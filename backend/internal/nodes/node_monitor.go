package users

import (
	"context"
	"crypto/tls"
	"database/sql"
	"errors"
	"fmt"
	"io"
	"math"
	"math/rand/v2"
	"strconv"
	"strings"
	"sync"
	"time"

	"exodus/internal/config"
	"exodus/internal/db"
	dbmanager "exodus/internal/db/manager"
	"exodus/internal/notifications"
	"exodus/internal/proto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

const nodeGRPCTokenHeader = "x-exodus-grpc-token"

// NodeMonitor dynamically manages node monitoring with status tracking.
type NodeMonitor struct {
	manager *dbmanager.DatabaseManager
	cfg     *config.BackendConfig

	// Active node contexts
	nodes     map[string]*nodeState
	nodesLock sync.RWMutex

	// Global context for shutdown
	globalCtx    context.Context
	globalCancel context.CancelFunc

	// Manual sync trigger
	syncNow chan struct{}
	// Manual deploy trigger
	deployNow chan deployRequest
	// Manual SRS sync trigger
	srsSyncNow chan struct{}

	// Runtime traffic metrics snapshots by node UUID.
	metricsByNodeUUID map[string]*NodeMetricsSnapshot
	metricsLock       sync.RWMutex

	usageRecorder NodeUserUsageRecorder
}

type deployRequest struct {
	Restart   bool
	NodeUUIDs []string
}

type TagTrafficCounters struct {
	UploadBytes   int64
	DownloadBytes int64
}

type userUsageDelta struct {
	UserID       int64
	Username     string
	TotalBytes   int64
	HistoryBytes int64
}

type NodeUserUsageRecorder interface {
	RecordNodeUserUsage(ctx context.Context, nodeID int64, userBytes map[int64]int64) error
}

type NodeMetricsSnapshot struct {
	NodeUUID    string
	UsersOnline int
	Inbounds    map[string]TagTrafficCounters
	Outbounds   map[string]TagTrafficCounters
	UpdatedAt   time.Time
}

type nodeState struct {
	nodeUUID      string
	nodeName      string
	address       string
	port          int
	apiSchema     string
	apiPath       string
	grpcAuthToken string
	ctx           context.Context
	cancel        context.CancelFunc
	conn          *grpc.ClientConn
	client        proto.NodeServiceClient
	stream        proto.NodeService_StreamNodeDataClient
	isConnected   bool
	isConnecting  bool
	lastError     string
	mutex         sync.RWMutex
}

func pathPrefixUnaryInterceptor(prefix, authToken string) grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req, reply any, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		if token := strings.TrimSpace(authToken); token != "" {
			ctx = metadata.AppendToOutgoingContext(ctx, nodeGRPCTokenHeader, token)
		}
		newMethod := prefix + method
		return invoker(ctx, newMethod, req, reply, cc, opts...)
	}
}

func pathPrefixStreamInterceptor(prefix, authToken string) grpc.StreamClientInterceptor {
	return func(ctx context.Context, desc *grpc.StreamDesc, cc *grpc.ClientConn, method string, streamer grpc.Streamer, opts ...grpc.CallOption) (grpc.ClientStream, error) {
		if token := strings.TrimSpace(authToken); token != "" {
			ctx = metadata.AppendToOutgoingContext(ctx, nodeGRPCTokenHeader, token)
		}
		newMethod := prefix + method
		return streamer(ctx, desc, cc, newMethod, opts...)
	}
}

// NewNodeMonitor creates a new NodeMonitor.
func NewNodeMonitor(manager *dbmanager.DatabaseManager, cfg *config.BackendConfig) *NodeMonitor {
	return &NodeMonitor{
		manager:           manager,
		cfg:               cfg,
		nodes:             make(map[string]*nodeState),
		syncNow:           make(chan struct{}, 1),
		deployNow:         make(chan deployRequest, 1),
		srsSyncNow:        make(chan struct{}, 1),
		metricsByNodeUUID: make(map[string]*NodeMetricsSnapshot),
	}
}

func (nm *NodeMonitor) SetNodeUserUsageRecorder(recorder NodeUserUsageRecorder) {
	nm.usageRecorder = recorder
}

// Start begins the node monitoring loop.
func (nm *NodeMonitor) Start(ctx context.Context, wg *sync.WaitGroup) {
	defer wg.Done()

	// Store global context
	nm.globalCtx = ctx

	// Initialize cancel function for Stop()
	nm.globalCancel = func() {
		// Cancel all node contexts
		nm.nodesLock.RLock()
		defer nm.nodesLock.RUnlock()
		for _, state := range nm.nodes {
			state.cancel()
		}
	}

	// Initial load and start
	nm.cfg.Logger.Trace("Node monitor initial sync")
	nm.syncNodes()

	// Periodic sync every 30 seconds
	syncTicker := time.NewTicker(30 * time.Second)
	defer syncTicker.Stop()

	for {
		select {
		case <-ctx.Done():
			nm.cfg.Logger.Info("Node monitor stopping")
			nm.stopAll()
			return
		case <-nm.syncNow:
			nm.cfg.Logger.Debug("Node monitor manual sync requested")
			nm.syncNodes()
		case deployReq := <-nm.deployNow:
			nm.cfg.Logger.Info("Node monitor deploy requested", "restart", deployReq.Restart, "node_targets", len(deployReq.NodeUUIDs))
			nm.deployToConnectedNodes(deployReq.Restart, deployReq.NodeUUIDs)
		case <-nm.srsSyncNow:
			nm.cfg.Logger.Info("Node monitor SRS sync requested")
			nm.syncSRSListsToConnectedNodes()
		case <-syncTicker.C:
			nm.syncNodes()
		}
	}
}

// syncNodes synchronizes monitored nodes with database.
func (nm *NodeMonitor) syncNodes() {
	// Load active nodes from DB
	dbNodes, err := nm.loadActiveNodes()
	if err != nil {
		nm.cfg.Logger.Warn("Failed to load nodes from DB", "error", err)
		return
	}

	// Build desired state
	desired := make(map[string]db.DBNode)
	for _, n := range dbNodes {
		desired[n.Name] = n
	}
	nm.cfg.Logger.Debug("Node monitor sync complete", "nodes", len(desired))

	nm.nodesLock.Lock()
	defer nm.nodesLock.Unlock()

	toStart := make(map[string]db.DBNode)

	// Stop removed or changed nodes
	for name, state := range nm.nodes {
		desiredNode, exists := desired[name]
		if !exists {
			nm.cfg.Logger.Debug("Node removed from DB, stopping monitor", "node", name)
			state.cancel()
			if state.conn != nil {
				state.conn.Close()
			}
			if state.nodeUUID != "" {
				nm.removeNodeMetrics(state.nodeUUID)
			}
			delete(nm.nodes, name)
			continue
		}
		if nodeConfigChanged(state, desiredNode) {
			nm.cfg.Logger.Info(
				"Node config changed, restarting monitor",
				"node", name,
				"old_address", state.address, "new_address", desiredNode.Address,
				"old_port", state.port, "new_port", desiredNode.Port,
				"old_schema", state.apiSchema, "new_schema", desiredNode.APISchema,
				"old_path", state.apiPath, "new_path", desiredNode.APIPath,
			)
			state.cancel()
			if state.conn != nil {
				state.conn.Close()
			}
			delete(nm.nodes, name)
			toStart[name] = desiredNode
		}
	}

	// Start new nodes
	for name, dbNode := range desired {
		if _, exists := nm.nodes[name]; !exists {
			toStart[name] = dbNode
		}
	}

	for _, dbNode := range toStart {
		nm.startNode(dbNode)
	}
}

// loadActiveNodes loads enabled nodes from database.
func (nm *NodeMonitor) loadActiveNodes() ([]db.DBNode, error) {
	nodes, err := db.LoadNodesFromDB(nm.manager, nm.cfg)
	if err != nil {
		return nil, err
	}

	grpcToken, err := nm.loadControlPlaneGRPCToken()
	if err != nil {
		return nil, err
	}
	for i := range nodes {
		nodes[i].APISchema = normalizeNodeSchema(nodes[i].APISchema)
		nodes[i].APIPath = normalizeNodePath(nodes[i].APIPath)
		if nodes[i].APISchema == "tls" {
			nodes[i].GRPCAuthToken = grpcToken
		}
	}
	return nodes, nil
}

func (nm *NodeMonitor) loadControlPlaneGRPCToken() (string, error) {
	var token sql.NullString
	err := nm.manager.ExecuteLowPriority(func(db dbmanager.DBExecutor) error {
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

// startNode starts monitoring a single node.
func (nm *NodeMonitor) startNode(dbNode db.DBNode) {
	ctx, cancel := context.WithCancel(nm.globalCtx)

	state := &nodeState{
		nodeUUID:      dbNode.UUID,
		nodeName:      dbNode.Name,
		address:       dbNode.Address,
		port:          dbNode.Port,
		apiSchema:     dbNode.APISchema,
		apiPath:       dbNode.APIPath,
		grpcAuthToken: dbNode.GRPCAuthToken,
		ctx:           ctx,
		cancel:        cancel,
	}

	nm.nodes[dbNode.Name] = state

	// Mark as connecting in DB
	nm.updateConnectionStatus(dbNode.Name, false, true, "Connecting...")

	go nm.monitorNode(state)

	nm.cfg.Logger.Info("Started monitoring node", "node", dbNode.Name, "address", dbNode.Address, "port", dbNode.Port, "schema", dbNode.APISchema, "path", dbNode.APIPath)
}

// monitorNode monitors a single node with reconnection logic.
func (nm *NodeMonitor) monitorNode(state *nodeState) {
	const (
		minBackoff = 2 * time.Second
		maxBackoff = 60 * time.Second
	)
	backoff := minBackoff

	for {
		if state.ctx.Err() != nil {
			nm.cfg.Logger.Debug("Node monitor stopped", "node", state.nodeName)
			return
		}

		nm.connectAndStream(state)
		if state.ctx.Err() != nil {
			nm.cfg.Logger.Debug("Node monitor stopped", "node", state.nodeName)
			return
		}

		wait := withJitter(backoff, 0.2)
		nm.cfg.Logger.Debug("Scheduling node reconnect", "node", state.nodeName, "wait", wait.String())

		timer := time.NewTimer(wait)
		select {
		case <-state.ctx.Done():
			timer.Stop()
			nm.cfg.Logger.Debug("Node monitor stopped", "node", state.nodeName)
			return
		case <-timer.C:
		}

		if backoff < maxBackoff {
			backoff = minDuration(maxBackoff, backoff*2)
		}
	}
}

// connectAndStream establishes connection and starts streaming.
func (nm *NodeMonitor) connectAndStream(state *nodeState) {
	state.mutex.Lock()
	if state.isConnecting {
		state.mutex.Unlock()
		return
	}
	state.isConnecting = true
	state.mutex.Unlock()

	// Create gRPC connection
	url := fmt.Sprintf("%s:%d", state.address, state.port)
	opts := []grpc.DialOption{}

	if state.address == "127.0.0.1" || strings.EqualFold(state.address, "localhost") {
		nm.cfg.Logger.Warn("Node address points to panel container loopback; use service name or host IP", "node", state.nodeName, "address", state.address)
	}

	apiSchema := normalizeNodeSchema(state.apiSchema)
	switch apiSchema {
	case "tls":
		if strings.TrimSpace(state.grpcAuthToken) == "" {
			nm.cfg.Logger.Warn("Failed to connect to node: missing global gRPC token", "node", state.nodeName)
			nm.updateConnectionStatus(state.nodeName, false, false, "Missing gRPC token in keygen")
			state.mutex.Lock()
			state.isConnecting = false
			state.mutex.Unlock()
			return
		}
		tlsCfg := &tls.Config{MinVersion: tls.VersionTLS12}
		if nm.cfg != nil && nm.cfg.Panel.AllowInsecureHTTP {
			tlsCfg.InsecureSkipVerify = true
			nm.cfg.Logger.Warn("Node TLS verification is disabled by EXODUS_ALLOW_INSECURE_HTTP", "node", state.nodeName)
		}
		opts = append(opts, grpc.WithTransportCredentials(credentials.NewTLS(tlsCfg)))
		nm.cfg.Logger.Debug("Using TLS + gRPC token for node gRPC", "node", state.nodeName, "address", state.address)
	case "mtls":
		tlsCfg, tlsErr := nm.loadNodeMTLSConfig(state.ctx)
		if tlsErr != nil {
			nm.cfg.Logger.Warn("Failed to build mTLS config for node", "node", state.nodeName, "error", tlsErr)
			nm.updateConnectionStatus(state.nodeName, false, false, fmt.Sprintf("mTLS config failed: %v", tlsErr))
			state.mutex.Lock()
			state.isConnecting = false
			state.mutex.Unlock()
			return
		}
		opts = append(opts, grpc.WithTransportCredentials(credentials.NewTLS(tlsCfg)))
		nm.cfg.Logger.Debug("Using mTLS for node gRPC", "node", state.nodeName, "address", state.address)
	default:
		nm.cfg.Logger.Warn("Node gRPC connection is insecure", "node", state.nodeName, "schema", apiSchema)
		opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))
	}

	cleanPath := normalizeNodePath(state.apiPath)
	opts = append(opts, grpc.WithUnaryInterceptor(pathPrefixUnaryInterceptor(cleanPath, state.grpcAuthToken)))
	opts = append(opts, grpc.WithStreamInterceptor(pathPrefixStreamInterceptor(cleanPath, state.grpcAuthToken)))
	if cleanPath != "" {
		nm.cfg.Logger.Debug("Using gRPC path prefix", "node", state.nodeName, "prefix", cleanPath)
	}

	conn, err := grpc.NewClient(url, opts...)
	if err != nil {
		nm.cfg.Logger.Warn("Failed to connect to node", "node", state.nodeName, "address", url, "error", err)
		nm.updateConnectionStatus(state.nodeName, false, false, fmt.Sprintf("Connection failed: %v", err))
		state.mutex.Lock()
		state.isConnecting = false
		state.mutex.Unlock()
		return
	}

	client := proto.NewNodeServiceClient(conn)

	// Create stream
	stream, err := client.StreamNodeData(state.ctx)
	if err != nil {
		nm.cfg.Logger.Warn("Failed to create stream", "node", state.nodeName, "error", err)
		conn.Close()
		nm.updateConnectionStatus(state.nodeName, false, false, fmt.Sprintf("Stream failed: %v", err))
		state.mutex.Lock()
		state.isConnecting = false
		state.mutex.Unlock()
		return
	}

	// Send initial config
	if err := stream.Send(&proto.NodeDataRequest{
		Request: &proto.NodeDataRequest_Config{
			Config: &proto.StreamConfig{
				IntervalSeconds: 20, // Default interval
			},
		},
	}); err != nil {
		nm.cfg.Logger.Warn("Failed to send config", "node", state.nodeName, "error", err)
		conn.Close()
		nm.updateConnectionStatus(state.nodeName, false, false, fmt.Sprintf("Config failed: %v", err))
		state.mutex.Lock()
		state.isConnecting = false
		state.mutex.Unlock()
		return
	}

	// Connection successful
	state.mutex.Lock()
	state.conn = conn
	state.client = client
	state.stream = stream
	state.isConnected = true
	state.isConnecting = false
	state.lastError = ""
	state.mutex.Unlock()

	// Update DB status (single update on connect)
	nm.updateConnectionStatus(state.nodeName, true, false, "Connected")

	nm.cfg.Logger.Info("Node connected", "node", state.nodeName)

	// Start stream receiver
	nm.receiveStream(state)
}

// receiveStream receives and processes stream data.
func (nm *NodeMonitor) receiveStream(state *nodeState) {
	for {
		resp, err := state.stream.Recv()
		if err == io.EOF {
			nm.cfg.Logger.Warn("Stream closed by node", "node", state.nodeName)
			nm.handleDisconnect(state, "Stream closed")
			return
		}
		if err != nil {
			if st, ok := status.FromError(err); ok {
				if st.Code() == codes.Canceled && state.ctx.Err() != nil {
					nm.cfg.Logger.Debug("Stream canceled", "node", state.nodeName)
					return
				}
				if st.Code() == codes.Unavailable {
					nm.cfg.Logger.Warn("Node unavailable", "node", state.nodeName)
					nm.handleDisconnect(state, "Node unavailable")
					return
				}
			}
			nm.cfg.Logger.Warn("Stream error", "node", state.nodeName, "error", err)
			nm.handleDisconnect(state, fmt.Sprintf("Stream error: %v", err))
			return
		}

		// Process response
		nm.processResponse(state.nodeName, resp)
	}
}

// handleDisconnect handles node disconnection.
func (nm *NodeMonitor) handleDisconnect(state *nodeState, reason string) {
	state.mutex.Lock()
	wasConnected := state.isConnected
	state.isConnected = false
	state.isConnecting = false
	state.lastError = reason
	if state.stream != nil {
		state.stream = nil
	}
	if state.conn != nil {
		state.conn.Close()
		state.conn = nil
	}
	state.mutex.Unlock()

	// Update DB status (single update on disconnect)
	if wasConnected {
		nm.updateConnectionStatus(state.nodeName, false, false, reason)
		nm.cfg.Logger.Info("Node disconnected", "node", state.nodeName, "reason", reason)
	} else {
		nm.cfg.Logger.Warn("Node disconnected before ready", "node", state.nodeName, "reason", reason)
	}

}

func withJitter(base time.Duration, factor float64) time.Duration {
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

func minDuration(a, b time.Duration) time.Duration {
	if a <= b {
		return a
	}
	return b
}

// processResponse processes node response data.
func (nm *NodeMonitor) processResponse(nodeName string, resp *proto.NodeDataResponse) {
	switch payload := resp.Response.(type) {
	case *proto.NodeDataResponse_Stats:
		nm.cfg.Logger.Trace("Node stats received", "node", nodeName)
		nm.updateNodeRuntimeFromStats(nodeName, payload.Stats.GetStats())
	case *proto.NodeDataResponse_Users:
		nm.cfg.Logger.Trace("Node users received", "node", nodeName)
	case *proto.NodeDataResponse_LogData:
		nm.cfg.Logger.Trace("Node log data received", "node", nodeName)
	default:
		nm.cfg.Logger.Trace("Node message received", "node", nodeName)
	}
}

func (nm *NodeMonitor) updateNodeRuntimeFromStats(nodeName string, stats []*proto.Stat) {
	if len(stats) == 0 {
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

	trafficDelta := extractTrafficStatsDelta(stats)

	singboxVersion := firstNonEmpty(values["singbox_version"])
	nodeVersion := firstNonEmpty(values["node_version"])
	singboxUptime := firstNonEmpty(values["singbox_uptime"])
	cpuCount := parseOptionalInt(values["cpu_count"])
	cpuModel := parseOptionalString(values["cpu_model"])
	totalRAM := parseOptionalString(values["total_ram"])
	usersOnline := trafficDelta.UsersOnline

	persistedNodeUUID := ""
	firstConnectedEvents := make([]notifications.Event, 0)
	err := nm.manager.ExecuteLowPriority(func(db dbmanager.DBExecutor) error {
		var (
			nodeUUID              string
			nodeID                int64
			consumptionMultiplier int64
		)
		if err := db.QueryRow(`SELECT uuid, id, consumption_multiplier FROM nodes WHERE name = ?`, nodeName).Scan(&nodeUUID, &nodeID, &consumptionMultiplier); err != nil {
			return err
		}
		persistedNodeUUID = nodeUUID

		query := `
			UPDATE nodes
			SET singbox_version = COALESCE(?, singbox_version),
			    node_version = COALESCE(?, node_version),
			    singbox_uptime = COALESCE(?, singbox_uptime),
			    cpu_count = COALESCE(?, cpu_count),
			    cpu_model = COALESCE(?, cpu_model),
			    total_ram = COALESCE(?, total_ram),
			    users_online = COALESCE(?, users_online),
			    updated_at = CURRENT_TIMESTAMP
			WHERE name = ?`
		if _, execErr := db.Exec(query, singboxVersion, nodeVersion, singboxUptime, cpuCount, cpuModel, totalRAM, usersOnline, nodeName); execErr != nil {
			return execErr
		}

		if trafficDelta.TotalUploadBytes > 0 || trafficDelta.TotalDownloadBytes > 0 {
			totalBytes := trafficDelta.TotalUploadBytes + trafficDelta.TotalDownloadBytes
			if _, execErr := db.Exec(`
				INSERT INTO nodes_usage_history (node_uuid, download_bytes, upload_bytes, total_bytes)
				VALUES (?, ?, ?, ?)
				ON CONFLICT (node_uuid, created_at)
				DO UPDATE SET
					download_bytes = nodes_usage_history.download_bytes + EXCLUDED.download_bytes,
					upload_bytes = nodes_usage_history.upload_bytes + EXCLUDED.upload_bytes,
					total_bytes = nodes_usage_history.total_bytes + EXCLUDED.total_bytes,
					updated_at = now()
			`, nodeUUID, trafficDelta.TotalDownloadBytes, trafficDelta.TotalUploadBytes, totalBytes); execErr != nil {
				return execErr
			}

			if _, execErr := db.Exec(`
				UPDATE nodes
				SET traffic_used_bytes = COALESCE(traffic_used_bytes, 0) + ?, updated_at = CURRENT_TIMESTAMP
				WHERE uuid = ?
			`, totalBytes, nodeUUID); execErr != nil {
				return execErr
			}
		}

		if len(trafficDelta.UserBytesByName) == 0 {
			return nil
		}

		usernames := make([]string, 0, len(trafficDelta.UserBytesByName))
		for username := range trafficDelta.UserBytesByName {
			if strings.TrimSpace(username) != "" {
				usernames = append(usernames, username)
			}
		}
		if len(usernames) == 0 {
			return nil
		}

		rows, queryErr := db.Query(`
			SELECT t_id, username
			FROM users
			WHERE status = 'ACTIVE' AND username = ANY(?)
		`, usernames)
		if queryErr != nil {
			return queryErr
		}
		defer rows.Close()

		userIDs := make(map[string]int64, len(usernames))
		for rows.Next() {
			var (
				userID   int64
				username string
			)
			if scanErr := rows.Scan(&userID, &username); scanErr != nil {
				return scanErr
			}
			userIDs[strings.TrimSpace(username)] = userID
		}
		if rowsErr := rows.Err(); rowsErr != nil {
			return rowsErr
		}

		firstConnectedByID := make(map[int64]bool, len(userIDs))
		if len(userIDs) > 0 {
			ids := make([]int64, 0, len(userIDs))
			for _, userID := range userIDs {
				ids = append(ids, userID)
			}
			firstRows, firstErr := db.Query(`
				SELECT t_id, first_connected_at IS NOT NULL
				FROM user_traffic
				WHERE t_id = ANY(?)
			`, ids)
			if firstErr != nil {
				return firstErr
			}
			defer firstRows.Close()
			for firstRows.Next() {
				var (
					userID            int64
					hasFirstConnected bool
				)
				if scanErr := firstRows.Scan(&userID, &hasFirstConnected); scanErr != nil {
					return scanErr
				}
				firstConnectedByID[userID] = hasFirstConnected
			}
			if rowsErr := firstRows.Err(); rowsErr != nil {
				return rowsErr
			}
		}

		usageDeltas := make([]userUsageDelta, 0, len(trafficDelta.UserBytesByName))
		for username, rawBytes := range trafficDelta.UserBytesByName {
			username = strings.TrimSpace(username)
			userID, ok := userIDs[username]
			if !ok {
				continue
			}
			if nm.cfg != nil && nm.cfg.Redis.UserUsageIgnoreBelowBytes > 0 && rawBytes < nm.cfg.Redis.UserUsageIgnoreBelowBytes {
				continue
			}

			effectiveBytes := applyConsumptionMultiplier(rawBytes, consumptionMultiplier)
			if effectiveBytes <= 0 {
				continue
			}

			usageDeltas = append(usageDeltas, userUsageDelta{
				UserID:       userID,
				Username:     username,
				TotalBytes:   effectiveBytes,
				HistoryBytes: rawBytes,
			})

			if !firstConnectedByID[userID] {
				firstConnectedEvents = append(firstConnectedEvents, notifications.Event{
					Scope: notifications.ScopeUser,
					Event: notifications.EventUserFirstConnected,
					Data: map[string]any{
						"tId":      userID,
						"username": username,
						"nodeUuid": nodeUUID,
						"nodeName": nodeName,
					},
				})
				firstConnectedByID[userID] = true
			}
		}

		if len(usageDeltas) == 0 {
			return nil
		}

		bulkCtx := nm.globalCtx
		if bulkCtx == nil {
			bulkCtx = context.Background()
		}

		if execErr := bulkUpsertUserTraffic(bulkCtx, db, usageDeltas, nodeUUID); execErr != nil {
			return execErr
		}

		if execErr := nm.recordNodeUserUsageHistory(bulkCtx, db, nodeID, usageDeltas); execErr != nil {
			return execErr
		}

		return nil
	})
	if err != nil {
		nm.cfg.Logger.Warn("Failed to persist node runtime stats", "node", nodeName, "error", err)
		return
	}

	if persistedNodeUUID != "" {
		nm.updateNodeMetricsSnapshot(persistedNodeUUID, usersOnline, trafficDelta)
	}
	for _, event := range firstConnectedEvents {
		notifications.Emit(context.Background(), nm.cfg, event)
	}
}

type trafficStatsDelta struct {
	TotalUploadBytes   int64
	TotalDownloadBytes int64
	UserBytesByName    map[string]int64
	UsersOnline        int
	InboundByTag       map[string]TagTrafficCounters
	OutboundByTag      map[string]TagTrafficCounters
}

func extractTrafficStatsDelta(stats []*proto.Stat) trafficStatsDelta {
	delta := trafficStatsDelta{
		UserBytesByName: make(map[string]int64),
		UsersOnline:     0,
		InboundByTag:    make(map[string]TagTrafficCounters),
		OutboundByTag:   make(map[string]TagTrafficCounters),
	}
	onlineUsers := make(map[string]struct{})

	for _, stat := range stats {
		if stat == nil {
			continue
		}

		name := strings.TrimSpace(stat.GetName())
		valueRaw := strings.TrimSpace(stat.GetValue())
		if name == "" || valueRaw == "" {
			continue
		}

		value, err := strconv.ParseInt(valueRaw, 10, 64)
		if err != nil || value <= 0 {
			continue
		}

		parts := strings.Split(name, ">>>")
		if len(parts) < 3 {
			continue
		}

		switch strings.ToLower(strings.TrimSpace(parts[0])) {
		case "inbound":
			if len(parts) != 4 || !strings.EqualFold(parts[2], "traffic") {
				continue
			}
			tag := strings.TrimSpace(parts[1])
			if tag == "" {
				continue
			}
			current := delta.InboundByTag[tag]
			switch strings.ToLower(strings.TrimSpace(parts[3])) {
			case "uplink":
				current.UploadBytes += value
			case "downlink":
				current.DownloadBytes += value
			default:
				continue
			}
			delta.InboundByTag[tag] = current
		case "outbound":
			if len(parts) != 4 || !strings.EqualFold(parts[2], "traffic") {
				continue
			}
			tag := strings.TrimSpace(parts[1])
			if tag == "" {
				continue
			}
			current := delta.OutboundByTag[tag]
			switch strings.ToLower(strings.TrimSpace(parts[3])) {
			case "uplink":
				delta.TotalUploadBytes += value
				current.UploadBytes += value
			case "downlink":
				delta.TotalDownloadBytes += value
				current.DownloadBytes += value
			default:
				continue
			}
			delta.OutboundByTag[tag] = current
		case "user":
			username := strings.TrimSpace(parts[1])
			if username == "" {
				continue
			}

			if len(parts) == 4 && strings.EqualFold(parts[2], "traffic") {
				switch strings.ToLower(strings.TrimSpace(parts[3])) {
				case "uplink", "downlink":
					delta.UserBytesByName[username] += value
					onlineUsers[username] = struct{}{}
				}
			}
		}
	}

	delta.UsersOnline = len(onlineUsers)
	return delta
}

func applyConsumptionMultiplier(totalBytes int64, multiplierNano int64) int64 {
	const nanoScale = int64(1_000_000_000)

	if totalBytes <= 0 || multiplierNano <= 0 {
		return 0
	}
	if multiplierNano == nanoScale {
		return totalBytes
	}

	scaled := math.Floor((float64(multiplierNano) / float64(nanoScale)) * float64(totalBytes))
	if scaled <= 0 {
		return 0
	}
	if scaled >= float64(math.MaxInt64) {
		return math.MaxInt64
	}
	return int64(scaled)
}

func (nm *NodeMonitor) recordNodeUserUsageHistory(ctx context.Context, db dbmanager.DBExecutor, nodeID int64, usageDeltas []userUsageDelta) error {
	if len(usageDeltas) == 0 {
		return nil
	}
	if nm.usageRecorder != nil {
		userBytes := make(map[int64]int64, len(usageDeltas))
		for _, delta := range usageDeltas {
			if delta.UserID <= 0 || delta.HistoryBytes <= 0 {
				continue
			}
			userBytes[delta.UserID] += delta.HistoryBytes
		}
		if len(userBytes) == 0 {
			return nil
		}
		if err := nm.usageRecorder.RecordNodeUserUsage(ctx, nodeID, userBytes); err == nil {
			return nil
		} else if nm.cfg != nil && nm.cfg.Logger != nil {
			nm.cfg.Logger.Warn("Failed to enqueue node user usage history in Redis, falling back to direct database write", "error", err)
		}
	}
	return bulkUpsertNodeUserUsageHistory(ctx, db, nodeID, usageDeltas)
}

func bulkUpsertUserTraffic(ctx context.Context, db dbmanager.DBExecutor, usageDeltas []userUsageDelta, nodeUUID string) error {
	const chunkSize = 1000

	for start := 0; start < len(usageDeltas); start += chunkSize {
		end := min(start+chunkSize, len(usageDeltas))
		chunk := usageDeltas[start:end]

		var query strings.Builder
		args := make([]any, 0, len(chunk)*3)

		query.WriteString(`
			INSERT INTO user_traffic (
				t_id, used_traffic_bytes, lifetime_used_traffic_bytes,
				online_at, last_connected_node_uuid, first_connected_at
			)
			SELECT
				v.t_id,
				v.total_bytes,
				v.total_bytes,
				now(),
				v.last_connected_node_uuid,
				now()
			FROM (VALUES `)

		for i, delta := range chunk {
			if i > 0 {
				query.WriteString(", ")
			}
			query.WriteString("(?::bigint, ?::bigint, ?::uuid)")
			args = append(args, delta.UserID, delta.TotalBytes, nodeUUID)
		}

		query.WriteString(`) AS v(t_id, total_bytes, last_connected_node_uuid)
			ON CONFLICT (t_id)
			DO UPDATE SET
				used_traffic_bytes = user_traffic.used_traffic_bytes + EXCLUDED.used_traffic_bytes,
				lifetime_used_traffic_bytes = user_traffic.lifetime_used_traffic_bytes + EXCLUDED.lifetime_used_traffic_bytes,
				online_at = now(),
				last_connected_node_uuid = EXCLUDED.last_connected_node_uuid,
				first_connected_at = COALESCE(user_traffic.first_connected_at, now())
		`)

		if _, err := db.ExecContext(ctx, query.String(), args...); err != nil {
			return err
		}
	}

	return nil
}

func bulkUpsertNodeUserUsageHistory(ctx context.Context, db dbmanager.DBExecutor, nodeID int64, usageDeltas []userUsageDelta) error {
	const chunkSize = 1000

	for start := 0; start < len(usageDeltas); start += chunkSize {
		end := min(start+chunkSize, len(usageDeltas))
		chunk := usageDeltas[start:end]

		var query strings.Builder
		args := make([]any, 0, len(chunk)*3)

		query.WriteString(`
			INSERT INTO nodes_user_usage_history (node_id, user_id, total_bytes)
			VALUES `)

		for i, delta := range chunk {
			if i > 0 {
				query.WriteString(", ")
			}
			query.WriteString("(?::bigint, ?::bigint, ?::bigint)")
			args = append(args, nodeID, delta.UserID, delta.HistoryBytes)
		}

		query.WriteString(`
			ON CONFLICT (node_id, created_at, user_id)
			DO UPDATE SET
				total_bytes = nodes_user_usage_history.total_bytes + EXCLUDED.total_bytes,
				updated_at = now()
		`)

		if _, err := db.ExecContext(ctx, query.String(), args...); err != nil {
			return err
		}
	}

	return nil
}

func parseOptionalInt(value string) any {
	if strings.TrimSpace(value) == "" {
		return nil
	}
	parsed, err := strconv.Atoi(strings.TrimSpace(value))
	if err != nil {
		return nil
	}
	return parsed
}

func parseOptionalString(value string) any {
	trimmed := strings.TrimSpace(value)
	if trimmed == "" {
		return nil
	}
	return trimmed
}

func firstNonEmpty(values ...string) any {
	for _, value := range values {
		trimmed := strings.TrimSpace(value)
		if trimmed != "" {
			return trimmed
		}
	}
	return nil
}

// updateConnectionStatus updates node connection status in database (only on change).
func (nm *NodeMonitor) updateConnectionStatus(nodeName string, isConnected, isConnecting bool, message string) {
	var notificationEvent notifications.Event
	err := nm.manager.ExecuteHighPriority(func(db dbmanager.DBExecutor) error {
		// Get current status
		var (
			nodeUUID          string
			nodeAddress       string
			nodePort          sql.NullInt64
			currentConnected  bool
			currentConnecting bool
			currentMessage    sql.NullString
		)

		err := db.QueryRow(`SELECT uuid, address, port, is_connected, is_connecting, last_status_message FROM nodes WHERE name = ?`, nodeName).
			Scan(&nodeUUID, &nodeAddress, &nodePort, &currentConnected, &currentConnecting, &currentMessage)
		if err != nil {
			if err == sql.ErrNoRows {
				nm.cfg.Logger.Debug("Node not found in DB", "node", nodeName)
				return nil
			}
			return fmt.Errorf("query node status: %w", err)
		}

		// Check if status actually changed
		msgStr := ""
		if currentMessage.Valid {
			msgStr = currentMessage.String
		}

		if currentConnected == isConnected && currentConnecting == isConnecting && msgStr == message {
			// No change, skip update
			return nil
		}

		// Update status
		_, err = db.Exec(`
			UPDATE nodes 
			SET is_connected = ?, is_connecting = ?, last_status_message = ?, last_status_change = CURRENT_TIMESTAMP
			WHERE name = ?`,
			isConnected, isConnecting, message, nodeName)

		if err != nil {
			return fmt.Errorf("update node status: %w", err)
		}

		if currentConnected != isConnected {
			eventName := notifications.EventNodeConnectionRestored
			if !isConnected {
				eventName = notifications.EventNodeConnectionLost
			}
			notificationEvent = notifications.Event{
				Scope: notifications.ScopeNode,
				Event: eventName,
				Data: map[string]any{
					"uuid":        nodeUUID,
					"name":        nodeName,
					"address":     nodeAddress,
					"port":        nodePort.Int64,
					"isConnected": isConnected,
					"message":     message,
				},
			}
		}

		nm.cfg.Logger.Debug("Node status updated", "node", nodeName, "connected", isConnected, "message", message)
		return nil
	})

	if err != nil {
		nm.cfg.Logger.Warn("Failed to update node status", "node", nodeName, "error", err)
		return
	}
	if notificationEvent.Event != "" {
		notifications.Emit(context.Background(), nm.cfg, notificationEvent)
	}
}

// stopAll stops all node monitors.
func (nm *NodeMonitor) stopAll() {
	nm.nodesLock.Lock()
	defer nm.nodesLock.Unlock()

	for name, state := range nm.nodes {
		state.cancel()
		if state.conn != nil {
			state.conn.Close()
		}
		delete(nm.nodes, name)
	}

	nm.cfg.Logger.Info("All node monitors stopped")
}

// Stop gracefully stops the monitor.
func (nm *NodeMonitor) Stop() {
	nm.globalCancel()
}

// RequestSync triggers an immediate DB sync (non-blocking).
func (nm *NodeMonitor) RequestSync() {
	if nm == nil {
		return
	}
	select {
	case nm.syncNow <- struct{}{}:
	default:
	}
}

// RequestDeploy triggers config deploy to connected nodes (non-blocking).
func (nm *NodeMonitor) RequestDeploy(restart bool, nodeUUIDs ...string) {
	if nm == nil {
		return
	}
	normalizedTargets := normalizeNodeUUIDTargets(nodeUUIDs)
	req := deployRequest{
		Restart:   restart,
		NodeUUIDs: normalizedTargets,
	}
	if nm.cfg != nil && nm.cfg.Logger != nil {
		nm.cfg.Logger.Debug("Node deploy requested", "restart", restart, "node_targets", len(normalizedTargets))
	}
	select {
	case nm.deployNow <- req:
		if nm.cfg != nil && nm.cfg.Logger != nil {
			nm.cfg.Logger.Trace("Node deploy enqueued", "restart", restart, "node_targets", len(normalizedTargets))
		}
	default:
		select {
		case <-nm.deployNow:
		default:
		}
		nm.deployNow <- req
		if nm.cfg != nil && nm.cfg.Logger != nil {
			nm.cfg.Logger.Debug("Node deploy queue replaced previous pending request", "restart", restart, "node_targets", len(normalizedTargets))
		}
	}
}

// RequestSRSDeploy triggers SRS list sync to connected nodes (non-blocking).
func (nm *NodeMonitor) RequestSRSDeploy() {
	if nm == nil {
		return
	}
	if nm.cfg != nil && nm.cfg.Logger != nil {
		nm.cfg.Logger.Debug("Node SRS sync requested")
	}
	select {
	case nm.srsSyncNow <- struct{}{}:
	default:
	}
}

func normalizeNodeSchema(value string) string {
	normalized := strings.ToLower(strings.TrimSpace(value))
	switch normalized {
	case "tls":
		return "tls"
	case "mtls", "grpc", "grpcs", "https", "":
		return "mtls"
	default:
		return "mtls"
	}
}

func normalizeNodePath(value string) string {
	trimmed := strings.TrimSpace(value)
	if trimmed == "" || trimmed == "/" {
		return ""
	}
	return "/" + strings.Trim(trimmed, "/")
}

func nodeConfigChanged(state *nodeState, desired db.DBNode) bool {
	if state == nil {
		return true
	}
	if state.address != desired.Address {
		return true
	}
	if state.port != desired.Port {
		return true
	}
	if normalizeNodeSchema(state.apiSchema) != normalizeNodeSchema(desired.APISchema) {
		return true
	}
	if normalizeNodePath(state.apiPath) != normalizeNodePath(desired.APIPath) {
		return true
	}
	if strings.TrimSpace(state.grpcAuthToken) != strings.TrimSpace(desired.GRPCAuthToken) {
		return true
	}
	return false
}

func normalizeNodeUUIDTargets(raw []string) []string {
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
	if len(result) == 0 {
		return nil
	}
	return result
}

func (nm *NodeMonitor) removeNodeMetrics(nodeUUID string) {
	if nm == nil || strings.TrimSpace(nodeUUID) == "" {
		return
	}
	nm.metricsLock.Lock()
	delete(nm.metricsByNodeUUID, nodeUUID)
	nm.metricsLock.Unlock()
}

func (nm *NodeMonitor) updateNodeMetricsSnapshot(nodeUUID string, usersOnline int, delta trafficStatsDelta) {
	if nm == nil || strings.TrimSpace(nodeUUID) == "" {
		return
	}

	nm.metricsLock.Lock()
	defer nm.metricsLock.Unlock()

	snapshot, exists := nm.metricsByNodeUUID[nodeUUID]
	if !exists || snapshot == nil {
		snapshot = &NodeMetricsSnapshot{
			NodeUUID:  nodeUUID,
			Inbounds:  make(map[string]TagTrafficCounters),
			Outbounds: make(map[string]TagTrafficCounters),
		}
		nm.metricsByNodeUUID[nodeUUID] = snapshot
	}

	snapshot.UsersOnline = usersOnline
	snapshot.UpdatedAt = time.Now().UTC()

	for tag, item := range delta.InboundByTag {
		current := snapshot.Inbounds[tag]
		current.UploadBytes += item.UploadBytes
		current.DownloadBytes += item.DownloadBytes
		snapshot.Inbounds[tag] = current
	}

	for tag, item := range delta.OutboundByTag {
		current := snapshot.Outbounds[tag]
		current.UploadBytes += item.UploadBytes
		current.DownloadBytes += item.DownloadBytes
		snapshot.Outbounds[tag] = current
	}
}

func (nm *NodeMonitor) SnapshotNodeMetrics() map[string]NodeMetricsSnapshot {
	if nm == nil {
		return map[string]NodeMetricsSnapshot{}
	}

	nm.metricsLock.RLock()
	defer nm.metricsLock.RUnlock()

	result := make(map[string]NodeMetricsSnapshot, len(nm.metricsByNodeUUID))
	for nodeUUID, source := range nm.metricsByNodeUUID {
		if source == nil {
			continue
		}

		copySnapshot := NodeMetricsSnapshot{
			NodeUUID:    source.NodeUUID,
			UsersOnline: source.UsersOnline,
			Inbounds:    make(map[string]TagTrafficCounters, len(source.Inbounds)),
			Outbounds:   make(map[string]TagTrafficCounters, len(source.Outbounds)),
			UpdatedAt:   source.UpdatedAt,
		}
		for tag, item := range source.Inbounds {
			copySnapshot.Inbounds[tag] = item
		}
		for tag, item := range source.Outbounds {
			copySnapshot.Outbounds[tag] = item
		}

		result[nodeUUID] = copySnapshot
	}
	return result
}
