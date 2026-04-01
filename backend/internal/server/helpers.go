package server

import (
	"html"
	"net"
	"net/http"
	"path/filepath"
	"strings"
)

func renderIndexTemplate(templateContent, metaTitle, metaDescription, panelData string) string {
	rendered := strings.ReplaceAll(templateContent, "<%- panelData %>", panelData)
	rendered = strings.ReplaceAll(rendered, "<%= metaDescription %>", html.EscapeString(metaDescription))
	rendered = strings.ReplaceAll(rendered, "<%- metaDescription %>", html.EscapeString(metaDescription))
	rendered = strings.ReplaceAll(rendered, "<%= metaTitle %>", html.EscapeString(metaTitle))
	rendered = strings.ReplaceAll(rendered, "<%- metaTitle %>", html.EscapeString(metaTitle))
	return rendered
}

func getRealIP(r *http.Request) string {
	forwardedFor := strings.TrimSpace(r.Header.Get("X-Forwarded-For"))
	if forwardedFor != "" {
		parts := strings.Split(forwardedFor, ",")
		if len(parts) > 0 {
			first := strings.TrimSpace(parts[0])
			if first != "" {
				return first
			}
		}
	}

	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err == nil && host != "" {
		return host
	}

	if r.RemoteAddr != "" {
		return r.RemoteAddr
	}

	return "0.0.0.0"
}

func cleanPath(rawPath string) string {
	if rawPath == "" {
		return "/"
	}

	cleaned := filepath.ToSlash(filepath.Clean("/" + rawPath))
	if cleaned == "." {
		return "/"
	}

	return cleaned
}

func splitSegments(requestPath string) []string {
	if requestPath == "/" || requestPath == "" {
		return nil
	}

	parts := strings.Split(strings.Trim(requestPath, "/"), "/")
	result := make([]string, 0, len(parts))
	for _, part := range parts {
		if part != "" {
			result = append(result, part)
		}
	}

	return result
}

func containsDotfile(requestPath string) bool {
	for _, segment := range splitSegments(requestPath) {
		if strings.HasPrefix(segment, ".") {
			return true
		}
	}

	return false
}

func closeConnection(w http.ResponseWriter) {
	if hijacker, ok := w.(http.Hijacker); ok {
		conn, _, err := hijacker.Hijack()
		if err == nil {
			_ = conn.Close()
			return
		}
	}
	panic(http.ErrAbortHandler)
}

func cookiePathForPrefix(prefix string) string {
	trimmed := strings.Trim(prefix, "/")
	if trimmed == "" {
		return "/"
	}

	return "/" + trimmed + "/"
}

func prefixAssetsInHTML(content, prefix string) string {
	trimmed := strings.Trim(prefix, "/")
	if trimmed == "" {
		return content
	}

	assetPrefix := "/" + trimmed + "/"
	replacer := strings.NewReplacer(
		`"/assets/`, `"`+assetPrefix+`assets/`,
		`'/assets/`, `'`+assetPrefix+`assets/`,
		`(/assets/`, `(`+assetPrefix+`assets/`,
		`"/locales/`, `"`+assetPrefix+`locales/`,
		`'/locales/`, `'`+assetPrefix+`locales/`,
		`(/locales/`, `(`+assetPrefix+`locales/`,
	)

	return replacer.Replace(content)
}
