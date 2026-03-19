"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.GetTopUsersByHwidDevicesCommand = void 0;
const zod_1 = require("zod");
const constants_1 = require("../../constants");
const api_1 = require("../../api");
var GetTopUsersByHwidDevicesCommand;
(function (GetTopUsersByHwidDevicesCommand) {
    GetTopUsersByHwidDevicesCommand.url = api_1.REST_API.HWID.TOP_USERS_BY_DEVICES;
    GetTopUsersByHwidDevicesCommand.TSQ_url = GetTopUsersByHwidDevicesCommand.url;
    GetTopUsersByHwidDevicesCommand.endpointDetails = (0, constants_1.getEndpointDetails)(api_1.HWID_ROUTES.TOP_USERS_BY_DEVICES, 'get', 'Get top users by HWID devices');
    GetTopUsersByHwidDevicesCommand.RequestQuerySchema = zod_1.z.object({
        start: zod_1.z.coerce
            .number()
            .default(0)
            .describe('Start index (offset) of the results to return, default is 0'),
        size: zod_1.z.coerce
            .number()
            .min(1, 'Size (limit) must be greater than 0')
            .max(100, 'Size (limit) must be less than 100')
            .describe('Number of results to return, no more than 100')
            .default(5),
    });
    GetTopUsersByHwidDevicesCommand.ResponseSchema = zod_1.z.object({
        response: zod_1.z.object({
            users: zod_1.z.array(zod_1.z.object({
                userUuid: zod_1.z.string().uuid(),
                id: zod_1.z.number(),
                username: zod_1.z.string(),
                devicesCount: zod_1.z.number(),
            })),
            total: zod_1.z.number(),
        }),
    });
})(GetTopUsersByHwidDevicesCommand || (exports.GetTopUsersByHwidDevicesCommand = GetTopUsersByHwidDevicesCommand = {}));
