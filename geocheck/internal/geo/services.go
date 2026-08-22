package geo

import (
	"context"
	"errors"
	"net/http"
	"regexp"
	"strings"

	"exodus-node/geocheck/internal/countries"
	"exodus-node/geocheck/internal/jsonx"
	"exodus-node/geocheck/internal/netx"
)

const (
	// googleConsentCookie pre-accepts the EU consent interstitial so the real
	// page (and its locale hints) is returned instead of the consent wall.
	googleConsentCookie = "SOCS=CAISNQgDEitib3FfaWRlbnRpdHlmcm9udGVuZHVpc2VydmVyXzIwMjUwNzMwLjA1X3AwGgJlbiACGgYIgPC_xAY"

	spotifyAPIKey   = "142b583129b2df829de3656f9eb484e6"
	spotifyClientID = "9a8d2f0ce77a4e248bb71fefcb557637"

	yes = "yes"
	no  = "no"
)

var (
	reGoogleLocaleUnderscore = regexp.MustCompile(`"[a-z]{2}_([A-Z]{2})"`)
	reGoogleLocaleDash       = regexp.MustCompile(`"[a-z]{2}-([A-Z]{2})"`)
	rePlayCountry            = regexp.MustCompile(`<div class="yVZQTb">([^<(]+)`)
	reYouTubeCountry         = regexp.MustCompile(`"countryCode":"(\w+)"`)
	reSpotifyGeo             = regexp.MustCompile(`"geoLocationCountryCode":"([^"]*)"`)
	reDeezerCountry          = regexp.MustCompile(`'country':\s*'([^']*)'`)
	rePrimeTerritory         = regexp.MustCompile(`"currentTerritory":"([^"]+)"`)
	reBingRegion             = regexp.MustCompile(`Region\s*:\s*"([^"]+)"`)
	reLiveCountry            = regexp.MustCompile(`"sRequestCountry":"([^"]*)"`)
	reSteamCountry           = regexp.MustCompile(`steamCountry=([^%;]*)`)
	rePSCountry              = regexp.MustCompile(`(?i)country=([A-Za-z]{2})`)
	reGoogleBlocked          = regexp.MustCompile(`(?i)unusual traffic from|is blocked|unaddressed abuse`)
)

// ServiceChecks returns the consumer-service probes, in display order.
func ServiceChecks() []Check {
	return []Check{
		google(),
		googleSearchCaptcha(),
		youtube(),
		youtubePremium(),
		youtubeMusic(),
		twitch(),
		chatGPT(),
		netflix(),
		spotify(),
		spotifySignup(),
		deezer(),
		reddit(),
		redditGuest(),
		amazonPrime(),
		apple(),
		steam(),
		playstation(),
		tiktok(),
		ookla(),
		jetbrains(),
		bing(),
	}
}

// google reads the locale Google's homepage renders for you, falling back to
// the country name printed on the Play Store landing page.
func google() Check {
	return Check{
		ID: "google", Name: "Google", Group: GroupServices, Kind: KindCountry,
		Run: func(ctx context.Context, env *Env, f netx.Family) (string, error) {
			resp, err := env.Stack.Do(ctx, f, netx.Request{URL: "https://www.google.com"})
			if err == nil {
				body := resp.Text()
				if m := reGoogleLocaleUnderscore.FindStringSubmatch(body); m != nil {
					return m[1], nil
				}
				// The dash form appears many times; the last one is the
				// document's own locale rather than an alternate hreflang.
				if all := reGoogleLocaleDash.FindAllStringSubmatch(body, -1); len(all) > 0 {
					return all[len(all)-1][1], nil
				}
			}

			play, perr := env.Stack.Do(ctx, f, netx.Request{
				URL: "https://play.google.com/",
				Headers: map[string]string{
					"Accept":          "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8",
					"Accept-Language": "en-US;q=0.9",
				},
			})
			if perr != nil {
				if err != nil {
					return "", err
				}
				return "", perr
			}
			if m := rePlayCountry.FindStringSubmatch(play.Text()); m != nil {
				return countries.Code(m[1]), nil
			}
			return "", nil
		},
	}
}

// googleSearchCaptcha detects whether the exit IP is flagged for abuse.
func googleSearchCaptcha() Check {
	return Check{
		ID: "google_captcha", Name: "Google Search captcha", Group: GroupServices, Kind: KindBlocked,
		Run: func(ctx context.Context, env *Env, f netx.Family) (string, error) {
			resp, err := env.Stack.Do(ctx, f, netx.Request{
				URL:     "https://www.google.com/search?q=cats",
				Headers: map[string]string{"Accept-Language": "en-US,en;q=0.9"},
			})
			if err != nil {
				return "", err
			}
			if resp.Status == http.StatusTooManyRequests || reGoogleBlocked.MatchString(resp.Text()) {
				return yes, nil
			}
			return no, nil
		},
	}
}

// youtube reads the region YouTube serves; it falls back to the Google answer
// because the two normally agree and YouTube's marker is not always present.
func youtube() Check {
	return Check{
		ID: "youtube", Name: "YouTube", Group: GroupServices, Kind: KindCountry,
		DependsOn: "google",
		Run: func(ctx context.Context, env *Env, f netx.Family) (string, error) {
			resp, err := env.Stack.Do(ctx, f, netx.Request{URL: "https://www.youtube.com"})
			if err == nil {
				if m := reYouTubeCountry.FindStringSubmatch(resp.Text()); m != nil && len(m[1]) == 2 {
					return m[1], nil
				}
			}
			if fallback := env.lookupShared("google", f); fallback != "" {
				return fallback, nil
			}
			return "", err
		},
	}
}

// availabilityByPhrase is the shape shared by the "is this product offered
// here" checks: fetch a page, look for the refusal sentence.
func availabilityByPhrase(id, name, url, phrase string, headers map[string]string) Check {
	needle := strings.ToLower(phrase)
	return Check{
		ID: id, Name: name, Group: GroupServices, Kind: KindAvailability,
		Run: func(ctx context.Context, env *Env, f netx.Family) (string, error) {
			resp, err := env.Stack.Do(ctx, f, netx.Request{URL: url, Headers: headers})
			if err != nil {
				return "", err
			}
			if len(resp.Body) == 0 {
				return "", nil
			}
			if strings.Contains(strings.ToLower(resp.Text()), needle) {
				return no, nil
			}
			return yes, nil
		},
	}
}

func youtubePremium() Check {
	return availabilityByPhrase(
		"youtube_premium", "YouTube Premium",
		"https://www.youtube.com/premium",
		"youtube premium is not available in your country",
		map[string]string{"Cookie": googleConsentCookie, "Accept-Language": "en-US,en;q=0.9"},
	)
}

func youtubeMusic() Check {
	return availabilityByPhrase(
		"youtube_music", "YouTube Music",
		"https://music.youtube.com/",
		"youtube music is not available in your area",
		map[string]string{"Cookie": googleConsentCookie, "Accept-Language": "en-US,en;q=0.9"},
	)
}

func twitch() Check {
	return Check{
		ID: "twitch", Name: "Twitch", Group: GroupServices, Kind: KindCountry,
		Run: func(ctx context.Context, env *Env, f netx.Family) (string, error) {
			resp, err := env.Stack.Do(ctx, f, netx.Request{
				Method:  http.MethodPost,
				URL:     "https://gql.twitch.tv/gql",
				Headers: map[string]string{"Client-Id": "kimne78kx3ncx6brgo4mv6wki5h1ko"},
				JSON: `[{"operationName":"VerifyEmail_CurrentUser","variables":{},"extensions":` +
					`{"persistedQuery":{"version":1,"sha256Hash":` +
					`"f9e7dcdf7e99c314c82d8f7f725fab5f99d1df3d7359b53c9ae122deec590198"}}}]`,
			})
			if err != nil {
				return "", err
			}
			return jsonx.String(resp.Body, "0.data.requestInfo.countryCode"), nil
		},
	}
}

func chatGPT() Check {
	return Check{
		ID: "chatgpt", Name: "ChatGPT", Group: GroupServices, Kind: KindCountry,
		Run: func(ctx context.Context, env *Env, f netx.Family) (string, error) {
			resp, err := env.Stack.Do(ctx, f, netx.Request{
				Method: http.MethodPost,
				URL:    "https://ab.chatgpt.com/v1/initialize",
				Headers: map[string]string{
					"Statsig-Api-Key": "client-zUdXdSTygXJdzoE0sWTkP8GKTVsUMF2IRM7ShVO2JAG",
				},
			})
			if err != nil {
				return "", err
			}
			return jsonx.String(resp.Body, "derived_fields.country"), nil
		},
	}
}

// fastAPIURL is Netflix's own speedtest bootstrap; it reports both the client
// location and the Open Connect target that would serve the stream.
const fastAPIURL = "https://api.fast.com/netflix/speedtest/v2?https=true&" +
	"token=YXNkZmFzZGxmbnNkYWZoYXNkZmhrYWxm&urlCount=1"

func netflix() Check {
	return Check{
		ID: "netflix", Name: "Netflix", Group: GroupServices, Kind: KindCountry,
		Run: func(ctx context.Context, env *Env, f netx.Family) (string, error) {
			resp, err := env.Stack.Do(ctx, f, netx.Request{URL: fastAPIURL})
			if err != nil {
				return "", err
			}
			return jsonx.String(resp.Body, "client.location.country"), nil
		},
	}
}

func spotify() Check {
	return Check{
		ID: "spotify", Name: "Spotify", Group: GroupServices, Kind: KindCountry,
		Run: func(ctx context.Context, env *Env, f netx.Family) (string, error) {
			resp, err := env.Stack.Do(ctx, f, netx.Request{URL: "https://accounts.spotify.com/status"})
			if err != nil {
				return "", err
			}
			if m := reSpotifyGeo.FindStringSubmatch(resp.Text()); m != nil {
				return m[1], nil
			}
			return "", nil
		},
	}
}

func spotifySignup() Check {
	return Check{
		ID: "spotify_signup", Name: "Spotify signup", Group: GroupServices, Kind: KindAvailability,
		Run: func(ctx context.Context, env *Env, f netx.Family) (string, error) {
			resp, err := env.Stack.Do(ctx, f, netx.Request{
				URL: "https://spclient.wg.spotify.com/signup/public/v1/account/?validate=1&key=" +
					spotifyAPIKey,
				Headers: map[string]string{"X-Client-Id": spotifyClientID},
			})
			if err != nil {
				return "", err
			}
			status := jsonx.String(resp.Body, "status")
			launched := jsonx.String(resp.Body, "is_country_launched")
			if status == "120" || status == "320" || launched == "false" {
				return no, nil
			}
			if status == "" && launched == "" {
				return "", nil
			}
			return yes, nil
		},
	}
}

func deezer() Check {
	return Check{
		ID: "deezer", Name: "Deezer", Group: GroupServices, Kind: KindCountry,
		Run: func(ctx context.Context, env *Env, f netx.Family) (string, error) {
			resp, err := env.Stack.Do(ctx, f, netx.Request{URL: "https://www.deezer.com/en/offers"})
			if err != nil {
				return "", err
			}
			if m := reDeezerCountry.FindStringSubmatch(resp.Text()); m != nil {
				return m[1], nil
			}
			return "", nil
		},
	}
}

// reddit exchanges an anonymous device token for an access token, then asks
// the federated GraphQL endpoint where it thinks the caller is.
// reRedditCountry pulls the served country out of the chat widget's markup.
var reRedditCountry = regexp.MustCompile(`country="([A-Za-z]{2})"`)

const redditBlocked = "blocked by network security"

func reddit() Check {
	return Check{
		ID: "reddit", Name: "Reddit", Group: GroupServices, Kind: KindCountry,
		Run: func(ctx context.Context, env *Env, f netx.Family) (string, error) {
			resp, err := env.Stack.Do(ctx, f, netx.Request{
				URL: "https://www.reddit.com/svc/shreddit/reddit-chat",
			})
			if err != nil {
				return "", err
			}
			body := resp.Text()
			if strings.Contains(strings.ToLower(body), redditBlocked) {
				return "", errors.New("blocked by Reddit's network security, " +
					"which turns away hosting ranges regardless of country")
			}
			if m := reRedditCountry.FindStringSubmatch(body); m != nil {
				return strings.ToUpper(m[1]), nil
			}
			return "", nil
		},
	}
}

func redditGuest() Check {
	return Check{
		ID: "reddit_guest", Name: "Reddit guest access", Group: GroupServices, Kind: KindAvailability,
		Run: func(ctx context.Context, env *Env, f netx.Family) (string, error) {
			resp, err := env.Stack.Do(ctx, f, netx.Request{URL: "https://www.reddit.com"})
			if err != nil {
				return no, nil
			}
			if resp.OK() && len(resp.Body) > 0 {
				return yes, nil
			}
			return no, nil
		},
	}
}

func amazonPrime() Check {
	return Check{
		ID: "prime", Name: "Amazon Prime Video", Group: GroupServices, Kind: KindCountry,
		Run: func(ctx context.Context, env *Env, f netx.Family) (string, error) {
			resp, err := env.Stack.Do(ctx, f, netx.Request{URL: "https://www.primevideo.com"})
			if err != nil {
				return "", err
			}
			if m := rePrimeTerritory.FindStringSubmatch(resp.Text()); m != nil {
				return m[1], nil
			}
			return "", nil
		},
	}
}

// apple asks the device-provisioning endpoint iOS itself uses; it answers with
// a bare country code.
func apple() Check {
	return Check{
		ID: "apple", Name: "Apple", Group: GroupServices, Kind: KindCountry,
		Run: func(ctx context.Context, env *Env, f netx.Family) (string, error) {
			resp, err := env.Stack.Do(ctx, f, netx.Request{
				URL: "https://gspe1-ssl.ls.apple.com/pep/gcc",
			})
			if err != nil {
				return "", err
			}
			return strings.TrimSpace(resp.Text()), nil
		},
	}
}

func steam() Check {
	return Check{
		ID: "steam", Name: "Steam", Group: GroupServices, Kind: KindCountry,
		Run: func(ctx context.Context, env *Env, f netx.Family) (string, error) {
			resp, err := env.Stack.Do(ctx, f, netx.Request{
				URL: "https://store.steampowered.com", HeadOnly: true,
			})
			if err != nil {
				return "", err
			}
			for _, c := range resp.Header.Values("Set-Cookie") {
				if m := reSteamCountry.FindStringSubmatch(c); m != nil {
					return m[1], nil
				}
			}
			return "", nil
		},
	}
}

func playstation() Check {
	return Check{
		ID: "playstation", Name: "PlayStation", Group: GroupServices, Kind: KindCountry,
		Run: func(ctx context.Context, env *Env, f netx.Family) (string, error) {
			resp, err := env.Stack.Do(ctx, f, netx.Request{
				URL: "https://www.playstation.com", HeadOnly: true,
			})
			if err != nil {
				return "", err
			}
			for _, c := range resp.Header.Values("Set-Cookie") {
				if strings.HasPrefix(c, "country=") {
					if m := rePSCountry.FindStringSubmatch(c); m != nil {
						return m[1], nil
					}
				}
			}
			return "", nil
		},
	}
}

func tiktok() Check {
	return Check{
		ID: "tiktok", Name: "TikTok", Group: GroupServices, Kind: KindCountry,
		Run: func(ctx context.Context, env *Env, f netx.Family) (string, error) {
			resp, err := env.Stack.Do(ctx, f, netx.Request{
				URL: "https://www.tiktok.com/api/v1/web-cookie-privacy/config?appId=1988",
			})
			if err != nil {
				return "", err
			}
			return jsonx.String(resp.Body, "body.appProps.region"), nil
		},
	}
}

func ookla() Check {
	return Check{
		ID: "ookla", Name: "Ookla Speedtest", Group: GroupServices, Kind: KindCountry,
		Run: func(ctx context.Context, env *Env, f netx.Family) (string, error) {
			resp, err := env.Stack.Do(ctx, f, netx.Request{
				URL: "https://www.speedtest.net/api/js/config-sdk",
			})
			if err != nil {
				return "", err
			}
			return jsonx.String(resp.Body, "location.countryCode"), nil
		},
	}
}

func jetbrains() Check {
	return Check{
		ID: "jetbrains", Name: "JetBrains", Group: GroupServices, Kind: KindCountry,
		Run: func(ctx context.Context, env *Env, f netx.Family) (string, error) {
			resp, err := env.Stack.Do(ctx, f, netx.Request{
				URL: "https://data.services.jetbrains.com/geo",
			})
			if err != nil {
				return "", err
			}
			return jsonx.String(resp.Body, "code"), nil
		},
	}
}

// bing reads the market Bing assigns. A redirect to cn.bing.com means the
// China market; a "WW" (worldwide) answer is unhelpful, so fall back to the
// Microsoft account endpoint which reports a real country.
func bing() Check {
	return Check{
		ID: "bing", Name: "Microsoft (Bing)", Group: GroupServices, Kind: KindCountry,
		Run: func(ctx context.Context, env *Env, f netx.Family) (string, error) {
			resp, err := env.Stack.Do(ctx, f, netx.Request{URL: "https://www.bing.com/search?q=cats"})
			if err != nil {
				return "", err
			}
			body := resp.Text()
			if strings.Contains(body, "cn.bing.com") {
				return "CN", nil
			}
			region := ""
			if m := reBingRegion.FindStringSubmatch(body); m != nil {
				region = m[1]
			}
			if region != "" && !strings.EqualFold(region, "WW") {
				return region, nil
			}

			live, err := env.Stack.Do(ctx, f, netx.Request{URL: "https://login.live.com"})
			if err != nil {
				return region, nil
			}
			if m := reLiveCountry.FindStringSubmatch(live.Text()); m != nil {
				return m[1], nil
			}
			return region, nil
		},
	}
}
