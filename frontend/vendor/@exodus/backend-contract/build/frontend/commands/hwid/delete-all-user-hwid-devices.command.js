"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.DeleteAllUserHwidDevicesCommand = void 0;
const zod_1 = require("zod");
const api_1 = require("../../api");
const constants_1 = require("../../constants");
const models_1 = require("../../models");
var DeleteAllUserHwidDevicesCommand;
(function (DeleteAllUserHwidDevicesCommand) {
    DeleteAllUserHwidDevicesCommand.url = api_1.REST_API.HWID.DELETE_ALL_USER_HWID_DEVICES;
    DeleteAllUserHwidDevicesCommand.TSQ_url = DeleteAllUserHwidDevicesCommand.url;
    DeleteAllUserHwidDevicesCommand.endpointDetails = (0, constants_1.getEndpointDetails)(api_1.HWID_ROUTES.DELETE_ALL_USER_HWID_DEVICES, 'post', 'Delete all user HWID devices', { scope: 'delete-all', kind: 'write' });
    DeleteAllUserHwidDevicesCommand.RequestBodySchema = zod_1.z.object({
        userId: zod_1.z.number(),
    });
    DeleteAllUserHwidDevicesCommand.ResponseSchema = zod_1.z.object({
        response: zod_1.z.object({
            total: zod_1.z.number(),
            devices: zod_1.z.array(models_1.HwidUserDeviceSchema),
        }),
    });
})(DeleteAllUserHwidDevicesCommand || (exports.DeleteAllUserHwidDevicesCommand = DeleteAllUserHwidDevicesCommand = {}));
