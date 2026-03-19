"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.GetSubscriptionTemplatesCommand = void 0;
const zod_1 = require("zod");
const api_1 = require("../../api");
const models_1 = require("../../models");
const constants_1 = require("../../constants");
var GetSubscriptionTemplatesCommand;
(function (GetSubscriptionTemplatesCommand) {
    GetSubscriptionTemplatesCommand.url = api_1.REST_API.SUBSCRIPTION_TEMPLATE.GET_ALL;
    GetSubscriptionTemplatesCommand.TSQ_url = GetSubscriptionTemplatesCommand.url;
    GetSubscriptionTemplatesCommand.endpointDetails = (0, constants_1.getEndpointDetails)(api_1.SUBSCRIPTION_TEMPLATE_ROUTES.GET_ALL, 'get', 'Get all subscription templates (wihout content)');
    GetSubscriptionTemplatesCommand.ResponseSchema = zod_1.z.object({
        response: zod_1.z.object({
            total: zod_1.z.number(),
            templates: zod_1.z.array(models_1.SubscriptionTemplateSchema),
        }),
    });
})(GetSubscriptionTemplatesCommand || (exports.GetSubscriptionTemplatesCommand = GetSubscriptionTemplatesCommand = {}));
