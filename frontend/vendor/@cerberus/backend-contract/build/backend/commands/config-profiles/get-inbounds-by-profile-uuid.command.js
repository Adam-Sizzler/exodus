"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.GetInboundsByProfileUuidCommand = void 0;
const zod_1 = require("zod");
const api_1 = require("../../api");
const models_1 = require("../../models");
const constants_1 = require("../../constants");
var GetInboundsByProfileUuidCommand;
(function (GetInboundsByProfileUuidCommand) {
    GetInboundsByProfileUuidCommand.url = api_1.REST_API.CONFIG_PROFILES.GET_INBOUNDS_BY_PROFILE_UUID;
    GetInboundsByProfileUuidCommand.TSQ_url = GetInboundsByProfileUuidCommand.url(':uuid');
    GetInboundsByProfileUuidCommand.endpointDetails = (0, constants_1.getEndpointDetails)(api_1.CONFIG_PROFILES_ROUTES.GET_INBOUNDS_BY_PROFILE_UUID(':uuid'), 'get', 'Get inbounds by profile uuid');
    GetInboundsByProfileUuidCommand.RequestSchema = zod_1.z.object({
        uuid: zod_1.z.string().uuid(),
    });
    GetInboundsByProfileUuidCommand.ResponseSchema = zod_1.z.object({
        response: zod_1.z.object({
            total: zod_1.z.number(),
            inbounds: zod_1.z.array(models_1.ConfigProfileInboundsSchema.extend({
                activeSquads: zod_1.z.array(zod_1.z.string().uuid()),
            })),
        }),
    });
})(GetInboundsByProfileUuidCommand || (exports.GetInboundsByProfileUuidCommand = GetInboundsByProfileUuidCommand = {}));
