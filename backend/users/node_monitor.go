package users

import (
	"context"
	"database/sql"
	"fmt"
	"io"
	"sync"
	"time"

	"v2ray-stat/backend/config"
	"v2ray-stat/backend/db"
	"v2ray-stat/backend/db/manager"
	"v2ray-stat/proto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
)

// NodeMonitor dynamically manages node monitoring with status tracking.
type NodeMonitor struct {
	manager *manager.DatabaseManager
	cfg     *config.BackendConfig

	// Active node contexts
	nodes     map[string]*nodeState
	nodesLock sync.RWMutex

	// Global context for shutdown
	globalCtx    context.Context
	globalCancel context.CancelFunc
}

type nodeState struct {
	nodeName      string
	address       string
	port          int
	apiSchema     string
	apiPath       string
	ctx           context.Context
	cancel        context.CancelFunc
	conn          *grpc.ClientConn
	client        proto.NodeServiceClient
	stream        proto.NodeService_StreamNodeDataClient
	isConnected   bool
	isConnecting  bool
	lastError     string
	mutex         sync.RWMutex
	reconnectChan chan struct{}
}

// NewNodeMonitor creates a new NodeMonitor.
func NewNodeMonitor(manager *manager.DatabaseManager, cfg *config.BackendConfig) *NodeMonitor {
	return &NodeMonitor{
		manager: manager,
		cfg:     cfg,
		nodes:   make(map[string]*nodeState),
	}
}

// Start begins the node monitoring loop.
func (nm *NodeMonitor) Start(ctx context.Context, wg *sync.WaitGroup) {
	defer wg.Done()

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
		nm.cfg.Logger.Debug("Failed to load nodes from DB", "error", err)
		return
	}

	// Build desired state
	desired := make(map[string]db.DBNode)
	for _, n := range dbNodes {
		desired[n.Name] = n
	}

	nm.nodesLock.Lock()
	defer nm.nodesLock.Unlock()

	// Stop removed nodes
	for name, state := range nm.nodes {
		if _, exists := desired[name]; !exists {
			nm.cfg.Logger.Debug("Node removed from DB, stopping monitor", "node", name)
			state.cancel()
			if state.conn != nil {
				state.conn.Close()
			}
			delete(nm.nodes, name)
		}
	}

	// Start new nodes
	for name, dbNode := range desired {
		if _, exists := nm.nodes[name]; !exists {
			nm.startNode(dbNode)
		}
	}
}

// loadActiveNodes loads enabled nodes from database.
func (nm *NodeMonitor) loadActiveNodes() ([]db.DBNode, error) {
	return db.LoadNodesFromDB(nm.manager, nm.cfg)
}

// startNode starts monitoring a single node.
func (nm *NodeMonitor) startNode(dbNode db.DBNode) {
	ctx, cancel := context.WithCancel(nm.globalCtx)

	state := &nodeState{
		nodeName:      dbNode.Name,
		address:       dbNode.Address,
		port:          dbNode.Port,
		apiSchema:     dbNode.APISchema,
		apiPath:       dbNode.APIPath,
		ctx:           ctx,
		cancel:        cancel,
		reconnectChan: make(chan struct{}, 1),
	}

	nm.nodes[dbNode.Name] = state

	// Mark as connecting in DB
	nm.updateConnectionStatus(dbNode.Name, false, true, "Connecting...")

	go nm.monitorNode(state)

	nm.cfg.Logger.Info("Started monitoring node", "node", dbNode.Name, "address", dbNode.Address)
}

// monitorNode monitors a single node with reconnection logic.
func (nm *NodeMonitor) monitorNode(state *nodeState) {
	// Initial connection attempt
	nm.connectAndStream(state)

	// Reconnection loop
	reconnectTicker := time.NewTicker(10 * time.Second)
	defer reconnectTicker.Stop()

	for {
		select {
		case <-state.ctx.Done():
			nm.cfg.Logger.Debug("Node monitor stopped", "node", state.nodeName)
			return

		case <-reconnectTicker.C:
			if !state.isConnected {
				nm.connectAndStream(state)
			}

		case <-state.reconnectChan:
			// Immediate reconnect requested
			nm.connectAndStream(state)
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
	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	}

	if state.apiPath != "" {
		// Path prefix interceptors if needed
	}

	conn, err := grpc.NewClient(url, opts...)
	if err != nil {
		nm.cfg.Logger.Debug("Failed to connect to node", "node", state.nodeName, "error", err)
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
		nm.cfg.Logger.Debug("Failed to create stream", "node", state.nodeName, "error", err)
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
				IntervalSeconds: 10, // Default interval
			},
		},
	}); err != nil {
		nm.cfg.Logger.Debug("Failed to send config", "node", state.nodeName, "error", err)
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
			nm.cfg.Logger.Debug("Stream closed by node", "node", state.nodeName)
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
					nm.cfg.Logger.Debug("Node unavailable", "node", state.nodeName)
					nm.handleDisconnect(state, "Node unavailable")
					return
				}
			}
			nm.cfg.Logger.Debug("Stream error", "node", state.nodeName, "error", err)
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
	}

	// Request reconnect
	select {
	case state.reconnectChan <- struct{}{}:
	default:
	}
}

// processResponse processes node response data.
func (nm *NodeMonitor) processResponse(nodeName string, resp *proto.NodeDataResponse) {
	switch r := resp.Response.(type) {
	case *proto.NodeDataResponse_Stats:
		apiData := convertProtoToApiResponse(r.Stats)
		if err := updateProxyStats(nm.manager, nodeName, apiData, nm.cfg); err != nil {
			nm.cfg.Logger.Debug("Failed to update proxy stats", "node", nodeName, "error", err)
		}
		if err := updateUserStats(nm.manager, nodeName, apiData, nm.cfg); err != nil {
			nm.cfg.Logger.Debug("Failed to update user stats", "node", nodeName, "error", err)
		}

	case *proto.NodeDataResponse_LogData:
		// Process log data if needed
	}
}

// updateConnectionStatus updates node connection status in database (only on change).
func (nm *NodeMonitor) updateConnectionStatus(nodeName string, isConnected, isConnecting bool, message string) {
	err := nm.manager.ExecuteHighPriority(func(db *sql.DB) error {
		// Get current status
		var currentConnected, currentConnecting bool
		var currentMessage sql.NullString

		err := db.QueryRow(`SELECT is_connected, is_connecting, last_status_message FROM nodes WHERE name = ?`, nodeName).
			Scan(&currentConnected, &currentConnecting, &currentMessage)
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

		nm.cfg.Logger.Debug("Node status updated", "node", nodeName, "connected", isConnected, "message", message)
		return nil
	})

	if err != nil {
		nm.cfg.Logger.Debug("Failed to update node status", "node", nodeName, "error", err)
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
