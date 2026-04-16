"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.GetSubscriptionPageConfigCommand = void 0;
const zod_1 = require("zod");
const api_1 = require("../../api");
const models_1 = require("../../models");
const constants_1 = require("../../constants");
var GetSubscriptionPageConfigCommand;
(function (GetSubscriptionPageConfigCommand) {
    GetSubscriptionPageConfigCommand.url = api_1.REST_API.SUBSCRIPTION_PAGE_CONFIGS.GET;
    GetSubscriptionPageConfigCommand.TSQ_url = GetSubscriptionPageConfigCommand.url(':uuid');
    GetSubscriptionPageConfigCommand.endpointDetails = (0, constants_1.getEndpointDetails)(api_1.SUBSCRIPTION_PAGE_CONFIGS_ROUTES.GET(':uuid'), 'get', 'Get subscription page config by uuid');
    GetSubscriptionPageConfigCommand.RequestSchema = zod_1.z.object({
        uuid: zod_1.z.string().uuid(),
    });
    GetSubscriptionPageConfigCommand.ResponseSchema = zod_1.z.object({
        response: models_1.SubscriptionPageConfigSchema.extend({
            config: zod_1.z.unknown(),
        }),
    });
})(GetSubscriptionPageConfigCommand || (exports.GetSubscriptionPageConfigCommand = GetSubscriptionPageConfigCommand = {}));
