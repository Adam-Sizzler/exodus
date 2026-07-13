package subscriptionconnections

import (
	"fmt"
	"net/url"
	"regexp"
	"strings"

	"exodus/internal/httpapi/shared"
	monitor "exodus/internal/subscriptionnodes"
)

var nodeTagRegex = regexp.MustCompile(`^[A-Z0-9_:]+$`)

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
	if req.PublicDomain != nil {
		if err := validatePublicDomain(*req.PublicDomain); err != nil {
			return err
		}
	}
	if req.APISchema != nil && strings.TrimSpace(*req.APISchema) == "" {
		return fmt.Errorf("apiSchema cannot be empty")
	}
	if req.APISchema != nil {
		switch strings.ToLower(strings.TrimSpace(*req.APISchema)) {
		case "mtls", "tls":
		default:
			return fmt.Errorf("apiSchema must be mtls or tls")
		}
	}
	if err := validateTags(req.Tags); err != nil {
		return err
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
	if req.PublicDomain.Set && req.PublicDomain.Value != nil {
		if err := validatePublicDomain(*req.PublicDomain.Value); err != nil {
			return err
		}
	}
	if req.APISchema != nil && strings.TrimSpace(*req.APISchema) == "" {
		return fmt.Errorf("apiSchema cannot be empty")
	}
	if req.APISchema != nil {
		switch strings.ToLower(strings.TrimSpace(*req.APISchema)) {
		case "mtls", "tls":
		default:
			return fmt.Errorf("apiSchema must be mtls or tls")
		}
	}
	if req.Tags != nil {
		if err := validateTags(*req.Tags); err != nil {
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
		if _, err := uuidValidate(value); err != nil {
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

func buildNodeResponses(records []nodeRecord, providersMap map[string]*providerResponse) []nodeAPI {
	response := make([]nodeAPI, 0, len(records))
	for _, record := range records {
		record = applySubNodeRuntimeSnapshot(record)

		var item nodeAPI
		item.UUID = record.UUID
		item.Name = record.Name
		item.Address = record.Address
		item.PublicDomain = record.PublicDomain
		item.Port = record.Port
		item.APISchema = normalizeSubNodeSchema(record.APISchema)
		item.APIPath = record.APIPath
		item.GRPCAuthToken = record.GRPCAuthToken
		item.SubpageConfigUUID = record.SubpageConfigUUID
		item.IsConnected = record.IsConnected
		item.IsDisabled = record.IsDisabled
		item.IsConnecting = record.IsConnecting
		item.LastStatusChange = record.LastStatusChange
		item.LastStatusMessage = record.LastStatusMessage
		item.SingboxVersion = record.SingboxVersion
		item.NodeVersion = record.NodeVersion
		item.SingboxUptime = shared.ParseUptimeSeconds(record.SingboxUptime)
		item.IsTrafficTrackingActive = record.IsTrafficTrackingActive
		item.TrafficResetDay = record.TrafficResetDay
		item.TrafficLimitBytes = record.TrafficLimitBytes
		item.TrafficUsedBytes = record.TrafficUsedBytes
		item.NotifyPercent = record.NotifyPercent
		item.UsersOnline = record.UsersOnline
		item.ViewPosition = record.ViewPosition
		item.CountryCode = record.CountryCode
		item.ConsumptionMultiplier = 1
		item.Tags = ensureStringSlice(record.Tags)
		item.CPUCount = record.CPUCount
		item.CPUModel = record.CPUModel
		item.TotalRAM = record.TotalRAM
		item.CreatedAt = record.CreatedAt
		item.UpdatedAt = record.UpdatedAt
		item.ConfigProfile.ActiveConfigProfileUUID = nil
		item.ConfigProfile.ActiveInbounds = ensureInboundSlice(nil)
		item.ProviderUUID = record.ProviderUUID
		if record.ProviderUUID != nil {
			item.Provider = providersMap[*record.ProviderUUID]
		}
		response = append(response, item)
	}

	return response
}

func applySubNodeRuntimeSnapshot(record nodeRecord) nodeRecord {
	snapshot, ok := monitor.GetSubNodeRuntimeSnapshot(record.Name)
	if !ok {
		return record
	}

	if snapshot.SingboxVersion != nil {
		value := strings.TrimSpace(*snapshot.SingboxVersion)
		if value != "" {
			record.SingboxVersion = &value
		}
	}
	if snapshot.NodeVersion != nil {
		value := strings.TrimSpace(*snapshot.NodeVersion)
		if value != "" {
			record.NodeVersion = &value
		}
	}
	if uptime := strings.TrimSpace(snapshot.SingboxUptime); uptime != "" {
		record.SingboxUptime = uptime
	}
	if snapshot.CPUCount != nil {
		value := *snapshot.CPUCount
		record.CPUCount = &value
	}
	if snapshot.CPUModel != nil {
		value := strings.TrimSpace(*snapshot.CPUModel)
		if value != "" {
			record.CPUModel = &value
		}
	}
	if snapshot.TotalRAM != nil {
		value := strings.TrimSpace(*snapshot.TotalRAM)
		if value != "" {
			record.TotalRAM = &value
		}
	}

	return record
}

func normalizeTags(tags []string) []string {
	if tags == nil {
		return []string{}
	}
	normalized := make([]string, 0, len(tags))
	for _, tag := range tags {
		tag = strings.ToUpper(strings.TrimSpace(tag))
		if tag == "" {
			continue
		}
		normalized = append(normalized, tag)
	}
	return dedupeStrings(normalized)
}

func normalizeNullableString(value *string) any {
	if value == nil || strings.TrimSpace(*value) == "" {
		return nil
	}
	return strings.TrimSpace(*value)
}

func normalizePublicDomain(value *string) any {
	if value == nil {
		return nil
	}

	trimmed := strings.TrimSpace(*value)
	if trimmed == "" {
		return nil
	}

	if !strings.Contains(trimmed, "://") {
		return strings.TrimSuffix(trimmed, "/")
	}

	parsed, err := url.Parse(trimmed)
	if err != nil || parsed.Host == "" {
		return strings.TrimSuffix(trimmed, "/")
	}

	parsed.Path = ""
	parsed.RawQuery = ""
	parsed.Fragment = ""
	parsed.User = nil
	return strings.TrimSuffix(parsed.String(), "/")
}

func normalizeNullableUUID(value *string) *string {
	if value == nil {
		return nil
	}
	trimmed := strings.TrimSpace(*value)
	if trimmed == "" {
		return nil
	}
	return &trimmed
}

func normalizeAPISchema(value *string) string {
	if value == nil {
		return "mtls"
	}
	normalized := strings.ToLower(strings.TrimSpace(*value))
	if normalized == "" {
		return "mtls"
	}
	switch normalized {
	case "tls":
		return "tls"
	case "mtls":
		return "mtls"
	default:
		return "mtls"
	}
}

func normalizeSubNodeSchema(value string) string {
	normalized := strings.ToLower(strings.TrimSpace(value))
	switch normalized {
	case "tls":
		return "tls"
	case "mtls":
		return "mtls"
	default:
		return "mtls"
	}
}

func normalizeAPIPath(value *string) string {
	if value == nil {
		return "/"
	}
	path := strings.TrimSpace(*value)
	if path == "" {
		return "/"
	}
	if !strings.HasPrefix(path, "/") {
		return "/" + path
	}
	return path
}

func validatePublicDomain(value string) error {
	trimmed := strings.TrimSpace(value)
	if trimmed == "" {
		return nil
	}
	if len(trimmed) > 255 {
		return fmt.Errorf("publicDomain must be less than or equal to 255 characters")
	}

	candidate := trimmed
	if !strings.Contains(candidate, "://") {
		candidate = "https://" + candidate
	}

	parsed, err := url.Parse(candidate)
	if err != nil || strings.TrimSpace(parsed.Host) == "" {
		return fmt.Errorf("publicDomain must be a valid domain or URL")
	}
	if parsed.User != nil || parsed.RawQuery != "" || parsed.Fragment != "" {
		return fmt.Errorf("publicDomain must not contain credentials, query or fragment")
	}
	if parsed.Path != "" && parsed.Path != "/" {
		return fmt.Errorf("publicDomain must not contain path")
	}

	return nil
}

func ensureStringSlice(values []string) []string {
	if values == nil {
		return []string{}
	}
	return values
}

func ensureInboundSlice(values []configProfileInboundResponse) []configProfileInboundResponse {
	if values == nil {
		return []configProfileInboundResponse{}
	}
	return values
}

func dedupeStrings(values []string) []string {
	if len(values) < 2 {
		return values
	}
	seen := make(map[string]struct{}, len(values))
	result := make([]string, 0, len(values))
	for _, value := range values {
		if _, ok := seen[value]; ok {
			continue
		}
		seen[value] = struct{}{}
		result = append(result, value)
	}
	return result
}

// Simple helper to avoid external uuid dependency in utils
func uuidValidate(uuidStr string) (string, error) {
	parsed := strings.TrimSpace(uuidStr)
	if len(parsed) != 36 {
		return "", fmt.Errorf("invalid uuid length")
	}
	return parsed, nil
}
