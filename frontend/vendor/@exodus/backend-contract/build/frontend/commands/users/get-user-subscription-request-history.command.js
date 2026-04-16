"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.GetUserSubscriptionRequestHistoryCommand = void 0;
const zod_1 = require("zod");
const constants_1 = require("../../constants");
const api_1 = require("../../api");
var GetUserSubscriptionRequestHistoryCommand;
(function (GetUserSubscriptionRequestHistoryCommand) {
    GetUserSubscriptionRequestHistoryCommand.url = api_1.REST_API.USERS.SUBSCRIPTION_REQUEST_HISTORY;
    GetUserSubscriptionRequestHistoryCommand.TSQ_url = GetUserSubscriptionRequestHistoryCommand.url(':uuid');
    GetUserSubscriptionRequestHistoryCommand.endpointDetails = (0, constants_1.getEndpointDetails)(api_1.USERS_ROUTES.SUBSCRIPTION_REQUEST_HISTORY(':uuid'), 'get', 'Get user subscription request history, recent 24 records');
    GetUserSubscriptionRequestHistoryCommand.RequestSchema = zod_1.z.object({
        uuid: zod_1.z.string().uuid(),
    });
    GetUserSubscriptionRequestHistoryCommand.ResponseSchema = zod_1.z.object({
        response: zod_1.z.object({
            total: zod_1.z.number(),
            records: zod_1.z.array(zod_1.z.object({
                id: zod_1.z.number(),
                userUuid: zod_1.z.string().uuid(),
                requestAt: zod_1.z
                    .string()
                    .datetime()
                    .transform((str) => new Date(str)),
                requestIp: zod_1.z.string().optional().nullable(),
                userAgent: zod_1.z.string().optional().nullable(),
            })),
        }),
    });
})(GetUserSubscriptionRequestHistoryCommand || (exports.GetUserSubscriptionRequestHistoryCommand = GetUserSubscriptionRequestHistoryCommand = {}));
