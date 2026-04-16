"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.IP_CONTROL_ROUTES = exports.IP_CONTROL_CONTROLLER = void 0;
exports.IP_CONTROL_CONTROLLER = 'ip-control';
exports.IP_CONTROL_ROUTES = {
    // POST /ip-control/fetch-ips/:userUuid
    FETCH_IPS: (uuid) => `fetch-ips/${uuid}`,
    // GET /ip-control/fetch-ips/result/:jobId
    GET_FETCH_IPS_RESULT: (jobId) => `fetch-ips/result/${jobId}`,
    // POST /ip-control/drop-connections
    DROP_CONNECTIONS: 'drop-connections',
};
