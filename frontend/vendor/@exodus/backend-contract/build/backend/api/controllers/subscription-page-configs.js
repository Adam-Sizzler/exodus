"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.SUBSCRIPTION_PAGE_CONFIGS_ROUTES = exports.SUBSCRIPTION_PAGE_CONFIGS_CONTROLLER = void 0;
exports.SUBSCRIPTION_PAGE_CONFIGS_CONTROLLER = 'subscription-page-configs';
const ACTIONS_ROUTE = 'actions';
exports.SUBSCRIPTION_PAGE_CONFIGS_ROUTES = {
    GET_ALL: '', // get
    GET: (uuid) => `${uuid}`, // get
    UPDATE: '', // patch
    DELETE: (uuid) => `${uuid}`, // delete
    CREATE: '', // post
    ACTIONS: {
        REORDER: `${ACTIONS_ROUTE}/reorder`,
        CLONE: `${ACTIONS_ROUTE}/clone`,
    },
};
