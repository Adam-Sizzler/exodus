package nodes

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"exodus/internal/config"
	dbmanager "exodus/internal/db/manager"
	"exodus/internal/httpapi/shared"

	"github.com/google/uuid"
)

func validateCreateRequest(req createNodeRequest) error {
	if len(strings.TrimSpace(req.Name)) < 3 || len(strings.TrimSpace(req.Name)) > 30 {
		return fmt.Errorf("name must be between 3 and 30 characters")
	}
	if len(strings.TrimSpace(req.Address)) < 2 {
		return fmt.Errorf("address must be at least 2 characters")
	}
	if req.Port != nil && (*req.Port < 1 || *req.Port > 65535) {
		return fmt.Errorf("invalid port")
	}
	if req.APISchema != nil && strings.TrimSpace(*req.APISchema) == "" {
		return fmt.Errorf("apiSchema cannot be empty")
	}
	if req.APISchema != nil && !isAllowedNodeAPISchema(*req.APISchema) {
		return fmt.Errorf("apiSchema must be one of: mtls, tls")
	}
	if req.TrafficLimitBytes != nil && *req.TrafficLimitBytes < 0 {
		return fmt.Errorf("trafficLimitBytes must be greater than or equal to 0")
	}
	if req.NotifyPercent != nil && (*req.NotifyPercent < 0 || *req.NotifyPercent > 100) {
		return fmt.Errorf("notifyPercent must be between 0 and 100")
	}
	if req.TrafficResetDay != nil && (*req.TrafficResetDay < 1 || *req.TrafficResetDay > 31) {
		return fmt.Errorf("trafficResetDay must be between 1 and 31")
	}
	if req.ConsumptionMultiplier != nil && (*req.ConsumptionMultiplier < 0 || *req.ConsumptionMultiplier > 100) {
		return fmt.Errorf("consumptionMultiplier must be between 0 and 100")
	}
	if req.ActivePluginUUID != nil && strings.TrimSpace(*req.ActivePluginUUID) != "" {
		if _, err := uuid.Parse(strings.TrimSpace(*req.ActivePluginUUID)); err != nil {
			return fmt.Errorf("invalid activePluginUuid")
		}
	}
	if _, err := uuid.Parse(req.ConfigProfile.ActiveConfigProfileUUID); err != nil {
		return fmt.Errorf("invalid activeConfigProfileUuid")
	}
	if err := validateUUIDs(req.ConfigProfile.ActiveInbounds); err != nil {
		return err
	}
	if err := validateTags(req.Tags); err != nil {
		return err
	}
	if req.CountryCode != nil {
		code := strings.TrimSpace(*req.CountryCode)
		if code != "" && len(code) != 2 {
			return fmt.Errorf("countryCode must be 2 characters")
		}
	}
	return nil
}

func validateUpdateRequest(req updateNodeRequest) error {
	if req.Name != nil {
		name := strings.TrimSpace(*req.Name)
		if len(name) < 3 || len(name) > 30 {
			return fmt.Errorf("name must be between 3 and 30 characters")
		}
	}
	if req.Address != nil && len(strings.TrimSpace(*req.Address)) < 2 {
		return fmt.Errorf("address must be at least 2 characters")
	}
	if req.Port != nil && (*req.Port < 1 || *req.Port > 65535) {
		return fmt.Errorf("invalid port")
	}
	if req.APISchema != nil && strings.TrimSpace(*req.APISchema) == "" {
		return fmt.Errorf("apiSchema cannot be empty")
	}
	if req.APISchema != nil && !isAllowedNodeAPISchema(*req.APISchema) {
		return fmt.Errorf("apiSchema must be one of: mtls, tls")
	}
	if req.TrafficLimitBytes != nil && *req.TrafficLimitBytes < 0 {
		return fmt.Errorf("trafficLimitBytes must be greater than or equal to 0")
	}
	if req.NotifyPercent != nil && (*req.NotifyPercent < 0 || *req.NotifyPercent > 100) {
		return fmt.Errorf("notifyPercent must be between 0 and 100")
	}
	if req.TrafficResetDay != nil && (*req.TrafficResetDay < 1 || *req.TrafficResetDay > 31) {
		return fmt.Errorf("trafficResetDay must be between 1 and 31")
	}
	if req.ConsumptionMultiplier != nil && (*req.ConsumptionMultiplier < 0 || *req.ConsumptionMultiplier > 100) {
		return fmt.Errorf("consumptionMultiplier must be between 0 and 100")
	}
	if req.ActivePluginUUID.Set && req.ActivePluginUUID.Value != nil && strings.TrimSpace(*req.ActivePluginUUID.Value) != "" {
		if _, err := uuid.Parse(strings.TrimSpace(*req.ActivePluginUUID.Value)); err != nil {
			return fmt.Errorf("invalid activePluginUuid")
		}
	}
	if req.ConfigProfile != nil {
		if _, err := uuid.Parse(req.ConfigProfile.ActiveConfigProfileUUID); err != nil {
			return fmt.Errorf("invalid activeConfigProfileUuid")
		}
		if err := validateUUIDs(req.ConfigProfile.ActiveInbounds); err != nil {
			return err
		}
	}
	if req.Tags != nil {
		if err := validateTags(*req.Tags); err != nil {
			return err
		}
	}
	if req.CountryCode != nil {
		code := strings.TrimSpace(*req.CountryCode)
		if code != "" && len(code) != 2 {
			return fmt.Errorf("countryCode must be 2 characters")
		}
	}
	return nil
}

func validateBulkUpdateRequest(req bulkUpdateNodesRequest) error {
	if err := validateUUIDs(req.UUIDs); err != nil {
		return err
	}
	fields := req.Fields
	if fields.ConsumptionMultiplier != nil && (*fields.ConsumptionMultiplier < 0 || *fields.ConsumptionMultiplier > 100) {
		return fmt.Errorf("consumptionMultiplier must be between 0 and 100")
	}
	if fields.CountryCode != nil {
		code := strings.TrimSpace(*fields.CountryCode)
		if code != "" && len(code) != 2 {
			return fmt.Errorf("countryCode must be 2 characters")
		}
	}
	if fields.ProviderUUID.Set && fields.ProviderUUID.Value != nil && strings.TrimSpace(*fields.ProviderUUID.Value) != "" {
		if _, err := uuid.Parse(strings.TrimSpace(*fields.ProviderUUID.Value)); err != nil {
			return fmt.Errorf("invalid providerUuid")
		}
	}
	if fields.ActivePluginUUID.Set && fields.ActivePluginUUID.Value != nil && strings.TrimSpace(*fields.ActivePluginUUID.Value) != "" {
		if _, err := uuid.Parse(strings.TrimSpace(*fields.ActivePluginUUID.Value)); err != nil {
			return fmt.Errorf("invalid activePluginUuid")
		}
	}
	if fields.Tags != nil {
		if err := validateTags(*fields.Tags); err != nil {
			return err
		}
	}
	return nil
}

func validateUUIDs(values []string) error {
	if len(values) == 0 {
		return fmt.Errorf("uuids cannot be empty")
	}
	for _, value := range values {
		if _, err := uuid.Parse(value); err != nil {
			return fmt.Errorf("invalid uuid value")
		}
	}
	return nil
}

func validateTags(tags []string) error {
	if len(tags) > 10 {
		return fmt.Errorf("maximum 10 tags")
	}
	for _, tag := range tags {
		if !nodeTagRegex.MatchString(tag) {
			return fmt.Errorf("tag can only contain uppercase letters, numbers, underscores and colons")
		}
		if len(tag) > 36 {
			return fmt.Errorf("each tag must be less than 36 characters")
		}
	}
	return nil
}

func handleConfigProfileValidationError(w http.ResponseWriter, err error, cfg *config.BackendConfig) {
	switch {
	case errors.Is(err, errConfigProfileNotFound):
		shared.SendError(w, http.StatusBadRequest, "config profile not found", nil, cfg)
	case errors.Is(err, errConfigProfileInboundInvalid):
		shared.SendError(w, http.StatusBadRequest, "config profile inbound not found in specified profile", nil, cfg)
	default:
		shared.SendError(w, http.StatusInternalServerError, "failed to validate config profile inbounds", err, cfg)
	}
}

func ensureConfigProfileInbounds(ctx context.Context, manager *dbmanager.DatabaseManager, profileUUID string, inboundUUIDs []string) error {
	return manager.ExecuteHighPriority(func(db dbmanager.DBExecutor) error {
		var exists int
		if err := db.QueryRowContext(ctx, `SELECT 1 FROM config_profiles WHERE uuid = ?`, profileUUID).Scan(&exists); err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return errConfigProfileNotFound
			}
			return err
		}

		found := make(map[string]struct{}, len(inboundUUIDs))
		rows, err := db.QueryContext(ctx, `
			SELECT uuid
			FROM config_profile_inbounds
			WHERE profile_uuid = ? AND uuid = ANY(?)
		`, profileUUID, inboundUUIDs)
		if err != nil {
			return err
		}
		defer rows.Close()
		for rows.Next() {
			var inboundUUID string
			if err := rows.Scan(&inboundUUID); err != nil {
				return err
			}
			found[inboundUUID] = struct{}{}
		}
		if err := rows.Err(); err != nil {
			return err
		}
		for _, inboundUUID := range inboundUUIDs {
			if _, ok := found[inboundUUID]; !ok {
				return errConfigProfileInboundInvalid
			}
		}
		return nil
	})
}
