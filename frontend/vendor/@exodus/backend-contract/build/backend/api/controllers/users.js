"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.USERS_ROUTES = exports.USERS_ACTIONS_ROUTE = exports.USERS_CONTROLLER = void 0;
exports.USERS_CONTROLLER = 'users';
exports.USERS_ACTIONS_ROUTE = 'actions';
exports.USERS_ROUTES = {
    CREATE: '',
    UPDATE: '',
    GET: '',
    DELETE: (uuid) => `${uuid}`,
    GET_BY_UUID: (uuid) => `${uuid}`,
    ACCESSIBLE_NODES: (uuid) => `${uuid}/accessible-nodes`,
    SUBSCRIPTION_REQUEST_HISTORY: (uuid) => `${uuid}/subscription-request-history`,
    ACTIONS: {
        ENABLE: (uuid) => `${uuid}/${exports.USERS_ACTIONS_ROUTE}/enable`,
        DISABLE: (uuid) => `${uuid}/${exports.USERS_ACTIONS_ROUTE}/disable`,
        RESET_TRAFFIC: (uuid) => `${uuid}/${exports.USERS_ACTIONS_ROUTE}/reset-traffic`,
        REVOKE_SUBSCRIPTION: (uuid) => `${uuid}/${exports.USERS_ACTIONS_ROUTE}/revoke`,
    },
    GET_BY: {
        ID: (id) => `by-id/${id}`,
        SHORT_UUID: (shortUuid) => `by-short-uuid/${shortUuid}`,
        USERNAME: (username) => `by-username/${username}`,
        SUBSCRIPTION_UUID: (subscriptionUuid) => `by-subscription-uuid/${subscriptionUuid}`,
        TELEGRAM_ID: (telegramId) => `by-telegram-id/${telegramId}`,
        EMAIL: (email) => `by-email/${email}`,
        TAG: (tag) => `by-tag/${tag}`,
    },
    BULK: {
        DELETE_BY_STATUS: 'bulk/delete-by-status',
        UPDATE: 'bulk/update',
        RESET_TRAFFIC: 'bulk/reset-traffic',
        REVOKE_SUBSCRIPTION: 'bulk/revoke-subscription',
        DELETE: 'bulk/delete',
        UPDATE_SQUADS: 'bulk/update-squads',
        EXTEND_EXPIRATION_DATE: 'bulk/extend-expiration-date',
        ALL: {
            UPDATE: 'bulk/all/update',
            RESET_TRAFFIC: 'bulk/all/reset-traffic',
            EXTEND_EXPIRATION_DATE: 'bulk/all/extend-expiration-date',
        },
    },
    TAGS: {
        GET: 'tags',
    },
};
