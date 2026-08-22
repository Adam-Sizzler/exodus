export declare const ROOT: "/api";
export declare const METRICS_ROOT: "/metrics";
export declare const HEALTH_ROOT: "/health";
export declare const BACKEND_TOOLS_ROOT: "/backend-tools";
export declare const BULLBOARD_ROOT: "/backend-tools/queues";
export declare const SWAGGER_ROOT: "/backend-tools/swagger";
export declare const SCALAR_ROOT: "/backend-tools/scalar";
export declare const REST_API: {
    readonly AUTH: {
        readonly LOGIN: "/api/auth/login";
        readonly REGISTER: "/api/auth/register";
        readonly GET_STATUS: "/api/auth/status";
        readonly OAUTH2: {
            readonly TELEGRAM_CALLBACK: "/api/auth/oauth2/tg/callback";
            readonly AUTHORIZE: "/api/auth/oauth2/authorize";
            readonly CALLBACK: "/api/auth/oauth2/callback";
        };
        readonly PASSKEY: {
            readonly GET_AUTHENTICATION_OPTIONS: "/api/auth/passkey/authentication/options";
            readonly VERIFY_AUTHENTICATION: "/api/auth/passkey/authentication/verify";
        };
    };
    readonly PASSKEYS: {
        readonly GET_ALL_PASSKEYS: "/api/passkeys/";
        readonly DELETE_PASSKEY: "/api/passkeys/";
        readonly UPDATE_PASSKEY: "/api/passkeys/";
        readonly GET_REGISTRATION_OPTIONS: "/api/passkeys/registration/options";
        readonly VERIFY_REGISTRATION: "/api/passkeys/registration/verify";
    };
    readonly API_TOKENS: {
        readonly CREATE: "/api/tokens/";
        readonly DELETE: (uuid: string) => string;
        readonly GET: "/api/tokens/";
        readonly GET_SCOPES: "/api/tokens/scopes";
        readonly OTT: "/api/tokens/ott";
    };
    readonly KEYGEN: {
        readonly GET: "/api/keygen/";
    };
    readonly NODES: {
        readonly CREATE: "/api/nodes/";
        readonly GET: "/api/nodes/";
        readonly GET_BY_UUID: (uuid: string) => string;
        readonly UPDATE: "/api/nodes/";
        readonly DELETE: (uuid: string) => string;
        readonly TAGS: {
            readonly GET: "/api/nodes/tags";
        };
        readonly ACTIONS: {
            readonly ENABLE: (uuid: string) => string;
            readonly DISABLE: (uuid: string) => string;
            readonly RESTART: (uuid: string) => string;
            readonly RESTART_ALL: "/api/nodes/actions/restart-all";
            readonly RESET_TRAFFIC: (uuid: string) => string;
            readonly REORDER: "/api/nodes/actions/reorder";
        };
        readonly BULK_ACTIONS: {
            readonly PROFILE_MODIFICATION: "/api/nodes/bulk-actions/profile-modification";
            readonly ACTIONS: "/api/nodes/bulk-actions";
            readonly UPDATE: "/api/nodes/bulk-actions/update";
        };
    };
    readonly USERS: {
        readonly CREATE: "/api/users/";
        readonly UPDATE: "/api/users/";
        readonly GET: "/api/users/";
        readonly STREAM: "/api/users/stream";
        readonly DELETE: (userId: string) => string;
        readonly GET_BY_ID: (userId: string) => string;
        readonly ACCESSIBLE_NODES: (userId: string) => string;
        readonly SUBSCRIPTION_REQUEST_HISTORY: (userId: string) => string;
        readonly ACTIONS: {
            readonly DISABLE: (userId: string) => string;
            readonly ENABLE: (userId: string) => string;
            readonly RESET_TRAFFIC: (userId: string) => string;
            readonly REVOKE_SUBSCRIPTION: (userId: string) => string;
            readonly EXTEND_EXPIRATION_DATE: (userId: string) => string;
        };
        readonly GET_BY: {
            readonly SHORT_UUID: (shortUuid: string) => string;
            readonly USERNAME: (username: string) => string;
        };
        readonly RESOLVE: "/api/users/resolve";
        readonly BULK: {
            readonly DELETE_BY_STATUS: "/api/users/bulk/delete-by-status";
            readonly UPDATE: "/api/users/bulk/update";
            readonly RESET_TRAFFIC: "/api/users/bulk/reset-traffic";
            readonly REVOKE_SUBSCRIPTION: "/api/users/bulk/revoke-subscription";
            readonly DELETE: "/api/users/bulk/delete";
            readonly UPDATE_SQUADS: "/api/users/bulk/update-squads";
            readonly EXTEND_EXPIRATION_DATE: "/api/users/bulk/extend-expiration-date";
            readonly ALL: {
                readonly UPDATE: "/api/users/bulk/all/update";
                readonly RESET_TRAFFIC: "/api/users/bulk/all/reset-traffic";
                readonly EXTEND_EXPIRATION_DATE: "/api/users/bulk/all/extend-expiration-date";
            };
        };
        readonly TAGS: {
            readonly GET: "/api/users/tags";
        };
    };
    readonly SUBSCRIPTION: {
        readonly GET: (shortUuid: string) => string;
        readonly GET_INFO: (shortUuid: string) => string;
    };
    readonly HOSTS: {
        readonly CREATE: "/api/hosts/";
        readonly UPDATE: "/api/hosts/";
        readonly GET: "/api/hosts/";
        readonly GET_BY_UUID: (uuid: string) => string;
        readonly DELETE: (uuid: string) => string;
        readonly ACTIONS: {
            readonly REORDER: "/api/hosts/actions/reorder";
        };
        readonly BULK: {
            readonly ENABLE_HOSTS: "/api/hosts/bulk/enable";
            readonly DISABLE_HOSTS: "/api/hosts/bulk/disable";
            readonly DELETE_HOSTS: "/api/hosts/bulk/delete";
            readonly UPDATE: "/api/hosts/bulk/update";
        };
        readonly TAGS: {
            readonly GET: "/api/hosts/tags";
        };
    };
    readonly SYSTEM: {
        readonly HEALTH: "/api/system/health";
        readonly METADATA: "/api/system/metadata";
        readonly CONFIGURATION: "/api/system/configuration";
        readonly STATS: {
            readonly SYSTEM_STATS: "/api/system/stats";
            readonly BANDWIDTH_STATS: "/api/system/stats/bandwidth";
            readonly NODES_STATS: "/api/system/stats/nodes";
            readonly NODES_METRICS: "/api/system/nodes/metrics";
            readonly RECAP: "/api/system/stats/recap";
            readonly DIGEST: "/api/system/stats/digest";
            readonly HTTP: "/api/system/stats/http";
        };
        readonly TOOLS: {
            readonly GENERATE_X25519: "/api/system/tools/x25519/generate";
        };
        readonly TESTERS: {
            readonly SRR_MATCHER: "/api/system/testers/srr-matcher";
        };
    };
    readonly SUBSCRIPTION_TEMPLATE: {
        readonly GET: (uuid: string) => string;
        readonly UPDATE: "/api/subscription-templates/";
        readonly DELETE: (uuid: string) => string;
        readonly GET_ALL: "/api/subscription-templates/";
        readonly CREATE: "/api/subscription-templates/";
        readonly ACTIONS: {
            readonly REORDER: "/api/subscription-templates/actions/reorder";
        };
    };
    readonly SUBSCRIPTION_SETTINGS: {
        readonly GET: "/api/subscription-settings/";
        readonly UPDATE: "/api/subscription-settings/";
    };
    readonly HWID: {
        readonly GET_ALL_HWID_DEVICES: "/api/hwid/devices";
        readonly CREATE_USER_HWID_DEVICE: "/api/hwid/devices";
        readonly GET_USER_HWID_DEVICES: (userId: string) => string;
        readonly DELETE_USER_HWID_DEVICE: "/api/hwid/devices/delete";
        readonly DELETE_ALL_USER_HWID_DEVICES: "/api/hwid/devices/delete-all";
        readonly STATS: "/api/hwid/devices/stats";
        readonly TOP_USERS_BY_DEVICES: "/api/hwid/devices/top-users";
    };
    readonly SUBSCRIPTIONS: {
        readonly GET: "/api/subscriptions/";
        readonly GET_BY: {
            readonly USERNAME: (username: string) => string;
            readonly SHORT_UUID: (shortUuid: string) => string;
            readonly SHORT_UUID_RAW: (shortUuid: string) => string;
            readonly ID: (userId: string) => string;
        };
        readonly SUBPAGE: {
            readonly GET_CONFIG: (shortUuid: string) => string;
        };
        readonly GET_CONNECTION_KEYS_BY_USER_ID: (userId: string) => string;
    };
    readonly CONFIG_PROFILES: {
        readonly GET: "/api/config-profiles/";
        readonly CREATE: "/api/config-profiles/";
        readonly UPDATE: "/api/config-profiles/";
        readonly GET_BY_UUID: (uuid: string) => string;
        readonly DELETE: (uuid: string) => string;
        readonly GET_INBOUNDS_BY_PROFILE_UUID: (uuid: string) => string;
        readonly GET_COMPUTED_CONFIG_BY_PROFILE_UUID: (uuid: string) => string;
        readonly GET_ALL_INBOUNDS: "/api/config-profiles/inbounds";
        readonly ACTIONS: {
            readonly REORDER: "/api/config-profiles/actions/reorder";
        };
    };
    readonly INTERNAL_SQUADS: {
        readonly GET: "/api/internal-squads/";
        readonly CREATE: "/api/internal-squads/";
        readonly UPDATE: "/api/internal-squads/";
        readonly GET_BY_UUID: (uuid: string) => string;
        readonly DELETE: (uuid: string) => string;
        readonly ACCESSIBLE_NODES: (uuid: string) => string;
        readonly BULK_ACTIONS: {
            readonly ADD_USERS: (uuid: string) => string;
            readonly REMOVE_USERS: (uuid: string) => string;
            readonly ADD_MANY_USERS: (uuid: string) => string;
            readonly REMOVE_MANY_USERS: (uuid: string) => string;
        };
        readonly ACTIONS: {
            readonly REORDER: "/api/internal-squads/actions/reorder";
        };
    };
    readonly INFRA_BILLING: {
        readonly GET_PROVIDERS: "/api/infra-billing/providers";
        readonly CREATE_PROVIDER: "/api/infra-billing/providers";
        readonly UPDATE_PROVIDER: "/api/infra-billing/providers";
        readonly DELETE_PROVIDER: (uuid: string) => string;
        readonly GET_PROVIDER_BY_UUID: (uuid: string) => string;
        readonly GET_BILLING_NODES: "/api/infra-billing/nodes";
        readonly CREATE_BILLING_NODE: "/api/infra-billing/nodes";
        readonly UPDATE_BILLING_NODE: "/api/infra-billing/nodes";
        readonly DELETE_BILLING_NODE: (uuid: string) => string;
        readonly GET_BILLING_HISTORY: "/api/infra-billing/history";
        readonly CREATE_BILLING_HISTORY: "/api/infra-billing/history";
        readonly DELETE_BILLING_HISTORY: (uuid: string) => string;
    };
    readonly SUBSCRIPTION_REQUEST_HISTORY: {
        readonly GET: "/api/subscription-request-history/";
        readonly STATS: "/api/subscription-request-history/stats";
    };
    readonly SNIPPETS: {
        readonly GET: "/api/snippets/";
        readonly CREATE: "/api/snippets/";
        readonly UPDATE: "/api/snippets/";
        readonly DELETE: "/api/snippets/";
        readonly ACTIONS: {
            readonly SYNC: "/api/snippets/actions/sync";
        };
    };
    readonly EXTERNAL_SQUADS: {
        readonly GET: "/api/external-squads/";
        readonly CREATE: "/api/external-squads/";
        readonly UPDATE: "/api/external-squads/";
        readonly GET_BY_UUID: (uuid: string) => string;
        readonly DELETE: (uuid: string) => string;
        readonly BULK_ACTIONS: {
            readonly ADD_USERS: (uuid: string) => string;
            readonly REMOVE_USERS: (uuid: string) => string;
        };
        readonly ACTIONS: {
            readonly REORDER: "/api/external-squads/actions/reorder";
        };
    };
    readonly EXODUS_SETTINGS: {
        readonly GET: "/api/exodus-settings/";
        readonly UPDATE: "/api/exodus-settings/";
    };
    readonly SUBSCRIPTION_PAGE_CONFIGS: {
        readonly GET: (uuid: string) => string;
        readonly GET_ALL: "/api/subscription-page-configs/";
        readonly UPDATE: "/api/subscription-page-configs/";
        readonly DELETE: (uuid: string) => string;
        readonly CREATE: "/api/subscription-page-configs/";
        readonly ACTIONS: {
            readonly REORDER: "/api/subscription-page-configs/actions/reorder";
            readonly CLONE: "/api/subscription-page-configs/actions/clone";
        };
    };
    readonly NODE_INTEGRATIONS: {
        readonly GET: (uuid: string) => string;
        readonly GET_ALL: "/api/node-integrations/";
        readonly UPDATE: "/api/node-integrations/";
        readonly DELETE: (uuid: string) => string;
        readonly CREATE: "/api/node-integrations/";
    };
    readonly NODE_PLUGINS: {
        readonly GET: (uuid: string) => string;
        readonly GET_ALL: "/api/node-plugins/";
        readonly UPDATE: "/api/node-plugins/";
        readonly DELETE: (uuid: string) => string;
        readonly CREATE: "/api/node-plugins/";
        readonly ACTIONS: {
            readonly REORDER: "/api/node-plugins/actions/reorder";
            readonly CLONE: "/api/node-plugins/actions/clone";
            readonly SYNC: "/api/node-plugins/actions/sync";
        };
        readonly EXECUTOR: "/api/node-plugins/executor";
        readonly SHARED_LISTS: {
            readonly GET_ALL: "/api/node-plugins/shared-lists";
            readonly GET: (name: string) => string;
            readonly CREATE: "/api/node-plugins/shared-lists";
            readonly UPDATE: "/api/node-plugins/shared-lists";
            readonly DELETE: (name: string) => string;
            readonly ACTIONS: {
                readonly SYNC: "/api/node-plugins/shared-lists/actions/sync";
            };
        };
        readonly TORRENT_BLOCKER: {
            readonly GET_REPORTS: "/api/node-plugins/torrent-blocker";
            readonly GET_REPORTS_STATS: "/api/node-plugins/torrent-blocker/stats";
            readonly TRUNCATE_REPORTS: "/api/node-plugins/torrent-blocker/truncate";
        };
    };
    readonly BANDWIDTH_STATS: {
        readonly NODES: {
            readonly GET: "/api/bandwidth-stats/nodes/";
            readonly GET_REALTIME: "/api/bandwidth-stats/nodes/realtime";
            readonly GET_USERS: (uuid: string) => string;
            readonly GET_USERS_BY_NODES: "/api/bandwidth-stats/nodes/users";
            readonly GET_USAGE: "/api/bandwidth-stats/nodes/usage";
        };
        readonly USERS: {
            readonly GET_BY_ID: (userId: string) => string;
        };
        readonly INTERNAL_SQUADS: {
            readonly GET_USAGE: (uuid: string) => string;
            readonly USER_USAGE: (squadUuid: string, userId: string) => string;
        };
    };
    readonly CONNECTIONS: {
        readonly CONNECTIONS_BY_USER: (userId: string) => string;
        readonly CONNECTIONS_BY_USER_RESULT: (jobId: string) => string;
        readonly CONNECTIONS_BY_NODE: (uuid: string) => string;
        readonly CONNECTIONS_BY_NODE_RESULT: (jobId: string) => string;
        readonly GEOCHECK_BY_NODE: (uuid: string) => string;
        readonly GEOCHECK_BY_NODE_RESULT: (jobId: string) => string;
        readonly DROP_CONNECTIONS: "/api/connections/drop";
    };
    readonly METADATA: {
        readonly NODE: {
            readonly GET: (uuid: string) => string;
            readonly UPSERT: (uuid: string) => string;
        };
        readonly USER: {
            readonly GET: (userId: string) => string;
            readonly UPSERT: (userId: string) => string;
        };
    };
};
//# sourceMappingURL=routes.d.ts.map