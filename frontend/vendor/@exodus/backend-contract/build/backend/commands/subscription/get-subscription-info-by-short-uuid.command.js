"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.GetSubscriptionInfoByShortUuidCommand = void 0;
const zod_1 = require("zod");
const api_1 = require("../../api");
const models_1 = require("../../models");
const constants_1 = require("../../constants");
var GetSubscriptionInfoByShortUuidCommand;
(function (GetSubscriptionInfoByShortUuidCommand) {
    GetSubscriptionInfoByShortUuidCommand.url = api_1.REST_API.SUBSCRIPTION.GET_INFO;
    GetSubscriptionInfoByShortUuidCommand.TSQ_url = GetSubscriptionInfoByShortUuidCommand.url(':shortUuid');
    GetSubscriptionInfoByShortUuidCommand.endpointDetails = (0, constants_1.getEndpointDetails)(api_1.SUBSCRIPTION_ROUTES.GET_INFO(':shortUuid'), 'get', 'Get Subscription Info by Short UUID');
    GetSubscriptionInfoByShortUuidCommand.RequestSchema = zod_1.z.object({
        shortUuid: zod_1.z.string(),
    });
    GetSubscriptionInfoByShortUuidCommand.ResponseSchema = zod_1.z.object({
        response: models_1.SubscriptionInfoSchema,
    });
})(GetSubscriptionInfoByShortUuidCommand || (exports.GetSubscriptionInfoByShortUuidCommand = GetSubscriptionInfoByShortUuidCommand = {}));
