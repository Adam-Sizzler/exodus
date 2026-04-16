"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.GetSubscriptionSettingsCommand = void 0;
const zod_1 = require("zod");
const api_1 = require("../../api");
const models_1 = require("../../models");
const constants_1 = require("../../constants");
var GetSubscriptionSettingsCommand;
(function (GetSubscriptionSettingsCommand) {
    GetSubscriptionSettingsCommand.url = api_1.REST_API.SUBSCRIPTION_SETTINGS.GET;
    GetSubscriptionSettingsCommand.TSQ_url = GetSubscriptionSettingsCommand.url;
    GetSubscriptionSettingsCommand.endpointDetails = (0, constants_1.getEndpointDetails)(api_1.SUBSCRIPTION_SETTINGS_ROUTES.GET, 'get', 'Get subscription settings');
    GetSubscriptionSettingsCommand.ResponseSchema = zod_1.z.object({
        response: models_1.SubscriptionSettingsSchema,
    });
})(GetSubscriptionSettingsCommand || (exports.GetSubscriptionSettingsCommand = GetSubscriptionSettingsCommand = {}));
