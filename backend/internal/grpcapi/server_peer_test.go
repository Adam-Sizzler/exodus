package grpcapi

import (
	"context"
	"net"
	"testing"

	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"
)

type dummyAddr string

func (d dummyAddr) Network() string { return "tcp" }
func (d dummyAddr) String() string  { return string(d) }

func TestStreamPeerAddress(t *testing.T) {
	t.Parallel()

	t.Run("empty context", func(t *testing.T) {
		got := streamPeerAddress(context.Background())
		if got != "unknown" {
			t.Fatalf("expected 'unknown', got %q", got)
		}
	})

	t.Run("peer addr with host:port", func(t *testing.T) {
		p := &peer.Peer{Addr: &net.TCPAddr{IP: net.ParseIP("172.18.0.2"), Port: 52998}}
		ctx := peer.NewContext(context.Background(), p)
		got := streamPeerAddress(ctx)
		if got != "172.18.0.2" {
			t.Fatalf("expected '172.18.0.2', got %q", got)
		}
	})

	t.Run("x-forwarded-for single", func(t *testing.T) {
		md := metadata.Pairs("x-forwarded-for", "203.0.113.195")
		ctx := metadata.NewIncomingContext(context.Background(), md)
		got := streamPeerAddress(ctx)
		if got != "203.0.113.195" {
			t.Fatalf("expected '203.0.113.195', got %q", got)
		}
	})

	t.Run("x-forwarded-for multiple", func(t *testing.T) {
		md := metadata.Pairs("x-forwarded-for", "203.0.113.195, 172.18.0.1")
		ctx := metadata.NewIncomingContext(context.Background(), md)
		got := streamPeerAddress(ctx)
		if got != "203.0.113.195" {
			t.Fatalf("expected '203.0.113.195', got %q", got)
		}
	})

	t.Run("x-real-ip header", func(t *testing.T) {
		md := metadata.Pairs("x-real-ip", "198.51.100.42")
		ctx := metadata.NewIncomingContext(context.Background(), md)
		got := streamPeerAddress(ctx)
		if got != "198.51.100.42" {
			t.Fatalf("expected '198.51.100.42', got %q", got)
		}
	})

	t.Run("cf-connecting-ip header takes precedence", func(t *testing.T) {
		md := metadata.Pairs(
			"cf-connecting-ip", "198.51.100.99",
			"x-forwarded-for", "10.0.0.1",
		)
		ctx := metadata.NewIncomingContext(context.Background(), md)
		got := streamPeerAddress(ctx)
		if got != "198.51.100.99" {
			t.Fatalf("expected '198.51.100.99', got %q", got)
		}
	})
}
