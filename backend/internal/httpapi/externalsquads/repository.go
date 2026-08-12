package externalsquads

import (
	"context"
	"database/sql"
	"encoding/json"
	"time"

	"exodus/internal/httpapi/shared"
)

type ExternalSquadRecord struct {
	UUID                  string          `json:"uuid"`
	ViewPosition          int             `json:"view_position"`
	Name                  string          `json:"name"`
	SubscriptionSettings  json.RawMessage `json:"subscription_settings,omitempty"`
	HostOverrides         json.RawMessage `json:"host_overrides,omitempty"`
	ResponseHeadersAdd    json.RawMessage `json:"response_headers_add,omitempty"`
	ResponseHeadersRemove json.RawMessage `json:"response_headers_remove,omitempty"`
	HWIDSettings          json.RawMessage `json:"hwid_settings,omitempty"`
	CustomRemarks         json.RawMessage `json:"custom_remarks,omitempty"`
	SubpageConfigUUID     *string         `json:"subpage_config_uuid,omitempty"`
	CreatedAt             time.Time       `json:"created_at"`
	UpdatedAt             time.Time       `json:"updated_at"`
}

func getExternalSquads(ctx context.Context, db *sql.DB) ([]ExternalSquadRecord, error) {
	rows, err := db.QueryContext(ctx, `
		SELECT uuid, view_position, name,
			subscription_settings, host_overrides, response_headers_add,
			array_to_json(COALESCE(response_headers_remove, ARRAY[]::text[]))::text AS response_headers_remove,
			hwid_settings, custom_remarks, subpage_config_uuid,
			created_at, updated_at
		FROM external_squads
		ORDER BY view_position ASC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	records := make([]ExternalSquadRecord, 0)
	for rows.Next() {
		rec, err := scanExternalSquad(rows)
		if err != nil {
			return nil, err
		}
		records = append(records, rec)
	}

	return records, rows.Err()
}

func getExternalSquadByUUID(ctx context.Context, db *sql.DB, squadUUID string) (ExternalSquadRecord, error) {
	row := db.QueryRowContext(ctx, `
		SELECT uuid, view_position, name,
			subscription_settings, host_overrides, response_headers_add,
			array_to_json(COALESCE(response_headers_remove, ARRAY[]::text[]))::text AS response_headers_remove,
			hwid_settings, custom_remarks, subpage_config_uuid,
			created_at, updated_at
		FROM external_squads
		WHERE uuid = $1
	`, squadUUID)

	return scanExternalSquad(row)
}

func scanExternalSquad(scanner shared.RowScanner) (ExternalSquadRecord, error) {
	var rec ExternalSquadRecord
	var subSettings, hostOverrides, respHeadersAdd, respHeadersRemove, hwidSettings, customRemarks sql.NullString
	var subpageConfigUUID sql.NullString

	err := scanner.Scan(
		&rec.UUID,
		&rec.ViewPosition,
		&rec.Name,
		&subSettings,
		&hostOverrides,
		&respHeadersAdd,
		&respHeadersRemove,
		&hwidSettings,
		&customRemarks,
		&subpageConfigUUID,
		&rec.CreatedAt,
		&rec.UpdatedAt,
	)
	if err != nil {
		return rec, err
	}

	rec.SubscriptionSettings = parseJSONRaw(subSettings)
	rec.HostOverrides = parseJSONRaw(hostOverrides)
	rec.ResponseHeadersAdd = parseJSONRaw(respHeadersAdd)
	rec.ResponseHeadersRemove = parseJSONRaw(respHeadersRemove)
	rec.HWIDSettings = parseJSONRaw(hwidSettings)
	rec.CustomRemarks = parseJSONRaw(customRemarks)
	if subpageConfigUUID.Valid {
		rec.SubpageConfigUUID = &subpageConfigUUID.String
	}

	return rec, nil
}
