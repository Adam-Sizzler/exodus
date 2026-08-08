"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.GetBandwidthStatsCommand = void 0;
const zod_1 = require("zod");
const api_1 = require("../../api");
const constants_1 = require("../../constants");
const base_stat_schema_1 = require("../../models/base-stat.schema");
var GetBandwidthStatsCommand;
(function (GetBandwidthStatsCommand) {
    GetBandwidthStatsCommand.url = api_1.REST_API.SYSTEM.STATS.BANDWIDTH_STATS;
    GetBandwidthStatsCommand.TSQ_url = GetBandwidthStatsCommand.url;
    GetBandwidthStatsCommand.endpointDetails = (0, constants_1.getEndpointDetails)(api_1.SYSTEM_ROUTES.STATS.BANDWIDTH_STATS, 'get', 'Get Bandwidth Stats', { scope: 'bandwidth-stats', kind: 'read' });
    GetBandwidthStatsCommand.RequestQuerySchema = zod_1.z.object({
        tz: zod_1.z.string().optional(),
    });
    GetBandwidthStatsCommand.ResponseSchema = zod_1.z.object({
        response: zod_1.z.object({
            bandwidthLastTwoDays: base_stat_schema_1.BaseStatSchema,
            bandwidthLastSevenDays: base_stat_schema_1.BaseStatSchema,
            bandwidthLast30Days: base_stat_schema_1.BaseStatSchema,
            bandwidthCalendarMonth: base_stat_schema_1.BaseStatSchema,
            bandwidthCurrentYear: base_stat_schema_1.BaseStatSchema,
        }),
    });

})(GetBandwidthStatsCommand || (exports.GetBandwidthStatsCommand = GetBandwidthStatsCommand = {}));
