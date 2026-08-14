"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.GetStatsNodesUsersUsageCommand = void 0;
const zod_1 = require("zod");
const api_1 = require("../../../api");
const constants_1 = require("../../../constants");
var GetStatsNodesUsersUsageCommand;
(function (GetStatsNodesUsersUsageCommand) {
    GetStatsNodesUsersUsageCommand.url = api_1.REST_API.BANDWIDTH_STATS.NODES.GET_USERS_BY_NODES;
    GetStatsNodesUsersUsageCommand.TSQ_url = GetStatsNodesUsersUsageCommand.url;
    GetStatsNodesUsersUsageCommand.endpointDetails = (0, constants_1.getEndpointDetails)(api_1.BANDWIDTH_STATS_ROUTES.NODES.GET_USERS_BY_NODES, 'post', 'Get Nodes Users Usage by Nodes UUIDs', { scope: 'nodes-users-usage', kind: 'read' });
    GetStatsNodesUsersUsageCommand.RequestBodySchema = zod_1.z.object({
        nodesUuids: zod_1.z.array(zod_1.z.uuid()).min(1),
    });
    GetStatsNodesUsersUsageCommand.RequestQuerySchema = zod_1.z.object({
        start: zod_1.z.iso.date().describe('Start date (YYYY-MM-DD)'),
        end: zod_1.z.iso.date().describe('End date (YYYY-MM-DD)'),
        topUsersLimit: zod_1.z.coerce.number().min(1).default(100),
    });
    GetStatsNodesUsersUsageCommand.ResponseSchema = zod_1.z.object({
        response: zod_1.z.object({
            categories: zod_1.z.array(zod_1.z.string()),
            sparklineData: zod_1.z.array(zod_1.z.number()),
            topUsers: zod_1.z.array(zod_1.z.object({
                color: zod_1.z.string(),
                username: zod_1.z.string(),
                total: zod_1.z.number(),
            })),
        }),
    });
})(GetStatsNodesUsersUsageCommand || (exports.GetStatsNodesUsersUsageCommand = GetStatsNodesUsersUsageCommand = {}));
