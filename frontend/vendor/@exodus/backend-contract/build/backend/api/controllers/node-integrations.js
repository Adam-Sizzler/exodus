"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.NODE_INTEGRATIONS_ROUTES = exports.NODE_INTEGRATIONS_CONTROLLER = void 0;
exports.NODE_INTEGRATIONS_CONTROLLER = 'node-integrations';
exports.NODE_INTEGRATIONS_ROUTES = {
    GET_ALL: '', // get
    GET: (uuid) => `${uuid}`, // get
    UPDATE: '', // patch
    DELETE: (uuid) => `${uuid}`, // delete
    CREATE: '', // post
};
