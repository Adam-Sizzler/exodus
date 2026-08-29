package server

import (
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"net/netip"
	"os"
	"sort"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/google/nftables"
	"github.com/google/nftables/expr"
	"golang.org/x/sys/unix"
)

const (
	linuxCAPNetAdmin = 12

	nftIPv4TableName = "exodus"
	nftIPv6TableName = "exodus6"

	nftIngressIPSet  = "ingress-filter-ip"
	nftEgressIPSet   = "egress-filter-ip"
	nftEgressPortSet = "egress-filter-port"
)

type NftIngressFilterPayload struct {
	Enabled     bool     `json:"enabled"`
	BlockedIPs  []string `json:"blocked_ips"`
	BlockedASNs []int    `json:"blocked_asns"`
}

type NftEgressFilterPayload struct {
	Enabled      bool     `json:"enabled"`
	BlockedIPs   []string `json:"blocked_ips"`
	BlockedPorts []int    `json:"blocked_ports"`
	BlockedASNs  []int    `json:"blocked_asns"`
}

type nftExecutorCommand struct {
	Command string          `json:"command"`
	IPs     json.RawMessage `json:"ips"`
}

type BlockedIPEntry struct {
	IP      string `json:"ip"`
	Timeout int    `json:"timeout"`
}

type nftIPElement struct {
	Value     string
	IsIPv6    bool
	Key       []byte
	KeyEnd    []byte
	IsNetwork bool
}

type nftTableSpec struct {
	Family        nftables.TableFamily
	Name          string
	IPType        nftables.SetDatatype
	AddrLen       uint32
	SrcAddrOffset uint32
	DstAddrOffset uint32
	IngressSet    string
	EgressSet     string
	EgressPortSet string
}

var nftTableSpecs = []nftTableSpec{
	{
		Family:        nftables.TableFamilyIPv4,
		Name:          nftIPv4TableName,
		IPType:        nftables.TypeIPAddr,
		AddrLen:       4,
		SrcAddrOffset: 12,
		DstAddrOffset: 16,
		IngressSet:    nftIngressIPSet,
		EgressSet:     nftEgressIPSet,
		EgressPortSet: nftEgressPortSet,
	},
	{
		Family:        nftables.TableFamilyIPv6,
		Name:          nftIPv6TableName,
		IPType:        nftables.TypeIP6Addr,
		AddrLen:       16,
		SrcAddrOffset: 8,
		DstAddrOffset: 24,
		IngressSet:    nftIngressIPSet + "6",
		EgressSet:     nftEgressIPSet + "6",
		EgressPortSet: nftEgressPortSet + "6",
	},
}

func applyNftablesModule(modules DeployModulesPayload, asnService *AsnLmdbService) error {
	if !hasCapNetAdmin() {
		return nil
	}

	if err := recreateNftablesTables(); err != nil {
		if isNftablesUnavailableError(err) {
			return nil
		}
		return err
	}

	if modules.IngressFilter.Enabled {
		blockedIPs := modules.IngressFilter.BlockedIPs
		if asnService != nil && len(modules.IngressFilter.BlockedASNs) > 0 {
			asnV4, asnV6 := asnService.ResolveASNs(modules.IngressFilter.BlockedASNs)
			blockedIPs = append(append([]string(nil), blockedIPs...), asnV4...)
			blockedIPs = append(blockedIPs, asnV6...)
		}
		if err := syncNftIPSet(nftIngressIPSet, blockedIPs, 0); err != nil {
			return fmt.Errorf("sync ingress filter: %w", err)
		}
	}

	if modules.EgressFilter.Enabled {
		blockedIPs := modules.EgressFilter.BlockedIPs
		if asnService != nil && len(modules.EgressFilter.BlockedASNs) > 0 {
			asnV4, asnV6 := asnService.ResolveASNs(modules.EgressFilter.BlockedASNs)
			blockedIPs = append(append([]string(nil), blockedIPs...), asnV4...)
			blockedIPs = append(blockedIPs, asnV6...)
		}
		if err := syncNftIPSet(nftEgressIPSet, blockedIPs, 0); err != nil {
			return fmt.Errorf("sync egress filter IPs: %w", err)
		}
		if err := syncNftPortSet(nftEgressPortSet, modules.EgressFilter.BlockedPorts); err != nil {
			return fmt.Errorf("sync egress filter ports: %w", err)
		}
	}

	return nil
}

func parseExecutorIPs(raw json.RawMessage) ([]BlockedIPEntry, error) {
	if len(raw) == 0 || string(raw) == "null" {
		return nil, nil
	}

	// Try parsing as array of objects [{ip, timeout}]
	var objEntries []BlockedIPEntry
	if err := json.Unmarshal(raw, &objEntries); err == nil && len(objEntries) > 0 && objEntries[0].IP != "" {
		return objEntries, nil
	}

	// Try parsing as array of strings ["1.2.3.4"]
	var strEntries []string
	if err := json.Unmarshal(raw, &strEntries); err == nil {
		res := make([]BlockedIPEntry, len(strEntries))
		for i, s := range strEntries {
			res[i] = BlockedIPEntry{IP: s}
		}
		return res, nil
	}

	// Try generic array
	var genericEntries []any
	if err := json.Unmarshal(raw, &genericEntries); err != nil {
		return nil, err
	}
	var res []BlockedIPEntry
	for _, item := range genericEntries {
		switch v := item.(type) {
		case string:
			res = append(res, BlockedIPEntry{IP: v})
		case map[string]any:
			ipStr, _ := v["ip"].(string)
			timeout := 0
			if t, ok := v["timeout"].(float64); ok {
				timeout = int(t)
			}
			res = append(res, BlockedIPEntry{IP: ipStr, Timeout: timeout})
		}
	}
	return res, nil
}

func ExecuteNodePluginCommand(raw json.RawMessage) (bool, error) {
	var command nftExecutorCommand
	if err := json.Unmarshal(raw, &command); err != nil {
		return false, fmt.Errorf("invalid node plugin executor payload: %w", err)
	}

	switch strings.TrimSpace(command.Command) {
	case "recreateTables":
		return executeNftCommand(recreateNftablesTables)
	case "blockIps":
		ips, err := parseExecutorIPs(command.IPs)
		if err != nil {
			return false, fmt.Errorf("invalid blockIps payload: %w", err)
		}
		return executeNftCommand(func() error {
			return blockNftIPs(ips)
		})
	case "unblockIps":
		ips, err := parseExecutorIPs(command.IPs)
		if err != nil {
			return false, fmt.Errorf("invalid unblockIps payload: %w", err)
		}
		return executeNftCommand(func() error {
			return unblockNftIPs(ips)
		})
	default:
		return false, fmt.Errorf("unsupported node plugin executor command: %s", command.Command)
	}
}

func blockNftIPs(items []BlockedIPEntry) error {
	for _, item := range items {
		ipStr := strings.TrimSpace(item.IP)
		if ipStr == "" {
			continue
		}
		ipElem, err := normalizeNftIP(ipStr)
		if err != nil {
			return err
		}
		if err := withNftConn(func(conn *nftables.Conn) error {
			for _, spec := range nftTableSpecs {
				if ipElem.IsIPv6 != (spec.Family == nftables.TableFamilyIPv6) {
					continue
				}
				setName := ipSetNameForSpec(nftIngressIPSet, spec)
				set := &nftables.Set{
					Table:      &nftables.Table{Family: spec.Family, Name: spec.Name},
					Name:       setName,
					HasTimeout: true,
				}
				elements := nftSetElementsForIP(ipElem, item.Timeout)
				if err := conn.SetAddElements(set, elements); err != nil {
					return fmt.Errorf("add IP to set %s/%s: %w", spec.Name, setName, err)
				}
			}
			return nil
		}); err != nil {
			return err
		}
	}
	return nil
}

func unblockNftIPs(items []BlockedIPEntry) error {
	for _, item := range items {
		ipStr := strings.TrimSpace(item.IP)
		if ipStr == "" {
			continue
		}
		ipElem, err := normalizeNftIP(ipStr)
		if err != nil {
			return err
		}
		if err := withNftConn(func(conn *nftables.Conn) error {
			for _, spec := range nftTableSpecs {
				if ipElem.IsIPv6 != (spec.Family == nftables.TableFamilyIPv6) {
					continue
				}
				setName := ipSetNameForSpec(nftIngressIPSet, spec)
				set := &nftables.Set{
					Table:      &nftables.Table{Family: spec.Family, Name: spec.Name},
					Name:       setName,
					HasTimeout: true,
				}
				elements := nftSetElementsForIP(ipElem, 0)
				if err := conn.SetDeleteElements(set, elements); err != nil && !isENOENT(err) {
					return fmt.Errorf("delete IP from set %s/%s: %w", spec.Name, setName, err)
				}
			}
			return nil
		}); err != nil {
			return err
		}
	}
	return nil
}

func executeNftCommand(fn func() error) (bool, error) {
	if !hasCapNetAdmin() {
		return false, nil
	}
	if err := ensureNftablesTables(); err != nil {
		if isNftablesUnavailableError(err) {
			return false, nil
		}
		return false, err
	}
	if err := fn(); err != nil {
		if isNftablesUnavailableError(err) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func recreateNftablesTables() error {
	for _, spec := range nftTableSpecs {
		conn, err := nftables.New()
		if err != nil {
			return fmt.Errorf("create nftables netlink connection: %w", err)
		}
		conn.DelTable(&nftables.Table{Family: spec.Family, Name: spec.Name})
		if err := conn.Flush(); err != nil && !isENOENT(err) {
			return fmt.Errorf("delete old nftables table %s: %w", spec.Name, err)
		}
	}

	conn, err := nftables.New()
	if err != nil {
		return fmt.Errorf("create nftables netlink connection: %w", err)
	}
	for _, spec := range nftTableSpecs {
		if err := addNftTableLayout(conn, spec); err != nil {
			return err
		}
	}
	if err := conn.Flush(); err != nil {
		return fmt.Errorf("create nftables tables: %w", err)
	}
	return nil
}

func ensureNftablesTables() error {
	conn, err := nftables.New()
	if err != nil {
		return fmt.Errorf("create nftables netlink connection: %w", err)
	}
	for _, spec := range nftTableSpecs {
		if _, err := conn.ListTableOfFamily(spec.Name, spec.Family); err != nil {
			return recreateNftablesTables()
		}
	}
	return nil
}

var (
	nftablesLogging            = true
	nftablesAcceptReplyTraffic = false
)

// ConfigureNftables updates global nftables settings for logging and reply traffic.
func ConfigureNftables(logging bool, acceptReplyTraffic bool) {
	nftablesLogging = logging
	nftablesAcceptReplyTraffic = acceptReplyTraffic
}

func addNftTableLayout(conn *nftables.Conn, spec nftTableSpec) error {
	table := conn.AddTable(&nftables.Table{Family: spec.Family, Name: spec.Name})
	policy := nftables.ChainPolicyAccept
	chains := map[string]*nftables.Chain{
		"input": conn.AddChain(&nftables.Chain{
			Name:     "input",
			Table:    table,
			Type:     nftables.ChainTypeFilter,
			Hooknum:  nftables.ChainHookInput,
			Priority: nftables.ChainPriorityFilter,
			Policy:   &policy,
		}),
		"forward": conn.AddChain(&nftables.Chain{
			Name:     "forward",
			Table:    table,
			Type:     nftables.ChainTypeFilter,
			Hooknum:  nftables.ChainHookForward,
			Priority: nftables.ChainPriorityFilter,
			Policy:   &policy,
		}),
		"output": conn.AddChain(&nftables.Chain{
			Name:     "output",
			Table:    table,
			Type:     nftables.ChainTypeFilter,
			Hooknum:  nftables.ChainHookOutput,
			Priority: nftables.ChainPriorityFilter,
			Policy:   &policy,
		}),
	}

	for _, setName := range []string{spec.IngressSet, spec.EgressSet} {
		if err := conn.AddSet(&nftables.Set{
			Table:      table,
			Name:       setName,
			KeyType:    spec.IPType,
			Interval:   true,
			AutoMerge:  true,
			HasTimeout: true,
		}, nil); err != nil {
			return fmt.Errorf("add nftables set %s/%s: %w", spec.Name, setName, err)
		}
	}
	if err := conn.AddSet(&nftables.Set{
		Table:     table,
		Name:      spec.EgressPortSet,
		KeyType:   nftables.TypeInetService,
		Interval:  true,
		AutoMerge: true,
	}, nil); err != nil {
		return fmt.Errorf("add nftables set %s/%s: %w", spec.Name, spec.EgressPortSet, err)
	}

	if nftablesAcceptReplyTraffic {
		addNftConntrackAcceptRule(conn, table, chains["input"])
		addNftConntrackAcceptRule(conn, table, chains["forward"])
	}

	addNftAddrDropRule(conn, table, chains["input"], spec, true, spec.IngressSet)
	addNftAddrDropRule(conn, table, chains["forward"], spec, true, spec.IngressSet)
	addNftAddrDropRule(conn, table, chains["forward"], spec, false, spec.EgressSet)
	addNftPortDropRule(conn, table, chains["forward"], spec.EgressPortSet, unix.IPPROTO_TCP)
	addNftPortDropRule(conn, table, chains["forward"], spec.EgressPortSet, unix.IPPROTO_UDP)
	addNftAddrDropRule(conn, table, chains["output"], spec, false, spec.EgressSet)
	addNftPortDropRule(conn, table, chains["output"], spec.EgressPortSet, unix.IPPROTO_TCP)
	addNftPortDropRule(conn, table, chains["output"], spec.EgressPortSet, unix.IPPROTO_UDP)

	return nil
}

func addNftConntrackAcceptRule(conn *nftables.Conn, table *nftables.Table, chain *nftables.Chain) {
	conn.AddRule(&nftables.Rule{
		Table: table,
		Chain: chain,
		Exprs: []expr.Any{
			&expr.Ct{Key: expr.CtKeySTATE, Register: 1},
			&expr.Bitwise{
				SourceRegister: 1,
				DestRegister:   1,
				Len:            4,
				Mask:           []byte{0, 0, 0, 6},
				Xor:            []byte{0, 0, 0, 0},
			},
			&expr.Cmp{
				Op:       expr.CmpOpNeq,
				Register: 1,
				Data:     []byte{0, 0, 0, 0},
			},
			&expr.Verdict{Kind: expr.VerdictAccept},
		},
	})
}

func addNftAddrDropRule(conn *nftables.Conn, table *nftables.Table, chain *nftables.Chain, spec nftTableSpec, source bool, setName string) {
	offset := spec.DstAddrOffset
	if source {
		offset = spec.SrcAddrOffset
	}
	exprs := []expr.Any{
		&expr.Payload{
			OperationType: expr.PayloadLoad,
			DestRegister:  1,
			Base:          expr.PayloadBaseNetworkHeader,
			Offset:        offset,
			Len:           spec.AddrLen,
		},
		&expr.Lookup{SourceRegister: 1, SetName: setName},
	}
	if nftablesLogging {
		exprs = append(exprs, &expr.Log{
			Key:  0x01,
			Data: []byte("[exodus-drop] "),
		})
	}
	exprs = append(exprs, &expr.Verdict{Kind: expr.VerdictDrop})

	conn.AddRule(&nftables.Rule{
		Table: table,
		Chain: chain,
		Exprs: exprs,
	})
}

func addNftPortDropRule(conn *nftables.Conn, table *nftables.Table, chain *nftables.Chain, setName string, proto byte) {
	exprs := []expr.Any{
		&expr.Payload{
			OperationType: expr.PayloadLoad,
			DestRegister:  1,
			Base:          expr.PayloadBaseTransportHeader,
			Offset:        2,
			Len:           2,
		},
		&expr.Meta{Key: expr.MetaKeyL4PROTO, Register: 2},
		&expr.Cmp{
			Op:       expr.CmpOpEq,
			Register: 2,
			Data:     []byte{proto},
		},
		&expr.Lookup{SourceRegister: 1, SetName: setName},
	}
	if nftablesLogging {
		exprs = append(exprs, &expr.Log{
			Key:  0x01,
			Data: []byte("[exodus-drop] "),
		})
	}
	exprs = append(exprs, &expr.Verdict{Kind: expr.VerdictDrop})

	conn.AddRule(&nftables.Rule{
		Table: table,
		Chain: chain,
		Exprs: exprs,
	})
}

func syncNftIPSet(baseSetName string, rawIPs []string, timeoutSeconds int) error {
	ips, err := normalizeNftIPs(rawIPs)
	if err != nil {
		return err
	}
	return withNftConn(func(conn *nftables.Conn) error {
		for _, spec := range nftTableSpecs {
			setName := ipSetNameForSpec(baseSetName, spec)
			set := &nftables.Set{
				Table:      &nftables.Table{Family: spec.Family, Name: spec.Name},
				Name:       setName,
				HasTimeout: true,
			}
			elements := make([]nftables.SetElement, 0, len(ips))
			for _, ip := range ips {
				if ip.IsIPv6 != (spec.Family == nftables.TableFamilyIPv6) {
					continue
				}
				elements = append(elements, nftSetElementsForIP(ip, timeoutSeconds)...)
			}
			if len(elements) == 0 {
				continue
			}
			if err := conn.SetAddElements(set, elements); err != nil {
				return fmt.Errorf("add IP elements to nft set %s/%s: %w", spec.Name, setName, err)
			}
		}
		return nil
	})
}

func syncNftPortSet(baseSetName string, rawPorts []int) error {
	ports, err := normalizeNftPorts(rawPorts)
	if err != nil {
		return err
	}
	return withNftConn(func(conn *nftables.Conn) error {
		for _, spec := range nftTableSpecs {
			setName := ipSetNameForSpec(baseSetName, spec)
			set := &nftables.Set{
				Table:     &nftables.Table{Family: spec.Family, Name: spec.Name},
				Name:      setName,
				Interval:  true,
				AutoMerge: true,
			}
			elements := make([]nftables.SetElement, 0, len(ports)*2)
			for _, port := range ports {
				start := uint16(port)
				end := uint32(start) + 1
				elements = append(elements,
					nftables.SetElement{Key: uint16Key(start)},
					nftables.SetElement{Key: uint16Key(uint16(end)), IntervalEnd: true},
				)
			}
			if len(elements) == 0 {
				continue
			}
			if err := conn.SetAddElements(set, elements); err != nil {
				return fmt.Errorf("add ports to nft set %s/%s: %w", spec.Name, set.Name, err)
			}
		}
		return nil
	})
}

func withNftConn(fn func(*nftables.Conn) error) error {
	if err := ensureNftablesTables(); err != nil {
		return err
	}
	conn, err := nftables.New()
	if err != nil {
		return fmt.Errorf("create nftables netlink connection: %w", err)
	}
	if err := fn(conn); err != nil {
		return err
	}
	if err := conn.Flush(); err != nil && !isEEXIST(err) {
		return fmt.Errorf("flush nftables transaction: %w", err)
	}
	return nil
}

func nftSetElementsForIP(ip nftIPElement, timeoutSeconds int) []nftables.SetElement {
	var timeout time.Duration
	if timeoutSeconds > 0 {
		timeout = time.Duration(timeoutSeconds) * time.Second
	}
	if len(ip.KeyEnd) > 0 {
		return []nftables.SetElement{
			{
				Key:     ip.Key,
				Timeout: timeout,
			},
			{
				Key:         ip.KeyEnd,
				IntervalEnd: true,
			},
		}
	}
	return []nftables.SetElement{
		{
			Key:     ip.Key,
			Timeout: timeout,
		},
	}
}

func ipSetNameForSpec(baseSetName string, spec nftTableSpec) string {
	if spec.Family != nftables.TableFamilyIPv6 {
		return baseSetName
	}
	switch baseSetName {
	case nftIngressIPSet, nftEgressIPSet, nftEgressPortSet:
		return baseSetName + "6"
	default:
		return baseSetName
	}
}

func normalizeNftIPs(rawIPs []string) ([]nftIPElement, error) {
	seen := make(map[string]struct{}, len(rawIPs))
	ips := make([]nftIPElement, 0, len(rawIPs))
	for _, rawIP := range rawIPs {
		ip, err := normalizeNftIP(rawIP)
		if err != nil {
			return nil, err
		}
		if _, ok := seen[ip.Value]; ok {
			continue
		}
		seen[ip.Value] = struct{}{}
		ips = append(ips, ip)
	}
	sort.Slice(ips, func(i, j int) bool { return ips[i].Value < ips[j].Value })
	return ips, nil
}

func normalizeNftIP(rawIP string) (nftIPElement, error) {
	value := strings.TrimSpace(rawIP)
	if value == "" {
		return nftIPElement{}, fmt.Errorf("empty IP address")
	}

	var prefix netip.Prefix
	if p, err := netip.ParsePrefix(value); err == nil {
		prefix = p.Masked()
	} else if addr, err := netip.ParseAddr(value); err == nil {
		prefix = netip.PrefixFrom(addr, addressBits(addr))
	} else {
		return nftIPElement{}, fmt.Errorf("invalid IP address or CIDR range %q", rawIP)
	}

	start, end := prefixRange(prefix)
	ip := nftIPElement{
		Value:     prefix.String(),
		IsIPv6:    prefix.Addr().Is6(),
		Key:       addrBytes(start),
		IsNetwork: prefix.Bits() != addressBits(prefix.Addr()),
	}
	if end.Next().IsValid() {
		ip.KeyEnd = addrBytes(end.Next())
	}
	return ip, nil
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

func prefixRange(prefix netip.Prefix) (netip.Addr, netip.Addr) {
	start := prefix.Masked().Addr()
	bits := addressBits(start)
	startBytes := addrBytes(start)
	endBytes := append([]byte(nil), startBytes...)
	for bit := prefix.Bits(); bit < bits; bit++ {
		byteIndex := bit / 8
		bitIndex := 7 - (bit % 8)
		endBytes[byteIndex] |= 1 << bitIndex
	}
	end := netip.AddrFrom16([16]byte{})
	if start.Is4() {
		var raw [4]byte
		copy(raw[:], endBytes)
		end = netip.AddrFrom4(raw)
	} else {
		var raw [16]byte
		copy(raw[:], endBytes)
		end = netip.AddrFrom16(raw)
	}
	return start, end
}

func addressBits(addr netip.Addr) int {
	if addr.Is4() {
		return 32
	}
	return 128
}

func addrBytes(addr netip.Addr) []byte {
	if addr.Is4() {
		raw := addr.As4()
		return append([]byte(nil), raw[:]...)
	}
	raw := addr.As16()
	return append([]byte(nil), raw[:]...)
}

func uint16Key(value uint16) []byte {
	out := make([]byte, 2)
	binary.BigEndian.PutUint16(out, value)
	return out
}

func hasCapNetAdmin() bool {
	status, err := os.ReadFile("/proc/self/status")
	if err != nil {
		return false
	}
	for _, line := range strings.Split(string(status), "\n") {
		if !strings.HasPrefix(line, "CapEff:") {
			continue
		}
		fields := strings.Fields(line)
		if len(fields) < 2 {
			return false
		}
		capEff, err := strconv.ParseUint(fields[1], 16, 64)
		if err != nil {
			return false
		}
		return capEff&(1<<linuxCAPNetAdmin) != 0
	}
	return false
}

func logNftablesAvailability(logger interface {
	Log(string, ...any)
	Warn(string, ...any)
}) {
	if logger == nil {
		return
	}
	if !hasCapNetAdmin() {
		logger.Warn("CAP_NET_ADMIN is not available.")
		logger.Warn("[PLUGIN] Ingress Filter: not available")
		logger.Warn("[PLUGIN] Egress Filter: not available")
		return
	}
	logger.Log("[OK] CAP_NET_ADMIN is available")
}

func isNftablesUnavailableError(err error) bool {
	return errors.Is(err, os.ErrPermission) ||
		errors.Is(err, syscall.EPERM) ||
		errors.Is(err, syscall.EACCES) ||
		errors.Is(err, syscall.ENOENT) ||
		errors.Is(err, unix.EOPNOTSUPP) ||
		strings.Contains(strings.ToLower(err.Error()), "operation not permitted") ||
		strings.Contains(strings.ToLower(err.Error()), "protocol not supported") ||
		strings.Contains(strings.ToLower(err.Error()), "not supported")
}

func isENOENT(err error) bool {
	return errors.Is(err, syscall.ENOENT) || strings.Contains(strings.ToLower(err.Error()), "no such file")
}

func isEEXIST(err error) bool {
	return errors.Is(err, syscall.EEXIST) ||
		errors.Is(err, os.ErrExist) ||
		strings.Contains(strings.ToLower(err.Error()), "file exists") ||
		strings.Contains(strings.ToLower(err.Error()), "already exists")
}
