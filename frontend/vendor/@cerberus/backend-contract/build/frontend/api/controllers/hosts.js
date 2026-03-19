"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.HOSTS_ROUTES = exports.HOST_ACTIONS_ROUTE = exports.HOSTS_CONTROLLER = void 0;
exports.HOSTS_CONTROLLER = 'hosts';
exports.HOST_ACTIONS_ROUTE = 'actions';
exports.HOSTS_ROUTES = {
    CREATE: '', // create
    UPDATE: '', // update host
    GET: '', // get all hosts
    GET_BY_UUID: (uuid) => `${uuid}`, // get by UUID
    DELETE: (uuid) => `${uuid}`, // delete by UUID
    ACTIONS: {
        REORDER: `${exports.HOST_ACTIONS_ROUTE}/reorder`,
    },
    BULK: {
        ENABLE_HOSTS: 'bulk/enable',
        DISABLE_HOSTS: 'bulk/disable',
        DELETE_HOSTS: 'bulk/delete',
        SET_INBOUND: 'bulk/set-inbound',
        SET_PORT: 'bulk/set-port',
    },
    TAGS: {
        GET: 'tags',
    },
};
