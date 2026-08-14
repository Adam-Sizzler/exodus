"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.GetSubscriptionRequestHistoryCommand = void 0;
const zod_1 = require("zod");
const api_1 = require("../../api");
const constants_1 = require("../../constants");
const models_1 = require("../../models");
var GetSubscriptionRequestHistoryCommand;
(function (GetSubscriptionRequestHistoryCommand) {
    GetSubscriptionRequestHistoryCommand.url = api_1.REST_API.SUBSCRIPTION_REQUEST_HISTORY.GET;
    GetSubscriptionRequestHistoryCommand.TSQ_url = GetSubscriptionRequestHistoryCommand.url;
    GetSubscriptionRequestHistoryCommand.endpointDetails = (0, constants_1.getEndpointDetails)(api_1.SUBSCRIPTION_REQUEST_HISTORY_ROUTES.GET, 'get', 'Get all subscription request history', { scope: 'list', kind: 'read' }, 'Please note that the filters here are primarily intended for use by the frontend and rely on expensive operators such as LIKE under the hood. Misusing these filters may negatively impact the performance of your database.');
    GetSubscriptionRequestHistoryCommand.RequestQuerySchema = models_1.TanstackQueryRequestQuerySchema;
    GetSubscriptionRequestHistoryCommand.ResponseSchema = zod_1.z.object({
        response: zod_1.z.object({
            records: zod_1.z.array(models_1.SubscriptionRequestHistorySchema),
            total: zod_1.z.number(),
        }),
    });
})(GetSubscriptionRequestHistoryCommand || (exports.GetSubscriptionRequestHistoryCommand = GetSubscriptionRequestHistoryCommand = {}));
