"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.INTERNAL_CACHE_KEYS_TTL = exports.INTERNAL_CACHE_KEYS = exports.CACHE_KEYS_TTL = exports.CACHE_KEYS = void 0;
exports.CACHE_KEYS = {
    SUBSCRIPTION_SETTINGS: 'subscription_settings',
    EXTERNAL_SQUAD_SETTINGS: (uuid) => `external_squad_settings:${uuid}`,
    SUBSCRIPTION_TEMPLATE: (name, type) => `subscription_template:${name}:${type}`,
    PASSKEY_REGISTRATION_OPTIONS: (uuid) => `passkey_registration_options:${uuid}`,
    PASSKEY_AUTHENTICATION_OPTIONS: (uuid) => `passkey_authentication_options:${uuid}`,
    EXODUS_SETTINGS: 'exodus_settings',
    NODE_SYSTEM_INFO: (uuid) => `node_system_info:${uuid}`,
    NODE_SYSTEM_STATS: (uuid) => `node_system_stats:${uuid}`,
    NODE_USERS_ONLINE: (uuid) => `node_users_online:${uuid}`,
    NODE_VERSIONS: (uuid) => `node_versions:${uuid}`,
    NODE_XRAY_UPTIME: (uuid) => `node_xray_uptime:${uuid}`,
    RAW_INBOUND: (uuid) => `raw_inbound:${uuid}`,
    XRAY_JSON_TEMPLATE: (uuid) => `xray_json_template:${uuid}`,
};
exports.CACHE_KEYS_TTL = {
    EXODUS_SETTINGS: 86400, // 1 day
    EXTERNAL_SQUAD_SETTINGS: 3600, // 1 hour
    SUBSCRIPTION_SETTINGS: 3600, // 1 hour
    NODE_SYSTEM_STATS: 30, // 30 seconds
    NODE_USERS_ONLINE: 16, // 16 seconds
    NODE_XRAY_UPTIME: 16, // 16 seconds
    RAW_INBOUND: 3600, // 1 hour
    XRAY_JSON_TEMPLATE: 3600, // 1 hour
};
exports.INTERNAL_CACHE_KEYS = {
    NODE_USER_USAGE_PREFIX: 'node_user_usage:',
    NODE_USER_USAGE: (nodeId) => `${exports.INTERNAL_CACHE_KEYS.NODE_USER_USAGE_PREFIX}${nodeId.toString()}`,
    NODE_USER_USAGE_KEYS: 'node_user_usage_keys',
    PROCESSING_POSTFIX: ':processing',
    RUNTIME_METRICS: 'runtime_metrics',
};
exports.INTERNAL_CACHE_KEYS_TTL = {
    NODE_USER_USAGE: 10800, // 3 hours in seconds
};
