"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.GetTorrentBlockerReportsCommand = void 0;
const zod_1 = require("zod");
const api_1 = require("../../../api");
const constants_1 = require("../../../constants");
const models_1 = require("../../../models");
var GetTorrentBlockerReportsCommand;
(function (GetTorrentBlockerReportsCommand) {
    GetTorrentBlockerReportsCommand.url = api_1.REST_API.NODE_PLUGINS.TORRENT_BLOCKER.GET_REPORTS;
    GetTorrentBlockerReportsCommand.TSQ_url = GetTorrentBlockerReportsCommand.url;
    GetTorrentBlockerReportsCommand.endpointDetails = (0, constants_1.getEndpointDetails)(api_1.NODE_PLUGINS_ROUTES.TORRENT_BLOCKER.GET_REPORTS, 'get', 'Get Torrent Blocker Reports', { scope: 'torrent-blocker-reports', kind: 'read' }, 'Please note that the filters here are primarily intended for use by the frontend and rely on expensive operators such as LIKE under the hood. Misusing these filters may negatively impact the performance of your database.');
    GetTorrentBlockerReportsCommand.RequestQuerySchema = models_1.TanstackQueryRequestQuerySchema;
    GetTorrentBlockerReportsCommand.ResponseSchema = zod_1.z.object({
        response: zod_1.z.object({
            records: zod_1.z.array(models_1.TorrentBlockerReportSchema),
            total: zod_1.z.number(),
        }),
    });

})(GetTorrentBlockerReportsCommand || (exports.GetTorrentBlockerReportsCommand = GetTorrentBlockerReportsCommand = {}));
