"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.CloneSubscriptionPageConfigCommand = void 0;
const zod_1 = require("zod");
const api_1 = require("../../../api");
const models_1 = require("../../../models");
const constants_1 = require("../../../constants");
var CloneSubscriptionPageConfigCommand;
(function (CloneSubscriptionPageConfigCommand) {
    CloneSubscriptionPageConfigCommand.url = api_1.REST_API.SUBSCRIPTION_PAGE_CONFIGS.ACTIONS.CLONE;
    CloneSubscriptionPageConfigCommand.TSQ_url = CloneSubscriptionPageConfigCommand.url;
    CloneSubscriptionPageConfigCommand.endpointDetails = (0, constants_1.getEndpointDetails)(api_1.SUBSCRIPTION_PAGE_CONFIGS_ROUTES.ACTIONS.CLONE, 'post', 'Clone subscription page config');
    CloneSubscriptionPageConfigCommand.RequestSchema = zod_1.z.object({
        cloneFromUuid: zod_1.z.string().uuid(),
    });
    CloneSubscriptionPageConfigCommand.ResponseSchema = zod_1.z.object({
        response: models_1.SubscriptionPageConfigSchema.extend({
            config: zod_1.z.unknown(),
        }),
    });
})(CloneSubscriptionPageConfigCommand || (exports.CloneSubscriptionPageConfigCommand = CloneSubscriptionPageConfigCommand = {}));
