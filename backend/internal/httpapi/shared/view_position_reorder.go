package shared

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"exodus/internal/config"

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

var uuidParse = func(v string) (string, error) {
	if _, err := uuid.Parse(v); err != nil {
		return "", err
	}
	return v, nil
}

func ApplyViewPositionReorder(ctx context.Context, dbConn *sql.DB, tableName string, orderedUUIDs []string, cfg *config.BackendConfig) error {
	if !isAllowedReorderTable(tableName) {
		return fmt.Errorf("unsupported table for reorder")
	}

	tx, err := dbConn.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer func() {
		_ = tx.Rollback()
	}()

	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM %s WHERE uuid = ANY($1)", tableName)
	var foundCount int
	if err := tx.QueryRowContext(ctx, countQuery, orderedUUIDs).Scan(&foundCount); err != nil {
		return err
	}
	if foundCount != len(orderedUUIDs) {
		return fmt.Errorf("some UUIDs not found in %s", tableName)
	}

	updateQuery := fmt.Sprintf("UPDATE %s SET view_position = $1 WHERE uuid = $2", tableName)
	for i, id := range orderedUUIDs {
		if _, err := tx.ExecContext(ctx, updateQuery, i, id); err != nil {
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
