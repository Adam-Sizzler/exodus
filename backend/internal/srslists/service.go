package srslists

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/exodus/subscription-page/backend/internal/logger"
)

const (
	defaultStorageDir  = "/opt/app/ruleset"
	storageDirEnv      = "SUB_SRS_DIR"
	manifestName       = "srs-lists.json"
	refreshInterval    = 7 * 24 * time.Hour
	rulesetRoutePrefix = "/ruleset"
	contentTypeRuleset = "application/octet-stream"
)

var (
	fileNameAllowedRe = regexp.MustCompile(`[^a-zA-Z0-9._-]+`)
	updaterOnce       sync.Once
)

type ListItem struct {
	Tag            string `json:"tag"`
	Format         string `json:"format"`
	URL            string `json:"url"`
	UpdateInterval string `json:"update_interval"`
	Path           string `json:"path,omitempty"`
}

type SyncSummary struct {
	Total      int
	Configured int
	Downloaded int
	Failed     int
}

func StorageDir() string {
	if value := strings.TrimSpace(os.Getenv(storageDirEnv)); value != "" {
		return filepath.Clean(value)
	}
	return defaultStorageDir
}

func IsRulesetRoute(routePath string) bool {
	cleaned := cleanRoutePath(routePath)
	return cleaned == rulesetRoutePrefix || strings.HasPrefix(cleaned, rulesetRoutePrefix+"/")
}

func ServeHTTP(w http.ResponseWriter, r *http.Request, routePath string) {
	cleaned := cleanRoutePath(routePath)
	fileName := strings.TrimPrefix(cleaned, rulesetRoutePrefix)
	fileName = strings.Trim(fileName, "/")
	fileName = sanitizeFileName(fileName)
	if fileName == "" || !strings.HasSuffix(strings.ToLower(fileName), ".srs") {
		http.NotFound(w, r)
		return
	}

	root := filepath.Clean(StorageDir())
	fullPath := filepath.Clean(filepath.Join(root, fileName))
	if fullPath != root && !strings.HasPrefix(fullPath, root+string(filepath.Separator)) {
		http.NotFound(w, r)
		return
	}

	info, err := os.Stat(fullPath)
	if err != nil || info.IsDir() {
		http.NotFound(w, r)
		return
	}

	w.Header().Set("Content-Type", contentTypeRuleset)
	w.Header().Set("Cache-Control", "public, max-age=3600")
	http.ServeFile(w, r, fullPath)
}

func StartAutoUpdater(ctx context.Context) {
	updaterOnce.Do(func() {
		log := logger.WithContext("SRSService")
		log.Info("SRS auto updater started", logger.String("interval", refreshInterval.String()), logger.String("dir", StorageDir()))
		go func() {
			RefreshFromManifest(ctx)

			ticker := time.NewTicker(refreshInterval)
			defer ticker.Stop()

			for {
				select {
				case <-ctx.Done():
					log.Info("SRS auto updater stopped")
					return
				case <-ticker.C:
					RefreshFromManifest(ctx)
				}
			}
		}()
	})
}

func RefreshFromManifest(ctx context.Context) {
	log := logger.WithContext("SRSService")
	manifestPath := filepath.Join(StorageDir(), manifestName)

	lists, err := loadManifest()
	if err != nil {
		if os.IsNotExist(err) {
			log.Debug("SRS manifest not found, skipping refresh", logger.String("path", manifestPath))
			return
		}
		log.Warn("Failed to read SRS manifest for refresh", logger.String("path", manifestPath), logger.String("error", err.Error()))
		return
	}
	if len(lists) == 0 {
		log.Debug("SRS manifest is empty, skipping refresh", logger.String("path", manifestPath))
		return
	}

	summary, syncErr := SyncLists(ctx, lists)
	if syncErr != nil {
		log.Warn("SRS refresh failed", logger.String("error", syncErr.Error()))
		return
	}

	log.Info(
		"SRS refresh completed",
		logger.Int("total", summary.Total),
		logger.Int("configured", summary.Configured),
		logger.Int("downloaded", summary.Downloaded),
		logger.Int("failed", summary.Failed),
	)
}

func SyncLists(ctx context.Context, lists []ListItem) (SyncSummary, error) {
	summary := SyncSummary{Total: len(lists)}
	if len(lists) == 0 {
		return summary, nil
	}

	root := StorageDir()
	manifestPath := filepath.Join(root, manifestName)
	if err := os.MkdirAll(root, 0o755); err != nil {
		return summary, fmt.Errorf("create srs dir: %w", err)
	}

	log := logger.WithContext("SRSService")
	normalized := normalizeLists(lists, func(err error) {
		log.Warn("Skip invalid SRS list", logger.String("error", err.Error()))
	})
	summary.Configured = len(normalized)

	client := &http.Client{Timeout: 90 * time.Second}
	for _, item := range normalized {
		if ctx != nil {
			select {
			case <-ctx.Done():
				return summary, ctx.Err()
			default:
			}
		}

		downloadStarted := time.Now()
		targetPath := filepath.Join(root, item.Path)
		if err := downloadFile(ctx, client, item.URL, targetPath); err != nil {
			summary.Failed++
			log.Warn(
				"Failed to download SRS file",
				logger.String("tag", item.Tag),
				logger.String("url", item.URL),
				logger.String("path", item.Path),
				logger.String("error", err.Error()),
				logger.Any("duration_ms", time.Since(downloadStarted).Milliseconds()),
			)
			continue
		}
		summary.Downloaded++
		log.Debug(
			"SRS file downloaded",
			logger.String("tag", item.Tag),
			logger.String("url", item.URL),
			logger.String("path", item.Path),
			logger.Any("duration_ms", time.Since(downloadStarted).Milliseconds()),
		)
	}

	manifestBytes, err := json.MarshalIndent(map[string]any{"srs_lists": normalized}, "", "  ")
	if err == nil {
		_ = os.WriteFile(manifestPath, manifestBytes, 0o644)
	}

	return summary, nil
}

func NormalizeListsForComparison(lists []ListItem) []ListItem {
	return normalizeLists(lists, nil)
}

func ListsEquivalent(a, b []ListItem) bool {
	a = NormalizeListsForComparison(a)
	b = NormalizeListsForComparison(b)
	if len(a) != len(b) {
		return false
	}
	if len(a) == 0 {
		return true
	}
	aKeys := make([]string, 0, len(a))
	for _, item := range a {
		aKeys = append(aKeys, listKey(item))
	}
	bKeys := make([]string, 0, len(b))
	for _, item := range b {
		bKeys = append(bKeys, listKey(item))
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

func normalizeLists(lists []ListItem, onInvalid func(error)) []ListItem {
	normalized := make([]ListItem, 0, len(lists))
	for _, raw := range lists {
		item, err := sanitizeItem(raw)
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

func sanitizeItem(item ListItem) (ListItem, error) {
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
	item.Path = sanitizeFileName(item.Path)
	if item.Path == "" {
		item.Path = deriveFileNameFromURL(item.URL)
	}
	if item.Path == "" {
		item.Path = fmt.Sprintf("%s.srs", item.Tag)
	}
	if !strings.HasSuffix(strings.ToLower(item.Path), ".srs") {
		item.Path += ".srs"
	}
	return item, nil
}

func deriveFileNameFromURL(rawURL string) string {
	u, err := url.Parse(strings.TrimSpace(rawURL))
	if err != nil {
		return ""
	}
	return sanitizeFileName(path.Base(u.Path))
}

func sanitizeFileName(name string) string {
	name = strings.TrimSpace(name)
	if name == "" {
		return ""
	}
	name = path.Base(strings.ReplaceAll(name, "\\", "/"))
	if name == "" || name == "." || name == "/" {
		return ""
	}
	name = fileNameAllowedRe.ReplaceAllString(name, "-")
	name = strings.Trim(name, "-._")
	return name
}

func listKey(item ListItem) string {
	return strings.Join([]string{
		item.Tag,
		item.Format,
		item.URL,
		item.UpdateInterval,
		item.Path,
	}, "\x00")
}

func loadManifest() ([]ListItem, error) {
	manifestPath := filepath.Join(StorageDir(), manifestName)
	payload, err := os.ReadFile(manifestPath)
	if err != nil {
		return nil, err
	}

	var parsed struct {
		SRSLists []ListItem `json:"srs_lists"`
	}
	if err := json.Unmarshal(payload, &parsed); err != nil {
		return nil, fmt.Errorf("parse %s: %w", manifestPath, err)
	}
	return parsed.SRSLists, nil
}

func downloadFile(ctx context.Context, client *http.Client, rawURL, targetPath string) error {
	if ctx == nil {
		ctx = context.Background()
	}
	downloadCtx, cancel := context.WithTimeout(ctx, 2*time.Minute)
	defer cancel()

	req, err := http.NewRequestWithContext(downloadCtx, http.MethodGet, rawURL, nil)
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

	if err := os.MkdirAll(filepath.Dir(targetPath), 0o755); err != nil {
		return fmt.Errorf("create target dir: %w", err)
	}

	tmpPath := targetPath + ".tmp"
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

func cleanRoutePath(rawPath string) string {
	if rawPath == "" {
		return "/"
	}
	cleaned := path.Clean("/" + rawPath)
	if cleaned == "." {
		return "/"
	}
	return cleaned
}
