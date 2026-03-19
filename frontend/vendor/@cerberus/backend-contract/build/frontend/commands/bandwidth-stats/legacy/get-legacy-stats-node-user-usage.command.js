"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.GetLegacyStatsNodeUserUsageCommand = void 0;
const zod_1 = require("zod");
const api_1 = require("../../../api");
const constants_1 = require("../../../constants");
var GetLegacyStatsNodeUserUsageCommand;
(function (GetLegacyStatsNodeUserUsageCommand) {
    GetLegacyStatsNodeUserUsageCommand.url = api_1.REST_API.BANDWIDTH_STATS.LEGACY.NODES.GET_USERS;
    GetLegacyStatsNodeUserUsageCommand.TSQ_url = GetLegacyStatsNodeUserUsageCommand.url(':uuid');
    GetLegacyStatsNodeUserUsageCommand.endpointDetails = (0, constants_1.getEndpointDetails)(api_1.BANDWIDTH_STATS_ROUTES.LEGACY.NODES.GET_USERS(':uuid'), 'get', 'Get Node User Usage by Range and Node UUID (Legacy)');
    GetLegacyStatsNodeUserUsageCommand.RequestSchema = zod_1.z.object({
        uuid: zod_1.z.string().uuid(),
    });
    GetLegacyStatsNodeUserUsageCommand.RequestQuerySchema = zod_1.z.object({
        start: zod_1.z.string(),
        end: zod_1.z.string(),
    });
    GetLegacyStatsNodeUserUsageCommand.ResponseSchema = zod_1.z.object({
        response: zod_1.z.array(zod_1.z.object({
            userUuid: zod_1.z.string().uuid(),
            username: zod_1.z.string(),
            nodeUuid: zod_1.z.string().uuid(),
            total: zod_1.z.number(),
            date: zod_1.z.string().transform((str) => new Date(str)),
        })),
    });
})(GetLegacyStatsNodeUserUsageCommand || (exports.GetLegacyStatsNodeUserUsageCommand = GetLegacyStatsNodeUserUsageCommand = {}));
