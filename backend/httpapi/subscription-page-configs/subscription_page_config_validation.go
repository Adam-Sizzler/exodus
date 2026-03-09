package subscriptionpageconfigs

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strings"
)

var localeCodeRegex = regexp.MustCompile(`^[a-z]{2}$`)

func normalizeAndValidateSubpageConfig(raw json.RawMessage) (map[string]interface{}, error) {
	var data map[string]interface{}
	if err := json.Unmarshal(raw, &data); err != nil {
		return nil, err
	}

	locales := normalizeLocales(data["locales"])
	if len(locales) == 0 {
		return nil, fmt.Errorf("locales is required")
	}
	data["locales"] = locales

	platforms, ok := data["platforms"].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("platforms must be an object")
	}

	baseTranslations, ok := data["baseTranslations"].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("baseTranslations must be an object")
	}

	svgLibrary, ok := data["svgLibrary"].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("svgLibrary must be an object")
	}

	if err := validateLocalizedTexts(platforms, locales, "platforms"); err != nil {
		return nil, err
	}
	if err := validateLocalizedTexts(baseTranslations, locales, "baseTranslations"); err != nil {
		return nil, err
	}
	if err := validateSvgReferences(svgLibrary, platforms); err != nil {
		return nil, err
	}

	data["platforms"] = cleanLocalizedTexts(platforms, locales)
	data["baseTranslations"] = cleanLocalizedTexts(baseTranslations, locales)

	return data, nil
}

func normalizeLocales(value interface{}) []string {
	allowed := map[string]struct{}{"en": {}, "ru": {}}
	out := make([]string, 0, 2)
	seen := map[string]struct{}{}

	slice, ok := value.([]interface{})
	if !ok {
		return out
	}
	for _, item := range slice {
		str, ok := item.(string)
		if !ok {
			continue
		}
		str = strings.ToLower(strings.TrimSpace(str))
		if !localeCodeRegex.MatchString(str) {
			continue
		}
		if _, allowedOK := allowed[str]; !allowedOK {
			continue
		}
		if _, exists := seen[str]; exists {
			continue
		}
		seen[str] = struct{}{}
		out = append(out, str)
	}
	return out
}

func isLocalizedText(obj map[string]interface{}) bool {
	if len(obj) == 0 {
		return false
	}
	hasLocale := false
	for key, val := range obj {
		if _, ok := val.(string); !ok {
			return false
		}
		if len(key) == 2 {
			hasLocale = true
		}
	}
	return hasLocale
}

func validateLocalizedTexts(obj interface{}, requiredLocales []string, path string) error {
	switch v := obj.(type) {
	case map[string]interface{}:
		if isLocalizedText(v) {
			for _, locale := range requiredLocales {
				val, ok := v[locale].(string)
				if !ok || strings.TrimSpace(val) == "" {
					return fmt.Errorf("missing required locale '%s' at %s", locale, path)
				}
			}
			return nil
		}
		for key, value := range v {
			if err := validateLocalizedTexts(value, requiredLocales, path+"."+key); err != nil {
				return err
			}
		}
	case []interface{}:
		for i, item := range v {
			if err := validateLocalizedTexts(item, requiredLocales, fmt.Sprintf("%s[%d]", path, i)); err != nil {
				return err
			}
		}
	}
	return nil
}

func validateSvgReferences(svgLibrary map[string]interface{}, platforms map[string]interface{}) error {
	validKeys := map[string]struct{}{}
	for key, value := range svgLibrary {
		if _, ok := value.(string); ok {
			validKeys[key] = struct{}{}
		}
	}

	return checkSvgReferences(platforms, validKeys, "platforms")
}

func checkSvgReferences(obj interface{}, validKeys map[string]struct{}, path string) error {
	switch v := obj.(type) {
	case map[string]interface{}:
		for key, value := range v {
			if key == "svgIconKey" {
				if str, ok := value.(string); ok {
					if _, exists := validKeys[str]; !exists {
						return fmt.Errorf("unknown svgIconKey '%s' at %s.%s", str, path, key)
					}
				}
				continue
			}

			if err := checkSvgReferences(value, validKeys, path+"."+key); err != nil {
				return err
			}
		}
	case []interface{}:
		for i, item := range v {
			if err := checkSvgReferences(item, validKeys, fmt.Sprintf("%s[%d]", path, i)); err != nil {
				return err
			}
		}
	}
	return nil
}

func cleanLocalizedTexts(obj interface{}, locales []string) interface{} {
	switch v := obj.(type) {
	case map[string]interface{}:
		if isLocalizedText(v) {
			allowed := map[string]struct{}{}
			for _, locale := range locales {
				allowed[locale] = struct{}{}
			}
			cleaned := map[string]interface{}{}
			for key, value := range v {
				str, ok := value.(string)
				if !ok {
					continue
				}
				if _, ok := allowed[key]; !ok {
					continue
				}
				str = strings.TrimSpace(str)
				if str != "" {
					cleaned[key] = str
				}
			}
			return cleaned
		}

		result := map[string]interface{}{}
		for key, value := range v {
			result[key] = cleanLocalizedTexts(value, locales)
		}
		return result
	case []interface{}:
		cleaned := make([]interface{}, 0, len(v))
		for _, item := range v {
			cleaned = append(cleaned, cleanLocalizedTexts(item, locales))
		}
		return cleaned
	default:
		return obj
	}
}
