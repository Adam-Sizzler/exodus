"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.USERS_ROUTES = exports.USERS_ACTIONS_ROUTE = exports.USERS_CONTROLLER = void 0;
exports.USERS_CONTROLLER = 'users';
exports.USERS_ACTIONS_ROUTE = 'actions';
exports.USERS_ROUTES = {
    CREATE: '',
    UPDATE: '',
    GET: '',
    STREAM: 'stream',
    DELETE: (userId) => `${userId}`,
    GET_BY_ID: (userId) => `${userId}`,
    ACCESSIBLE_NODES: (userId) => `${userId}/accessible-nodes`,
    SUBSCRIPTION_REQUEST_HISTORY: (userId) => `${userId}/subscription-request-history`,
    ACTIONS: {
        ENABLE: (userId) => `${userId}/${exports.USERS_ACTIONS_ROUTE}/enable`,
        DISABLE: (userId) => `${userId}/${exports.USERS_ACTIONS_ROUTE}/disable`,
        RESET_TRAFFIC: (userId) => `${userId}/${exports.USERS_ACTIONS_ROUTE}/reset-traffic`,
        REVOKE_SUBSCRIPTION: (userId) => `${userId}/${exports.USERS_ACTIONS_ROUTE}/revoke`,
        EXTEND_EXPIRATION_DATE: (userId) => `${userId}/${exports.USERS_ACTIONS_ROUTE}/extend`,
    },
    GET_BY: {
        SHORT_UUID: (shortUuid) => `by-short-uuid/${shortUuid}`,
        USERNAME: (username) => `by-username/${username}`,
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
    RESOLVE: 'resolve',
};
