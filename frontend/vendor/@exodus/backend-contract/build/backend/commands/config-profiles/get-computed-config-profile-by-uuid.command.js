"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.GetComputedConfigProfileByUuidCommand = void 0;
const zod_1 = require("zod");
const api_1 = require("../../api");
const constants_1 = require("../../constants");
const models_1 = require("../../models");
var GetComputedConfigProfileByUuidCommand;
(function (GetComputedConfigProfileByUuidCommand) {
    GetComputedConfigProfileByUuidCommand.url = api_1.REST_API.CONFIG_PROFILES.GET_COMPUTED_CONFIG_BY_PROFILE_UUID;
    GetComputedConfigProfileByUuidCommand.TSQ_url = GetComputedConfigProfileByUuidCommand.url(':uuid');
    GetComputedConfigProfileByUuidCommand.endpointDetails = (0, constants_1.getEndpointDetails)(api_1.CONFIG_PROFILES_ROUTES.GET_COMPUTED_CONFIG_BY_PROFILE_UUID(':uuid'), 'get', 'Get computed config profile by uuid');
    GetComputedConfigProfileByUuidCommand.RequestSchema = zod_1.z.object({
        uuid: zod_1.z.string().uuid(),
    });
    GetComputedConfigProfileByUuidCommand.ResponseSchema = zod_1.z.object({
        response: models_1.ConfigProfileSchema,
    });
})(GetComputedConfigProfileByUuidCommand || (exports.GetComputedConfigProfileByUuidCommand = GetComputedConfigProfileByUuidCommand = {}));
