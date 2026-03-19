"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.CreateSubscriptionPageConfigCommand = void 0;
const zod_1 = require("zod");
const api_1 = require("../../api");
const models_1 = require("../../models");
const constants_1 = require("../../constants");
var CreateSubscriptionPageConfigCommand;
(function (CreateSubscriptionPageConfigCommand) {
    CreateSubscriptionPageConfigCommand.url = api_1.REST_API.SUBSCRIPTION_PAGE_CONFIGS.CREATE;
    CreateSubscriptionPageConfigCommand.TSQ_url = CreateSubscriptionPageConfigCommand.url;
    CreateSubscriptionPageConfigCommand.endpointDetails = (0, constants_1.getEndpointDetails)(api_1.SUBSCRIPTION_PAGE_CONFIGS_ROUTES.CREATE, 'post', 'Create subscription page config');
    CreateSubscriptionPageConfigCommand.RequestSchema = zod_1.z.object({
        name: zod_1.z
            .string()
            .min(2, 'Name must be at least 2 characters')
            .max(30, 'Name must be less than 30 characters')
            .regex(/^[A-Za-z0-9_\s-]+$/, 'Name can only contain letters, numbers, underscores, dashes and spaces'),
    });
    CreateSubscriptionPageConfigCommand.ResponseSchema = zod_1.z.object({
        response: models_1.SubscriptionPageConfigSchema,
    });
})(CreateSubscriptionPageConfigCommand || (exports.CreateSubscriptionPageConfigCommand = CreateSubscriptionPageConfigCommand = {}));
