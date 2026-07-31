package system

import (
	"bytes"
	"database/sql"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
	"github.com/prometheus/common/expfmt"
)

// newRuntimeMetricsRegistry builds a dedicated Prometheus registry for DB
// connection-pool stats and Go/process runtime metrics, using the official
// client library instead of hand-built text (unlike the legacy per-node
// metrics below, which stay as-is for now).
//
// Collector registration must happen exactly once per process — Register
// panics on a duplicate collector — so this must be called once at handler
// construction time (inside MetricsHandler, which itself runs once at
// startup), never per-request.
func newRuntimeMetricsRegistry(interactive, background *sql.DB) *prometheus.Registry {
	registry := prometheus.NewRegistry()

	if interactive != nil {
		registry.MustRegister(collectors.NewDBStatsCollector(interactive, "interactive"))
	}
	if background != nil {
		registry.MustRegister(collectors.NewDBStatsCollector(background, "background"))
	}
	registry.MustRegister(collectors.NewGoCollector())
	registry.MustRegister(collectors.NewProcessCollector(collectors.ProcessCollectorOpts{}))

	return registry
}

// renderRegistry serializes a registry to the standard Prometheus text
// exposition format, so it can be concatenated with the legacy hand-built
// node-metrics payload in renderPrometheusMetricsCached.
func renderRegistry(registry *prometheus.Registry) (string, error) {
	families, err := registry.Gather()
	if err != nil {
		return "", err
	}

	var buf bytes.Buffer
	encoder := expfmt.NewEncoder(&buf, expfmt.NewFormat(expfmt.TypeTextPlain))
	for _, family := range families {
		if err := encoder.Encode(family); err != nil {
			return "", err
		}
	}
	return buf.String(), nil
}
