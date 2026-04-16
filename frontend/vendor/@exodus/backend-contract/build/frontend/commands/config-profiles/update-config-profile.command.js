"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.UpdateConfigProfileCommand = void 0;
const zod_1 = require("zod");
const api_1 = require("../../api");
const constants_1 = require("../../constants");
const models_1 = require("../../models");
var UpdateConfigProfileCommand;
(function (UpdateConfigProfileCommand) {
    UpdateConfigProfileCommand.url = api_1.REST_API.CONFIG_PROFILES.UPDATE;
    UpdateConfigProfileCommand.TSQ_url = UpdateConfigProfileCommand.url;
    UpdateConfigProfileCommand.endpointDetails = (0, constants_1.getEndpointDetails)(api_1.CONFIG_PROFILES_ROUTES.UPDATE, 'patch', 'Update Core Config in specific config profile');
    UpdateConfigProfileCommand.RequestSchema = zod_1.z.object({
        uuid: zod_1.z.string().uuid('UUID must be a valid UUID'),
        name: zod_1.z
            .string()
            .min(2, 'Name must be at least 2 characters')
            .max(30, 'Name must be less than 30 characters')
            .regex(/^[A-Za-z0-9_\s-]+$/, 'Name can only contain letters, numbers, underscores, dashes and spaces')
            .optional(),
        config: zod_1.z.object({}).passthrough().optional(),
    });
    UpdateConfigProfileCommand.ResponseSchema = zod_1.z.object({
        response: models_1.ConfigProfileSchema,
    });
})(UpdateConfigProfileCommand || (exports.UpdateConfigProfileCommand = UpdateConfigProfileCommand = {}));
