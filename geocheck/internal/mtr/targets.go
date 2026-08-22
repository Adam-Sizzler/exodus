package mtr

import "strings"

// Target is one destination whose path we measure.
type Target struct {
	ID   string
	Name string // display label
	Host string // hostname or literal address
	Net  string // owning network, for grouping

	// ASN is the autonomous system the destination is expected to live in.
	// A path that ends anywhere else means DNS or the route was intercepted.
	ASN int

	// Anycast marks destinations announced from many sites at once. Their RTT
	// says how far the nearest announcement is, not where "the server" is, so
	// a geographic sanity check against the address is meaningless.
	Anycast bool

	// Port is used for TCP probing when ICMP is filtered.
	Port int

	// ICMPSilent marks hosts that never answer echo requests; they are probed
	// over TCP from the start.
	ICMPSilent bool

	// Tags group targets into the sets selectable with --targets.
	Tags []string
}

// HasTag reports whether the target carries a tag.
func (t Target) HasTag(tag string) bool {
	for _, x := range t.Tags {
		if strings.EqualFold(x, tag) {
			return true
		}
	}
	return false
}

// Catalog is the default destination set: the networks that carry the bulk of
// consumer traffic, plus the public resolvers that make good anycast probes.
//
// ICMPSilent is set from measurement, not from the brand. Responsiveness is a
// property of the individual address: Meta answers on 31.13.72.36 but never on
// 31.13.72.174 in the same /24, and several "Microsoft" or "Apple" hostnames
// actually resolve into Akamai. Where a company has both a silent front door
// and a responsive one, the responsive host is the target.
var Catalog = []Target{
	// --- Google ---
	{
		ID: "google", Name: "google.com", Host: "google.com", Net: "Google",
		ASN: 15169, Port: 443, Tags: []string{"default", "web", "google"},
	},
	{
		ID: "youtube", Name: "youtube.com", Host: "youtube.com", Net: "Google",
		ASN: 15169, Port: 443, Tags: []string{"default", "video", "google"},
	},
	{
		ID: "google_dns", Name: "dns.google (8.8.8.8)", Host: "8.8.8.8", Net: "Google",
		ASN: 15169, Anycast: true, Port: 443, Tags: []string{"default", "dns", "google"},
	},

	// --- Cloudflare ---
	{
		ID: "cloudflare", Name: "cloudflare.com", Host: "cloudflare.com", Net: "Cloudflare",
		ASN: 13335, Anycast: true, Port: 443, Tags: []string{"default", "web", "cdn"},
	},
	{
		ID: "cloudflare_dns", Name: "one.one.one.one (1.1.1.1)", Host: "1.1.1.1", Net: "Cloudflare",
		ASN: 13335, Anycast: true, Port: 443, Tags: []string{"default", "dns"},
	},

	// --- Streaming and social ---
	{
		// The apex resolves to an AWS load balancer that drops ICMP; the www
		// host is Netflix's own front-end tier and answers.
		ID: "netflix", Name: "netflix.com", Host: "www.netflix.com", Net: "Netflix",
		ASN: 40027, Port: 443, Tags: []string{"default", "video"},
	},
	{
		ID: "facebook", Name: "facebook.com", Host: "facebook.com", Net: "Meta",
		ASN: 32934, Port: 443, Tags: []string{"default", "social"},
	},
	{
		// Shares AS32934 with facebook.com, and both enter Meta at the same
		// border, so the verdict usually matches. It earns its place anyway:
		// Instagram is blocked and throttled independently of Facebook in
		// several countries, and the addresses it serves from are a different
		// set, so it can fail while facebook.com stays clean.
		//
		// Every address its DNS rotated through answered ICMP when measured,
		// unlike the Meta addresses that motivated the note above.
		ID: "instagram", Name: "instagram.com", Host: "www.instagram.com", Net: "Meta",
		ASN: 32934, Port: 443, Tags: []string{"default", "social"},
	},
	{
		ID: "whatsapp", Name: "whatsapp.net", Host: "www.whatsapp.com", Net: "Meta",
		ASN: 32934, Port: 443, ICMPSilent: true, Tags: []string{"social", "messaging"},
	},
	// Telegram's five data centres. Which one serves you is fixed by where the
	// account was registered, not by where you are, so a user in Europe can
	// legitimately be pinned to DC1.
	//
	// They occupy three networks, not five: DC1 and DC3 share AS59930, DC2 and
	// DC4 share AS62041, and DC5 is alone on AS62014.
	//
	// The default set traces DC2 and DC5 only. Note what that means: AS59930 is
	// not measured at all by a plain run, so a problem confined to DC1 and DC3
	// will not show up. That is a deliberate trade for a shorter default, not
	// an oversight — `-T telegram` traces all five.
	//
	// The apex telegram.org drops ICMP, so the DC addresses are used directly.
	{
		ID: "telegram_dc1", Name: "Telegram DC1", Host: "149.154.175.50", Net: "Telegram",
		ASN: 59930, Port: 443, Tags: []string{"telegram", "messaging"},
	},
	{
		ID: "telegram", Name: "Telegram DC2", Host: "149.154.167.51", Net: "Telegram",
		ASN: 62041, Port: 443, Tags: []string{"default", "telegram", "messaging"},
	},
	{
		ID: "telegram_dc3", Name: "Telegram DC3", Host: "149.154.175.100", Net: "Telegram",
		ASN: 59930, Port: 443, Tags: []string{"telegram", "messaging"},
	},
	{
		ID: "telegram_dc4", Name: "Telegram DC4", Host: "149.154.167.91", Net: "Telegram",
		ASN: 62041, Port: 443, Tags: []string{"telegram", "messaging"},
	},
	{
		ID: "telegram_dc5", Name: "Telegram DC5", Host: "91.108.56.130", Net: "Telegram",
		ASN: 62014, Port: 443, Tags: []string{"default", "telegram", "messaging"},
	},
	{
		// The gateway answers ICMP consistently; discord.com itself is flaky.
		ID: "discord", Name: "gateway.discord.gg", Host: "gateway.discord.gg",
		Net: "Cloudflare / Discord", ASN: 13335, Anycast: true, Port: 443,
		Tags: []string{"default", "messaging"},
	},
	{
		ID: "tiktok", Name: "tiktok.com", Host: "www.tiktok.com", Net: "Akamai (TikTok)",
		ASN: 20940, Port: 443, Tags: []string{"default", "social", "video"},
	},
	{
		ID: "twitch", Name: "twitch.tv", Host: "www.twitch.tv", Net: "Amazon",
		ASN: 16509, Port: 443, ICMPSilent: true, Tags: []string{"video"},
	},

	// --- Clouds and CDNs ---
	{
		// A CloudFront edge, which answers ICMP; the regional ec2.* endpoints
		// do not, in any region.
		ID: "aws", Name: "awsstatic.com (CloudFront)", Host: "d1.awsstatic.com", Net: "Amazon",
		ASN: 16509, Port: 443, Tags: []string{"default", "cloud", "cdn"},
	},
	{
		// Exchange Online answers ICMP; the microsoft.com anycast edge does not.
		ID: "microsoft", Name: "outlook.office365.com", Host: "outlook.office365.com",
		Net: "Microsoft", ASN: 8075, Port: 443, Tags: []string{"cloud"},
	},
	{
		ID: "akamai", Name: "akamai.com", Host: "www.akamai.com", Net: "Akamai",
		ASN: 20940, Port: 443, Tags: []string{"default", "cdn"},
	},
	{
		ID: "fastly", Name: "fastly.com", Host: "www.fastly.com", Net: "Fastly",
		ASN: 54113, Anycast: true, Port: 443, Tags: []string{"default", "cdn"},
	},
	{
		// The apex is Apple's own AS714; www.apple.com is served by Akamai.
		ID: "apple", Name: "apple.com", Host: "apple.com", Net: "Apple",
		ASN: 714, Port: 443, Tags: []string{"default", "web"},
	},

	// --- AI and dev ---
	{
		ID: "openai", Name: "api.openai.com", Host: "api.openai.com", Net: "Cloudflare / OpenAI",
		ASN: 13335, Anycast: true, Port: 443, Tags: []string{"default", "ai"},
	},
	{
		ID: "anthropic", Name: "claude.ai", Host: "claude.ai", Net: "Anthropic",
		ASN: 399358, Anycast: true, Port: 443, Tags: []string{"default", "ai"},
	},
	{
		ID: "github", Name: "github.com", Host: "github.com", Net: "GitHub / Azure",
		ASN: 8075, Port: 443, Tags: []string{"default", "dev"},
	},
	{
		// Valve's own content network. The steampowered.com front ends are
		// Akamai and Fastly, so they would measure those instead.
		ID: "steam", Name: "steamcontent.com", Host: "cache1-fra1.steamcontent.com",
		Net: "Valve", ASN: 32590, Port: 443, Tags: []string{"default", "gaming"},
	},

	// --- Anycast resolvers: cheap, reliable, honest latency references ---
	{
		ID: "quad9", Name: "Quad9 (9.9.9.9)", Host: "9.9.9.9", Net: "Quad9",
		ASN: 19281, Anycast: true, Port: 443, Tags: []string{"default", "dns"},
	},
	{
		ID: "opendns", Name: "OpenDNS (208.67.222.222)", Host: "208.67.222.222", Net: "Cisco OpenDNS",
		ASN: 36692, Anycast: true, Port: 443, Tags: []string{"dns"},
	},
}

// DefaultTargets returns the targets probed when no selection is given.
func DefaultTargets() []Target { return Select("default") }

// Select filters the catalog by tag. The tag "all" returns everything.
func Select(tags ...string) []Target {
	var out []Target
	for _, t := range Catalog {
		for _, tag := range tags {
			if strings.EqualFold(tag, "all") || t.HasTag(tag) || strings.EqualFold(tag, t.ID) {
				out = append(out, t)
				break
			}
		}
	}
	return out
}

// Tags lists every tag in the catalog, for help output.
func Tags() []string {
	seen := map[string]bool{}
	var out []string
	for _, t := range Catalog {
		for _, tag := range t.Tags {
			if !seen[tag] {
				seen[tag] = true
				out = append(out, tag)
			}
		}
	}
	return out
}
