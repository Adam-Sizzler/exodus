package panelsettings

// DefaultBrandingSettings returns the default canonical map for branding settings (null title and null logoUrl).
func DefaultBrandingSettings() map[string]any {
	return map[string]any{
		"title":   nil,
		"logoUrl": nil,
	}
}

// DefaultPasswordSettings returns the default password auth settings.
func DefaultPasswordSettings() map[string]any {
	return map[string]any{
		"enabled": true,
	}
}

// DefaultPasskeySettings returns the default passkey settings.
func DefaultPasskeySettings() map[string]any {
	return map[string]any{
		"enabled": false,
		"rpId":    nil,
		"origin":  nil,
	}
}

// DefaultOAuth2Settings returns the default oauth2 providers settings.
func DefaultOAuth2Settings() map[string]any {
	return map[string]any{
		"github": map[string]any{
			"enabled":       false,
			"clientId":      nil,
			"clientSecret":  nil,
			"allowedEmails": []string{},
		},
		"pocketid": map[string]any{
			"enabled":        false,
			"clientId":       nil,
			"clientSecret":   nil,
			"plainDomain":    nil,
			"frontendDomain": nil,
			"allowedEmails":  []string{},
		},
		"yandex": map[string]any{
			"enabled":       false,
			"clientId":      nil,
			"clientSecret":  nil,
			"allowedEmails": []string{},
		},
		"keycloak": map[string]any{
			"enabled":        false,
			"realm":          nil,
			"clientId":       nil,
			"clientSecret":   nil,
			"keycloakDomain": nil,
			"frontendDomain": nil,
			"allowedEmails":  []string{},
		},
		"generic": map[string]any{
			"enabled":          false,
			"clientId":         nil,
			"clientSecret":     nil,
			"withPkce":         false,
			"authorizationUrl": nil,
			"tokenUrl":         nil,
			"frontendDomain":   nil,
			"allowedEmails":    []string{},
		},
		"telegram": map[string]any{
			"enabled":        false,
			"clientId":       nil,
			"clientSecret":   nil,
			"allowedIds":     []string{},
			"frontendDomain": nil,
		},
	}
}

const (
	DefaultBrandingSettingsJSON = `{"title":null,"logoUrl":null}`
	DefaultPasswordSettingsJSON = `{"enabled":true}`
	DefaultPasskeySettingsJSON  = `{"rpId":null,"origin":null,"enabled":false}`
	DefaultOAuth2SettingsJSON   = `{"github":{"enabled":false,"clientId":null,"clientSecret":null,"allowedEmails":[]},"yandex":{"enabled":false,"clientId":null,"clientSecret":null,"allowedEmails":[]},"generic":{"enabled":false,"clientId":null,"tokenUrl":null,"withPkce":false,"clientSecret":null,"allowedEmails":[],"frontendDomain":null,"authorizationUrl":null},"keycloak":{"realm":null,"enabled":false,"clientId":null,"clientSecret":null,"allowedEmails":[],"frontendDomain":null,"keycloakDomain":null},"pocketid":{"enabled":false,"clientId":null,"plainDomain":null,"frontendDomain":null,"clientSecret":null,"allowedEmails":[]},"telegram":{"enabled":false,"clientId":null,"clientSecret":null,"allowedIds":[],"frontendDomain":null}}`
)
