package cerberus

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/cerberus/subscription-page/backend/internal/config"
)

const RealIPHeader = "x-cerberus-real-ip"

var ignoredProxyHeaders = map[string]struct{}{
	"accept-encoding":              {},
	"alt-svc":                      {},
	"authorization":                {},
	"cache-control":                {},
	"cf-access-client-id":          {},
	"cf-access-client-secret":      {},
	"cf-cache-status":              {},
	"cf-ray":                       {},
	"connection":                   {},
	"content-length":               {},
	"content-security-policy":      {},
	"cross-origin-opener-policy":   {},
	"cross-origin-resource-policy": {},
	"expires":                      {},
	"host":                         {},
	"keep-alive":                   {},
	"nel":                          {},
	"origin-agent-cluster":         {},
	"pragma":                       {},
	"proxy-authenticate":           {},
	"proxy-authorization":          {},
	"report-to":                    {},
	"server":                       {},
	"te":                           {},
	"trailer":                      {},
	"transfer-encoding":            {},
	"upgrade":                      {},
	"x-api-key":                    {},
	"x-forwarded-for":              {},
	"x-forwarded-proto":            {},
	"x-forwarded-scheme":           {},
}

type Client struct {
	httpClient *http.Client
	baseURL    string
	headers    http.Header
}

type SubpageConfigEnvelope struct {
	WebpageAllowed    bool   `json:"webpageAllowed"`
	SubpageConfigUUID string `json:"subpageConfigUuid"`
}

type metadataResponse struct {
	Response struct {
		Version string `json:"version"`
	} `json:"response"`
}

type subpageConfigListResponse struct {
	Response struct {
		Configs []struct {
			UUID string `json:"uuid"`
		} `json:"configs"`
	} `json:"response"`
}

type subpageConfigByUUIDResponse struct {
	Response struct {
		Config json.RawMessage `json:"config"`
	} `json:"response"`
}

type subpageConfigByShortResponse struct {
	Response SubpageConfigEnvelope `json:"response"`
}

type userByUsernameResponse struct {
	Response struct {
		ShortUUID string `json:"shortUuid"`
	} `json:"response"`
}

func NewClient(cfg config.Config) *Client {
	headers := make(http.Header)
	headers.Set("User-Agent", "Cerberus Subscription Page")
	headers.Set("Authorization", "Bearer "+cfg.CerberusAPIToken)

	if cfg.CaddyAuthAPIToken != "" {
		headers.Set("X-Api-Key", cfg.CaddyAuthAPIToken)
	}

	if cfg.CloudflareZeroTrustClientID != "" && cfg.CloudflareZeroTrustClientSecret != "" {
		headers.Set("CF-Access-Client-Id", cfg.CloudflareZeroTrustClientID)
		headers.Set("CF-Access-Client-Secret", cfg.CloudflareZeroTrustClientSecret)
	}

	if strings.HasPrefix(cfg.CerberusPanelURL, "http://") {
		headers.Set("X-Forwarded-For", "127.0.0.1")
		headers.Set("X-Forwarded-Proto", "https")
	}

	return &Client{
		httpClient: &http.Client{Timeout: 10 * time.Second},
		baseURL:    cfg.CerberusPanelURL,
		headers:    headers,
	}
}

func (c *Client) GetMetadata(ctx context.Context) (string, error) {
	body, _, status, err := c.request(ctx, http.MethodGet, "/api/system/metadata", nil, nil)
	if err != nil {
		return "", err
	}

	if status < 200 || status >= 300 {
		return "", fmt.Errorf("unexpected metadata status: %d", status)
	}

	var response metadataResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return "", err
	}

	if response.Response.Version == "" {
		return "", fmt.Errorf("metadata response version is empty")
	}

	return response.Response.Version, nil
}

func (c *Client) GetSubscriptionPageConfigList(ctx context.Context) ([]string, error) {
	body, _, status, err := c.request(ctx, http.MethodGet, "/api/subscription-page-configs", nil, nil)
	if err != nil {
		return nil, err
	}

	if status == http.StatusNotFound {
		return nil, fmt.Errorf("this subscription page requires Cerberus Panel version >= 2.4.0")
	}

	if status < 200 || status >= 300 {
		return nil, fmt.Errorf("unexpected subpage config list status: %d", status)
	}

	var response subpageConfigListResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, err
	}

	configs := make([]string, 0, len(response.Response.Configs))
	for _, configEntry := range response.Response.Configs {
		if configEntry.UUID != "" {
			configs = append(configs, configEntry.UUID)
		}
	}

	return configs, nil
}

func (c *Client) GetSubscriptionPageConfigByUUID(ctx context.Context, uuid string) (json.RawMessage, error) {
	body, _, status, err := c.request(
		ctx,
		http.MethodGet,
		"/api/subscription-page-configs/"+url.PathEscape(uuid),
		nil,
		nil,
	)
	if err != nil {
		return nil, err
	}

	if status < 200 || status >= 300 {
		return nil, fmt.Errorf("unexpected subpage config status: %d", status)
	}

	var response subpageConfigByUUIDResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, err
	}

	if len(response.Response.Config) == 0 {
		return nil, fmt.Errorf("empty config payload for %s", uuid)
	}

	return response.Response.Config, nil
}

func (c *Client) GetUserByUsername(ctx context.Context, clientIP, username string) (string, error) {
	paths := []string{
		"/api/users/by-username/" + url.PathEscape(username),
		"/api/users/username/" + url.PathEscape(username),
		"/api/users/find/" + url.PathEscape(username),
	}

	headers := http.Header{}
	headers.Set(RealIPHeader, clientIP)

	for _, candidate := range paths {
		body, _, status, err := c.request(ctx, http.MethodGet, candidate, nil, headers)
		if err != nil {
			continue
		}

		if status == http.StatusNotFound || status < 200 || status >= 300 {
			continue
		}

		var response userByUsernameResponse
		if err := json.Unmarshal(body, &response); err != nil {
			continue
		}

		if response.Response.ShortUUID != "" {
			return response.Response.ShortUUID, nil
		}
	}

	return "", fmt.Errorf("user %q not found", username)
}

func (c *Client) GetSubscriptionInfo(ctx context.Context, clientIP, shortUUID string) (map[string]any, error) {
	paths := []string{
		"/api/sub/" + url.PathEscape(shortUUID) + "/info",
		"/api/subscriptions/" + url.PathEscape(shortUUID) + "/info",
		"/api/subscriptions/by-short-uuid/" + url.PathEscape(shortUUID),
		"/api/subscriptions/short-uuid/" + url.PathEscape(shortUUID),
		"/api/subscriptions/" + url.PathEscape(shortUUID),
	}

	headers := http.Header{}
	headers.Set(RealIPHeader, clientIP)

	for _, candidate := range paths {
		body, _, status, err := c.request(ctx, http.MethodGet, candidate, nil, headers)
		if err != nil {
			continue
		}

		if status < 200 || status >= 300 {
			continue
		}

		var payload map[string]any
		if err := json.Unmarshal(body, &payload); err != nil {
			continue
		}

		return payload, nil
	}

	return nil, fmt.Errorf("subscription info for %s not found", shortUUID)
}

func (c *Client) GetSubpageConfig(
	ctx context.Context,
	shortUUID string,
	requestHeaders http.Header,
) (SubpageConfigEnvelope, error) {
	payload := map[string]any{
		"requestHeaders": flattenHeadersForJSON(requestHeaders),
	}

	body, _, status, err := c.request(
		ctx,
		http.MethodPost,
		"/api/subscriptions/subpage-config/"+url.PathEscape(shortUUID),
		payload,
		nil,
	)
	if err != nil {
		return SubpageConfigEnvelope{}, err
	}

	if status < 200 || status >= 300 {
		return SubpageConfigEnvelope{}, fmt.Errorf(
			"unexpected subpage config by short uuid status: %d",
			status,
		)
	}

	var response subpageConfigByShortResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return SubpageConfigEnvelope{}, err
	}

	return response.Response, nil
}

func (c *Client) GetSubscription(
	ctx context.Context,
	clientIP, shortUUID, clientType string,
	incomingHeaders http.Header,
) ([]byte, http.Header, error) {
	endpoint := "/api/sub/" + url.PathEscape(shortUUID)
	if clientType != "" {
		endpoint += "/" + url.PathEscape(clientType)
	}

	headers := http.Header{}
	for key, values := range incomingHeaders {
		if _, ignored := ignoredProxyHeaders[strings.ToLower(key)]; ignored {
			continue
		}

		for _, value := range values {
			headers.Add(key, value)
		}
	}

	headers.Set("Accept", "*/*")
	headers.Set("Cache-Control", "no-cache, no-store, must-revalidate, private, max-age=0")
	headers.Set("Pragma", "no-cache")
	headers.Set("Expires", "0")
	headers.Set(RealIPHeader, clientIP)

	body, responseHeaders, status, err := c.request(ctx, http.MethodGet, endpoint, nil, headers)
	if err != nil {
		return nil, nil, err
	}

	if status == http.StatusNotFound {
		return nil, nil, fmt.Errorf("subscription %s not found in Cerberus", shortUUID)
	}

	if status < 200 || status >= 300 {
		return nil, nil, fmt.Errorf("unexpected subscription status: %d", status)
	}

	return body, responseHeaders, nil
}

func (c *Client) request(
	ctx context.Context,
	method, requestPath string,
	body any,
	extraHeaders http.Header,
) ([]byte, http.Header, int, error) {
	var reader io.Reader
	if body != nil {
		encoded, err := json.Marshal(body)
		if err != nil {
			return nil, nil, 0, err
		}
		reader = strings.NewReader(string(encoded))
	}

	request, err := http.NewRequestWithContext(ctx, method, c.baseURL+requestPath, reader)
	if err != nil {
		return nil, nil, 0, err
	}

	for key, values := range c.headers {
		for _, value := range values {
			request.Header.Add(key, value)
		}
	}

	for key, values := range extraHeaders {
		for _, value := range values {
			request.Header.Add(key, value)
		}
	}

	if body != nil {
		request.Header.Set("Content-Type", "application/json")
	}

	response, err := c.httpClient.Do(request)
	if err != nil {
		return nil, nil, 0, err
	}
	defer response.Body.Close()

	responseBody, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, nil, response.StatusCode, err
	}

	return responseBody, response.Header.Clone(), response.StatusCode, nil
}

func flattenHeadersForJSON(headers http.Header) map[string]any {
	result := make(map[string]any, len(headers))
	for key, values := range headers {
		if len(values) == 0 {
			result[key] = ""
			continue
		}

		if len(values) == 1 {
			result[key] = values[0]
			continue
		}

		clone := make([]string, 0, len(values))
		clone = append(clone, values...)
		result[key] = clone
	}

	return result
}
