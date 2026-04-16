"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.CreateExternalSquadCommand = void 0;
const zod_1 = require("zod");
const api_1 = require("../../api");
const constants_1 = require("../../constants");
const models_1 = require("../../models");
var CreateExternalSquadCommand;
(function (CreateExternalSquadCommand) {
    CreateExternalSquadCommand.url = api_1.REST_API.EXTERNAL_SQUADS.CREATE;
    CreateExternalSquadCommand.TSQ_url = CreateExternalSquadCommand.url;
    CreateExternalSquadCommand.endpointDetails = (0, constants_1.getEndpointDetails)(api_1.EXTERNAL_SQUADS_ROUTES.CREATE, 'post', 'Create external squad');
    CreateExternalSquadCommand.RequestSchema = zod_1.z.object({
        name: zod_1.z
            .string()
            .min(2, 'Name must be at least 2 characters')
            .max(30, 'Name must be less than 30 characters')
            .regex(/^[A-Za-z0-9_\s-]+$/, 'Name can only contain letters, numbers, underscores, dashes and spaces'),
    });
    CreateExternalSquadCommand.ResponseSchema = zod_1.z.object({
        response: models_1.ExternalSquadSchema,
    });
})(CreateExternalSquadCommand || (exports.CreateExternalSquadCommand = CreateExternalSquadCommand = {}));
