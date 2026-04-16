"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.GetStatsCommand = void 0;
const zod_1 = require("zod");
const constants_1 = require("../../constants");
const api_1 = require("../../api");
const constants_2 = require("../../constants");
var GetStatsCommand;
(function (GetStatsCommand) {
    GetStatsCommand.url = api_1.REST_API.SYSTEM.STATS.SYSTEM_STATS;
    GetStatsCommand.TSQ_url = GetStatsCommand.url;
    GetStatsCommand.endpointDetails = (0, constants_1.getEndpointDetails)(api_1.SYSTEM_ROUTES.STATS.SYSTEM_STATS, 'get', 'Get Stats');
    GetStatsCommand.RequestQuerySchema = zod_1.z.object({
        tz: zod_1.z.string().optional(),
    });
    GetStatsCommand.ResponseSchema = zod_1.z.object({
        response: zod_1.z.object({
            cpu: zod_1.z.object({
                cores: zod_1.z.number(),
                physicalCores: zod_1.z.number(),
            }),
            memory: zod_1.z.object({
                total: zod_1.z.number(),
                free: zod_1.z.number(),
                used: zod_1.z.number(),
                active: zod_1.z.number(),
                available: zod_1.z.number(),
            }),
            uptime: zod_1.z.number(),
            timestamp: zod_1.z.number(),
            users: zod_1.z.object({
                statusCounts: zod_1.z.record(zod_1.z.enum(Object.values(constants_2.USERS_STATUS)), zod_1.z.number()),
                totalUsers: zod_1.z.number(),
            }),
            onlineStats: zod_1.z.object({
                lastDay: zod_1.z.number(),
                lastWeek: zod_1.z.number(),
                neverOnline: zod_1.z.number(),
                onlineNow: zod_1.z.number(),
            }),
            nodes: zod_1.z.object({
                totalOnline: zod_1.z.number(),
                totalBytesLifetime: zod_1.z.string(),
            }),
        }),
    });
})(GetStatsCommand || (exports.GetStatsCommand = GetStatsCommand = {}));
