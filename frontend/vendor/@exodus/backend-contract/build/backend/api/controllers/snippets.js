"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.SNIPPETS_ROUTES = exports.SNIPPETS_CONTROLLER = void 0;
exports.SNIPPETS_CONTROLLER = 'snippets';
const ACTIONS_ROUTE = 'actions';
exports.SNIPPETS_ROUTES = {
    GET: '', // Get list of all snippets // get
    CREATE: '', // Create new snippet // post
    UPDATE: '', // Update snippet by name // patch
    DELETE: '', // Delete snippet by name // delete
    ACTIONS: {
        SYNC: `${ACTIONS_ROUTE}/sync`, // post
    },
};
