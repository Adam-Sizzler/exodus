import { TSubscriptionTemplateType } from '../subscription-template';
export declare const CACHE_KEYS: {
    readonly SUBSCRIPTION_SETTINGS: "subscription_settings";
    readonly EXTERNAL_SQUAD_SETTINGS: (uuid: string) => string;
    readonly SUBSCRIPTION_TEMPLATE: (name: string, type: TSubscriptionTemplateType) => string;
    readonly PASSKEY_REGISTRATION_OPTIONS: (uuid: string) => string;
    readonly PASSKEY_AUTHENTICATION_OPTIONS: (uuid: string) => string;
    readonly EXODUS_SETTINGS: "exodus_settings";
    readonly SHORT_UUID_RANGE: "short_uuid_range";
};
export declare const CACHE_KEYS_TTL: {
    readonly EXODUS_SETTINGS: 86400000;
    readonly EXTERNAL_SQUAD_SETTINGS: 3600000;
    readonly SUBSCRIPTION_SETTINGS: 3600000;
    readonly SHORT_UUID_RANGE: 86400000;
};
export declare const INTERNAL_CACHE_KEYS: {
    readonly NODE_USER_USAGE_PREFIX: "node_user_usage:";
    readonly NODE_USER_USAGE: (nodeId: bigint) => string;
    readonly NODE_USER_USAGE_KEYS: "node_user_usage_keys";
    readonly PROCESSING_POSTFIX: ":processing";
};
export declare const INTERNAL_CACHE_KEYS_TTL: {
    readonly NODE_USER_USAGE: 10800;
};
//# sourceMappingURL=cache-keys.constants.d.ts.map