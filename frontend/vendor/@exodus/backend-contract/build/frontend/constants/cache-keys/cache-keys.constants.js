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
    SHORT_UUID_RANGE: 'short_uuid_range',
};
exports.CACHE_KEYS_TTL = {
    EXODUS_SETTINGS: 86400000, // 1 day
    EXTERNAL_SQUAD_SETTINGS: 3600000, // 1 hour
    SUBSCRIPTION_SETTINGS: 3600000, // 1 hour
    SHORT_UUID_RANGE: 86400000, // 1 day
};
exports.INTERNAL_CACHE_KEYS = {
    NODE_USER_USAGE_PREFIX: 'node_user_usage:',
    NODE_USER_USAGE: (nodeId) => `${exports.INTERNAL_CACHE_KEYS.NODE_USER_USAGE_PREFIX}${nodeId.toString()}`,
    NODE_USER_USAGE_KEYS: 'node_user_usage_keys',
    PROCESSING_POSTFIX: ':processing',
};
exports.INTERNAL_CACHE_KEYS_TTL = {
    NODE_USER_USAGE: 10800, // 3 hours in seconds
};
