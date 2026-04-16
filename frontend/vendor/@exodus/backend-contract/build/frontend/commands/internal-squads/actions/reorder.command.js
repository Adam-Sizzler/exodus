"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.ReorderInternalSquadCommand = void 0;
const zod_1 = require("zod");
const api_1 = require("../../../api");
const constants_1 = require("../../../constants");
const models_1 = require("../../../models");
var ReorderInternalSquadCommand;
(function (ReorderInternalSquadCommand) {
    ReorderInternalSquadCommand.url = api_1.REST_API.INTERNAL_SQUADS.ACTIONS.REORDER;
    ReorderInternalSquadCommand.TSQ_url = ReorderInternalSquadCommand.url;
    ReorderInternalSquadCommand.endpointDetails = (0, constants_1.getEndpointDetails)(api_1.INTERNAL_SQUADS_ROUTES.ACTIONS.REORDER, 'post', 'Reorder internal squads');
    ReorderInternalSquadCommand.RequestSchema = zod_1.z.object({
        items: zod_1.z.array(models_1.InternalSquadSchema.pick({
            viewPosition: true,
            uuid: true,
        })),
    });
    ReorderInternalSquadCommand.ResponseSchema = zod_1.z.object({
        response: zod_1.z.object({
            total: zod_1.z.number(),
            internalSquads: zod_1.z.array(models_1.InternalSquadSchema),
        }),
    });
})(ReorderInternalSquadCommand || (exports.ReorderInternalSquadCommand = ReorderInternalSquadCommand = {}));
