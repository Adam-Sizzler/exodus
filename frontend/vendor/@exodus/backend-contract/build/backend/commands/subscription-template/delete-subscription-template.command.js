"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.DeleteSubscriptionTemplateCommand = void 0;
const zod_1 = require("zod");
const api_1 = require("../../api");
const constants_1 = require("../../constants");
var DeleteSubscriptionTemplateCommand;
(function (DeleteSubscriptionTemplateCommand) {
    DeleteSubscriptionTemplateCommand.url = api_1.REST_API.SUBSCRIPTION_TEMPLATE.DELETE;
    DeleteSubscriptionTemplateCommand.TSQ_url = DeleteSubscriptionTemplateCommand.url(':uuid');
    DeleteSubscriptionTemplateCommand.endpointDetails = (0, constants_1.getEndpointDetails)(api_1.SUBSCRIPTION_TEMPLATE_ROUTES.DELETE(':uuid'), 'delete', 'Delete subscription template');
    DeleteSubscriptionTemplateCommand.RequestSchema = zod_1.z.object({
        uuid: zod_1.z.string().uuid(),
    });
    DeleteSubscriptionTemplateCommand.ResponseSchema = zod_1.z.object({
        response: zod_1.z.object({
            isDeleted: zod_1.z.boolean(),
        }),
    });
})(DeleteSubscriptionTemplateCommand || (exports.DeleteSubscriptionTemplateCommand = DeleteSubscriptionTemplateCommand = {}));
