package passkeys

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	dbmanager "exodus/internal/db/manager"

	gowebauthn "github.com/go-webauthn/webauthn/webauthn"
)

func loadPasskeySettings(ctx context.Context, manager *dbmanager.DatabaseManager) (passkeySettings, error) {
	var settings passkeySettings
	err := manager.ExecuteHighPriority(func(db dbmanager.DBExecutor) error {
		var rawConfig sql.NullString
		row := db.QueryRowContext(ctx, `
			SELECT passkey_settings::text
			FROM exodus_settings
			WHERE id = 1
			LIMIT 1
		`)
		if scanErr := row.Scan(&rawConfig); scanErr != nil {
			return scanErr
		}
		if !rawConfig.Valid || strings.TrimSpace(rawConfig.String) == "" {
			return nil
		}
		return json.Unmarshal([]byte(rawConfig.String), &settings)
	})
	return settings, err
}

func loadWebAuthnAdmin(ctx context.Context, manager *dbmanager.DatabaseManager, adminUUID string) (*webAuthnAdmin, error) {
	admin := &webAuthnAdmin{}
	err := manager.ExecuteHighPriority(func(db dbmanager.DBExecutor) error {
		var row *sql.Row
		if strings.TrimSpace(adminUUID) == "" {
			row = db.QueryRowContext(ctx, `
				SELECT uuid, username
				FROM admin
				ORDER BY created_at ASC
				LIMIT 1
			`)
		} else {
			row = db.QueryRowContext(ctx, `
				SELECT uuid, username
				FROM admin
				WHERE uuid = ?
				LIMIT 1
			`, adminUUID)
		}

		if err := row.Scan(&admin.uuid, &admin.username); err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return errAdminNotFound
			}
			return err
		}

		rows, err := db.QueryContext(ctx, `
			SELECT id, public_key, counter, COALESCE(transports, ''), COALESCE(device_type, ''), backed_up
			FROM passkeys
			WHERE admin_uuid = ?
			ORDER BY created_at ASC
		`, admin.uuid)
		if err != nil {
			return err
		}
		defer rows.Close()

		for rows.Next() {
			var (
				id        string
				publicKey []byte
				counter   int64
				transRaw  string
				device    string
				backedUp  bool
			)
			if err := rows.Scan(&id, &publicKey, &counter, &transRaw, &device, &backedUp); err != nil {
				return err
			}

			rawID, err := decodeCredentialID(id)
			if err != nil {
				return fmt.Errorf("decode credential id: %w", err)
			}
			backupEligible := strings.EqualFold(device, "multiDevice") || backedUp
			admin.credentials = append(admin.credentials, gowebauthn.Credential{
				ID:        rawID,
				PublicKey: publicKey,
				Transport: parseTransports(transRaw),
				Flags: gowebauthn.CredentialFlags{
					BackupEligible: backupEligible,
					BackupState:    backedUp,
				},
				Authenticator: gowebauthn.Authenticator{
					SignCount: uint32Counter(counter),
				},
			})
		}
		return rows.Err()
	})
	if err != nil {
		return nil, err
	}
	if strings.TrimSpace(admin.uuid) == "" {
		return nil, errAdminNotFound
	}
	return admin, nil
}

func saveNewCredential(ctx context.Context, manager *dbmanager.DatabaseManager, adminUUID string, credential *gowebauthn.Credential) error {
	credentialID := encodeCredentialID(credential.ID)
	transports := transportsToCSV(credential.Transport)
	deviceType := "singleDevice"
	if credential.Flags.BackupEligible {
		deviceType = "multiDevice"
	}
	provider := passkeyProviderName(credential.Transport)

	return manager.ExecuteHighPriority(func(db dbmanager.DBExecutor) error {
		_, err := db.ExecContext(ctx, `
			INSERT INTO passkeys (
				id, admin_uuid, public_key, counter, device_type, backed_up, transports, passkey_provider
			) VALUES (?, ?, ?, ?, ?, ?, ?, ?)
		`, credentialID, adminUUID, credential.PublicKey, int64(credential.Authenticator.SignCount), deviceType, credential.Flags.BackupState, transports, provider)
		return err
	})
}

func updateCredentialUsage(ctx context.Context, manager *dbmanager.DatabaseManager, adminUUID string, credential *gowebauthn.Credential) error {
	credentialID := encodeCredentialID(credential.ID)
	deviceType := "singleDevice"
	if credential.Flags.BackupEligible {
		deviceType = "multiDevice"
	}

	var rowsAffected int64
	err := manager.ExecuteHighPriority(func(db dbmanager.DBExecutor) error {
		result, err := db.ExecContext(ctx, `
			UPDATE passkeys
			SET counter = ?, device_type = ?, backed_up = ?, updated_at = CURRENT_TIMESTAMP
			WHERE id = ? AND admin_uuid = ?
		`, int64(credential.Authenticator.SignCount), deviceType, credential.Flags.BackupState, credentialID, adminUUID)
		if err != nil {
			return err
		}
		rowsAffected, err = result.RowsAffected()
		return err
	})
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return errAdminNotFound
	}
	return nil
}

func listPasskeysForAdmin(ctx context.Context, manager *dbmanager.DatabaseManager, adminUUID string) ([]passkeyRecord, error) {
	records := make([]passkeyRecord, 0)
	err := manager.ExecuteHighPriority(func(db dbmanager.DBExecutor) error {
		rows, err := db.QueryContext(ctx, `
			SELECT id,
			       COALESCE(NULLIF(passkey_provider, ''), id) AS name,
			       created_at,
			       updated_at
			FROM passkeys
			WHERE admin_uuid = ?
			ORDER BY created_at DESC
		`, adminUUID)
		if err != nil {
			return err
		}
		defer rows.Close()

		for rows.Next() {
			var item passkeyRecord
			if scanErr := rows.Scan(&item.ID, &item.Name, &item.CreatedAt, &item.LastUsedAt); scanErr != nil {
				return scanErr
			}
			records = append(records, item)
		}
		return rows.Err()
	})
	return records, err
}
