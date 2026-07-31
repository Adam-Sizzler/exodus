package seed

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
)

// defaultResponseRules defines the initial Subscription Response Rules in pretty-formatted JSON structure.
const defaultResponseRules = `{
  "version": "1",
  "rules": [
    {
      "name": "Browser Subscription",
      "enabled": true,
      "operator": "AND",
      "conditions": [
        {
          "headerName": "accept",
          "operator": "CONTAINS",
          "value": "text/html",
          "caseSensitive": true
        }
      ],
      "description": "System critical: do not delete or disable this rule.",
      "responseType": "BROWSER"
    },
    {
      "name": "Mihomo Clients",
      "enabled": true,
      "operator": "AND",
      "conditions": [
        {
          "headerName": "user-agent",
          "operator": "REGEX",
          "value": "^(?:FlClash|FlClashX|Flowvy|[Cc]lash-[Vv]erge|[Kk]oala-[Cc]lash|[Cc]lash-?[Mm]eta|[Mm]urge|[Cc]lashX [Mm]eta|[Mm]ihomo|[Cc]lash-nyanpasu|clash.meta|prizrak-box)",
          "caseSensitive": false
        }
      ],
      "description": "Response with generated YAML config (Mihomo Template)",
      "responseType": "MIHOMO"
    },
    {
      "name": "Stash (iOS, macOS)",
      "enabled": true,
      "operator": "AND",
      "conditions": [
        {
          "headerName": "user-agent",
          "operator": "REGEX",
          "value": "^stash",
          "caseSensitive": false
        }
      ],
      "description": "Response with generated YAML config (Stash Template)",
      "responseType": "STASH"
    },
    {
      "name": "Sing-box clients",
      "enabled": true,
      "operator": "AND",
      "conditions": [
        {
          "headerName": "user-agent",
          "operator": "REGEX",
          "value": "^sfa|sfi|sfm|sft|karing|singbox|rabbithole",
          "caseSensitive": false
        }
      ],
      "description": "Resonse with generated JSON config (Singbox template)",
      "responseType": "SINGBOX"
    },
    {
      "name": "Clash Core Clients",
      "enabled": true,
      "operator": "AND",
      "conditions": [
        {
          "headerName": "user-agent",
          "operator": "REGEX",
          "value": "^clash",
          "caseSensitive": false
        }
      ],
      "description": "Response with generated YAML config (Clash Template)",
      "responseType": "CLASH"
    },
    {
      "name": "Fallback Base64",
      "enabled": true,
      "operator": "AND",
      "conditions": [],
      "description": "System critical: do not delete or disable this rule.",
      "responseType": "XRAY_BASE64"
    }
  ]
}`

// PrevResponseRulesHash stores the canonical SHA-256 hash of the PREVIOUS defaultResponseRules.
// IMPORTANT: When updating defaultResponseRules, calculate its current canonical hash and set it as PrevResponseRulesHash FIRST before updating defaultResponseRules.
const PrevResponseRulesHash = "4e61e34171f1ad37b23b2ef9dc22a441dca4beef75bc98fe2a54fe3685c52ad8"

const defaultHWIDSettings = `{
  "enabled": false,
  "maxDevicesAnnounce": null,
  "fallbackDeviceLimit": 999
}`

const defaultCustomRemarks = `{
  "emptyHosts": [
    "→ exodus",
    "→ No hosts found",
    "→ Check Hosts tab",
    "→ Check Internal Squads tab"
  ],
  "expiredUsers": [
    "⌛ Subscription expired",
    "Contact support"
  ],
  "limitedUsers": [
    "🚧 Subscription limited",
    "Contact support"
  ],
  "disabledUsers": [
    "🚫 Subscription disabled",
    "Contact support"
  ],
  "HWIDNotSupported": [
    "App not supported"
  ],
  "HWIDMaxDevicesExceeded": [
    "Limit of devices reached"
  ]
}`

const defaultSubpageConfigUUID = "00000000-0000-0000-0000-000000000000"

const defaultSingboxConfig = `{
  "log": {
    "level": "info"
  },
  "dns": {
    "servers": [
      {
        "tag": "dns-remote",
        "type": "udp",
        "server": "1.1.1.1",
        "detour": "direct"
      }
    ]
  },
  "inbounds": [
    {
      "type": "shadowsocks",
      "tag": "ss-in",
      "listen": "127.0.0.1",
      "listen_port": 2080,
      "method": "chacha20-ietf-poly1305"
    },
    {
      "type": "trojan",
      "tag": "trojan-in",
      "listen": "127.0.0.1",
      "listen_port": 2443,
      "users": []
    }
  ],
  "outbounds": [
    {
      "type": "direct",
      "tag": "direct"
    },
    {
      "type": "block",
      "tag": "block"
    }
  ],
  "route": {
    "final": "direct"
  }
}`

const defaultPasskeySettings = `{
  "enabled": false,
  "origin": null,
  "rpId": null
}`

const defaultOAuth2Settings = `{
  "generic": {
    "allowedEmails": [],
    "authorizationUrl": null,
    "clientId": null,
    "clientSecret": null,
    "enabled": false,
    "frontendDomain": null,
    "tokenUrl": null,
    "withPkce": false
  },
  "github": {
    "allowedEmails": [],
    "clientId": null,
    "clientSecret": null,
    "enabled": false
  },
  "keycloak": {
    "allowedEmails": [],
    "clientId": null,
    "clientSecret": null,
    "enabled": false,
    "frontendDomain": null,
    "keycloakDomain": null,
    "realm": null
  },
  "pocketid": {
    "allowedEmails": [],
    "clientId": null,
    "clientSecret": null,
    "enabled": false,
    "plainDomain": null
  },
  "telegram": {
    "allowedIds": [],
    "clientId": null,
    "clientSecret": null,
    "enabled": false,
    "frontendDomain": null
  },
  "yandex": {
    "allowedEmails": [],
    "clientId": null,
    "clientSecret": null,
    "enabled": false
  }
}`

const defaultPasswordSettings = `{
  "enabled": true
}`

const defaultBrandingSettings = `{
  "logoUrl": null,
  "title": "EXODUS"
}`

func canonicalHash(rawJSON string) (string, error) {
	var v any
	if err := json.Unmarshal([]byte(rawJSON), &v); err != nil {
		return "", err
	}
	canonical, err := json.Marshal(v)
	if err != nil {
		return "", err
	}
	sum := sha256.Sum256(canonical)
	return hex.EncodeToString(sum[:]), nil
}
