package configprofiles

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
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
	return validateConfigStructure(req.Config)
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
		return validateConfigStructure(*req.Config)
	}
	return nil
}

func validateConfigStructure(configRaw json.RawMessage) error {
	var parsed map[string]any
	if err := json.Unmarshal(configRaw, &parsed); err != nil {
		return fmt.Errorf("config must be valid JSON")
	}
	outboundsRaw, ok := parsed["outbounds"]
	if !ok {
		return fmt.Errorf("Config doesn't have outbounds.")
	}
	outboundsArray, ok := outboundsRaw.([]any)
	if !ok || len(outboundsArray) == 0 {
		return fmt.Errorf("Config doesn't have outbounds.")
	}
	inboundsRaw, ok := parsed["inbounds"]
	if !ok {
		return nil
	}
	inboundsArray, ok := inboundsRaw.([]any)
	if !ok {
		return fmt.Errorf("inbounds must be an array")
	}
	seenTags := make(map[string]struct{})
	for _, inboundRaw := range inboundsArray {
		inboundMap, ok := inboundRaw.(map[string]any)
		if !ok {
			continue
		}
		tag, _ := inboundMap["tag"].(string)
		tag = strings.TrimSpace(tag)
		if tag == "" {
			return fmt.Errorf("all inbounds must have a non-empty tag")
		}
		if strings.Contains(tag, ",") {
			return fmt.Errorf("character ',' is not allowed in inbound tag %q", tag)
		}
		if _, ok := seenTags[tag]; ok {
			return fmt.Errorf("duplicate inbound tag %q found. All inbound tags must be unique", tag)
		}
		seenTags[tag] = struct{}{}

		if err := validateShadowsocksInbound(inboundMap); err != nil {
			return err
		}
	}
	return nil
}

func validateShadowsocksInbound(inboundMap map[string]any) error {
	inboundType, _ := inboundMap["type"].(string)
	if inboundType == "" {
		inboundType, _ = inboundMap["protocol"].(string)
	}
	inboundType = strings.ToLower(strings.TrimSpace(inboundType))
	if inboundType != "shadowsocks" && inboundType != "ss" {
		return nil
	}

	method, _ := inboundMap["method"].(string)
	password, _ := inboundMap["password"].(string)

	if settings, ok := inboundMap["settings"].(map[string]any); ok {
		if m, ok := settings["method"].(string); ok && strings.TrimSpace(m) != "" {
			method = m
		}
		if p, ok := settings["password"].(string); ok && strings.TrimSpace(p) != "" {
			password = p
		}
	}

	method = strings.TrimSpace(method)
	password = strings.TrimSpace(password)

	if strings.HasPrefix(method, "2022-blake3-") {
		if password == "" {
			return fmt.Errorf("Shadowsocks password is required for 2022-blake3-* methods. (inbound → settings → password – generate with: openssl rand -base64 32)")
		}
		decoded, err := base64.StdEncoding.DecodeString(password)
		if err != nil || len(decoded) != 32 {
			return fmt.Errorf("Shadowsocks password for %q must be a base64 string that decodes to exactly 32 bytes. (inbound → settings → password – generate with: openssl rand -base64 32)", method)
		}
	}

	return nil
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
