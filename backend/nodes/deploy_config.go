package users

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	dbmanager "v2ray-stat/backend/db/manager"
	"v2ray-stat/proto"

	"google.golang.org/grpc/codes"
)

type deployTaskPayload struct {
	Config  json.RawMessage `json:"config"`
	Restart *bool           `json:"restart,omitempty"`
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

func (nm *NodeMonitor) deployToConnectedNodes(restart bool) {
	if nm == nil {
		return
	}
	nm.cfg.Logger.Info("Deploying node configs", "restart", restart)

	dbNodes, err := nm.loadActiveNodes()
	if err != nil {
		nm.cfg.Logger.Warn("Failed to load nodes for deploy", "error", err)
		return
	}

	nodesByName := make(map[string]string, len(dbNodes))
	for _, n := range dbNodes {
		nodesByName[n.Name] = n.UUID
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
		targets = append(targets, deployTarget{
			name:   nodeName,
			uuid:   nodeUUID,
			client: client,
		})
	}
	nm.nodesLock.RUnlock()
	nm.cfg.Logger.Debug("Prepared deploy targets", "active_nodes", len(dbNodes), "connected_targets", len(targets), "restart", restart)

	if len(targets) == 0 {
		nm.cfg.Logger.Warn("No connected nodes to deploy")
		return
	}

	for _, target := range targets {
		nm.cfg.Logger.Debug("Building deploy config for node", "node", target.name, "node_uuid", target.uuid)
		configJSON, err := nm.buildNodeConfigForDeploy(nm.globalCtx, target.uuid)
		if err != nil {
			nm.cfg.Logger.Warn("Failed to build node deploy config", "node", target.name, "node_uuid", target.uuid, "error", err)
			continue
		}

		restartFlag := restart
		taskPayload, err := json.Marshal(deployTaskPayload{
			Config:  configJSON,
			Restart: &restartFlag,
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
		nm.cfg.Logger.Debug("Submitting deploy task", "node", target.name, "payload_bytes", len(taskPayload), "restart", restart)
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

		nm.cfg.Logger.Info("Node config deployed", "node", target.name, "restart", restart, "message", resp.Message)
	}
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
		activeTags[b.Tag] = struct{}{}
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
			if !ok || strings.TrimSpace(binding.Tag) == "" || strings.TrimSpace(user.Username) == "" {
				continue
			}

			if dedup[binding.Tag] == nil {
				dedup[binding.Tag] = make(map[string]struct{})
			}
			if _, exists := dedup[binding.Tag][user.Username]; exists {
				continue
			}
			dedup[binding.Tag][user.Username] = struct{}{}
			usersByTag[binding.Tag] = append(usersByTag[binding.Tag], user)
		}

		return rows.Err()
	})
	if err != nil {
		return nil, err
	}

	var parsed map[string]any
	if err := json.Unmarshal(profileConfig, &parsed); err != nil {
		return nil, fmt.Errorf("invalid profile config json: %w", err)
	}

	rawInbounds, ok := parsed["inbounds"].([]any)
	if !ok {
		return nil, fmt.Errorf("profile config has no valid inbounds array")
	}

	filteredInbounds := make([]any, 0, len(rawInbounds))
	for _, raw := range rawInbounds {
		inbound, ok := raw.(map[string]any)
		if !ok {
			continue
		}

		tag, _ := inbound["tag"].(string)
		inboundType := normalizeInboundType(inbound)
		_, isActiveTag := activeTags[tag]

		if !isActiveTag && !isUnsecureInbound(inboundType) {
			continue
		}
		if isActiveTag {
			inbound["users"] = buildInboundUsers(inboundType, usersByTag[tag])
		}

		filteredInbounds = append(filteredInbounds, inbound)
	}

	parsed["inbounds"] = filteredInbounds

	finalConfig, err := json.Marshal(parsed)
	if err != nil {
		return nil, fmt.Errorf("marshal deploy config: %w", err)
	}
	return finalConfig, nil
}

func normalizeInboundType(inbound map[string]any) string {
	if inbound == nil {
		return ""
	}
	if value, ok := inbound["type"].(string); ok && strings.TrimSpace(value) != "" {
		return strings.ToLower(strings.TrimSpace(value))
	}
	if value, ok := inbound["protocol"].(string); ok && strings.TrimSpace(value) != "" {
		return strings.ToLower(strings.TrimSpace(value))
	}
	return ""
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
