"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.UpdateSubscriptionPageConfigCommand = void 0;
const zod_1 = require("zod");
const api_1 = require("../../api");
const models_1 = require("../../models");
const constants_1 = require("../../constants");
var UpdateSubscriptionPageConfigCommand;
(function (UpdateSubscriptionPageConfigCommand) {
    UpdateSubscriptionPageConfigCommand.url = api_1.REST_API.SUBSCRIPTION_PAGE_CONFIGS.UPDATE;
    UpdateSubscriptionPageConfigCommand.TSQ_url = UpdateSubscriptionPageConfigCommand.url;
    UpdateSubscriptionPageConfigCommand.endpointDetails = (0, constants_1.getEndpointDetails)(api_1.SUBSCRIPTION_PAGE_CONFIGS_ROUTES.UPDATE, 'patch', 'Update subscription page config');
    UpdateSubscriptionPageConfigCommand.RequestSchema = zod_1.z.object({
        uuid: zod_1.z.string().uuid(),
        name: zod_1.z
            .string()
            .min(2, 'Name must be at least 2 characters')
            .max(30, 'Name must be less than 30 characters')
            .regex(/^[A-Za-z0-9_\s-]+$/, 'Name can only contain letters, numbers, underscores, dashes and spaces')
            .optional(),
        config: zod_1.z.optional(zod_1.z.unknown()),
    });
    UpdateSubscriptionPageConfigCommand.ResponseSchema = zod_1.z.object({
        response: models_1.SubscriptionPageConfigSchema.extend({
            config: zod_1.z.unknown(),
        }),
    });
})(UpdateSubscriptionPageConfigCommand || (exports.UpdateSubscriptionPageConfigCommand = UpdateSubscriptionPageConfigCommand = {}));
