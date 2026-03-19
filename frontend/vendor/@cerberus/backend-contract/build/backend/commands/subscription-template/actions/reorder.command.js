"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.ReorderSubscriptionTemplateCommand = void 0;
const zod_1 = require("zod");
const api_1 = require("../../../api");
const models_1 = require("../../../models");
const constants_1 = require("../../../constants");
var ReorderSubscriptionTemplateCommand;
(function (ReorderSubscriptionTemplateCommand) {
    ReorderSubscriptionTemplateCommand.url = api_1.REST_API.SUBSCRIPTION_TEMPLATE.ACTIONS.REORDER;
    ReorderSubscriptionTemplateCommand.TSQ_url = ReorderSubscriptionTemplateCommand.url;
    ReorderSubscriptionTemplateCommand.endpointDetails = (0, constants_1.getEndpointDetails)(api_1.SUBSCRIPTION_TEMPLATE_ROUTES.ACTIONS.REORDER, 'post', 'Reorder subscription templates');
    ReorderSubscriptionTemplateCommand.RequestSchema = zod_1.z.object({
        items: zod_1.z.array(models_1.SubscriptionTemplateSchema.pick({
            viewPosition: true,
            uuid: true,
        })),
    });
    ReorderSubscriptionTemplateCommand.ResponseSchema = zod_1.z.object({
        response: zod_1.z.object({
            total: zod_1.z.number(),
            templates: zod_1.z.array(models_1.SubscriptionTemplateSchema),
        }),
    });
})(ReorderSubscriptionTemplateCommand || (exports.ReorderSubscriptionTemplateCommand = ReorderSubscriptionTemplateCommand = {}));
