"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.GetConfigProfileByUuidCommand = void 0;
const zod_1 = require("zod");
const api_1 = require("../../api");
const constants_1 = require("../../constants");
const models_1 = require("../../models");
var GetConfigProfileByUuidCommand;
(function (GetConfigProfileByUuidCommand) {
    GetConfigProfileByUuidCommand.url = api_1.REST_API.CONFIG_PROFILES.GET_BY_UUID;
    GetConfigProfileByUuidCommand.TSQ_url = GetConfigProfileByUuidCommand.url(':uuid');
    GetConfigProfileByUuidCommand.endpointDetails = (0, constants_1.getEndpointDetails)(api_1.CONFIG_PROFILES_ROUTES.GET_BY_UUID(':uuid'), 'get', 'Get config profile by uuid');
    GetConfigProfileByUuidCommand.RequestSchema = zod_1.z.object({
        uuid: zod_1.z.string().uuid(),
    });
    GetConfigProfileByUuidCommand.ResponseSchema = zod_1.z.object({
        response: models_1.ConfigProfileSchema,
    });
})(GetConfigProfileByUuidCommand || (exports.GetConfigProfileByUuidCommand = GetConfigProfileByUuidCommand = {}));
