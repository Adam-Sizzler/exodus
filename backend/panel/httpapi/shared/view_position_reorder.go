package shared

import (
	"context"
	"fmt"
	"strings"

	"v2ray-stat/backend/panel/config"
	dbmanager "v2ray-stat/backend/panel/db/manager"

	"github.com/google/uuid"
)

type ViewPositionReorderRequest struct {
	OrderedUUIDs []string `json:"ordered_uuids"`
}

func (r *ViewPositionReorderRequest) Validate() error {
	if len(r.OrderedUUIDs) == 0 {
		return fmt.Errorf("ordered_uuids is required")
	}

	seen := make(map[string]struct{}, len(r.OrderedUUIDs))
	for _, id := range r.OrderedUUIDs {
		id = strings.TrimSpace(id)
		if id == "" {
			return fmt.Errorf("ordered_uuids contains empty value")
		}
		if _, err := parseUUID(id); err != nil {
			return fmt.Errorf("ordered_uuids contains invalid UUID")
		}
		if _, ok := seen[id]; ok {
			return fmt.Errorf("ordered_uuids contains duplicates")
		}
		seen[id] = struct{}{}
	}

	return nil
}

func parseUUID(v string) (string, error) {
	u, err := uuidParse(v)
	if err != nil {
		return "", err
	}
	return u, nil
}

// isolated for testability and to avoid import loops in callers.
var uuidParse = func(v string) (string, error) {
	if _, err := uuid.Parse(v); err != nil {
		return "", err
	}
	return v, nil
}

func ApplyViewPositionReorder(ctx context.Context, db dbmanager.DBExecutor, tableName string, orderedUUIDs []string, cfg *config.BackendConfig) error {
	if !isAllowedReorderTable(tableName) {
		return fmt.Errorf("unsupported table for reorder")
	}

	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	inPlaceholders := strings.TrimSuffix(strings.Repeat("?,", len(orderedUUIDs)), ",")
	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM %s WHERE uuid IN (%s)", tableName, inPlaceholders)
	countArgs := make([]interface{}, 0, len(orderedUUIDs))
	for _, id := range orderedUUIDs {
		countArgs = append(countArgs, id)
	}

	var foundCount int
	if err := tx.QueryRow(countQuery, countArgs...).Scan(&foundCount); err != nil {
		_ = tx.Rollback()
		return err
	}
	if foundCount != len(orderedUUIDs) {
		_ = tx.Rollback()
		return fmt.Errorf("some UUIDs not found in %s", tableName)
	}

	updateQuery := fmt.Sprintf("UPDATE %s SET view_position = ? WHERE uuid = ?", tableName)
	for i, id := range orderedUUIDs {
		if _, err := tx.Exec(updateQuery, i, id); err != nil {
			_ = tx.Rollback()
			return err
		}
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	cfg.Logger.Info("view_position reordered", "table", tableName, "count", len(orderedUUIDs))
	return nil
}

func isAllowedReorderTable(tableName string) bool {
	switch tableName {
	case "hosts", "config_profiles", "internal_squads", "nodes", "users", "subscription_templates":
		return true
	default:
		return false
	}
}
