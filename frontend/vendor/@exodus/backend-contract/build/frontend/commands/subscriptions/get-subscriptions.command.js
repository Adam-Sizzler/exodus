"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.GetSubscriptionsCommand = void 0;
const zod_1 = require("zod");
const api_1 = require("../../api");
const constants_1 = require("../../constants");
var GetSubscriptionsCommand;
(function (GetSubscriptionsCommand) {
    GetSubscriptionsCommand.url = api_1.REST_API.SUBSCRIPTIONS.GET;
    GetSubscriptionsCommand.TSQ_url = GetSubscriptionsCommand.url;
    GetSubscriptionsCommand.endpointDetails = (0, constants_1.getEndpointDetails)(api_1.SUBSCRIPTIONS_ROUTES.GET, 'get', 'Get all subscriptions', { scope: 'list', kind: 'read' });
    GetSubscriptionsCommand.RequestQuerySchema = zod_1.z.object({
        start: zod_1.z.coerce
            .number()
            .default(0)
            .describe('Start index (offset) of the users to return, default is 0'),
        size: zod_1.z.coerce
            .number()
            .min(1)
            .max(500)
            .describe('Number of subscriptions to return, no more than 500')
            .default(25),
    });
    GetSubscriptionsCommand.ResponseSchema = zod_1.z.object({
        response: zod_1.z.object({
            subscriptions: zod_1.z.array(zod_1.z.object({
                isFound: zod_1.z.boolean(),
                user: zod_1.z.object({
                    shortUuid: zod_1.z.string(),
                    daysLeft: zod_1.z.number(),
                    trafficUsed: zod_1.z.string(),
                    trafficLimit: zod_1.z.string(),
                    lifetimeTrafficUsed: zod_1.z.string(),
                    trafficUsedBytes: zod_1.z.string(),
                    trafficLimitBytes: zod_1.z.string(),
                    lifetimeTrafficUsedBytes: zod_1.z.string(),
                    username: zod_1.z.string(),
                    expiresAt: zod_1.z.iso.datetime().transform((str) => new Date(str)),
                    isActive: zod_1.z.boolean(),
                    userStatus: zod_1.z.enum(constants_1.USERS_STATUS),
                    trafficLimitStrategy: zod_1.z.enum(constants_1.RESET_PERIODS),
                }),
                links: zod_1.z.array(zod_1.z.string()),
                ssConfLinks: zod_1.z.record(zod_1.z.string(), zod_1.z.string()),
                subscriptionUrl: zod_1.z.string(),
            })),
            total: zod_1.z.number(),
        }),
    });

})(GetSubscriptionsCommand || (exports.GetSubscriptionsCommand = GetSubscriptionsCommand = {}));
