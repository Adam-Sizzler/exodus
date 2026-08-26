package nodes

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"exodus/internal/config"
	"exodus/internal/httpapi/shared"
	"exodus/internal/util"

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
	if err := validateProxyURL(req.ProxyURL); err != nil {
		return err
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
	if req.NodeConsumptionMultiplier != nil && (*req.NodeConsumptionMultiplier < 0 || *req.NodeConsumptionMultiplier > 100) {
		return fmt.Errorf("nodeConsumptionMultiplier must be between 0 and 100")
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
	if err := validateNodeIPs(req.IPs); err != nil {
		return err
	}
	if err := validateNote(req.Note); err != nil {
		return err
	}
	return nil
}

func validateNodeIPs(items []NodeIPItem) error {
	if len(items) > 64 {
		return fmt.Errorf("ips cannot contain more than 64 items")
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
	if req.ProxyURL.Set {
		if err := validateProxyURL(req.ProxyURL.Value); err != nil {
			return err
		}
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
	if req.NodeConsumptionMultiplier != nil && (*req.NodeConsumptionMultiplier < 0 || *req.NodeConsumptionMultiplier > 100) {
		return fmt.Errorf("nodeConsumptionMultiplier must be between 0 and 100")
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
	if req.IPs != nil {
		if err := validateNodeIPs(*req.IPs); err != nil {
			return err
		}
	}
	if req.Note.Set {
		if err := validateNote(req.Note.Value); err != nil {
			return err
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
	if fields.NodeConsumptionMultiplier != nil && (*fields.NodeConsumptionMultiplier < 0 || *fields.NodeConsumptionMultiplier > 100) {
		return fmt.Errorf("nodeConsumptionMultiplier must be between 0 and 100")
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
	if fields.Note.Set {
		if err := validateNote(fields.Note.Value); err != nil {
			return err
		}
	}
	return nil
}

func validateProxyURL(value *string) error {
	if value == nil || strings.TrimSpace(*value) == "" {
		return nil
	}
	if !nodeProxyURLRegex.MatchString(strings.TrimSpace(*value)) {
		return fmt.Errorf("proxyUrl must match socks5://[user:pass@]host:port")
	}
	return nil
}

func validateNote(value *string) error {
	if value == nil {
		return nil
	}
	if len(strings.TrimSpace(*value)) > 255 {
		return fmt.Errorf("note must be less than 255 characters")
	}
	return nil
}

func validateUUIDs(values []string) error {
	return util.ValidateUUIDs(values)
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
		shared.SendAPIError(w, shared.ErrConfigProfileNotFound, cfg)
	case errors.Is(err, errConfigProfileInboundInvalid):
		shared.SendAPIError(w, shared.ErrConfigProfileInboundNotFoundInProfile, cfg)
	default:
		shared.SendAPIError(w, shared.ErrValidateConfigProfileInboundsFailed.WithCause(err), cfg)
	}
}
