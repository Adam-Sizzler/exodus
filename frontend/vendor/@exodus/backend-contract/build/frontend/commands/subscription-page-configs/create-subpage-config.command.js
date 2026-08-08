"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.CreateSubpageConfigCommand = void 0;
const zod_1 = require("zod");
const api_1 = require("../../api");
const constants_1 = require("../../constants");
const models_1 = require("../../models");
var CreateSubpageConfigCommand;
(function (CreateSubpageConfigCommand) {
    CreateSubpageConfigCommand.url = api_1.REST_API.SUBSCRIPTION_PAGE_CONFIGS.CREATE;
    CreateSubpageConfigCommand.TSQ_url = CreateSubpageConfigCommand.url;
    CreateSubpageConfigCommand.endpointDetails = (0, constants_1.getEndpointDetails)(api_1.SUBSCRIPTION_PAGE_CONFIGS_ROUTES.CREATE, 'post', 'Create subscription page config', { scope: 'create', kind: 'write' });
    CreateSubpageConfigCommand.RequestBodySchema = zod_1.z.object({
        name: zod_1.z
            .string()
            .min(2, 'Name must be at least 2 characters')
            .max(30, 'Name must be less than 30 characters')
            .regex(/^[A-Za-z0-9_\s-]+$/, 'Name can only contain letters, numbers, underscores, dashes and spaces'),
    });
    CreateSubpageConfigCommand.ResponseSchema = zod_1.z.object({
        response: models_1.SubscriptionPageConfigSchema,
    });

})(CreateSubpageConfigCommand || (exports.CreateSubpageConfigCommand = CreateSubpageConfigCommand = {}));
