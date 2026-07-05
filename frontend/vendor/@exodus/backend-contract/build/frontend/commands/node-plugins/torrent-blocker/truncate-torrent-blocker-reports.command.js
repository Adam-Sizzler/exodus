"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.TruncateTorrentBlockerReportsCommand = void 0;
const zod_1 = require("zod");
const models_1 = require("../../../models");
const api_1 = require("../../../api");
const constants_1 = require("../../../constants");
var TruncateTorrentBlockerReportsCommand;
(function (TruncateTorrentBlockerReportsCommand) {
    TruncateTorrentBlockerReportsCommand.url = api_1.REST_API.NODE_PLUGINS.TORRENT_BLOCKER.TRUNCATE_REPORTS;
    TruncateTorrentBlockerReportsCommand.TSQ_url = TruncateTorrentBlockerReportsCommand.url;
    TruncateTorrentBlockerReportsCommand.endpointDetails = (0, constants_1.getEndpointDetails)(api_1.NODE_PLUGINS_ROUTES.TORRENT_BLOCKER.TRUNCATE_REPORTS, 'delete', 'Truncate Torrent Blocker Reports');
    TruncateTorrentBlockerReportsCommand.ResponseSchema = zod_1.z.object({
        response: zod_1.z.object({
            records: zod_1.z.array(models_1.TorrentBlockerReportSchema),
            total: zod_1.z.number(),
        }),
    });
})(TruncateTorrentBlockerReportsCommand || (exports.TruncateTorrentBlockerReportsCommand = TruncateTorrentBlockerReportsCommand = {}));
