package subscriptionnodes

import (
	"regexp"
	"time"
)

const exodusRealIPHeader = "x-exodus-real-ip"

const (
	bridgeOperationSubscriptionInfo    = "subscription_info"
	bridgeOperationSubscriptionContent = "subscription_content"
	bridgeOperationSubpageByShortUUID  = "subpage_config_for_short"
	bridgeOperationSubpageByUUID       = "subpage_config_by_uuid"
)

const (
	subNodeRuntimeStatVersion  = "sub_node_version"
	subNodeRuntimeStatUptime   = "sub_node_uptime"
	subNodeRuntimeStatCPUCount = "cpu_count"
	subNodeRuntimeStatCPUModel = "cpu_model"
	subNodeRuntimeStatTotalRAM = "total_ram"
)

var subNodeVersionPattern = regexp.MustCompile(`^[vV]?\d+\.\d+\.\d+(?:[-+][0-9A-Za-z\.-]+)?$`)

const (
	subNodeStreamInterval      = 20 * time.Second
	subNodeStreamIdleTimeout   = 75 * time.Second
	subNodeStreamWatchInterval = 5 * time.Second
)

type subpageConfigPushCommand struct {
	uuid        string
	config      []byte
	targetUUIDs []string
}

type dbSubNode struct {
	UUID              string
	Name              string
	Address           string
	Port              int
	APISchema         string
	APIPath           string
	GRPCAuthToken     string
	SubpageConfigUUID string
}

type SubNodeRuntimeSnapshot struct {
	SingboxVersion *string
	NodeVersion    *string
	SingboxUptime  string
	CPUCount       *int
	CPUModel       *string
	TotalRAM       *string
}
