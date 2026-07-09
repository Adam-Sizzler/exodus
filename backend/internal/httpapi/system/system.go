package system

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"math/big"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"runtime"
	"runtime/debug"
	"strconv"
	"strings"
	"time"

	"exodus/internal/config"
	"exodus/internal/constant"
	dbmanager "exodus/internal/db/manager"
	"exodus/internal/httpapi/shared"
	"exodus/internal/nodehotcache"
)

type usageRange struct {
	start time.Time
	end   time.Time
}

type bandwidthStat struct {
	Current    string `json:"current"`
	Difference string `json:"difference"`
	Previous   string `json:"previous"`
}

type usersRecap struct {
	total             int64
	newUsersThisMonth int64
}

type nodesRecap struct {
	total             int64
	totalRam          string
	totalCpuCores     int64
	distinctCountries int64
}

type nodesMetricsResponse struct {
	Nodes []nodeMetricsItem `json:"nodes"`
}

type nodeMetricsItem struct {
	NodeUUID       string            `json:"nodeUuid"`
	NodeName       string            `json:"nodeName"`
	CountryEmoji   string            `json:"countryEmoji"`
	ProviderName   string            `json:"providerName"`
	UsersOnline    int               `json:"usersOnline"`
	InboundsStats  []nodeTrafficStat `json:"inboundsStats"`
	OutboundsStats []nodeTrafficStat `json:"outboundsStats"`
}

type nodeTrafficStat struct {
	Tag      string `json:"tag"`
	Upload   string `json:"upload"`
	Download string `json:"download"`
}

func MetadataHandler(cfg *config.BackendConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			shared.WriteJSONError(w, http.StatusMethodNotAllowed, "method not allowed")
			return
		}

		version, backendSHA, branch, buildTime := readBuildMetadata()
		if version == "" || version == "unknown" {
			version = constant.Version
		}
		version = normalizeVersion(version)

		frontendSHA := firstMetadataEnv("EXODUS_FRONTEND_COMMIT")
		if frontendSHA == "" {
			frontendSHA = backendSHA
		}
		buildNumber := firstMetadataEnv("EXODUS_BUILD_NUMBER", "BUILD_NUMBER", "GITHUB_RUN_NUMBER")
		if buildNumber == "" {
			buildNumber = "unknown"
		}

		repositoryURL := readRepositoryURL()
		backendCommitURL := buildCommitURL(repositoryURL, backendSHA)
		frontendCommitURL := buildCommitURL(repositoryURL, frontendSHA)

		payload := map[string]any{
			"response": map[string]any{
				"version": version,
				"build": map[string]string{
					"time":   buildTime,
					"number": buildNumber,
				},
				"git": map[string]any{
					"repositoryUrl": repositoryURL,
					"backend": map[string]string{
						"commitSha": backendSHA,
						"branch":    branch,
						"commitUrl": backendCommitURL,
					},
					"frontend": map[string]string{
						"commitSha": frontendSHA,
						"commitUrl": frontendCommitURL,
					},
				},
			},
		}

		cfg.Logger.Trace("System metadata requested", "remote_addr", r.RemoteAddr, "version", version)
		shared.WriteJSON(w, http.StatusOK, payload)
	}
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

func StatsHandler(manager *dbmanager.DatabaseManager, cfg *config.BackendConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			shared.WriteJSONError(w, http.StatusMethodNotAllowed, "method not allowed")
			return
		}

		cpuCores := runtime.NumCPU()
		physicalCores := detectPhysicalCores(cpuCores)
		mem := readMemStats()
		uptime := readUptimeSeconds()
		timestamp := time.Now().UnixMilli()

		statusCounts, totalUsers, err := readUsersStatusStats(r.Context(), manager)
		if err != nil {
			shared.SendError(w, http.StatusInternalServerError, "failed to read users stats", err, cfg)
			return
		}

		onlineStats, err := readOnlineStats(r.Context(), manager)
		if err != nil {
			shared.SendError(w, http.StatusInternalServerError, "failed to read online stats", err, cfg)
			return
		}

		totalOnline, err := readTotalOnlineOnNodes(r.Context(), manager, cfg)
		if err != nil {
			shared.SendError(w, http.StatusInternalServerError, "failed to read nodes online stats", err, cfg)
			return
		}

		lifetimeBytes, err := readLifetimeTrafficBytes(r.Context(), manager)
		if err != nil {
			shared.SendError(w, http.StatusInternalServerError, "failed to read lifetime traffic", err, cfg)
			return
		}

		payload := map[string]any{
			"response": map[string]any{
				"cpu": map[string]int{
					"cores":         cpuCores,
					"physicalCores": physicalCores,
				},
				"memory": map[string]int64{
					"total": mem.total,
					"free":  mem.free,
					"used":  mem.used,
				},
				"uptime":    uptime,
				"timestamp": timestamp,
				"users": map[string]any{
					"statusCounts": statusCounts,
					"totalUsers":   totalUsers,
				},
				"onlineStats": map[string]int64{
					"lastDay":     onlineStats.lastDay,
					"lastWeek":    onlineStats.lastWeek,
					"neverOnline": onlineStats.neverOnline,
					"onlineNow":   onlineStats.onlineNow,
				},
				"nodes": map[string]any{
					"totalOnline":        totalOnline,
					"totalBytesLifetime": lifetimeBytes,
				},
			},
		}

		cfg.Logger.Trace("System stats requested", "remote_addr", r.RemoteAddr, "total_users", totalUsers)
		shared.WriteJSON(w, http.StatusOK, payload)
	}
}

func RecapHandler(manager *dbmanager.DatabaseManager, cfg *config.BackendConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			shared.WriteJSONError(w, http.StatusMethodNotAllowed, "method not allowed")
			return
		}

		now := time.Now().UTC()
		startOfMonth := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, time.UTC)
		endOfMonth := startOfMonth.AddDate(0, 1, 0).Add(-time.Nanosecond)

		users, err := readUsersRecap(r.Context(), manager, startOfMonth)
		if err != nil {
			shared.SendError(w, http.StatusInternalServerError, "failed to read users recap", err, cfg)
			return
		}

		nodes, err := readNodesRecap(r.Context(), manager, cfg)
		if err != nil {
			shared.SendError(w, http.StatusInternalServerError, "failed to read nodes recap", err, cfg)
			return
		}

		lifetimeTraffic, err := readLifetimeTrafficBytes(r.Context(), manager)
		if err != nil {
			shared.SendError(w, http.StatusInternalServerError, "failed to read lifetime traffic recap", err, cfg)
			return
		}

		monthTraffic, err := readUsageBytesTextByRange(r.Context(), manager, usageRange{start: startOfMonth, end: endOfMonth})
		if err != nil {
			shared.SendError(w, http.StatusInternalServerError, "failed to read month traffic recap", err, cfg)
			return
		}

		initDate, err := readInitDate(r.Context(), manager)
		if err != nil {
			shared.SendError(w, http.StatusInternalServerError, "failed to read init date recap", err, cfg)
			return
		}

		version, _, _, _ := readBuildMetadata()
		if version == "" || version == "unknown" {
			version = constant.Version
		}
		version = normalizeVersion(version)

		payload := map[string]any{
			"response": map[string]any{
				"thisMonth": map[string]any{
					"users":   users.newUsersThisMonth,
					"traffic": monthTraffic,
				},
				"total": map[string]any{
					"users":             users.total,
					"nodes":             nodes.total,
					"traffic":           lifetimeTraffic,
					"nodesRam":          nodes.totalRam,
					"nodesCpuCores":     nodes.totalCpuCores,
					"distinctCountries": nodes.distinctCountries,
				},
				"version":  version,
				"initDate": initDate.UTC().Format(time.RFC3339),
			},
		}

		cfg.Logger.Trace("System recap requested", "remote_addr", r.RemoteAddr, "total_users", users.total, "total_nodes", nodes.total)
		shared.WriteJSON(w, http.StatusOK, payload)
	}
}

func BandwidthStatsHandler(manager *dbmanager.DatabaseManager, cfg *config.BackendConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			shared.WriteJSONError(w, http.StatusMethodNotAllowed, "method not allowed")
			return
		}

		tz := strings.TrimSpace(r.URL.Query().Get("tz"))
		loc := resolveLocation(tz)
		now := time.Now().In(loc)

		lastTwoDays, err := readUsageComparison(r.Context(), manager, getLastTwoDaysRanges(now))
		if err != nil {
			shared.SendError(w, http.StatusInternalServerError, "failed to read last two days bandwidth", err, cfg)
			return
		}
		lastSevenDays, err := readUsageComparison(r.Context(), manager, getLastSevenDaysRanges(now))
		if err != nil {
			shared.SendError(w, http.StatusInternalServerError, "failed to read last seven days bandwidth", err, cfg)
			return
		}
		last30Days, err := readUsageComparison(r.Context(), manager, getLast30DaysRanges(now))
		if err != nil {
			shared.SendError(w, http.StatusInternalServerError, "failed to read last 30 days bandwidth", err, cfg)
			return
		}
		calendarMonth, err := readUsageComparison(r.Context(), manager, getCalendarMonthRanges(now))
		if err != nil {
			shared.SendError(w, http.StatusInternalServerError, "failed to read calendar month bandwidth", err, cfg)
			return
		}
		currentYear, err := readUsageComparison(r.Context(), manager, getCalendarYearRanges(now))
		if err != nil {
			shared.SendError(w, http.StatusInternalServerError, "failed to read current year bandwidth", err, cfg)
			return
		}

		payload := map[string]any{
			"response": map[string]any{
				"bandwidthLastTwoDays":   lastTwoDays,
				"bandwidthLastSevenDays": lastSevenDays,
				"bandwidthLast30Days":    last30Days,
				"bandwidthCalendarMonth": calendarMonth,
				"bandwidthCurrentYear":   currentYear,
			},
		}

		cfg.Logger.Trace("System bandwidth stats requested", "remote_addr", r.RemoteAddr, "tz", loc.String())
		shared.WriteJSON(w, http.StatusOK, payload)
	}
}

func HealthHandler(cfg *config.BackendConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			shared.WriteJSONError(w, http.StatusMethodNotAllowed, "method not allowed")
			return
		}

		health := buildGoHealthResponse()
		payload := map[string]any{
			"response": health,
		}

		if len(health.RuntimeMetrics) > 0 {
			metric := health.RuntimeMetrics[0]
			cfg.Logger.Debug(
				"System Go runtime health requested",
				"remote_addr", r.RemoteAddr,
				"pid", metric.PID,
				"rss_bytes", metric.Memory.RSSBytes,
				"heap_alloc_bytes", metric.Memory.HeapAllocBytes,
				"goroutines", metric.Scheduler.Goroutines,
				"cpu_percent", metric.CPU.ProcessPercent,
			)
		} else {
			cfg.Logger.Debug("System Go runtime health requested", "remote_addr", r.RemoteAddr)
		}

		shared.WriteJSON(w, http.StatusOK, payload)
	}
}

func NodesStatsHandler(manager *dbmanager.DatabaseManager, cfg *config.BackendConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			shared.WriteJSONError(w, http.StatusMethodNotAllowed, "method not allowed")
			return
		}

		type nodeDayStat struct {
			NodeName   string `json:"nodeName"`
			Date       string `json:"date"`
			TotalBytes string `json:"totalBytes"`
		}

		stats := make([]nodeDayStat, 0)
		err := manager.ExecuteHighPriority(func(db dbmanager.DBExecutor) error {
			rows, err := db.QueryContext(r.Context(), `
				SELECT
					n.name AS node_name,
					DATE_TRUNC('day', nu.created_at)::date AS date,
					COALESCE(SUM(nu.total_bytes), 0)::text AS total_bytes
				FROM nodes_usage_history nu
				JOIN nodes n ON nu.node_uuid = n.uuid
				WHERE nu.created_at >= NOW() - INTERVAL '7 days'
				GROUP BY n.name, DATE_TRUNC('day', nu.created_at)::date
				ORDER BY date ASC
			`)
			if err != nil {
				return err
			}
			defer rows.Close()

			for rows.Next() {
				var (
					name       string
					date       time.Time
					totalBytes string
				)
				if err := rows.Scan(&name, &date, &totalBytes); err != nil {
					return err
				}
				stats = append(stats, nodeDayStat{
					NodeName:   name,
					Date:       date.UTC().Format(time.RFC3339),
					TotalBytes: totalBytes,
				})
			}
			return rows.Err()
		})
		if err != nil {
			shared.SendError(w, http.StatusInternalServerError, "failed to read nodes statistics", err, cfg)
			return
		}

		shared.WriteJSON(w, http.StatusOK, map[string]any{
			"response": map[string]any{
				"lastSevenDays": stats,
			},
		})
	}
}

func NodesMetricsHandler(manager *dbmanager.DatabaseManager, cfg *config.BackendConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			shared.WriteJSONError(w, http.StatusMethodNotAllowed, "method not allowed")
			return
		}

		nodes, err := loadNodesMetricsViaPrometheus(r.Context(), manager, cfg)
		if err != nil {
			cfg.Logger.Warn("Failed to load nodes metrics via prometheus endpoint", "error", err)
			nodes = []nodeMetricsItem{}
		}

		shared.WriteJSON(w, http.StatusOK, map[string]any{
			"response": nodesMetricsResponse{
				Nodes: nodes,
			},
		})
	}
}

type memStats struct {
	total     int64
	free      int64
	used      int64
	active    int64
	available int64
}

type onlineStats struct {
	onlineNow   int64
	lastDay     int64
	lastWeek    int64
	neverOnline int64
}

func readUsersStatusStats(ctx context.Context, manager *dbmanager.DatabaseManager) (map[string]int64, int64, error) {
	statusCounts := map[string]int64{
		"ACTIVE":   0,
		"DISABLED": 0,
		"LIMITED":  0,
		"EXPIRED":  0,
	}
	var total int64

	err := manager.ExecuteHighPriority(func(db dbmanager.DBExecutor) error {
		rows, err := db.QueryContext(ctx, `
			SELECT status, COUNT(*)
			FROM users
			WHERE status IN ('ACTIVE', 'DISABLED', 'LIMITED', 'EXPIRED')
			GROUP BY status
		`)
		if err != nil {
			return err
		}
		defer rows.Close()

		for rows.Next() {
			var status string
			var count int64
			if err := rows.Scan(&status, &count); err != nil {
				return err
			}
			statusCounts[status] = count
			total += count
		}
		return rows.Err()
	})

	return statusCounts, total, err
}

func readOnlineStats(ctx context.Context, manager *dbmanager.DatabaseManager) (onlineStats, error) {
	nowUTC := time.Now().UTC()
	thresholdOnline := nowUTC.Add(-30 * time.Second)
	thresholdDay := nowUTC.Add(-24 * time.Hour)
	thresholdWeek := nowUTC.Add(-7 * 24 * time.Hour)

	stats := onlineStats{}
	err := manager.ExecuteHighPriority(func(db dbmanager.DBExecutor) error {
		return db.QueryRowContext(ctx, `
			SELECT
				COUNT(t_id) FILTER (WHERE online_at >= ?) AS online_now,
				COUNT(t_id) FILTER (WHERE online_at >= ?) AS last_day,
				COUNT(t_id) FILTER (WHERE online_at >= ?) AS last_week,
				COUNT(t_id) FILTER (WHERE online_at IS NULL) AS never_online
			FROM user_traffic
		`, thresholdOnline, thresholdDay, thresholdWeek).Scan(
			&stats.onlineNow,
			&stats.lastDay,
			&stats.lastWeek,
			&stats.neverOnline,
		)
	})

	return stats, err
}

func readTotalOnlineOnNodes(ctx context.Context, manager *dbmanager.DatabaseManager, cfg *config.BackendConfig) (int64, error) {
	uuids := make([]string, 0)
	err := manager.ExecuteHighPriority(func(db dbmanager.DBExecutor) error {
		rows, err := db.QueryContext(ctx, `
			SELECT uuid
			FROM nodes
			WHERE is_connected = TRUE
		`)
		if err != nil {
			return err
		}
		defer rows.Close()
		for rows.Next() {
			var uuid string
			if err := rows.Scan(&uuid); err != nil {
				return err
			}
			uuids = append(uuids, uuid)
		}
		return rows.Err()
	})
	if err != nil {
		return 0, err
	}
	cache, _ := nodehotcache.Default(cfg).GetMany(ctx, uuids)
	var total int64
	for _, uuid := range uuids {
		total += int64(cache[uuid].UsersOnline)
	}
	return total, err
}

func readLifetimeTrafficBytes(ctx context.Context, manager *dbmanager.DatabaseManager) (string, error) {
	var total string
	err := manager.ExecuteHighPriority(func(db dbmanager.DBExecutor) error {
		return db.QueryRowContext(ctx, `
			SELECT COALESCE(SUM(total_bytes), 0)::text
			FROM nodes_usage_history
		`).Scan(&total)
	})
	return total, err
}

func readUsersRecap(ctx context.Context, manager *dbmanager.DatabaseManager, startOfMonth time.Time) (usersRecap, error) {
	recap := usersRecap{}
	err := manager.ExecuteHighPriority(func(db dbmanager.DBExecutor) error {
		return db.QueryRowContext(ctx, `
			SELECT
				COUNT(*)::bigint AS total,
				COUNT(*) FILTER (WHERE created_at >= ?)::bigint AS new_users_this_month
			FROM users
		`, startOfMonth).Scan(&recap.total, &recap.newUsersThisMonth)
	})
	return recap, err
}

func readNodesRecap(ctx context.Context, manager *dbmanager.DatabaseManager, cfg *config.BackendConfig) (nodesRecap, error) {
	recap := nodesRecap{}
	uuids := make([]string, 0)
	err := manager.ExecuteHighPriority(func(db dbmanager.DBExecutor) error {
		if err := db.QueryRowContext(ctx, `
			SELECT
				COUNT(*)::bigint AS total,
				COUNT(DISTINCT CASE
					WHEN country_code IS NOT NULL AND country_code <> '' AND country_code <> 'XX'
					THEN country_code
					ELSE NULL
				END)::bigint AS distinct_countries
			FROM nodes
		`).Scan(&recap.total, &recap.distinctCountries); err != nil {
			return err
		}

		rows, err := db.QueryContext(ctx, `SELECT uuid FROM nodes`)
		if err != nil {
			return err
		}
		defer rows.Close()
		for rows.Next() {
			var uuid string
			if err := rows.Scan(&uuid); err != nil {
				return err
			}
			uuids = append(uuids, uuid)
		}
		return rows.Err()
	})
	if err != nil {
		return recap, err
	}
	cache, _ := nodehotcache.Default(cfg).GetMany(ctx, uuids)
	var totalRAM int64
	var totalCPUCores int64
	for _, uuid := range uuids {
		system := cache[uuid].System
		if system == nil || len(system.Info) == 0 {
			continue
		}
		var info struct {
			MemoryTotal int64 `json:"memoryTotal"`
			CPUs        int64 `json:"cpus"`
		}
		if err := json.Unmarshal(system.Info, &info); err != nil {
			continue
		}
		totalRAM += info.MemoryTotal
		totalCPUCores += info.CPUs
	}
	recap.totalRam = fmt.Sprintf("%d", totalRAM)
	recap.totalCpuCores = totalCPUCores
	return recap, err
}

func readUsageBytesTextByRange(ctx context.Context, manager *dbmanager.DatabaseManager, dtRange usageRange) (string, error) {
	var total string
	err := manager.ExecuteHighPriority(func(db dbmanager.DBExecutor) error {
		return db.QueryRowContext(ctx, `
			SELECT COALESCE(SUM(total_bytes), 0)::text
			FROM nodes_usage_history
			WHERE created_at >= ? AND created_at <= ?
		`, dtRange.start, dtRange.end).Scan(&total)
	})
	return total, err
}

func readInitDate(ctx context.Context, manager *dbmanager.DatabaseManager) (time.Time, error) {
	var initDate time.Time
	err := manager.ExecuteHighPriority(func(db dbmanager.DBExecutor) error {
		return db.QueryRowContext(ctx, `
			SELECT COALESCE(
				(SELECT started_at FROM schema_migrations ORDER BY started_at ASC LIMIT 1),
				NOW()
			)
		`).Scan(&initDate)
	})
	if err != nil {
		return time.Now().UTC(), err
	}
	if initDate.IsZero() {
		return time.Now().UTC(), nil
	}
	return initDate, nil
}

func readUsageComparison(ctx context.Context, manager *dbmanager.DatabaseManager, ranges [2]usageRange) (bandwidthStat, error) {
	previousBytes, err := readUsageBytesByRange(ctx, manager, ranges[0])
	if err != nil {
		return bandwidthStat{}, err
	}
	currentBytes, err := readUsageBytesByRange(ctx, manager, ranges[1])
	if err != nil {
		return bandwidthStat{}, err
	}

	difference := new(big.Int).Sub(currentBytes, previousBytes)

	return bandwidthStat{
		Current:    formatBigBytes(currentBytes),
		Previous:   formatBigBytes(previousBytes),
		Difference: formatBigBytes(difference),
	}, nil
}

func readUsageBytesByRange(ctx context.Context, manager *dbmanager.DatabaseManager, dtRange usageRange) (*big.Int, error) {
	var totalText string
	err := manager.ExecuteHighPriority(func(db dbmanager.DBExecutor) error {
		return db.QueryRowContext(ctx, `
			SELECT COALESCE(SUM(total_bytes), 0)::text
			FROM nodes_usage_history
			WHERE created_at >= ? AND created_at <= ?
		`, dtRange.start, dtRange.end).Scan(&totalText)
	})
	if err != nil {
		return nil, err
	}

	result, ok := new(big.Int).SetString(strings.TrimSpace(totalText), 10)
	if !ok {
		return nil, fmt.Errorf("invalid bigint value: %s", totalText)
	}
	return result, nil
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

	// "used" is reported as the actively-used working set (Active from
	// /proc/meminfo) rather than the naive total-minus-free calculation.
	// total-free counts reclaimable page cache/buffers as "used", which
	// makes RAM usage look almost full on any Linux host even when the
	// system is barely loaded. Active reflects memory the kernel considers
	// recently/actually in use by processes, which is a much more honest
	// number to show on the dashboard.
	result.used = result.active
	if result.used <= 0 {
		// Fallback for kernels/environments without an "Active" line.
		result.used = result.total - result.free
	}
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
