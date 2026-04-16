"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.GetSubscriptionRequestHistoryStatsCommand = void 0;
const zod_1 = require("zod");
const api_1 = require("../../api");
const constants_1 = require("../../constants");
var GetSubscriptionRequestHistoryStatsCommand;
(function (GetSubscriptionRequestHistoryStatsCommand) {
    GetSubscriptionRequestHistoryStatsCommand.url = api_1.REST_API.SUBSCRIPTION_REQUEST_HISTORY.STATS;
    GetSubscriptionRequestHistoryStatsCommand.TSQ_url = GetSubscriptionRequestHistoryStatsCommand.url;
    GetSubscriptionRequestHistoryStatsCommand.endpointDetails = (0, constants_1.getEndpointDetails)(api_1.SUBSCRIPTION_REQUEST_HISTORY_ROUTES.STATS, 'get', 'Get subscription request history stats');
    GetSubscriptionRequestHistoryStatsCommand.ResponseSchema = zod_1.z.object({
        response: zod_1.z.object({
            byParsedApp: zod_1.z.array(zod_1.z.object({
                app: zod_1.z.string(),
                count: zod_1.z.number(),
            })),
            hourlyRequestStats: zod_1.z.array(zod_1.z.object({
                dateTime: zod_1.z
                    .string()
                    .datetime()
                    .transform((str) => new Date(str)),
                requestCount: zod_1.z.number(),
            })),
        }),
    });
})(GetSubscriptionRequestHistoryStatsCommand || (exports.GetSubscriptionRequestHistoryStatsCommand = GetSubscriptionRequestHistoryStatsCommand = {}));
