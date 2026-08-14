"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.GetSubscriptionTemplateCommand = void 0;
const zod_1 = require("zod");
const api_1 = require("../../api");
const constants_1 = require("../../constants");
const models_1 = require("../../models");
var GetSubscriptionTemplateCommand;
(function (GetSubscriptionTemplateCommand) {
    GetSubscriptionTemplateCommand.url = api_1.REST_API.SUBSCRIPTION_TEMPLATE.GET;
    GetSubscriptionTemplateCommand.TSQ_url = GetSubscriptionTemplateCommand.url(':uuid');
    GetSubscriptionTemplateCommand.endpointDetails = (0, constants_1.getEndpointDetails)(api_1.SUBSCRIPTION_TEMPLATE_ROUTES.GET(':uuid'), 'get', 'Get subscription template by uuid', { scope: 'get', kind: 'read' });
    GetSubscriptionTemplateCommand.RequestParamSchema = zod_1.z.object({
        uuid: zod_1.z.uuid(),
    });
    GetSubscriptionTemplateCommand.ResponseSchema = zod_1.z.object({
        response: models_1.SubscriptionTemplateSchema,
    });
})(GetSubscriptionTemplateCommand || (exports.GetSubscriptionTemplateCommand = GetSubscriptionTemplateCommand = {}));
