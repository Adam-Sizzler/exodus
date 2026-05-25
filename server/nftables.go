package server

import (
	"encoding/json"
	"fmt"
	"net/netip"
	"os/exec"
	"sort"
	"strings"
)

const (
	nftTableFamily       = "inet"
	nftTableName         = "exodus"
	nftTorrentBlockerSet = "torrent-blocker"
	nftIngressIPSet      = "ingress-filter-ip"
	nftEgressIPSet       = "egress-filter-ip"
	nftEgressPortSet     = "egress-filter-port"
)

type NftIngressFilterPayload struct {
	Enabled    bool     `json:"enabled"`
	BlockedIPs []string `json:"blocked_ips"`
}

type NftEgressFilterPayload struct {
	Enabled      bool     `json:"enabled"`
	BlockedIPs   []string `json:"blocked_ips"`
	BlockedPorts []int    `json:"blocked_ports"`
}

type nftExecutorCommand struct {
	Command string `json:"command"`
	IPs     []struct {
		IP      string `json:"ip"`
		Timeout int    `json:"timeout"`
	} `json:"ips"`
}

func applyNftablesModule(modules DeployModulesPayload) error {
	if err := recreateNftablesTable(); err != nil {
		return err
	}

	if modules.IngressFilter.Enabled {
		if err := syncNftIPSet(nftIngressIPSet, modules.IngressFilter.BlockedIPs, 0); err != nil {
			return fmt.Errorf("sync ingress filter: %w", err)
		}
	}

	if modules.EgressFilter.Enabled {
		if err := syncNftIPSet(nftEgressIPSet, modules.EgressFilter.BlockedIPs, 0); err != nil {
			return fmt.Errorf("sync egress filter IPs: %w", err)
		}
		if err := syncNftPortSet(nftEgressPortSet, modules.EgressFilter.BlockedPorts); err != nil {
			return fmt.Errorf("sync egress filter ports: %w", err)
		}
	}

	return nil
}

func ExecuteNodePluginCommand(raw json.RawMessage) (bool, error) {
	var command nftExecutorCommand
	if err := json.Unmarshal(raw, &command); err != nil {
		return false, fmt.Errorf("invalid node plugin executor payload: %w", err)
	}

	switch strings.TrimSpace(command.Command) {
	case "recreateTables":
		return true, recreateNftablesTable()
	case "blockIps":
		if err := ensureNftablesTable(); err != nil {
			return false, err
		}
		for _, item := range command.IPs {
			if err := addNftIPElement(nftTorrentBlockerSet, item.IP, item.Timeout); err != nil {
				return false, err
			}
		}
		return true, nil
	case "unblockIps":
		if err := ensureNftablesTable(); err != nil {
			return false, err
		}
		ips := make([]string, 0, len(command.IPs))
		for _, item := range command.IPs {
			ips = append(ips, item.IP)
		}
		if err := removeNftIPs(nftTorrentBlockerSet, ips); err != nil {
			return false, err
		}
		if err := removeNftIPs(nftIngressIPSet, ips); err != nil {
			return false, err
		}
		return true, nil
	default:
		return false, fmt.Errorf("unsupported node plugin executor command: %s", command.Command)
	}
}

func recreateNftablesTable() error {
	_ = runNft("delete", "table", nftTableFamily, nftTableName)

	script := fmt.Sprintf(`
table %s %s {
	set %s {
		type ipv4_addr
		flags interval,timeout
	}

	set %s {
		type ipv4_addr
		flags interval,timeout
	}

	set %s {
		type ipv4_addr
		flags interval,timeout
	}

	set %s {
		type inet_service
	}

	chain input {
		type filter hook input priority filter; policy accept;
		ip saddr @%s drop
		ip saddr @%s drop
	}

	chain output {
		type filter hook output priority filter; policy accept;
		ip daddr @%s drop
		tcp dport @%s drop
		udp dport @%s drop
	}
}
`,
		nftTableFamily,
		nftTableName,
		nftTorrentBlockerSet,
		nftIngressIPSet,
		nftEgressIPSet,
		nftEgressPortSet,
		nftTorrentBlockerSet,
		nftIngressIPSet,
		nftEgressIPSet,
		nftEgressPortSet,
		nftEgressPortSet,
	)

	if err := runNftScript(script); err != nil {
		return fmt.Errorf("create nftables table %s %s: %w", nftTableFamily, nftTableName, err)
	}
	return nil
}

func ensureNftablesTable() error {
	if err := runNft("list", "table", nftTableFamily, nftTableName); err == nil {
		return nil
	}
	return recreateNftablesTable()
}

func syncNftIPSet(setName string, rawIPs []string, timeoutSeconds int) error {
	ips, err := normalizeNftIPs(rawIPs)
	if err != nil {
		return err
	}
	for _, ip := range ips {
		if err := addNftIPElement(setName, ip, timeoutSeconds); err != nil {
			return err
		}
	}
	return nil
}

func syncNftPortSet(setName string, rawPorts []int) error {
	ports, err := normalizeNftPorts(rawPorts)
	if err != nil {
		return err
	}
	for _, port := range ports {
		script := fmt.Sprintf("add element %s %s %s { %d }\n", nftTableFamily, nftTableName, setName, port)
		if err := runNftScript(script); err != nil {
			return fmt.Errorf("add port %d to nft set %s: %w", port, setName, err)
		}
	}
	return nil
}

func addNftIPElement(setName string, rawIP string, timeoutSeconds int) error {
	ip, err := normalizeNftIP(rawIP)
	if err != nil {
		return err
	}

	_ = removeNftIP(setName, ip)

	timeoutPart := ""
	if timeoutSeconds > 0 {
		timeoutPart = fmt.Sprintf(" timeout %ds", timeoutSeconds)
	}
	script := fmt.Sprintf("add element %s %s %s { %s%s }\n", nftTableFamily, nftTableName, setName, ip, timeoutPart)
	if err := runNftScript(script); err != nil {
		return fmt.Errorf("add IP %s to nft set %s: %w", ip, setName, err)
	}
	return nil
}

func removeNftIPs(setName string, rawIPs []string) error {
	ips, err := normalizeNftIPs(rawIPs)
	if err != nil {
		return err
	}
	for _, ip := range ips {
		_ = removeNftIP(setName, ip)
	}
	return nil
}

func removeNftIP(setName string, ip string) error {
	script := fmt.Sprintf("delete element %s %s %s { %s }\n", nftTableFamily, nftTableName, setName, ip)
	return runNftScript(script)
}

func normalizeNftIPs(rawIPs []string) ([]string, error) {
	seen := make(map[string]struct{}, len(rawIPs))
	ips := make([]string, 0, len(rawIPs))
	for _, rawIP := range rawIPs {
		ip, err := normalizeNftIP(rawIP)
		if err != nil {
			return nil, err
		}
		if _, ok := seen[ip]; ok {
			continue
		}
		seen[ip] = struct{}{}
		ips = append(ips, ip)
	}
	sort.Strings(ips)
	return ips, nil
}

func normalizeNftIP(rawIP string) (string, error) {
	value := strings.TrimSpace(rawIP)
	if prefix, err := netip.ParsePrefix(value); err == nil {
		prefix = prefix.Masked()
		if !prefix.Addr().Is4() {
			return "", fmt.Errorf("only IPv4 addresses and CIDR ranges are supported by nft set %s", nftTableName)
		}
		return prefix.String(), nil
	}

	addr, err := netip.ParseAddr(value)
	if err != nil {
		return "", fmt.Errorf("invalid IP address or CIDR range %q", rawIP)
	}
	if !addr.Is4() {
		return "", fmt.Errorf("only IPv4 addresses and CIDR ranges are supported by nft set %s", nftTableName)
	}
	return addr.String(), nil
}

func normalizeNftPorts(rawPorts []int) ([]int, error) {
	seen := make(map[int]struct{}, len(rawPorts))
	ports := make([]int, 0, len(rawPorts))
	for _, port := range rawPorts {
		if port < 1 || port > 65535 {
			return nil, fmt.Errorf("invalid egress port %d", port)
		}
		if _, ok := seen[port]; ok {
			continue
		}
		seen[port] = struct{}{}
		ports = append(ports, port)
	}
	sort.Ints(ports)
	return ports, nil
}

func runNft(args ...string) error {
	output, err := exec.Command("nft", args...).CombinedOutput()
	if err != nil {
		return fmt.Errorf("nft %s failed: %w: %s", strings.Join(args, " "), err, strings.TrimSpace(string(output)))
	}
	return nil
}

func runNftScript(script string) error {
	cmd := exec.Command("nft", "-f", "-")
	cmd.Stdin = strings.NewReader(script)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("nft script failed: %w: %s", err, strings.TrimSpace(string(output)))
	}
	return nil
}
