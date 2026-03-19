"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.INTERNAL_SQUADS_ROUTES = exports.INTERNAL_SQUADS_CONTROLLER = void 0;
exports.INTERNAL_SQUADS_CONTROLLER = 'internal-squads';
const BULK_ACTIONS_ROUTE = 'bulk-actions';
const ACTIONS_ROUTE = 'actions';
exports.INTERNAL_SQUADS_ROUTES = {
    GET: '', // Get list of all internal squads // get
    CREATE: '', // Create new internal squad // post
    UPDATE: '', // Update internal squad by uuid // patch
    GET_BY_UUID: (uuid) => `${uuid}`, // Get internal squad by uuid // get
    DELETE: (uuid) => `${uuid}`, // Delete internal squad by uuid // delete
    ACCESSIBLE_NODES: (uuid) => `${uuid}/accessible-nodes`, // Get accessible nodes for internal squad // get
    BULK_ACTIONS: {
        ADD_USERS: (uuid) => `${uuid}/${BULK_ACTIONS_ROUTE}/add-users`, // Add users to internal squad // post
        REMOVE_USERS: (uuid) => `${uuid}/${BULK_ACTIONS_ROUTE}/remove-users`, // Remove users from internal squad // delete
    },
    ACTIONS: {
        REORDER: `${ACTIONS_ROUTE}/reorder`,
    },
};
