"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.SUBSCRIPTION_ROUTES = exports.SUBSCRIPTION_CONTROLLER = void 0;
exports.SUBSCRIPTION_CONTROLLER = 'sub';
exports.SUBSCRIPTION_ROUTES = {
    GET: '',
    GET_OUTLINE: '/outline',
    GET_INFO: (shortUuid) => `${shortUuid}/info`,
};
