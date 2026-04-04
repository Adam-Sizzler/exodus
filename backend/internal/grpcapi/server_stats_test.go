package grpcapi

import (
	"strconv"
	"testing"
)

func TestBuildStats_UsesSubNodeMetricNames(t *testing.T) {
	server := NewNodeServer("sub-test")
	stats := server.buildStats()

	values := make(map[string]string, len(stats))
	for _, stat := range stats {
		if stat == nil {
			continue
		}
		values[stat.GetName()] = stat.GetValue()
	}

	requiredKeys := []string{
		subNodeStatVersion,
		subNodeStatUptime,
		subNodeStatCPUCount,
		subNodeStatCPUModel,
		subNodeStatTotalRAM,
	}
	for _, key := range requiredKeys {
		if _, ok := values[key]; !ok {
			t.Fatalf("missing required stat key: %s", key)
		}
	}

	if _, hasSingbox := values["singbox_uptime"]; hasSingbox {
		t.Fatalf("unexpected legacy stat key: singbox_uptime")
	}
	if _, hasXray := values["xray_uptime"]; hasXray {
		t.Fatalf("unexpected legacy stat key: xray_uptime")
	}

	uptimeValue := values[subNodeStatUptime]
	parsedUptime, err := strconv.ParseInt(uptimeValue, 10, 64)
	if err != nil {
		t.Fatalf("uptime is not numeric: %q", uptimeValue)
	}
	if parsedUptime < 0 {
		t.Fatalf("uptime must be >= 0, got: %d", parsedUptime)
	}
}
