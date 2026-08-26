package server

import (
	"io"
	"time"

	"exodus-node/proto"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const defaultStreamInterval = 15 * time.Second

// StreamNodeData handles bidirectional streaming for node data.
func (s *NodeServer) StreamNodeData(stream proto.NodeService_StreamNodeDataServer) error {
	remoteAddr := ExtractClientIP(stream.Context())
	log := s.Cfg.LoggerFor("StreamService")
	log.Info("Panel stream connected successfully", "remote_addr", remoteAddr)
	defer func() {
		log.Warn("Panel stream disconnected", "remote_addr", remoteAddr)
	}()

	reqCh := make(chan *proto.NodeDataRequest)
	recvErrCh := make(chan error, 1)
	sendErrCh := make(chan error, 1)
	updateIntervalCh := make(chan time.Duration, 1)

	go receiveStreamRequests(stream, reqCh, recvErrCh)
	go s.sendStreamStats(stream, sendErrCh, updateIntervalCh)

	for {
		select {
		case <-stream.Context().Done():
			log.Debug("Panel stream context canceled", "remote_addr", remoteAddr, "error", stream.Context().Err())
			return stream.Context().Err()
		case err := <-sendErrCh:
			log.Error("Failed to send stream data", "remote_addr", remoteAddr, "error", err)
			return err
		case err := <-recvErrCh:
			if err == io.EOF {
				log.Info("Panel stream closed by panel", "remote_addr", remoteAddr)
				return nil
			}
			log.Warn("Panel stream receive failed", "remote_addr", remoteAddr, "error", err)
			return err
		case req, ok := <-reqCh:
			if !ok || req == nil {
				continue
			}
			if err := s.handleStreamRequest(req, updateIntervalCh); err != nil {
				return err
			}
		}
	}
}

func receiveStreamRequests(stream proto.NodeService_StreamNodeDataServer, reqCh chan<- *proto.NodeDataRequest, recvErrCh chan<- error) {
	defer close(reqCh)
	for {
		req, err := stream.Recv()
		if err != nil {
			recvErrCh <- err
			return
		}
		select {
		case reqCh <- req:
		case <-stream.Context().Done():
			return
		}
	}
}

func (s *NodeServer) sendStreamStats(
	stream proto.NodeService_StreamNodeDataServer,
	sendErrCh chan<- error,
	updateIntervalCh <-chan time.Duration,
) {
	ticker := time.NewTicker(defaultStreamInterval)
	defer ticker.Stop()

	for {
		stats, err := s.GetApiStats(stream.Context(), &proto.GetApiStatsRequest{})
		if err != nil {
			// Local stats/core collection errors must not close the stream. The stream is
			// the recovery channel used by the panel to keep talking to the node and push
			// a corrected config. Only real stream Send/Recv errors close it.
			s.Cfg.LoggerFor("StatsService").Warn("Failed to collect stats for stream; sending degraded response", "error", err)
			stats = &proto.GetApiStatsResponse{Stats: []*proto.Stat{
				{Name: "core_status", Value: "error"},
				{Name: "core_error", Value: err.Error()},
			}}
		}

		if err := stream.Send(&proto.NodeDataResponse{
			Response: &proto.NodeDataResponse_Stats{Stats: stats},
		}); err != nil {
			sendErrCh <- err
			return
		}

		for {
			select {
			case <-stream.Context().Done():
				return
			case next := <-updateIntervalCh:
				ticker.Reset(next)
				continue
			case <-ticker.C:
			}
			break
		}
	}
}

func (s *NodeServer) handleStreamRequest(req *proto.NodeDataRequest, updateIntervalCh chan time.Duration) error {
	switch req.Request.(type) {
	case *proto.NodeDataRequest_Config:
		interval := req.GetConfig().IntervalSeconds
		if interval <= 0 {
			return status.Errorf(codes.InvalidArgument, "invalid interval: %d", interval)
		}
		next := time.Duration(interval) * time.Second
		select {
		case updateIntervalCh <- next:
		default:
			<-updateIntervalCh
			updateIntervalCh <- next
		}
		s.Cfg.LoggerFor("StatsService").Debug("Stream interval updated", "interval_seconds", interval)
		return nil
	case *proto.NodeDataRequest_ListUsers:
		return status.Error(codes.Unimplemented, statsOnlyMessage)
	default:
		return status.Errorf(codes.InvalidArgument, "unknown request type")
	}
}
