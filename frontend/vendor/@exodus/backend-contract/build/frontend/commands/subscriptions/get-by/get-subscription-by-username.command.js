"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.GetSubscriptionByUsernameCommand = void 0;
const zod_1 = require("zod");
const api_1 = require("../../../api");
const models_1 = require("../../../models");
const constants_1 = require("../../../constants");
var GetSubscriptionByUsernameCommand;
(function (GetSubscriptionByUsernameCommand) {
    GetSubscriptionByUsernameCommand.url = api_1.REST_API.SUBSCRIPTIONS.GET_BY.USERNAME;
    GetSubscriptionByUsernameCommand.TSQ_url = GetSubscriptionByUsernameCommand.url(':username');
    GetSubscriptionByUsernameCommand.endpointDetails = (0, constants_1.getEndpointDetails)(api_1.SUBSCRIPTIONS_ROUTES.GET_BY.USERNAME(':username'), 'get', 'Get subscription by username');
    GetSubscriptionByUsernameCommand.RequestSchema = zod_1.z.object({
        username: zod_1.z.string(),
    });
    GetSubscriptionByUsernameCommand.ResponseSchema = zod_1.z.object({
        response: models_1.SubscriptionInfoSchema,
    });
})(GetSubscriptionByUsernameCommand || (exports.GetSubscriptionByUsernameCommand = GetSubscriptionByUsernameCommand = {}));
