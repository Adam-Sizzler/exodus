import { TSubscriptionTemplateType } from '../subscription-template';
export declare const CACHE_KEYS: {
    readonly SUBSCRIPTION_SETTINGS: "subscription_settings";
    readonly EXTERNAL_SQUAD_SETTINGS: (uuid: string) => string;
    readonly SUBSCRIPTION_TEMPLATE: (name: string, type: TSubscriptionTemplateType) => string;
    readonly PASSKEY_REGISTRATION_OPTIONS: (uuid: string) => string;
    readonly PASSKEY_AUTHENTICATION_OPTIONS: (uuid: string) => string;
    readonly EXODUS_SETTINGS: "exodus_settings";
    readonly NODE_SYSTEM_INFO: (uuid: string) => string;
    readonly NODE_SYSTEM_STATS: (uuid: string) => string;
    readonly NODE_USERS_ONLINE: (uuid: string) => string;
    readonly NODE_VERSIONS: (uuid: string) => string;
    readonly NODE_XRAY_UPTIME: (uuid: string) => string;
    readonly RAW_INBOUND: (uuid: string) => string;
    readonly XRAY_JSON_TEMPLATE: (uuid: string) => string;
    readonly EXTERNAL_SQUAD_TEMPLATE_NAME: (uuid: string, type: TSubscriptionTemplateType) => string;
};
export declare const CACHE_KEYS_TTL: {
    readonly EXODUS_SETTINGS: 86400;
    readonly EXTERNAL_SQUAD_SETTINGS: 3600;
    readonly SUBSCRIPTION_SETTINGS: 3600;
    readonly NODE_SYSTEM_STATS: 30;
    readonly NODE_USERS_ONLINE: 16;
    readonly NODE_XRAY_UPTIME: 16;
    readonly RAW_INBOUND: 3600;
    readonly XRAY_JSON_TEMPLATE: 3600;
    readonly EXTERNAL_SQUAD_TEMPLATE_NAME: 3600;
};
export declare const INTERNAL_CACHE_KEYS: {
    readonly NODE_USER_USAGE_PREFIX: "node_user_usage:";
    readonly NODE_USER_USAGE: (nodeId: bigint) => string;
    readonly NODE_USER_USAGE_KEYS: "node_user_usage_keys";
    readonly PROCESSING_POSTFIX: ":processing";
    readonly RUNTIME_METRICS: "runtime_metrics";
};
export declare const INTERNAL_CACHE_KEYS_TTL: {
    readonly NODE_USER_USAGE: 10800;
};
export declare const EXPORT_TO_STREAM_KEYS: {
    readonly PREFIX: "ioraw:";
    readonly USER_USAGE: "export:user_usage";
    readonly SUBSCRIPTION_REQUESTS: "export:subscription_requests";
    readonly NODE_CONNECTIONS: "export:node_connections";
};
//# sourceMappingURL=cache-keys.constants.d.ts.map