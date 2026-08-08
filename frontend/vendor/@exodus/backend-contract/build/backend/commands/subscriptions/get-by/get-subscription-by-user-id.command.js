"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.GetSubscriptionByIdCommand = void 0;
const zod_1 = require("zod");
const api_1 = require("../../../api");
const constants_1 = require("../../../constants");
const models_1 = require("../../../models");
var GetSubscriptionByIdCommand;
(function (GetSubscriptionByIdCommand) {
    GetSubscriptionByIdCommand.url = api_1.REST_API.SUBSCRIPTIONS.GET_BY.ID;
    GetSubscriptionByIdCommand.TSQ_url = GetSubscriptionByIdCommand.url(':userId');
    GetSubscriptionByIdCommand.endpointDetails = (0, constants_1.getEndpointDetails)(api_1.SUBSCRIPTIONS_ROUTES.GET_BY.ID(':userId'), 'get', 'Get subscription by User ID', { scope: 'by-id', kind: 'read' });
    GetSubscriptionByIdCommand.RequestParamSchema = zod_1.z.object({
        userId: models_1.numberParamSchema.describe('User ID'),
    });
    GetSubscriptionByIdCommand.ResponseSchema = zod_1.z.object({
        response: models_1.SubscriptionInfoSchema,
    });

})(GetSubscriptionByIdCommand || (exports.GetSubscriptionByIdCommand = GetSubscriptionByIdCommand = {}));
