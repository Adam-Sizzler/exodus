package exodus

import (
	"fmt"
	"html"
	"net/http"
	"os"
	urlpath "path"
	"path/filepath"
	"strings"
)

var panelStaticPathPrefixes = []string{
	"assets/",
	"favicons/",
	"lotties/",
	"splash_screens/",
}

var panelStaticPathFiles = map[string]struct{}{
	"site.webmanifest": {},
}

func serveAppConfigJS(w http.ResponseWriter, basePath string) {
	w.Header().Set("Content-Type", "application/javascript; charset=utf-8")
	w.Header().Set("Cache-Control", "no-store")
	_, _ = fmt.Fprintf(
		w,
		"window.__EXODUS_RUNTIME__={basePath:%q};\n",
		basePath,
	)
}

func servePanelIndex(w http.ResponseWriter, indexPath string, basePathWithSlash, basePath string) {
	indexBytes, err := os.ReadFile(indexPath)
	if err != nil {
		http.Error(w, "panel index not found", http.StatusNotFound)
		return
	}

	page := renderPanelIndex(string(indexBytes), basePathWithSlash, basePath)

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(page))
}

func renderPanelIndex(page string, basePathWithSlash, basePath string) string {
	injected := fmt.Sprintf(
		"<base href=\"%s\" />\n<script>window.__EXODUS_RUNTIME__={basePath:%q};</script>",
		html.EscapeString(basePathWithSlash),
		basePath,
	)

	// Normalize accidental ".//" asset prefixes from the built index template.
	page = strings.ReplaceAll(page, ".//", "./")
	page = rewritePanelStaticAssetPaths(page, basePath)
	if strings.Contains(page, "<head>") {
		page = strings.Replace(page, "<head>", "<head>\n"+injected, 1)
	} else {
		page = injected + page
	}

	return page
}

func rewritePanelStaticAssetPaths(page string, basePath string) string {
	basePath = strings.TrimSuffix(strings.TrimSpace(basePath), "/")
	if basePath == "" || basePath == "/" {
		return strings.ReplaceAll(page, "%BASE_URL%", "")
	}

	page = strings.ReplaceAll(page, "%BASE_URL%", basePath)

	replacements := make([]string, 0, (len(panelStaticPathPrefixes)+len(panelStaticPathFiles))*8)
	addPrefixReplacement := func(staticPath string) {
		target := basePath + "/" + staticPath
		replacements = append(replacements,
			`="/`+staticPath, `="`+target,
			`='/`+staticPath, `='`+target,
			`("/`+staticPath, `("`+target,
			`('/`+staticPath, `('`+target,
		)
	}

	for _, prefix := range panelStaticPathPrefixes {
		addPrefixReplacement(prefix)
	}
	for file := range panelStaticPathFiles {
		addPrefixReplacement(file)
	}

	return strings.NewReplacer(replacements...).Replace(page)
}

func servePanelStaticFile(w http.ResponseWriter, r *http.Request, staticFS http.Handler, uiDir, requestPath string) bool {
	cleanPath, ok := cleanPanelStaticPath(requestPath)
	if !ok {
		return false
	}

	targetPath := filepath.Join(uiDir, filepath.FromSlash(cleanPath))
	if info, err := os.Stat(targetPath); err == nil && !info.IsDir() {
		staticReq := r.Clone(r.Context())
		staticReq.URL.Path = "/" + cleanPath
		staticReq.URL.RawPath = ""
		staticFS.ServeHTTP(w, staticReq)
		return true
	}

	http.NotFound(w, r)
	return true
}

func cleanPanelStaticPath(requestPath string) (string, bool) {
	cleanPath := strings.TrimPrefix(urlpath.Clean("/"+strings.TrimPrefix(requestPath, "/")), "/")
	if cleanPath == "." || cleanPath == "" {
		return "", false
	}

	if _, ok := panelStaticPathFiles[cleanPath]; ok {
		return cleanPath, true
	}
	for _, prefix := range panelStaticPathPrefixes {
		if strings.HasPrefix(cleanPath, prefix) {
			return cleanPath, true
		}
	}

	return "", false
}
