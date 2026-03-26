package grpcapi

import (
	"context"
	"fmt"
	"io"
	"os"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/cerberus/subscription-page/backend/internal/proto"

	statuspb "google.golang.org/genproto/googleapis/rpc/status"
	"google.golang.org/grpc/codes"
)

type NodeServer struct {
	proto.UnimplementedNodeServiceServer

	startedAt time.Time
	version   string
	cpuCount  int
	cpuModel  string
	totalRAM  string

	streamMu     sync.RWMutex
	stream       proto.NodeService_StreamNodeDataServer
	streamSendMu sync.Mutex

	pendingMu sync.Mutex
	pending   map[string]chan *proto.SubscriptionBridgeResponse
	reqSeq    atomic.Uint64

	subpageConfigsMu sync.RWMutex
	subpageConfigs   map[string][]byte
}

func NewNodeServer(version string) *NodeServer {
	trimmedVersion := strings.TrimSpace(version)
	if trimmedVersion == "" {
		trimmedVersion = "sub"
	}

	return &NodeServer{
		startedAt:      time.Now().UTC(),
		version:        trimmedVersion,
		cpuCount:       runtime.NumCPU(),
		cpuModel:       detectCPUModel(),
		totalRAM:       detectTotalRAM(),
		pending:        make(map[string]chan *proto.SubscriptionBridgeResponse),
		subpageConfigs: make(map[string][]byte),
	}
}

func (s *NodeServer) StreamNodeData(stream proto.NodeService_StreamNodeDataServer) error {
	s.attachStream(stream)
	defer s.detachStream(stream)

	intervalSeconds := 20
	ticker := time.NewTicker(time.Duration(intervalSeconds) * time.Second)
	defer ticker.Stop()

	intervalUpdates := make(chan int, 1)
	requestUsers := make(chan struct{}, 1)
	recvErrCh := make(chan error, 1)

	go func() {
		for {
			req, err := stream.Recv()
			if err != nil {
				recvErrCh <- err
				return
			}
			s.handleStreamRequest(req, intervalUpdates, requestUsers)
		}
	}()

	sendStats := func() error {
		return s.sendResponse(stream, &proto.NodeDataResponse{
			Response: &proto.NodeDataResponse_Stats{
				Stats: &proto.GetApiStatsResponse{
					Stats: s.buildStats(),
				},
			},
		})
	}

	sendUsers := func() error {
		return s.sendResponse(stream, &proto.NodeDataResponse{
			Response: &proto.NodeDataResponse_Users{
				Users: &proto.ListUsersResponse{Users: []*proto.User{}},
			},
		})
	}

	if err := sendStats(); err != nil {
		return err
	}

	for {
		select {
		case <-stream.Context().Done():
			return nil
		case recvErr := <-recvErrCh:
			if recvErr == io.EOF {
				return nil
			}
			return recvErr
		case nextInterval := <-intervalUpdates:
			nextInterval = clampInterval(nextInterval)
			if nextInterval == intervalSeconds {
				continue
			}
			intervalSeconds = nextInterval
			ticker.Reset(time.Duration(intervalSeconds) * time.Second)
		case <-requestUsers:
			if err := sendUsers(); err != nil {
				return err
			}
		case <-ticker.C:
			if err := sendStats(); err != nil {
				return err
			}
		}
	}
}

func (s *NodeServer) QueryPanel(ctx context.Context, req *proto.SubscriptionBridgeRequest) (*proto.SubscriptionBridgeResponse, error) {
	if req == nil {
		return nil, fmt.Errorf("subscription bridge request is nil")
	}
	if strings.TrimSpace(req.Operation) == "" {
		return nil, fmt.Errorf("subscription bridge operation is required")
	}

	if strings.TrimSpace(req.RequestId) == "" {
		req.RequestId = s.nextRequestID()
	}

	stream := s.currentStream()
	if stream == nil {
		return nil, fmt.Errorf("panel stream is not connected")
	}

	responseCh := make(chan *proto.SubscriptionBridgeResponse, 1)
	s.pendingMu.Lock()
	s.pending[req.RequestId] = responseCh
	s.pendingMu.Unlock()

	cleanup := func() {
		s.pendingMu.Lock()
		delete(s.pending, req.RequestId)
		s.pendingMu.Unlock()
	}

	if err := s.sendResponse(stream, &proto.NodeDataResponse{
		Response: &proto.NodeDataResponse_SubscriptionRequest{SubscriptionRequest: req},
	}); err != nil {
		cleanup()
		return nil, err
	}

	select {
	case <-ctx.Done():
		cleanup()
		return nil, ctx.Err()
	case <-stream.Context().Done():
		cleanup()
		return nil, fmt.Errorf("panel stream context closed")
	case resp := <-responseCh:
		if resp == nil {
			return nil, fmt.Errorf("empty bridge response")
		}
		return resp, nil
	}
}

func (s *NodeServer) GetCachedSubpageConfig(uuid string) ([]byte, bool) {
	uuid = strings.TrimSpace(uuid)
	if uuid == "" {
		return nil, false
	}

	s.subpageConfigsMu.RLock()
	config, ok := s.subpageConfigs[uuid]
	s.subpageConfigsMu.RUnlock()
	if !ok {
		return nil, false
	}

	clone := make([]byte, len(config))
	copy(clone, config)
	return clone, true
}

func (s *NodeServer) ListUsers(context.Context, *proto.ListUsersRequest) (*proto.ListUsersResponse, error) {
	return &proto.ListUsersResponse{Users: []*proto.User{}}, nil
}

func (s *NodeServer) GetApiStats(context.Context, *proto.GetApiStatsRequest) (*proto.GetApiStatsResponse, error) {
	return &proto.GetApiStatsResponse{Stats: s.buildStats()}, nil
}

func (s *NodeServer) GetLogData(context.Context, *proto.GetLogDataRequest) (*proto.GetLogDataResponse, error) {
	return &proto.GetLogDataResponse{UserLogData: map[string]*proto.UserLogData{}}, nil
}

func (s *NodeServer) AddUsers(context.Context, *proto.AddUsersRequest) (*proto.OperationResponse, error) {
	return &proto.OperationResponse{Status: okStatus("users accepted")}, nil
}

func (s *NodeServer) DeleteUsers(context.Context, *proto.DeleteUsersRequest) (*proto.OperationResponse, error) {
	return &proto.OperationResponse{Status: okStatus("users removed")}, nil
}

func (s *NodeServer) SetUserEnabled(context.Context, *proto.SetUserEnabledRequest) (*proto.OperationResponse, error) {
	return &proto.OperationResponse{Status: okStatus("user status updated")}, nil
}

func (s *NodeServer) SubmitTask(_ context.Context, task *proto.NodeTask) (*statuspb.Status, error) {
	message := "task accepted"
	if task != nil && strings.TrimSpace(task.Operation) != "" {
		message = fmt.Sprintf("task accepted: %s", strings.TrimSpace(task.Operation))
	}
	return okStatus(message), nil
}

func (s *NodeServer) GetTaskStatus(context.Context, *proto.TaskStatusRequest) (*proto.TaskStatusResponse, error) {
	return &proto.TaskStatusResponse{Status: "success"}, nil
}

func (s *NodeServer) handleStreamRequest(req *proto.NodeDataRequest, intervalUpdates chan<- int, requestUsers chan<- struct{}) {
	if req == nil {
		return
	}

	switch payload := req.GetRequest().(type) {
	case *proto.NodeDataRequest_Config:
		if payload.Config == nil {
			return
		}
		if payload.Config.IntervalSeconds > 0 {
			select {
			case intervalUpdates <- int(payload.Config.IntervalSeconds):
			default:
			}
		}
	case *proto.NodeDataRequest_ListUsers:
		select {
		case requestUsers <- struct{}{}:
		default:
		}
	case *proto.NodeDataRequest_SubscriptionResponse:
		s.handleBridgeResponse(payload.SubscriptionResponse)
	case *proto.NodeDataRequest_SubpageConfigUpdate:
		s.applySubpageConfigUpdate(payload.SubpageConfigUpdate)
	}
}

func (s *NodeServer) handleBridgeResponse(resp *proto.SubscriptionBridgeResponse) {
	if resp == nil {
		return
	}
	requestID := strings.TrimSpace(resp.GetRequestId())
	if requestID == "" {
		return
	}

	s.pendingMu.Lock()
	responseCh, ok := s.pending[requestID]
	if ok {
		delete(s.pending, requestID)
	}
	s.pendingMu.Unlock()
	if !ok {
		return
	}

	select {
	case responseCh <- resp:
	default:
	}
}

func (s *NodeServer) applySubpageConfigUpdate(update *proto.SubpageConfigUpdate) {
	if update == nil {
		return
	}
	uuid := strings.TrimSpace(update.GetUuid())
	if uuid == "" {
		s.subpageConfigsMu.Lock()
		s.subpageConfigs = make(map[string][]byte)
		s.subpageConfigsMu.Unlock()
		return
	}

	payload := update.GetConfig()
	clone := make([]byte, len(payload))
	copy(clone, payload)

	s.subpageConfigsMu.Lock()
	s.subpageConfigs[uuid] = clone
	s.subpageConfigsMu.Unlock()
}

func (s *NodeServer) buildStats() []*proto.Stat {
	uptimeSeconds := int64(time.Since(s.startedAt).Seconds())
	if uptimeSeconds < 0 {
		uptimeSeconds = 0
	}

	return []*proto.Stat{
		{Name: "singbox_version", Value: "subscription-page"},
		{Name: "node_version", Value: s.version},
		{Name: "singbox_uptime", Value: strconv.FormatInt(uptimeSeconds, 10)},
		{Name: "cpu_count", Value: strconv.Itoa(s.cpuCount)},
		{Name: "cpu_model", Value: s.cpuModel},
		{Name: "total_ram", Value: s.totalRAM},
	}
}

func (s *NodeServer) currentStream() proto.NodeService_StreamNodeDataServer {
	s.streamMu.RLock()
	stream := s.stream
	s.streamMu.RUnlock()
	return stream
}

func (s *NodeServer) attachStream(stream proto.NodeService_StreamNodeDataServer) {
	s.streamMu.Lock()
	s.stream = stream
	s.streamMu.Unlock()

	s.failPending("panel stream replaced")
}

func (s *NodeServer) detachStream(stream proto.NodeService_StreamNodeDataServer) {
	s.streamMu.Lock()
	if s.stream == stream {
		s.stream = nil
	}
	s.streamMu.Unlock()

	s.failPending("panel stream disconnected")
}

func (s *NodeServer) sendResponse(stream proto.NodeService_StreamNodeDataServer, resp *proto.NodeDataResponse) error {
	s.streamSendMu.Lock()
	defer s.streamSendMu.Unlock()

	current := s.currentStream()
	if current == nil || current != stream {
		return fmt.Errorf("panel stream is not active")
	}

	return stream.Send(resp)
}

func (s *NodeServer) failPending(reason string) {
	s.pendingMu.Lock()
	pending := s.pending
	s.pending = make(map[string]chan *proto.SubscriptionBridgeResponse)
	s.pendingMu.Unlock()

	for requestID, responseCh := range pending {
		resp := &proto.SubscriptionBridgeResponse{
			RequestId:  requestID,
			StatusCode: 503,
			Error:      reason,
		}
		select {
		case responseCh <- resp:
		default:
		}
	}
}

func (s *NodeServer) nextRequestID() string {
	next := s.reqSeq.Add(1)
	return fmt.Sprintf("sub-%d-%d", time.Now().UnixNano(), next)
}

func okStatus(message string) *statuspb.Status {
	return &statuspb.Status{Code: int32(codes.OK), Message: message}
}

func clampInterval(seconds int) int {
	if seconds < 5 {
		return 5
	}
	if seconds > 300 {
		return 300
	}
	return seconds
}

func detectCPUModel() string {
	data, err := os.ReadFile("/proc/cpuinfo")
	if err != nil {
		return "unknown"
	}
	for _, line := range strings.Split(string(data), "\n") {
		if strings.HasPrefix(strings.ToLower(line), "model name") {
			parts := strings.SplitN(line, ":", 2)
			if len(parts) != 2 {
				continue
			}
			value := strings.TrimSpace(parts[1])
			if value != "" {
				return value
			}
		}
	}
	return "unknown"
}

func detectTotalRAM() string {
	data, err := os.ReadFile("/proc/meminfo")
	if err != nil {
		return ""
	}
	for _, line := range strings.Split(string(data), "\n") {
		if !strings.HasPrefix(line, "MemTotal:") {
			continue
		}
		fields := strings.Fields(line)
		if len(fields) < 2 {
			continue
		}
		kb, parseErr := strconv.ParseInt(fields[1], 10, 64)
		if parseErr != nil || kb <= 0 {
			continue
		}
		return formatIECBytes(kb * 1024)
	}
	return ""
}

func formatIECBytes(bytes int64) string {
	if bytes <= 0 {
		return "0 B"
	}
	units := []string{"B", "KiB", "MiB", "GiB", "TiB", "PiB"}
	value := float64(bytes)
	unitIdx := 0
	for value >= 1024 && unitIdx < len(units)-1 {
		value /= 1024
		unitIdx++
	}
	if unitIdx == 0 {
		return fmt.Sprintf("%d %s", bytes, units[unitIdx])
	}
	return fmt.Sprintf("%.2f %s", value, units[unitIdx])
}
