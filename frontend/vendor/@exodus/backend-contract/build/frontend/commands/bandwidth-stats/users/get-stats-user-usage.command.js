"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.GetStatsUserUsageCommand = void 0;
const zod_1 = require("zod");
const api_1 = require("../../../api");
const constants_1 = require("../../../constants");
var GetStatsUserUsageCommand;
(function (GetStatsUserUsageCommand) {
    GetStatsUserUsageCommand.url = api_1.REST_API.BANDWIDTH_STATS.USERS.GET_BY_UUID;
    GetStatsUserUsageCommand.TSQ_url = GetStatsUserUsageCommand.url(':uuid');
    GetStatsUserUsageCommand.endpointDetails = (0, constants_1.getEndpointDetails)(api_1.BANDWIDTH_STATS_ROUTES.USERS.GET_BY_UUID(':uuid'), 'get', 'Get User Usage by Range');
    GetStatsUserUsageCommand.RequestSchema = zod_1.z.object({
        uuid: zod_1.z.string().uuid(),
    });
    GetStatsUserUsageCommand.RequestQuerySchema = zod_1.z.object({
        start: zod_1.z.string().date(),
        end: zod_1.z.string().date(),
        topNodesLimit: zod_1.z.coerce.number().min(1).default(20),
    });
    GetStatsUserUsageCommand.ResponseSchema = zod_1.z.object({
        response: zod_1.z.object({
            categories: zod_1.z.array(zod_1.z.string()),
            sparklineData: zod_1.z.array(zod_1.z.number()),
            topNodes: zod_1.z.array(zod_1.z.object({
                uuid: zod_1.z.string().uuid(),
                color: zod_1.z.string(),
                name: zod_1.z.string(),
                countryCode: zod_1.z.string(),
                total: zod_1.z.number(),
            })),
            series: zod_1.z.array(zod_1.z.object({
                uuid: zod_1.z.string().uuid(),
                name: zod_1.z.string(),
                color: zod_1.z.string(),
                countryCode: zod_1.z.string(),
                total: zod_1.z.number(),
                data: zod_1.z.array(zod_1.z.number()),
            })),
        }),
    });
})(GetStatsUserUsageCommand || (exports.GetStatsUserUsageCommand = GetStatsUserUsageCommand = {}));
