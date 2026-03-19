"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.GetSubscriptionByUuidCommand = void 0;
const zod_1 = require("zod");
const api_1 = require("../../../api");
const models_1 = require("../../../models");
const constants_1 = require("../../../constants");
var GetSubscriptionByUuidCommand;
(function (GetSubscriptionByUuidCommand) {
    GetSubscriptionByUuidCommand.url = api_1.REST_API.SUBSCRIPTIONS.GET_BY.UUID;
    GetSubscriptionByUuidCommand.TSQ_url = GetSubscriptionByUuidCommand.url(':uuid');
    GetSubscriptionByUuidCommand.endpointDetails = (0, constants_1.getEndpointDetails)(api_1.SUBSCRIPTIONS_ROUTES.GET_BY.UUID(':uuid'), 'get', 'Get subscription by uuid');
    GetSubscriptionByUuidCommand.RequestSchema = zod_1.z.object({
        uuid: zod_1.z.string(),
    });
    GetSubscriptionByUuidCommand.ResponseSchema = zod_1.z.object({
        response: models_1.SubscriptionInfoSchema,
    });
})(GetSubscriptionByUuidCommand || (exports.GetSubscriptionByUuidCommand = GetSubscriptionByUuidCommand = {}));
