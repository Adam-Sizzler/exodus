package subscription

import (
	"bytes"
	"encoding/json"
	"fmt"
	"math/rand"
	"sort"
	"strconv"
	"strings"

	"gopkg.in/yaml.v3"
)

func generateYAMLConfig(templateYAML []byte, hosts []SubscriptionHost, user SubscriptionUser) (string, error) {
	var root yaml.Node
	if len(templateYAML) > 0 {
		if err := yaml.Unmarshal(templateYAML, &root); err != nil {
			root = yaml.Node{}
		}
	}
	topLevelSpacing := extractYAMLTopLevelSpacing(templateYAML)
	config := ensureYAMLDocumentMappingNode(&root)
	proxiesNode := ensureYAMLMappingSequenceValue(config, "proxies")
	proxyNames := []string{}
	trailingSelectorProxyNames := []string{}
	for _, host := range hosts {
		proxy := buildMihomoProxy(host, user)
		if proxy == nil {
			continue
		}
		proxiesNode.Content = append(proxiesNode.Content, buildOrderedYAMLValueNode("proxy", proxy))
		if name, ok := proxy["name"].(string); ok && name != "" {
			proxyNames = append(proxyNames, name)
			trailingSelectorProxyNames = append(trailingSelectorProxyNames, name)
		}
	}
	groupsNode := ensureYAMLMappingSequenceValue(config, "proxy-groups")
	for _, group := range groupsNode.Content {
		if group == nil || group.Kind != yaml.MappingNode {
			continue
		}
		groupProxies := ensureYAMLMappingSequenceValue(group, "proxies")
		exodusNode := yamlMappingNode(group, "exodus")
		deleteYAMLMappingKey(group, "exodus")
		if exodusNode != nil {
			if v, ok := yamlMappingBool(exodusNode, "include-proxies"); ok && !v {
				continue
			}
			if v, ok := yamlMappingBool(exodusNode, "select-random-proxy"); ok && v {
				if len(proxyNames) > 0 {
					picked := proxyNames[rand.Intn(len(proxyNames))]
					existing := yamlSequenceStrings(groupProxies)
					setYAMLSequenceStrings(groupProxies, appendUniqueStrings(existing, picked))
				}
				continue
			}
			if v, ok := yamlMappingBool(exodusNode, "shuffle-proxies-order"); ok && v {
				shuffled := make([]string, len(proxyNames))
				copy(shuffled, proxyNames)
				rand.Shuffle(len(shuffled), func(i, j int) { shuffled[i], shuffled[j] = shuffled[j], shuffled[i] })
				existing := yamlSequenceStrings(groupProxies)
				setYAMLSequenceStrings(groupProxies, appendUniqueStrings(existing, shuffled...))
				continue
			}
			existing := yamlSequenceStrings(groupProxies)
			setYAMLSequenceStrings(groupProxies, appendUniqueStrings(existing, proxyNames...))
			continue
		}
		groupType := strings.ToLower(strings.TrimSpace(yamlMappingString(group, "type")))
		switch groupType {
		case "select":
			existingEntries := yamlSequenceStrings(groupProxies)
			hostNameSet := make(map[string]struct{}, len(proxyNames))
			for _, name := range proxyNames {
				hostNameSet[name] = struct{}{}
			}
			middleEntries := make([]string, 0, len(existingEntries))
			for _, entry := range existingEntries {
				if _, isHost := hostNameSet[entry]; !isHost {
					middleEntries = append(middleEntries, entry)
				}
			}
			finalEntries := make([]string, 0, len(middleEntries)+len(trailingSelectorProxyNames))
			finalEntries = append(finalEntries, middleEntries...)
			finalEntries = append(finalEntries, trailingSelectorProxyNames...)
			setYAMLSequenceStrings(groupProxies, finalEntries)
		case "url-test", "urltest":
			setYAMLSequenceStrings(groupProxies, proxyNames)
		default:
			finalEntries := appendUniqueStrings(yamlSequenceStrings(groupProxies), proxyNames...)
			setYAMLSequenceStrings(groupProxies, finalEntries)
		}
	}
	providersNode := yamlMappingNode(config, "proxy-providers")
	if providersNode != nil {
		for i := 0; i+1 < len(providersNode.Content); i += 2 {
			providerNode := providersNode.Content[i+1]
			if providerNode == nil || providerNode.Kind != yaml.MappingNode {
				continue
			}
			exodusNode := yamlMappingNode(providerNode, "exodus")
			if exodusNode == nil {
				continue
			}
			deleteYAMLMappingKey(providerNode, "exodus")
			v, ok := yamlMappingBool(exodusNode, "include-proxies")
			if !ok || !v {
				continue
			}
			payloadNode := &yaml.Node{Kind: yaml.SequenceNode, Tag: "!!seq"}
			for _, host := range hosts {
				proxy := buildMihomoProxy(host, user)
				if proxy == nil {
					continue
				}
				payloadNode.Content = append(payloadNode.Content, buildOrderedYAMLValueNode("proxy", proxy))
			}
			replaced := false
			for j := 0; j+1 < len(providerNode.Content); j += 2 {
				if providerNode.Content[j] != nil && providerNode.Content[j].Value == "payload" {
					providerNode.Content[j+1] = payloadNode
					replaced = true
					break
				}
			}
			if !replaced {
				providerNode.Content = append(providerNode.Content,
					newYAMLScalarNode("payload"), payloadNode)
			}
		}
	}
	var buf bytes.Buffer
	encoder := yaml.NewEncoder(&buf)
	encoder.SetIndent(2)
	err := encoder.Encode(config)
	closeErr := encoder.Close()
	if err != nil {
		return "", err
	}
	if closeErr != nil {
		return "", closeErr
	}
	rendered := applyYAMLTopLevelSpacing(buf.String(), topLevelSpacing)
	return rendered, nil
}

func ensureYAMLDocumentMappingNode(root *yaml.Node) *yaml.Node {
	if root.Kind != yaml.DocumentNode {
		*root = yaml.Node{Kind: yaml.DocumentNode}
	}
	if len(root.Content) == 0 || root.Content[0] == nil || root.Content[0].Kind != yaml.MappingNode {
		root.Content = []*yaml.Node{{Kind: yaml.MappingNode, Tag: "!!map"}}
	}
	return root.Content[0]
}

func ensureYAMLMappingSequenceValue(mapping *yaml.Node, key string) *yaml.Node {
	if mapping == nil || mapping.Kind != yaml.MappingNode {
		return &yaml.Node{Kind: yaml.SequenceNode, Tag: "!!seq"}
	}
	for i := 0; i+1 < len(mapping.Content); i += 2 {
		keyNode := mapping.Content[i]
		if keyNode == nil || keyNode.Value != key {
			continue
		}
		valueNode := mapping.Content[i+1]
		if valueNode == nil || valueNode.Kind != yaml.SequenceNode {
			valueNode = &yaml.Node{Kind: yaml.SequenceNode, Tag: "!!seq"}
			mapping.Content[i+1] = valueNode
		}
		return valueNode
	}
	valueNode := &yaml.Node{Kind: yaml.SequenceNode, Tag: "!!seq"}
	mapping.Content = append(mapping.Content, newYAMLScalarNode(key), valueNode)
	return valueNode
}

func yamlMappingString(mapping *yaml.Node, key string) string {
	if mapping == nil || mapping.Kind != yaml.MappingNode {
		return ""
	}
	for i := 0; i+1 < len(mapping.Content); i += 2 {
		keyNode := mapping.Content[i]
		if keyNode == nil || keyNode.Value != key {
			continue
		}
		valueNode := mapping.Content[i+1]
		if valueNode == nil || valueNode.Kind != yaml.ScalarNode {
			return ""
		}
		return strings.TrimSpace(valueNode.Value)
	}
	return ""
}

func yamlMappingBool(mapping *yaml.Node, key string) (bool, bool) {
	if mapping == nil || mapping.Kind != yaml.MappingNode {
		return false, false
	}
	for i := 0; i+1 < len(mapping.Content); i += 2 {
		if mapping.Content[i] == nil || mapping.Content[i].Value != key {
			continue
		}
		v := mapping.Content[i+1]
		if v == nil || v.Kind != yaml.ScalarNode {
			return false, false
		}
		val := strings.ToLower(strings.TrimSpace(v.Value))
		return val == "true", true
	}
	return false, false
}

func yamlMappingNode(mapping *yaml.Node, key string) *yaml.Node {
	if mapping == nil || mapping.Kind != yaml.MappingNode {
		return nil
	}
	for i := 0; i+1 < len(mapping.Content); i += 2 {
		if mapping.Content[i] == nil || mapping.Content[i].Value != key {
			continue
		}
		v := mapping.Content[i+1]
		if v != nil && v.Kind == yaml.MappingNode {
			return v
		}
		return nil
	}
	return nil
}

func deleteYAMLMappingKey(mapping *yaml.Node, key string) {
	if mapping == nil || mapping.Kind != yaml.MappingNode {
		return
	}
	for i := 0; i+1 < len(mapping.Content); i += 2 {
		if mapping.Content[i] == nil || mapping.Content[i].Value != key {
			continue
		}
		mapping.Content = append(mapping.Content[:i], mapping.Content[i+2:]...)
		return
	}
}

func yamlSequenceStrings(sequence *yaml.Node) []string {
	if sequence == nil || sequence.Kind != yaml.SequenceNode {
		return nil
	}
	values := make([]string, 0, len(sequence.Content))
	for _, item := range sequence.Content {
		if item == nil || item.Kind != yaml.ScalarNode {
			continue
		}
		value := strings.TrimSpace(item.Value)
		if value == "" {
			continue
		}
		values = append(values, value)
	}
	return values
}

func setYAMLSequenceStrings(sequence *yaml.Node, values []string) {
	if sequence == nil {
		return
	}
	sequence.Kind = yaml.SequenceNode
	sequence.Tag = "!!seq"
	sequence.Content = sequence.Content[:0]
	for _, value := range values {
		sequence.Content = append(sequence.Content, newYAMLScalarNode(value))
	}
}

func buildOrderedYAMLValueNode(parentKey string, value interface{}) *yaml.Node {
	switch v := value.(type) {
	case map[string]interface{}:
		return buildOrderedYAMLMappingNode(parentKey, v)
	case []string:
		node := &yaml.Node{Kind: yaml.SequenceNode, Tag: "!!seq"}
		for _, item := range v {
			node.Content = append(node.Content, newYAMLScalarNode(item))
		}
		return node
	case []interface{}:
		node := &yaml.Node{Kind: yaml.SequenceNode, Tag: "!!seq"}
		for _, item := range v {
			node.Content = append(node.Content, buildOrderedYAMLValueNode("", item))
		}
		return node
	case string:
		return newYAMLScalarNode(v)
	case bool:
		return &yaml.Node{Kind: yaml.ScalarNode, Tag: "!!bool", Value: strconv.FormatBool(v)}
	case int:
		return &yaml.Node{Kind: yaml.ScalarNode, Tag: "!!int", Value: strconv.Itoa(v)}
	case int64:
		return &yaml.Node{Kind: yaml.ScalarNode, Tag: "!!int", Value: strconv.FormatInt(v, 10)}
	case float64:
		return &yaml.Node{Kind: yaml.ScalarNode, Tag: "!!float", Value: strconv.FormatFloat(v, 'f', -1, 64)}
	case nil:
		return &yaml.Node{Kind: yaml.ScalarNode, Tag: "!!null"}
	default:
		return newYAMLScalarNode(fmt.Sprint(v))
	}
}

func buildOrderedYAMLMappingNode(parentKey string, values map[string]interface{}) *yaml.Node {
	node := &yaml.Node{Kind: yaml.MappingNode, Tag: "!!map"}
	usedKeys := make(map[string]struct{}, len(values))
	for _, key := range preferredYAMLKeyOrder(parentKey) {
		value, exists := values[key]
		if !exists {
			continue
		}
		appendYAMLMappingEntry(node, key, buildOrderedYAMLValueNode(key, value))
		usedKeys[key] = struct{}{}
	}
	remainingKeys := make([]string, 0, len(values))
	for key := range values {
		if _, exists := usedKeys[key]; exists {
			continue
		}
		remainingKeys = append(remainingKeys, key)
	}
	sort.Strings(remainingKeys)
	for _, key := range remainingKeys {
		appendYAMLMappingEntry(node, key, buildOrderedYAMLValueNode(key, values[key]))
	}
	return node
}

func preferredYAMLKeyOrder(parentKey string) []string {
	switch parentKey {
	case "proxy":
		return []string{
			"name", "type", "server", "port", "udp", "network",
			"uuid", "password", "cipher", "tls", "skip-cert-verify",
			"servername", "client-fingerprint", "alpn", "packet-encoding",
			"flow", "encryption", "ws-opts", "grpc-opts", "smux",
		}
	case "ws-opts":
		return []string{"path", "headers", "max-early-data", "early-data-header-name", "v2ray-http-upgrade", "v2ray-http-upgrade-fast-open"}
	case "headers":
		return []string{"Host"}
	case "grpc-opts":
		return []string{"grpc-service-name"}
	case "smux":
		return []string{"enabled", "protocol", "max-connections", "padding"}
	default:
		return nil
	}
}

func appendYAMLMappingEntry(node *yaml.Node, key string, value *yaml.Node) {
	node.Content = append(node.Content, newYAMLScalarNode(key), value)
}

func newYAMLScalarNode(value string) *yaml.Node {
	return &yaml.Node{Kind: yaml.ScalarNode, Tag: "!!str", Value: value}
}

func extractYAMLTopLevelSpacing(templateYAML []byte) map[string]bool {
	if len(templateYAML) == 0 {
		return nil
	}
	lines := strings.Split(strings.ReplaceAll(string(templateYAML), "\r\n", "\n"), "\n")
	spacing := map[string]bool{}
	previousBlank := false
	for _, line := range lines {
		if strings.TrimSpace(line) == "" {
			previousBlank = true
			continue
		}
		if key, ok := extractYAMLTopLevelKey(line); ok && previousBlank {
			spacing[key] = true
		}
		previousBlank = false
	}
	return spacing
}

func applyYAMLTopLevelSpacing(rendered string, spacing map[string]bool) string {
	if rendered == "" || len(spacing) == 0 {
		return rendered
	}
	hasTrailingNewline := strings.HasSuffix(rendered, "\n")
	lines := strings.Split(strings.TrimSuffix(rendered, "\n"), "\n")
	output := make([]string, 0, len(lines)+len(spacing))
	for index, line := range lines {
		if key, ok := extractYAMLTopLevelKey(line); ok && index > 0 && spacing[key] {
			if len(output) > 0 && strings.TrimSpace(output[len(output)-1]) != "" {
				output = append(output, "")
			}
		}
		output = append(output, line)
	}
	result := strings.Join(output, "\n")
	if hasTrailingNewline {
		result += "\n"
	}
	return result
}

func extractYAMLTopLevelKey(line string) (string, bool) {
	if line == "" {
		return "", false
	}
	if line[0] == ' ' || line[0] == '\t' || strings.HasPrefix(strings.TrimSpace(line), "- ") {
		return "", false
	}
	index := strings.IndexByte(line, ':')
	if index <= 0 {
		return "", false
	}
	key := strings.TrimSpace(line[:index])
	if key == "" {
		return "", false
	}
	return key, true
}

func buildMihomoProxy(host SubscriptionHost, user SubscriptionUser) map[string]interface{} {
	protocol := normalizedHostProtocol(host)
	if protocol == "" {
		return nil
	}
	name := host.Remark
	if name == "" {
		name = host.Address
	}
	proxy := map[string]interface{}{
		"name":   name,
		"type":   protocol,
		"server": host.Address,
		"port":   host.Port,
		"udp":    true,
	}
	if host.MihomoIPVersion != nil && strings.TrimSpace(*host.MihomoIPVersion) != "" {
		proxy["ip-version"] = strings.TrimSpace(*host.MihomoIPVersion)
	}
	network := "tcp"
	if host.InboundNetwork != nil && *host.InboundNetwork != "" {
		network = *host.InboundNetwork
	}
	switch protocol {
	case "vless":
		credential := effectiveProtocolCredential(host, user)
		if credential == "" {
			return nil
		}
		proxy["uuid"] = credential
	case "trojan":
		credential := effectiveProtocolCredential(host, user)
		if credential == "" {
			return nil
		}
		proxy["password"] = credential
	case "shadowsocks":
		credential := effectiveProtocolCredential(host, user)
		if credential == "" {
			return nil
		}
		proxy["type"] = "ss"
		proxy["password"] = credential
		method := extractShadowsocksMethod(host.InboundRaw)
		if method == "" {
			method = "aes-128-gcm"
		}
		proxy["cipher"] = method
	}
	security := "none"
	if host.InboundSecurity != nil && *host.InboundSecurity != "" {
		security = *host.InboundSecurity
	} else {
		switch strings.ToUpper(host.SecurityLayer) {
		case "TLS":
			security = "tls"
		case "NONE":
			security = "none"
		}
	}

	var mihomoSNI string
	if host.KeepSNIBlank {
		mihomoSNI = ""
	} else if host.OverrideSNIFromAddress {
		if host.SNI != nil && *host.SNI != "" {
			mihomoSNI = *host.SNI
		} else {
			mihomoSNI = host.Address
		}
	} else {
		nativeSNI := extractMihomoNativeSNI(host.InboundRaw)
		if nativeSNI != "" {
			mihomoSNI = nativeSNI
		} else {
			mihomoSNI = host.Address
		}
	}

	if security == "tls" {
		proxy["tls"] = true
		proxy["skip-cert-verify"] = host.AllowInsecure
		if mihomoSNI != "" {
			proxy["servername"] = mihomoSNI
		}
	}
	if host.Fingerprint != nil && *host.Fingerprint != "" {
		proxy["client-fingerprint"] = *host.Fingerprint
	}
	if host.ALPN != nil && *host.ALPN != "" {
		proxy["alpn"] = strings.Split(*host.ALPN, ",")
	}
	if network != "" {
		proxy["network"] = network
	}
	if network == "ws" {
		wsOpts := map[string]interface{}{}
		if host.Path != nil && *host.Path != "" {
			wsOpts["path"] = *host.Path
		}
		headers := map[string]interface{}{}
		if host.Host != nil && *host.Host != "" {
			headers["Host"] = *host.Host
		}
		if len(headers) > 0 {
			wsOpts["headers"] = headers
		}
		if len(wsOpts) > 0 {
			proxy["ws-opts"] = wsOpts
		}
	}
	if network == "grpc" {
		grpcOpts := map[string]interface{}{}
		if host.Path != nil && *host.Path != "" {
			grpcOpts["grpc-service-name"] = *host.Path
		}
		if len(grpcOpts) > 0 {
			proxy["grpc-opts"] = grpcOpts
		}
	}
	if host.ClashMuxParams != nil {
		if mux := parseMihomoMuxParams(*host.ClashMuxParams); mux != nil {
			proxy["smux"] = mux
		}
	}
	return proxy
}

func extractMihomoNativeSNI(inboundRaw []byte) string {
	if len(inboundRaw) == 0 {
		return ""
	}
	var raw map[string]interface{}
	if err := json.Unmarshal(inboundRaw, &raw); err != nil {
		return ""
	}
	streamSettings, _ := raw["streamSettings"].(map[string]interface{})
	if streamSettings == nil {
		return ""
	}
	security, _ := streamSettings["security"].(string)
	switch strings.ToLower(strings.TrimSpace(security)) {
	case "tls":
		tlsSettings, _ := streamSettings["tlsSettings"].(map[string]interface{})
		if tlsSettings != nil {
			if sn, ok := tlsSettings["serverName"].(string); ok {
				return strings.TrimSpace(sn)
			}
		}
	case "reality":
		realitySettings, _ := streamSettings["realitySettings"].(map[string]interface{})
		if realitySettings != nil {
			if sns, ok := realitySettings["serverNames"].([]interface{}); ok && len(sns) > 0 {
				if sn, ok := sns[0].(string); ok {
					return strings.TrimSpace(sn)
				}
			}
		}
	}
	return ""
}

func parseMihomoMuxParams(rawStr string) map[string]interface{} {
	rawStr = strings.TrimSpace(rawStr)
	if rawStr == "" {
		return nil
	}
	var result map[string]interface{}
	if err := json.Unmarshal([]byte(rawStr), &result); err == nil && result != nil {
		return normalizeMihomoMuxKeys(result)
	}
	if err := yaml.Unmarshal([]byte(rawStr), &result); err == nil && result != nil {
		return normalizeMihomoMuxKeys(result)
	}
	return nil
}

func normalizeMihomoMuxKeys(mux map[string]interface{}) map[string]interface{} {
	if mux == nil {
		return nil
	}
	normalized := make(map[string]interface{})
	for k, v := range mux {
		nk := strings.ReplaceAll(k, "_", "-")
		normalized[nk] = v
	}
	return normalized
}
