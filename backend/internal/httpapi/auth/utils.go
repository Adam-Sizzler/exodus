package auth

import (
	"context"
	"net"
	"net/http"
	"path"
	"strings"
	"sync"
	"time"

	"exodus/internal/config"
	"exodus/internal/httpapi/middleware"
	"exodus/internal/jobqueue"
	"exodus/internal/notifications"
	"exodus/internal/security"

	"github.com/redis/go-redis/v9"
)

// --- Login brute-force rate limiting -------------------------------------
const (
	AuthRateLimitWindow      = 60 * time.Second
	AuthRateLimitMaxAttempts = 5
	AuthFailedLoginDelay     = 5 * time.Second
)

type AuthRateLimiter struct {
	mu        sync.Mutex
	attempts  map[string][]time.Time
	redis     *redis.Client
	redisOnce sync.Once
}

var globalAuthRateLimiter = &AuthRateLimiter{
	attempts: make(map[string][]time.Time),
}

func (l *AuthRateLimiter) getRedis(cfg *config.BackendConfig) *redis.Client {
	l.redisOnce.Do(func() {
		if cfg != nil {
			client, err := jobqueue.NewRedisClient(cfg)
			if err == nil {
				l.redis = client
			}
		}
	})
	return l.redis
}

func (l *AuthRateLimiter) Allow(ctx context.Context, ip string, cfg *config.BackendConfig) (allowed bool, remaining int, retryAfter time.Duration) {
	if ip == "" {
		return true, AuthRateLimitMaxAttempts, 0
	}

	redisClient := l.getRedis(cfg)
	if redisClient != nil {
		key := "ratelimit:auth:login:" + ip
		val, err := redisClient.Get(ctx, key).Int()
		if err == nil && val >= AuthRateLimitMaxAttempts {
			ttl, _ := redisClient.TTL(ctx, key).Result()
			if ttl <= 0 {
				ttl = AuthRateLimitWindow
			}
			return false, 0, ttl
		}
		if err == nil {
			rem := AuthRateLimitMaxAttempts - val
			if rem < 0 {
				rem = 0
			}
			return true, rem, 0
		}
	}

	// Fallback to in-memory sliding window
	l.mu.Lock()
	defer l.mu.Unlock()

	now := time.Now()
	cutoff := now.Add(-AuthRateLimitWindow)
	existing := l.attempts[ip]
	kept := existing[:0]
	for _, t := range existing {
		if t.After(cutoff) {
			kept = append(kept, t)
		}
	}
	l.attempts[ip] = kept

	if len(kept) >= AuthRateLimitMaxAttempts {
		oldest := kept[0]
		retry := AuthRateLimitWindow - now.Sub(oldest)
		if retry <= 0 {
			retry = time.Second
		}
		return false, 0, retry
	}

	remaining = AuthRateLimitMaxAttempts - len(kept)
	return true, remaining, 0
}

func (l *AuthRateLimiter) RecordFailedAttempt(ctx context.Context, ip string, cfg *config.BackendConfig) {
	if ip == "" {
		return
	}

	redisClient := l.getRedis(cfg)
	if redisClient != nil {
		key := "ratelimit:auth:login:" + ip
		pipe := redisClient.Pipeline()
		pipe.Incr(ctx, key)
		pipe.Expire(ctx, key, AuthRateLimitWindow)
		_, _ = pipe.Exec(ctx)
	}

	// Also record in-memory fallback
	l.mu.Lock()
	defer l.mu.Unlock()
	l.attempts[ip] = append(l.attempts[ip], time.Now())
}

func (l *AuthRateLimiter) Reset(ctx context.Context, ip string, cfg *config.BackendConfig) {
	if ip == "" {
		return
	}

	redisClient := l.getRedis(cfg)
	if redisClient != nil {
		key := "ratelimit:auth:login:" + ip
		_ = redisClient.Del(ctx, key).Err()
	}

	l.mu.Lock()
	defer l.mu.Unlock()
	delete(l.attempts, ip)
}

func loginRateLimitKey(r *http.Request, cfg *config.BackendConfig) string {
	ip := middleware.GetClientIP(r, cfg)
	if ip != "" {
		return ip
	}
	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil || host == "" {
		return r.RemoteAddr
	}
	return host
}

// --- Timing Safety -------------------------------------------------------
var (
	dummyPasswordHashOnce sync.Once
	dummyPasswordHashVal  string
)

func getDummyPasswordHash(secret string) string {
	dummyPasswordHashOnce.Do(func() {
		dummyPasswordHashVal = mustDummyPasswordHash(secret)
	})
	return dummyPasswordHashVal
}

func mustDummyPasswordHash(secret string) string {
	hash, err := security.HashPassword("exodus-timing-safety-placeholder-do-not-use", secret)
	if err != nil {
		return "00000000000000000000000000000000:" +
			"00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000"
	}
	return hash
}

// --- Context Principal ---------------------------------------------------
func CurrentAuthPrincipal(ctx context.Context) (*AuthPrincipal, bool) {
	if ctx == nil {
		return nil, false
	}
	val := ctx.Value(authPrincipalContextKey)
	if val == nil {
		return nil, false
	}
	principal, ok := val.(*AuthPrincipal)
	if !ok || principal == nil {
		return nil, false
	}
	return principal, true
}

// --- API Scopes ----------------------------------------------------------
func requireAPITokenScope(principal *AuthPrincipal, r *http.Request) bool {
	if principal == nil {
		return false
	}
	if principal.TokenType != "jwt_api_token" {
		return true // UI sessions bypass API scope validation
	}
	if len(principal.Scopes) == 0 {
		return false
	}

	resource, expectedScope, kind := apiTokenScopeForRequest(r.Method, r.URL.Path)
	if resource == "" || expectedScope == "" {
		return false
	}

	if hasScope(principal.Scopes, expectedScope) {
		return true
	}

	if kind == "item" && strings.HasSuffix(expectedScope, ":write") {
		wildcardScope := resource + ":write"
		if hasScope(principal.Scopes, wildcardScope) {
			return true
		}
	}
	if kind == "item" && strings.HasSuffix(expectedScope, ":read") {
		wildcardScope := resource + ":read"
		if hasScope(principal.Scopes, wildcardScope) {
			return true
		}
	}

	return false
}

func normalizeAPITokenPrincipalScopes(scopes []string) []string {
	seen := make(map[string]struct{}, len(scopes))
	result := make([]string, 0, len(scopes))
	for _, raw := range scopes {
		scope := strings.ToLower(strings.TrimSpace(raw))
		if scope == "" {
			continue
		}
		if _, ok := seen[scope]; ok {
			continue
		}
		seen[scope] = struct{}{}
		result = append(result, scope)
	}
	return result
}

func hasScope(scopes []string, expected string) bool {
	for _, scope := range scopes {
		if strings.EqualFold(scope, expected) {
			return true
		}
	}
	return false
}

func apiTokenScopeForRequest(method, requestPath string) (resource, endpointScope, kind string) {
	clean := cleanPath(requestPath)
	if !strings.HasPrefix(clean, "/api/") {
		return "", "", ""
	}
	trimmed := clean[len("/api/"):]

	switch {
	case strings.HasPrefix(trimmed, "hosts/") || trimmed == "hosts":
		return genericCRUDScope("hosts", method, clean, "/api/hosts")
	case strings.HasPrefix(trimmed, "internal-squads/") || trimmed == "internal-squads":
		return genericCRUDScope("internal_squads", method, clean, "/api/internal-squads")
	case strings.HasPrefix(trimmed, "subscription-templates/") || trimmed == "subscription-templates":
		return genericCRUDScope("subscription_templates", method, clean, "/api/subscription-templates")
	case strings.HasPrefix(trimmed, "system/metadata") || trimmed == "system/metadata":
		return genericCRUDScope("system_metadata", method, clean, "/api/system/metadata")
	case strings.HasPrefix(trimmed, "nodes/") || trimmed == "nodes":
		return genericCRUDScope("nodes", method, clean, "/api/nodes")
	case strings.HasPrefix(trimmed, "config-profiles/") || trimmed == "config-profiles":
		return genericCRUDScope("config_profiles", method, clean, "/api/config-profiles")
	case strings.HasPrefix(trimmed, "subscription-connections/") || trimmed == "subscription-connections":
		return genericCRUDScope("subscription_connections", method, clean, "/api/subscription-connections")
	case strings.HasPrefix(trimmed, "squads/") || trimmed == "squads":
		return genericCRUDScope("squads", method, clean, "/api/squads")
	default:
		return "", "", ""
	}
}

func crudEndpointScope(resource, method, kind string) (string, string, string) {
	switch method {
	case http.MethodGet, http.MethodHead:
		return resource, resource + ":read", kind
	case http.MethodPost, http.MethodPut, http.MethodPatch, http.MethodDelete:
		return resource, resource + ":write", kind
	default:
		return "", "", ""
	}
}

func genericCRUDScope(resource, method, requestPath, collectionPath string) (string, string, string) {
	if requestPath == collectionPath || requestPath == collectionPath+"/" {
		return crudEndpointScope(resource, method, "collection")
	}
	return crudEndpointScope(resource, method, "item")
}

// --- Notifications & Emitters ---------------------------------------------
func emitLoginNotification(ctx context.Context, cfg *config.BackendConfig, event, username, adminUUID, password, reason string, r *http.Request) {
	var (
		clientIP  = notificationClientIP(cfg, r)
		userAgent = strings.TrimSpace(r.Header.Get("User-Agent"))
	)

	notifications.Emit(ctx, cfg, notifications.Event{
		Scope: notifications.ScopeService,
		Event: event,
		Data: map[string]any{
			"username":  username,
			"adminUuid": adminUUID,
			"password":  password,
			"ip":        clientIP,
			"userAgent": userAgent,
			"reason":    reason,
		},
	})
}

func notificationClientIP(cfg *config.BackendConfig, r *http.Request) string {
	return middleware.GetClientIP(r, cfg)
}

func cleanPath(rawPath string) string {
	cleaned := path.Clean(strings.TrimSpace(rawPath))
	if cleaned == "." {
		return "/"
	}
	return cleaned
}

func readBearerToken(authHeader string) string {
	trimmed := strings.TrimSpace(authHeader)
	if trimmed == "" {
		return ""
	}

	parts := strings.SplitN(trimmed, " ", 2)
	if len(parts) == 2 && strings.EqualFold(parts[0], "bearer") {
		return strings.TrimSpace(parts[1])
	}

	return ""
}

func AuthTTLMinutes() int {
	return 2880 // 48 hours
}

func nullableString(value string) *string {
	trimmed := strings.TrimSpace(value)
	if trimmed == "" {
		return nil
	}
	return &trimmed
}
