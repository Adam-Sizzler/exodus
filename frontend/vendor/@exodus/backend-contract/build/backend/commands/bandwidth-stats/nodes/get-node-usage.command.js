"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.GetNodeUsageCommand = void 0;
const zod_1 = require("zod");
const api_1 = require("../../../api");
const constants_1 = require("../../../constants");
var GetNodeUsageCommand;
(function (GetNodeUsageCommand) {
    GetNodeUsageCommand.url = api_1.REST_API.BANDWIDTH_STATS.NODES.GET_USAGE;
    GetNodeUsageCommand.TSQ_url = GetNodeUsageCommand.url;
    GetNodeUsageCommand.endpointDetails = (0, constants_1.getEndpointDetails)(api_1.BANDWIDTH_STATS_ROUTES.NODES.GET_USAGE, 'post', 'Get users exceeding a traffic threshold on the given nodes for a period', { scope: 'node-usage', kind: 'read' }, 'Returns users whose total usage over the period on the given nodes is >= minTotalBytes. Underlying usage data is flushed to the database roughly every 2 minutes.');
    GetNodeUsageCommand.RequestBodySchema = zod_1.z.object({
        nodesUuids: zod_1.z.array(zod_1.z.uuid()).min(1).describe('Node UUIDs to include'),
    });
    GetNodeUsageCommand.RequestQuerySchema = zod_1.z.object({
        start: zod_1.z.iso.date().describe('Start date (YYYY-MM-DD)'),
        end: zod_1.z.iso.date().describe('End date (YYYY-MM-DD)'),
        minTotalBytes: zod_1.z.coerce
            .number()
            .min(0)
            .optional()
            .default(0)
            .describe('Only include users whose total usage over the period is >= this (bytes)'),
    });
    GetNodeUsageCommand.ResponseSchema = zod_1.z.object({
        response: zod_1.z.object({
            nodes: zod_1.z.array(zod_1.z.object({
                uuid: zod_1.z.uuid(),
                users: zod_1.z.array(zod_1.z.object({
                    id: zod_1.z.number(),
                    totalBytes: zod_1.z
                        .number()
                        .describe('Total used bytes over the period (raw bytes)'),
                })),
            })),
        }),
    });
})(GetNodeUsageCommand || (exports.GetNodeUsageCommand = GetNodeUsageCommand = {}));
