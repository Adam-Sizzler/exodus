"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.BANDWIDTH_STATS_ROUTES = exports.BANDWIDTH_STATS_INTERNAL_SQUADS_CONTROLLER = exports.BANDWIDTH_STATS_USERS_CONTROLLER = exports.BANDWIDTH_STATS_NODES_CONTROLLER = exports.BANDWIDTH_STATS_INTERNAL_SQUADS_ROUTE = exports.BANDWIDTH_STATS_USERS_ROUTE = exports.BANDWIDTH_STATS_NODES_ROUTE = exports.BANDWIDTH_STATS_CONTROLLER = void 0;
exports.BANDWIDTH_STATS_CONTROLLER = 'bandwidth-stats';
exports.BANDWIDTH_STATS_NODES_ROUTE = 'nodes';
exports.BANDWIDTH_STATS_USERS_ROUTE = 'users';
exports.BANDWIDTH_STATS_INTERNAL_SQUADS_ROUTE = 'internal-squads';
exports.BANDWIDTH_STATS_NODES_CONTROLLER = `${exports.BANDWIDTH_STATS_CONTROLLER}/${exports.BANDWIDTH_STATS_NODES_ROUTE}`;
exports.BANDWIDTH_STATS_USERS_CONTROLLER = `${exports.BANDWIDTH_STATS_CONTROLLER}/${exports.BANDWIDTH_STATS_USERS_ROUTE}`;
exports.BANDWIDTH_STATS_INTERNAL_SQUADS_CONTROLLER = `${exports.BANDWIDTH_STATS_CONTROLLER}/${exports.BANDWIDTH_STATS_INTERNAL_SQUADS_ROUTE}`;
exports.BANDWIDTH_STATS_ROUTES = {
    NODES: {
        GET: '',
        GET_REALTIME: 'realtime',
        GET_USERS: (uuid) => `${uuid}/users`,
        GET_USERS_BY_NODES: 'users',
        GET_USAGE: 'usage',
    },
    USERS: {
        GET_BY_ID: (userId) => `${userId}`,
    },
    INTERNAL_SQUADS: {
        GET_USAGE: (uuid) => `${uuid}/usage`,
        USER_USAGE: (squadUuid, userId) => `${squadUuid}/users/${userId}/usage`,
    },
};
