"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.DeleteUserHwidDeviceCommand = void 0;
const zod_1 = require("zod");
const api_1 = require("../../api");
const constants_1 = require("../../constants");
const models_1 = require("../../models");
var DeleteUserHwidDeviceCommand;
(function (DeleteUserHwidDeviceCommand) {
    DeleteUserHwidDeviceCommand.url = api_1.REST_API.HWID.DELETE_USER_HWID_DEVICE;
    DeleteUserHwidDeviceCommand.TSQ_url = DeleteUserHwidDeviceCommand.url;
    DeleteUserHwidDeviceCommand.endpointDetails = (0, constants_1.getEndpointDetails)(api_1.HWID_ROUTES.DELETE_USER_HWID_DEVICE, 'post', 'Delete a user HWID device', { scope: 'delete', kind: 'write' });
    DeleteUserHwidDeviceCommand.RequestBodySchema = zod_1.z.object({
        userId: zod_1.z.number(),
        hwid: zod_1.z.string(),
    });
    DeleteUserHwidDeviceCommand.ResponseSchema = zod_1.z.object({
        response: zod_1.z.object({
            total: zod_1.z.number(),
            devices: zod_1.z.array(models_1.HwidUserDeviceSchema),
        }),
    });

})(DeleteUserHwidDeviceCommand || (exports.DeleteUserHwidDeviceCommand = DeleteUserHwidDeviceCommand = {}));
