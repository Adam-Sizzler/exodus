"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.GetAllHwidDevicesCommand = void 0;
const zod_1 = require("zod");
const models_1 = require("../../models");
const constants_1 = require("../../constants");
const api_1 = require("../../api");
var GetAllHwidDevicesCommand;
(function (GetAllHwidDevicesCommand) {
    GetAllHwidDevicesCommand.url = api_1.REST_API.HWID.GET_ALL_HWID_DEVICES;
    GetAllHwidDevicesCommand.TSQ_url = GetAllHwidDevicesCommand.url;
    GetAllHwidDevicesCommand.endpointDetails = (0, constants_1.getEndpointDetails)(api_1.HWID_ROUTES.GET_ALL_HWID_DEVICES, 'get', 'Get all HWID devices');
    GetAllHwidDevicesCommand.RequestQuerySchema = models_1.TanstackQueryRequestQuerySchema;
    GetAllHwidDevicesCommand.ResponseSchema = zod_1.z.object({
        response: zod_1.z.object({
            devices: zod_1.z.array(models_1.HwidUserDeviceSchema),
            total: zod_1.z.number(),
        }),
    });
})(GetAllHwidDevicesCommand || (exports.GetAllHwidDevicesCommand = GetAllHwidDevicesCommand = {}));
