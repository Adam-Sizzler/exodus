"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.GetStatsNodeUsersUsageCommand = void 0;
const zod_1 = require("zod");
const api_1 = require("../../../api");
const constants_1 = require("../../../constants");
var GetStatsNodeUsersUsageCommand;
(function (GetStatsNodeUsersUsageCommand) {
    GetStatsNodeUsersUsageCommand.url = api_1.REST_API.BANDWIDTH_STATS.NODES.GET_USERS;
    GetStatsNodeUsersUsageCommand.TSQ_url = GetStatsNodeUsersUsageCommand.url(':uuid');
    GetStatsNodeUsersUsageCommand.endpointDetails = (0, constants_1.getEndpointDetails)(api_1.BANDWIDTH_STATS_ROUTES.NODES.GET_USERS(':uuid'), 'get', 'Get Node Users Usage by Node UUID');
    GetStatsNodeUsersUsageCommand.RequestQuerySchema = zod_1.z.object({
        start: zod_1.z.string().date(),
        end: zod_1.z.string().date(),
        topUsersLimit: zod_1.z.coerce.number().min(1).default(100),
    });
    GetStatsNodeUsersUsageCommand.RequestSchema = zod_1.z.object({
        uuid: zod_1.z.string().uuid(),
    });
    GetStatsNodeUsersUsageCommand.ResponseSchema = zod_1.z.object({
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
})(GetStatsNodeUsersUsageCommand || (exports.GetStatsNodeUsersUsageCommand = GetStatsNodeUsersUsageCommand = {}));
