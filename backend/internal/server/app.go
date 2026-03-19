package server

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/cerberus/subscription-page/backend/internal/assets"
	"github.com/cerberus/subscription-page/backend/internal/config"
	"github.com/cerberus/subscription-page/backend/internal/cerberus"
	"github.com/cerberus/subscription-page/backend/internal/security"
	"github.com/cerberus/subscription-page/backend/internal/subpages"
)

const sessionCookieName = "session"

var (
	appConfigPaths = map[string]struct{}{
		"/assets/.app-config-v2.json": {},
	}

	allowedClientTypes = map[string]struct{}{
		"stash":      {},
		"singbox":    {},
		"mihomo":     {},
		"json":       {},
		"v2ray-json": {},
		"clash":      {},
	}

	browserKeywords = []string{
		"Mozilla",
		"Chrome",
		"Safari",
		"Firefox",
		"Opera",
		"Edge",
		"TelegramBot",
		"WhatsApp",
	}

	genericPathFragments = []string{
		"favicon.ico",
		"robots.txt",
		".png",
		".jpg",
		".jpeg",
		".gif",
		".svg",
		".webp",
		".ico",
	}
)

type App struct {
	cfg        config.Config
	client     *cerberus.Client
	assetsPath string
	subpages   *subpages.Store
}

func New(cfg config.Config) (*App, error) {
	client := cerberus.NewClient(cfg)

	version, err := client.GetMetadata(context.Background())
	if err != nil {
		return nil, fmt.Errorf("connection to Cerberus Panel failed: %w", err)
	}

	log.Printf("[OK] Connected to Cerberus v%s", version)

	assetsPath, err := assets.DetectPath()
	if err != nil {
		return nil, err
	}

	store := subpages.NewStore(cfg)
	if err := store.Load(context.Background(), client); err != nil {
		return nil, err
	}

	log.Printf("[CONFIG] assets path: %s", assetsPath)

	return &App{
		cfg:        cfg,
		client:     client,
		assetsPath: assetsPath,
		subpages:   store,
	}, nil
}

func (a *App) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("X-Robots-Tag", "noindex, nofollow, noarchive, nosnippet, noimageindex")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET")
	w.Header().Set("Referrer-Policy", "no-referrer")
	w.Header().Set("X-Content-Type-Options", "nosniff")

	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	if r.Method != http.MethodGet && r.Method != http.MethodHead {
		closeConnection(w)
		return
	}

	if !a.cfg.IsDevelopment() {
		if strings.TrimSpace(r.Header.Get("X-Forwarded-For")) == "" ||
			r.Header.Get("X-Forwarded-Proto") != "https" {
			log.Printf(
				"[ERROR] Reverse proxy and HTTPS are required. X-Forwarded-For=%q X-Forwarded-Proto=%q",
				r.Header.Get("X-Forwarded-For"),
				r.Header.Get("X-Forwarded-Proto"),
			)
			closeConnection(w)
			return
		}
	}

	requestPath := cleanPath(r.URL.Path)
	routePath, ok := a.applyCustomPrefix(requestPath)
	if !ok {
		closeConnection(w)
		return
	}

	if a.isAppConfigPath(routePath) {
		a.handleAppConfig(w, r)
		return
	}

	if strings.HasPrefix(routePath, "/assets") || strings.HasPrefix(routePath, "/locales") {
		if _, err := a.verifySessionCookie(r); err != nil {
			log.Printf("[DEBUG] static asset session verification failed: %v", err)
			closeConnection(w)
			return
		}

		a.serveStatic(w, r, routePath)
		return
	}

	segments := splitSegments(routePath)
	if len(segments) != 1 && len(segments) != 2 {
		closeConnection(w)
		return
	}

	shortUUID := segments[0]
	clientType := ""
	if len(segments) == 2 {
		clientType = segments[1]
		if _, ok := allowedClientTypes[clientType]; !ok {
			log.Printf("[ERROR] invalid client type: %s", clientType)
			closeConnection(w)
			return
		}
	}

	if a.isGenericPath(routePath) {
		closeConnection(w)
		return
	}

	clientIP := getRealIP(r)
	resolvedShortUUID, err := a.resolveShortUUID(r.Context(), clientIP, shortUUID)
	if err != nil {
		log.Printf("[DEBUG] short uuid resolution failed for %s: %v", shortUUID, err)
		closeConnection(w)
		return
	}

	if a.isBrowser(r.UserAgent()) {
		a.returnWebpage(clientIP, resolvedShortUUID, w, r)
		return
	}

	a.proxySubscription(clientIP, resolvedShortUUID, clientType, w, r)
}

func (a *App) handleAppConfig(w http.ResponseWriter, r *http.Request) {
	claims, err := a.verifySessionCookie(r)
	if err != nil {
		log.Printf("[DEBUG] app config session verification failed: %v", err)
		closeConnection(w)
		return
	}

	encryptedUUID, _ := claims["su"].(string)
	if encryptedUUID == "" {
		closeConnection(w)
		return
	}

	subpageConfig, err := a.subpages.ResolveRawConfig(encryptedUUID)
	if err != nil {
		log.Printf("[ERROR] %v", err)
		closeConnection(w)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(subpageConfig)
}

func (a *App) serveStatic(w http.ResponseWriter, r *http.Request, requestPath string) {
	if containsDotfile(requestPath) {
		closeConnection(w)
		return
	}

	trimmedPath := strings.TrimPrefix(requestPath, "/")
	fullPath := filepath.Join(a.assetsPath, filepath.FromSlash(trimmedPath))
	cleanAssetsRoot := filepath.Clean(a.assetsPath)
	cleanFullPath := filepath.Clean(fullPath)

	if cleanFullPath != cleanAssetsRoot &&
		!strings.HasPrefix(cleanFullPath, cleanAssetsRoot+string(filepath.Separator)) {
		closeConnection(w)
		return
	}

	info, err := os.Stat(cleanFullPath)
	if err != nil || info.IsDir() {
		closeConnection(w)
		return
	}

	// Frontend bundle requests app-config from "/assets/.app-config-v2.json".
	// Under prefix mode we rewrite it to "/<prefix>/assets/.app-config-v2.json"
	// so Nginx with only "location /sub/" can serve it correctly.
	if strings.Trim(a.cfg.CustomSubPrefix, "/") != "" && strings.HasSuffix(cleanFullPath, ".js") {
		content, readErr := os.ReadFile(cleanFullPath)
		if readErr == nil {
			prefix := "/" + strings.Trim(a.cfg.CustomSubPrefix, "/")
			rewritten := strings.ReplaceAll(
				string(content),
				`"/assets/.app-config-v2.json"`,
				fmt.Sprintf("%q", prefix+"/assets/.app-config-v2.json"),
			)
			w.Header().Set("Content-Type", "application/javascript; charset=utf-8")
			if r.Method != http.MethodHead {
				_, _ = io.WriteString(w, rewritten)
			}
			return
		}
	}

	http.ServeFile(w, r, cleanFullPath)
}

func (a *App) proxySubscription(
	clientIP, shortUUID, clientType string,
	w http.ResponseWriter,
	r *http.Request,
) {
	body, headers, err := a.client.GetSubscription(
		r.Context(),
		clientIP,
		shortUUID,
		clientType,
		r.Header,
	)
	if err != nil {
		log.Printf("[ERROR] get subscription failed for %s: %v", shortUUID, err)
		closeConnection(w)
		return
	}

	for key, values := range headers {
		for _, value := range values {
			w.Header().Add(key, value)
		}
	}

	w.WriteHeader(http.StatusOK)
	if r.Method != http.MethodHead {
		_, _ = w.Write(body)
	}
}

func (a *App) returnWebpage(clientIP, shortUUID string, w http.ResponseWriter, r *http.Request) {
	subscriptionData, err := a.client.GetSubscriptionInfo(r.Context(), clientIP, shortUUID)
	if err != nil {
		log.Printf("[ERROR] get subscription info failed for %s: %v", shortUUID, err)
		closeConnection(w)
		return
	}

	subpageResponse, err := a.client.GetSubpageConfig(r.Context(), shortUUID, r.Header)
	if err != nil {
		log.Printf("[ERROR] get subpage config failed for %s: %v", shortUUID, err)
		closeConnection(w)
		return
	}

	if !subpageResponse.WebpageAllowed {
		log.Printf("Webpage access is not allowed by Cerberus's SRR.")
		a.proxySubscription(clientIP, shortUUID, "", w, r)
		return
	}

	settings := a.subpages.BaseSettingsFor(subpageResponse.SubpageConfigUUID)
	if !settings.ShowConnectionKeys {
		hideConnectionKeys(subscriptionData)
	}

	encryptedUUID, err := a.subpages.EncryptResolvedUUID(subpageResponse.SubpageConfigUUID)
	if err != nil {
		log.Printf("[ERROR] failed to encrypt subpage config uuid: %v", err)
		closeConnection(w)
		return
	}

	sessionToken, err := security.SignJWT(security.Claims{
		"sessionId": security.RandomToken(32),
		"su":        encryptedUUID,
		"exp":       time.Now().Add(33 * time.Minute).Unix(),
	}, a.cfg.SessionSecret)
	if err != nil {
		log.Printf("[ERROR] failed to build session jwt: %v", err)
		closeConnection(w)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     sessionCookieName,
		Value:    sessionToken,
		HttpOnly: true,
		Secure:   true,
		// Use root cookie path to support frontends that request app-config/assets
		// from both "/assets/*" and "/<prefix>/assets/*".
		Path:   "/",
		MaxAge: 1800,
	})

	panelData, err := json.Marshal(subscriptionData)
	if err != nil {
		log.Printf("[ERROR] failed to marshal subscription info: %v", err)
		closeConnection(w)
		return
	}

	indexHTML, err := os.ReadFile(filepath.Join(a.assetsPath, "index.html"))
	if err != nil {
		log.Printf("[ERROR] failed to read frontend index.html: %v", err)
		closeConnection(w)
		return
	}

	rendered := renderIndexTemplate(
		string(indexHTML),
		settings.MetaTitle,
		settings.MetaDescription,
		base64.StdEncoding.EncodeToString(panelData),
	)
	rendered = prefixAssetsInHTML(rendered, a.cfg.CustomSubPrefix)
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	if r.Method != http.MethodHead {
		_, _ = io.WriteString(w, rendered)
	}
}

func (a *App) verifySessionCookie(r *http.Request) (security.Claims, error) {
	cookie, err := r.Cookie(sessionCookieName)
	if err != nil {
		return nil, err
	}

	return security.VerifyJWT(cookie.Value, a.cfg.SessionSecret)
}

func (a *App) applyCustomPrefix(requestPath string) (string, bool) {
	if a.cfg.CustomSubPrefix == "" {
		return requestPath, true
	}

	// Compatibility mode: some frontend bundles request these routes from root
	// even when CUSTOM_SUB_PREFIX is enabled.
	if requestPath == "/assets" ||
		requestPath == "/locales" ||
		strings.HasPrefix(requestPath, "/assets/") ||
		strings.HasPrefix(requestPath, "/locales/") {
		return requestPath, true
	}

	prefix := "/" + a.cfg.CustomSubPrefix
	if requestPath == prefix {
		return "/", true
	}

	if !strings.HasPrefix(requestPath, prefix+"/") {
		// Compatibility mode: allow direct short UUID routes without prefix.
		segments := splitSegments(requestPath)
		if len(segments) == 1 || len(segments) == 2 {
			return requestPath, true
		}
		return "", false
	}

	return strings.TrimPrefix(requestPath, prefix), true
}

func (a *App) isAppConfigPath(requestPath string) bool {
	_, ok := appConfigPaths[requestPath]
	return ok
}

func (a *App) isBrowser(userAgent string) bool {
	for _, keyword := range browserKeywords {
		if strings.Contains(userAgent, keyword) {
			return true
		}
	}

	return false
}

func (a *App) isGenericPath(requestPath string) bool {
	for _, fragment := range genericPathFragments {
		if strings.Contains(requestPath, fragment) {
			return true
		}
	}

	return false
}

func (a *App) resolveShortUUID(ctx context.Context, clientIP, shortUUID string) (string, error) {
	return shortUUID, nil
}

func hideConnectionKeys(subscriptionData map[string]any) {
	responseValue, ok := subscriptionData["response"]
	if !ok {
		return
	}

	responseMap, ok := responseValue.(map[string]any)
	if !ok {
		return
	}

	responseMap["links"] = []any{}
	responseMap["ssConfLinks"] = map[string]any{}
}
