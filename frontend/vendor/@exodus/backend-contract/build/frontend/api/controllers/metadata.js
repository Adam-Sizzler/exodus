"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.METADATA_ROUTES = exports.METADATA_CONTROLLER = void 0;
exports.METADATA_CONTROLLER = 'metadata';
exports.METADATA_ROUTES = {
    NODE: {
        GET: (uuid) => `node/${uuid}`,
        UPSERT: (uuid) => `node/${uuid}`,
    },
    USER: {
        GET: (uuid) => `user/${uuid}`,
        UPSERT: (uuid) => `user/${uuid}`,
    },
};
