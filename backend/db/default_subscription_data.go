package db

import (
	"context"
	"database/sql"
	"fmt"

	"v2ray-stat/backend/config"

	"github.com/google/uuid"
)

type defaultSubscriptionTemplate struct {
	templateType string
	templateYAML *string
	templateJSON *string
	name         string
	viewPosition int
}

var defaultSubscriptionTemplates = []defaultSubscriptionTemplate{
	{
		templateType: "XRAY_JSON",
		templateJSON: strPtr(`{"dns":{"servers":["1.1.1.1","1.0.0.1"],"queryStrategy":"UseIP"},"routing":{"rules":[{"type":"field","protocol":["bittorrent"],"outboundTag":"direct"}],"domainMatcher":"hybrid","domainStrategy":"IPIfNonMatch"},"inbounds":[{"tag":"socks","port":10808,"listen":"127.0.0.1","protocol":"socks","settings":{"udp":true,"auth":"noauth"},"sniffing":{"enabled":true,"routeOnly":false,"destOverride":["http","tls","quic"]}},{"tag":"http","port":10809,"listen":"127.0.0.1","protocol":"http","settings":{"allowTransparent":false},"sniffing":{"enabled":true,"routeOnly":false,"destOverride":["http","tls","quic"]}}],"outbounds":[{"tag":"direct","protocol":"freedom"},{"tag":"block","protocol":"blackhole"}]}`),
		name:         "Default",
		viewPosition: 1,
	},
	{
		templateType: "MIHOMO",
		templateYAML: strPtr(`mixed-port: 7890
socks-port: 7891
redir-port: 7892
allow-lan: true
mode: global
log-level: info
external-controller: 127.0.0.1:9090
dns:
  enable: true
  use-hosts: true
  enhanced-mode: fake-ip
  fake-ip-range: 198.18.0.1/16
  default-nameserver:
    - 1.1.1.1
    - 8.8.8.8
  nameserver:
    - 1.1.1.1
    - 8.8.8.8
  fake-ip-filter:
    - '*.lan'
    - stun.*.*.*
    - stun.*.*
    - time.windows.com
    - time.nist.gov
    - time.apple.com
    - time.asia.apple.com
    - '*.openwrt.pool.ntp.org'
    - pool.ntp.org
    - ntp.ubuntu.com
    - time1.apple.com
    - time2.apple.com
    - time3.apple.com
    - time4.apple.com
    - time5.apple.com
    - time6.apple.com
    - time7.apple.com
    - time1.google.com
    - time2.google.com
    - time3.google.com
    - time4.google.com
    - api.joox.com
    - joox.com
    - '*.xiami.com'
    - '*.msftconnecttest.com'
    - '*.msftncsi.com'
    - '+.xboxlive.com'
    - '*.*.stun.playstation.net'
    - xbox.*.*.microsoft.com
    - '*.ipv6.microsoft.com'
    - speedtest.cros.wr.pvp.net

proxies: # LEAVE THIS LINE!

proxy-groups:
  - name: '→ Remnawave'
    type: 'select'
    proxies: # LEAVE THIS LINE!

rules:
  - MATCH,→ Remnawave
`),
		name:         "Default",
		viewPosition: 2,
	},
	{
		templateType: "STASH",
		templateYAML: strPtr(`proxy-groups:
  - name: → Remnawave
    type: select
    proxies: # LEAVE THIS LINE!

proxies: # LEAVE THIS LINE!

rules:
  - SCRIPT,quic,REJECT
  - DOMAIN-SUFFIX,iphone-ld.apple.com,DIRECT
  - DOMAIN-SUFFIX,lcdn-locator.apple.com,DIRECT
  - DOMAIN-SUFFIX,lcdn-registration.apple.com,DIRECT
  - DOMAIN-SUFFIX,push.apple.com,DIRECT
  - PROCESS-NAME,v2ray,DIRECT
  - PROCESS-NAME,Surge,DIRECT
  - PROCESS-NAME,ss-local,DIRECT
  - PROCESS-NAME,privoxy,DIRECT
  - PROCESS-NAME,trojan,DIRECT
  - PROCESS-NAME,trojan-go,DIRECT
  - PROCESS-NAME,naive,DIRECT
  - PROCESS-NAME,CloudflareWARP,DIRECT
  - PROCESS-NAME,Cloudflare WARP,DIRECT
  - IP-CIDR,162.159.193.0/24,DIRECT,no-resolve
  - PROCESS-NAME,p4pclient,DIRECT
  - PROCESS-NAME,Thunder,DIRECT
  - PROCESS-NAME,DownloadService,DIRECT
  - PROCESS-NAME,qbittorrent,DIRECT
  - PROCESS-NAME,Transmission,DIRECT
  - PROCESS-NAME,fdm,DIRECT
  - PROCESS-NAME,aria2c,DIRECT
  - PROCESS-NAME,Folx,DIRECT
  - PROCESS-NAME,NetTransport,DIRECT
  - PROCESS-NAME,uTorrent,DIRECT
  - PROCESS-NAME,WebTorrent,DIRECT
  - GEOIP,LAN,DIRECT
  - MATCH,→ Remnawave
script:
  shortcuts:
    quic: network == 'udp' and dst_port == 443
dns:
  default-nameserver:
    - 1.1.1.1
    - 1.0.0.1
  nameserver:
    - 1.1.1.1
    - 1.0.0.1
log-level: warning
mode: rule
`),
		name:         "Default",
		viewPosition: 3,
	},
	{
		templateType: "CLASH",
		templateYAML: strPtr(`mixed-port: 7890
socks-port: 7891
redir-port: 7892
allow-lan: true
mode: global
log-level: info
external-controller: 127.0.0.1:9090
dns:
  enable: true
  use-hosts: true
  enhanced-mode: fake-ip
  fake-ip-range: 198.18.0.1/16
  default-nameserver:
    - 1.1.1.1
    - 8.8.8.8
  nameserver:
    - 1.1.1.1
    - 8.8.8.8
  fake-ip-filter:
    - '*.lan'
    - stun.*.*.*
    - stun.*.*
    - time.windows.com
    - time.nist.gov
    - time.apple.com
    - time.asia.apple.com
    - '*.openwrt.pool.ntp.org'
    - pool.ntp.org
    - ntp.ubuntu.com
    - time1.apple.com
    - time2.apple.com
    - time3.apple.com
    - time4.apple.com
    - time5.apple.com
    - time6.apple.com
    - time7.apple.com
    - time1.google.com
    - time2.google.com
    - time3.google.com
    - time4.google.com
    - api.joox.com
    - joox.com
    - '*.xiami.com'
    - '*.msftconnecttest.com'
    - '*.msftncsi.com'
    - '+.xboxlive.com'
    - '*.*.stun.playstation.net'
    - xbox.*.*.microsoft.com
    - '*.ipv6.microsoft.com'
    - speedtest.cros.wr.pvp.net

proxies: # LEAVE THIS LINE!

proxy-groups:
  - name: '→ Remnawave'
    type: 'select'
    proxies: # LEAVE THIS LINE!

rules:
  - MATCH,→ Remnawave`),
		name:         "Default",
		viewPosition: 4,
	},
	{
		templateType: "SINGBOX",
		templateJSON: strPtr(`{"dns":{"rules":[{"server":"remote","query_type":["A","AAAA"]},{"server":"local","outbound":"any"}],"fakeip":{"enabled":true,"inet4_range":"198.18.0.0/15","inet6_range":"fc00::/18"},"servers":[{"tag":"cf-dns","address":"tls://1.1.1.1"},{"tag":"local","detour":"direct","address":"tcp://1.1.1.1","strategy":"ipv4_only","address_strategy":"prefer_ipv4"},{"tag":"remote","address":"fakeip"}],"independent_cache":true},"log":{"level":"debug","disabled":true,"timestamp":true},"route":{"rules":[{"action":"sniff"},{"mode":"or","type":"logical","rules":[{"protocol":"dns"},{"port":53}],"action":"hijack-dns"},{"outbound":"direct","ip_is_private":true}],"override_android_vpn":true,"auto_detect_interface":true},"inbounds":[{"mtu":9000,"tag":"tun-in","type":"tun","sniff":true,"stack":"mixed","platform":{"http_proxy":{"server":"127.0.0.1","enabled":true,"server_port":2412}},"auto_route":true,"strict_route":true,"inet4_address":"172.19.0.1/30","inet6_address":"fdfe:dcba:9876::1/126","interface_name":"tun125","endpoint_independent_nat":true},{"tag":"mixed-in","type":"mixed","sniff":true,"users":[],"listen":"127.0.0.1","listen_port":2412,"set_system_proxy":false}],"outbounds":[{"tag":"→ Remnawave","type":"selector","outbounds":null,"interrupt_exist_connections":true},{"tag":"direct","type":"direct"}],"experimental":{"clash_api":{"external_ui":"yacd","default_mode":"rule","external_controller":"127.0.0.1:9090","external_ui_download_url":"https://github.com/MetaCubeX/Yacd-meta/archive/gh-pages.zip","external_ui_download_detour":"direct"},"cache_file":{"path":"remnawave.db","enabled":true,"cache_id":"remnawave","store_fakeip":true}}}`),
		name:         "Default",
		viewPosition: 5,
	},
}

const defaultResponseRules = `{"rules":[{"name":"Browser Subscription","enabled":true,"operator":"AND","conditions":[{"value":"text/html","operator":"CONTAINS","headerName":"accept","caseSensitive":true}],"description":"System critical: do not delete or disable this rule.","responseType":"BROWSER"},{"name":"Mihomo Clients","enabled":true,"operator":"AND","conditions":[{"value":"^(?:FlClash|FlClashX|Flowvy|[Cc]lash-[Vv]erge|[Kk]oala-[Cc]lash|[Cc]lash-?[Mm]eta|[Mm]urge|[Cc]lashX [Mm]eta|[Mm]ihomo|[Cc]lash-nyanpasu|clash.meta|prizrak-box)","operator":"REGEX","headerName":"user-agent","caseSensitive":false}],"description":"Response with generated YAML config (Mihomo Template)","responseType":"MIHOMO"},{"name":"Stash (iOS, macOS)","enabled":true,"operator":"AND","conditions":[{"value":"^stash","operator":"REGEX","headerName":"user-agent","caseSensitive":false}],"description":"Response with generated YAML config (Stash Template)","responseType":"STASH"},{"name":"Sing-box clients","enabled":true,"operator":"AND","conditions":[{"value":"^sfa|sfi|sfm|sft|karing|singbox|rabbithole","operator":"REGEX","headerName":"user-agent","caseSensitive":false}],"description":"Resonse with generated JSON config (Singbox template)","responseType":"SINGBOX"},{"name":"Clash Core Clients","enabled":true,"operator":"AND","conditions":[{"value":"^clash","operator":"REGEX","headerName":"user-agent","caseSensitive":false}],"description":"Response with generated YAML config (Clash Template)","responseType":"CLASH"},{"name":"Fallback Base64","enabled":true,"operator":"AND","conditions":[],"description":"System critical: do not delete or disable this rule.","responseType":"XRAY_BASE64"}],"version":"1"}`

const defaultHWIDSettings = `{"enabled":false,"maxDevicesAnnounce":null,"fallbackDeviceLimit":999}`

const defaultCustomRemarks = `{"emptyHosts":["→ v2raystat","→ No hosts found","→ Check Hosts tab","→ Check Internal Squads tab"],"expiredUsers":["⌛ Subscription expired","Contact support"],"limitedUsers":["🚧 Subscription limited","Contact support"],"disabledUsers":["🚫 Subscription disabled","Contact support"],"HWIDNotSupported":["App not supported"],"HWIDMaxDevicesExceeded":["Limit of devices reached"]}`

func ensureDefaultSubscriptionData(fileDB *sql.DB, cfg *config.BackendConfig) error {
	tx, err := fileDB.BeginTx(context.Background(), nil)
	if err != nil {
		return fmt.Errorf("begin defaults transaction: %w", err)
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback()
		}
	}()

	if err = ensureDefaultSubscriptionSettings(tx); err != nil {
		return err
	}

	if err = ensureDefaultTemplates(tx); err != nil {
		return err
	}

	if err = tx.Commit(); err != nil {
		return fmt.Errorf("commit defaults transaction: %w", err)
	}

	cfg.Logger.Info("Default subscription data ensured")
	return nil
}

func ensureDefaultSubscriptionSettings(tx *sql.Tx) error {
	var count int
	if err := tx.QueryRow(`SELECT COUNT(*) FROM subscription_settings`).Scan(&count); err != nil {
		return fmt.Errorf("count subscription_settings: %w", err)
	}
	if count > 0 {
		return nil
	}

	_, err := tx.Exec(`
		INSERT INTO subscription_settings (
			uuid, profile_title, support_link, profile_update_interval,
			address, port, api_schema, api_path,
			is_profile_webpage_url_enabled, serve_json_at_base_subscription,
			happ_announce, happ_routing, is_show_custom_remarks,
			custom_remarks, custom_response_headers, randomize_hosts,
			response_rules, hwid_settings
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		uuid.NewString(),
		"v2raystat",
		"https://github.com",
		12,
		"",
		9263,
		"grpc",
		"",
		true,
		false,
		"",
		"",
		true,
		defaultCustomRemarks,
		`{}`,
		false,
		defaultResponseRules,
		defaultHWIDSettings,
	)
	if err != nil {
		return fmt.Errorf("insert default subscription_settings: %w", err)
	}

	return nil
}

func ensureDefaultTemplates(tx *sql.Tx) error {
	for _, tmpl := range defaultSubscriptionTemplates {
		var count int
		err := tx.QueryRow(`SELECT COUNT(*) FROM subscription_templates WHERE template_type = ?`, tmpl.templateType).Scan(&count)
		if err != nil {
			return fmt.Errorf("count template %s: %w", tmpl.templateType, err)
		}
		if count > 0 {
			continue
		}

		_, err = tx.Exec(`
			INSERT INTO subscription_templates (
				uuid, view_position, name, template_type, template_yaml, template_json
			) VALUES (?, ?, ?, ?, ?, ?)`,
			uuid.NewString(),
			tmpl.viewPosition,
			tmpl.name,
			tmpl.templateType,
			tmpl.templateYAML,
			tmpl.templateJSON,
		)
		if err != nil {
			return fmt.Errorf("insert template %s: %w", tmpl.templateType, err)
		}
	}

	return nil
}

func strPtr(v string) *string {
	return &v
}
