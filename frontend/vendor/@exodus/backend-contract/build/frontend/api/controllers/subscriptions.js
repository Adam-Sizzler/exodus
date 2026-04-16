"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.SUBSCRIPTIONS_ROUTES = exports.SUBSCRIPTIONS_CONTROLLER = void 0;
exports.SUBSCRIPTIONS_CONTROLLER = 'subscriptions';
exports.SUBSCRIPTIONS_ROUTES = {
    GET: '',
    GET_BY: {
        USERNAME: (username) => `by-username/${username}`,
        UUID: (uuid) => `by-uuid/${uuid}`,
        SHORT_UUID: (shortUuid) => `by-short-uuid/${shortUuid}`,
        SHORT_UUID_RAW: (shortUuid) => `by-short-uuid/${shortUuid}/raw`,
    },
    SUBPAGE: {
        GET_CONFIG: (shortUuid) => `subpage-config/${shortUuid}`,
    },
};
