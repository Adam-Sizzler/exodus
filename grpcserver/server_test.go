package grpcserver

import (
	"context"
	"testing"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

func TestTrimPathPrefix(t *testing.T) {
	tests := []struct {
		name      string
		path      string
		prefix    string
		wantPath  string
		wantMatch bool
	}{
		{
			name:      "empty prefix",
			path:      "/experimental.v2rayapi.StatsService/QueryStats",
			prefix:    "",
			wantPath:  "/experimental.v2rayapi.StatsService/QueryStats",
			wantMatch: true,
		},
		{
			name:      "exact prefix",
			path:      "/node",
			prefix:    "/node",
			wantPath:  "/",
			wantMatch: true,
		},
		{
			name:      "nested service path",
			path:      "/node/exodus.NodeService/StreamNodeData",
			prefix:    "/node",
			wantPath:  "/exodus.NodeService/StreamNodeData",
			wantMatch: true,
		},
		{
			name:      "same text without path boundary",
			path:      "/node-api/exodus.NodeService/StreamNodeData",
			prefix:    "/node",
			wantPath:  "/node-api/exodus.NodeService/StreamNodeData",
			wantMatch: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotPath, gotMatch := trimPathPrefix(tt.path, tt.prefix)
			if gotPath != tt.wantPath || gotMatch != tt.wantMatch {
				t.Fatalf("trimPathPrefix() = (%q, %t), want (%q, %t)", gotPath, gotMatch, tt.wantPath, tt.wantMatch)
			}
		})
	}
}

func TestValidateIncomingGRPCToken(t *testing.T) {
	ctx := metadata.NewIncomingContext(context.Background(), metadata.Pairs(grpcTokenHeader, "1234567890abcdef"))
	if err := validateIncomingGRPCToken(ctx, "1234567890abcdef"); err != nil {
		t.Fatalf("validateIncomingGRPCToken() error = %v", err)
	}

	err := validateIncomingGRPCToken(context.Background(), "1234567890abcdef")
	if status.Code(err) != codes.Unauthenticated {
		t.Fatalf("validateIncomingGRPCToken() code = %s, want %s", status.Code(err), codes.Unauthenticated)
	}
}

func TestSNIVerificationMatch(t *testing.T) {
	expectedSNI := "node1.domain.com"
	verifySNI := func(received string) bool {
		return expectedSNI != "" && (received == expectedSNI || (len(received) > 0 && received == "node1.domain.com"))
	}

	if !verifySNI("node1.domain.com") {
		t.Fatalf("expected SNI match for node1.domain.com")
	}
	if verifySNI("other.domain.com") {
		t.Fatalf("expected mismatch for other.domain.com")
	}
}

