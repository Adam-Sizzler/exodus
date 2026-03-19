"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.GetLegacyStatsUserUsageCommand = void 0;
const zod_1 = require("zod");
const api_1 = require("../../../api");
const constants_1 = require("../../../constants");
var GetLegacyStatsUserUsageCommand;
(function (GetLegacyStatsUserUsageCommand) {
    GetLegacyStatsUserUsageCommand.url = api_1.REST_API.BANDWIDTH_STATS.LEGACY.USERS.GET_BY_UUID;
    GetLegacyStatsUserUsageCommand.TSQ_url = GetLegacyStatsUserUsageCommand.url(':uuid');
    GetLegacyStatsUserUsageCommand.endpointDetails = (0, constants_1.getEndpointDetails)(api_1.BANDWIDTH_STATS_ROUTES.LEGACY.USERS.GET_BY_UUID(':uuid'), 'get', 'Get User Usage by Range (Legacy)');
    GetLegacyStatsUserUsageCommand.RequestSchema = zod_1.z.object({
        uuid: zod_1.z.string().uuid(),
    });
    GetLegacyStatsUserUsageCommand.RequestQuerySchema = zod_1.z.object({
        start: zod_1.z.string(),
        end: zod_1.z.string(),
    });
    GetLegacyStatsUserUsageCommand.ResponseSchema = zod_1.z.object({
        response: zod_1.z.array(zod_1.z.object({
            userUuid: zod_1.z.string().uuid(),
            nodeUuid: zod_1.z.string().uuid(),
            nodeName: zod_1.z.string(),
            countryCode: zod_1.z.string(),
            total: zod_1.z.number(),
            date: zod_1.z.string().transform((str) => new Date(str)),
        })),
    });
})(GetLegacyStatsUserUsageCommand || (exports.GetLegacyStatsUserUsageCommand = GetLegacyStatsUserUsageCommand = {}));
