"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.GetNodesMetricsCommand = void 0;
const zod_1 = require("zod");
const constants_1 = require("../../constants");
const api_1 = require("../../api");
var GetNodesMetricsCommand;
(function (GetNodesMetricsCommand) {
    GetNodesMetricsCommand.url = api_1.REST_API.SYSTEM.STATS.NODES_METRICS;
    GetNodesMetricsCommand.TSQ_url = GetNodesMetricsCommand.url;
    GetNodesMetricsCommand.endpointDetails = (0, constants_1.getEndpointDetails)(api_1.SYSTEM_ROUTES.STATS.NODES_METRICS, 'get', 'Get Nodes Metrics');
    GetNodesMetricsCommand.ResponseSchema = zod_1.z.object({
        response: zod_1.z.object({
            nodes: zod_1.z.array(zod_1.z.object({
                nodeUuid: zod_1.z.string(),
                nodeName: zod_1.z.string(),
                countryEmoji: zod_1.z.string(),
                providerName: zod_1.z.string(),
                usersOnline: zod_1.z.number(),
                inboundsStats: zod_1.z.array(zod_1.z.object({
                    tag: zod_1.z.string(),
                    upload: zod_1.z.string(),
                    download: zod_1.z.string(),
                })),
                outboundsStats: zod_1.z.array(zod_1.z.object({
                    tag: zod_1.z.string(),
                    upload: zod_1.z.string(),
                    download: zod_1.z.string(),
                })),
            })),
        }),
    });
})(GetNodesMetricsCommand || (exports.GetNodesMetricsCommand = GetNodesMetricsCommand = {}));
