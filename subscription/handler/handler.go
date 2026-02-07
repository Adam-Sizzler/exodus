package handler

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"maps"
	"net/http"
	"os"
	"path/filepath"
	"slices"
	"strings"

	"v2ray-stat/subscription/api"
	"v2ray-stat/subscription/config"
	"v2ray-stat/subscription/templates"
)

var clientFormats = map[string]struct {
	Format string
	Header map[string]string
}{
	"xray":    {"json", map[string]string{"Content-Type": "application/json"}},
	"singbox": {"json", map[string]string{"Content-Type": "application/json"}},
	"mihomo":  {"yaml", map[string]string{"Content-Type": "text/yaml"}},
	"edg":     {"json", map[string]string{"Content-Type": "application/json"}},
	"unknown": {"json", map[string]string{"Content-Type": "application/json"}},
}

func SubscriptionHandler(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	mode, client, user := q.Get("mode"), q.Get("client"), q.Get("user")

	userAgent := strings.ToLower(r.Header.Get("User-Agent"))
	accept := r.Header.Get("Accept")

	isBrowser := strings.Contains(userAgent, "mozilla/5.0") || strings.Contains(userAgent, "edg")
	isHtmlRequest := strings.Contains(accept, "text/html")

	if client == "" {
		switch {
		case strings.Contains(userAgent, "mihomo"),
			strings.Contains(userAgent, "clash"),
			strings.Contains(userAgent, "flclash"),
			strings.Contains(userAgent, "flclashx"),
			strings.Contains(userAgent, "verge"),
			strings.Contains(userAgent, "clash-verge"),
			strings.Contains(userAgent, "clash verge"),
			strings.Contains(userAgent, "koala-clash"),
			strings.Contains(userAgent, "clash-meta"),
			strings.Contains(userAgent, "clashmeta"),
			strings.Contains(userAgent, "merge"),
			strings.Contains(userAgent, "clashx meta"),
			strings.Contains(userAgent, "clash-nyanpasu"),
			strings.Contains(userAgent, "clash.meta"),
			strings.Contains(userAgent, "prizrak-box"),
			strings.Contains(userAgent, "rabbithole"),
			strings.Contains(userAgent, "flowvy"):
			client = "mihomo"

		case strings.Contains(userAgent, "sfa"),
			strings.Contains(userAgent, "sfi"),
			strings.Contains(userAgent, "sfm"),
			strings.Contains(userAgent, "sft"),
			strings.Contains(userAgent, "karing"),
			strings.Contains(userAgent, "singbox"):
			client = "singbox"

		case strings.Contains(userAgent, "happ"),
			strings.Contains(userAgent, "v2rayng"),
			strings.Contains(userAgent, "io.github.saeeddev94.xray"),
			strings.Contains(userAgent, "v2rayn"):
			client = "xray"

		case strings.Contains(userAgent, "edg"),
			strings.Contains(userAgent, "telegrambot"):
			client = "unknown"

		default:
			client = "xray"
		}
	}

	if mode == "" {
		mode = "advanced"
	}

	cfg := config.GetConfig()
	cfg.Logger.Debug("Received subscription request", "client", client, "user", user, "mode", mode, "user_agent", r.Header.Get("User-Agent"))

	if user == "" {
		cfg.Logger.Warn("Missing user parameter in request")
		http.Error(w, "missing user parameter", http.StatusBadRequest)
		return
	}

	if _, ok := clientFormats[client]; !ok {
		cfg.Logger.Warn("Invalid client requested", "client", client)
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	if mode != "base" && mode != "advanced" {
		cfg.Logger.Warn("Invalid mode requested", "mode", mode)
		http.Error(w, "invalid mode, must be 'base' or 'advanced'", http.StatusBadRequest)
		return
	}

	userConfig, ok := cfg.Subscription.UserMap[user]
	if !ok {
		cfg.Logger.Debug("User not found in UserMap, applying defaults", "user", user)
		userConfig = config.UserConfig{
			Clients:       slices.Clone(cfg.Subscription.Defaults.Clients),
			IncludeNodes:  slices.Clone(cfg.Subscription.Defaults.IncludeNodes),
			NodeTemplates: make(map[string]map[string]string),
			Headers:       make(map[string]string),
		}
		for mode, templates := range cfg.Subscription.Defaults.NodeTemplates {
			userConfig.NodeTemplates[mode] = make(map[string]string)
			maps.Copy(userConfig.NodeTemplates[mode], templates)
		}
		maps.Copy(userConfig.Headers, cfg.Subscription.Defaults.Headers)
	}

	cfg.Logger.Trace("User config before merging", "user", user, "config", fmt.Sprintf("%+v", userConfig))

	if userConfig.Group != "" {
		if groupConfig, ok := cfg.Subscription.Groups[userConfig.Group]; ok {
			cfg.Logger.Debug("Merging group config", "group", userConfig.Group)
			if groupConfig.Clients != nil {
				userConfig.Clients = slices.Clone(groupConfig.Clients)
			}
			if groupConfig.IncludeNodes != nil {
				userConfig.IncludeNodes = slices.Clone(groupConfig.IncludeNodes)
			}
			if groupConfig.NodeTemplates != nil {
				userConfig.NodeTemplates = make(map[string]map[string]string)
				for mode, templates := range groupConfig.NodeTemplates {
					userConfig.NodeTemplates[mode] = make(map[string]string)
					if templates != nil {
						maps.Copy(userConfig.NodeTemplates[mode], templates)
					}
				}
			}
			if groupConfig.Headers != nil {
				userConfig.Headers = make(map[string]string)
				maps.Copy(userConfig.Headers, groupConfig.Headers)
			}
		} else {
			cfg.Logger.Warn("Group not found for user", "group", userConfig.Group, "user", user)
			http.Error(w, fmt.Sprintf("group %s not found for user %s", userConfig.Group, user), http.StatusBadRequest)
			return
		}
	}

	if len(userConfig.Clients) == 0 {
		userConfig.Clients = slices.Clone(cfg.Subscription.Defaults.Clients)
		cfg.Logger.Debug("Applied default clients for user", "user", user)
	}

	if len(userConfig.IncludeNodes) == 0 {
		userConfig.IncludeNodes = slices.Clone(cfg.Subscription.Defaults.IncludeNodes)
		cfg.Logger.Debug("Applied default nodes for user", "user", user)
	}

	if len(userConfig.NodeTemplates) == 0 {
		userConfig.NodeTemplates = make(map[string]map[string]string)
		for mode, templates := range cfg.Subscription.Defaults.NodeTemplates {
			userConfig.NodeTemplates[mode] = make(map[string]string)
			maps.Copy(userConfig.NodeTemplates[mode], templates)
		}
		cfg.Logger.Debug("Applied default templates for user", "user", user)
	}

	if len(userConfig.Headers) == 0 {
		userConfig.Headers = make(map[string]string)
		maps.Copy(userConfig.Headers, cfg.Subscription.Defaults.Headers)
		cfg.Logger.Debug("Applied default headers for user", "user", user)
	}

	if !slices.Contains(userConfig.Clients, client) {
		cfg.Logger.Warn("App not supported")
		http.Error(w, fmt.Sprintf("client %s not supported for user %s", client, user), http.StatusBadRequest)
		return
	}

	if !isBrowser && !slices.Contains(userConfig.Clients, client) {
		cfg.Logger.Warn("App not supported")
		http.Error(w, fmt.Sprintf("client %s not supported for user %s", client, user), http.StatusBadRequest)
		return
	}

	if isBrowser && !slices.Contains(userConfig.Clients, client) {
		if len(userConfig.Clients) > 0 {
			client = userConfig.Clients[0]
		}
	}

	userIDs, err := api.GetUserIDs(&cfg, user)
	if err != nil {
		cfg.Logger.Error("Failed to fetch user IDs", "user", user, "error", err)
		http.Error(w, fmt.Sprintf("failed to fetch user IDs: %v", err), http.StatusInternalServerError)
		return
	}

	cfg.Logger.Debug("Fetched user IDs", "user", user, "count", len(userIDs))

	var configs []string
	var totalUplink, totalDownlink, maxSubEnd, maxTrafficCap int64
	tmpls := templates.GetTemplates()

	for _, node := range userConfig.IncludeNodes {
		cfg.Logger.Trace("Processing node for user", "node", node, "user", user)

		modeTemplates, modeOk := userConfig.NodeTemplates[mode]
		if !modeOk || modeTemplates == nil {
			cfg.Logger.Debug("No templates for mode", "mode", mode, "user", user)
			continue
		}

		templateName, ok := modeTemplates[node]
		if !ok || templateName == "" {
			cfg.Logger.Warn("No template specified for node in mode", "node", node, "mode", mode, "user", user)
			continue
		}

		cfg.Logger.Trace("Found template for node", "node", node, "templateName", templateName)

		templatePath := filepath.Join(mode, templateName)
		template, ok := tmpls[client][templatePath]
		if !ok {
			cfg.Logger.Warn("Template not found for client for user", "template", templatePath, "client", client, "user", user)
			continue
		}

		var userID string
		var nodeTraffic api.UserID
		for _, uid := range userIDs {
			if uid.NodeName == node && uid.User == user {
				userID = uid.ID
				nodeTraffic = uid
				break
			}
		}

		if userID == "" {
			cfg.Logger.Warn("No user ID found for user on node", "user", user, "node", node)
			continue
		}

		meta, ok := cfg.NodeMetadata[node]
		if !ok {
			cfg.Logger.Warn("No metadata specified for node", "node", node, "user", user)
			continue
		}

		totalUplink += nodeTraffic.Uplink
		totalDownlink += nodeTraffic.Downlink
		if nodeTraffic.SubEnd > maxSubEnd {
			maxSubEnd = nodeTraffic.SubEnd
		}
		if nodeTraffic.TrafficCap > maxTrafficCap {
			maxTrafficCap = nodeTraffic.TrafficCap
		}

		configStr := strings.ReplaceAll(template, "{user_id}", userID)
		if meta.DomainPlaceholder != "" {
			configStr = strings.ReplaceAll(configStr, "{domain}", meta.DomainPlaceholder)
		}
		if meta.IPPlaceholder != "" {
			configStr = strings.ReplaceAll(configStr, "{ip}", meta.IPPlaceholder)
		}
		if meta.PortPlaceholder != "" {
			configStr = strings.ReplaceAll(configStr, "{port}", meta.PortPlaceholder)
		}
		if meta.RemarkPlaceholder != "" {
			configStr = strings.ReplaceAll(configStr, "{remark}", meta.RemarkPlaceholder)
		}

		configs = append(configs, configStr)
		cfg.Logger.Trace("Generated config for node", "node", node, "user", user, "domain", meta.DomainPlaceholder)
	}

	hasBase := false
	hasAdvanced := false

	baseConfigs := generateConfigsForMode(cfg, userConfig, client, user, "base", userIDs, tmpls)
	if len(baseConfigs) > 0 {
		hasBase = true
	}

	advConfigs := generateConfigsForMode(cfg, userConfig, client, user, "advanced", userIDs, tmpls)
	if len(advConfigs) > 0 {
		hasAdvanced = true
	}

	if isBrowser && isHtmlRequest {
		serveHtmlPage(w, totalUplink, totalDownlink, maxTrafficCap, maxSubEnd, userConfig, hasBase, hasAdvanced)
		return
	}

	if len(configs) == 0 {
		cfg.Logger.Warn("No configurations generated for user", "user", user)
		http.Error(w, "no configurations generated", http.StatusNotFound)
		return
	}

	cfg.Logger.Debug("Generated configurations", "count", len(configs), "user", user)

	var output any
	if mode == "base" {
		combinedConfig := strings.Join(configs, "\n")
		encodedConfig := base64.StdEncoding.EncodeToString([]byte(combinedConfig))
		output = encodedConfig
		w.Header().Set("Content-Type", "text/plain")
		cfg.Logger.Trace("Prepared base64-encoded config for base mode", "user", user, "client", client)
	} else if client == "mihomo" {
		baseTemplatePath := filepath.Join(mode, "base")
		baseTemplate, ok := tmpls[client][baseTemplatePath]
		if !ok {
			cfg.Logger.Warn("Base template not found for client for user", "template", baseTemplatePath, "client", client, "user", user)
			http.Error(w, "base template not found", http.StatusInternalServerError)
			return
		}

		var proxyStrings []string
		for _, config := range configs {
			lines := strings.Split(strings.TrimSpace(config), "\n")
			if len(lines) == 0 {
				continue
			}
			indented := []string{"  - " + strings.TrimSpace(lines[0])}
			for i := 1; i < len(lines); i++ {
				trimmed := strings.TrimSpace(lines[i])
				if trimmed != "" {
					indented = append(indented, "    "+trimmed)
				}
			}
			proxyStrings = append(proxyStrings, strings.Join(indented, "\n"))
		}

		combinedProxies := strings.Join(proxyStrings, "\n")
		configStr := strings.Replace(baseTemplate, "proxies: []", fmt.Sprintf("proxies:\n%s", combinedProxies), 1)
		output = []string{configStr}
		cfg.Logger.Trace("Combined proxies for mihomo", "user", user)
	} else {
		var parsedConfigs []any
		for _, configStr := range configs {
			var config any
			if clientFormats[client].Format == "json" {
				if err := json.Unmarshal([]byte(configStr), &config); err != nil {
					cfg.Logger.Error("Failed to parse template for client for user", "client", client, "user", user, "error", err)
					continue
				}
			} else {
				config = configStr
			}
			parsedConfigs = append(parsedConfigs, config)
		}

		if len(parsedConfigs) == 0 {
			cfg.Logger.Warn("No valid configurations generated for user", "user", user)
			http.Error(w, "no valid configurations generated", http.StatusNotFound)
			return
		}

		if client == "singbox" {
			if len(parsedConfigs) > 1 {
				cfg.Logger.Warn("Multiple configs for singbox - user first one", "user", user, "count", len(parsedConfigs))
			}
			output = parsedConfigs[0]
		} else {
			output = parsedConfigs
		}
		cfg.Logger.Debug("Parsed configurations for non-mihomo client", "client", client, "count", len(parsedConfigs))
	}

	for k, v := range clientFormats[client].Header {
		if mode == "base" && k == "Content-Type" {
			continue
		}
		w.Header().Set(k, v)
	}

	userInfo := fmt.Sprintf("upload=%d; download=%d; total=%d; expire=%d", totalUplink, totalDownlink, maxTrafficCap, maxSubEnd)
	w.Header().Set("Subscription-Userinfo", userInfo)

	for k, v := range userConfig.Headers {
		w.Header().Set(k, v)
	}

	cfg.Logger.Trace("Set response headers", "user", user)

	switch {
	case mode == "base":
		if _, err := w.Write([]byte(output.(string))); err != nil {
			cfg.Logger.Error("Failed to write base64-encoded response", "user", user, "error", err)
			http.Error(w, fmt.Sprintf("failed to write response: %v", err), http.StatusInternalServerError)
			return
		}
		cfg.Logger.Debug("Sent base64-encoded response", "user", user)
	case clientFormats[client].Format == "json":
		formattedJSON, err := json.MarshalIndent(output, "", "  ")
		if err != nil {
			cfg.Logger.Error("Failed to marshal JSON for client for user", "client", client, "user", user, "error", err)
			http.Error(w, fmt.Sprintf("failed to encode response: %v", err), http.StatusInternalServerError)
			return
		}
		if _, err := w.Write(formattedJSON); err != nil {
			cfg.Logger.Error("Failed to write JSON response", "user", user, "error", err)
			http.Error(w, fmt.Sprintf("failed to write response: %v", err), http.StatusInternalServerError)
			return
		}
		cfg.Logger.Debug("Sent JSON response", "user", user)
	case clientFormats[client].Format == "yaml":
		if _, err := w.Write([]byte(output.([]string)[0])); err != nil {
			cfg.Logger.Error("Failed to write YAML response", "user", user, "error", err)
			http.Error(w, fmt.Sprintf("failed to write response: %v", err), http.StatusInternalServerError)
			return
		}
		cfg.Logger.Debug("Sent YAML response", "user", user)
	}
}

func generateConfigsForMode(cfg config.Config, userConfig config.UserConfig, client string, user string, mode string, userIDs []api.UserID, tmpls map[string]map[string]string) []string {
	var configs []string

	for _, node := range userConfig.IncludeNodes {
		modeTemplates, modeOk := userConfig.NodeTemplates[mode]
		if !modeOk || modeTemplates == nil {
			continue
		}

		templateName, ok := modeTemplates[node]
		if !ok || templateName == "" {
			continue
		}

		templatePath := filepath.Join(mode, templateName)
		template, ok := tmpls[client][templatePath]
		if !ok {
			continue
		}

		var userID string
		for _, uid := range userIDs {
			if uid.NodeName == node && uid.User == user {
				userID = uid.ID
				break
			}
		}

		if userID == "" {
			continue
		}

		meta, ok := cfg.NodeMetadata[node]
		if !ok {
			continue
		}

		configStr := strings.ReplaceAll(template, "{user_id}", userID)
		if meta.DomainPlaceholder != "" {
			configStr = strings.ReplaceAll(configStr, "{domain}", meta.DomainPlaceholder)
		}
		if meta.IPPlaceholder != "" {
			configStr = strings.ReplaceAll(configStr, "{ip}", meta.IPPlaceholder)
		}
		if meta.PortPlaceholder != "" {
			configStr = strings.ReplaceAll(configStr, "{port}", meta.PortPlaceholder)
		}
		if meta.RemarkPlaceholder != "" {
			configStr = strings.ReplaceAll(configStr, "{remark}", meta.RemarkPlaceholder)
		}

		configs = append(configs, configStr)
	}

	return configs
}

func serveHtmlPage(w http.ResponseWriter, uplink int64, downlink int64, trafficCap int64, subEnd int64, userConfig config.UserConfig, hasBase bool, hasAdvanced bool) {
	cfg := config.GetConfig()

	// 1. Читаем HTML файл
	htmlPath := "/app/assets/index.html"
	htmlBytes, err := os.ReadFile(htmlPath)
	if err != nil {
		cfg.Logger.Error("Failed to read HTML file", "error", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	html := string(htmlBytes)

	// 2. Подготовка данных для JS
	userData := map[string]any{
		"upload":           uplink,
		"download":         downlink,
		"total":            trafficCap,
		"expire":           subEnd,
		"profileTitle":     safeBase64Decode(userConfig.Headers["Profile-Title"]),
		"announce":         safeBase64Decode(userConfig.Headers["Announce"]),
		"supportUrl":       userConfig.Headers["Support-Url"],
		"availableClients": userConfig.Clients,
	}

	userDataJSON, err := json.Marshal(userData)
	if err != nil {
		cfg.Logger.Error("Failed to marshal user data", "error", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	validationJSON := fmt.Sprintf(`{"base": %t, "advanced": %t}`, hasBase, hasAdvanced)

	http.SetCookie(w, &http.Cookie{
		Name:     "session_token",
		Value:    "authenticated_user",
		Path:     "/",
		Secure:   true,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
		MaxAge:   3600,
	})

	embedScript := fmt.Sprintf(`
        <script>
            window.userData = %s;
            window.validationState = %s;
        </script>
    `, userDataJSON, validationJSON)

	html = strings.Replace(html, "</body>", embedScript+"</body>", 1)

	w.Header().Set("Content-Type", "text/html")
	w.Write([]byte(html))
}

func safeBase64Decode(str string) string {
	if str == "" {
		return ""
	}
	if strings.Contains(str, "base64:") {
		str = strings.Split(str, "base64:")[1]
	}
	decoded, err := base64.StdEncoding.DecodeString(strings.TrimSpace(str))
	if err != nil {
		return ""
	}
	return string(decoded)
}
