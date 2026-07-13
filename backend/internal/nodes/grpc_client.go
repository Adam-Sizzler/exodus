package users

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
)

type nodeState struct {
	nodeUUID      string
	nodeName      string
	address       string
	port          int
	proxyURL      string
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

	targetAddr := fmt.Sprintf("%s:%d", state.address, state.port)

	if state.address == "127.0.0.1" || strings.EqualFold(state.address, "localhost") {
		nm.cfg.Logger.Warn("Node address points to panel container loopback; use service name or host IP", "node", state.nodeName, "address", state.address)
	}

	apiSchema := normalizeNodeSchema(state.apiSchema)
	var useMTLS, useTLS, skipVerify bool
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
		useTLS = true
		if nm.cfg != nil && nm.cfg.Panel.AllowInsecureHTTP {
			skipVerify = true
			nm.cfg.Logger.Warn("Node TLS verification is disabled by EXODUS_ALLOW_INSECURE_HTTP", "node", state.nodeName)
		}
		nm.cfg.Logger.Debug("Using TLS + gRPC token for node gRPC", "node", state.nodeName, "address", state.address)
	case "mtls":
		useMTLS = true
		nm.cfg.Logger.Debug("Using mTLS for node gRPC", "node", state.nodeName, "address", state.address)
	default:
		nm.cfg.Logger.Warn("Node gRPC connection is insecure", "node", state.nodeName, "schema", apiSchema)
	}

	cleanPath := normalizeNodePath(state.apiPath)
	if cleanPath != "" {
		nm.cfg.Logger.Debug("Using gRPC path prefix", "node", state.nodeName, "prefix", cleanPath)
	}

	opts, err := grpcauth.GetDialOptions(state.ctx, nm.manager, useMTLS, useTLS, skipVerify, cleanPath, state.grpcAuthToken, state.proxyURL)
	if err != nil {
		nm.cfg.Logger.Warn("Failed to build gRPC options for node", "node", state.nodeName, "error", err)
		nm.updateConnectionStatus(state.nodeName, false, false, err.Error())
		state.mutex.Lock()
		state.isConnecting = false
		state.mutex.Unlock()
		return
	}

	conn, err := grpc.NewClient(targetAddr, opts...)
	if err != nil {
		nm.cfg.Logger.Error("Failed to connect to node", "node", state.nodeName, "address", targetAddr, "error", err)
		nm.updateConnectionStatus(state.nodeName, false, false, fmt.Sprintf("Connection failed: %v", err))
		state.mutex.Lock()
		state.isConnecting = false
		state.mutex.Unlock()
		return
	}

	client := proto.NewNodeServiceClient(conn)

	stream, err := client.StreamNodeData(state.ctx)
	if err != nil {
		nm.cfg.Logger.Error("Failed to create stream", "node", state.nodeName, "error", err)
		conn.Close()
		nm.updateConnectionStatus(state.nodeName, false, false, fmt.Sprintf("Stream failed: %v", err))
		state.mutex.Lock()
		state.isConnecting = false
		state.mutex.Unlock()
		return
	}

	if err := stream.Send(&proto.NodeDataRequest{
		Request: &proto.NodeDataRequest_Config{
			Config: &proto.StreamConfig{
				IntervalSeconds: 15,
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

	state.mutex.Lock()
	state.conn = conn
	state.client = client
	state.stream = stream
	state.isConnected = true
	state.isConnecting = false
	state.lastError = ""
	state.mutex.Unlock()

	nm.updateConnectionStatus(state.nodeName, false, true, "")

	nm.cfg.Logger.Info("Node control-plane connected", "node", state.nodeName)
	nm.RequestDeploy(true, state.nodeUUID)

	nm.receiveStream(state)
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
