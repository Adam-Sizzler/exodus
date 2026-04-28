package grpcapi

import (
	"context"
	"io"
	"testing"
	"time"

	"github.com/exodus/subscription-page/backend/internal/proto"

	"google.golang.org/grpc/metadata"
)

type fakeNodeStream struct {
	ctx  context.Context
	sent chan *proto.NodeDataResponse
}

func (f *fakeNodeStream) SetHeader(metadata.MD) error { return nil }

func (f *fakeNodeStream) SendHeader(metadata.MD) error { return nil }

func (f *fakeNodeStream) SetTrailer(metadata.MD) {}

func (f *fakeNodeStream) Context() context.Context { return f.ctx }

func (f *fakeNodeStream) Send(resp *proto.NodeDataResponse) error {
	select {
	case f.sent <- resp:
	default:
	}
	return nil
}

func (f *fakeNodeStream) Recv() (*proto.NodeDataRequest, error) {
	<-f.ctx.Done()
	return nil, io.EOF
}

func (f *fakeNodeStream) SendMsg(any) error { return nil }

func (f *fakeNodeStream) RecvMsg(any) error { return nil }

func TestQueryPanelWaitsForStreamReconnect(t *testing.T) {
	t.Parallel()

	server := NewNodeServer("sub-test")
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	req := &proto.SubscriptionBridgeRequest{
		RequestId: "req-1",
		Operation: "subscription_info",
		ShortUuid: "short-uuid",
	}

	type result struct {
		resp *proto.SubscriptionBridgeResponse
		err  error
	}

	resultCh := make(chan result, 1)
	go func() {
		resp, err := server.QueryPanel(ctx, req)
		resultCh <- result{resp: resp, err: err}
	}()

	time.Sleep(20 * time.Millisecond)

	streamCtx, streamCancel := context.WithCancel(context.Background())
	defer streamCancel()

	stream := &fakeNodeStream{
		ctx:  streamCtx,
		sent: make(chan *proto.NodeDataResponse, 1),
	}
	server.attachStream(stream)
	defer server.detachStream(stream)

	select {
	case outbound := <-stream.sent:
		bridgeReq := outbound.GetSubscriptionRequest()
		if bridgeReq == nil {
			t.Fatal("expected outbound subscription bridge request")
		}
		if bridgeReq.GetRequestId() != req.GetRequestId() {
			t.Fatalf("unexpected request id: got %q want %q", bridgeReq.GetRequestId(), req.GetRequestId())
		}
	case <-time.After(250 * time.Millisecond):
		t.Fatal("expected QueryPanel to send request after stream attach")
	}

	server.handleBridgeResponse(&proto.SubscriptionBridgeResponse{
		RequestId:  req.GetRequestId(),
		StatusCode: 200,
		Payload:    []byte("ok"),
	})

	select {
	case result := <-resultCh:
		if result.err != nil {
			t.Fatalf("QueryPanel returned error: %v", result.err)
		}
		if string(result.resp.GetPayload()) != "ok" {
			t.Fatalf("unexpected payload: %q", string(result.resp.GetPayload()))
		}
	case <-time.After(250 * time.Millisecond):
		t.Fatal("expected QueryPanel to finish after bridge response")
	}
}
