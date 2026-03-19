"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.BANDWIDTH_STATS_ROUTES = exports.BANDWIDTH_STATS_USERS_CONTROLLER = exports.BANDWIDTH_STATS_NODES_CONTROLLER = exports.BANDWIDTH_STATS_USERS_ROUTE = exports.BANDWIDTH_STATS_NODES_ROUTE = exports.BANDWIDTH_STATS_CONTROLLER = void 0;
exports.BANDWIDTH_STATS_CONTROLLER = 'bandwidth-stats';
exports.BANDWIDTH_STATS_NODES_ROUTE = 'nodes';
exports.BANDWIDTH_STATS_USERS_ROUTE = 'users';
exports.BANDWIDTH_STATS_NODES_CONTROLLER = `${exports.BANDWIDTH_STATS_CONTROLLER}/${exports.BANDWIDTH_STATS_NODES_ROUTE}`;
exports.BANDWIDTH_STATS_USERS_CONTROLLER = `${exports.BANDWIDTH_STATS_CONTROLLER}/${exports.BANDWIDTH_STATS_USERS_ROUTE}`;
// Variants:
// 1. Nodes -> Metrics
// 2. Nodes -> Management -> Show usage (!need legacy)
// 3. Users -> User -> Show Usage (!need legacy)
exports.BANDWIDTH_STATS_ROUTES = {
    NODES: {
        // GET /bandwidth-stats/nodes –– Nodes -> Metrics
        GET: '',
        // GET /bandwidth-stats/nodes/realtime –– Nodes -> Management -> Metric Cards
        GET_REALTIME: 'realtime',
        // GET /bandwidth-stats/nodes/:nodeUuid/users –– Nodes -> Management -> Show usage
        GET_USERS: (uuid) => `${uuid}/users`,
    },
    USERS: {
        // GET /bandwidth-stats/users/:userUuid –– Users -> User -> Show Usage
        GET_BY_UUID: (uuid) => `${uuid}`,
    },
    LEGACY: {
        NODES: {
            // GET /bandwidth-stats/nodes/:nodeUuid/users/legacy –– Nodes -> Management -> Show usage (legacy)
            GET_USERS: (uuid) => `${uuid}/users/legacy`,
        },
        USERS: {
            // GET /bandwidth-stats/users/:userUuid/legacy –– Users -> User -> Show Usage (legacy)
            GET_BY_UUID: (uuid) => `${uuid}/legacy`,
        },
    },
};
