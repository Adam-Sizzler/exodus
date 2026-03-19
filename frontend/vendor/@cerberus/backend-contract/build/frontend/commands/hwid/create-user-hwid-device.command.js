"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.CreateUserHwidDeviceCommand = void 0;
const zod_1 = require("zod");
const constants_1 = require("../../constants");
const models_1 = require("../../models");
const api_1 = require("../../api");
var CreateUserHwidDeviceCommand;
(function (CreateUserHwidDeviceCommand) {
    CreateUserHwidDeviceCommand.url = api_1.REST_API.HWID.CREATE_USER_HWID_DEVICE;
    CreateUserHwidDeviceCommand.TSQ_url = CreateUserHwidDeviceCommand.url;
    CreateUserHwidDeviceCommand.endpointDetails = (0, constants_1.getEndpointDetails)(api_1.HWID_ROUTES.CREATE_USER_HWID_DEVICE, 'post', 'Create a user HWID device');
    CreateUserHwidDeviceCommand.RequestSchema = zod_1.z.object({
        hwid: zod_1.z.string(),
        userUuid: zod_1.z.string().uuid(),
        platform: zod_1.z.optional(zod_1.z.string()),
        osVersion: zod_1.z.optional(zod_1.z.string()),
        deviceModel: zod_1.z.optional(zod_1.z.string()),
        userAgent: zod_1.z.optional(zod_1.z.string()),
    });
    CreateUserHwidDeviceCommand.ResponseSchema = zod_1.z.object({
        response: zod_1.z.object({
            total: zod_1.z.number(),
            devices: zod_1.z.array(models_1.HwidUserDeviceSchema),
        }),
    });
})(CreateUserHwidDeviceCommand || (exports.CreateUserHwidDeviceCommand = CreateUserHwidDeviceCommand = {}));
