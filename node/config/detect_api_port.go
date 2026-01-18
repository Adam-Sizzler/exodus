package config

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

func detectApiPort(configPath string) (string, string) {
	data, err := os.ReadFile(configPath)
	if err != nil {
		return "", ""
	}

	var raw map[string]any
	if err := json.Unmarshal(data, &raw); err != nil {
		return "", ""
	}

	// Поиск в "api": { "listen": "127.0.0.1:10812" } (Стандарт Xray)
	if api, ok := raw["api"].(map[string]any); ok {
		if listen, ok := api["listen"].(string); ok {
			return splitAddr(listen)
		}
	}

	// Поиск в "experimental": { "v2ray_api": { "listen": "..." } } (Sing-box)
	if exp, ok := raw["experimental"].(map[string]any); ok {
		if v2api, ok := exp["v2ray-api"].(map[string]any); ok {
			if listen, ok := v2api["listen"].(string); ok {
				return splitAddr(listen)
			}
		}
	}

	// Поиск в "inbounds": [ { "tag": "api", "port": 9953 } ]
	if inbounds, ok := raw["inbounds"].([]any); ok {
		for _, ib := range inbounds {
			if m, ok := ib.(map[string]any); ok {
				if m["tag"] == "api" {
					if port, ok := m["port"]; ok {
						portStr := fmt.Sprintf("%v", port)
						addr := "127.0.0.1"
						if l, ok := m["listen"].(string); ok && l != "0.0.0.0" {
							addr = l
						}
						return addr, portStr
					}
				}
			}
		}
	}

	return "", ""
}

func splitAddr(addr string) (string, string) {
	if !strings.Contains(addr, ":") {
		return addr, ""
	}
	parts := strings.Split(addr, ":")
	return parts[0], parts[1]
}
