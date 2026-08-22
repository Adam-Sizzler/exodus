package geo

import (
	"context"
	"net/http"
	"net/url"
	"strings"

	"exodus-node/geocheck/internal/jsonx"
	"exodus-node/geocheck/internal/netx"
)

// apiProvider is a GeoIP endpoint that answers with the country of an IP.
type apiProvider struct {
	id      string
	name    string
	url     string // {ip} is substituted with the public address
	path    string // dot path to the country code in the JSON reply
	headers map[string]string

	// v6OverV4 asks about the IPv6 address over an IPv4 transport, for
	// providers whose API endpoint has no AAAA record.
	v6OverV4 bool
}

// geoIPProviders mirrors the endpoint set ipregion probes, expressed
// declaratively. Keys embedded in URLs are the providers' own public demo
// keys, exactly as their websites use them.
var geoIPProviders = []apiProvider{
	{
		id: "maxmind", name: "maxmind.com",
		url:     "https://geoip.maxmind.com/geoip/v2.1/city/me",
		path:    "country.iso_code",
		headers: map[string]string{"Referer": "https://www.maxmind.com"},
	},
	{
		id: "ripe", name: "rdap.db.ripe.net",
		url:  "https://rdap.db.ripe.net/ip/{ip}",
		path: "country",
	},
	{
		id: "ipinfo", name: "ipinfo.io",
		url: "https://ipinfo.io/{ip}/json", path: "country",
		v6OverV4: true,
	},
	{
		id: "ipregistry", name: "ipregistry.co",
		url:     "https://api.ipregistry.co/{ip}?hostname=true&key=sb69ksjcajfs4c",
		path:    "location.country.code",
		headers: map[string]string{"Origin": "https://ipregistry.co"},
	},
	{
		id: "ipapi_co", name: "ipapi.co",
		url: "https://ipapi.co/{ip}/json", path: "country",
		headers:  map[string]string{"Referer": "https://ipapi.co/"},
		v6OverV4: true,
	},
	{
		id: "country_is", name: "country.is",
		url: "https://api.country.is/{ip}", path: "country",
	},
	{
		id: "geoapify", name: "geoapify.com",
		url:  "https://api.geoapify.com/v1/ipinfo?&ip={ip}&apiKey=b8568cb9afc64fad861a69edbddb2658",
		path: "country.iso_code",
	},
	{
		id: "geojs", name: "geojs.io",
		url: "https://get.geojs.io/v1/ip/country.json?ip={ip}", path: "0.country",
	},
	{
		id: "ipapi_is", name: "ipapi.is",
		url: "https://api.ipapi.is/?q={ip}", path: "cc",
		v6OverV4: true,
	},
	{
		id: "ipbase", name: "ipbase.com",
		url: "https://api.ipbase.com/v2/info?ip={ip}", path: "data.location.country.alpha2",
	},
	{
		id: "ipquery", name: "ipquery.io",
		url: "https://api.ipquery.io/{ip}", path: "location.country_code",
	},
	{
		id: "ipwho", name: "ipwho.is",
		url: "https://ipwho.is/{ip}", path: "country_code",
	},
	{
		id: "ip_api_com", name: "ip-api.com",
		url:     "https://demo.ip-api.com/json/{ip}?fields=countryCode",
		path:    "countryCode",
		headers: map[string]string{"Origin": "https://ip-api.com"},
	},
	{
		id: "2ip", name: "2ip.io",
		url: "https://api.2ip.io", path: "code",
	},
}

func (p apiProvider) check() Check {
	return Check{
		ID: p.id, Name: p.name, Group: GroupGeoIP, Kind: KindCountry,
		Run: func(ctx context.Context, env *Env, f netx.Family) (string, error) {
			ip := env.PublicIP(f)
			transport := f
			if f == netx.V6 && p.v6OverV4 {
				transport = netx.V4
			}
			resp, err := env.Stack.Do(ctx, transport, netx.Request{
				URL:     strings.ReplaceAll(p.url, "{ip}", ip.String()),
				Headers: p.headers,
			})
			if err != nil {
				return "", err
			}
			if rateLimited(resp) {
				return "", nil
			}
			return jsonx.String(resp.Body, p.path), nil
		},
	}
}

// cloudflareTrace reads the edge PoP's view of you from the /cdn-cgi/trace
// key=value document.
func cloudflareTrace() Check {
	return Check{
		ID: "cloudflare", Name: "cloudflare.com", Group: GroupGeoIP, Kind: KindCountry,
		Run: func(ctx context.Context, env *Env, f netx.Family) (string, error) {
			resp, err := env.Stack.Do(ctx, f, netx.Request{
				URL: "https://www.cloudflare.com/cdn-cgi/trace",
			})
			if err != nil {
				return "", err
			}
			for _, line := range strings.Split(resp.Text(), "\n") {
				if k, v, ok := strings.Cut(line, "="); ok && k == "loc" {
					return v, nil
				}
			}
			return "", nil
		},
	}
}

// ipLocationCom posts the address to the site's own form endpoint.
func ipLocationCom() Check {
	return Check{
		ID: "iplocation", Name: "iplocation.com", Group: GroupGeoIP, Kind: KindCountry,
		Run: func(ctx context.Context, env *Env, f netx.Family) (string, error) {
			resp, err := env.Stack.Do(ctx, f, netx.Request{
				Method: http.MethodPost,
				URL:    "https://iplocation.com",
				Form:   url.Values{"ip": {env.PublicIP(f).String()}},
			})
			if err != nil {
				return "", err
			}
			return jsonx.String(resp.Body, "country_code"), nil
		},
	}
}

// DatabaseChecks returns the GeoIP database/API probes.
func DatabaseChecks() []Check {
	checks := make([]Check, 0, 2+len(geoIPProviders))
	checks = append(checks, cloudflareTrace(), ipLocationCom())
	for _, p := range geoIPProviders {
		checks = append(checks, p.check())
	}
	return checks
}

// rateLimited reports whether the provider pushed us away; those replies carry
// no useful country and should read as "no answer" rather than an error.
func rateLimited(r *netx.Response) bool {
	return r.Status == http.StatusForbidden || r.Status == http.StatusTooManyRequests
}
