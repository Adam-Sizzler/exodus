package httpapi

import (
	"database/sql"
	"net/http"
	"strings"

	"exodus/internal/config"
	"exodus/internal/db"
	"exodus/internal/httpapi/asynqmon"
	"exodus/internal/httpapi/auth"
	"exodus/internal/httpapi/bandwidthstats"
	"exodus/internal/httpapi/configprofiles"
	"exodus/internal/httpapi/externalsquads"
	"exodus/internal/httpapi/health"
	"exodus/internal/httpapi/hosts"
	"exodus/internal/httpapi/hwiduserdevices"
	"exodus/internal/httpapi/infrabilling"
	"exodus/internal/httpapi/keygen"
	"exodus/internal/httpapi/metadata"
	"exodus/internal/httpapi/middleware"
	"exodus/internal/httpapi/nodeplugins"
	"exodus/internal/httpapi/nodes"
	"exodus/internal/httpapi/panelsettings"
	"exodus/internal/httpapi/passkeys"
	"exodus/internal/httpapi/squads"
	"exodus/internal/httpapi/srslists"
	"exodus/internal/httpapi/subscription"
	subscriptionconnections "exodus/internal/httpapi/subscriptionconnections"
	"exodus/internal/httpapi/subscriptionpageconfigs"
	subscriptionrequesthistory "exodus/internal/httpapi/subscriptionrequesthistory"
	"exodus/internal/httpapi/subscriptionsettings"
	"exodus/internal/httpapi/subscriptiontemplate"
	"exodus/internal/httpapi/system"
	"exodus/internal/httpapi/users"
)

func NewAPIHandler(pools *db.Pools, cfg *config.BackendConfig) http.Handler {
	mainMux := http.NewServeMux()

	// 1. Public routes (unprotected) with optional auth parsing
	publicMux := http.NewServeMux()
	RegisterPublicRoutes(publicMux, pools.Interactive, pools.Background, cfg)
	publicHandler := auth.WithOptionalPanelAuth(pools.Interactive, cfg, publicMux)

	// 2. Protected routes with strict Auth enforcement
	protectedMux := http.NewServeMux()
	RegisterProtectedRoutes(protectedMux, pools.Interactive, pools.Background, cfg)
	protectedHandler := auth.WithPanelAuth(pools.Interactive, cfg, protectedMux)

	// Mount protected routes under /api/ first, then fall back to public routes / mainMux
	mainMux.Handle("/api/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check if request matches a public route
		if isPublicPath(r.URL.Path, cfg) {
			publicHandler.ServeHTTP(w, r)
		} else {
			protectedHandler.ServeHTTP(w, r)
		}
	}))

	// Non-/api/ public routes (e.g. /health)
	mainMux.HandleFunc("/health", health.HealthHandler())

	handler := middleware.WithRequestLogging(cfg, "api", mainMux)
	return middleware.WithCORS(cfg, handler)
}

func isPublicPath(path string, cfg *config.BackendConfig) bool {
	if strings.HasPrefix(path, "/api/subscriptions/connection-keys/") {
		return false
	}
	switch {
	case path == "/api/auth/status",
		path == "/api/auth/bootstrap",
		path == "/api/auth/setup",
		path == "/api/auth/login",
		path == "/api/auth/register",
		path == "/api/auth/oauth2/authorize",
		path == "/api/auth/oauth2/callback",
		path == "/api/auth/passkey/authentication/options",
		path == "/api/auth/passkey/authentication/verify",
		path == "/api/sub",
		strings.HasPrefix(path, "/api/sub/"),
		path == "/api/subscriptions/subpage-config",
		strings.HasPrefix(path, "/api/subscriptions/subpage-config/"),
		path == "/api/system/metadata",
		path == "/api/system/health",
		path == "/api/health":
		return true
	}
	if cfg != nil && cfg.Docs.IsEnabled {
		if path == cfg.Docs.ScalarPath || path == cfg.Docs.ScalarPath+"/openapi.json" ||
			path == cfg.Docs.SwaggerPath || path == cfg.Docs.SwaggerPath+"/openapi.json" {
			return true
		}
	}
	return false
}

func RegisterPublicRoutes(mux *http.ServeMux, db, backgroundDB *sql.DB, cfg *config.BackendConfig) {
	mux.HandleFunc("/api/auth/bootstrap", auth.AuthBootstrapHandler(db, cfg))
	mux.HandleFunc("/api/auth/setup", auth.AuthSetupHandler(db, cfg))
	mux.HandleFunc("/api/auth/status", auth.AuthStatusHandler(db, cfg))
	mux.HandleFunc("/api/auth/register", auth.AuthRegisterCompatHandler(db, cfg))
	mux.HandleFunc("/api/auth/login", auth.AuthLoginCompatHandler(db, cfg))
	mux.HandleFunc("/api/auth/oauth2/authorize", auth.OAuth2AuthorizeHandler(db, cfg))
	mux.HandleFunc("/api/auth/oauth2/callback", auth.OAuth2CallbackHandler(db, cfg))
	mux.HandleFunc("/api/auth/passkey/authentication/options", passkeys.AuthenticationOptionsHandler(db, cfg))
	mux.HandleFunc("/api/auth/passkey/authentication/verify", passkeys.VerifyAuthenticationHandler(db, cfg))

	// API docs (Scalar UI + Swagger UI + raw OpenAPI JSON). Only active when IS_DOCS_ENABLED=true.
	mux.HandleFunc(cfg.Docs.ScalarPath, panelsettings.DocsScalarHandler(cfg))
	mux.HandleFunc(cfg.Docs.ScalarPath+"/openapi.json", panelsettings.DocsOpenAPIHandler(cfg))
	mux.HandleFunc(cfg.Docs.SwaggerPath, panelsettings.DocsSwaggerHandler(cfg))
	mux.HandleFunc(cfg.Docs.SwaggerPath+"/openapi.json", panelsettings.DocsOpenAPIHandler(cfg))

	mux.HandleFunc("/api/sub", subscription.SubscriptionPublicHandler(db, backgroundDB, cfg))
	mux.HandleFunc("/api/sub/", subscription.SubscriptionPublicHandler(db, backgroundDB, cfg))

	mux.HandleFunc("/api/subscriptions/subpage-config/", subscription.SubpageConfigPublicHandler(db, backgroundDB, cfg))
	mux.HandleFunc("/api/subscriptions/subpage-config", subscription.SubpageConfigPublicHandler(db, backgroundDB, cfg))

	mux.HandleFunc("/api/system/metadata", system.MetadataHandler(cfg))
	mux.HandleFunc("/api/system/health", system.HealthHandler(cfg))
	mux.HandleFunc("/api/health", health.HealthHandler())
}

func RegisterProtectedRoutes(mux *http.ServeMux, db, backgroundDB *sql.DB, cfg *config.BackendConfig) {
	mux.HandleFunc("/api/auth/logout", auth.AuthLogoutHandler(db, cfg))
	mux.HandleFunc("/api/auth/me", auth.AuthMeHandler(db, cfg))

	mux.HandleFunc("/api/exodus-settings", auth.RequireAdminRole(panelsettings.ExodusSettingsHandler(db, cfg)))
	mux.HandleFunc("/api/exodus-settings/", auth.RequireAdminRole(panelsettings.ExodusSettingsHandler(db, cfg)))
	mux.HandleFunc("/api/tokens/scopes", auth.RequireAdminRole(panelsettings.PanelAPITokenScopesHandler(db, cfg)))
	mux.HandleFunc("/api/tokens", auth.RequireAdminRole(panelsettings.PanelAPITokensHandler(db, cfg)))
	mux.HandleFunc("/api/tokens/", auth.RequireAdminRole(panelsettings.PanelAPITokenByUUIDHandler(db, cfg)))

	mux.HandleFunc("/api/nodes", nodes.NodesHandler(db, cfg))
	mux.HandleFunc("/api/nodes/", nodes.NodeByUUIDHandler(db, cfg))
	mux.HandleFunc("/api/nodes/actions/", nodes.NodesActionsHandler(db, cfg))
	mux.HandleFunc("/api/nodes/bulk-actions", nodes.NodesBulkActionsHandler(db, cfg))
	mux.HandleFunc("/api/nodes/bulk-actions/", nodes.NodesBulkActionsHandler(db, cfg))
	mux.HandleFunc("/api/nodes/tags", nodes.NodesTagsHandler(db, cfg))
	mux.HandleFunc("/api/nodes-with-config", squads.NodesWithConfigHandler(db, cfg))
	mux.HandleFunc("/api/node-plugins", nodeplugins.Handler(db, cfg))
	mux.HandleFunc("/api/node-plugins/", nodeplugins.Handler(db, cfg))
	mux.HandleFunc("/api/metadata/user/", auth.RequireAdminRole(metadata.UserHandler(db, cfg)))
	mux.HandleFunc("/api/metadata/node/", auth.RequireAdminRole(metadata.NodeHandler(db, cfg)))

	mux.HandleFunc("/api/subscription-connections", subscriptionconnections.NodesHandler(db, cfg))
	mux.HandleFunc("/api/subscription-connections/", subscriptionconnections.NodeByUUIDHandler(db, cfg))
	mux.HandleFunc("/api/subscription-connections/actions/", subscriptionconnections.NodesActionsHandler(db, cfg))
	mux.HandleFunc("/api/subscription-connections/bulk-actions", subscriptionconnections.NodesBulkActionsHandler(db, cfg))
	mux.HandleFunc("/api/subscription-connections/bulk-actions/", subscriptionconnections.NodesBulkActionsHandler(db, cfg))
	mux.HandleFunc("/api/subscription-connections/tags", subscriptionconnections.NodesTagsHandler(db, cfg))

	mux.HandleFunc("/api/hosts", hosts.HostsHandler(db, cfg))
	mux.HandleFunc("/api/hosts/", hosts.HostByUUIDHandler(db, cfg))
	mux.HandleFunc("/api/hosts/actions/", hosts.HostsActionsHandler(db, cfg))
	mux.HandleFunc("/api/hosts/bulk/", hosts.HostsBulkHandler(db, cfg))
	mux.HandleFunc("/api/hosts/tags", hosts.HostsTagsHandler(db, cfg))

	mux.HandleFunc("/api/users", users.UsersHandler(db, cfg))
	mux.HandleFunc("/api/users/", users.UserByUUIDHandler(db, cfg))
	mux.HandleFunc("/api/users/bulk/", users.UsersBulkHandler(db, cfg))
	mux.HandleFunc("/api/users/tags", users.UsersTagsHandler(db, cfg))

	mux.HandleFunc("/api/keygen", keygen.KeygenHandler(db, cfg))
	mux.HandleFunc("/api/keygen/", keygen.KeygenHandler(db, cfg))

	if asynqmonHandler, err := asynqmon.NewAsynqmon(cfg); err == nil {
		mux.Handle("/api/queues/static/", asynqmonHandler)
		mux.Handle("/api/queues/", auth.RequireAdminRoleHandler(asynqmonHandler))
		mux.Handle("/api/queues", auth.RequireAdminRoleHandler(asynqmonHandler))
	}
	mux.HandleFunc("/api/passkeys/registration/options", passkeys.RegistrationOptionsHandler(db, cfg))
	mux.HandleFunc("/api/passkeys/registration/verify", passkeys.VerifyRegistrationHandler(db, cfg))
	mux.HandleFunc("/api/passkeys", passkeys.PasskeysHandler(db, cfg))
	mux.HandleFunc("/api/passkeys/", passkeys.PasskeysHandler(db, cfg))

	mux.HandleFunc("/api/bandwidth-stats/nodes", bandwidthstats.NodesHandler(db, cfg))
	mux.HandleFunc("/api/bandwidth-stats/nodes/", bandwidthstats.NodesHandler(db, cfg))
	mux.HandleFunc("/api/bandwidth-stats/users", bandwidthstats.UsersHandler(db, cfg))
	mux.HandleFunc("/api/bandwidth-stats/users/", bandwidthstats.UsersHandler(db, cfg))

	mux.HandleFunc("/api/config-profiles", configprofiles.ConfigProfilesHandler(db, cfg))
	mux.HandleFunc("/api/config-profiles/", configprofiles.ConfigProfileByUUIDHandler(db, cfg))
	mux.HandleFunc("/api/config-profiles/actions/", configprofiles.ConfigProfilesActionsHandler(db, cfg))
	mux.HandleFunc("/api/config-profiles/inbounds", configprofiles.ConfigProfilesInboundsHandler(db, cfg))
	mux.HandleFunc("/api/config-profiles/snippets", configprofiles.ConfigProfileSnippetsHandler(db, cfg))
	mux.HandleFunc("/api/snippets", configprofiles.ConfigProfileSnippetsHandler(db, cfg))
	mux.HandleFunc("/api/snippets/", configprofiles.ConfigProfileSnippetsHandler(db, cfg))
	mux.HandleFunc("/api/config-profiles-with-inbounds", squads.ConfigProfilesWithInboundsHandler(db, cfg))

	mux.HandleFunc("/api/internal-squads", squads.InternalSquadsHandler(db, cfg))
	mux.HandleFunc("/api/internal-squads/", squads.InternalSquadByUUIDHandler(db, cfg))
	mux.HandleFunc("/api/internal-squads/actions/reorder", squads.InternalSquadsReorderHandler(db, cfg))
	mux.HandleFunc("/api/internal-squads/reorder", squads.InternalSquadsReorderHandler(db, cfg))
	mux.HandleFunc("/api/squads-summary", squads.AllSquadsSummaryHandler(db, cfg))
	mux.HandleFunc("/api/squad-inbounds", squads.SquadInboundsHandler(db, cfg))
	mux.HandleFunc("/api/squad-members", squads.SquadMembersHandler(db, cfg))
	mux.HandleFunc("/api/squad-details/", squads.SquadDetailsHandler(db, cfg))

	mux.HandleFunc("/api/inbound-assignments", squads.InboundAssignmentsHandler(db, cfg))
	mux.HandleFunc("/api/inbounds-with-profiles", squads.InboundsWithProfilesHandler(db, cfg))

	mux.HandleFunc("/api/external-squads", externalsquads.ExternalSquadsHandler(db, cfg))
	mux.HandleFunc("/api/external-squads/", externalsquads.ExternalSquadByUUIDHandler(db, cfg))
	mux.HandleFunc("/api/external-squads/actions/reorder", externalsquads.ExternalSquadsReorderHandler(db, cfg))
	mux.HandleFunc("/api/external-squads/reorder", externalsquads.ExternalSquadsReorderHandler(db, cfg))

	mux.HandleFunc("/api/srs-lists", srslists.SRSListsHandler(db, cfg))
	mux.HandleFunc("/api/srs-lists/", srslists.SRSListByUUIDHandler(db, cfg))
	mux.HandleFunc("/api/srs-lists/actions/", srslists.SRSListsActionsHandler(db, cfg))
	mux.HandleFunc("/api/srs-lists/bulk/", srslists.SRSListsBulkHandler(db, cfg))

	mux.HandleFunc("/api/hwid/devices/delete-all", hwiduserdevices.HWIDCompatDeleteAllUserDevicesHandler(db, cfg))
	mux.HandleFunc("/api/hwid/devices/delete", hwiduserdevices.HWIDCompatDevicesHandler(db, cfg))
	mux.HandleFunc("/api/hwid/devices", hwiduserdevices.HWIDCompatDevicesHandler(db, cfg))
	mux.HandleFunc("/api/hwid/devices/", hwiduserdevices.HWIDCompatDevicesHandler(db, cfg))
	mux.HandleFunc("/api/hwid/devices/stats", hwiduserdevices.HWIDCompatStatsHandler(db, cfg))
	mux.HandleFunc("/api/hwid/devices/top-users", hwiduserdevices.HWIDCompatTopUsersHandler(db, cfg))

	mux.HandleFunc("/api/subscription-settings", subscriptionsettings.SubscriptionSettingsHandler(db, cfg))
	mux.HandleFunc("/api/subscription-settings/", subscriptionsettings.SubscriptionSettingsHandler(db, cfg))

	mux.HandleFunc("/api/infra-billing/providers", infrabilling.ProvidersHandler(db, cfg))
	mux.HandleFunc("/api/infra-billing/providers/", infrabilling.ProvidersHandler(db, cfg))
	mux.HandleFunc("/api/infra-billing/nodes", infrabilling.BillingNodesHandler(db, cfg))
	mux.HandleFunc("/api/infra-billing/nodes/", infrabilling.BillingNodesHandler(db, cfg))
	mux.HandleFunc("/api/infra-billing/history", infrabilling.BillingHistoryHandler(db, cfg))
	mux.HandleFunc("/api/infra-billing/history/", infrabilling.BillingHistoryHandler(db, cfg))

	mux.HandleFunc("/api/subscription-templates", subscriptiontemplate.SubscriptionTemplatesHandler(db, cfg))
	mux.HandleFunc("/api/subscription-templates/", subscriptiontemplate.SubscriptionTemplateByUUIDHandler(db, cfg))
	mux.HandleFunc("/api/subscription-templates/actions/", subscriptiontemplate.SubscriptionTemplatesActionsHandler(db, cfg))

	mux.HandleFunc("/api/subscriptions/connection-keys/", subscription.SubscriptionByUUIDHandler(db, backgroundDB, cfg))
	mux.HandleFunc("/api/subscriptions", subscription.SubscriptionsHandler(db, backgroundDB, cfg))
	mux.HandleFunc("/api/subscriptions/", subscription.SubscriptionByUUIDHandler(db, backgroundDB, cfg))

	mux.HandleFunc("/api/subscription-page-configs", subscriptionpageconfigs.SubscriptionPageConfigsHandler(db, cfg))
	mux.HandleFunc("/api/subscription-page-configs/", subscriptionpageconfigs.SubscriptionPageConfigByUUIDHandler(db, cfg))
	mux.HandleFunc("/api/subscription-page-configs/actions/", subscriptionpageconfigs.SubscriptionPageConfigsActionsHandler(db, cfg))
	mux.HandleFunc("/api/subscription-request-history", subscriptionrequesthistory.SubscriptionRequestHistoryHandler(db, cfg))
	mux.HandleFunc("/api/subscription-request-history/", subscriptionrequesthistory.SubscriptionRequestHistoryHandler(db, cfg))
	mux.HandleFunc("/api/subscription-request-history/stats", subscriptionrequesthistory.SubscriptionRequestHistoryStatsHandler(db, cfg))

	mux.HandleFunc("/api/system/stats", system.StatsHandler(db, cfg))
	mux.HandleFunc("/api/system/stats/bandwidth", system.BandwidthStatsHandler(db, cfg))
	mux.HandleFunc("/api/system/stats/nodes", system.NodesStatsHandler(db, cfg))
	mux.HandleFunc("/api/system/stats/recap", system.RecapHandler(db, cfg))
	mux.HandleFunc("/api/system/nodes/metrics", system.NodesMetricsHandler(db, cfg))
	mux.HandleFunc("/api/system/testers/srr-matcher", system.TestSRRMatcherHandler(cfg))
	mux.HandleFunc("/api/system/tools/happ/encrypt", system.EncryptHappCryptoLinkHandler(cfg))
	mux.HandleFunc("/api/system/tools/x25519/generate", system.GenerateX25519Handler(cfg))

	mux.Handle("/api/", http.NotFoundHandler())
}

func RegisterRoutes(mux *http.ServeMux, db, backgroundDB *sql.DB, cfg *config.BackendConfig) {
	RegisterPublicRoutes(mux, db, backgroundDB, cfg)
	RegisterProtectedRoutes(mux, db, backgroundDB, cfg)
}
