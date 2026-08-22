package netx

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"net/url"
	"strings"
)

const maxBodyBytes = 4 << 20 // 4 MiB: enough for the HTML pages we scrape

// Request describes one probe against a service endpoint.
type Request struct {
	Method    string
	URL       string
	Headers   map[string]string
	UserAgent string

	// Exactly one body form may be set.
	JSON string
	Form url.Values
	Raw  []byte

	// HeadOnly issues a HEAD request and keeps only the response headers.
	HeadOnly bool

	// NoRedirect returns the first response instead of following it. Captive
	// portal detection depends on this: the redirect *is* the finding, and
	// following it would report the portal's page as if it were the endpoint's.
	NoRedirect bool
}

// Response is the part of an HTTP reply the checks care about.
type Response struct {
	Status int
	Header http.Header
	Body   []byte

	// FinalURL is the URL the response actually came from, after redirects.
	// Some services encode the answer in where they send you rather than in
	// the page: Netflix redirects to a locale-prefixed path to say which
	// catalogue you are being served.
	FinalURL string
}

// Text returns the body as a string.
func (r *Response) Text() string {
	if r == nil {
		return ""
	}
	return string(r.Body)
}

// OK reports whether the status is in the 2xx range.
func (r *Response) OK() bool { return r != nil && r.Status >= 200 && r.Status < 300 }

// Do executes req over the given address family.
func (s *Stack) Do(ctx context.Context, f Family, req Request) (*Response, error) {
	method := req.Method
	if method == "" {
		method = http.MethodGet
	}
	if req.HeadOnly {
		method = http.MethodHead
	}

	var body io.Reader
	contentType := ""
	switch {
	case req.JSON != "":
		body = strings.NewReader(req.JSON)
		contentType = "application/json"
	case req.Form != nil:
		body = strings.NewReader(req.Form.Encode())
		contentType = "application/x-www-form-urlencoded"
	case req.Raw != nil:
		body = bytes.NewReader(req.Raw)
	}

	httpReq, err := http.NewRequestWithContext(ctx, method, req.URL, body)
	if err != nil {
		return nil, err
	}
	if contentType != "" {
		httpReq.Header.Set("Content-Type", contentType)
	}
	ua := req.UserAgent
	if ua == "" {
		ua = s.opts.UserAgent
	}
	httpReq.Header.Set("User-Agent", ua)
	// Accept-Encoding is deliberately left alone: net/http only decompresses
	// transparently when it added the header itself.
	for k, v := range req.Headers {
		httpReq.Header.Set(k, v)
	}

	client := s.HTTP(f)
	if req.NoRedirect {
		// A shallow copy shares the transport, so this costs nothing and keeps
		// connection reuse intact.
		noFollow := *client
		noFollow.CheckRedirect = func(*http.Request, []*http.Request) error {
			return http.ErrUseLastResponse
		}
		client = &noFollow
	}

	resp, err := client.Do(httpReq)
	if err != nil {
		return nil, err
	}
	// The body is only read, so a close error carries nothing actionable.
	defer func() { _ = resp.Body.Close() }()

	out := &Response{Status: resp.StatusCode, Header: resp.Header}
	if resp.Request != nil && resp.Request.URL != nil {
		out.FinalURL = resp.Request.URL.String()
	}
	if !req.HeadOnly {
		out.Body, err = io.ReadAll(io.LimitReader(resp.Body, maxBodyBytes))
		if err != nil {
			return nil, err
		}
	} else {
		_, _ = io.Copy(io.Discard, io.LimitReader(resp.Body, 4096))
	}
	return out, nil
}
