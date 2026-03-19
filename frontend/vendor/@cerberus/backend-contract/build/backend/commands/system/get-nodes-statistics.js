"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.GetNodesStatisticsCommand = void 0;
const zod_1 = require("zod");
const constants_1 = require("../../constants");
const api_1 = require("../../api");
var GetNodesStatisticsCommand;
(function (GetNodesStatisticsCommand) {
    GetNodesStatisticsCommand.url = api_1.REST_API.SYSTEM.STATS.NODES_STATS;
    GetNodesStatisticsCommand.TSQ_url = GetNodesStatisticsCommand.url;
    GetNodesStatisticsCommand.endpointDetails = (0, constants_1.getEndpointDetails)(api_1.SYSTEM_ROUTES.STATS.NODES_STATS, 'get', 'Get Nodes Statistics');
    GetNodesStatisticsCommand.RequestQuerySchema = zod_1.z.object({
        tz: zod_1.z.string().optional(),
    });
    GetNodesStatisticsCommand.ResponseSchema = zod_1.z.object({
        response: zod_1.z.object({
            lastSevenDays: zod_1.z.array(zod_1.z.object({
                nodeName: zod_1.z.string(),
                date: zod_1.z.string(),
                totalBytes: zod_1.z.string(),
            })),
        }),
    });
})(GetNodesStatisticsCommand || (exports.GetNodesStatisticsCommand = GetNodesStatisticsCommand = {}));
