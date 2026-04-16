"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.DeleteSubscriptionPageConfigCommand = void 0;
const zod_1 = require("zod");
const api_1 = require("../../api");
const constants_1 = require("../../constants");
var DeleteSubscriptionPageConfigCommand;
(function (DeleteSubscriptionPageConfigCommand) {
    DeleteSubscriptionPageConfigCommand.url = api_1.REST_API.SUBSCRIPTION_PAGE_CONFIGS.DELETE;
    DeleteSubscriptionPageConfigCommand.TSQ_url = DeleteSubscriptionPageConfigCommand.url(':uuid');
    DeleteSubscriptionPageConfigCommand.endpointDetails = (0, constants_1.getEndpointDetails)(api_1.SUBSCRIPTION_PAGE_CONFIGS_ROUTES.DELETE(':uuid'), 'delete', 'Delete subscription page config');
    DeleteSubscriptionPageConfigCommand.RequestSchema = zod_1.z.object({
        uuid: zod_1.z.string().uuid(),
    });
    DeleteSubscriptionPageConfigCommand.ResponseSchema = zod_1.z.object({
        response: zod_1.z.object({
            isDeleted: zod_1.z.boolean(),
        }),
    });
})(DeleteSubscriptionPageConfigCommand || (exports.DeleteSubscriptionPageConfigCommand = DeleteSubscriptionPageConfigCommand = {}));
