"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.NODE_SSH_ROUTES = exports.NODE_SSH_CONTROLLER = void 0;
exports.NODE_SSH_CONTROLLER = 'node-ssh';
exports.NODE_SSH_ROUTES = {
    CREATE_TICKET: (uuid) => `${uuid}/ticket`, // post
    EVALUATE_VAULT: 'vault/evaluate', // post
};
