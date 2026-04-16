"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.HWID_ROUTES = exports.HWID_CONTROLLER = void 0;
exports.HWID_CONTROLLER = 'hwid';
exports.HWID_ROUTES = {
    GET_ALL_HWID_DEVICES: 'devices', // get
    CREATE_USER_HWID_DEVICE: 'devices',
    GET_USER_HWID_DEVICES: (userUuid) => `devices/${userUuid}`,
    DELETE_USER_HWID_DEVICE: 'devices/delete',
    DELETE_ALL_USER_HWID_DEVICES: 'devices/delete-all',
    STATS: 'devices/stats', // get
    TOP_USERS_BY_DEVICES: 'devices/top-users', // get
};
