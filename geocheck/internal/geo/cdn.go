package geo

import (
	"context"
	"net/http"
	"net/url"
	"strings"

	"exodus-node/geocheck/internal/jsonx"
	"exodus-node/geocheck/internal/netx"
)

// CDNChecks reports which edge location actually serves you. This is the most
// direct evidence of where your traffic physically lands: a Frankfurt exit that
// is served from a Moscow PoP is not really in Frankfurt.
func CDNChecks() []Check {
	return []Check{cloudflareColo(), youtubeCDN(), netflixCDN()}
}

// cloudflareColo reads the `colo=` field of /cdn-cgi/trace, the IATA code of
// the Cloudflare PoP terminating the connection.
func cloudflareColo() Check {
	return Check{
		ID: "cdn_cloudflare", Name: "Cloudflare edge", Group: GroupCDN, Kind: KindCountry,
		Run: func(ctx context.Context, env *Env, f netx.Family) (string, error) {
			resp, err := env.Stack.Do(ctx, f, netx.Request{
				URL: "https://www.cloudflare.com/cdn-cgi/trace",
			})
			if err != nil {
				return "", err
			}
			for _, line := range strings.Split(resp.Text(), "\n") {
				if k, v, ok := strings.Cut(line, "="); ok && k == "colo" {
					return iataCountry(ctx, env, f, v), nil
				}
			}
			return "", nil
		},
	}
}

// youtubeCDN asks the googlevideo redirector which Google edge cluster your
// prefix maps to; the cluster name starts with the IATA code of its city.
func youtubeCDN() Check {
	return Check{
		ID: "cdn_youtube", Name: "YouTube / GGC edge", Group: GroupCDN, Kind: KindCountry,
		Run: func(ctx context.Context, env *Env, f netx.Family) (string, error) {
			resp, err := env.Stack.Do(ctx, f, netx.Request{
				URL: "https://redirector.googlevideo.com/report_mapping?di=no",
			})
			if err != nil {
				return "", err
			}
			code := parseGGCCluster(resp.Text())
			if code == "" {
				return "", nil
			}
			return iataCountry(ctx, env, f, code), nil
		},
	}
}

// parseGGCCluster pulls the IATA prefix out of a "prefix => cluster" mapping
// line, e.g. "192.0.2.0/24 => fra16s52" -> "FRA".
func parseGGCCluster(body string) string {
	for _, line := range strings.Split(body, "\n") {
		fields := strings.Fields(line)
		if len(fields) < 3 {
			continue
		}
		cluster := fields[2]
		if i := strings.Index(cluster, "-"); i >= 0 {
			cluster = cluster[i+1:]
		}
		if len(cluster) < 3 {
			continue
		}
		code := strings.ToUpper(cluster[:3])
		if isAlpha3(code) {
			return code
		}
	}
	return ""
}

func isAlpha3(s string) bool {
	if len(s) != 3 {
		return false
	}
	for i := 0; i < 3; i++ {
		if s[i] < 'A' || s[i] > 'Z' {
			return false
		}
	}
	return true
}

// netflixCDN reports the country of the Open Connect appliance that would
// serve a stream.
func netflixCDN() Check {
	return Check{
		ID: "cdn_netflix", Name: "Netflix Open Connect", Group: GroupCDN, Kind: KindCountry,
		Run: func(ctx context.Context, env *Env, f netx.Family) (string, error) {
			resp, err := env.Stack.Do(ctx, f, netx.Request{URL: fastAPIURL})
			if err != nil {
				return "", err
			}
			return jsonx.String(resp.Body, "targets.0.location.country"), nil
		},
	}
}

// iataCountry maps an airport code to a country, preferring the embedded table
// and falling back to a public lookup for codes we do not carry.
func iataCountry(ctx context.Context, env *Env, f netx.Family, code string) string {
	code = strings.ToUpper(strings.TrimSpace(code))
	if code == "" {
		return ""
	}
	if c, ok := iataToCountry[code]; ok {
		return c
	}
	resp, err := env.Stack.Do(ctx, f, netx.Request{
		Method: http.MethodPost,
		URL:    "https://www.air-port-codes.com/api/v1/single",
		Headers: map[string]string{
			"APC-Auth": "96dc04b3fb",
			"Referer":  "https://www.air-port-codes.com/",
		},
		Form: url.Values{"iata": {code}},
	})
	if err != nil {
		return ""
	}
	return jsonx.String(resp.Body, "airport.country.iso")
}
