"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.CreateInternalSquadCommand = void 0;
const zod_1 = require("zod");
const api_1 = require("../../api");
const constants_1 = require("../../constants");
const models_1 = require("../../models");
var CreateInternalSquadCommand;
(function (CreateInternalSquadCommand) {
    CreateInternalSquadCommand.url = api_1.REST_API.INTERNAL_SQUADS.CREATE;
    CreateInternalSquadCommand.TSQ_url = CreateInternalSquadCommand.url;
    CreateInternalSquadCommand.endpointDetails = (0, constants_1.getEndpointDetails)(api_1.INTERNAL_SQUADS_ROUTES.CREATE, 'post', 'Create internal squad');
    CreateInternalSquadCommand.RequestSchema = zod_1.z.object({
        name: zod_1.z
            .string()
            .min(2, 'Name must be at least 2 characters')
            .max(30, 'Name must be less than 30 characters')
            .regex(/^[A-Za-z0-9_\s-]+$/, 'Name can only contain letters, numbers, underscores, dashes and spaces'),
        inbounds: zod_1.z.array(zod_1.z.string().uuid()),
    });
    CreateInternalSquadCommand.ResponseSchema = zod_1.z.object({
        response: models_1.InternalSquadSchema,
    });
})(CreateInternalSquadCommand || (exports.CreateInternalSquadCommand = CreateInternalSquadCommand = {}));
