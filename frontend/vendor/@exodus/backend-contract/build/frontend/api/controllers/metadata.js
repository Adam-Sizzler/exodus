"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.METADATA_ROUTES = exports.METADATA_CONTROLLER = void 0;
exports.METADATA_CONTROLLER = 'metadata';
exports.METADATA_ROUTES = {
    NODE: {
        GET: (nodeId) => `node/${nodeId}`,
        UPSERT: (nodeId) => `node/${nodeId}`,
    },
    USER: {
        GET: (userId) => `user/${userId}`,
        UPSERT: (userId) => `user/${userId}`,
    },
};
