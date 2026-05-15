package users

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	dbmanager "exodus/internal/db/manager"
	"exodus/internal/proto"
	srscore "exodus/internal/srslists"

	"github.com/iancoleman/orderedmap"
	"google.golang.org/grpc/codes"
)

type deployTaskPayload struct {
	Config       json.RawMessage         `json:"config"`
	Restart      *bool                   `json:"restart,omitempty"`
	ForceRestart *bool                   `json:"force_restart,omitempty"`
	SRSLists     []srscore.NodeSyncItem  `json:"srs_lists,omitempty"`
	Modules      *deployModulesTaskBlock `json:"modules,omitempty"`
}

type deployModulesTaskBlock struct {
	HaproxyEnabled bool                    `json:"haproxy_enabled"`
	HaproxyUsers   []deployHaproxyUserItem `json:"haproxy_users,omitempty"`
}

type deployHaproxyUserItem struct {
	Username       string `json:"username"`
	VLESSUUID      string `json:"vless_uuid"`
	TrojanPassword string `json:"trojan_password"`
}

type deployTarget struct {
	name   string
	uuid   string
	client proto.NodeServiceClient
}

type nodeInboundBinding struct {
	InboundUUID string
	Tag         string
}

type inboundUserCredentials struct {
	Username       string
	VLESSUUID      string
	TrojanPassword string
	SSPassword     string
}

func normalizeTagValue(tag string) string {
	return strings.TrimSpace(tag)
}

func (nm *NodeMonitor) deployToConnectedNodes(restart bool, forceRestart bool, requestedNodeUUIDs []string) {
	if nm == nil {
		return
	}
	nm.cfg.Logger.Info("Deploying node configs", "restart", restart, "force_restart", forceRestart, "requested_node_targets", len(requestedNodeUUIDs))

	dbNodes, err := nm.loadActiveNodes()
	if err != nil {
		nm.cfg.Logger.Warn("Failed to load nodes for deploy", "error", err)
		return
	}

	nodesByName := make(map[string]string, len(dbNodes))
	for _, n := range dbNodes {
		nodesByName[n.Name] = n.UUID
	}

	targetFilter := make(map[string]struct{}, len(requestedNodeUUIDs))
	for _, nodeUUID := range requestedNodeUUIDs {
		trimmed := strings.TrimSpace(nodeUUID)
		if trimmed == "" {
			continue
		}
		targetFilter[trimmed] = struct{}{}
	}

	targets := make([]deployTarget, 0)
	nm.nodesLock.RLock()
	for nodeName, state := range nm.nodes {
		if state == nil {
			continue
		}
		state.mutex.RLock()
		isReady := state.isConnected && state.client != nil
		client := state.client
		state.mutex.RUnlock()
		if !isReady {
			continue
		}
		nodeUUID, ok := nodesByName[nodeName]
		if !ok {
			continue
		}
		if len(targetFilter) > 0 {
			if _, allowed := targetFilter[nodeUUID]; !allowed {
				continue
			}
		}
		targets = append(targets, deployTarget{
			name:   nodeName,
			uuid:   nodeUUID,
			client: client,
		})
	}
	nm.nodesLock.RUnlock()
	nm.cfg.Logger.Debug("Prepared deploy targets", "active_nodes", len(dbNodes), "connected_targets", len(targets), "requested_node_targets", len(targetFilter), "restart", restart, "force_restart", forceRestart)

	if len(targets) == 0 {
		nm.cfg.Logger.Warn("No connected nodes to deploy")
		return
	}

	srsLists, srsErr := srscore.LoadNodeSyncItems(context.Background(), nm.manager)
	if srsErr != nil {
		nm.cfg.Logger.Warn("Failed to load SRS lists for deploy payload", "error", srsErr)
	}
	haproxyEnabled, modulesErr := nm.loadHaproxyModuleEnabled(context.Background())
	if modulesErr != nil {
		nm.cfg.Logger.Warn("Failed to load HAPROXY module settings for deploy payload", "error", modulesErr)
	}

	for _, target := range targets {
		nm.cfg.Logger.Debug("Building deploy config for node", "node", target.name, "node_uuid", target.uuid)
		configJSON, err := nm.buildNodeConfigForDeploy(nm.globalCtx, target.uuid)
		if err != nil {
			nm.cfg.Logger.Warn("Failed to build node deploy config", "node", target.name, "node_uuid", target.uuid, "error", err)
			continue
		}

		modules := &deployModulesTaskBlock{HaproxyEnabled: haproxyEnabled}
		if haproxyEnabled {
			haproxyUsers, usersErr := nm.loadNodeHaproxyUsers(nm.globalCtx, target.uuid)
			if usersErr != nil {
				nm.cfg.Logger.Warn("Failed to load node users for HAPROXY payload", "node", target.name, "node_uuid", target.uuid, "error", usersErr)
			} else {
				modules.HaproxyUsers = haproxyUsers
			}
		}

		restartFlag := restart
		forceRestartFlag := forceRestart
		taskPayload, err := json.Marshal(deployTaskPayload{
			Config:       configJSON,
			Restart:      &restartFlag,
			ForceRestart: &forceRestartFlag,
			SRSLists:     srsLists,
			Modules:      modules,
		})
		if err != nil {
			nm.cfg.Logger.Warn("Failed to serialize deploy payload", "node", target.name, "error", err)
			continue
		}

		ctxBase := nm.globalCtx
		if ctxBase == nil {
			ctxBase = context.Background()
		}
		ctx, cancel := context.WithTimeout(ctxBase, 60*time.Second)
		nm.cfg.Logger.Debug("Submitting deploy task", "node", target.name, "payload_bytes", len(taskPayload), "restart", restart, "force_restart", forceRestart)
		resp, err := target.client.SubmitTask(ctx, &proto.NodeTask{
			TaskId:    fmt.Sprintf("deploy-%d", time.Now().UnixNano()),
			Operation: "deploy_config",
			Payload:   taskPayload,
		})
		cancel()

		if err != nil {
			nm.cfg.Logger.Warn("Deploy task failed", "node", target.name, "error", err)
			continue
		}
		if resp == nil || resp.Code != int32(codes.OK) {
			if resp == nil {
				nm.cfg.Logger.Warn("Deploy task returned nil status", "node", target.name)
			} else {
				nm.cfg.Logger.Warn("Deploy task rejected", "node", target.name, "code", resp.Code, "message", resp.Message)
			}
			continue
		}

		nm.cfg.Logger.Info("Node config deployed", "node", target.name, "restart", restart, "force_restart", forceRestart, "message", resp.Message)
	}
}

func (nm *NodeMonitor) loadHaproxyModuleEnabled(ctx context.Context) (bool, error) {
	if ctx == nil {
		ctx = context.Background()
	}

	enabled := false
	err := nm.manager.ExecuteHighPriority(func(db dbmanager.DBExecutor) error {
		row := db.QueryRowContext(ctx, `SELECT haproxy_enabled FROM modules_settings WHERE id = 1`)
		if err := row.Scan(&enabled); err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return nil
			}
			return err
		}
		return nil
	})
	return enabled, err
}

func (nm *NodeMonitor) loadNodeHaproxyUsers(ctx context.Context, nodeUUID string) ([]deployHaproxyUserItem, error) {
	if strings.TrimSpace(nodeUUID) == "" {
		return nil, fmt.Errorf("node uuid is empty")
	}
	if ctx == nil {
		ctx = context.Background()
	}

	items := make([]deployHaproxyUserItem, 0)
	err := nm.manager.ExecuteHighPriority(func(db dbmanager.DBExecutor) error {
		rows, err := db.QueryContext(ctx, `
				SELECT
					u.username,
					CASE
						WHEN bool_or(lower(cpi.type) = 'vless') THEN u.vless_uuid::text
						ELSE ''
					END AS vless_uuid,
				CASE
					WHEN bool_or(lower(cpi.type) = 'trojan') THEN u.trojan_password
					ELSE ''
				END AS trojan_password
				FROM config_profile_inbounds_to_nodes cpitn
				JOIN config_profile_inbounds cpi ON cpi.uuid = cpitn.config_profile_inbound_uuid
				JOIN internal_squad_inbounds isi ON isi.inbound_uuid = cpitn.config_profile_inbound_uuid
				JOIN internal_squad_members ism ON ism.internal_squad_uuid = isi.internal_squad_uuid
				JOIN users u ON u.t_id = ism.user_id
			WHERE cpitn.node_uuid::text = ? AND u.status = 'ACTIVE'
				GROUP BY u.t_id, u.username, u.vless_uuid, u.trojan_password
				HAVING bool_or(lower(cpi.type) IN ('vless', 'trojan'))
				ORDER BY u.t_id ASC
		`, nodeUUID)
		if err != nil {
			return err
		}
		defer rows.Close()

		for rows.Next() {
			var item deployHaproxyUserItem
			if err := rows.Scan(&item.Username, &item.VLESSUUID, &item.TrojanPassword); err != nil {
				return err
			}
			item.Username = strings.TrimSpace(item.Username)
			if item.Username == "" {
				continue
			}
			items = append(items, item)
		}
		return rows.Err()
	})

	return items, err
}

func (nm *NodeMonitor) buildNodeConfigForDeploy(ctx context.Context, nodeUUID string) (json.RawMessage, error) {
	if strings.TrimSpace(nodeUUID) == "" {
		return nil, fmt.Errorf("node uuid is empty")
	}

	if ctx == nil {
		ctx = context.Background()
	}

	var profileConfig json.RawMessage
	bindings := make([]nodeInboundBinding, 0)

	err := nm.manager.ExecuteHighPriority(func(db dbmanager.DBExecutor) error {
		row := db.QueryRowContext(ctx, `
			SELECT cp.config
			FROM nodes n
			JOIN config_profiles cp ON cp.uuid = n.active_config_profile_uuid
			WHERE n.uuid = ? AND n.is_disabled = false
		`, nodeUUID)
		if err := row.Scan(&profileConfig); err != nil {
			return err
		}

		rows, err := db.QueryContext(ctx, `
			SELECT cpi.uuid, cpi.tag
			FROM config_profile_inbounds_to_nodes cpitn
			JOIN config_profile_inbounds cpi ON cpi.uuid = cpitn.config_profile_inbound_uuid
			WHERE cpitn.node_uuid = ?
		`, nodeUUID)
		if err != nil {
			return err
		}
		defer rows.Close()

		for rows.Next() {
			var item nodeInboundBinding
			if err := rows.Scan(&item.InboundUUID, &item.Tag); err != nil {
				return err
			}
			bindings = append(bindings, item)
		}
		return rows.Err()
	})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("node %s has no active config profile", nodeUUID)
		}
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

	err = nm.manager.ExecuteHighPriority(func(db dbmanager.DBExecutor) error {
		rows, err := db.QueryContext(ctx, `
			SELECT isi.inbound_uuid, u.username, u.vless_uuid, u.trojan_password, u.ss_password
			FROM internal_squad_inbounds isi
			JOIN internal_squad_members ism ON ism.internal_squad_uuid = isi.internal_squad_uuid
			JOIN users u ON u.t_id = ism.user_id
			WHERE isi.inbound_uuid = ANY(?) AND u.status = 'ACTIVE'
			ORDER BY u.t_id ASC
		`, inboundUUIDs)
		if err != nil {
			return err
		}
		defer rows.Close()

		for rows.Next() {
			var (
				inboundUUID string
				user        inboundUserCredentials
			)
			if err := rows.Scan(&inboundUUID, &user.Username, &user.VLESSUUID, &user.TrojanPassword, &user.SSPassword); err != nil {
				return err
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

		return rows.Err()
	})
	if err != nil {
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

	// Fallback guard: if selected inbound tags do not match config tags, keep the profile inbounds as-is
	// instead of dropping secure inbounds and deploying a broken runtime config.
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
	switch strings.ToLower(strings.TrimSpace(inboundType)) {
	case "vless":
		for _, user := range users {
			result = append(result, map[string]any{
				"name": user.Username,
				"uuid": user.VLESSUUID,
			})
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
	default:
		return []any{}
	}
	return result
}

func (nm *NodeMonitor) syncSRSListsToConnectedNodes() {
	if nm == nil {
		return
	}

	srsLists, err := srscore.LoadNodeSyncItems(context.Background(), nm.manager)
	if err != nil {
		nm.cfg.Logger.Warn("Failed to load SRS lists for node sync", "error", err)
		return
	}

	payload, err := json.Marshal(map[string]any{
		"srs_lists": srsLists,
	})
	if err != nil {
		nm.cfg.Logger.Warn("Failed to marshal SRS sync payload", "error", err)
		return
	}

	nm.nodesLock.RLock()
	targets := make([]deployTarget, 0, len(nm.nodes))
	for nodeName, state := range nm.nodes {
		if state == nil {
			continue
		}
		state.mutex.RLock()
		isReady := state.isConnected && state.client != nil
		client := state.client
		state.mutex.RUnlock()
		if !isReady {
			continue
		}
		targets = append(targets, deployTarget{name: nodeName, client: client})
	}
	nm.nodesLock.RUnlock()

	for _, target := range targets {
		ctxBase := nm.globalCtx
		if ctxBase == nil {
			ctxBase = context.Background()
		}
		ctx, cancel := context.WithTimeout(ctxBase, 30*time.Second)
		resp, submitErr := target.client.SubmitTask(ctx, &proto.NodeTask{
			TaskId:    fmt.Sprintf("sync-srs-%d", time.Now().UnixNano()),
			Operation: "sync_srs_lists",
			Payload:   payload,
		})
		cancel()

		if submitErr != nil {
			nm.cfg.Logger.Warn("SRS sync task failed", "node", target.name, "error", submitErr)
			continue
		}
		if resp == nil || resp.Code != int32(codes.OK) {
			if resp == nil {
				nm.cfg.Logger.Warn("SRS sync returned nil status", "node", target.name)
			} else {
				nm.cfg.Logger.Warn("SRS sync rejected", "node", target.name, "code", resp.Code, "message", resp.Message)
			}
			continue
		}

		nm.cfg.Logger.Info("SRS lists synced to node", "node", target.name, "lists", len(srsLists))
	}
}
