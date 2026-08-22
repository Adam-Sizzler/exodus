package portal

import "strings"

// Catalog lists the connectivity-check endpoints, each with the exact answer
// its vendor specifies.
//
// The expectations here are assertions, so a wrong one produces a false
// "altered" verdict, which is worse than not checking at all. Every entry was
// fetched and its body counted in bytes; anything that could not be confirmed
// is listed at the bottom of this comment rather than guessed.
//
// Endpoints are chosen to span distinct networks as well as distinct vendors.
// Several of these names are not hosted by the company they belong to —
// www.msftconnecttest.com and conn1.oppomobile.com are Akamai, connect.rom.miui.com
// is Azure — so picking by brand would silently stack three checks onto one CDN
// that fails together. The default set below covers eight distinct autonomous
// systems.
//
// Watch the trailing newlines. Apple's two endpoints differ by exactly one byte
// (captive.apple.com sends 69 with \n, www.apple.com 68 without), and Debian's
// /nm sends 25 bytes with \n where GNOME's sends the same sentence as 24
// without. Comparison tolerates surrounding whitespace so these stay robust.
//
// Deliberately absent, having been checked and rejected:
//   - qualcomm.com and qualcomm.cn: the whole site answers 403 "error code: 1009"
//     outside permitted regions, so no path can be relied on.
//   - msftconnecttest.com and msftncsi.com bare apexes: no A and no AAAA record.
//   - connectivity.lineageos.org: no DNS records; lineageos.org/generate_204 is 404.
//   - www.msftconnecttest.com/redirect: a 302 by design, being the portal landing
//     page rather than a check.
//   - test.steampowered.com/generate_204: 404. The working paths are / and /204.
var Catalog = []Endpoint{
	// --- Default set: eight distinct autonomous systems ---
	{
		ID: "google", Name: "Google (Android)", Vendor: "Google",
		URL:        "http://connectivitycheck.gstatic.com/generate_204",
		WantStatus: 204,
		Tags:       []string{"default", "http"},
	},
	{
		ID: "cloudflare", Name: "Cloudflare", Vendor: "Cloudflare",
		URL:        "http://cp.cloudflare.com/generate_204",
		WantStatus: 204,
		Tags:       []string{"default", "http"},
	},
	{
		ID: "apple", Name: "Apple (iOS/macOS)", Vendor: "Apple",
		URL:        "http://captive.apple.com/hotspot-detect.html",
		WantStatus: 200,
		WantBody:   "<HTML><HEAD><TITLE>Success</TITLE></HEAD><BODY>Success</BODY></HTML>\n",
		Tags:       []string{"default", "http"},
	},
	{
		// Akamai, despite the name. The only endpoint on Microsoft's own
		// network is edge_microsoft below.
		ID: "microsoft", Name: "Microsoft NCSI", Vendor: "Microsoft",
		URL:        "http://www.msftconnecttest.com/connecttest.txt",
		WantStatus: 200,
		WantBody:   "Microsoft Connect Test",
		Tags:       []string{"default", "http"},
	},
	{
		ID: "firefox", Name: "Mozilla Firefox", Vendor: "Mozilla",
		URL:        "http://detectportal.firefox.com/success.txt",
		WantStatus: 200,
		WantBody:   "success\n",
		Tags:       []string{"default", "http"},
	},
	{
		ID: "ubuntu", Name: "Ubuntu", Vendor: "Canonical",
		URL:        "http://connectivity-check.ubuntu.com/",
		WantStatus: 204,
		Tags:       []string{"default", "http", "linux"},
	},
	{
		// Single-homed on Hetzner rather than a CDN, so it stays meaningful
		// when a CDN-wide anomaly takes the others out together.
		ID: "arch", Name: "Arch Linux", Vendor: "Arch Linux",
		URL:        "http://ping.archlinux.org/",
		WantStatus: 200,
		WantBody:   "This domain is used for connectivity checking (captive portal detection).\n",
		Tags:       []string{"default", "http", "linux"},
	},
	{
		ID: "edge_microsoft", Name: "Microsoft Edge", Vendor: "Microsoft",
		URL:        "http://edge-http.microsoft.com/captiveportal/generate_204",
		WantStatus: 204,
		Tags:       []string{"default", "http"},
	},
	{
		// HTTPS cannot be rewritten without breaking TLS, so a failure here
		// beside a working HTTP check points at TLS interception.
		ID: "google_tls", Name: "Google over HTTPS", Vendor: "Google",
		URL:        "https://connectivitycheck.gstatic.com/generate_204",
		WantStatus: 204,
		Tags:       []string{"default", "tls"},
	},

	// --- Additional vendors, selectable by tag ---
	{
		ID: "cloudflare_tls", Name: "Cloudflare over HTTPS", Vendor: "Cloudflare",
		URL:        "https://cp.cloudflare.com/generate_204",
		WantStatus: 204,
		Tags:       []string{"tls"},
	},
	{
		ID: "google_clients3", Name: "Google (clients3)", Vendor: "Google",
		URL:        "http://clients3.google.com/generate_204",
		WantStatus: 204,
		Tags:       []string{"http"},
	},
	{
		ID: "msftncsi", Name: "Microsoft NCSI (legacy)", Vendor: "Microsoft",
		URL:        "http://www.msftncsi.com/ncsi.txt",
		WantStatus: 200,
		WantBody:   "Microsoft NCSI",
		Tags:       []string{"http"},
	},
	{
		ID: "steam", Name: "Steam", Vendor: "Valve",
		URL:        "http://test.steampowered.com/204",
		WantStatus: 204,
		Tags:       []string{"http"},
	},
	{
		ID: "gnome", Name: "GNOME NetworkManager", Vendor: "GNOME",
		URL:        "http://nmcheck.gnome.org/check_network_status.txt",
		WantStatus: 200,
		WantBody:   "NetworkManager is online",
		Tags:       []string{"http", "linux"},
	},
	{
		ID: "debian", Name: "Debian", Vendor: "Debian",
		URL:        "http://network-test.debian.org/nm",
		WantStatus: 200,
		WantBody:   "NetworkManager is online\n",
		Tags:       []string{"http", "linux"},
	},
	{
		ID: "fedora", Name: "Fedora", Vendor: "Fedora",
		URL:        "http://fedoraproject.org/static/hotspot.txt",
		WantStatus: 200,
		WantBody:   "OK",
		Tags:       []string{"http", "linux"},
	},
	{
		ID: "kde", Name: "KDE", Vendor: "KDE",
		URL:        "http://networkcheck.kde.org/",
		WantStatus: 200,
		WantBody:   "OK",
		Tags:       []string{"http", "linux"},
	},
	{
		ID: "kindle", Name: "Amazon Kindle", Vendor: "Amazon",
		URL:        "http://spectrum.s3.amazonaws.com/kindle-wifi/wifistub.html",
		WantStatus: 200,
		// The embedded UUID is the stable part of the page.
		BodyContains: "81ce4465-7167-4dcb-835b-dcc9e44c112a",
		Tags:         []string{"http"},
	},
	{
		ID: "v2ex", Name: "V2EX", Vendor: "V2EX",
		// The plain-HTTP form is a 301 to HTTPS, so only the HTTPS one is usable.
		URL:        "https://v2ex.com/generate_204",
		WantStatus: 204,
		Tags:       []string{"tls"},
	},
	{
		ID: "kuketz", Name: "Kuketz", Vendor: "Kuketz",
		URL:        "http://www.kuketz.de/generate_204",
		WantStatus: 204,
		Tags:       []string{"http"},
	},

	// --- Handset vendors. Reachable worldwide unless noted. ---
	{
		// Hosted on Azure, so it answers quickly from anywhere.
		ID: "xiaomi", Name: "Xiaomi (MIUI)", Vendor: "Xiaomi",
		URL:        "http://connect.rom.miui.com/generate_204",
		WantStatus: 204,
		Tags:       []string{"http", "cn"},
	},
	{
		ID: "huawei", Name: "Huawei", Vendor: "Huawei",
		URL:        "http://connectivitycheck.platform.hicloud.com/generate_204",
		WantStatus: 204,
		Tags:       []string{"http", "cn"},
	},
	{
		ID: "honor", Name: "Honor", Vendor: "Honor",
		URL:        "http://connectivitycheck.platform.hihonorcloud.com/generate_204",
		WantStatus: 204,
		Tags:       []string{"http", "cn"},
	},
	{
		ID: "oppo", Name: "OPPO", Vendor: "OPPO",
		URL:        "http://conn1.oppomobile.com/generate_204",
		WantStatus: 204,
		Tags:       []string{"http", "cn"},
	},
	{
		ID: "vivo", Name: "Vivo", Vendor: "Vivo",
		URL:        "http://wifi.vivo.com.cn/generate_204",
		WantStatus: 204,
		Tags:       []string{"http", "cn"},
	},
	{
		// Genuinely inside China (China Telecom, Shanghai), so it is slow from
		// elsewhere. Useful for a China vantage point, not as a global default.
		ID: "samsung", Name: "Samsung", Vendor: "Samsung",
		URL:        "http://connectivity.samsung.com.cn/generate_204",
		WantStatus: 204,
		Tags:       []string{"http", "cn"},
	},
}

// httpsUnsafe lists endpoints whose HTTPS form fails for reasons that have
// nothing to do with interception: a shared CDN certificate that does not carry
// the hostname, or a handshake that does not complete on some anycast nodes.
// Checking them over TLS would report a legitimate quirk as a MITM, so the
// catalogue only ever fetches them over plain HTTP.
var httpsUnsafe = map[string]string{
	"microsoft":      "serves Akamai's generic a248.e.akamai.net certificate",
	"msftncsi":       "serves Akamai's generic a248.e.akamai.net certificate",
	"debian":         "serves Fastly's default sni-644 certificate",
	"edge_microsoft": "answers 400 to HTTPS while HTTP returns the 204",
	"firefox":        "some Fastly anycast nodes never complete the handshake",
}

// HTTPSUnsafe reports whether an endpoint must not be fetched over TLS, and why.
func HTTPSUnsafe(id string) (string, bool) {
	reason, ok := httpsUnsafe[id]
	return reason, ok
}

// DefaultEndpoints returns the set probed when no selection is given: broad
// vendor coverage across distinct networks, without probing six URLs that share
// one anycast path.
func DefaultEndpoints() []Endpoint { return Select("default") }

// Select filters the catalogue by tag or id. The tag "all" returns everything.
func Select(tags ...string) []Endpoint {
	var out []Endpoint
	for _, ep := range Catalog {
		for _, tag := range tags {
			if strings.EqualFold(tag, "all") || ep.ID == tag || hasTag(ep, tag) {
				out = append(out, ep)
				break
			}
		}
	}
	return out
}

func hasTag(ep Endpoint, tag string) bool {
	for _, t := range ep.Tags {
		if strings.EqualFold(t, tag) {
			return true
		}
	}
	return false
}

// Tags lists every tag in the catalogue, for help output.
func Tags() []string {
	seen := map[string]bool{}
	var out []string
	for _, ep := range Catalog {
		for _, t := range ep.Tags {
			if !seen[t] {
				seen[t] = true
				out = append(out, t)
			}
		}
	}
	return out
}
