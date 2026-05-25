package users

import (
	"encoding/json"
	"fmt"
	"sync"
)

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

// RequestNodeDeployWithForce triggers config deploy and forces core reload on target nodes.
func RequestNodeDeployWithForce(restart bool, forceRestart bool, nodeUUIDs ...string) {
	globalMonitorMu.RLock()
	nm := globalMonitor
	globalMonitorMu.RUnlock()
	if nm != nil {
		nm.RequestDeployWithForce(restart, forceRestart, nodeUUIDs...)
	}
}

// RequestNodePluginExecutor sends a node-plugin runtime command to connected nodes.
func RequestNodePluginExecutor(command json.RawMessage, nodeUUIDs ...string) error {
	globalMonitorMu.RLock()
	nm := globalMonitor
	globalMonitorMu.RUnlock()
	if nm == nil {
		return fmt.Errorf("node monitor is not ready")
	}
	return nm.ExecuteNodePluginCommand(nil, command, nodeUUIDs)
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
