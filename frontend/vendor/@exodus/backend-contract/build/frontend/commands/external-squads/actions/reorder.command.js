"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.ReorderExternalSquadCommand = void 0;
const zod_1 = require("zod");
const api_1 = require("../../../api");
const constants_1 = require("../../../constants");
const models_1 = require("../../../models");
var ReorderExternalSquadCommand;
(function (ReorderExternalSquadCommand) {
    ReorderExternalSquadCommand.url = api_1.REST_API.EXTERNAL_SQUADS.ACTIONS.REORDER;
    ReorderExternalSquadCommand.TSQ_url = ReorderExternalSquadCommand.url;
    ReorderExternalSquadCommand.endpointDetails = (0, constants_1.getEndpointDetails)(api_1.EXTERNAL_SQUADS_ROUTES.ACTIONS.REORDER, 'post', 'Reorder external squads');
    ReorderExternalSquadCommand.RequestSchema = zod_1.z.object({
        items: zod_1.z.array(models_1.ExternalSquadSchema.pick({
            viewPosition: true,
            uuid: true,
        })),
    });
    ReorderExternalSquadCommand.ResponseSchema = zod_1.z.object({
        response: zod_1.z.object({
            total: zod_1.z.number(),
            externalSquads: zod_1.z.array(models_1.ExternalSquadSchema),
        }),
    });
})(ReorderExternalSquadCommand || (exports.ReorderExternalSquadCommand = ReorderExternalSquadCommand = {}));
