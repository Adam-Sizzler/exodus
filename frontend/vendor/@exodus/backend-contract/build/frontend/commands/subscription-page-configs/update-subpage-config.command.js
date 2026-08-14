"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.UpdateSubpageConfigCommand = void 0;
const zod_1 = require("zod");
const api_1 = require("../../api");
const constants_1 = require("../../constants");
const models_1 = require("../../models");
var UpdateSubpageConfigCommand;
(function (UpdateSubpageConfigCommand) {
    UpdateSubpageConfigCommand.url = api_1.REST_API.SUBSCRIPTION_PAGE_CONFIGS.UPDATE;
    UpdateSubpageConfigCommand.TSQ_url = UpdateSubpageConfigCommand.url;
    UpdateSubpageConfigCommand.endpointDetails = (0, constants_1.getEndpointDetails)(api_1.SUBSCRIPTION_PAGE_CONFIGS_ROUTES.UPDATE, 'patch', 'Update subscription page config', { scope: 'update', kind: 'write' });
    UpdateSubpageConfigCommand.RequestBodySchema = zod_1.z.object({
        uuid: zod_1.z.uuid(),
        name: zod_1.z
            .string()
            .min(2, 'Name must be at least 2 characters')
            .max(30, 'Name must be less than 30 characters')
            .regex(/^[A-Za-z0-9_\s-]+$/, 'Name can only contain letters, numbers, underscores, dashes and spaces')
            .optional(),
        config: zod_1.z.optional(zod_1.z.unknown()),
    });
    UpdateSubpageConfigCommand.ResponseSchema = zod_1.z.object({
        response: models_1.SubscriptionPageConfigSchema.extend({
            config: zod_1.z.unknown(),
        }),
    });
})(UpdateSubpageConfigCommand || (exports.UpdateSubpageConfigCommand = UpdateSubpageConfigCommand = {}));
