"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.GetConfigProfilesCommand = void 0;
const zod_1 = require("zod");
const api_1 = require("../../api");
const constants_1 = require("../../constants");
const models_1 = require("../../models");
var GetConfigProfilesCommand;
(function (GetConfigProfilesCommand) {
    GetConfigProfilesCommand.url = api_1.REST_API.CONFIG_PROFILES.GET;
    GetConfigProfilesCommand.TSQ_url = GetConfigProfilesCommand.url;
    GetConfigProfilesCommand.endpointDetails = (0, constants_1.getEndpointDetails)(api_1.CONFIG_PROFILES_ROUTES.GET, 'get', 'Get config profiles');
    GetConfigProfilesCommand.ResponseSchema = zod_1.z.object({
        response: zod_1.z.object({
            total: zod_1.z.number(),
            configProfiles: zod_1.z.array(models_1.ConfigProfileSchema),
        }),
    });
})(GetConfigProfilesCommand || (exports.GetConfigProfilesCommand = GetConfigProfilesCommand = {}));
