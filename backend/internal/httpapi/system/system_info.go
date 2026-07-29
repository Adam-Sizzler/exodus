package system

import (
	"bufio"
	"fmt"
	"math/big"
	"net/url"
	"os"
	"regexp"
	"runtime/debug"
	"strconv"
	"strings"
	"time"

	"exodus/internal/constant"
)

type memStats struct {
	total     int64
	free      int64
	used      int64
	active    int64
	available int64
}

func normalizeVersion(value string) string {
	candidates := []string{value, constant.Version}
	for _, raw := range candidates {
		normalized, ok := normalizeVersionCandidate(raw)
		if ok {
			return normalized
		}
	}
	return "unknown"
}

var gitDescribeVersionPattern = regexp.MustCompile(`^(.+)-\d+-g[0-9a-f]{7,}(?:-dirty)?$`)

func normalizeVersionCandidate(value string) (string, bool) {
	raw := strings.TrimSpace(value)
	if raw == "" {
		return "", false
	}
	if strings.EqualFold(raw, "latest") || strings.EqualFold(raw, "unknown") || raw == "(devel)" {
		return "", false
	}

	if matches := gitDescribeVersionPattern.FindStringSubmatch(raw); len(matches) == 2 {
		base := strings.TrimSpace(matches[1])
		if isKnownMetadataValue(base) && !strings.EqualFold(base, "latest") {
			return base, true
		}
	}

	return raw, true
}

func getLastTwoDaysRanges(now time.Time) [2]usageRange {
	today := dayRange(now, 0)
	yesterday := dayRange(now, 1)
	return [2]usageRange{yesterday, today}
}

func getLastSevenDaysRanges(now time.Time) [2]usageRange {
	currentStart := startOfDay(now.AddDate(0, 0, -6))
	currentEnd := endOfDay(now)
	previousEnd := currentStart.Add(-time.Nanosecond)
	previousStart := startOfDay(previousEnd.AddDate(0, 0, -6))

	return [2]usageRange{
		{start: previousStart.UTC(), end: previousEnd.UTC()},
		{start: currentStart.UTC(), end: currentEnd.UTC()},
	}
}

func getLast30DaysRanges(now time.Time) [2]usageRange {
	currentStart := startOfDay(now.AddDate(0, 0, -29))
	currentEnd := endOfDay(now)
	previousEnd := currentStart.Add(-time.Nanosecond)
	previousStart := startOfDay(previousEnd.AddDate(0, 0, -29))

	return [2]usageRange{
		{start: previousStart.UTC(), end: previousEnd.UTC()},
		{start: currentStart.UTC(), end: currentEnd.UTC()},
	}
}

func getCalendarMonthRanges(now time.Time) [2]usageRange {
	currentStart := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
	currentEnd := endOfDay(now)

	previousMonthNow := now.AddDate(0, -1, 0)
	previousStart := time.Date(previousMonthNow.Year(), previousMonthNow.Month(), 1, 0, 0, 0, 0, now.Location())
	previousEnd := currentStart.Add(-time.Nanosecond)

	return [2]usageRange{
		{start: previousStart.UTC(), end: previousEnd.UTC()},
		{start: currentStart.UTC(), end: currentEnd.UTC()},
	}
}

func getCalendarYearRanges(now time.Time) [2]usageRange {
	currentStart := time.Date(now.Year(), time.January, 1, 0, 0, 0, 0, now.Location())
	currentEnd := endOfDay(now)
	previousStart := time.Date(now.Year()-1, time.January, 1, 0, 0, 0, 0, now.Location())
	previousEnd := currentStart.Add(-time.Nanosecond)

	return [2]usageRange{
		{start: previousStart.UTC(), end: previousEnd.UTC()},
		{start: currentStart.UTC(), end: currentEnd.UTC()},
	}
}

func dayRange(now time.Time, subtractDays int) usageRange {
	target := now.AddDate(0, 0, -subtractDays)
	start := startOfDay(target).UTC()
	end := endOfDay(target).UTC()
	return usageRange{start: start, end: end}
}

func startOfDay(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
}

func endOfDay(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), 23, 59, 59, int(time.Second-time.Nanosecond), t.Location())
}

func resolveLocation(requestTZ string) *time.Location {
	tz := strings.TrimSpace(requestTZ)
	if tz == "" {
		return time.UTC
	}

	loc, err := time.LoadLocation(tz)
	if err != nil {
		return time.UTC
	}
	return loc
}

func readBuildMetadata() (version string, commitSHA string, branch string, buildTime string) {
	version = "unknown"
	commitSHA = "unknown"
	branch = "unknown"
	buildTime = "unknown"

	if info, ok := debug.ReadBuildInfo(); ok {
		if info.Main.Version != "" && info.Main.Version != "(devel)" {
			version = info.Main.Version
		}
		for _, setting := range info.Settings {
			switch setting.Key {
			case "vcs.revision":
				if setting.Value != "" {
					commitSHA = setting.Value
				}
			case "vcs.branch":
				if setting.Value != "" {
					branch = setting.Value
				}
			case "vcs.time":
				if setting.Value != "" {
					buildTime = setting.Value
				}
			}
		}
	}

	if value := firstMetadataEnv("EXODUS_VERSION"); value != "" {
		version = value
	}
	if value := firstMetadataEnv("EXODUS_BACKEND_COMMIT", "EXODUS_REVISION", "EXODUS_COMMIT", "GITHUB_SHA"); value != "" {
		commitSHA = value
	}
	if value := firstMetadataEnv("EXODUS_GIT_BRANCH", "EXODUS_BRANCH", "GITHUB_REF_NAME", "GITHUB_HEAD_REF"); value != "" {
		branch = value
	}
	if value := firstMetadataEnv("EXODUS_BUILD_TIME"); value != "" {
		buildTime = value
	}

	if commitSHA == "unknown" && isKnownMetadataValue(constant.Revision) {
		commitSHA = constant.Revision
	}

	return version, commitSHA, branch, buildTime
}

func readRepositoryURL() string {
	if value := normalizeRepositoryURL(firstMetadataEnv("EXODUS_REPOSITORY_URL", "EXODUS_GIT_REMOTE", "GITHUB_REPOSITORY_URL")); value != "" {
		return value
	}

	repository := firstMetadataEnv("GITHUB_REPOSITORY")
	if repository == "" {
		return "unknown"
	}

	serverURL := firstMetadataEnv("GITHUB_SERVER_URL")
	if serverURL == "" {
		serverURL = "https://github.com"
	}

	value := normalizeRepositoryURL(strings.TrimRight(serverURL, "/") + "/" + repository)
	if value == "" {
		return "unknown"
	}
	return value
}

func firstMetadataEnv(keys ...string) string {
	for _, key := range keys {
		value := strings.TrimSpace(os.Getenv(key))
		if isKnownMetadataValue(value) {
			return value
		}
	}
	return ""
}

func isKnownMetadataValue(value string) bool {
	trimmed := strings.TrimSpace(value)
	if trimmed == "" {
		return false
	}
	return !strings.EqualFold(trimmed, "unknown") && trimmed != "(devel)"
}

func normalizeRepositoryURL(raw string) string {
	value := strings.TrimSpace(raw)
	if !isKnownMetadataValue(value) {
		return ""
	}

	value = strings.TrimPrefix(value, "git+")
	switch {
	case strings.HasPrefix(value, "git@"):
		parts := strings.SplitN(strings.TrimPrefix(value, "git@"), ":", 2)
		if len(parts) == 2 {
			value = "https://" + parts[0] + "/" + parts[1]
		}
	case strings.HasPrefix(value, "ssh://git@"):
		value = "https://" + strings.TrimPrefix(value, "ssh://git@")
	case !strings.Contains(value, "://"):
		host, _, ok := strings.Cut(value, "/")
		if ok && strings.Contains(host, ".") {
			value = "https://" + value
		}
	}

	if parsed, err := url.Parse(value); err == nil && parsed.Scheme != "" && parsed.Host != "" {
		parsed.User = nil
		value = parsed.String()
	}

	value = strings.TrimRight(value, "/")
	value = strings.TrimSuffix(value, ".git")
	return value
}

func buildCommitURL(repositoryURL, sha string) string {
	trimmed := strings.TrimSpace(sha)
	if !isKnownMetadataValue(trimmed) {
		return "unknown"
	}

	repositoryURL = normalizeRepositoryURL(repositoryURL)
	if repositoryURL == "" {
		return "unknown"
	}
	return repositoryURL + "/commit/" + trimmed
}

func readMemStats() memStats {
	result := memStats{}
	file, err := os.Open("/proc/meminfo")
	if err != nil {
		return result
	}
	defer file.Close()

	values := map[string]int64{}
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.SplitN(line, ":", 2)
		if len(parts) != 2 {
			continue
		}
		key := strings.TrimSpace(parts[0])
		fields := strings.Fields(strings.TrimSpace(parts[1]))
		if len(fields) == 0 {
			continue
		}
		amountKB, err := strconv.ParseInt(fields[0], 10, 64)
		if err != nil {
			continue
		}
		values[key] = amountKB * 1024
	}

	result.total = values["MemTotal"]
	result.free = values["MemFree"]
	result.available = values["MemAvailable"]
	result.active = values["Active"]

	if err := scanner.Err(); err != nil {
		return result
	}

	if result.available == 0 {
		result.available = result.free
	}
	result.used = result.total - result.available
	if result.used < 0 {
		result.used = 0
	}
	if result.used > result.total {
		result.used = result.total
	}

	return result
}

func readUptimeSeconds() int64 {
	data, err := os.ReadFile("/proc/uptime")
	if err != nil {
		return 0
	}
	fields := strings.Fields(string(data))
	if len(fields) == 0 {
		return 0
	}
	value, err := strconv.ParseFloat(fields[0], 64)
	if err != nil {
		return 0
	}
	return int64(value)
}

func detectPhysicalCores(fallback int) int {
	file, err := os.Open("/proc/cpuinfo")
	if err != nil {
		return fallback
	}
	defer file.Close()

	physicalCoresByPackage := map[string]struct{}{}
	scanner := bufio.NewScanner(file)

	var packageID, coreID string
	flush := func() {
		if packageID != "" && coreID != "" {
			physicalCoresByPackage[packageID+":"+coreID] = struct{}{}
		}
		packageID = ""
		coreID = ""
	}

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			flush()
			continue
		}
		parts := strings.SplitN(line, ":", 2)
		if len(parts) != 2 {
			continue
		}
		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])
		switch key {
		case "physical id":
			packageID = value
		case "core id":
			coreID = value
		}
	}
	flush()

	if err := scanner.Err(); err != nil {
		return fallback
	}

	if len(physicalCoresByPackage) == 0 {
		return fallback
	}
	return len(physicalCoresByPackage)
}

func formatBigBytes(value *big.Int) string {
	if value == nil || value.Sign() == 0 {
		return "0 B"
	}

	sign := ""
	abs := new(big.Int).Set(value)
	if abs.Sign() < 0 {
		sign = "-"
		abs.Abs(abs)
	}

	byteFloat, _ := new(big.Float).SetInt(abs).Float64()
	units := []string{"B", "KiB", "MiB", "GiB", "TiB", "PiB", "EiB"}
	idx := 0
	for byteFloat >= 1024 && idx < len(units)-1 {
		byteFloat /= 1024
		idx++
	}

	return fmt.Sprintf("%s%.2f %s", sign, byteFloat, units[idx])
}
