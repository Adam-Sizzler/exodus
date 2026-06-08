package server

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"time"

	"exodus-node/config"
)

var srsFileNameAllowedRe = regexp.MustCompile(`[^a-zA-Z0-9._-]+`)

const srsManifestName = "srs-lists.json"

type SRSListItem struct {
	Tag            string `json:"tag"`
	Format         string `json:"format"`
	URL            string `json:"url"`
	UpdateInterval string `json:"update_interval"`
	Path           string `json:"path,omitempty"`
}

type SRSSyncSummary struct {
	Total      int
	Configured int
	Downloaded int
	Failed     int
}

func normalizeSRSListsForSync(lists []SRSListItem, onInvalid func(error)) []SRSListItem {
	normalized := make([]SRSListItem, 0, len(lists))
	for _, raw := range lists {
		item, err := sanitizeSRSItem(raw)
		if err != nil {
			if onInvalid != nil {
				onInvalid(err)
			}
			continue
		}
		normalized = append(normalized, item)
	}
	return normalized
}

func srsListKey(item SRSListItem) string {
	return strings.Join([]string{
		item.Tag,
		item.Format,
		item.URL,
		item.UpdateInterval,
		item.Path,
	}, "\x00")
}

func areSRSListsEquivalent(a, b []SRSListItem) bool {
	if len(a) != len(b) {
		return false
	}
	if len(a) == 0 {
		return true
	}
	aKeys := make([]string, 0, len(a))
	for _, item := range a {
		aKeys = append(aKeys, srsListKey(item))
	}
	bKeys := make([]string, 0, len(b))
	for _, item := range b {
		bKeys = append(bKeys, srsListKey(item))
	}
	sort.Strings(aKeys)
	sort.Strings(bKeys)
	for i := range aKeys {
		if aKeys[i] != bKeys[i] {
			return false
		}
	}
	return true
}

func allSRSFilesExist(lists []SRSListItem) bool {
	for _, item := range lists {
		targetPath := filepath.Join(config.FixedSingboxDir, item.Path)
		if _, err := os.Stat(targetPath); err != nil {
			return false
		}
	}
	return true
}

func sanitizeSRSItem(item SRSListItem) (SRSListItem, error) {
	item.Tag = strings.TrimSpace(item.Tag)
	if item.Tag == "" {
		return item, fmt.Errorf("tag is required")
	}
	if strings.TrimSpace(item.URL) == "" {
		return item, fmt.Errorf("url is required")
	}
	u, err := url.Parse(item.URL)
	if err != nil || (u.Scheme != "http" && u.Scheme != "https") {
		return item, fmt.Errorf("invalid url %q", item.URL)
	}
	item.Format = strings.ToLower(strings.TrimSpace(item.Format))
	if item.Format == "" {
		item.Format = "binary"
	}
	item.UpdateInterval = strings.TrimSpace(item.UpdateInterval)
	if item.UpdateInterval == "" {
		item.UpdateInterval = "1d"
	}
	item.Path = sanitizeSRSFileName(item.Path)
	if item.Path == "" {
		item.Path = deriveSRSFileNameFromURL(item.URL)
	}
	if item.Path == "" {
		item.Path = fmt.Sprintf("%s.srs", item.Tag)
	}
	if !strings.HasSuffix(strings.ToLower(item.Path), ".srs") {
		item.Path += ".srs"
	}
	return item, nil
}

func deriveSRSFileNameFromURL(rawURL string) string {
	u, err := url.Parse(strings.TrimSpace(rawURL))
	if err != nil {
		return ""
	}
	return sanitizeSRSFileName(pathBase(u.Path))
}

func pathBase(rawPath string) string {
	base := strings.TrimSpace(filepath.Base(strings.TrimSpace(rawPath)))
	if base == "" || base == "." || base == "/" {
		return ""
	}
	return base
}

func sanitizeSRSFileName(name string) string {
	name = pathBase(name)
	if name == "" {
		return ""
	}
	name = srsFileNameAllowedRe.ReplaceAllString(name, "-")
	name = strings.Trim(name, "-._")
	return name
}

func applySRSListsToConfig(cfg map[string]any, lists []SRSListItem) {
	if cfg == nil {
		return
	}

	route := mapAny(cfg["route"])
	rawRuleSets, _ := route["rule_set"].([]any)
	incomingTags := make(map[string]struct{}, len(lists))
	incomingPaths := make(map[string]struct{}, len(lists))
	for _, item := range lists {
		normalized, err := sanitizeSRSItem(item)
		if err != nil {
			continue
		}
		incomingTags[normalized.Tag] = struct{}{}
		incomingPaths[normalized.Path] = struct{}{}
	}

	filtered := make([]any, 0, len(rawRuleSets))
	for _, raw := range rawRuleSets {
		rs, ok := raw.(map[string]any)
		if !ok {
			filtered = append(filtered, raw)
			continue
		}
		tag, _ := rs["tag"].(string)
		if strings.HasPrefix(strings.TrimSpace(tag), "exodus-") {
			continue
		}
		pathValue, _ := rs["path"].(string)
		if _, shouldReplace := incomingTags[strings.TrimSpace(tag)]; shouldReplace {
			continue
		}
		if _, shouldReplace := incomingPaths[strings.TrimSpace(pathValue)]; shouldReplace {
			continue
		}
		filtered = append(filtered, rs)
	}

	for _, item := range lists {
		normalized, err := sanitizeSRSItem(item)
		if err != nil {
			continue
		}
		localPath := filepath.Join(config.FixedSingboxDir, normalized.Path)
		if _, statErr := os.Stat(localPath); statErr != nil {
			continue
		}
		ruleSet := map[string]any{
			"type":   "local",
			"tag":    normalized.Tag,
			"format": normalized.Format,
			"path":   normalized.Path,
		}
		filtered = append(filtered, ruleSet)
	}

	route["rule_set"] = filtered
	cfg["route"] = route
}

func (s *NodeServer) SyncSRSLists(lists []SRSListItem) (SRSSyncSummary, error) {
	summary := SRSSyncSummary{Total: len(lists)}
	if len(lists) == 0 {
		return summary, nil
	}

	manifestPath := filepath.Join(config.FixedSingboxDir, srsManifestName)
	if err := os.MkdirAll(config.FixedSingboxDir, 0o755); err != nil {
		return summary, fmt.Errorf("create singbox dir: %w", err)
	}

	normalized := normalizeSRSListsForSync(lists, func(err error) {
		s.Cfg.LoggerFor("SRSService").Warn("Skip invalid SRS list", "error", err)
	})
	client := &http.Client{Timeout: 90 * time.Second}
	for _, item := range normalized {
		downloadStarted := time.Now()
		targetPath := filepath.Join(config.FixedSingboxDir, item.Path)
		if err := downloadSRSFile(client, item.URL, targetPath); err != nil {
			summary.Failed++
			s.Cfg.LoggerFor("SRSService").Warn("Failed to download SRS file", "tag", item.Tag, "url", item.URL, "path", item.Path, "error", err, "duration_ms", time.Since(downloadStarted).Milliseconds())
			continue
		}
		summary.Downloaded++
		summary.Configured++
		s.Cfg.LoggerFor("SRSService").Debug("SRS file downloaded", "tag", item.Tag, "url", item.URL, "path", item.Path, "duration_ms", time.Since(downloadStarted).Milliseconds())
	}

	manifestBytes, err := json.MarshalIndent(map[string]any{"srs_lists": normalized}, "", "  ")
	if err == nil {
		_ = os.WriteFile(manifestPath, manifestBytes, 0o644)
	}

	return summary, nil
}

// SyncSRSListsIfChanged skips downloads when manifest has the same SRS set and all files already exist.
// Returns summary, whether download was performed, and error.
func (s *NodeServer) SyncSRSListsIfChanged(lists []SRSListItem) (SRSSyncSummary, bool, error) {
	summary := SRSSyncSummary{Total: len(lists)}
	if len(lists) == 0 {
		return summary, false, nil
	}

	normalized := normalizeSRSListsForSync(lists, func(err error) {
		s.Cfg.LoggerFor("SRSService").Warn("Skip invalid SRS list", "error", err)
	})
	summary.Configured = len(normalized)

	currentManifest, err := loadSRSManifest()
	if err == nil {
		currentNormalized := normalizeSRSListsForSync(currentManifest, nil)
		if areSRSListsEquivalent(normalized, currentNormalized) && allSRSFilesExist(normalized) {
			return summary, false, nil
		}
	}

	summary, err = s.SyncSRSLists(lists)
	if err != nil {
		return summary, false, err
	}
	return summary, true, nil
}

func loadSRSManifest() ([]SRSListItem, error) {
	manifestPath := filepath.Join(config.FixedSingboxDir, srsManifestName)
	payload, err := os.ReadFile(manifestPath)
	if err != nil {
		return nil, err
	}

	var parsed struct {
		SRSLists []SRSListItem `json:"srs_lists"`
	}
	if err := json.Unmarshal(payload, &parsed); err != nil {
		return nil, fmt.Errorf("parse %s: %w", manifestPath, err)
	}
	return parsed.SRSLists, nil
}

func downloadSRSFile(client *http.Client, rawURL, targetPath string) error {
	tmpPath := targetPath + ".tmp"
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, rawURL, nil)
	if err != nil {
		return fmt.Errorf("build request: %w", err)
	}

	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("unexpected status: %d", resp.StatusCode)
	}

	f, err := os.OpenFile(tmpPath, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0o644)
	if err != nil {
		return fmt.Errorf("open temp file: %w", err)
	}
	defer func() {
		_ = f.Close()
	}()

	if _, err := io.Copy(f, resp.Body); err != nil {
		_ = os.Remove(tmpPath)
		return fmt.Errorf("write file: %w", err)
	}

	if err := f.Sync(); err != nil {
		_ = os.Remove(tmpPath)
		return fmt.Errorf("fsync file: %w", err)
	}

	if err := os.Rename(tmpPath, targetPath); err != nil {
		_ = os.Remove(tmpPath)
		return fmt.Errorf("replace file: %w", err)
	}

	return nil
}
