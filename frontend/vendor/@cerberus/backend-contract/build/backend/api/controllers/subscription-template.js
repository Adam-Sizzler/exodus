"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.SUBSCRIPTION_TEMPLATE_ROUTES = exports.SUBSCRIPTION_TEMPLATE_CONTROLLER = void 0;
exports.SUBSCRIPTION_TEMPLATE_CONTROLLER = 'subscription-templates';
const ACTIONS_ROUTE = 'actions';
exports.SUBSCRIPTION_TEMPLATE_ROUTES = {
    GET_ALL: '', // get
    GET: (uuid) => `${uuid}`, // get
    UPDATE: '', // patch
    DELETE: (uuid) => `${uuid}`, // delete
    CREATE: '', // post
    ACTIONS: {
        REORDER: `${ACTIONS_ROUTE}/reorder`,
    },
};
