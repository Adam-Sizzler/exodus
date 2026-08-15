package subscription

/*
CROSS-CUTTING RULES / НЕЯВНЫЕ ЗАВИСИМОСТИ:
1. Использование host.Path:
   Значение `host.Path` используется как для WebSocket (`ws.path`), так и для gRPC (`grpc.service_name`).
   В билдерах (Mihomo, Xray, Singbox) при проверке пути gRPC нужно использовать `host.Path`.
2. Обработка Reality:
   Для совместимости с xray/mihomo билдерами, Reality принудительно схлопывается в `"tls"`.
*/

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/url"
	"strings"

	"gopkg.in/yaml.v3"
)

func buildRawHost(host SubscriptionHost) RawHost {
	protocol := ""
	if host.InboundType != nil {
		protocol = *host.InboundType
	}
	return RawHost{
		UUID:          host.UUID,
		Remark:        host.Remark,
		Address:       host.Address,
		Port:          host.Port,
		Protocol:      protocol,
		Network:       host.InboundNetwork,
		Security:      host.InboundSecurity,
		Path:          host.Path,
		SNI:           host.SNI,
		Host:          host.Host,
		ALPN:          host.ALPN,
		Fingerprint:   host.Fingerprint,
		IsDisabled:    host.IsDisabled,
		IsHidden:      host.IsHidden,
	}
}

func encodeRemark(remark string) string {
	return url.PathEscape(remark)
}

func buildSubscriptionLinks(hosts []SubscriptionHost, user SubscriptionUser) ([]string, map[string]string) {
	links := []string{}
	ssConfLinks := map[string]string{}
	for _, host := range hosts {
		link, protocol := buildHostLink(host, user)
		if link == "" {
			continue
		}
		links = append(links, link)
		if protocol == "shadowsocks" || protocol == "ss" {
			remark := host.Remark
			if remark == "" {
				remark = host.Address
			}
			encoded := base64.RawURLEncoding.EncodeToString([]byte(remark))
			domain := host.Address
			ssConfLinks[remark] = fmt.Sprintf("ssconf://%s/%s/ss/%s#%s", domain, user.ShortUUID, encoded, encodeRemark(remark))
		}
	}
	return links, ssConfLinks
}

func buildHostLink(host SubscriptionHost, user SubscriptionUser) (string, string) {
	protocol := normalizedHostProtocol(host)
	if protocol == "" {
		return "", ""
	}
	switch protocol {
	case "vless":
		return buildVlessLink(host, user), protocol
	case "trojan":
		return buildTrojanLink(host, user), protocol
	case "shadowsocks", "ss":
		return buildShadowsocksLink(host, user), protocol
	case "hysteria2", "hy2", "hysteria":
		return buildHysteria2Link(host, user), protocol
	case "anytls":
		return buildAnytlsLink(host, user), protocol
	case "tuic":
		return buildTuicLink(host, user), protocol
	case "vmess":
		return buildVmessLink(host, user), protocol
	default:
		return "", protocol
	}
}

func normalizedHostProtocol(host SubscriptionHost) string {
	if host.InboundType == nil {
		return ""
	}
	protocol := strings.ToLower(strings.TrimSpace(*host.InboundType))
	if protocol == "ss" {
		return "shadowsocks"
	}
	if protocol == "hy2" {
		return "hysteria2"
	}
	return protocol
}

func effectiveProtocolCredential(host SubscriptionHost, user SubscriptionUser) string {
	if host.OverrideProtocolCredential && host.ProtocolCredential != nil && *host.ProtocolCredential != "" {
		return *host.ProtocolCredential
	}
	protocol := normalizedHostProtocol(host)
	switch protocol {
	case "vless", "vmess":
		return user.VlessUUID
	case "trojan", "tuic":
		return user.TrojanPassword
	case "shadowsocks", "ss":
		return user.SSPassword
	case "hysteria2", "hy2", "hysteria":
		return user.Hysteria2Password
	case "anytls":
		return user.AnytlsPassword
	case "naive":
		return user.NaivePassword
	case "shadowtls":
		return user.ShadowtlsPassword
	default:
		return user.VlessUUID
	}
}

func effectiveNaiveUsername(user SubscriptionUser) string {
	return firstNonEmpty(user.Username, user.ShortUUID, user.UUID)
}

func buildVlessLink(host SubscriptionHost, user SubscriptionUser) string {
	credential := effectiveProtocolCredential(host, user)
	if credential == "" {
		return ""
	}
	params := url.Values{}
	params.Set("encryption", "none")
	applyTransportParams(&params, host)
	remark := host.Remark
	if remark == "" {
		remark = host.Address
	}
	return fmt.Sprintf("vless://%s@%s:%d?%s#%s", credential, host.Address, host.Port, params.Encode(), encodeRemark(remark))
}

func buildTrojanLink(host SubscriptionHost, user SubscriptionUser) string {
	credential := effectiveProtocolCredential(host, user)
	if credential == "" {
		return ""
	}
	params := url.Values{}
	applyTransportParams(&params, host)
	remark := host.Remark
	if remark == "" {
		remark = host.Address
	}
	return fmt.Sprintf("trojan://%s@%s:%d?%s#%s", url.QueryEscape(credential), host.Address, host.Port, params.Encode(), encodeRemark(remark))
}

func buildShadowsocksLink(host SubscriptionHost, user SubscriptionUser) string {
	method := extractShadowsocksMethod(host.InboundRaw)
	if method == "" {
		method = "aes-128-gcm"
	}
	credential := effectiveProtocolCredential(host, user)
	if credential == "" {
		return ""
	}
	creds := fmt.Sprintf("%s:%s", method, credential)
	encoded := base64.RawURLEncoding.EncodeToString([]byte(creds))
	remark := host.Remark
	if remark == "" {
		remark = host.Address
	}
	return fmt.Sprintf("ss://%s@%s:%d#%s", encoded, host.Address, host.Port, encodeRemark(remark))
}

func buildHysteria2Link(host SubscriptionHost, user SubscriptionUser) string {
	credential := effectiveProtocolCredential(host, user)
	if credential == "" {
		return ""
	}
	params := url.Values{}
	sni := ""
	if host.SNI != nil {
		sni = *host.SNI
	}
	if sni == "" && host.OverrideSNIFromAddress {
		sni = host.Address
	}
	if sni != "" {
		params.Set("sni", sni)
	}
	if host.ALPN != nil && *host.ALPN != "" {
		params.Set("alpn", *host.ALPN)
	}
	remark := host.Remark
	if remark == "" {
		remark = host.Address
	}
	query := params.Encode()
	if query != "" {
		return fmt.Sprintf("hysteria2://%s@%s:%d?%s#%s", url.QueryEscape(credential), host.Address, host.Port, query, encodeRemark(remark))
	}
	return fmt.Sprintf("hysteria2://%s@%s:%d#%s", url.QueryEscape(credential), host.Address, host.Port, encodeRemark(remark))
}

func buildAnytlsLink(host SubscriptionHost, user SubscriptionUser) string {
	credential := effectiveProtocolCredential(host, user)
	if credential == "" {
		return ""
	}
	params := url.Values{}
	applyTransportParams(&params, host)
	remark := host.Remark
	if remark == "" {
		remark = host.Address
	}
	query := params.Encode()
	if query != "" {
		return fmt.Sprintf("anytls://%s@%s:%d?%s#%s", url.QueryEscape(credential), host.Address, host.Port, query, encodeRemark(remark))
	}
	return fmt.Sprintf("anytls://%s@%s:%d#%s", url.QueryEscape(credential), host.Address, host.Port, encodeRemark(remark))
}

func buildTuicLink(host SubscriptionHost, user SubscriptionUser) string {
	credential := effectiveProtocolCredential(host, user)
	if credential == "" {
		return ""
	}
	uuidStr := user.VlessUUID
	params := url.Values{}
	sni := ""
	if host.SNI != nil {
		sni = *host.SNI
	}
	if sni == "" && host.OverrideSNIFromAddress {
		sni = host.Address
	}
	if sni != "" {
		params.Set("sni", sni)
	}
	if host.ALPN != nil && *host.ALPN != "" {
		params.Set("alpn", *host.ALPN)
	}
	remark := host.Remark
	if remark == "" {
		remark = host.Address
	}
	query := params.Encode()
	if query != "" {
		return fmt.Sprintf("tuic://%s:%s@%s:%d?%s#%s", uuidStr, credential, host.Address, host.Port, query, encodeRemark(remark))
	}
	return fmt.Sprintf("tuic://%s:%s@%s:%d#%s", uuidStr, credential, host.Address, host.Port, encodeRemark(remark))
}

func buildVmessLink(_ SubscriptionHost, _ SubscriptionUser) string {
	// VMess is not supported yet in this implementation.
	return ""
}

func applyTransportParams(params *url.Values, host SubscriptionHost) {
	network := "tcp"
	if host.InboundNetwork != nil && *host.InboundNetwork != "" {
		network = *host.InboundNetwork
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
	params.Set("type", network)
	if security != "" {
		params.Set("security", security)
	}
	sni := ""
	if host.SNI != nil {
		sni = *host.SNI
	}
	if sni == "" && host.OverrideSNIFromAddress {
		sni = host.Address
	}
	if sni != "" {
		params.Set("sni", sni)
	}
	if host.ALPN != nil && *host.ALPN != "" {
		params.Set("alpn", *host.ALPN)
	}
	if host.Fingerprint != nil && *host.Fingerprint != "" {
		params.Set("fp", *host.Fingerprint)
	}
	if host.Path != nil && *host.Path != "" {
		params.Set("path", *host.Path)
	}
	if host.Host != nil && *host.Host != "" {
		params.Set("host", *host.Host)
	}
}

func extractShadowsocksMethod(raw json.RawMessage) string {
	if len(raw) == 0 {
		return ""
	}
	var obj map[string]interface{}
	if err := json.Unmarshal(raw, &obj); err != nil {
		return ""
	}
	if settings, ok := obj["settings"].(map[string]interface{}); ok {
		if method, ok := settings["method"].(string); ok {
			return method
		}
	}
	if method, ok := obj["method"].(string); ok {
		return method
	}
	return ""
}

func parseJSONMapString(raw *string) map[string]interface{} {
	if raw == nil {
		return nil
	}
	trimmed := strings.TrimSpace(*raw)
	if trimmed == "" || trimmed == "null" {
		return nil
	}
	var parsed map[string]interface{}
	if err := json.Unmarshal([]byte(trimmed), &parsed); err == nil {
		if len(parsed) > 0 {
			return parsed
		}
		return nil
	}
	var yamlPayload string
	if err := json.Unmarshal([]byte(trimmed), &yamlPayload); err != nil {
		return nil
	}
	yamlPayload = strings.TrimSpace(yamlPayload)
	if yamlPayload == "" {
		return nil
	}
	var yamlParsed map[string]interface{}
	if err := yaml.Unmarshal([]byte(yamlPayload), &yamlParsed); err != nil {
		return nil
	}
	if len(yamlParsed) == 0 {
		return nil
	}
	return yamlParsed
}

func generateXrayJSONConfig(templateJSON []byte, hosts []SubscriptionHost, user SubscriptionUser) (string, error) {
	baseConfig := map[string]interface{}{}
	if len(templateJSON) > 0 {
		if err := json.Unmarshal(templateJSON, &baseConfig); err != nil {
			baseConfig = map[string]interface{}{}
		}
	}
	outbounds := []interface{}{}
	if existing, ok := baseConfig["outbounds"].([]interface{}); ok {
		outbounds = existing
	}
	for _, host := range hosts {
		outbound := buildXrayOutbound(host, user)
		if outbound != nil {
			outbounds = append(outbounds, outbound)
		}
	}
	baseConfig["outbounds"] = outbounds
	return marshalJSONWithTemplateTopLevelOrder(templateJSON, baseConfig)
}

func buildXrayOutbound(host SubscriptionHost, user SubscriptionUser) map[string]interface{} {
	protocol := normalizedHostProtocol(host)
	if protocol == "" {
		return nil
	}
	remark := host.Remark
	if remark == "" {
		remark = host.Address
	}
	network := "tcp"
	if host.InboundNetwork != nil && *host.InboundNetwork != "" {
		network = *host.InboundNetwork
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
	streamSettings := map[string]interface{}{
		"network":  network,
		"security": security,
	}
	sni := ""
	if host.SNI != nil {
		sni = *host.SNI
	}
	if sni == "" && host.OverrideSNIFromAddress {
		sni = host.Address
	}
	if security == "tls" {
		tlsSettings := map[string]interface{}{}
		if sni != "" {
			tlsSettings["serverName"] = sni
		}
		if host.ALPN != nil && *host.ALPN != "" {
			tlsSettings["alpn"] = strings.Split(*host.ALPN, ",")
		}
		if host.Fingerprint != nil && *host.Fingerprint != "" {
			tlsSettings["fingerprint"] = *host.Fingerprint
		}
		streamSettings["tlsSettings"] = tlsSettings
	}
	if network == "ws" {
		wsSettings := map[string]interface{}{}
		if host.Path != nil && *host.Path != "" {
			wsSettings["path"] = *host.Path
		}
		if host.Host != nil && *host.Host != "" {
			wsSettings["headers"] = map[string]interface{}{"Host": *host.Host}
		}
		streamSettings["wsSettings"] = wsSettings
	}
	if network == "grpc" {
		grpcSettings := map[string]interface{}{}
		if host.Path != nil && *host.Path != "" {
			grpcSettings["serviceName"] = *host.Path
		}
		streamSettings["grpcSettings"] = grpcSettings
	}
	if network == "xhttp" {
		xhttpSettings := readMap(readMap(parseInboundRaw(host.InboundRaw), "streamSettings"), "xhttpSettings")
		if host.Path != nil && *host.Path != "" {
			xhttpSettings["path"] = *host.Path
		}
		if host.Host != nil && *host.Host != "" {
			xhttpSettings["host"] = *host.Host
		}
		if extra := parseJSONMapString(host.XHTTPExtraParams); extra != nil {
			xhttpSettings["extra"] = extra
		}
		streamSettings["xhttpSettings"] = xhttpSettings
	}
	if sockopt := parseJSONMapString(host.SockoptParams); sockopt != nil {
		streamSettings["sockopt"] = sockopt
	}
	outbound := map[string]interface{}{
		"tag":            remark,
		"protocol":       protocol,
		"streamSettings": streamSettings,
	}
	switch protocol {
	case "vless":
		credential := effectiveProtocolCredential(host, user)
		if credential == "" {
			return nil
		}
		outbound["settings"] = map[string]interface{}{
			"vnext": []interface{}{
				map[string]interface{}{
					"address": host.Address,
					"port":    host.Port,
					"users": []interface{}{
						map[string]interface{}{
							"id":         credential,
							"encryption": "none",
						},
					},
				},
			},
		}
	case "trojan":
		credential := effectiveProtocolCredential(host, user)
		if credential == "" {
			return nil
		}
		outbound["settings"] = map[string]interface{}{
			"servers": []interface{}{
				map[string]interface{}{
					"address":  host.Address,
					"port":     host.Port,
					"password": credential,
				},
			},
		}
	case "shadowsocks":
		credential := effectiveProtocolCredential(host, user)
		if credential == "" {
			return nil
		}
		method := extractShadowsocksMethod(host.InboundRaw)
		if method == "" {
			method = "aes-128-gcm"
		}
		outbound["protocol"] = "shadowsocks"
		outbound["settings"] = map[string]interface{}{
			"servers": []interface{}{
				map[string]interface{}{
					"address":  host.Address,
					"port":     host.Port,
					"method":   method,
					"password": credential,
				},
			},
		}
	default:
		return nil
	}
	if mux := parseJSONMapString(host.MuxParams); mux != nil {
		outbound["mux"] = mux
	}
	return outbound
}
