"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.CreateConfigProfileCommand = void 0;
const zod_1 = require("zod");
const api_1 = require("../../api");
const constants_1 = require("../../constants");
const models_1 = require("../../models");
var CreateConfigProfileCommand;
(function (CreateConfigProfileCommand) {
    CreateConfigProfileCommand.url = api_1.REST_API.CONFIG_PROFILES.CREATE;
    CreateConfigProfileCommand.TSQ_url = CreateConfigProfileCommand.url;
    CreateConfigProfileCommand.endpointDetails = (0, constants_1.getEndpointDetails)(api_1.CONFIG_PROFILES_ROUTES.CREATE, 'post', 'Create config profile');
    CreateConfigProfileCommand.RequestSchema = zod_1.z.object({
        name: zod_1.z
            .string()
            .min(2, 'Name must be at least 2 characters')
            .max(30, 'Name must be less than 30 characters')
            .regex(/^[A-Za-z0-9_\s-]+$/, 'Name can only contain letters, numbers, underscores, dashes and spaces'),
        config: zod_1.z.object({}).passthrough(),
    });
    CreateConfigProfileCommand.ResponseSchema = zod_1.z.object({
        response: models_1.ConfigProfileSchema,
    });
})(CreateConfigProfileCommand || (exports.CreateConfigProfileCommand = CreateConfigProfileCommand = {}));
