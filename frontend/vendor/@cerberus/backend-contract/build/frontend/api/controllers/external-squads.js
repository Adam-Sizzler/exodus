"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.EXTERNAL_SQUADS_ROUTES = exports.EXTERNAL_SQUADS_CONTROLLER = void 0;
exports.EXTERNAL_SQUADS_CONTROLLER = 'external-squads';
const BULK_ACTIONS_ROUTE = 'bulk-actions';
const ACTIONS_ROUTE = 'actions';
exports.EXTERNAL_SQUADS_ROUTES = {
    GET: '', // Get list of all external squads // get
    CREATE: '', // Create new external squad // post
    UPDATE: '', // Update external squad by uuid // patch
    GET_BY_UUID: (uuid) => `${uuid}`, // Get external squad by uuid // get
    DELETE: (uuid) => `${uuid}`, // Delete external squad by uuid // delete
    BULK_ACTIONS: {
        ADD_USERS: (uuid) => `${uuid}/${BULK_ACTIONS_ROUTE}/add-users`, // Add users to external squad // post
        REMOVE_USERS: (uuid) => `${uuid}/${BULK_ACTIONS_ROUTE}/remove-users`, // Remove users from external squad // delete
    },
    ACTIONS: {
        REORDER: `${ACTIONS_ROUTE}/reorder`,
    },
};
