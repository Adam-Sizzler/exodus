package exodus

import (
	"fmt"
	"html"
	"net/http"
	"os"
	"strings"
)

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

	injected := fmt.Sprintf(
		"<base href=\"%s\" />\n<script>window.__EXODUS_RUNTIME__={basePath:%q};</script>",
		html.EscapeString(basePathWithSlash),
		basePath,
	)

	page := string(indexBytes)
	// Normalize accidental ".//" asset prefixes from the built index template.
	page = strings.ReplaceAll(page, ".//", "./")
	if strings.Contains(page, "<head>") {
		page = strings.Replace(page, "<head>", "<head>\n"+injected, 1)
	} else {
		page = injected + page
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(page))
}
