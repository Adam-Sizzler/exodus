"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.TruncateTorrentBlockerReportsCommand = void 0;
const api_1 = require("../../../api");
const constants_1 = require("../../../constants");
var TruncateTorrentBlockerReportsCommand;
(function (TruncateTorrentBlockerReportsCommand) {
    TruncateTorrentBlockerReportsCommand.url = api_1.REST_API.NODE_PLUGINS.TORRENT_BLOCKER.TRUNCATE_REPORTS;
    TruncateTorrentBlockerReportsCommand.TSQ_url = TruncateTorrentBlockerReportsCommand.url;
    TruncateTorrentBlockerReportsCommand.endpointDetails = (0, constants_1.getEndpointDetails)(api_1.NODE_PLUGINS_ROUTES.TORRENT_BLOCKER.TRUNCATE_REPORTS, 'delete', 'Truncate Torrent Blocker Reports', { scope: 'truncate', kind: 'write' });
})(TruncateTorrentBlockerReportsCommand || (exports.TruncateTorrentBlockerReportsCommand = TruncateTorrentBlockerReportsCommand = {}));
