package server

import (
	"context"
	"net"
	"strings"

	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"
)

// ExtractClientIP retrieves the real client IP from incoming gRPC metadata headers
// (e.g. CF-Connecting-IP, X-Real-IP, X-Forwarded-For) or falls back to the peer address.
func ExtractClientIP(ctx context.Context) string {
	if md, ok := metadata.FromIncomingContext(ctx); ok {
		for _, key := range []string{"cf-connecting-ip", "x-real-ip", "true-client-ip", "x-client-ip"} {
			if vals := md.Get(key); len(vals) > 0 {
				val := strings.TrimSpace(vals[0])
				if val != "" {
					return val
				}
			}
		}
		if vals := md.Get("x-forwarded-for"); len(vals) > 0 {
			raw := strings.TrimSpace(vals[0])
			if raw != "" {
				parts := strings.Split(raw, ",")
				for _, p := range parts {
					cleaned := strings.TrimSpace(p)
					if cleaned != "" {
						return cleaned
					}
				}
			}
		}
	}

	if p, ok := peer.FromContext(ctx); ok && p.Addr != nil {
		addrStr := p.Addr.String()
		if host, _, err := net.SplitHostPort(addrStr); err == nil && host != "" {
			return host
		}
		return addrStr
	}
	return "unknown"
}
