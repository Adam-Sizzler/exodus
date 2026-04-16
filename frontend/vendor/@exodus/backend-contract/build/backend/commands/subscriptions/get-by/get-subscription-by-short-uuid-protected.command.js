"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.GetSubscriptionByShortUuidProtectedCommand = void 0;
const zod_1 = require("zod");
const api_1 = require("../../../api");
const models_1 = require("../../../models");
const constants_1 = require("../../../constants");
var GetSubscriptionByShortUuidProtectedCommand;
(function (GetSubscriptionByShortUuidProtectedCommand) {
    GetSubscriptionByShortUuidProtectedCommand.url = api_1.REST_API.SUBSCRIPTIONS.GET_BY.SHORT_UUID;
    GetSubscriptionByShortUuidProtectedCommand.TSQ_url = GetSubscriptionByShortUuidProtectedCommand.url(':shortUuid');
    GetSubscriptionByShortUuidProtectedCommand.endpointDetails = (0, constants_1.getEndpointDetails)(api_1.SUBSCRIPTIONS_ROUTES.GET_BY.SHORT_UUID(':shortUuid'), 'get', 'Get subscription by short uuid (protected route)');
    GetSubscriptionByShortUuidProtectedCommand.RequestSchema = zod_1.z.object({
        shortUuid: zod_1.z.string(),
    });
    GetSubscriptionByShortUuidProtectedCommand.ResponseSchema = zod_1.z.object({
        response: models_1.SubscriptionInfoSchema,
    });
})(GetSubscriptionByShortUuidProtectedCommand || (exports.GetSubscriptionByShortUuidProtectedCommand = GetSubscriptionByShortUuidProtectedCommand = {}));
