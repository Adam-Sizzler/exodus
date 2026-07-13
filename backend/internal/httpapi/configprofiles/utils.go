package configprofiles

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/google/uuid"
)

var (
	errConfigProfileSnippetEmpty                = errors.New("config profile snippet is empty")
	errConfigProfileSnippetContainsEmptyObjects = errors.New("config profile snippet contains empty objects")
	errConfigProfileSnippetInvalidJSON          = errors.New("config profile snippet is invalid json")
)

func validateCreateConfigProfileRequest(req createConfigProfileRequest) error {
	name := strings.TrimSpace(req.Name)
	if len(name) < 2 || len(name) > 30 {
		return fmt.Errorf("name must be between 2 and 30 characters")
	}
	if len(req.Config) == 0 {
		return fmt.Errorf("config is required")
	}
	var parsed map[string]any
	if err := json.Unmarshal(req.Config, &parsed); err != nil {
		return fmt.Errorf("config must be valid JSON")
	}
	return nil
}

func validateUpdateConfigProfileRequest(req updateConfigProfileRequest) error {
	if req.Name == nil && req.Config == nil {
		return fmt.Errorf("no fields to update")
	}
	if req.Name != nil {
		name := strings.TrimSpace(*req.Name)
		if len(name) < 2 || len(name) > 30 {
			return fmt.Errorf("name must be between 2 and 30 characters")
		}
	}
	if req.Config != nil {
		var parsed map[string]any
		if err := json.Unmarshal(*req.Config, &parsed); err != nil {
			return fmt.Errorf("config must be valid JSON")
		}
	}
	return nil
}

func extractInboundNetwork(inboundMap map[string]any) string {
	if transport, ok := inboundMap["transport"].(map[string]any); ok {
		if transportType, ok := transport["type"].(string); ok && transportType != "" {
			return transportType
		}
	}
	if networkValue, ok := inboundMap["network"].(string); ok && networkValue != "" {
		return networkValue
	}
	return ""
}

func extractInboundSecurity(inboundMap map[string]any) string {
	if tls, ok := inboundMap["tls"].(map[string]any); ok {
		if enabled, _ := tls["enabled"].(bool); enabled {
			return "tls"
		}
		return ""
	}
	if securityValue, ok := inboundMap["security"].(string); ok && securityValue != "" {
		return securityValue
	}
	return ""
}

func parseConfigInbounds(profileUUID string, configJSON json.RawMessage) ([]ConfigProfileInbound, error) {
	var configData map[string]any
	if err := json.Unmarshal(configJSON, &configData); err != nil {
		return nil, fmt.Errorf("failed to parse config JSON: %w", err)
	}

	inboundsRaw, ok := configData["inbounds"]
	if !ok {
		return []ConfigProfileInbound{}, nil
	}
	inboundsArray, ok := inboundsRaw.([]any)
	if !ok {
		return nil, fmt.Errorf("inbounds must be an array")
	}

	seenTags := make(map[string]struct{})
	result := make([]ConfigProfileInbound, 0, len(inboundsArray))

	for _, inboundRaw := range inboundsArray {
		inboundMap, ok := inboundRaw.(map[string]any)
		if !ok {
			continue
		}

		tag, _ := inboundMap["tag"].(string)
		if strings.TrimSpace(tag) == "" {
			continue
		}
		if _, ok := seenTags[tag]; ok {
			continue
		}
		seenTags[tag] = struct{}{}

		item := ConfigProfileInbound{
			UUID:         uuid.NewString(),
			ProfileUUID:  profileUUID,
			Tag:          tag,
			ActiveSquads: []string{},
		}
		if typeValue, ok := inboundMap["type"].(string); ok {
			item.Type = typeValue
		} else if protocolValue, ok := inboundMap["protocol"].(string); ok {
			item.Type = protocolValue
		}
		if networkValue := extractInboundNetwork(inboundMap); networkValue != "" {
			item.Network = &networkValue
		}
		if securityValue := extractInboundSecurity(inboundMap); securityValue != "" {
			item.Security = &securityValue
		}
		if portValue, ok := inboundMap["listen_port"].(float64); ok {
			p := int(portValue)
			item.Port = &p
		} else if portValue, ok := inboundMap["port"].(float64); ok {
			p := int(portValue)
			item.Port = &p
		}
		rawInbound, err := json.Marshal(inboundMap)
		if err != nil {
			continue
		}
		item.RawInbound = rawInbound
		result = append(result, item)
	}

	return result, nil
}

func validateConfigProfileSnippet(raw json.RawMessage) error {
	if len(raw) == 0 {
		return errConfigProfileSnippetEmpty
	}
	if !json.Valid(raw) {
		return errConfigProfileSnippetInvalidJSON
	}

	var items []json.RawMessage
	if err := json.Unmarshal(raw, &items); err != nil {
		return errConfigProfileSnippetEmpty
	}
	if len(items) == 0 {
		return errConfigProfileSnippetEmpty
	}

	for _, item := range items {
		var object map[string]any
		if err := json.Unmarshal(item, &object); err != nil {
			return errConfigProfileSnippetContainsEmptyObjects
		}
		if len(object) == 0 {
			return errConfigProfileSnippetContainsEmptyObjects
		}
	}

	return nil
}

func dedupeStrings(values []string) []string {
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
