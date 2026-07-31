package users

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/iancoleman/orderedmap"
)

func (nm *NodeMonitor) loadNodePluginRuntimeConfig(ctx context.Context, nodeUUID string) (activeNodePluginRuntimeConfig, error) {
	if strings.TrimSpace(nodeUUID) == "" {
		return activeNodePluginRuntimeConfig{}, fmt.Errorf("node uuid is empty")
	}
	if ctx == nil {
		ctx = context.Background()
	}

	var pluginConfig activeNodePluginRuntimeConfig
	var rawConfig sql.NullString
	row := nm.db.QueryRowContext(ctx, `
		SELECT np.plugin_config::text
		FROM nodes n
		JOIN node_plugin np ON np.uuid = n.active_plugin_uuid
		WHERE n.uuid::text = $1
		LIMIT 1
	`, nodeUUID)
	if err := row.Scan(&rawConfig); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return pluginConfig, nil
		}
		return pluginConfig, err
	}
	if !rawConfig.Valid || strings.TrimSpace(rawConfig.String) == "" {
		return pluginConfig, nil
	}
	if err := json.Unmarshal([]byte(rawConfig.String), &pluginConfig); err != nil {
		return pluginConfig, err
	}
	return pluginConfig, nil
}

func (cfg activeNodePluginRuntimeConfig) sharedIPLists() map[string][]string {
	result := make(map[string][]string, len(cfg.SharedLists))
	for _, list := range cfg.SharedLists {
		if strings.TrimSpace(list.Type) != "ipList" {
			continue
		}
		name := strings.TrimSpace(list.Name)
		if name == "" {
			continue
		}
		result[name] = append([]string(nil), list.Items...)
	}
	return result
}

func resolvePluginIPRefs(raw []string, sharedLists map[string][]string) []string {
	result := make([]string, 0, len(raw))
	for _, item := range raw {
		value := strings.TrimSpace(item)
		if value == "" {
			continue
		}
		if strings.HasPrefix(value, "ext:") {
			result = append(result, sharedLists[value]...)
			continue
		}
		result = append(result, value)
	}
	return result
}

func normalizeStringSlice(raw []string) []string {
	seen := make(map[string]struct{}, len(raw))
	result := make([]string, 0, len(raw))
	for _, item := range raw {
		value := strings.TrimSpace(item)
		if value == "" {
			continue
		}
		if _, ok := seen[value]; ok {
			continue
		}
		seen[value] = struct{}{}
		result = append(result, value)
	}
	return result
}

func normalizePortSlice(raw []int) []int {
	seen := make(map[int]struct{}, len(raw))
	result := make([]int, 0, len(raw))
	for _, port := range raw {
		if port < 1 || port > 65535 {
			continue
		}
		if _, ok := seen[port]; ok {
			continue
		}
		seen[port] = struct{}{}
		result = append(result, port)
	}
	return result
}

func (nm *NodeMonitor) loadNodeHaproxyUsers(ctx context.Context, nodeUUID string, inboundTags []string) ([]deployHaproxyUserItem, bool, error) {
	if strings.TrimSpace(nodeUUID) == "" {
		return nil, false, fmt.Errorf("node uuid is empty")
	}
	if ctx == nil {
		ctx = context.Background()
	}

	inboundTags = normalizeHaproxyInboundTags(inboundTags)
	if len(inboundTags) == 0 {
		return nil, false, nil
	}

	matchAll := haproxyUsesAllInboundTags(inboundTags)
	tagFilterSQL := ""
	matchArgs := []any{nodeUUID}
	usersArgs := []any{nodeUUID}
	if !matchAll {
		tagFilterSQL = " AND cpi.tag = ANY($2)"
		matchArgs = append(matchArgs, inboundTags)
		usersArgs = append(usersArgs, inboundTags)
	}

	items := make([]deployHaproxyUserItem, 0)
	matched := false

	matchQuery := fmt.Sprintf(`
		SELECT EXISTS (
			SELECT 1
			FROM config_profile_inbounds_to_nodes cpitn
			JOIN config_profile_inbounds cpi ON cpi.uuid = cpitn.config_profile_inbound_uuid
			WHERE cpitn.node_uuid::text = $1
				AND lower(cpi.type) IN ('vless', 'trojan', 'naive', 'anytls')%s
		)`, tagFilterSQL)
	if err := nm.db.QueryRowContext(ctx, matchQuery, matchArgs...).Scan(&matched); err != nil {
		return nil, false, err
	}
	if !matched {
		return nil, false, nil
	}

	usersQuery := fmt.Sprintf(`
		SELECT
			u.username,
			CASE
				WHEN bool_or(lower(cpi.type) = 'vless') THEN COALESCE(u.vless_uuid::text, '')
				ELSE ''
			END AS vless_uuid,
			CASE
				WHEN bool_or(lower(cpi.type) = 'trojan') THEN COALESCE(u.trojan_password, '')
				ELSE ''
			END AS trojan_password,
			CASE
				WHEN bool_or(lower(cpi.type) = 'naive') THEN COALESCE(u.naive_password, '')
				ELSE ''
			END AS naive_password,
			CASE
				WHEN bool_or(lower(cpi.type) = 'anytls') THEN COALESCE(u.anytls_password, '')
				ELSE ''
			END AS anytls_password
		FROM config_profile_inbounds_to_nodes cpitn
		JOIN config_profile_inbounds cpi ON cpi.uuid = cpitn.config_profile_inbound_uuid
		JOIN internal_squad_inbounds isi ON isi.inbound_uuid = cpitn.config_profile_inbound_uuid
		JOIN internal_squad_members ism ON ism.internal_squad_uuid = isi.internal_squad_uuid
		JOIN users u ON u.t_id = ism.user_id
		WHERE cpitn.node_uuid::text = $1
			AND u.status = 'ACTIVE'
			AND lower(cpi.type) IN ('vless', 'trojan', 'naive', 'anytls')%s
		GROUP BY u.t_id, u.username, u.vless_uuid, u.trojan_password, u.naive_password, u.anytls_password
		ORDER BY u.t_id ASC
	`, tagFilterSQL)
	rows, err := nm.db.QueryContext(ctx, usersQuery, usersArgs...)
	if err != nil {
		return nil, false, err
	}
	defer rows.Close()

	for rows.Next() {
		var item deployHaproxyUserItem
		if err := rows.Scan(&item.Username, &item.VLESSUUID, &item.TrojanPassword, &item.NaivePassword, &item.AnytlsPassword); err != nil {
			return nil, false, err
		}
		item.Username = strings.TrimSpace(item.Username)
		if item.Username == "" {
			continue
		}
		items = append(items, item)
	}
	return items, matched, rows.Err()
}

func (nm *NodeMonitor) buildNodeConfigForDeploy(ctx context.Context, nodeUUID string) (json.RawMessage, error) {
	if strings.TrimSpace(nodeUUID) == "" {
		return nil, fmt.Errorf("node uuid is empty")
	}

	if ctx == nil {
		ctx = context.Background()
	}

	var profileConfig json.RawMessage
	row := nm.db.QueryRowContext(ctx, `
		SELECT cp.config
		FROM nodes n
		JOIN config_profiles cp ON cp.uuid = n.active_config_profile_uuid
		WHERE n.uuid = $1 AND n.is_disabled = false
	`, nodeUUID)
	if err := row.Scan(&profileConfig); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("node %s has no active config profile", nodeUUID)
		}
		return nil, err
	}

	rows, err := nm.db.QueryContext(ctx, `
		SELECT cpi.uuid, cpi.tag
		FROM config_profile_inbounds_to_nodes cpitn
		JOIN config_profile_inbounds cpi ON cpi.uuid = cpitn.config_profile_inbound_uuid
		WHERE cpitn.node_uuid = $1
	`, nodeUUID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	bindings := make([]nodeInboundBinding, 0)
	for rows.Next() {
		var item nodeInboundBinding
		if err := rows.Scan(&item.InboundUUID, &item.Tag); err != nil {
			return nil, err
		}
		bindings = append(bindings, item)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	if len(bindings) == 0 {
		return nil, fmt.Errorf("node %s has no active inbounds", nodeUUID)
	}

	bindingByInboundUUID := make(map[string]nodeInboundBinding, len(bindings))
	activeTags := make(map[string]struct{}, len(bindings))
	inboundUUIDs := make([]string, 0, len(bindings))
	for _, b := range bindings {
		bindingByInboundUUID[b.InboundUUID] = b
		normTag := normalizeTagValue(b.Tag)
		if normTag != "" {
			activeTags[normTag] = struct{}{}
		}
		inboundUUIDs = append(inboundUUIDs, b.InboundUUID)
	}

	usersByTag := make(map[string][]inboundUserCredentials, len(activeTags))
	dedup := make(map[string]map[string]struct{}, len(activeTags))

	userRows, err := nm.db.QueryContext(ctx, `
		SELECT
			isi.inbound_uuid,
			u.username,
			COALESCE(u.vless_uuid::text, ''),
			COALESCE(u.trojan_password, ''),
			COALESCE(u.ss_password, ''),
			COALESCE(u.naive_password, ''),
			COALESCE(u.shadowtls_password, ''),
			COALESCE(u.hysteria2_password, ''),
			COALESCE(u.anytls_password, '')
		FROM internal_squad_inbounds isi
		JOIN internal_squad_members ism ON ism.internal_squad_uuid = isi.internal_squad_uuid
		JOIN users u ON u.t_id = ism.user_id
		WHERE isi.inbound_uuid = ANY($1) AND u.status = 'ACTIVE'
		ORDER BY u.t_id ASC
	`, inboundUUIDs)
	if err != nil {
		return nil, err
	}
	defer userRows.Close()

	for userRows.Next() {
		var (
			inboundUUID string
			user        inboundUserCredentials
		)
		if err := userRows.Scan(
			&inboundUUID,
			&user.Username,
			&user.VLESSUUID,
			&user.TrojanPassword,
			&user.SSPassword,
			&user.NaivePassword,
			&user.ShadowTLSPass,
			&user.Hysteria2Pass,
			&user.AnytlsPassword,
		); err != nil {
			return nil, err
		}
		binding, ok := bindingByInboundUUID[inboundUUID]
		tag := normalizeTagValue(binding.Tag)
		if !ok || tag == "" || strings.TrimSpace(user.Username) == "" {
			continue
		}

		if dedup[tag] == nil {
			dedup[tag] = make(map[string]struct{})
		}
		if _, exists := dedup[tag][user.Username]; exists {
			continue
		}
		dedup[tag][user.Username] = struct{}{}
		usersByTag[tag] = append(usersByTag[tag], user)
	}
	if err := userRows.Err(); err != nil {
		return nil, err
	}

	parsed := orderedmap.New()
	if err := json.Unmarshal(profileConfig, parsed); err != nil {
		return nil, fmt.Errorf("invalid profile config json: %w", err)
	}

	rawInboundsRaw, ok := parsed.Get("inbounds")
	if !ok {
		return nil, fmt.Errorf("profile config has no valid inbounds array")
	}
	rawInbounds, ok := rawInboundsRaw.([]any)
	if !ok {
		return nil, fmt.Errorf("profile config has no valid inbounds array")
	}

	matchedActiveTags := 0
	for _, raw := range rawInbounds {
		tag := getFieldString(raw, "tag")
		if _, isActiveTag := activeTags[normalizeTagValue(tag)]; isActiveTag {
			matchedActiveTags++
		}
	}

	useFallbackKeepAll := matchedActiveTags == 0 && len(activeTags) > 0
	if useFallbackKeepAll {
		nm.cfg.Logger.Warn("No selected inbound tags matched config inbounds; keeping all profile inbounds", "node_uuid", nodeUUID, "selected_tags", len(activeTags), "config_inbounds", len(rawInbounds))
	}

	filteredInbounds := make([]any, 0, len(rawInbounds))
	for _, raw := range rawInbounds {
		tag := getFieldString(raw, "tag")
		normTag := normalizeTagValue(tag)
		inboundType := normalizeInboundType(raw)
		_, isActiveTag := activeTags[normTag]

		if !useFallbackKeepAll && !isActiveTag && !isUnsecureInbound(inboundType) {
			continue
		}
		if isActiveTag {
			raw = setField(raw, "users", buildInboundUsers(inboundType, usersByTag[normTag]))
		}

		filteredInbounds = append(filteredInbounds, raw)
	}

	parsed.Set("inbounds", filteredInbounds)

	finalConfig, err := json.Marshal(parsed)
	if err != nil {
		return nil, fmt.Errorf("marshal deploy config: %w", err)
	}
	return finalConfig, nil
}

func normalizeInboundType(inbound any) string {
	if value := getFieldString(inbound, "type"); strings.TrimSpace(value) != "" {
		return strings.ToLower(strings.TrimSpace(value))
	}
	if value := getFieldString(inbound, "protocol"); strings.TrimSpace(value) != "" {
		return strings.ToLower(strings.TrimSpace(value))
	}
	return ""
}

func getField(v any, key string) (any, bool) {
	switch m := v.(type) {
	case map[string]any:
		value, ok := m[key]
		return value, ok
	case orderedmap.OrderedMap:
		return m.Get(key)
	case *orderedmap.OrderedMap:
		if m == nil {
			return nil, false
		}
		return m.Get(key)
	default:
		return nil, false
	}
}

func getFieldString(v any, key string) string {
	value, ok := getField(v, key)
	if !ok {
		return ""
	}
	s, _ := value.(string)
	return s
}

func setField(v any, key string, value any) any {
	switch m := v.(type) {
	case map[string]any:
		m[key] = value
		return m
	case orderedmap.OrderedMap:
		m.Set(key, value)
		return m
	case *orderedmap.OrderedMap:
		if m == nil {
			n := orderedmap.New()
			n.Set(key, value)
			return *n
		}
		m.Set(key, value)
		return *m
	default:
		return v
	}
}

func isUnsecureInbound(inboundType string) bool {
	switch inboundType {
	case "dokodemo-door", "http", "mixed", "wireguard":
		return true
	default:
		return false
	}
}

func buildInboundUsers(inboundType string, users []inboundUserCredentials) []any {
	result := make([]any, 0, len(users))
	normalizedType := strings.ToLower(strings.TrimSpace(inboundType))
	if normalizedType == "ss" {
		normalizedType = "shadowsocks"
	}
	if normalizedType == "hy2" {
		normalizedType = "hysteria2"
	}
	switch normalizedType {
	case "vless", "vmess":
		for _, user := range users {
			item := map[string]any{
				"name": user.Username,
				"uuid": user.VLESSUUID,
			}
			if normalizedType == "vmess" {
				item["alterId"] = 0
			}
			result = append(result, item)
		}
	case "trojan":
		for _, user := range users {
			result = append(result, map[string]any{
				"name":     user.Username,
				"password": user.TrojanPassword,
			})
		}
	case "shadowsocks":
		for _, user := range users {
			result = append(result, map[string]any{
				"name":     user.Username,
				"password": user.SSPassword,
			})
		}
	case "naive":
		for _, user := range users {
			result = append(result, map[string]any{
				"username": user.Username,
				"password": user.NaivePassword,
			})
		}
	case "anytls":
		for _, user := range users {
			result = append(result, map[string]any{
				"name":     user.Username,
				"password": user.AnytlsPassword,
			})
		}
	case "shadowtls":
		for _, user := range users {
			result = append(result, map[string]any{
				"name":     user.Username,
				"password": user.ShadowTLSPass,
			})
		}
	case "hysteria":
		for _, user := range users {
			result = append(result, map[string]any{
				"name":     user.Username,
				"auth_str": user.Hysteria2Pass,
			})
		}
	case "hysteria2":
		for _, user := range users {
			result = append(result, map[string]any{
				"name":     user.Username,
				"password": user.Hysteria2Pass,
			})
		}
	case "tuic":
		for _, user := range users {
			result = append(result, map[string]any{
				"name":     user.Username,
				"uuid":     user.VLESSUUID,
				"password": user.TrojanPassword,
			})
		}
	default:
		return []any{}
	}
	return result
}
