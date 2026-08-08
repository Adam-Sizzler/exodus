"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.GetHwidDevicesCommand = void 0;
const zod_1 = require("zod");
const api_1 = require("../../api");
const constants_1 = require("../../constants");
const models_1 = require("../../models");
var GetHwidDevicesCommand;
(function (GetHwidDevicesCommand) {
    GetHwidDevicesCommand.url = api_1.REST_API.HWID.GET_ALL_HWID_DEVICES;
    GetHwidDevicesCommand.TSQ_url = GetHwidDevicesCommand.url;
    GetHwidDevicesCommand.endpointDetails = (0, constants_1.getEndpointDetails)(api_1.HWID_ROUTES.GET_ALL_HWID_DEVICES, 'get', 'Get HWID devices', { scope: 'list', kind: 'read' }, 'Please note that the filters here are primarily intended for use by the frontend and rely on expensive operators such as LIKE under the hood. Misusing these filters may negatively impact the performance of your database.');
    GetHwidDevicesCommand.RequestQuerySchema = models_1.TanstackQueryRequestQuerySchema;
    GetHwidDevicesCommand.ResponseSchema = zod_1.z.object({
        response: zod_1.z.object({
            devices: zod_1.z.array(models_1.HwidUserDeviceSchema),
            total: zod_1.z.number(),
        }),
    });

})(GetHwidDevicesCommand || (exports.GetHwidDevicesCommand = GetHwidDevicesCommand = {}));
