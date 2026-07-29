package srslists

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"exodus/internal/config"
	srscore "exodus/internal/srslists"

	"github.com/google/uuid"
)

type srsListAPI struct {
	UUID           string     `json:"uuid"`
	Tag            string     `json:"tag"`
	Format         string     `json:"format"`
	URL            string     `json:"url"`
	UpdateInterval string     `json:"updateInterval"`
	Path           *string    `json:"path"`
	FileName       string     `json:"fileName"`
	ShortName      string     `json:"shortName"`
	ViewPosition   int        `json:"viewPosition"`
	IsEnabled      bool       `json:"isEnabled"`
	IsAvailable    bool       `json:"isAvailable"`
	LastCheckedAt  *time.Time `json:"lastCheckedAt"`
	LastError      *string    `json:"lastError"`
	CreatedAt      time.Time  `json:"createdAt"`
	UpdatedAt      time.Time  `json:"updatedAt"`
}

type createSRSListsRequest struct {
	URL            string   `json:"url"`
	URLs           []string `json:"urls"`
	Tag            string   `json:"tag"`
	Format         string   `json:"format"`
	UpdateInterval string   `json:"updateInterval"`
	Path           string   `json:"path"`
	IsEnabled      *bool    `json:"isEnabled"`
}

type updateSRSListRequest struct {
	UUID           string  `json:"uuid"`
	URL            *string `json:"url"`
	Tag            *string `json:"tag"`
	Format         *string `json:"format"`
	UpdateInterval *string `json:"updateInterval"`
	Path           *string `json:"path"`
	IsEnabled      *bool   `json:"isEnabled"`
}

type reorderSRSListsRequest struct {
	Items []struct {
		UUID         string `json:"uuid"`
		ViewPosition int    `json:"viewPosition"`
	} `json:"items"`
}

type bulkDeleteRequest struct {
	UUIDs []string `json:"uuids"`
}

type bulkEnableRequest struct {
	UUIDs []string `json:"uuids"`
}

type bulkSetIntervalRequest struct {
	UUIDs          []string `json:"uuids"`
	UpdateInterval string   `json:"updateInterval"`
}

type checkListsRequest struct {
	UUIDs []string `json:"uuids"`
}

func checkSelectedLists(ctx context.Context, db *sql.DB, cfg *config.BackendConfig, uuids []string) error {
	clean, err := normalizeUUIDs(uuids)
	if err != nil {
		return err
	}
	if len(clean) == 0 {
		return nil
	}

	itemsByUUID := make(map[string]string, len(clean))
	rows, err := db.QueryContext(ctx, `SELECT uuid, url FROM srs_lists WHERE uuid = ANY($1)`, clean)
	if err != nil {
		return err
	}
	defer rows.Close()
	for rows.Next() {
		var id, rawURL string
		if err := rows.Scan(&id, &rawURL); err != nil {
			return err
		}
		itemsByUUID[id] = rawURL
	}
	if err := rows.Err(); err != nil {
		return err
	}

	for id, rawURL := range itemsByUUID {
		err := srscore.CheckOneURL(ctx, rawURL)
		isAvailable := err == nil
		var errText any
		if err != nil {
			errText = err.Error()
		}
		if _, writeErr := db.ExecContext(ctx, `
			UPDATE srs_lists
			SET is_available = $1,
				last_checked_at = CURRENT_TIMESTAMP,
				last_error = $2,
				updated_at = CURRENT_TIMESTAMP
			WHERE uuid = $3
		`, isAvailable, errText, id); writeErr != nil {
			if cfg != nil && cfg.Logger != nil {
				cfg.Logger.Warn("Failed to write selected srs check result", "uuid", id, "error", writeErr)
			}
		}
	}

	return nil
}

func normalizeUUIDs(values []string) ([]string, error) {
	clean := make([]string, 0, len(values))
	seen := make(map[string]struct{}, len(values))
	for _, raw := range values {
		u := strings.TrimSpace(raw)
		if u == "" {
			continue
		}
		if _, err := uuid.Parse(u); err != nil {
			return nil, fmt.Errorf("invalid uuid: %s", u)
		}
		if _, ok := seen[u]; ok {
			continue
		}
		seen[u] = struct{}{}
		clean = append(clean, u)
	}
	return clean, nil
}
