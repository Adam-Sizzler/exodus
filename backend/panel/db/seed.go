package db

import (
	"context"
	"database/sql"
	"fmt"

	"v2ray-stat/backend/panel/config"
	"v2ray-stat/backend/panel/dbutil"

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
  - name: '→ V2RS'
    type: 'select'
    proxies: # LEAVE THIS LINE!

rules:
  - MATCH,→ V2RS
`),
		name:         "Default",
		viewPosition: 2,
	},
	{
		templateType: "STASH",
		templateYAML: strPtr(`proxy-groups:
  - name: → V2RS
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
  - MATCH,→ V2RS
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
  - name: '→ V2RS'
    type: 'select'
    proxies: # LEAVE THIS LINE!

rules:
  - MATCH,→ V2RS
`),
		name:         "Default",
		viewPosition: 4,
	},
}

const defaultResponseRules = `{"rules":[{"name":"Browser Subscription","enabled":true,"operator":"AND","conditions":[{"value":"text/html","operator":"CONTAINS","headerName":"accept","caseSensitive":true}],"description":"System critical: do not delete or disable this rule.","responseType":"BROWSER"},{"name":"Mihomo Clients","enabled":true,"operator":"AND","conditions":[{"value":"^(?:FlClash|FlClashX|Flowvy|[Cc]lash-[Vv]erge|[Kk]oala-[Cc]lash|[Cc]lash-?[Mm]eta|[Mm]urge|[Cc]lashX [Mm]eta|[Mm]ihomo|[Cc]lash-nyanpasu|clash.meta|prizrak-box)","operator":"REGEX","headerName":"user-agent","caseSensitive":false}],"description":"Response with generated YAML config (Mihomo Template)","responseType":"MIHOMO"},{"name":"Stash (iOS, macOS)","enabled":true,"operator":"AND","conditions":[{"value":"^stash","operator":"REGEX","headerName":"user-agent","caseSensitive":false}],"description":"Response with generated YAML config (Stash Template)","responseType":"STASH"},{"name":"Sing-box clients","enabled":true,"operator":"AND","conditions":[{"value":"^sfa|sfi|sfm|sft|karing|singbox|rabbithole","operator":"REGEX","headerName":"user-agent","caseSensitive":false}],"description":"Resonse with generated JSON config (Singbox template)","responseType":"SINGBOX"},{"name":"Clash Core Clients","enabled":true,"operator":"AND","conditions":[{"value":"^clash","operator":"REGEX","headerName":"user-agent","caseSensitive":false}],"description":"Response with generated YAML config (Clash Template)","responseType":"CLASH"},{"name":"Fallback Base64","enabled":true,"operator":"AND","conditions":[],"description":"System critical: do not delete or disable this rule.","responseType":"XRAY_BASE64"}],"version":"1"}`
const defaultHWIDSettings = `{"enabled":false,"maxDevicesAnnounce":null,"fallbackDeviceLimit":999}`
const defaultCustomRemarks = `{"emptyHosts":["→ v2rs","→ No hosts found","→ Check Hosts tab","→ Check Internal Squads tab"],"expiredUsers":["⌛ Subscription expired","Contact support"],"limitedUsers":["🚧 Subscription limited","Contact support"],"disabledUsers":["🚫 Subscription disabled","Contact support"],"HWIDNotSupported":["App not supported"],"HWIDMaxDevicesExceeded":["Limit of devices reached"]}`
const defaultSubpageConfigUUID = "00000000-0000-0000-0000-000000000000"
const defaultSubscriptionPageConfig = `{
  "version": "1",
  "locales": ["en", "ru"],
  "brandingSettings": {
    "title": "Subscription",
    "logoUrl": "",
    "supportUrl": "https://github.com"
  },
  "uiConfig": {
    "subscriptionInfoBlockType": "expanded",
    "installationGuidesBlockType": "cards"
  },
  "baseSettings": {
    "metaTitle": "Subscription",
    "metaDescription": "Subscription",
    "showConnectionKeys": false,
    "hideGetLinkButton": false
  },
  "baseTranslations": {
    "installationGuideHeader": { "en": "Installation Guide", "ru": "Инструкция по установке" },
    "connectionKeysHeader": { "en": "Connection Keys", "ru": "Ключи подключения" },
    "linkCopied": { "en": "Link Copied", "ru": "Ссылка скопирована" },
    "linkCopiedToClipboard": { "en": "Link copied to clipboard", "ru": "Ссылка скопирована в буфер" },
    "getLink": { "en": "Get Link", "ru": "Получить ссылку" },
    "scanQrCode": { "en": "Scan QR Code", "ru": "Сканировать QR-код" },
    "scanQrCodeDescription": { "en": "Scan QR code with your client", "ru": "Отсканируйте QR-код в приложении" },
    "copyLink": { "en": "Copy Link", "ru": "Скопировать ссылку" },
    "name": { "en": "Name", "ru": "Имя" },
    "status": { "en": "Status", "ru": "Статус" },
    "active": { "en": "Active", "ru": "Активен" },
    "inactive": { "en": "Inactive", "ru": "Неактивен" },
    "expires": { "en": "Expires", "ru": "Истекает" },
    "bandwidth": { "en": "Bandwidth", "ru": "Трафик" },
    "scanToImport": { "en": "Scan to Import", "ru": "Сканировать для импорта" },
    "expiresIn": { "en": "Expires In", "ru": "Истекает через" },
    "expired": { "en": "Expired", "ru": "Истек" },
    "unknown": { "en": "Unknown", "ru": "Неизвестно" },
    "indefinitely": { "en": "Indefinitely", "ru": "Бессрочно" }
  },
  "svgLibrary": {
    "default": "<svg xmlns=\\"http://www.w3.org/2000/svg\\" viewBox=\\"0 0 24 24\\"><circle cx=\\"12\\" cy=\\"12\\" r=\\"10\\" fill=\\"currentColor\\"/></svg>"
  },
  "platforms": {
    "android": {
      "displayName": { "en": "Android", "ru": "Android" },
      "svgIconKey": "default",
      "apps": [
        {
          "name": "V2RayNG",
          "svgIconKey": "default",
          "featured": true,
          "blocks": [
            {
              "svgIconKey": "default",
              "svgIconColor": "blue",
              "title": { "en": "Install the app", "ru": "Установите приложение" },
              "description": { "en": "Download and install the client.", "ru": "Скачайте и установите клиент." },
              "buttons": [
                {
                  "link": "https://play.google.com/store/apps/details?id=com.v2ray.ang",
                  "type": "external",
                  "text": { "en": "Download", "ru": "Скачать" },
                  "svgIconKey": "default"
                },
                {
                  "link": "{{SUBSCRIPTION_LINK}}",
                  "type": "subscriptionLink",
                  "text": { "en": "Get Link", "ru": "Получить ссылку" },
                  "svgIconKey": "default"
                }
              ]
            }
          ]
        }
      ]
    },
    "ios": {
      "displayName": { "en": "iOS", "ru": "iOS" },
      "svgIconKey": "default",
      "apps": [
        {
          "name": "Stash",
          "svgIconKey": "default",
          "featured": true,
          "blocks": [
            {
              "svgIconKey": "default",
              "svgIconColor": "green",
              "title": { "en": "Install the app", "ru": "Установите приложение" },
              "description": { "en": "Download and install the client.", "ru": "Скачайте и установите клиент." },
              "buttons": [
                {
                  "link": "https://apps.apple.com/app/stash/id1596063349",
                  "type": "external",
                  "text": { "en": "Download", "ru": "Скачать" },
                  "svgIconKey": "default"
                },
                {
                  "link": "{{SUBSCRIPTION_LINK}}",
                  "type": "subscriptionLink",
                  "text": { "en": "Get Link", "ru": "Получить ссылку" },
                  "svgIconKey": "default"
                }
              ]
            }
          ]
        }
      ]
    }
  }
}`

// SeedDefaults inserts base settings and templates if they do not exist.
func SeedDefaults(ctx context.Context, dbConn *sql.DB, cfg *config.BackendConfig) error {
	if dbConn == nil {
		return fmt.Errorf("database connection is nil")
	}

	tx, err := dbConn.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("begin defaults transaction: %w", err)
	}

	if err := ensureV2rsSettings(tx); err != nil {
		_ = tx.Rollback()
		return err
	}
	if err := ensureDefaultSubscriptionSettings(tx); err != nil {
		_ = tx.Rollback()
		return err
	}
	if err := ensureDefaultTemplates(tx); err != nil {
		_ = tx.Rollback()
		return err
	}
	if err := ensureDefaultSubscriptionPageConfig(tx); err != nil {
		_ = tx.Rollback()
		return err
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit defaults transaction: %w", err)
	}

	cfg.Logger.Info("Default configuration data ensured")
	return nil
}

func ensureV2rsSettings(tx *sql.Tx) error {
	var exists int
	if err := tx.QueryRowContext(context.Background(), `SELECT COUNT(*) FROM v2rs_settings WHERE id = 1`).Scan(&exists); err != nil {
		return fmt.Errorf("check v2rs_settings row: %w", err)
	}
	if exists > 0 {
		return nil
	}

	passkeySettings := `{"rpId":null,"origin":null,"enabled":false}`
	oauth2Settings := `{"github":{"enabled":false,"clientId":null,"clientSecret":null,"allowedEmails":[]},"yandex":{"enabled":false,"clientId":null,"clientSecret":null,"allowedEmails":[]},"generic":{"enabled":false,"clientId":null,"tokenUrl":null,"withPkce":false,"clientSecret":null,"allowedEmails":[],"frontendDomain":null,"authorizationUrl":null},"keycloak":{"realm":null,"enabled":false,"clientId":null,"clientSecret":null,"allowedEmails":[],"frontendDomain":null,"keycloakDomain":null},"pocketid":{"enabled":false,"clientId":null,"plainDomain":null,"clientSecret":null,"allowedEmails":[]}}`
	tgAuthSettings := `{"enabled":false,"adminIds":[],"botToken":null}`
	passwordSettings := `{"enabled":true}`
	brandingSettings := `{"title":"V2RS","logoUrl":null}`

	query := dbutil.Rebind(`
		INSERT INTO v2rs_settings (
			id, passkey_settings, oauth2_settings, tg_auth_settings, password_settings, branding_settings
		) VALUES (1, ?, ?, ?, ?, ?)
	`)
	if _, err := tx.ExecContext(context.Background(), query, passkeySettings, oauth2Settings, tgAuthSettings, passwordSettings, brandingSettings); err != nil {
		return fmt.Errorf("insert default v2rs_settings row: %w", err)
	}

	return nil
}

func ensureDefaultSubscriptionSettings(tx *sql.Tx) error {
	var count int
	if err := tx.QueryRowContext(context.Background(), `SELECT COUNT(*) FROM subscription_settings`).Scan(&count); err != nil {
		return fmt.Errorf("count subscription_settings: %w", err)
	}
	if count > 0 {
		return nil
	}

	query := dbutil.Rebind(`
		INSERT INTO subscription_settings (
			uuid, profile_title, support_link, profile_update_interval,
			address, port, api_schema, api_path,
			is_profile_webpage_url_enabled, serve_json_at_base_subscription,
			happ_announce, happ_routing, is_show_custom_remarks,
			custom_remarks, custom_response_headers, randomize_hosts,
			response_rules, hwid_settings
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`)
	_, err := tx.ExecContext(context.Background(), query,
		uuid.NewString(),
		"v2rs",
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
		if err := tx.QueryRowContext(context.Background(), dbutil.Rebind(`SELECT COUNT(*) FROM subscription_templates WHERE template_type = ?`), tmpl.templateType).Scan(&count); err != nil {
			return fmt.Errorf("count template %s: %w", tmpl.templateType, err)
		}
		if count > 0 {
			continue
		}

		query := dbutil.Rebind(`
			INSERT INTO subscription_templates (
				uuid, view_position, name, template_type, template_yaml, template_json
			) VALUES (?, ?, ?, ?, ?, ?)
		`)
		if _, err := tx.ExecContext(context.Background(), query,
			uuid.NewString(),
			tmpl.viewPosition,
			tmpl.name,
			tmpl.templateType,
			tmpl.templateYAML,
			tmpl.templateJSON,
		); err != nil {
			return fmt.Errorf("insert template %s: %w", tmpl.templateType, err)
		}
	}

	return nil
}

func ensureDefaultSubscriptionPageConfig(tx *sql.Tx) error {
	var count int
	if err := tx.QueryRowContext(context.Background(), `SELECT COUNT(*) FROM subscription_page_config`).Scan(&count); err != nil {
		return fmt.Errorf("count subscription_page_config: %w", err)
	}
	if count > 0 {
		return nil
	}

	query := dbutil.Rebind(`
		INSERT INTO subscription_page_config (
			uuid, view_position, name, config
		) VALUES (?, ?, ?, ?)
	`)
	if _, err := tx.ExecContext(context.Background(), query, defaultSubpageConfigUUID, 1, "Default", defaultSubscriptionPageConfig); err != nil {
		return fmt.Errorf("insert default subscription_page_config: %w", err)
	}

	return nil
}

func strPtr(v string) *string {
	if v == "" {
		return nil
	}
	return &v
}
