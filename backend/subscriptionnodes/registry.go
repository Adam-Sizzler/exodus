package subscriptionnodes

import "sync"

var (
	globalSubMonitorMu sync.RWMutex
	globalSubMonitor   *SubNodeMonitor
)

func RegisterGlobalSubNodeMonitor(monitor *SubNodeMonitor) {
	globalSubMonitorMu.Lock()
	globalSubMonitor = monitor
	globalSubMonitorMu.Unlock()
}

func RequestSubNodeSync() {
	globalSubMonitorMu.RLock()
	monitor := globalSubMonitor
	globalSubMonitorMu.RUnlock()
	if monitor != nil {
		monitor.RequestSync()
	}
}

func RequestSubNodeDeploy(nodeUUIDs ...string) {
	globalSubMonitorMu.RLock()
	monitor := globalSubMonitor
	globalSubMonitorMu.RUnlock()
	if monitor != nil {
		monitor.RequestDeploy(nodeUUIDs...)
	}
}

func RequestSubNodeSubpageConfigPush(uuid string, config []byte, nodeUUIDs ...string) {
	globalSubMonitorMu.RLock()
	monitor := globalSubMonitor
	globalSubMonitorMu.RUnlock()
	if monitor != nil {
		monitor.RequestSubpageConfigPush(uuid, config, nodeUUIDs...)
	}
}

func GetSubNodeRuntimeSnapshot(nodeName string) (SubNodeRuntimeSnapshot, bool) {
	globalSubMonitorMu.RLock()
	monitor := globalSubMonitor
	globalSubMonitorMu.RUnlock()
	if monitor == nil {
		return SubNodeRuntimeSnapshot{}, false
	}
	return monitor.RuntimeSnapshot(nodeName)
}
