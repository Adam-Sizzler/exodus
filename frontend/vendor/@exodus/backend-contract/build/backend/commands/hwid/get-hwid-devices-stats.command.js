"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.GetHwidDevicesStatsCommand = void 0;
const zod_1 = require("zod");
const constants_1 = require("../../constants");
const api_1 = require("../../api");
var GetHwidDevicesStatsCommand;
(function (GetHwidDevicesStatsCommand) {
    GetHwidDevicesStatsCommand.url = api_1.REST_API.HWID.STATS;
    GetHwidDevicesStatsCommand.TSQ_url = GetHwidDevicesStatsCommand.url;
    GetHwidDevicesStatsCommand.endpointDetails = (0, constants_1.getEndpointDetails)(api_1.HWID_ROUTES.STATS, 'get', 'Get HWID devices stats');
    GetHwidDevicesStatsCommand.ResponseSchema = zod_1.z.object({
        response: zod_1.z.object({
            byPlatform: zod_1.z.array(zod_1.z.object({
                platform: zod_1.z.string(),
                count: zod_1.z.number(),
            })),
            byApp: zod_1.z.array(zod_1.z.object({
                app: zod_1.z.string(),
                count: zod_1.z.number(),
            })),
            stats: zod_1.z.object({
                totalUniqueDevices: zod_1.z.number(),
                totalHwidDevices: zod_1.z.number(),
                averageHwidDevicesPerUser: zod_1.z.number(),
            }),
        }),
    });
})(GetHwidDevicesStatsCommand || (exports.GetHwidDevicesStatsCommand = GetHwidDevicesStatsCommand = {}));
