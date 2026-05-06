package users

import "sync"

var (
	globalMonitorMu sync.RWMutex
	globalMonitor   *NodeMonitor
)

// RegisterGlobalNodeMonitor stores a global reference for API-triggered syncs.
func RegisterGlobalNodeMonitor(nm *NodeMonitor) {
	globalMonitorMu.Lock()
	globalMonitor = nm
	globalMonitorMu.Unlock()
}

// RequestNodeSync triggers an immediate sync if a monitor is registered.
func RequestNodeSync() {
	globalMonitorMu.RLock()
	nm := globalMonitor
	globalMonitorMu.RUnlock()
	if nm != nil {
		nm.RequestSync()
	}
}

// RequestNodeDeploy triggers config deploy to connected nodes.
func RequestNodeDeploy(restart bool, nodeUUIDs ...string) {
	globalMonitorMu.RLock()
	nm := globalMonitor
	globalMonitorMu.RUnlock()
	if nm != nil {
		nm.RequestDeploy(restart, nodeUUIDs...)
	}
}

// RequestSRSDeploy triggers SRS lists sync to connected nodes.
func RequestSRSDeploy() {
	globalMonitorMu.RLock()
	nm := globalMonitor
	globalMonitorMu.RUnlock()
	if nm != nil {
		nm.RequestSRSDeploy()
	}
}

// GetNodeMetricsSnapshot returns current per-node traffic metrics snapshots.
func GetNodeMetricsSnapshot() map[string]NodeMetricsSnapshot {
	globalMonitorMu.RLock()
	nm := globalMonitor
	globalMonitorMu.RUnlock()
	if nm == nil {
		return map[string]NodeMetricsSnapshot{}
	}
	return nm.SnapshotNodeMetrics()
}
