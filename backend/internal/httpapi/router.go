package httpapi

import (
	"net/http"

	"exodus/internal/config"
	dbmanager "exodus/internal/db/manager"
	"exodus/internal/httpapi/auth"
	"exodus/internal/httpapi/bandwidthstats"
	"exodus/internal/httpapi/configprofiles"
	"exodus/internal/httpapi/externalsquads"
	"exodus/internal/httpapi/health"
	"exodus/internal/httpapi/hosts"
	"exodus/internal/httpapi/hwiduserdevices"
	"exodus/internal/httpapi/infrabilling"
	"exodus/internal/httpapi/keygen"
	"exodus/internal/httpapi/middleware"
	"exodus/internal/httpapi/modulessettings"
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

func NewAPIHandler(manager *dbmanager.DatabaseManager, cfg *config.BackendConfig) http.Handler {
	mux := http.NewServeMux()
	RegisterRoutes(mux, manager, cfg)
	handler := auth.WithPanelAuth(manager, cfg, mux)
	handler = middleware.WithRequestLogging(cfg, "api", handler)
	return middleware.WithCORS(cfg, handler)
}

func RegisterRoutes(mux *http.ServeMux, manager *dbmanager.DatabaseManager, cfg *config.BackendConfig) {
	mux.HandleFunc("/api/auth/bootstrap", auth.AuthBootstrapHandler(manager, cfg))
	mux.HandleFunc("/api/auth/setup", auth.AuthSetupHandler(manager, cfg))
	mux.HandleFunc("/api/auth/status", auth.AuthStatusHandler(manager, cfg))
	mux.HandleFunc("/api/auth/register", auth.AuthRegisterCompatHandler(manager, cfg))
	mux.HandleFunc("/api/auth/login", auth.AuthLoginCompatHandler(manager, cfg))
	mux.HandleFunc("/api/auth/logout", auth.AuthLogoutHandler(manager, cfg))
	mux.HandleFunc("/api/auth/me", auth.AuthMeHandler(manager, cfg))
	mux.HandleFunc("/api/auth/oauth2/authorize", auth.OAuth2AuthorizeHandler(manager, cfg))
	mux.HandleFunc("/api/auth/oauth2/callback", auth.OAuth2CallbackHandler(manager, cfg))
	mux.HandleFunc("/api/auth/oauth2/tg/callback", auth.TelegramCallbackHandler(manager, cfg))
	mux.HandleFunc("/api/auth/passkey/authentication/options", passkeys.AuthenticationOptionsHandler(manager, cfg))
	mux.HandleFunc("/api/auth/passkey/authentication/verify", passkeys.VerifyAuthenticationHandler(manager, cfg))

	mux.HandleFunc("/api/exodus-settings", panelsettings.ExodusSettingsHandler(manager, cfg))
	mux.HandleFunc("/api/exodus-settings/", panelsettings.ExodusSettingsHandler(manager, cfg))
	mux.HandleFunc("/api/tokens", panelsettings.PanelAPITokensHandler(manager, cfg))
	mux.HandleFunc("/api/tokens/", panelsettings.PanelAPITokenByUUIDHandler(manager, cfg))
	mux.HandleFunc("/api/modules-settings", modulessettings.ModulesSettingsHandler(manager, cfg))
	mux.HandleFunc("/api/modules-settings/", modulessettings.ModulesSettingsHandler(manager, cfg))

	mux.HandleFunc("/api/nodes", nodes.NodesHandler(manager, cfg))
	mux.HandleFunc("/api/nodes/", nodes.NodeByUUIDHandler(manager, cfg))
	mux.HandleFunc("/api/nodes/actions/", nodes.NodesActionsHandler(manager, cfg))
	mux.HandleFunc("/api/nodes/bulk-actions", nodes.NodesBulkActionsHandler(manager, cfg))
	mux.HandleFunc("/api/nodes/bulk-actions/", nodes.NodesBulkActionsHandler(manager, cfg))
	mux.HandleFunc("/api/nodes/tags", nodes.NodesTagsHandler(manager, cfg))
	mux.HandleFunc("/api/nodes-with-config", squads.NodesWithConfigHandler(manager, cfg))
	mux.HandleFunc("/api/node-plugins", nodeplugins.Handler(manager, cfg))
	mux.HandleFunc("/api/node-plugins/", nodeplugins.Handler(manager, cfg))

	mux.HandleFunc("/api/subscription-connections", subscriptionconnections.NodesHandler(manager, cfg))
	mux.HandleFunc("/api/subscription-connections/", subscriptionconnections.NodeByUUIDHandler(manager, cfg))
	mux.HandleFunc("/api/subscription-connections/actions/", subscriptionconnections.NodesActionsHandler(manager, cfg))
	mux.HandleFunc("/api/subscription-connections/bulk-actions", subscriptionconnections.NodesBulkActionsHandler(manager, cfg))
	mux.HandleFunc("/api/subscription-connections/bulk-actions/", subscriptionconnections.NodesBulkActionsHandler(manager, cfg))
	mux.HandleFunc("/api/subscription-connections/tags", subscriptionconnections.NodesTagsHandler(manager, cfg))

	mux.HandleFunc("/api/hosts", hosts.HostsHandler(manager, cfg))
	mux.HandleFunc("/api/hosts/", hosts.HostByUUIDHandler(manager, cfg))
	mux.HandleFunc("/api/hosts/actions/", hosts.HostsActionsHandler(manager, cfg))
	mux.HandleFunc("/api/hosts/bulk/", hosts.HostsBulkHandler(manager, cfg))
	mux.HandleFunc("/api/hosts/tags", hosts.HostsTagsHandler(manager, cfg))

	mux.HandleFunc("/api/users", users.UsersHandler(manager, cfg))
	mux.HandleFunc("/api/users/", users.UserByUUIDHandler(manager, cfg))
	mux.HandleFunc("/api/users/bulk/", users.UsersBulkHandler(manager, cfg))
	mux.HandleFunc("/api/users/tags", users.UsersTagsHandler(manager, cfg))

	mux.HandleFunc("/api/keygen", keygen.KeygenHandler(manager, cfg))
	mux.HandleFunc("/api/keygen/", keygen.KeygenHandler(manager, cfg))
	mux.HandleFunc("/api/passkeys/registration/options", passkeys.RegistrationOptionsHandler(manager, cfg))
	mux.HandleFunc("/api/passkeys/registration/verify", passkeys.VerifyRegistrationHandler(manager, cfg))
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
	mux.HandleFunc("/api/config-profiles-with-inbounds", squads.ConfigProfilesWithInboundsHandler(manager, cfg))

	mux.HandleFunc("/api/internal-squads", squads.InternalSquadsHandler(manager, cfg))
	mux.HandleFunc("/api/internal-squads/", squads.InternalSquadByUUIDHandler(manager, cfg))
	mux.HandleFunc("/api/internal-squads/actions/reorder", squads.InternalSquadsReorderHandler(manager, cfg))
	mux.HandleFunc("/api/internal-squads/reorder", squads.InternalSquadsReorderHandler(manager, cfg))
	mux.HandleFunc("/api/squads-summary", squads.AllSquadsSummaryHandler(manager, cfg))
	mux.HandleFunc("/api/squad-inbounds", squads.SquadInboundsHandler(manager, cfg))
	mux.HandleFunc("/api/squad-members", squads.SquadMembersHandler(manager, cfg))
	mux.HandleFunc("/api/squad-details/", squads.SquadDetailsHandler(manager, cfg))

	mux.HandleFunc("/api/inbound-assignments", squads.InboundAssignmentsHandler(manager, cfg))
	mux.HandleFunc("/api/inbounds-with-profiles", squads.InboundsWithProfilesHandler(manager, cfg))

	mux.HandleFunc("/api/external-squads", externalsquads.ExternalSquadsHandler(manager, cfg))
	mux.HandleFunc("/api/external-squads/", externalsquads.ExternalSquadByUUIDHandler(manager, cfg))
	mux.HandleFunc("/api/external-squads/actions/reorder", externalsquads.ExternalSquadsReorderHandler(manager, cfg))
	mux.HandleFunc("/api/external-squads/reorder", externalsquads.ExternalSquadsReorderHandler(manager, cfg))

	mux.HandleFunc("/api/srs-lists", srslists.SRSListsHandler(manager, cfg))
	mux.HandleFunc("/api/srs-lists/", srslists.SRSListByUUIDHandler(manager, cfg))
	mux.HandleFunc("/api/srs-lists/actions/", srslists.SRSListsActionsHandler(manager, cfg))
	mux.HandleFunc("/api/srs-lists/bulk/", srslists.SRSListsBulkHandler(manager, cfg))

	mux.HandleFunc("/api/hwid/devices/delete-all", hwiduserdevices.HWIDCompatDeleteAllUserDevicesHandler(manager, cfg))
	mux.HandleFunc("/api/hwid/devices/delete", hwiduserdevices.HWIDCompatDevicesHandler(manager, cfg))
	mux.HandleFunc("/api/hwid/devices", hwiduserdevices.HWIDCompatDevicesHandler(manager, cfg))
	mux.HandleFunc("/api/hwid/devices/", hwiduserdevices.HWIDCompatDevicesHandler(manager, cfg))
	mux.HandleFunc("/api/hwid/devices/stats", hwiduserdevices.HWIDCompatStatsHandler(manager, cfg))
	mux.HandleFunc("/api/hwid/devices/top-users", hwiduserdevices.HWIDCompatTopUsersHandler(manager, cfg))

	mux.HandleFunc("/api/subscription-settings", subscriptionsettings.SubscriptionSettingsHandler(manager, cfg))
	mux.HandleFunc("/api/subscription-settings/", subscriptionsettings.SubscriptionSettingsHandler(manager, cfg))

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
	mux.HandleFunc("/api/system/stats/recap", system.RecapHandler(manager, cfg))
	mux.HandleFunc("/api/system/nodes/metrics", system.NodesMetricsHandler(manager, cfg))
	mux.HandleFunc("/api/system/health", system.HealthHandler(cfg))
	mux.HandleFunc("/api/system/testers/srr-matcher", system.TestSRRMatcherHandler(cfg))
	mux.HandleFunc("/api/system/tools/happ/encrypt", system.EncryptHappCryptoLinkHandler(cfg))
	mux.HandleFunc("/api/system/tools/x25519/generate", system.GenerateX25519Handler(cfg))

	mux.HandleFunc("/api/health", health.HealthHandler())

	mux.Handle("/api/", http.NotFoundHandler())
}
