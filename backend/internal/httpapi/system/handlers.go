package system

import (
	"database/sql"
	"net/http"
	"runtime"
	"strings"
	"time"

	"exodus/internal/config"
	"exodus/internal/constant"
	"exodus/internal/httpapi/shared"
)

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

		frontendSHA := firstMetadataEnv("__EX_METADATA_GIT_FRONTEND_COMMIT")
		if frontendSHA == "" {
			frontendSHA = backendSHA
		}
		buildNumber := firstMetadataEnv("__EX_METADATA_BUILD_NUMBER", "BUILD_NUMBER", "GITHUB_RUN_NUMBER")
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

func StatsHandler(db *sql.DB, cfg *config.BackendConfig) http.HandlerFunc {
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

		statusCounts, totalUsers, err := readUsersStatusStats(r.Context(), db)
		if err != nil {
			shared.SendError(w, http.StatusInternalServerError, "failed to read users stats", err, cfg)
			return
		}

		onlineStats, err := readOnlineStats(r.Context(), db)
		if err != nil {
			shared.SendError(w, http.StatusInternalServerError, "failed to read online stats", err, cfg)
			return
		}

		totalOnline, err := readTotalOnlineOnNodes(r.Context(), db, cfg)
		if err != nil {
			shared.SendError(w, http.StatusInternalServerError, "failed to read nodes online stats", err, cfg)
			return
		}

		lifetimeBytes, err := readLifetimeTrafficBytes(r.Context(), db)
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

func RecapHandler(db *sql.DB, cfg *config.BackendConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			shared.WriteJSONError(w, http.StatusMethodNotAllowed, "method not allowed")
			return
		}

		now := time.Now().UTC()
		startOfMonth := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, time.UTC)
		endOfMonth := startOfMonth.AddDate(0, 1, 0).Add(-time.Nanosecond)

		users, err := readUsersRecap(r.Context(), db, startOfMonth)
		if err != nil {
			shared.SendError(w, http.StatusInternalServerError, "failed to read users recap", err, cfg)
			return
		}

		nodes, err := readNodesRecap(r.Context(), db, cfg)
		if err != nil {
			shared.SendError(w, http.StatusInternalServerError, "failed to read nodes recap", err, cfg)
			return
		}

		lifetimeTraffic, err := readLifetimeTrafficBytes(r.Context(), db)
		if err != nil {
			shared.SendError(w, http.StatusInternalServerError, "failed to read lifetime traffic recap", err, cfg)
			return
		}

		monthTraffic, err := readUsageBytesTextByRange(r.Context(), db, usageRange{start: startOfMonth, end: endOfMonth})
		if err != nil {
			shared.SendError(w, http.StatusInternalServerError, "failed to read month traffic recap", err, cfg)
			return
		}

		initDate, err := readInitDate(r.Context(), db)
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

func BandwidthStatsHandler(db *sql.DB, cfg *config.BackendConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			shared.WriteJSONError(w, http.StatusMethodNotAllowed, "method not allowed")
			return
		}

		tz := strings.TrimSpace(r.URL.Query().Get("tz"))
		loc := resolveLocation(tz)
		now := time.Now().In(loc)

		lastTwoDays, err := readUsageComparison(r.Context(), db, getLastTwoDaysRanges(now))
		if err != nil {
			shared.SendError(w, http.StatusInternalServerError, "failed to read last two days bandwidth", err, cfg)
			return
		}
		lastSevenDays, err := readUsageComparison(r.Context(), db, getLastSevenDaysRanges(now))
		if err != nil {
			shared.SendError(w, http.StatusInternalServerError, "failed to read last seven days bandwidth", err, cfg)
			return
		}
		last30Days, err := readUsageComparison(r.Context(), db, getLast30DaysRanges(now))
		if err != nil {
			shared.SendError(w, http.StatusInternalServerError, "failed to read last 30 days bandwidth", err, cfg)
			return
		}
		calendarMonth, err := readUsageComparison(r.Context(), db, getCalendarMonthRanges(now))
		if err != nil {
			shared.SendError(w, http.StatusInternalServerError, "failed to read calendar month bandwidth", err, cfg)
			return
		}
		currentYear, err := readUsageComparison(r.Context(), db, getCalendarYearRanges(now))
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

func NodesStatsHandler(db *sql.DB, cfg *config.BackendConfig) http.HandlerFunc {
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
			shared.SendError(w, http.StatusInternalServerError, "failed to read nodes statistics", err, cfg)
			return
		}
		defer rows.Close()

		for rows.Next() {
			var (
				name       string
				date       time.Time
				totalBytes string
			)
			if scanErr := rows.Scan(&name, &date, &totalBytes); scanErr != nil {
				shared.SendError(w, http.StatusInternalServerError, "failed to read nodes statistics", scanErr, cfg)
				return
			}
			stats = append(stats, nodeDayStat{
				NodeName:   name,
				Date:       date.UTC().Format(time.RFC3339),
				TotalBytes: totalBytes,
			})
		}
		if err := rows.Err(); err != nil {
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

func NodesMetricsHandler(db *sql.DB, cfg *config.BackendConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			shared.WriteJSONError(w, http.StatusMethodNotAllowed, "method not allowed")
			return
		}

		nodes, err := loadNodesMetricsViaPrometheus(r.Context(), db, cfg)
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
