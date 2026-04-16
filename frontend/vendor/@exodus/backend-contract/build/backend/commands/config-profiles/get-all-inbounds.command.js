"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.GetAllInboundsCommand = void 0;
const zod_1 = require("zod");
const api_1 = require("../../api");
const models_1 = require("../../models");
const constants_1 = require("../../constants");
var GetAllInboundsCommand;
(function (GetAllInboundsCommand) {
    GetAllInboundsCommand.url = api_1.REST_API.CONFIG_PROFILES.GET_ALL_INBOUNDS;
    GetAllInboundsCommand.TSQ_url = GetAllInboundsCommand.url;
    GetAllInboundsCommand.endpointDetails = (0, constants_1.getEndpointDetails)(api_1.CONFIG_PROFILES_ROUTES.GET_ALL_INBOUNDS, 'get', 'Get all inbounds from all config profiles');
    GetAllInboundsCommand.ResponseSchema = zod_1.z.object({
        response: zod_1.z.object({
            total: zod_1.z.number(),
            inbounds: zod_1.z.array(models_1.ConfigProfileInboundsSchema.extend({
                activeSquads: zod_1.z.array(zod_1.z.string().uuid()),
            })),
        }),
    });
})(GetAllInboundsCommand || (exports.GetAllInboundsCommand = GetAllInboundsCommand = {}));
