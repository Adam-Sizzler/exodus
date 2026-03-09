package httpapi

import (
	"net/http"

	"v2ray-stat/backend/config"
	dbmanager "v2ray-stat/backend/db/manager"
	"v2ray-stat/backend/httpapi/auth"
	"v2ray-stat/backend/httpapi/bandwidthstats"
	"v2ray-stat/backend/httpapi/configprofiles"
	"v2ray-stat/backend/httpapi/externalsquads"
	"v2ray-stat/backend/httpapi/health"
	"v2ray-stat/backend/httpapi/hosts"
	"v2ray-stat/backend/httpapi/hwiduserdevices"
	"v2ray-stat/backend/httpapi/infrabilling"
	"v2ray-stat/backend/httpapi/keygen"
	"v2ray-stat/backend/httpapi/middleware"
	"v2ray-stat/backend/httpapi/nodes"
	"v2ray-stat/backend/httpapi/panelsettings"
	"v2ray-stat/backend/httpapi/passkeys"
	"v2ray-stat/backend/httpapi/squads"
	"v2ray-stat/backend/httpapi/subscription"
	subscriptionpageconfigs "v2ray-stat/backend/httpapi/subscription-page-configs"
	subscriptionsettings "v2ray-stat/backend/httpapi/subscription-settings"
	subscriptiontemplate "v2ray-stat/backend/httpapi/subscription-template"
	subscriptionrequesthistory "v2ray-stat/backend/httpapi/subscriptionrequesthistory"
	"v2ray-stat/backend/httpapi/system"
	"v2ray-stat/backend/httpapi/users"
)

func NewAPIHandler(manager *dbmanager.DatabaseManager, cfg *config.BackendConfig) http.Handler {
	mux := http.NewServeMux()
	RegisterRoutes(mux, manager, cfg)
	handler := auth.WithPanelAuth(manager, cfg, mux)
	handler = middleware.WithRequestLogging(cfg, "api", handler)
	return middleware.WithCORS(cfg, handler)
}

func RegisterRoutes(mux *http.ServeMux, manager *dbmanager.DatabaseManager, cfg *config.BackendConfig) {
	mux.HandleFunc("/api/v1/auth/bootstrap", auth.AuthBootstrapHandler(manager, cfg))
	mux.HandleFunc("/api/v1/auth/setup", auth.AuthSetupHandler(manager, cfg))
	mux.HandleFunc("/api/v1/auth/login", auth.AuthLoginHandler(manager, cfg))
	mux.HandleFunc("/api/v1/auth/logout", auth.AuthLogoutHandler(manager, cfg))
	mux.HandleFunc("/api/v1/auth/me", auth.AuthMeHandler(manager, cfg))
	mux.HandleFunc("/api/auth/bootstrap", auth.AuthBootstrapHandler(manager, cfg))
	mux.HandleFunc("/api/auth/setup", auth.AuthSetupHandler(manager, cfg))
	mux.HandleFunc("/api/auth/status", auth.AuthStatusHandler(manager, cfg))
	mux.HandleFunc("/api/auth/register", auth.AuthRegisterCompatHandler(manager, cfg))
	mux.HandleFunc("/api/auth/login", auth.AuthLoginCompatHandler(manager, cfg))
	mux.HandleFunc("/api/auth/logout", auth.AuthLogoutHandler(manager, cfg))
	mux.HandleFunc("/api/auth/me", auth.AuthMeHandler(manager, cfg))

	mux.HandleFunc("/api/v1/settings", panelsettings.PanelSettingsHandler(manager, cfg))
	mux.HandleFunc("/api/remnawave-settings", panelsettings.RemnawaveSettingsHandler(manager, cfg))
	mux.HandleFunc("/api/remnawave-settings/", panelsettings.RemnawaveSettingsHandler(manager, cfg))
	mux.HandleFunc("/api/v1/api-tokens", panelsettings.PanelAPITokensHandler(manager, cfg))
	mux.HandleFunc("/api/v1/api-tokens/", panelsettings.PanelAPITokenByUUIDHandler(manager, cfg))
	mux.HandleFunc("/api/tokens", panelsettings.PanelAPITokensHandler(manager, cfg))
	mux.HandleFunc("/api/tokens/", panelsettings.PanelAPITokenByUUIDHandler(manager, cfg))

	mux.HandleFunc("/api/nodes", nodes.NodesHandler(manager, cfg))
	mux.HandleFunc("/api/nodes/", nodes.NodeByUUIDHandler(manager, cfg))
	mux.HandleFunc("/api/nodes/actions/", nodes.NodesActionsHandler(manager, cfg))
	mux.HandleFunc("/api/nodes/bulk-actions", nodes.NodesBulkActionsHandler(manager, cfg))
	mux.HandleFunc("/api/nodes/bulk-actions/", nodes.NodesBulkActionsHandler(manager, cfg))
	mux.HandleFunc("/api/nodes/tags", nodes.NodesTagsHandler(manager, cfg))
	mux.HandleFunc("/api/v1/nodes", nodes.NodesHandler(manager, cfg))
	mux.HandleFunc("/api/v1/nodes/", nodes.NodeByUUIDHandler(manager, cfg))
	mux.HandleFunc("/api/v1/nodes-with-config", squads.NodesWithConfigHandler(manager, cfg))

	mux.HandleFunc("/api/hosts", hosts.HostsHandler(manager, cfg))
	mux.HandleFunc("/api/hosts/", hosts.HostByUUIDHandler(manager, cfg))
	mux.HandleFunc("/api/hosts/actions/", hosts.HostsActionsHandler(manager, cfg))
	mux.HandleFunc("/api/hosts/bulk/", hosts.HostsBulkHandler(manager, cfg))
	mux.HandleFunc("/api/hosts/tags", hosts.HostsTagsHandler(manager, cfg))

	mux.HandleFunc("/api/users", users.UsersHandler(manager, cfg))
	mux.HandleFunc("/api/users/", users.UserByUUIDHandler(manager, cfg))
	mux.HandleFunc("/api/users/bulk/", users.UsersBulkHandler(manager, cfg))
	mux.HandleFunc("/api/users/tags", users.UsersTagsHandler(manager, cfg))
	mux.HandleFunc("/api/v1/users-list", users.UsersHandler(manager, cfg))
	mux.HandleFunc("/api/v1/users-list/", users.UserByUUIDHandler(manager, cfg))

	mux.HandleFunc("/api/keygen", keygen.KeygenHandler(manager, cfg))
	mux.HandleFunc("/api/keygen/", keygen.KeygenHandler(manager, cfg))
	mux.HandleFunc("/api/passkeys", passkeys.PasskeysHandler(manager, cfg))
	mux.HandleFunc("/api/passkeys/", passkeys.PasskeysHandler(manager, cfg))

	mux.HandleFunc("/api/bandwidth-stats/nodes", bandwidthstats.NodesHandler(manager, cfg))
	mux.HandleFunc("/api/bandwidth-stats/nodes/", bandwidthstats.NodesHandler(manager, cfg))
	mux.HandleFunc("/api/bandwidth-stats/users", bandwidthstats.UsersHandler(manager, cfg))
	mux.HandleFunc("/api/bandwidth-stats/users/", bandwidthstats.UsersHandler(manager, cfg))

	mux.HandleFunc("/api/config-profiles", configprofiles.ConfigProfilesHandler(manager, cfg))
	mux.HandleFunc("/api/config-profiles/", configprofiles.ConfigProfileByUUIDHandler(manager, cfg))
	mux.HandleFunc("/api/config-profiles/actions/", configprofiles.ConfigProfilesActionsHandler(manager, cfg))
	mux.HandleFunc("/api/config-profiles/inbounds", configprofiles.ConfigProfilesInboundsHandler(manager, cfg))
	mux.HandleFunc("/api/config-profiles/snippets", configprofiles.ConfigProfileSnippetsHandler(manager, cfg))
	mux.HandleFunc("/api/snippets", configprofiles.ConfigProfileSnippetsHandler(manager, cfg))
	mux.HandleFunc("/api/snippets/", configprofiles.ConfigProfileSnippetsHandler(manager, cfg))
	mux.HandleFunc("/api/v1/config-profiles", configprofiles.ConfigProfilesHandler(manager, cfg))
	mux.HandleFunc("/api/v1/config-profiles/", configprofiles.ConfigProfileByUUIDHandler(manager, cfg))
	mux.HandleFunc("/api/v1/config-profiles/reorder", configprofiles.ConfigProfilesActionsHandler(manager, cfg))
	mux.HandleFunc("/api/v1/snippets", configprofiles.ConfigProfileSnippetsHandler(manager, cfg))
	mux.HandleFunc("/api/v1/config-profiles-with-inbounds", squads.ConfigProfilesWithInboundsHandler(manager, cfg))

	mux.HandleFunc("/api/v1/internal-squads", squads.InternalSquadsHandler(manager, cfg))
	mux.HandleFunc("/api/v1/internal-squads/", squads.InternalSquadByUUIDHandler(manager, cfg))
	mux.HandleFunc("/api/v1/internal-squads/actions/reorder", squads.InternalSquadsReorderHandler(manager, cfg))
	mux.HandleFunc("/api/v1/internal-squads/reorder", squads.InternalSquadsReorderHandler(manager, cfg))
	mux.HandleFunc("/api/v1/squads-summary", squads.AllSquadsSummaryHandler(manager, cfg))
	mux.HandleFunc("/api/v1/squad-inbounds", squads.SquadInboundsHandler(manager, cfg))
	mux.HandleFunc("/api/v1/squad-members", squads.SquadMembersHandler(manager, cfg))
	mux.HandleFunc("/api/v1/squad-details/", squads.SquadDetailsHandler(manager, cfg))
	mux.HandleFunc("/api/internal-squads", squads.InternalSquadsHandler(manager, cfg))
	mux.HandleFunc("/api/internal-squads/", squads.InternalSquadByUUIDHandler(manager, cfg))
	mux.HandleFunc("/api/internal-squads/actions/reorder", squads.InternalSquadsReorderHandler(manager, cfg))
	mux.HandleFunc("/api/internal-squads/reorder", squads.InternalSquadsReorderHandler(manager, cfg))
	mux.HandleFunc("/api/squads-summary", squads.AllSquadsSummaryHandler(manager, cfg))
	mux.HandleFunc("/api/squad-inbounds", squads.SquadInboundsHandler(manager, cfg))
	mux.HandleFunc("/api/squad-members", squads.SquadMembersHandler(manager, cfg))
	mux.HandleFunc("/api/squad-details/", squads.SquadDetailsHandler(manager, cfg))

	mux.HandleFunc("/api/v1/inbound-assignments", squads.InboundAssignmentsHandler(manager, cfg))
	mux.HandleFunc("/api/v1/inbounds-with-profiles", squads.InboundsWithProfilesHandler(manager, cfg))

	mux.HandleFunc("/api/v1/external-squads", externalsquads.ExternalSquadsHandler(manager, cfg))
	mux.HandleFunc("/api/v1/external-squads/", externalsquads.ExternalSquadByUUIDHandler(manager, cfg))
	mux.HandleFunc("/api/v1/external-squads/actions/reorder", externalsquads.ExternalSquadsReorderHandler(manager, cfg))
	mux.HandleFunc("/api/v1/external-squads/reorder", externalsquads.ExternalSquadsReorderHandler(manager, cfg))
	mux.HandleFunc("/api/external-squads", externalsquads.ExternalSquadsHandler(manager, cfg))
	mux.HandleFunc("/api/external-squads/", externalsquads.ExternalSquadByUUIDHandler(manager, cfg))
	mux.HandleFunc("/api/external-squads/actions/reorder", externalsquads.ExternalSquadsReorderHandler(manager, cfg))
	mux.HandleFunc("/api/external-squads/reorder", externalsquads.ExternalSquadsReorderHandler(manager, cfg))

	mux.HandleFunc("/api/v1/hwid-user-devices", hwiduserdevices.HWIDUserDevicesHandler(manager, cfg))
	mux.HandleFunc("/api/v1/hwid-user-devices/", hwiduserdevices.HWIDDeviceByUUIDHandler(manager, cfg))
	mux.HandleFunc("/api/v1/hwid-user-devices/user/", hwiduserdevices.HWIDUserDevicesByUserUUIDHandler(manager, cfg))
	mux.HandleFunc("/api/v1/hwid-user-devices/check", hwiduserdevices.HWIDCheckHandler(manager, cfg))
	mux.HandleFunc("/api/hwid/devices", hwiduserdevices.HWIDCompatDevicesHandler(manager, cfg))
	mux.HandleFunc("/api/hwid/devices/", hwiduserdevices.HWIDCompatDevicesHandler(manager, cfg))
	mux.HandleFunc("/api/hwid/devices/stats", hwiduserdevices.HWIDCompatStatsHandler(manager, cfg))
	mux.HandleFunc("/api/hwid/devices/top-users", hwiduserdevices.HWIDCompatTopUsersHandler(manager, cfg))

	mux.HandleFunc("/api/subscription-settings", subscriptionsettings.SubscriptionSettingsHandler(manager, cfg))
	mux.HandleFunc("/api/subscription-settings/", subscriptionsettings.SubscriptionSettingsHandler(manager, cfg))
	// mux.HandleFunc("/api/v1/subscription-settings/", subscriptionsettings.SubscriptionSettingsByUUIDHandler(manager, cfg))

	mux.HandleFunc("/api/infra-billing/providers", infrabilling.ProvidersHandler(manager, cfg))
	mux.HandleFunc("/api/infra-billing/providers/", infrabilling.ProvidersHandler(manager, cfg))
	mux.HandleFunc("/api/infra-billing/nodes", infrabilling.BillingNodesHandler(manager, cfg))
	mux.HandleFunc("/api/infra-billing/nodes/", infrabilling.BillingNodesHandler(manager, cfg))
	mux.HandleFunc("/api/infra-billing/history", infrabilling.BillingHistoryHandler(manager, cfg))
	mux.HandleFunc("/api/infra-billing/history/", infrabilling.BillingHistoryHandler(manager, cfg))

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
	mux.HandleFunc("/api/subscription-request-history", subscriptionrequesthistory.SubscriptionRequestHistoryHandler(manager, cfg))
	mux.HandleFunc("/api/subscription-request-history/", subscriptionrequesthistory.SubscriptionRequestHistoryHandler(manager, cfg))
	mux.HandleFunc("/api/subscription-request-history/stats", subscriptionrequesthistory.SubscriptionRequestHistoryStatsHandler(manager, cfg))

	mux.HandleFunc("/api/system/metadata", system.MetadataHandler(cfg))
	mux.HandleFunc("/api/system/stats", system.StatsHandler(manager, cfg))
	mux.HandleFunc("/api/system/stats/bandwidth", system.BandwidthStatsHandler(manager, cfg))
	mux.HandleFunc("/api/system/stats/nodes", system.NodesStatsHandler(manager, cfg))
	mux.HandleFunc("/api/system/nodes/metrics", system.NodesMetricsHandler(cfg))
	mux.HandleFunc("/api/system/health", system.HealthHandler(cfg))

	mux.HandleFunc("/api/v1/system/metadata", system.MetadataHandler(cfg))
	mux.HandleFunc("/api/v1/system/stats", system.StatsHandler(manager, cfg))
	mux.HandleFunc("/api/v1/system/stats/bandwidth", system.BandwidthStatsHandler(manager, cfg))
	mux.HandleFunc("/api/v1/system/stats/nodes", system.NodesStatsHandler(manager, cfg))
	mux.HandleFunc("/api/v1/system/nodes/metrics", system.NodesMetricsHandler(cfg))
	mux.HandleFunc("/api/v1/system/health", system.HealthHandler(cfg))

	mux.HandleFunc("/api/health", health.HealthHandler())
	mux.HandleFunc("/api/v1/health", health.HealthHandler())

	mux.Handle("/api/", http.NotFoundHandler())
}
