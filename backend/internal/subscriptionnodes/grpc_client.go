package subscriptionnodes

import (
	"context"
	"fmt"
	"math/rand/v2"
	"strings"
	"sync"
	"time"

	"exodus/internal/nodes/grpcauth"
	"exodus/internal/proto"

	"google.golang.org/grpc"
	_ "google.golang.org/grpc/encoding/gzip"
)

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

	isConnected  bool
	isConnecting bool
	lastError    string
	mutex        sync.RWMutex
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

	apiSchema := strings.ToLower(strings.TrimSpace(state.apiSchema))
	var useMTLS, useTLS, skipVerify bool
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
		useTLS = true
		if sm.cfg != nil && sm.cfg.Panel.AllowInsecureHTTP {
			skipVerify = true
			sm.cfg.Logger.Warn("Subscription TLS verification is disabled by EXODUS_ALLOW_INSECURE_HTTP", "node", state.nodeName)
		}
	default:
		useMTLS = true
	}

	cleanPath := normalizeSubPath(state.apiPath)
	opts, err := grpcauth.GetDialOptions(state.ctx, sm.manager, useMTLS, useTLS, skipVerify, cleanPath, state.grpcAuthToken, "")
	if err != nil {
		sm.cfg.Logger.Warn("Failed to build mTLS config for subscription node", "node", state.nodeName, "error", err)
		sm.updateConnectionStatus(state.nodeName, false, false, fmt.Sprintf("mTLS config failed: %v", err))
		state.mutex.Lock()
		state.isConnecting = false
		state.mutex.Unlock()
		return
	}

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

func (sm *SubNodeMonitor) sendNodeRequest(state *subNodeState, req *proto.NodeDataRequest) error {
	state.mutex.RLock()
	stream := state.stream
	connected := state.isConnected
	state.mutex.RUnlock()
	if !connected || stream == nil {
		return fmt.Errorf("stream is not connected")
	}

	return stream.Send(req)
}

func (sm *SubNodeMonitor) markStreamActivity(state *subNodeState) {
	state.mutex.Lock()
	if state.stream != nil {
		state.lastResponseAt = time.Now()
	}
	state.mutex.Unlock()
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
