"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.CONNECTIONS_ROUTES = exports.CONNECTIONS_CONTROLLER = void 0;
exports.CONNECTIONS_CONTROLLER = 'connections';
exports.CONNECTIONS_ROUTES = {
    // POST
    CONNECTIONS_BY_USER: (userId) => `by-user/${userId}`,
    // GET
    CONNECTIONS_BY_USER_RESULT: (jobId) => `by-user/${jobId}`,
    // POST
    CONNECTIONS_BY_NODE: (uuid) => `by-node/${uuid}`,
    // GET
    CONNECTIONS_BY_NODE_RESULT: (jobId) => `by-node/${jobId}`,
    // POST
    DROP_CONNECTIONS: 'drop',
};
