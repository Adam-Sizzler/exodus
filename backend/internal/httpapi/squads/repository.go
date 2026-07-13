package squads

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	dbmanager "exodus/internal/db/manager"
	"exodus/internal/dbutil"
)

type SquadRepository struct {
	manager *dbmanager.DatabaseManager
}

func NewSquadRepository(manager *dbmanager.DatabaseManager) *SquadRepository {
	return &SquadRepository{manager: manager}
}

func (r *SquadRepository) getSquads(ctx context.Context) ([]InternalSquad, error) {
	var squads []InternalSquad
	err := r.manager.ExecuteHighPriority(func(db dbmanager.DBExecutor) error {
		query := `
			SELECT uuid, view_position, name, created_at, updated_at
			FROM internal_squads
			ORDER BY view_position ASC, name ASC`

		rows, err := db.QueryContext(ctx, query)
		if err != nil {
			return err
		}
		defer rows.Close()

		for rows.Next() {
			squad, err := scanInternalSquad(rows)
			if err != nil {
				return err
			}
			squads = append(squads, squad)
		}
		return rows.Err()
	})
	return squads, err
}

func (r *SquadRepository) getSquadByUUID(ctx context.Context, squadUUID string) (InternalSquad, error) {
	var squad InternalSquad
	err := r.manager.ExecuteHighPriority(func(db dbmanager.DBExecutor) error {
		query := `SELECT uuid, view_position, name, created_at, updated_at
				  FROM internal_squads WHERE uuid = ?`
		row := db.QueryRowContext(ctx, query, squadUUID)
		var scanErr error
		squad, scanErr = scanInternalSquad(row)
		return scanErr
	})
	return squad, err
}

func (r *SquadRepository) getSquadMembersCount(ctx context.Context, squadUUID string) (int, error) {
	var count int
	err := r.manager.ExecuteHighPriority(func(db dbmanager.DBExecutor) error {
		return db.QueryRowContext(ctx,
			`SELECT COUNT(*) FROM internal_squad_members WHERE internal_squad_uuid = ?`,
			squadUUID).Scan(&count)
	})
	return count, err
}

func (r *SquadRepository) getSquadInbounds(ctx context.Context, squadUUID string) ([]InternalSquadInboundAPI, error) {
	var inbounds []InternalSquadInboundAPI
	err := r.manager.ExecuteHighPriority(func(db dbmanager.DBExecutor) error {
		rows, err := db.QueryContext(ctx, `
			SELECT cpi.uuid, cpi.profile_uuid, cpi.tag, cpi.type, cpi.network, cpi.security, cpi.port, cpi.raw_inbound
			FROM internal_squad_inbounds isi
			JOIN config_profile_inbounds cpi ON cpi.uuid = isi.inbound_uuid
			WHERE isi.internal_squad_uuid = ?
			ORDER BY cpi.tag ASC
		`, squadUUID)
		if err != nil {
			return err
		}
		defer rows.Close()

		for rows.Next() {
			var inbound InternalSquadInboundAPI
			var network, security, rawInbound sql.NullString
			var port sql.NullInt64
			if scanErr := rows.Scan(
				&inbound.UUID,
				&inbound.ProfileUUID,
				&inbound.Tag,
				&inbound.Type,
				&network,
				&security,
				&port,
				&rawInbound,
			); scanErr != nil {
				return scanErr
			}
			if network.Valid {
				inbound.Network = &network.String
			}
			if security.Valid {
				inbound.Security = &security.String
			}
			if port.Valid {
				p := int(port.Int64)
				inbound.Port = &p
			}
			if rawInbound.Valid {
				inbound.RawInbound = json.RawMessage(rawInbound.String)
			}
			inbounds = append(inbounds, inbound)
		}
		return rows.Err()
	})
	return inbounds, err
}

func (r *SquadRepository) getSquadAccessibleNodes(ctx context.Context, squadUUID string) ([]InternalSquadAccessibleNode, error) {
	nodes := make([]InternalSquadAccessibleNode, 0)
	err := r.manager.ExecuteHighPriority(func(db dbmanager.DBExecutor) error {
		var exists int
		if err := db.QueryRowContext(ctx, `SELECT 1 FROM internal_squads WHERE uuid = ?`, squadUUID).Scan(&exists); err != nil {
			if err == sql.ErrNoRows {
				return errInternalSquadNotFound
			}
			return err
		}

		rows, err := db.QueryContext(ctx, `
			SELECT
				n.uuid,
				n.name,
				n.country_code,
				cp.uuid,
				cp.name,
				cpi.tag
			FROM internal_squad_inbounds isi
			INNER JOIN config_profile_inbounds cpi ON cpi.uuid = isi.inbound_uuid
			INNER JOIN config_profiles cp ON cp.uuid = cpi.profile_uuid
			INNER JOIN config_profile_inbounds_to_nodes cpin
				ON cpin.config_profile_inbound_uuid = cpi.uuid
			INNER JOIN nodes n
				ON n.uuid = cpin.node_uuid
				AND n.active_config_profile_uuid = cp.uuid
			WHERE isi.internal_squad_uuid = ?
			ORDER BY n.view_position ASC, n.name ASC, cpi.tag ASC
		`, squadUUID)
		if err != nil {
			return err
		}
		defer rows.Close()

		indexByNode := make(map[string]int)
		inboundSeenByNode := make(map[string]map[string]bool)
		for rows.Next() {
			var nodeUUID, nodeName, countryCode, profileUUID, profileName, inboundTag string
			if err := rows.Scan(&nodeUUID, &nodeName, &countryCode, &profileUUID, &profileName, &inboundTag); err != nil {
				return err
			}
			idx, ok := indexByNode[nodeUUID]
			if !ok {
				nodes = append(nodes, InternalSquadAccessibleNode{
					UUID:              nodeUUID,
					NodeName:          nodeName,
					CountryCode:       countryCode,
					ConfigProfileUUID: profileUUID,
					ConfigProfileName: profileName,
					ActiveInbounds:    make([]string, 0),
				})
				idx = len(nodes) - 1
				indexByNode[nodeUUID] = idx
				inboundSeenByNode[nodeUUID] = make(map[string]bool)
			}
			if !inboundSeenByNode[nodeUUID][inboundTag] {
				nodes[idx].ActiveInbounds = append(nodes[idx].ActiveInbounds, inboundTag)
				inboundSeenByNode[nodeUUID][inboundTag] = true
			}
		}
		return rows.Err()
	})
	return nodes, err
}

func (r *SquadRepository) getConfigProfilesWithInbounds(ctx context.Context) ([]ConfigProfileWithInbounds, error) {
	var profiles []ConfigProfileWithInbounds
	err := r.manager.ExecuteHighPriority(func(db dbmanager.DBExecutor) error {
		profileQuery := `
			SELECT uuid, view_position, name, config, created_at, updated_at
			FROM config_profiles
			ORDER BY view_position ASC, name ASC`

		profileRows, err := db.QueryContext(ctx, profileQuery)
		if err != nil {
			return err
		}
		defer profileRows.Close()

		type tempProfile struct {
			UUID         string
			ViewPosition int
			Name         string
			Config       json.RawMessage
			CreatedAt    time.Time
			UpdatedAt    time.Time
		}
		var tempProfiles []tempProfile

		for profileRows.Next() {
			var tp tempProfile
			var viewPosition sql.NullInt64
			var configStr sql.NullString

			if err := profileRows.Scan(&tp.UUID, &viewPosition, &tp.Name, &configStr, &tp.CreatedAt, &tp.UpdatedAt); err != nil {
				return err
			}

			if viewPosition.Valid {
				tp.ViewPosition = int(viewPosition.Int64)
			}
			if configStr.Valid {
				tp.Config = json.RawMessage(configStr.String)
			}

			tempProfiles = append(tempProfiles, tp)
		}
		profileRows.Close()

		inboundQuery := `
			SELECT uuid, profile_uuid, tag, type, network, security, port, raw_inbound
			FROM config_profile_inbounds
			ORDER BY profile_uuid, tag`

		inboundRows, err := db.QueryContext(ctx, inboundQuery)
		if err != nil {
			return err
		}
		defer inboundRows.Close()

		inboundsByProfile := make(map[string][]InboundInfo)
		for inboundRows.Next() {
			var ib InboundInfo
			var profileUUID string
			var network, security, rawInbound sql.NullString
			var port sql.NullInt64

			if err := inboundRows.Scan(&ib.UUID, &profileUUID, &ib.Tag, &ib.Type, &network, &security, &port, &rawInbound); err != nil {
				return err
			}

			if network.Valid {
				ib.Network = &network.String
			}
			if security.Valid {
				ib.Security = &security.String
			}
			if port.Valid {
				p := int(port.Int64)
				ib.Port = &p
			}
			if rawInbound.Valid {
				ib.RawInbound = json.RawMessage(rawInbound.String)
			}

			inboundsByProfile[profileUUID] = append(inboundsByProfile[profileUUID], ib)
		}
		inboundRows.Close()

		for _, tp := range tempProfiles {
			ibs := inboundsByProfile[tp.UUID]
			if ibs == nil {
				ibs = make([]InboundInfo, 0)
			}
			profiles = append(profiles, ConfigProfileWithInbounds{
				UUID:         tp.UUID,
				Name:         tp.Name,
				ViewPosition: tp.ViewPosition,
				Config:       tp.Config,
				Inbounds:     ibs,
				CreatedAt:    tp.CreatedAt,
				UpdatedAt:    tp.UpdatedAt,
			})
		}

		return nil
	})
	return profiles, err
}

func (r *SquadRepository) getInboundAssignments(ctx context.Context, nodeUUID string) ([]ConfigProfileInboundToNode, error) {
	var assignments []ConfigProfileInboundToNode
	err := r.manager.ExecuteHighPriority(func(db dbmanager.DBExecutor) error {
		var query string
		var args []interface{}

		if nodeUUID != "" {
			query = `
				SELECT config_profile_inbound_uuid, node_uuid, created_at
				FROM config_profile_inbounds_to_nodes
				WHERE node_uuid = ?
				ORDER BY config_profile_inbound_uuid`
			args = []interface{}{nodeUUID}
		} else {
			query = `
				SELECT config_profile_inbound_uuid, node_uuid, created_at
				FROM config_profile_inbounds_to_nodes
				ORDER BY node_uuid, config_profile_inbound_uuid`
		}

		rows, err := db.QueryContext(ctx, query, args...)
		if err != nil {
			return err
		}
		defer rows.Close()

		for rows.Next() {
			var a ConfigProfileInboundToNode
			if err := rows.Scan(&a.ConfigProfileInboundUUID, &a.NodeUUID, &a.CreatedAt); err != nil {
				return err
			}
			assignments = append(assignments, a)
		}
		return rows.Err()
	})
	return assignments, err
}

func (r *SquadRepository) createSquad(ctx context.Context, squadUUID string, req InternalSquadCreateRequest) error {
	return r.manager.ExecuteHighPriority(func(db dbmanager.DBExecutor) error {
		query := `
			INSERT INTO internal_squads (
				uuid, view_position, name, created_at, updated_at
			) VALUES (?, ?, ?, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)`

		_, err := db.ExecContext(ctx, query, squadUUID, req.ViewPosition, req.Name)
		return err
	})
}

func (r *SquadRepository) updateSquad(ctx context.Context, squadUUID string, clauses []string, args []any, inboundUUIDs []string) error {
	return r.manager.ExecuteHighPriority(func(db dbmanager.DBExecutor) error {
		tx, err := db.BeginTx(ctx, nil)
		if err != nil {
			return err
		}

		if len(clauses) > 0 {
			args = append(args, squadUUID)
			query := fmt.Sprintf("UPDATE internal_squads SET %s, updated_at = CURRENT_TIMESTAMP WHERE uuid = ?", strings.Join(clauses, ", "))
			result, err := tx.ExecContext(ctx, query, args...)
			if err != nil {
				_ = tx.Rollback()
				return err
			}
			rows, err := result.RowsAffected()
			if err != nil {
				_ = tx.Rollback()
				return err
			}
			if rows == 0 {
				_ = tx.Rollback()
				return sql.ErrNoRows
			}
		}

		if inboundUUIDs != nil {
			if _, err := tx.ExecContext(ctx, `DELETE FROM internal_squad_inbounds WHERE internal_squad_uuid = ?`, squadUUID); err != nil {
				_ = tx.Rollback()
				return err
			}

			seen := make(map[string]struct{}, len(inboundUUIDs))
			for _, inboundUUID := range inboundUUIDs {
				cleanInboundUUID := strings.TrimSpace(inboundUUID)
				if cleanInboundUUID == "" {
					continue
				}
				if _, ok := seen[cleanInboundUUID]; ok {
					continue
				}
				seen[cleanInboundUUID] = struct{}{}

				var inboundExists int
				if err := tx.QueryRowContext(ctx, `SELECT 1 FROM config_profile_inbounds WHERE uuid = ?`, cleanInboundUUID).Scan(&inboundExists); err != nil {
					_ = tx.Rollback()
					if err == sql.ErrNoRows {
						return fmt.Errorf("inbound not found")
					}
					return err
				}

				if _, err := tx.ExecContext(
					ctx,
					`INSERT INTO internal_squad_inbounds (internal_squad_uuid, inbound_uuid) VALUES (?, ?)`,
					squadUUID,
					cleanInboundUUID,
				); err != nil {
					_ = tx.Rollback()
					return err
				}
			}
		}

		return tx.Commit()
	})
}

func (r *SquadRepository) deleteSquad(ctx context.Context, squadUUID string) (string, error) {
	var squadName string
	err := r.manager.ExecuteHighPriority(func(db dbmanager.DBExecutor) error {
		if err := db.QueryRowContext(ctx, "SELECT name FROM internal_squads WHERE uuid = ?", squadUUID).Scan(&squadName); err != nil {
			return err
		}
		_, err := db.ExecContext(ctx, "DELETE FROM internal_squads WHERE uuid = ?", squadUUID)
		return err
	})
	return squadName, err
}

func (r *SquadRepository) reorderSquads(ctx context.Context, items []reorderSquadItem) error {
	return r.manager.ExecuteHighPriority(func(db dbmanager.DBExecutor) error {
		tx, err := db.BeginTx(ctx, nil)
		if err != nil {
			return err
		}
		for _, item := range items {
			if _, err := tx.ExecContext(ctx, `UPDATE internal_squads SET view_position = ? WHERE uuid = ?`, item.ViewPosition, item.UUID); err != nil {
				_ = tx.Rollback()
				return err
			}
		}
		if _, err := tx.ExecContext(ctx, `SELECT setval('internal_squads_view_position_seq', (SELECT COALESCE(MAX(view_position), 0) FROM internal_squads) + 1)`); err != nil {
			_ = tx.Rollback()
			return err
		}
		return tx.Commit()
	})
}

func (r *SquadRepository) setInboundAssignments(ctx context.Context, nodeUUID string, inboundUUIDs []string) error {
	return r.manager.ExecuteHighPriority(func(db dbmanager.DBExecutor) error {
		tx, err := db.BeginTx(ctx, nil)
		if err != nil {
			return err
		}

		if _, err := tx.ExecContext(ctx, `DELETE FROM config_profile_inbounds_to_nodes WHERE node_uuid = ?`, nodeUUID); err != nil {
			_ = tx.Rollback()
			return err
		}

		for _, inboundUUID := range inboundUUIDs {
			_, err := tx.ExecContext(ctx, `
				INSERT INTO config_profile_inbounds_to_nodes (config_profile_inbound_uuid, node_uuid)
				VALUES (?, ?)`, inboundUUID, nodeUUID)
			if err != nil {
				_ = tx.Rollback()
				return err
			}
		}

		return tx.Commit()
	})
}

func (r *SquadRepository) getSquadInboundBindings(ctx context.Context, squadUUID string) ([]InternalSquadInbound, error) {
	var squadInbounds []InternalSquadInbound
	err := r.manager.ExecuteHighPriority(func(db dbmanager.DBExecutor) error {
		var query string
		var args []interface{}

		if squadUUID != "" {
			query = `
				SELECT internal_squad_uuid, inbound_uuid
				FROM internal_squad_inbounds
				WHERE internal_squad_uuid = ?
				ORDER BY inbound_uuid`
			args = []interface{}{squadUUID}
		} else {
			query = `
				SELECT internal_squad_uuid, inbound_uuid
				FROM internal_squad_inbounds
				ORDER BY internal_squad_uuid, inbound_uuid`
		}

		rows, err := db.QueryContext(ctx, query, args...)
		if err != nil {
			return err
		}
		defer rows.Close()

		for rows.Next() {
			var si InternalSquadInbound
			if err := rows.Scan(&si.InternalSquadUUID, &si.InboundUUID); err != nil {
				return err
			}
			squadInbounds = append(squadInbounds, si)
		}
		return rows.Err()
	})
	return squadInbounds, err
}

func (r *SquadRepository) setSquadInbounds(ctx context.Context, squadUUID string, inboundUUIDs []string) error {
	return r.manager.ExecuteHighPriority(func(db dbmanager.DBExecutor) error {
		tx, err := db.BeginTx(ctx, nil)
		if err != nil {
			return err
		}

		var squadID string
		err = tx.QueryRowContext(ctx, "SELECT uuid FROM internal_squads WHERE uuid = ?", squadUUID).Scan(&squadID)
		if err != nil {
			_ = tx.Rollback()
			if err == sql.ErrNoRows {
				return fmt.Errorf("squad not found")
			}
			return err
		}

		_, err = tx.ExecContext(ctx, "DELETE FROM internal_squad_inbounds WHERE internal_squad_uuid = ?", squadUUID)
		if err != nil {
			_ = tx.Rollback()
			return fmt.Errorf("failed to clear existing inbounds: %w", err)
		}

		for _, inboundUUID := range inboundUUIDs {
			var inboundID string
			err := tx.QueryRowContext(ctx, "SELECT uuid FROM config_profile_inbounds WHERE uuid = ?", inboundUUID).Scan(&inboundID)
			if err != nil {
				_ = tx.Rollback()
				if err == sql.ErrNoRows {
					return fmt.Errorf("inbound not found: %s", inboundUUID)
				}
				return err
			}

			_, err = tx.ExecContext(ctx,
				"INSERT INTO internal_squad_inbounds (internal_squad_uuid, inbound_uuid) VALUES (?, ?)",
				squadUUID, inboundUUID)
			if err != nil {
				_ = tx.Rollback()
				return fmt.Errorf("failed to insert inbound: %w", err)
			}
		}

		return tx.Commit()
	})
}

func (r *SquadRepository) getSquadMembers(ctx context.Context, squadUUID string) ([]InternalSquadMember, error) {
	var squadMembers []InternalSquadMember
	err := r.manager.ExecuteHighPriority(func(db dbmanager.DBExecutor) error {
		var query string
		var args []interface{}

		if squadUUID != "" {
			query = `
				SELECT m.internal_squad_uuid, m.user_id, u.username
				FROM internal_squad_members m
				JOIN users u ON m.user_id = u.t_id
				WHERE m.internal_squad_uuid = ?
				ORDER BY m.user_id`
			args = []interface{}{squadUUID}
		} else {
			query = `
				SELECT m.internal_squad_uuid, m.user_id, u.username
				FROM internal_squad_members m
				JOIN users u ON m.user_id = u.t_id
				ORDER BY m.internal_squad_uuid, m.user_id`
		}

		rows, err := db.QueryContext(ctx, query, args...)
		if err != nil {
			return err
		}
		defer rows.Close()

		for rows.Next() {
			var sm InternalSquadMember
			if err := rows.Scan(&sm.InternalSquadUUID, &sm.UserID, &sm.Username); err != nil {
				return err
			}
			squadMembers = append(squadMembers, sm)
		}
		return rows.Err()
	})
	return squadMembers, err
}

func (r *SquadRepository) setSquadMembers(ctx context.Context, squadUUID string, userIDs []int64) error {
	return r.manager.ExecuteHighPriority(func(db dbmanager.DBExecutor) error {
		tx, err := db.BeginTx(ctx, nil)
		if err != nil {
			return err
		}

		var squadID string
		err = tx.QueryRowContext(ctx, "SELECT uuid FROM internal_squads WHERE uuid = ?", squadUUID).Scan(&squadID)
		if err != nil {
			_ = tx.Rollback()
			if err == sql.ErrNoRows {
				return fmt.Errorf("squad not found")
			}
			return err
		}

		_, err = tx.ExecContext(ctx, "DELETE FROM internal_squad_members WHERE internal_squad_uuid = ?", squadUUID)
		if err != nil {
			_ = tx.Rollback()
			return fmt.Errorf("failed to clear existing members: %w", err)
		}

		for _, userID := range userIDs {
			var exists bool
			err := tx.QueryRowContext(ctx, "SELECT EXISTS(SELECT 1 FROM users WHERE t_id = ?)", userID).Scan(&exists)
			if err != nil {
				_ = tx.Rollback()
				return err
			}
			if !exists {
				_ = tx.Rollback()
				return fmt.Errorf("user not found: %d", userID)
			}

			_, err = tx.ExecContext(ctx,
				"INSERT INTO internal_squad_members (internal_squad_uuid, user_id) VALUES (?, ?)",
				squadUUID, userID)
			if err != nil {
				_ = tx.Rollback()
				return fmt.Errorf("failed to insert member: %w", err)
			}
		}

		return tx.Commit()
	})
}

func (r *SquadRepository) getSquadDetails(ctx context.Context, squadUUID string) (SquadDetails, error) {
	var squad SquadDetails
	err := r.manager.ExecuteHighPriority(func(db dbmanager.DBExecutor) error {
		query := `
			SELECT uuid, view_position, name, created_at, updated_at
			FROM internal_squads
			WHERE uuid = ?`

		var viewPosition sql.NullInt64
		err := db.QueryRowContext(ctx, query, squadUUID).Scan(
			&squad.UUID, &viewPosition, &squad.Name, &squad.CreatedAt, &squad.UpdatedAt)
		if err != nil {
			if err == sql.ErrNoRows {
				return fmt.Errorf("squad not found")
			}
			return err
		}

		if viewPosition.Valid {
			squad.ViewPosition = int(viewPosition.Int64)
		}

		inboundQuery := `
			SELECT i.uuid, i.tag, i.type, i.network, i.security, i.port, i.raw_inbound
			FROM internal_squad_inbounds si
			JOIN config_profile_inbounds i ON si.inbound_uuid = i.uuid
			WHERE si.internal_squad_uuid = ?
			ORDER BY i.tag`

		rows, err := db.QueryContext(ctx, inboundQuery, squadUUID)
		if err != nil {
			return err
		}
		defer rows.Close()

		for rows.Next() {
			var inbound InboundInfo
			var network, security sql.NullString
			var port sql.NullInt64
			var rawInbound string

			if err := rows.Scan(&inbound.UUID, &inbound.Tag, &inbound.Type, &network, &security, &port, &rawInbound); err != nil {
				return err
			}

			if network.Valid {
				inbound.Network = &network.String
			}
			if security.Valid {
				inbound.Security = &security.String
			}
			if port.Valid {
				p := int(port.Int64)
				inbound.Port = &p
			}
			inbound.RawInbound = json.RawMessage(rawInbound)

			squad.Inbounds = append(squad.Inbounds, inbound)
		}

		memberQuery := `
			SELECT m.user_id, u.username
			FROM internal_squad_members m
			JOIN users u ON m.user_id = u.t_id
			WHERE m.internal_squad_uuid = ?
			ORDER BY u.username`

		memberRows, err := db.QueryContext(ctx, memberQuery, squadUUID)
		if err != nil {
			return err
		}
		defer memberRows.Close()

		for memberRows.Next() {
			var member SquadMemberInfo
			if err := memberRows.Scan(&member.UserID, &member.Username); err != nil {
				return err
			}
			squad.Members = append(squad.Members, member)
		}

		return nil
	})
	return squad, err
}

func (r *SquadRepository) getAllSquadsSummary(ctx context.Context) ([]SquadSummary, error) {
	var squads []SquadSummary
	err := r.manager.ExecuteHighPriority(func(db dbmanager.DBExecutor) error {
		query := `
			SELECT 
				s.uuid,
				s.view_position,
				s.name,
				COUNT(DISTINCT m.user_id) as members_count,
				COUNT(DISTINCT si.inbound_uuid) as inbounds_count
			FROM internal_squads s
			LEFT JOIN internal_squad_members m ON s.uuid = m.internal_squad_uuid
			LEFT JOIN internal_squad_inbounds si ON s.uuid = si.internal_squad_uuid
			GROUP BY s.uuid, s.view_position, s.name
			ORDER BY s.view_position ASC, s.name ASC`

		rows, err := db.QueryContext(ctx, query)
		if err != nil {
			return err
		}
		defer rows.Close()

		for rows.Next() {
			var squad SquadSummary
			var viewPosition sql.NullInt64
			var membersCount, inboundsCount sql.NullInt64

			if err := rows.Scan(&squad.UUID, &viewPosition, &squad.Name, &membersCount, &inboundsCount); err != nil {
				return err
			}

			if viewPosition.Valid {
				squad.ViewPosition = int(viewPosition.Int64)
			}
			if membersCount.Valid {
				squad.MembersCount = int(membersCount.Int64)
			}
			if inboundsCount.Valid {
				squad.InboundsCount = int(inboundsCount.Int64)
			}

			squads = append(squads, squad)
		}
		return rows.Err()
	})
	return squads, err
}

func (r *SquadRepository) getNodesWithConfig(ctx context.Context) ([]NodeWithConfig, error) {
	var nodes []NodeWithConfig
	err := r.manager.ExecuteHighPriority(func(db dbmanager.DBExecutor) error {
		query := `
			SELECT
				n.uuid,
				n.id,
				n.name,
				n.address,
				n.port,
				n.active_config_profile_uuid,
				COALESCE(cp.name, '') as config_profile_name,
				n.is_connected,
				n.is_disabled,
				n.view_position,
				n.country_code,
				n.tags as tags
			FROM nodes n
			LEFT JOIN config_profiles cp ON n.active_config_profile_uuid = cp.uuid
			ORDER BY n.view_position ASC, n.name ASC`

		rows, err := db.QueryContext(ctx, query)
		if err != nil {
			return err
		}
		defer rows.Close()

		for rows.Next() {
			var node NodeWithConfig
			var id sql.NullInt64
			var activeConfigProfileUUID, configProfileName sql.NullString
			var tags dbutil.StringArray

			if err := rows.Scan(&node.UUID, &id, &node.Name, &node.Address, &node.Port,
				&activeConfigProfileUUID, &configProfileName, &node.IsConnected, &node.IsDisabled,
				&node.ViewPosition, &node.CountryCode, &tags); err != nil {
				return err
			}

			if id.Valid {
				node.ID = &id.Int64
			}

			if activeConfigProfileUUID.Valid && activeConfigProfileUUID.String != "" {
				node.ActiveConfigProfileUUID = &activeConfigProfileUUID.String
				if configProfileName.Valid && configProfileName.String != "" {
					node.ConfigProfileName = &configProfileName.String
				}
			}

			if len(tags) > 0 {
				node.Tags = tags.Slice()
			} else {
				node.Tags = []string{}
			}

			nodes = append(nodes, node)
		}
		return rows.Err()
	})
	return nodes, err
}

func (r *SquadRepository) getInboundsWithProfiles(ctx context.Context) ([]InboundWithProfile, error) {
	var inbounds []InboundWithProfile
	err := r.manager.ExecuteHighPriority(func(db dbmanager.DBExecutor) error {
		query := `
			SELECT 
				p.uuid AS profile_uuid,
				p.name AS profile_name,
				i.uuid AS inbound_uuid,
				i.tag AS inbound_tag,
				i.type AS inbound_type,
				i.port AS inbound_port
			FROM config_profile_inbounds i
			JOIN config_profiles p ON p.uuid = i.profile_uuid
			ORDER BY p.name, i.tag`

		rows, err := db.QueryContext(ctx, query)
		if err != nil {
			return err
		}
		defer rows.Close()

		for rows.Next() {
			var ib InboundWithProfile
			var port sql.NullInt64

			if err := rows.Scan(&ib.ProfileUUID, &ib.ProfileName, &ib.InboundUUID, &ib.InboundTag, &ib.InboundType, &port); err != nil {
				return err
			}

			if port.Valid {
				p := int(port.Int64)
				ib.InboundPort = &p
			}

			inbounds = append(inbounds, ib)
		}
		return rows.Err()
	})
	return inbounds, err
}
