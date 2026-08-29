package server

import (
	"testing"
)

func TestConfigureNftables(t *testing.T) {
	ConfigureNftables(false, true)
	if nftablesLogging {
		t.Errorf("nftablesLogging = true, want false")
	}
	if !nftablesAcceptReplyTraffic {
		t.Errorf("nftablesAcceptReplyTraffic = false, want true")
	}

	// Reset to defaults
	ConfigureNftables(true, false)
	if !nftablesLogging {
		t.Errorf("nftablesLogging = false, want true")
	}
	if nftablesAcceptReplyTraffic {
		t.Errorf("nftablesAcceptReplyTraffic = true, want false")
	}
}
