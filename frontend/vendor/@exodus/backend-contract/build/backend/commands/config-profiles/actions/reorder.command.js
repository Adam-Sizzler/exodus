"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.ReorderConfigProfileCommand = void 0;
const zod_1 = require("zod");
const api_1 = require("../../../api");
const constants_1 = require("../../../constants");
const models_1 = require("../../../models");
var ReorderConfigProfileCommand;
(function (ReorderConfigProfileCommand) {
    ReorderConfigProfileCommand.url = api_1.REST_API.CONFIG_PROFILES.ACTIONS.REORDER;
    ReorderConfigProfileCommand.TSQ_url = ReorderConfigProfileCommand.url;
    ReorderConfigProfileCommand.endpointDetails = (0, constants_1.getEndpointDetails)(api_1.CONFIG_PROFILES_ROUTES.ACTIONS.REORDER, 'post', 'Reorder config profiles');
    ReorderConfigProfileCommand.RequestSchema = zod_1.z.object({
        items: zod_1.z.array(models_1.ConfigProfileSchema.pick({
            viewPosition: true,
            uuid: true,
        })),
    });
    ReorderConfigProfileCommand.ResponseSchema = zod_1.z.object({
        response: zod_1.z.object({
            total: zod_1.z.number(),
            configProfiles: zod_1.z.array(models_1.ConfigProfileSchema),
        }),
    });
})(ReorderConfigProfileCommand || (exports.ReorderConfigProfileCommand = ReorderConfigProfileCommand = {}));
