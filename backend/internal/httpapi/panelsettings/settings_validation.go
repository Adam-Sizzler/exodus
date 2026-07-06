package panelsettings

import (
	"fmt"
	"net/mail"
	"strings"
)

func validateAuthenticationSettings(settings map[string]any) error {
	passkey := mapValue(settings["passkey_settings"])
	password := mapValue(settings["password_settings"])
	oauth := normalizeOAuth2Settings(settings["oauth2_settings"])

	passkeyEnabled := boolValue(passkey["enabled"])
	passwordEnabled := boolValue(password["enabled"])
	oauthEnabled := map[string]bool{
		"github":   providerEnabled(oauth, "github"),
		"pocketid": providerEnabled(oauth, "pocketid"),
		"yandex":   providerEnabled(oauth, "yandex"),
		"keycloak": providerEnabled(oauth, "keycloak"),
		"generic":  providerEnabled(oauth, "generic"),
		"telegram": providerEnabled(oauth, "telegram"),
	}

	hasOAuth := false
	for _, enabled := range oauthEnabled {
		if enabled {
			hasOAuth = true
			break
		}
	}
	if !passkeyEnabled && !passwordEnabled && !hasOAuth {
		return fmt.Errorf("At least one authentication method must be enabled")
	}

	if passkeyEnabled && (stringValue(passkey["rpId"]) == "" || stringValue(passkey["origin"]) == "") {
		return fmt.Errorf("[Passkey] RP ID and origin must be set in order to use passkey authentication.")
	}

	for _, provider := range []string{"github", "yandex"} {
		cfg := mapValue(oauth[provider])
		if !boolValue(cfg["enabled"]) {
			continue
		}
		if stringValue(cfg["clientId"]) == "" || stringValue(cfg["clientSecret"]) == "" {
			return fmt.Errorf("[OAuth2] ClientID and ClientSecret must be set in order to use authentication. Please set required fields or disable misconfigured OAuth2 provider.")
		}
		allowedEmails := stringSliceValue(cfg["allowedEmails"])
		if len(allowedEmails) == 0 {
			return fmt.Errorf("[OAuth2] At least one email must be set in order to use authentication. Please set required fields or disable misconfigured OAuth2 provider.")
		}
		if err := validateEmails(allowedEmails); err != nil {
			return err
		}
	}

	for _, provider := range []string{"pocketid", "keycloak", "generic"} {
		cfg := mapValue(oauth[provider])
		if !boolValue(cfg["enabled"]) {
			continue
		}
		if stringValue(cfg["clientId"]) == "" || stringValue(cfg["clientSecret"]) == "" {
			return fmt.Errorf("[OAuth2] ClientID and ClientSecret must be set in order to use authentication. Please set required fields or disable misconfigured OAuth2 provider.")
		}
		if err := validateEmails(stringSliceValue(cfg["allowedEmails"])); err != nil {
			return err
		}
	}

	pocketID := mapValue(oauth["pocketid"])
	if boolValue(pocketID["enabled"]) && stringValue(pocketID["plainDomain"]) == "" {
		return fmt.Errorf("[PocketID] Plain domain must be set in order to use PocketID authentication.")
	}

	keycloak := mapValue(oauth["keycloak"])
	if boolValue(keycloak["enabled"]) &&
		(stringValue(keycloak["frontendDomain"]) == "" ||
			stringValue(keycloak["keycloakDomain"]) == "" ||
			stringValue(keycloak["realm"]) == "") {
		return fmt.Errorf("[Keycloak] Frontend domain, Keycloak domain and Realm must be set in order to use Keycloak authentication.")
	}

	generic := mapValue(oauth["generic"])
	if boolValue(generic["enabled"]) &&
		(stringValue(generic["authorizationUrl"]) == "" ||
			stringValue(generic["tokenUrl"]) == "" ||
			stringValue(generic["frontendDomain"]) == "") {
		return fmt.Errorf("[Generic OAuth2] Authorization URL, token URL and frontend domain must be set in order to use Generic OAuth2 authentication.")
	}

	telegram := mapValue(oauth["telegram"])
	if boolValue(telegram["enabled"]) {
		if stringValue(telegram["clientId"]) == "" ||
			stringValue(telegram["clientSecret"]) == "" ||
			stringValue(telegram["frontendDomain"]) == "" {
			return fmt.Errorf("[Telegram OAuth2] Client ID, client secret and frontend domain must be set in order to use Telegram OAuth2 authentication.")
		}
		if len(stringSliceValue(telegram["allowedIds"])) == 0 {
			return fmt.Errorf("[Telegram OAuth2] At least one admin ID must be set in order to use Telegram OAuth2 authentication.")
		}
	}

	return nil
}

func providerEnabled(oauth map[string]any, provider string) bool {
	return boolValue(mapValue(oauth[provider])["enabled"])
}

func validateEmails(emails []string) error {
	for _, email := range emails {
		email = strings.TrimSpace(email)
		if email == "" {
			continue
		}
		parsed, err := mail.ParseAddress(email)
		if err != nil || parsed.Address != email || !strings.Contains(email, "@") {
			return fmt.Errorf("[OAuth2] Email %s is not a valid email address.", email)
		}
	}
	return nil
}

func mapValue(value any) map[string]any {
	if value == nil {
		return map[string]any{}
	}
	if typed, ok := value.(map[string]any); ok {
		return typed
	}
	return map[string]any{}
}

func boolValue(value any) bool {
	typed, ok := value.(bool)
	return ok && typed
}

func stringValue(value any) string {
	typed, ok := value.(string)
	if !ok {
		return ""
	}
	return strings.TrimSpace(typed)
}

func stringSliceValue(value any) []string {
	switch typed := value.(type) {
	case []string:
		out := make([]string, 0, len(typed))
		for _, item := range typed {
			if trimmed := strings.TrimSpace(item); trimmed != "" {
				out = append(out, trimmed)
			}
		}
		return out
	case []any:
		out := make([]string, 0, len(typed))
		for _, item := range typed {
			if str, ok := item.(string); ok {
				if trimmed := strings.TrimSpace(str); trimmed != "" {
					out = append(out, trimmed)
				}
			}
		}
		return out
	default:
		return nil
	}
}
