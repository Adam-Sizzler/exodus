package httpapi

import (
	"net/http"

	"v2ray-stat/backend/panel/config"
	dbmanager "v2ray-stat/backend/panel/db/manager"
	"v2ray-stat/backend/panel/httpapi/auth"
	"v2ray-stat/backend/panel/httpapi/configprofiles"
	"v2ray-stat/backend/panel/httpapi/health"
	"v2ray-stat/backend/panel/httpapi/hosts"
	"v2ray-stat/backend/panel/httpapi/middleware"
	"v2ray-stat/backend/panel/httpapi/nodes"
	"v2ray-stat/backend/panel/httpapi/panelsettings"
	"v2ray-stat/backend/panel/httpapi/squads"
	"v2ray-stat/backend/panel/httpapi/subscription"
	"v2ray-stat/backend/panel/httpapi/templates"
	"v2ray-stat/backend/panel/httpapi/users"
)

func NewAPIHandler(manager *dbmanager.DatabaseManager, cfg *config.BackendConfig) http.Handler {
	mux := http.NewServeMux()
	RegisterRoutes(mux, manager, cfg)
	return middleware.WithCORS(cfg, auth.WithPanelAuth(manager, cfg, mux))
}

func RegisterRoutes(mux *http.ServeMux, manager *dbmanager.DatabaseManager, cfg *config.BackendConfig) {
	mux.HandleFunc("/api/v2/auth/bootstrap", auth.AuthBootstrapHandler(manager, cfg))
	mux.HandleFunc("/api/v2/auth/setup", auth.AuthSetupHandler(manager, cfg))
	mux.HandleFunc("/api/v2/auth/login", auth.AuthLoginHandler(manager, cfg))
	mux.HandleFunc("/api/v2/auth/logout", auth.AuthLogoutHandler(manager, cfg))
	mux.HandleFunc("/api/v2/auth/me", auth.AuthMeHandler(manager, cfg))

	mux.HandleFunc("/api/v2/panel/settings", panelsettings.PanelSettingsHandler(manager, cfg))
	mux.HandleFunc("/api/v2/panel/api-tokens", panelsettings.PanelAPITokensHandler(manager, cfg))
	mux.HandleFunc("/api/v2/panel/api-tokens/", panelsettings.PanelAPITokenByUUIDHandler(manager, cfg))

	mux.HandleFunc("/api/v2/nodes", nodes.NodesHandler(manager, cfg))
	mux.HandleFunc("/api/v2/nodes/", nodes.NodeByUUIDHandler(manager, cfg))
	mux.HandleFunc("/api/v2/nodes/summary", nodes.NodesSummaryHandler(manager, cfg))
	mux.HandleFunc("/api/v2/nodes/reorder", nodes.NodesReorderHandler(manager, cfg))
	mux.HandleFunc("/api/v2/nodes-with-config", squads.NodesWithConfigHandler(manager, cfg))

	mux.HandleFunc("/api/v2/hosts", hosts.HostsHandler(manager, cfg))
	mux.HandleFunc("/api/v2/hosts/", hosts.HostByUUIDHandler(manager, cfg))
	mux.HandleFunc("/api/v2/hosts/reorder", hosts.HostsReorderHandler(manager, cfg))
	mux.HandleFunc("/api/v2/hosts-to-nodes", hosts.HostNodeAssignmentsHandler(manager, cfg))

	mux.HandleFunc("/api/v2/users-list", users.UsersAPIHandler(manager, cfg))
	mux.HandleFunc("/api/v2/users-list/", users.UserByUUIDHandler(manager, cfg))
	mux.HandleFunc("/api/v2/users-list/reorder", users.UsersReorderHandler(manager, cfg))
	mux.HandleFunc("/api/v2/users-list/summary", users.UsersListSummaryHandler(manager, cfg))
	mux.HandleFunc("/api/v2/users-list/create", users.UsersCreateHandler(manager, cfg))

	mux.HandleFunc("/api/v2/config-profiles", configprofiles.ConfigProfilesHandler(manager, cfg))
	mux.HandleFunc("/api/v2/config-profiles/", configprofiles.ConfigProfileByUUIDHandler(manager, cfg))
	mux.HandleFunc("/api/v2/config-profiles/reorder", configprofiles.ConfigProfilesReorderHandler(manager, cfg))
	mux.HandleFunc("/api/v2/config-profiles-with-inbounds", squads.ConfigProfilesWithInboundsHandler(manager, cfg))

	mux.HandleFunc("/api/v2/internal-squads", squads.InternalSquadsHandler(manager, cfg))
	mux.HandleFunc("/api/v2/internal-squads/", squads.InternalSquadByUUIDHandler(manager, cfg))
	mux.HandleFunc("/api/v2/internal-squads/reorder", squads.InternalSquadsReorderHandler(manager, cfg))
	mux.HandleFunc("/api/v2/squads-summary", squads.AllSquadsSummaryHandler(manager, cfg))
	mux.HandleFunc("/api/v2/squad-inbounds", squads.SquadInboundsHandler(manager, cfg))
	mux.HandleFunc("/api/v2/squad-members", squads.SquadMembersHandler(manager, cfg))
	mux.HandleFunc("/api/v2/squad-details/", squads.SquadDetailsHandler(manager, cfg))

	mux.HandleFunc("/api/v2/inbound-assignments", squads.InboundAssignmentsHandler(manager, cfg))
	mux.HandleFunc("/api/v2/inbounds-with-profiles", squads.InboundsWithProfilesHandler(manager, cfg))

	mux.HandleFunc("/api/v2/subscription-settings", subscription.SubscriptionSettingsHandler(manager, cfg))
	mux.HandleFunc("/api/v2/subscription-settings/", subscription.SubscriptionSettingsByUUIDHandler(manager, cfg))

	mux.HandleFunc("/api/v2/templates", templates.SubscriptionTemplatesHandler(manager, cfg))
	mux.HandleFunc("/api/v2/templates/", templates.SubscriptionTemplateByUUIDHandler(manager, cfg))
	mux.HandleFunc("/api/v2/templates/reorder", templates.SubscriptionTemplatesReorderHandler(manager, cfg))

	mux.HandleFunc("/api/v2/sub", subscription.SubscriptionPublicHandler(manager, cfg))
	mux.HandleFunc("/api/v2/sub/", subscription.SubscriptionPublicHandler(manager, cfg))

	mux.HandleFunc("/api/v2/subscriptions", subscription.SubscriptionsHandler(manager, cfg))
	mux.HandleFunc("/api/v2/subscriptions/", subscription.SubscriptionsByPathHandler(manager, cfg))

	mux.HandleFunc("/api/v2/subscription-page-configs", subscription.SubscriptionPageConfigsHandler(manager, cfg))
	mux.HandleFunc("/api/v2/subscription-page-configs/", subscription.SubscriptionPageConfigByUUIDHandler(manager, cfg))

	mux.HandleFunc("/api/v2/health", health.HealthHandler())

	mux.Handle("/api/", http.NotFoundHandler())
}
