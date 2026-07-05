"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.NODE_PLUGINS_ROUTES = exports.NODE_PLUGINS_CONTROLLER = void 0;
exports.NODE_PLUGINS_CONTROLLER = 'node-plugins';
const ACTIONS_ROUTE = 'actions';
const TORRENT_BLOCKER_ROUTE = 'torrent-blocker';
exports.NODE_PLUGINS_ROUTES = {
    GET_ALL: '', // get
    GET: (uuid) => `${uuid}`, // get
    UPDATE: '', // patch
    DELETE: (uuid) => `${uuid}`, // delete
    CREATE: '', // post,
    ACTIONS: {
        REORDER: `${ACTIONS_ROUTE}/reorder`,
        CLONE: `${ACTIONS_ROUTE}/clone`,
    },
    EXECUTOR: 'executor',
    TORRENT_BLOCKER: {
        GET_REPORTS: `${TORRENT_BLOCKER_ROUTE}`,
        GET_REPORTS_STATS: `${TORRENT_BLOCKER_ROUTE}/stats`,
        TRUNCATE_REPORTS: `${TORRENT_BLOCKER_ROUTE}/truncate`,
    },
};
