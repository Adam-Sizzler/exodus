package access

import (
	"context"
	"net/url"
	"regexp"
	"strings"

	"exodus-node/geocheck/internal/netx"
)

const browserUA = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 " +
	"(KHTML, like Gecko) Chrome/151.0.0.0 Safari/537.36"

// Checks returns the service-availability probes, in display order.
func Checks() []Check {
	return []Check{
		chatGPTWeb(),
		chatGPTApp(),
		youTubePremium(),
		netflix(),
		claude(),
		tiktok(),
		gemini(),
		notebookLM(),
	}
}

// youTubePremium confirms availability positively rather than by absence.
// The refusal sentence names the country; the "ad-free" wording only appears on
// the real offer page, so requiring it means a redirect to some other page
// cannot be mistaken for success.
func youTubePremium() Check {
	return Check{
		ID: "youtube_premium_access", Name: "YouTube Premium",
		Run: func(ctx context.Context, env Env) Result {
			resp, err := env.Stack.Do(ctx, env.Family, netx.Request{
				URL:       "https://www.youtube.com/premium",
				UserAgent: browserUA,
				Headers:   map[string]string{"Accept-Language": "en-US,en;q=0.9"},
			})
			if err != nil {
				return Result{State: StateError, Detail: "request failed", Err: err}
			}

			body := strings.ToLower(resp.Text())
			switch {
			case strings.Contains(body, "youtube premium is not available in your country"):
				return Result{State: StateBlocked, Detail: "not offered in this country"}
			case strings.Contains(body, "ad-free"):
				return Result{State: StateAvailable}
			}
			return Result{
				State:  StateError,
				Detail: "neither the offer nor the refusal wording was present",
			}
		},
	}
}

// Netflix title IDs used to separate a full catalogue from an originals-only one.
//
// The first is a licensed third-party title, which Netflix only serves where it
// holds regional rights. The second is a Netflix original, available anywhere
// Netflix operates at all. Reaching the original but not the licensed title is
// the signature of an address Netflix serves in "originals only" mode, which is
// what it does for traffic it believes comes from a proxy or hosting provider.
const (
	netflixLicensedTitle = "70143836"
	netflixOriginalTitle = "80197526"
)

func netflix() Check {
	return Check{
		ID: "netflix_access", Name: "Netflix",
		Run: func(ctx context.Context, env Env) Result {
			region, ok, err := netflixTitle(ctx, env, netflixLicensedTitle)
			if err != nil {
				return Result{State: StateError, Detail: "request failed", Err: err}
			}
			if ok {
				return Result{
					State: StateAvailable, Region: region,
					Detail: "full catalogue",
				}
			}

			region, ok, err = netflixTitle(ctx, env, netflixOriginalTitle)
			if err != nil {
				return Result{State: StateError, Detail: "request failed", Err: err}
			}
			if ok {
				return Result{
					State: StateRestricted, Region: region,
					Detail: "originals only; licensed titles are withheld",
				}
			}
			return Result{State: StateBlocked, Detail: "not available"}
		},
	}
}

// netflixBase is the origin the title checks run against; tests point it at a
// local server.
const netflixBase = "https://www.netflix.com"

func netflixTitle(ctx context.Context, env Env, id string) (region string, ok bool, err error) {
	return netflixTitleAt(ctx, env, netflixBase, id)
}

// netflixTitleAt fetches a title page and reads the region out of where Netflix
// redirects. A locale prefix such as /fr-en/title/... names the catalogue being
// served; a bare /title/... means the US one.
func netflixTitleAt(ctx context.Context, env Env, base, id string) (region string, ok bool, err error) {
	resp, err := env.Stack.Do(ctx, env.Family, netx.Request{
		URL:       base + "/title/" + id,
		UserAgent: browserUA,
	})
	if err != nil {
		return "", false, err
	}
	if resp.Status == 404 {
		// Netflix answers 404 for a title it will not serve here.
		return "", false, nil
	}
	if !resp.OK() {
		return "", false, nil
	}

	final := resp.FinalURL
	if final == "" {
		return "", false, nil
	}
	u, perr := url.Parse(final)
	if perr != nil {
		return "", false, nil
	}

	// Path is either /title/<id> or /<locale>/title/<id>.
	seg := strings.SplitN(strings.TrimPrefix(u.Path, "/"), "/", 2)
	if len(seg) == 0 || seg[0] == "" {
		return "", false, nil
	}
	if seg[0] == "title" {
		return "US", true, nil
	}
	// A locale looks like "fr-en"; the country is the part before the dash.
	locale, _, _ := strings.Cut(seg[0], "-")
	if len(locale) != 2 {
		return "", false, nil
	}
	return strings.ToUpper(locale), true, nil
}

// chatGPTWeb uses the endpoint the web client calls before showing the login
// form; it names an unsupported country outright.
func chatGPTWeb() Check {
	return Check{
		ID: "chatgpt_web", Name: "ChatGPT (web)",
		Run: func(ctx context.Context, env Env) Result {
			resp, err := env.Stack.Do(ctx, env.Family, netx.Request{
				URL:       "https://api.openai.com/compliance/cookie_requirements",
				UserAgent: browserUA,
			})
			if err != nil {
				return Result{State: StateError, Detail: "request failed", Err: err}
			}
			if containsFold(resp.Text(), "unsupported_country") {
				return Result{State: StateBlocked, Detail: "unsupported country"}
			}
			return Result{State: StateAvailable}
		},
	}
}

// chatGPTApp probes the host the iOS app talks to, which rejects hosting and
// proxy address space specifically — a different gate from the country check.
func chatGPTApp() Check {
	return Check{
		ID: "chatgpt_app", Name: "ChatGPT (app)",
		Run: func(ctx context.Context, env Env) Result {
			resp, err := env.Stack.Do(ctx, env.Family, netx.Request{
				URL:       "https://ios.chat.openai.com",
				UserAgent: browserUA,
			})
			if err != nil {
				return Result{State: StateError, Detail: "request failed", Err: err}
			}
			body := resp.Text()
			if isChallenge(resp.Status, strings.ToLower(body)) {
				return Result{
					State:  StateError,
					Detail: "Cloudflare challenged the request, so availability was never tested",
				}
			}
			switch {
			case containsFold(body, "disallowed isp"):
				return Result{
					State:  StateBlocked,
					Detail: "disallowed ISP; the exit address is classed as hosting or proxy space",
				}
			case containsFold(body, "been blocked"):
				return Result{State: StateBlocked, Detail: "blocked"}
			}
			return Result{State: StateAvailable}
		},
	}
}

var reCFLoc = regexp.MustCompile(`loc=([A-Z]{2})`)

// claudeUnavailableMarkers are the strings the region-refusal page carries. The
// apostrophe appears both literally and HTML-escaped depending on how the page
// is rendered, so both forms are matched.
var claudeUnavailableMarkers = []string{
	"app unavailable in region",
	"/app-unavailable-in-region",
	"unfortunately, claude isn't available here.",
	"unfortunately, claude isn&apos;t available here.",
	"unfortunately, claude isn&#39;t available here.",
}

// challengeMarkers identify Cloudflare's interstitial. The first two are the
// variable names its script sets, which are the most stable part of the page;
// the rest are the wording, which changes more often.
var challengeMarkers = []string{
	"cf_chl_opt",
	"_cf_chl",
	"just a moment",
	"cf-browser-verification",
	"enable javascript and cookies to continue",
}

// isChallenge reports whether a response came from Cloudflare's edge rather
// than from the service behind it.
//
// This matters because the challenge is issued before the origin ever sees the
// request, so its body says nothing about whether the service would serve this
// region. It is also why the check below tests for it first: the challenge page
// echoes the requested path back inside its own markup, so a request for
// /app-unavailable-in-region returns a Cloudflare page containing the string
// "app-unavailable-in-region" — which would otherwise read as proof of exactly
// the region block that was never actually established.
func isChallenge(status int, lowerBody string) bool {
	switch status {
	case 403, 429, 503:
	default:
		return false
	}
	for _, m := range challengeMarkers {
		if strings.Contains(lowerBody, m) {
			return true
		}
	}
	return false
}

// classifyClaude turns one response into a verdict. It is separate from the
// request so the decision can be tested against captured pages.
func classifyClaude(status int, body, region string) Result {
	lower := strings.ToLower(body)

	if isChallenge(status, lower) {
		return Result{
			State: StateError, Region: region,
			Detail: "Cloudflare challenged the request, so availability was never tested",
		}
	}
	for _, marker := range claudeUnavailableMarkers {
		if strings.Contains(lower, marker) {
			return Result{State: StateBlocked, Region: region, Detail: "not available in this region"}
		}
	}
	if status == 403 {
		// Refused, but without the page that would attribute it to geography.
		// Reporting "blocked" here would state as fact something the response
		// does not show, so this stays an unknown.
		return Result{
			State: StateError, Region: region,
			Detail: "HTTP 403 without the region page, so the cause is unknown",
		}
	}
	if status >= 200 && status < 400 {
		return Result{State: StateAvailable, Region: region}
	}
	return Result{
		State: StateError, Region: region,
		Detail: "unexpected HTTP " + itoa(status),
	}
}

func claude() Check {
	return Check{
		ID: "claude_access", Name: "Claude",
		Run: func(ctx context.Context, env Env) Result {
			region := cloudflareLoc(ctx, env, "https://claude.ai/cdn-cgi/trace")

			resp, err := env.Stack.Do(ctx, env.Family, netx.Request{
				URL:       "https://claude.ai/",
				UserAgent: browserUA,
				Headers: map[string]string{
					"Accept":                    "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,*/*;q=0.8",
					"Accept-Language":           "en-US,en;q=0.9",
					"Sec-Fetch-Dest":            "document",
					"Sec-Fetch-Mode":            "navigate",
					"Sec-Fetch-Site":            "none",
					"Sec-Fetch-User":            "?1",
					"Upgrade-Insecure-Requests": "1",
				},
			})
			if err != nil {
				return Result{State: StateError, Region: region, Detail: "request failed", Err: err}
			}
			return classifyClaude(resp.Status, resp.Text(), region)
		},
	}
}

// tiktokBlockMarkers are the phrasings TikTok uses to refuse a region.
var tiktokBlockMarkers = []string{
	"service is currently unavailable in your region",
	"tiktok is not available in your country",
	"tiktok is unavailable in your country",
	"not available in your region",
}

func tiktok() Check {
	return Check{
		ID: "tiktok_access", Name: "TikTok",
		Run: func(ctx context.Context, env Env) Result {
			resp, err := env.Stack.Do(ctx, env.Family, netx.Request{
				URL:       "https://www.tiktok.com/",
				UserAgent: browserUA,
				Headers: map[string]string{
					"Accept":          "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8",
					"Accept-Language": "en-US,en;q=0.9",
				},
			})
			if err != nil {
				return Result{State: StateError, Detail: "request failed", Err: err}
			}

			body := strings.ToLower(resp.Text())
			for _, marker := range tiktokBlockMarkers {
				if strings.Contains(body, marker) {
					return Result{State: StateBlocked, Detail: "blocked by region"}
				}
			}
			if resp.Status >= 200 && resp.Status < 400 {
				return Result{State: StateAvailable}
			}
			return Result{State: StateError, Detail: "unexpected HTTP " + itoa(resp.Status)}
		},
	}
}

// cloudflareLoc reads the country from a Cloudflare /cdn-cgi/trace document.
// It is best effort: the availability verdict does not depend on it.
func cloudflareLoc(ctx context.Context, env Env, traceURL string) string {
	resp, err := env.Stack.Do(ctx, env.Family, netx.Request{
		URL:       traceURL,
		UserAgent: browserUA,
		Headers:   map[string]string{"Accept": "text/plain,*/*;q=0.8"},
	})
	if err != nil {
		return ""
	}
	if m := reCFLoc.FindStringSubmatch(resp.Text()); m != nil {
		return m[1]
	}
	return ""
}

func containsFold(haystack, needle string) bool {
	return strings.Contains(strings.ToLower(haystack), strings.ToLower(needle))
}

func itoa(n int) string {
	if n == 0 {
		return "0"
	}
	neg := n < 0
	if neg {
		n = -n
	}
	var b [12]byte
	i := len(b)
	for n > 0 {
		i--
		b[i] = byte('0' + n%10)
		n /= 10
	}
	if neg {
		i--
		b[i] = '-'
	}
	return string(b[i:])
}

// reGeminiRegion pulls the served country out of the configuration block that
// Google's account bar embeds in every one of its pages. The code is ISO 3166-1
// alpha-3 and sits at a fixed position in that array, so the two leading
// numbers are matched loosely rather than pinned.
var reGeminiRegion = regexp.MustCompile(`,\d+,\d+,200,"([A-Z]{3})"`)

// geminiUnsupported are the countries where Google does not offer Gemini.
var geminiUnsupported = map[string]bool{
	"RUS": true,
	"BLR": true,
	"CHN": true,
	"PRK": true,
	"IRN": true,
	"CUB": true,
	"SYR": true,
}

// classifyGemini is separated from the request so the decision can be tested
// against captured pages.
func classifyGemini(status int, body string) Result {
	if isChallenge(status, strings.ToLower(body)) {
		return Result{
			State:  StateError,
			Detail: "Cloudflare challenged the request, so availability was never tested",
		}
	}
	if status < 200 || status >= 400 {
		return Result{State: StateError, Detail: "unexpected HTTP " + itoa(status)}
	}

	m := reGeminiRegion.FindStringSubmatch(body)
	if m == nil {
		return Result{
			State:  StateError,
			Detail: "could not read the served region from the page",
		}
	}

	region := m[1]
	if geminiUnsupported[region] {
		return Result{State: StateBlocked, Region: region, Detail: "not offered in this country"}
	}
	return Result{State: StateAvailable, Region: region}
}

func gemini() Check {
	return Check{
		ID: "gemini_access", Name: "Gemini",
		Run: func(ctx context.Context, env Env) Result {
			resp, err := env.Stack.Do(ctx, env.Family, netx.Request{
				URL:       "https://gemini.google.com/app",
				UserAgent: browserUA,
				Headers: map[string]string{
					"Accept":          "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8",
					"Accept-Language": "en-US,en;q=0.9",
				},
			})
			if err != nil {
				return Result{State: StateError, Detail: "request failed", Err: err}
			}
			return classifyGemini(resp.Status, resp.Text())
		},
	}
}

// classifyNotebookLM reads the answer out of where the request was sent rather
// than out of a page. NotebookLM states the reason in the redirect itself —
// `?location=unsupported` — which is a far better signal than anything its
// sibling Gemini offers, and the reason this check needs no country list.
//
// A served region is redirected to a Google sign-in instead, because the
// product needs an account. That is not a failure: reaching the sign-in is the
// evidence that the region is served.
func classifyNotebookLM(status int, finalURL, body string) Result {
	if isChallenge(status, strings.ToLower(body)) {
		return Result{
			State:  StateError,
			Detail: "Cloudflare challenged the request, so availability was never tested",
		}
	}

	if strings.Contains(finalURL, "location=unsupported") {
		return Result{State: StateBlocked, Detail: "not offered in this country"}
	}

	u, err := url.Parse(finalURL)
	if err != nil || finalURL == "" {
		return Result{State: StateError, Detail: "no destination to judge"}
	}
	switch {
	case u.Host == "accounts.google.com",
		strings.HasSuffix(u.Host, ".google.com") && strings.Contains(u.Path, "/login"):
		return Result{State: StateAvailable, Detail: "sign-in required"}
	}

	if status >= 200 && status < 400 {
		// Somewhere unexpected, but not the refusal. Better to say the
		// destination was unrecognised than to read it as either verdict.
		return Result{
			State:  StateError,
			Detail: "unexpected destination " + u.Host,
		}
	}
	return Result{State: StateError, Detail: "unexpected HTTP " + itoa(status)}
}

func notebookLM() Check {
	return Check{
		ID: "notebooklm_access", Name: "NotebookLM",
		Run: func(ctx context.Context, env Env) Result {
			resp, err := env.Stack.Do(ctx, env.Family, netx.Request{
				URL:       "https://notebooklm.google.com/",
				UserAgent: browserUA,
				Headers: map[string]string{
					"Accept":          "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8",
					"Accept-Language": "en-US,en;q=0.9",
				},
			})
			if err != nil {
				return Result{State: StateError, Detail: "request failed", Err: err}
			}
			return classifyNotebookLM(resp.Status, resp.FinalURL, resp.Text())
		},
	}
}
