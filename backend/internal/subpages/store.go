package subpages

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/cerberus/subscription-page/backend/internal/config"
	"github.com/cerberus/subscription-page/backend/internal/cerberus"
	"github.com/cerberus/subscription-page/backend/internal/security"
)

type BaseSettings struct {
	MetaTitle          string
	MetaDescription    string
	ShowConnectionKeys bool
	HideGetLinkButton  bool
}

type storedConfig struct {
	Raw          json.RawMessage
	BaseSettings BaseSettings
}

type Store struct {
	subpageConfigUUID string
	sessionSecret     string
	configs           map[string]storedConfig
}

func NewStore(cfg config.Config) *Store {
	return &Store{
		subpageConfigUUID: cfg.SubpageConfigUUID,
		sessionSecret:     cfg.SessionSecret,
		configs:           make(map[string]storedConfig),
	}
}

func (s *Store) Load(ctx context.Context, client *cerberus.Client) error {
	configIDs, err := client.GetSubscriptionPageConfigList(ctx)
	if err != nil {
		return fmt.Errorf("subscription page config list cannot be fetched: %w", err)
	}

	if len(configIDs) == 0 {
		return fmt.Errorf("subscription page config list is empty")
	}

	log.Printf("Found %d subscription page configs.", len(configIDs))

	for _, configID := range configIDs {
		rawConfig, err := client.GetSubscriptionPageConfigByUUID(ctx, configID)
		if err != nil {
			return fmt.Errorf("error while fetching subpage config %s: %w", configID, err)
		}

		s.configs[configID] = storedConfig{
			Raw:          rawConfig,
			BaseSettings: extractBaseSettings(rawConfig),
		}

		log.Printf("[OK] %s", configID)
	}

	if len(s.configs) == 0 {
		return fmt.Errorf("at least one subpage config must be valid")
	}

	log.Printf("[OK] Subpage configs are loaded successfully.")
	return nil
}

func (s *Store) ResolveRawConfig(encryptedSubpageConfigUUID string) (json.RawMessage, error) {
	decrypted, err := security.DecryptUUID(encryptedSubpageConfigUUID, s.sessionSecret)
	if err != nil {
		return nil, err
	}

	subpageConfig, ok := s.configs[decrypted]
	if !ok {
		return nil, fmt.Errorf("subpage config %s not found", decrypted)
	}

	return subpageConfig.Raw, nil
}

func (s *Store) EncryptResolvedUUID(subpageConfigUUID string) (string, error) {
	return security.EncryptUUID(s.finalUUID(subpageConfigUUID), s.sessionSecret)
}

func (s *Store) BaseSettingsFor(subpageConfigUUID string) BaseSettings {
	finalUUID := s.finalUUID(subpageConfigUUID)
	if subpageConfig, ok := s.configs[finalUUID]; ok {
		return subpageConfig.BaseSettings
	}

	return defaultBaseSettings()
}

func (s *Store) finalUUID(subpageConfigUUID string) string {
	if s.subpageConfigUUID == config.SubpageDefaultConfigUUID && subpageConfigUUID != "" {
		return subpageConfigUUID
	}

	return s.subpageConfigUUID
}

func extractBaseSettings(rawConfig json.RawMessage) BaseSettings {
	type envelope struct {
		BaseSettings struct {
			MetaTitle          string `json:"metaTitle"`
			MetaDescription    string `json:"metaDescription"`
			ShowConnectionKeys bool   `json:"showConnectionKeys"`
			HideGetLinkButton  bool   `json:"hideGetLinkButton"`
		} `json:"baseSettings"`
	}

	var parsed envelope
	if err := json.Unmarshal(rawConfig, &parsed); err != nil {
		return defaultBaseSettings()
	}

	metaTitle := parsed.BaseSettings.MetaTitle
	if metaTitle == "" {
		metaTitle = "Subscription Page"
	}

	metaDescription := parsed.BaseSettings.MetaDescription
	if metaDescription == "" {
		metaDescription = "Subscription Page"
	}

	return BaseSettings{
		MetaTitle:          metaTitle,
		MetaDescription:    metaDescription,
		ShowConnectionKeys: parsed.BaseSettings.ShowConnectionKeys,
		HideGetLinkButton:  parsed.BaseSettings.HideGetLinkButton,
	}
}

func defaultBaseSettings() BaseSettings {
	return BaseSettings{
		MetaTitle:          "Subscription Page",
		MetaDescription:    "Subscription Page",
		ShowConnectionKeys: false,
		HideGetLinkButton:  false,
	}
}
