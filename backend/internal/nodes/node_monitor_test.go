package users

import (
	"testing"

	"exodus/internal/proto"
)

func TestExtractTrafficStatsDelta(t *testing.T) {
	stats := []*proto.Stat{
		{Name: "outbound>>>direct>>>traffic>>>uplink", Value: "10"},
		{Name: "outbound>>>direct>>>traffic>>>downlink", Value: "20"},
		{Name: "user>>>alice>>>traffic>>>uplink", Value: "3"},
		{Name: "user>>>alice>>>traffic>>>downlink", Value: "7"},
		{Name: "singbox_version", Value: "1.13.3"},
		{Name: "outbound>>>warp>>>traffic>>>uplink", Value: "bad"},
	}

	delta := extractTrafficStatsDelta(stats)

	if delta.TotalUploadBytes != 10 {
		t.Fatalf("unexpected total upload: got %d want %d", delta.TotalUploadBytes, 10)
	}
	if delta.TotalDownloadBytes != 20 {
		t.Fatalf("unexpected total download: got %d want %d", delta.TotalDownloadBytes, 20)
	}

	if got := delta.UserBytesByName["alice"]; got != 10 {
		t.Fatalf("unexpected alice bytes: got %d want %d", got, 10)
	}
	if _, ok := delta.UserBytesByName["bob"]; ok {
		t.Fatalf("bob should not have traffic bytes")
	}

	// Online users are derived only from non-zero traffic during the interval.
	if delta.UsersOnline != 1 {
		t.Fatalf("unexpected users_online: got %d want %d", delta.UsersOnline, 1)
	}
}

func TestApplyConsumptionMultiplier(t *testing.T) {
	if got := applyConsumptionMultiplier(100, 1_000_000_000); got != 100 {
		t.Fatalf("unexpected 1.0 multiplier result: got %d want %d", got, 100)
	}
	if got := applyConsumptionMultiplier(101, 500_000_000); got != 50 {
		t.Fatalf("unexpected 0.5 multiplier result: got %d want %d", got, 50)
	}
	if got := applyConsumptionMultiplier(100, 0); got != 0 {
		t.Fatalf("unexpected 0 multiplier result: got %d want %d", got, 0)
	}
}
