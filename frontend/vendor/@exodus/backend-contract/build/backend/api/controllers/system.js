"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.SYSTEM_ROUTES = exports.SYSTEM_CONTROLLER = void 0;
exports.SYSTEM_CONTROLLER = 'system';
exports.SYSTEM_ROUTES = {
    STATS: {
        SYSTEM_STATS: 'stats',
        BANDWIDTH_STATS: 'stats/bandwidth',
        NODES_STATS: 'stats/nodes',
        NODES_METRICS: 'nodes/metrics',
    },
    TOOLS: {
        GENERATE_X25519: 'tools/x25519/generate',
        ENCRYPT_HAPP_CRYPTO_LINK: 'tools/happ/encrypt',
    },
    HEALTH: 'health',
    METADATA: 'metadata',
    TESTERS: {
        SRR_MATCHER: 'testers/srr-matcher',
    },
};
