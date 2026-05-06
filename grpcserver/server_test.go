package grpcserver

import "testing"

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
