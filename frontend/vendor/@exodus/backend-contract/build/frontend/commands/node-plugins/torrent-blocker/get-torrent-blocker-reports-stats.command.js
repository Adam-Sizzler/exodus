"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.GetTorrentBlockerReportsStatsCommand = void 0;
const zod_1 = require("zod");
const api_1 = require("../../../api");
const constants_1 = require("../../../constants");
var GetTorrentBlockerReportsStatsCommand;
(function (GetTorrentBlockerReportsStatsCommand) {
    GetTorrentBlockerReportsStatsCommand.url = api_1.REST_API.NODE_PLUGINS.TORRENT_BLOCKER.GET_REPORTS_STATS;
    GetTorrentBlockerReportsStatsCommand.TSQ_url = GetTorrentBlockerReportsStatsCommand.url;
    GetTorrentBlockerReportsStatsCommand.endpointDetails = (0, constants_1.getEndpointDetails)(api_1.NODE_PLUGINS_ROUTES.TORRENT_BLOCKER.GET_REPORTS_STATS, 'get', 'Get Torrent Blocker Reports Stats', { scope: 'torrent-blocker-stats', kind: 'read' });
    GetTorrentBlockerReportsStatsCommand.ResponseSchema = zod_1.z.object({
        response: zod_1.z.object({
            stats: zod_1.z.object({
                distinctNodes: zod_1.z.number(),
                distinctUsers: zod_1.z.number(),
                totalReports: zod_1.z.number(),
                reportsLast24Hours: zod_1.z.number(),
            }),
            topUsers: zod_1.z.array(zod_1.z.object({
                userId: zod_1.z.number(),
                color: zod_1.z.string(),
                username: zod_1.z.string(),
                total: zod_1.z.number(),
            })),
            topNodes: zod_1.z.array(zod_1.z.object({
                uuid: zod_1.z.uuid(),
                countryCode: zod_1.z.string(),
                color: zod_1.z.string(),
                name: zod_1.z.string(),
                total: zod_1.z.number(),
            })),
        }),
    });

})(GetTorrentBlockerReportsStatsCommand || (exports.GetTorrentBlockerReportsStatsCommand = GetTorrentBlockerReportsStatsCommand = {}));
