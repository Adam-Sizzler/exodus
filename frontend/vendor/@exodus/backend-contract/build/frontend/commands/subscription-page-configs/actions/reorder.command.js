"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.ReorderSubscriptionPageConfigsCommand = void 0;
const zod_1 = require("zod");
const api_1 = require("../../../api");
const models_1 = require("../../../models");
const constants_1 = require("../../../constants");
var ReorderSubscriptionPageConfigsCommand;
(function (ReorderSubscriptionPageConfigsCommand) {
    ReorderSubscriptionPageConfigsCommand.url = api_1.REST_API.SUBSCRIPTION_PAGE_CONFIGS.ACTIONS.REORDER;
    ReorderSubscriptionPageConfigsCommand.TSQ_url = ReorderSubscriptionPageConfigsCommand.url;
    ReorderSubscriptionPageConfigsCommand.endpointDetails = (0, constants_1.getEndpointDetails)(api_1.SUBSCRIPTION_PAGE_CONFIGS_ROUTES.ACTIONS.REORDER, 'post', 'Reorder subscription page configs');
    ReorderSubscriptionPageConfigsCommand.RequestSchema = zod_1.z.object({
        items: zod_1.z.array(models_1.SubscriptionPageConfigSchema.pick({
            viewPosition: true,
            uuid: true,
        })),
    });
    ReorderSubscriptionPageConfigsCommand.ResponseSchema = zod_1.z.object({
        response: zod_1.z.object({
            total: zod_1.z.number(),
            configs: zod_1.z.array(models_1.SubscriptionPageConfigSchema),
        }),
    });
})(ReorderSubscriptionPageConfigsCommand || (exports.ReorderSubscriptionPageConfigsCommand = ReorderSubscriptionPageConfigsCommand = {}));
