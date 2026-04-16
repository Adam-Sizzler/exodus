"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.GetSubscriptionPageConfigsCommand = void 0;
const zod_1 = require("zod");
const api_1 = require("../../api");
const models_1 = require("../../models");
const constants_1 = require("../../constants");
var GetSubscriptionPageConfigsCommand;
(function (GetSubscriptionPageConfigsCommand) {
    GetSubscriptionPageConfigsCommand.url = api_1.REST_API.SUBSCRIPTION_PAGE_CONFIGS.GET_ALL;
    GetSubscriptionPageConfigsCommand.TSQ_url = GetSubscriptionPageConfigsCommand.url;
    GetSubscriptionPageConfigsCommand.endpointDetails = (0, constants_1.getEndpointDetails)(api_1.SUBSCRIPTION_PAGE_CONFIGS_ROUTES.GET_ALL, 'get', 'Get all subscription page configs');
    GetSubscriptionPageConfigsCommand.ResponseSchema = zod_1.z.object({
        response: zod_1.z.object({
            total: zod_1.z.number(),
            configs: zod_1.z.array(models_1.SubscriptionPageConfigSchema),
        }),
    });
})(GetSubscriptionPageConfigsCommand || (exports.GetSubscriptionPageConfigsCommand = GetSubscriptionPageConfigsCommand = {}));
