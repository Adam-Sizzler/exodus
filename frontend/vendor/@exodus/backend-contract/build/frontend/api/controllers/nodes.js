"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.NODES_ROUTES = exports.NODE_ACTIONS_ROUTE = exports.NODES_CONTROLLER = void 0;
exports.NODES_CONTROLLER = 'nodes';
exports.NODE_ACTIONS_ROUTE = 'actions';
const BULK_ACTIONS_ROUTE = 'bulk-actions';
exports.NODES_ROUTES = {
    CREATE: '', // create
    GET: '', // get all nodes
    GET_BY_UUID: (uuid) => `${uuid}`, // get by UUID
    UPDATE: '', // update, patch
    DELETE: (uuid) => `${uuid}`, // delete by UUID
    ACTIONS: {
        ENABLE: (uuid) => `${uuid}/${exports.NODE_ACTIONS_ROUTE}/enable`,
        DISABLE: (uuid) => `${uuid}/${exports.NODE_ACTIONS_ROUTE}/disable`,
        RESTART: (uuid) => `${uuid}/${exports.NODE_ACTIONS_ROUTE}/restart`,
        RESET_TRAFFIC: (uuid) => `${uuid}/${exports.NODE_ACTIONS_ROUTE}/reset-traffic`,
        RESTART_ALL: `${exports.NODE_ACTIONS_ROUTE}/restart-all`,
        REORDER: `${exports.NODE_ACTIONS_ROUTE}/reorder`,
    },
    BULK_ACTIONS: {
        PROFILE_MODIFICATION: `${BULK_ACTIONS_ROUTE}/profile-modification`,
    },
    TAGS: {
        GET: 'tags',
    },
};
