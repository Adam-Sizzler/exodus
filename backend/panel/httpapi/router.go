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
	subscriptionpageconfigs "v2ray-stat/backend/panel/httpapi/subscription-page-configs"
	subscriptionsettings "v2ray-stat/backend/panel/httpapi/subscription-settings"
	subscriptiontemplate "v2ray-stat/backend/panel/httpapi/subscription-template"
	"v2ray-stat/backend/panel/httpapi/users"
)

func NewAPIHandler(manager *dbmanager.DatabaseManager, cfg *config.BackendConfig) http.Handler {
	mux := http.NewServeMux()
	RegisterRoutes(mux, manager, cfg)
	return middleware.WithCORS(cfg, auth.WithPanelAuth(manager, cfg, mux))
}

func RegisterRoutes(mux *http.ServeMux, manager *dbmanager.DatabaseManager, cfg *config.BackendConfig) {
	mux.HandleFunc("/api/v1/auth/bootstrap", auth.AuthBootstrapHandler(manager, cfg))
	mux.HandleFunc("/api/v1/auth/setup", auth.AuthSetupHandler(manager, cfg))
	mux.HandleFunc("/api/v1/auth/login", auth.AuthLoginHandler(manager, cfg))
	mux.HandleFunc("/api/v1/auth/logout", auth.AuthLogoutHandler(manager, cfg))
	mux.HandleFunc("/api/v1/auth/me", auth.AuthMeHandler(manager, cfg))

	mux.HandleFunc("/api/v1/panel/settings", panelsettings.PanelSettingsHandler(manager, cfg))
	mux.HandleFunc("/api/v1/panel/api-tokens", panelsettings.PanelAPITokensHandler(manager, cfg))
	mux.HandleFunc("/api/v1/panel/api-tokens/", panelsettings.PanelAPITokenByUUIDHandler(manager, cfg))

	mux.HandleFunc("/api/v1/nodes", nodes.NodesHandler(manager, cfg))
	mux.HandleFunc("/api/v1/nodes/", nodes.NodeByUUIDHandler(manager, cfg))
	mux.HandleFunc("/api/v1/nodes/summary", nodes.NodesSummaryHandler(manager, cfg))
	mux.HandleFunc("/api/v1/nodes/reorder", nodes.NodesReorderHandler(manager, cfg))
	mux.HandleFunc("/api/v1/nodes-with-config", squads.NodesWithConfigHandler(manager, cfg))

	mux.HandleFunc("/api/hosts", hosts.HostsHandler(manager, cfg))
	mux.HandleFunc("/api/hosts/", hosts.HostByUUIDHandler(manager, cfg))
	mux.HandleFunc("/api/hosts/actions/", hosts.HostsActionsHandler(manager, cfg))
	mux.HandleFunc("/api/hosts/bulk/", hosts.HostsBulkHandler(manager, cfg))
	mux.HandleFunc("/api/hosts/tags", hosts.HostsTagsHandler(manager, cfg))

	mux.HandleFunc("/api/v1/users-list", users.UsersAPIHandler(manager, cfg))
	mux.HandleFunc("/api/v1/users-list/", users.UserByUUIDHandler(manager, cfg))
	mux.HandleFunc("/api/v1/users-list/reorder", users.UsersReorderHandler(manager, cfg))
	mux.HandleFunc("/api/v1/users-list/summary", users.UsersListSummaryHandler(manager, cfg))
	mux.HandleFunc("/api/v1/users-list/create", users.UsersCreateHandler(manager, cfg))

	mux.HandleFunc("/api/v1/config-profiles", configprofiles.ConfigProfilesHandler(manager, cfg))
	mux.HandleFunc("/api/v1/config-profiles/", configprofiles.ConfigProfileByUUIDHandler(manager, cfg))
	mux.HandleFunc("/api/v1/config-profiles/reorder", configprofiles.ConfigProfilesReorderHandler(manager, cfg))
	mux.HandleFunc("/api/v1/config-profiles-with-inbounds", squads.ConfigProfilesWithInboundsHandler(manager, cfg))

	mux.HandleFunc("/api/v1/internal-squads", squads.InternalSquadsHandler(manager, cfg))
	mux.HandleFunc("/api/v1/internal-squads/", squads.InternalSquadByUUIDHandler(manager, cfg))
	mux.HandleFunc("/api/v1/internal-squads/reorder", squads.InternalSquadsReorderHandler(manager, cfg))
	mux.HandleFunc("/api/v1/squads-summary", squads.AllSquadsSummaryHandler(manager, cfg))
	mux.HandleFunc("/api/v1/squad-inbounds", squads.SquadInboundsHandler(manager, cfg))
	mux.HandleFunc("/api/v1/squad-members", squads.SquadMembersHandler(manager, cfg))
	mux.HandleFunc("/api/v1/squad-details/", squads.SquadDetailsHandler(manager, cfg))

	mux.HandleFunc("/api/v1/inbound-assignments", squads.InboundAssignmentsHandler(manager, cfg))
	mux.HandleFunc("/api/v1/inbounds-with-profiles", squads.InboundsWithProfilesHandler(manager, cfg))

	mux.HandleFunc("/api/subscription-settings", subscriptionsettings.SubscriptionSettingsHandler(manager, cfg))
	// mux.HandleFunc("/api/v1/subscription-settings/", subscriptionsettings.SubscriptionSettingsByUUIDHandler(manager, cfg))

	mux.HandleFunc("/api/subscription-templates", subscriptiontemplate.SubscriptionTemplatesHandler(manager, cfg))
	mux.HandleFunc("/api/subscription-templates/", subscriptiontemplate.SubscriptionTemplateByUUIDHandler(manager, cfg))
	mux.HandleFunc("/api/subscription-templates/actions/", subscriptiontemplate.SubscriptionTemplatesActionsHandler(manager, cfg))

	mux.HandleFunc("/api/sub", subscription.SubscriptionPublicHandler(manager, cfg))
	mux.HandleFunc("/api/sub/", subscription.SubscriptionPublicHandler(manager, cfg))

	mux.HandleFunc("/api/subscriptions", subscription.SubscriptionsHandler(manager, cfg))
	mux.HandleFunc("/api/subscriptions/", subscription.SubscriptionsByPathHandler(manager, cfg))

	mux.HandleFunc("/api/subscription-page-configs", subscriptionpageconfigs.SubscriptionPageConfigsHandler(manager, cfg))
	mux.HandleFunc("/api/subscription-page-configs/", subscriptionpageconfigs.SubscriptionPageConfigByUUIDHandler(manager, cfg))
	mux.HandleFunc("/api/subscription-page-configs/actions/", subscriptionpageconfigs.SubscriptionPageConfigsActionsHandler(manager, cfg))

	mux.HandleFunc("/api/health", health.HealthHandler())
	mux.HandleFunc("/api/v1/health", health.HealthHandler())

	mux.Handle("/api/", http.NotFoundHandler())
}
