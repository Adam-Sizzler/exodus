package seed

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"exodus/internal/config"
)

func ensureDefaultInternalSquad(ctx context.Context, tx *sql.Tx, _ *config.BackendConfig) error {
	var count int
	if err := tx.QueryRowContext(ctx, `SELECT COUNT(*) FROM internal_squads`).Scan(&count); err != nil {
		return fmt.Errorf("count internal_squads: %w", err)
	}
	if count > 0 {
		fmt.Println("ℹ Default internal squad already exists")
		return nil
	}

	var profileUUID string
	if err := tx.QueryRowContext(ctx, `SELECT uuid FROM config_profiles ORDER BY view_position ASC LIMIT 1`).Scan(&profileUUID); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			fmt.Println("ℹ Default internal squad skipped (no config profiles found)")
			return nil
		}
		return fmt.Errorf("read config profile for internal squad: %w", err)
	}

	rows, err := tx.QueryContext(ctx, `SELECT uuid FROM config_profile_inbounds WHERE profile_uuid = $1`, profileUUID)
	if err != nil {
		return fmt.Errorf("read config profile inbounds: %w", err)
	}
	defer rows.Close()

	var inboundUUIDs []string
	for rows.Next() {
		var inbUUID string
		if err := rows.Scan(&inbUUID); err != nil {
			return err
		}
		inboundUUIDs = append(inboundUUIDs, inbUUID)
	}
	if err := rows.Err(); err != nil {
		return err
	}

	squadUUID := "00000000-0000-0000-0000-000000000000"
	if _, err := tx.ExecContext(ctx, `
		INSERT INTO internal_squads (uuid, view_position, name, created_at, updated_at)
		VALUES ($1, $2, $3, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
	`, squadUUID, 1, "Default-Squad"); err != nil {
		return fmt.Errorf("insert default internal squad: %w", err)
	}

	for _, inbUUID := range inboundUUIDs {
		if _, err := tx.ExecContext(ctx, `
			INSERT INTO internal_squad_inbounds (internal_squad_uuid, inbound_uuid)
			VALUES ($1, $2) ON CONFLICT DO NOTHING
		`, squadUUID, inbUUID); err != nil {
			return fmt.Errorf("link inbound %s to default squad: %w", inbUUID, err)
		}
	}

	fmt.Printf("✔ Default internal squad seeded (inbounds: %d)\n", len(inboundUUIDs))
	return nil
}
