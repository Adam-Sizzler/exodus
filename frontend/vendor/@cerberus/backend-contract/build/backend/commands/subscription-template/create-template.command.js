"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.CreateSubscriptionTemplateCommand = void 0;
const zod_1 = require("zod");
const constants_1 = require("../../constants");
const api_1 = require("../../api");
const models_1 = require("../../models");
var CreateSubscriptionTemplateCommand;
(function (CreateSubscriptionTemplateCommand) {
    CreateSubscriptionTemplateCommand.url = api_1.REST_API.SUBSCRIPTION_TEMPLATE.CREATE;
    CreateSubscriptionTemplateCommand.TSQ_url = CreateSubscriptionTemplateCommand.url;
    CreateSubscriptionTemplateCommand.endpointDetails = (0, constants_1.getEndpointDetails)(api_1.SUBSCRIPTION_TEMPLATE_ROUTES.CREATE, 'post', 'Create subscription template');
    CreateSubscriptionTemplateCommand.RequestSchema = zod_1.z.object({
        name: zod_1.z
            .string()
            .min(2, 'Name must be at least 2 characters')
            .max(255, 'Name must be less than 255 characters')
            .regex(/^[A-Za-z0-9_\s-]+$/, 'Name can only contain letters, numbers, underscores, dashes and spaces'),
        templateType: zod_1.z.nativeEnum(constants_1.SUBSCRIPTION_TEMPLATE_TYPE),
    });
    CreateSubscriptionTemplateCommand.ResponseSchema = zod_1.z.object({
        response: models_1.SubscriptionTemplateSchema,
    });
})(CreateSubscriptionTemplateCommand || (exports.CreateSubscriptionTemplateCommand = CreateSubscriptionTemplateCommand = {}));
