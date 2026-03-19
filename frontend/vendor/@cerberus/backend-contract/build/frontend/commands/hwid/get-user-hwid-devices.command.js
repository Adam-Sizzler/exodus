"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.GetUserHwidDevicesCommand = void 0;
const zod_1 = require("zod");
const constants_1 = require("../../constants");
const models_1 = require("../../models");
const api_1 = require("../../api");
var GetUserHwidDevicesCommand;
(function (GetUserHwidDevicesCommand) {
    GetUserHwidDevicesCommand.url = api_1.REST_API.HWID.GET_USER_HWID_DEVICES;
    GetUserHwidDevicesCommand.TSQ_url = GetUserHwidDevicesCommand.url(':userUuid');
    GetUserHwidDevicesCommand.endpointDetails = (0, constants_1.getEndpointDetails)(api_1.HWID_ROUTES.GET_USER_HWID_DEVICES(':userUuid'), 'get', 'Get user HWID devices');
    GetUserHwidDevicesCommand.RequestSchema = zod_1.z.object({
        userUuid: zod_1.z.string().uuid(),
    });
    GetUserHwidDevicesCommand.ResponseSchema = zod_1.z.object({
        response: zod_1.z.object({
            total: zod_1.z.number(),
            devices: zod_1.z.array(models_1.HwidUserDeviceSchema),
        }),
    });
})(GetUserHwidDevicesCommand || (exports.GetUserHwidDevicesCommand = GetUserHwidDevicesCommand = {}));
