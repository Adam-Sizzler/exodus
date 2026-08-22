package server

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"strings"
	"sync"
	"time"

	"exodus-node/geocheck"
)

const (
	taskOperationGeocheck = "geocheck"
	geocheckTimeout       = 45 * time.Second
)

var (
	geocheckMu         sync.Mutex
	geocheckInProgress bool
)

type GeocheckPayload struct {
	IP        string `json:"ip,omitempty"`
	Interface string `json:"interface,omitempty"`
}

func ExecuteGeocheck(ctx context.Context, payloadBytes []byte) (string, error) {
	var payload GeocheckPayload
	if len(payloadBytes) > 0 {
		_ = json.Unmarshal(payloadBytes, &payload)
	}

	ip := strings.TrimSpace(payload.IP)
	iface := strings.TrimSpace(payload.Interface)

	if ip != "" {
		if net.ParseIP(ip) == nil {
			return "", fmt.Errorf("geocheck: %q is not a valid IP address", ip)
		}
	}

	geocheckMu.Lock()
	if geocheckInProgress {
		geocheckMu.Unlock()
		return "", errors.New("geocheck: a run is already in progress")
	}
	geocheckInProgress = true
	geocheckMu.Unlock()

	defer func() {
		geocheckMu.Lock()
		geocheckInProgress = false
		geocheckMu.Unlock()
	}()

	execCtx, cancel := context.WithTimeout(ctx, geocheckTimeout)
	defer cancel()

	out, err := geocheck.Run(execCtx, iface, ip)
	if err != nil {
		if execCtx.Err() == context.DeadlineExceeded {
			return "", fmt.Errorf("geocheck exceeded timeout of %v", geocheckTimeout)
		}
		return "", fmt.Errorf("geocheck execution failed: %w", err)
	}

	stdoutStr := strings.TrimSpace(string(out))
	if stdoutStr == "" {
		return "", errors.New("geocheck returned empty output")
	}

	return stdoutStr, nil
}
