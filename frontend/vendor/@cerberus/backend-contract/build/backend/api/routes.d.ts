export declare const ROOT: "/api";
export declare const METRICS_ROOT: "/metrics";
export declare const HEALTH_ROOT: "/health";
export declare const BULLBOARD_ROOT: "/queues";
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
        };
    };
    readonly USERS: {
        readonly CREATE: "/api/users/";
        readonly UPDATE: "/api/users/";
        readonly GET: "/api/users/";
        readonly DELETE: (uuid: string) => string;
        readonly GET_BY_UUID: (uuid: string) => string;
        readonly ACCESSIBLE_NODES: (uuid: string) => string;
        readonly SUBSCRIPTION_REQUEST_HISTORY: (uuid: string) => string;
        readonly ACTIONS: {
            readonly DISABLE: (uuid: string) => string;
            readonly ENABLE: (uuid: string) => string;
            readonly RESET_TRAFFIC: (uuid: string) => string;
            readonly REVOKE_SUBSCRIPTION: (uuid: string) => string;
        };
        readonly GET_BY: {
            readonly ID: (id: string) => string;
            readonly SHORT_UUID: (shortUuid: string) => string;
            readonly USERNAME: (username: string) => string;
            readonly SUBSCRIPTION_UUID: (subscriptionUuid: string) => string;
            readonly TELEGRAM_ID: (telegramId: string) => string;
            readonly EMAIL: (email: string) => string;
            readonly TAG: (tag: string) => string;
        };
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
        readonly GET_OUTLINE: (shortUuid: string) => string;
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
            readonly SET_INBOUND: "/api/hosts/bulk/set-inbound";
            readonly SET_PORT: "/api/hosts/bulk/set-port";
        };
        readonly TAGS: {
            readonly GET: "/api/hosts/tags";
        };
    };
    readonly SYSTEM: {
        readonly HEALTH: "/api/system/health";
        readonly METADATA: "/api/system/metadata";
        readonly STATS: {
            readonly SYSTEM_STATS: "/api/system/stats";
            readonly BANDWIDTH_STATS: "/api/system/stats/bandwidth";
            readonly NODES_STATS: "/api/system/stats/nodes";
            readonly NODES_METRICS: "/api/system/nodes/metrics";
        };
        readonly TOOLS: {
            readonly GENERATE_X25519: "/api/system/tools/x25519/generate";
            readonly ENCRYPT_HAPP_CRYPTO_LINK: "/api/system/tools/happ/encrypt";
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
        readonly GET_USER_HWID_DEVICES: (userUuid: string) => string;
        readonly DELETE_USER_HWID_DEVICE: "/api/hwid/devices/delete";
        readonly DELETE_ALL_USER_HWID_DEVICES: "/api/hwid/devices/delete-all";
        readonly STATS: "/api/hwid/devices/stats";
        readonly TOP_USERS_BY_DEVICES: "/api/hwid/devices/top-users";
    };
    readonly SUBSCRIPTIONS: {
        readonly GET: "/api/subscriptions/";
        readonly GET_BY: {
            readonly USERNAME: (username: string) => string;
            readonly UUID: (uuid: string) => string;
            readonly SHORT_UUID: (shortUuid: string) => string;
            readonly SHORT_UUID_RAW: (shortUuid: string) => string;
        };
        readonly SUBPAGE: {
            readonly GET_CONFIG: (shortUuid: string) => string;
        };
        readonly GET_CONNECTION_KEYS_BY_UUID: (uuid: string) => string;
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
    readonly REMNAAWAVE_SETTINGS: {
        readonly GET: "/api/cerberus-settings/";
        readonly UPDATE: "/api/cerberus-settings/";
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
    readonly BANDWIDTH_STATS: {
        readonly NODES: {
            readonly GET: "/api/bandwidth-stats/nodes/";
            readonly GET_REALTIME: "/api/bandwidth-stats/nodes/realtime";
            readonly GET_USERS: (uuid: string) => string;
        };
        readonly USERS: {
            readonly GET_BY_UUID: (uuid: string) => string;
        };
        readonly LEGACY: {
            readonly NODES: {
                readonly GET_USERS: (uuid: string) => string;
            };
            readonly USERS: {
                readonly GET_BY_UUID: (uuid: string) => string;
            };
        };
    };
    readonly IP_CONTROL: {
        readonly FETCH_IPS: (uuid: string) => string;
        readonly GET_FETCH_IPS_RESULT: (jobId: string) => string;
        readonly DROP_CONNECTIONS: "/api/ip-control/drop-connections";
    };
};
//# sourceMappingURL=routes.d.ts.map