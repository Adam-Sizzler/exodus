"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.GetStatsNodesUsageCommand = void 0;
const zod_1 = require("zod");
const api_1 = require("../../../api");
const constants_1 = require("../../../constants");
var GetStatsNodesUsageCommand;
(function (GetStatsNodesUsageCommand) {
    GetStatsNodesUsageCommand.url = api_1.REST_API.BANDWIDTH_STATS.NODES.GET;
    GetStatsNodesUsageCommand.TSQ_url = GetStatsNodesUsageCommand.url;
    GetStatsNodesUsageCommand.endpointDetails = (0, constants_1.getEndpointDetails)(api_1.BANDWIDTH_STATS_ROUTES.NODES.GET, 'get', 'Get Nodes Usage by Range', { scope: 'nodes-usage', kind: 'read' });
    GetStatsNodesUsageCommand.RequestQuerySchema = zod_1.z.object({
        start: zod_1.z.iso.date().describe('Start date (YYYY-MM-DD)'),
        end: zod_1.z.iso.date().describe('End date (YYYY-MM-DD)'),
        topNodesLimit: zod_1.z.coerce
            .number()
            .min(1)
            .default(20)
            .describe('Limit of top nodes to return'),
    });
    GetStatsNodesUsageCommand.ResponseSchema = zod_1.z.object({
        response: zod_1.z.object({
            categories: zod_1.z.array(zod_1.z.string()),
            sparklineData: zod_1.z.array(zod_1.z.number()),
            topNodes: zod_1.z.array(zod_1.z.object({
                uuid: zod_1.z.uuid(),
                color: zod_1.z.string(),
                name: zod_1.z.string(),
                countryCode: zod_1.z.string(),
                total: zod_1.z.number(),
            })),
            series: zod_1.z.array(zod_1.z.object({
                uuid: zod_1.z.uuid(),
                name: zod_1.z.string(),
                color: zod_1.z.string(),
                countryCode: zod_1.z.string(),
                total: zod_1.z.number(),
                data: zod_1.z.array(zod_1.z.number()),
            })),
        }),
    });

})(GetStatsNodesUsageCommand || (exports.GetStatsNodesUsageCommand = GetStatsNodesUsageCommand = {}));
