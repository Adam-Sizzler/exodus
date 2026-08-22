package ai

import (
	"context"
	"strconv"
	"strings"

	"exodus-node/geocheck/internal/netx"
)

// Checks returns the endpoint probes, in display order.
func Checks() []Check {
	out := make([]Check, 0, len(catalog))
	for _, e := range catalog {
		out = append(out, e.check())
	}
	return out
}

// endpoint describes one API in the terms the probe needs.
type endpoint struct {
	id, name, vendor string
	url              string
	headers          map[string]string

	// authStatus are answers that prove reachability on their own.
	authStatus []int
	// regionMarkers are the phrases a service uses when it refuses a region.
	// Matched case-insensitively against the body.
	regionMarkers []string
}

var catalog = []endpoint{
	{
		id: "openai", name: "OpenAI", vendor: "OpenAI",
		url:        "https://api.openai.com/v1/models",
		authStatus: []int{401},
		regionMarkers: []string{
			"unsupported_country_region_territory",
			"country, region, or territory not supported",
		},
	},
	{
		id: "anthropic", name: "Anthropic", vendor: "Anthropic",
		url:        "https://api.anthropic.com/v1/models",
		headers:    map[string]string{"anthropic-version": "2023-06-01"},
		authStatus: []int{401},
		regionMarkers: []string{
			"request not allowed",
			"unsupported_country",
		},
	},
	{
		id: "gemini", name: "Google Gemini", vendor: "Google",
		url:        "https://generativelanguage.googleapis.com/v1beta/models",
		authStatus: []int{401, 403},
	},
	{
		id: "deepseek", name: "DeepSeek", vendor: "DeepSeek",
		url:        "https://api.deepseek.com/models",
		authStatus: []int{401},
	},
	{
		id: "qwen_intl", name: "Qwen (Singapore)", vendor: "Alibaba",
		url:        "https://dashscope-intl.aliyuncs.com/compatible-mode/v1/models",
		authStatus: []int{401},
	},
	{
		id: "qwen_cn", name: "Qwen (Beijing)", vendor: "Alibaba",
		url:        "https://dashscope.aliyuncs.com/compatible-mode/v1/models",
		authStatus: []int{401},
	},
	{
		id: "moonshot", name: "Kimi / Moonshot", vendor: "Moonshot",
		url:        "https://api.moonshot.cn/v1/models",
		authStatus: []int{401},
	},
	{
		id: "zhipu", name: "Zhipu GLM", vendor: "Zhipu",
		url:        "https://open.bigmodel.cn/api/paas/v4/models",
		authStatus: []int{401},
	},
}

func (e endpoint) check() Check {
	return Check{
		ID: e.id, Name: e.name, Vendor: e.vendor,
		Run: func(ctx context.Context, env Env) Result {
			resp, err := env.Stack.Do(ctx, env.Family, netx.Request{
				URL:     e.url,
				Headers: e.headers,
			})
			if err != nil {
				return Result{State: StateError, Detail: "request failed", Err: err}
			}
			return e.classify(resp.Status, resp.Text())
		},
	}
}

func (e endpoint) classify(status int, body string) Result {
	lower := strings.ToLower(body)

	if isChallenge(status, lower) {
		return Result{
			State: StateError, Status: status,
			Detail: "challenged before reaching the API",
		}
	}

	for _, m := range e.regionMarkers {
		if strings.Contains(lower, strings.ToLower(m)) {
			return Result{
				State: StateBlocked, Status: status,
				Detail: "the API refused this region",
			}
		}
	}

	for _, s := range e.authStatus {
		if status == s {
			return Result{
				State: StateReachable, Status: status,
				Detail: "answered, credentials not supplied",
			}
		}
	}

	if status >= 200 && status < 300 {
		return Result{State: StateReachable, Status: status, Detail: "answered"}
	}

	return Result{
		State: StateError, Status: status,
		Detail: "unexpected HTTP " + strconv.Itoa(status),
	}
}

func isChallenge(status int, lowerBody string) bool {
	switch status {
	case 403, 429, 503:
	default:
		return false
	}
	for _, m := range []string{"cf_chl_opt", "_cf_chl", "just a moment", "cf-browser-verification"} {
		if strings.Contains(lowerBody, m) {
			return true
		}
	}
	return false
}
