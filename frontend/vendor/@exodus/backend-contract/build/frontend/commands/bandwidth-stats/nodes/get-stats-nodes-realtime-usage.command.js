"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.GetStatsNodesRealtimeUsageCommand = void 0;
const zod_1 = require("zod");
const api_1 = require("../../../api");
const constants_1 = require("../../../constants");
var GetStatsNodesRealtimeUsageCommand;
(function (GetStatsNodesRealtimeUsageCommand) {
    GetStatsNodesRealtimeUsageCommand.url = api_1.REST_API.BANDWIDTH_STATS.NODES.GET_REALTIME;
    GetStatsNodesRealtimeUsageCommand.TSQ_url = GetStatsNodesRealtimeUsageCommand.url;
    GetStatsNodesRealtimeUsageCommand.endpointDetails = (0, constants_1.getEndpointDetails)(api_1.BANDWIDTH_STATS_ROUTES.NODES.GET_REALTIME, 'get', 'Get Nodes Realtime Usage');
    GetStatsNodesRealtimeUsageCommand.ResponseSchema = zod_1.z.object({
        response: zod_1.z.array(zod_1.z.object({
            nodeUuid: zod_1.z.string().uuid(),
            nodeName: zod_1.z.string(),
            countryCode: zod_1.z.string(),
            downloadBytes: zod_1.z.number(),
            uploadBytes: zod_1.z.number(),
            totalBytes: zod_1.z.number(),
            downloadSpeedBps: zod_1.z.number(),
            uploadSpeedBps: zod_1.z.number(),
            totalSpeedBps: zod_1.z.number(),
        })),
    });
})(GetStatsNodesRealtimeUsageCommand || (exports.GetStatsNodesRealtimeUsageCommand = GetStatsNodesRealtimeUsageCommand = {}));
