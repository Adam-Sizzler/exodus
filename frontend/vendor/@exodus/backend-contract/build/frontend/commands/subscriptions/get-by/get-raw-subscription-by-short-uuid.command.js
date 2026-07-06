"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.GetRawSubscriptionByShortUuidCommand = void 0;
const zod_1 = require("zod");
const api_1 = require("../../../api");
const constants_1 = require("../../../constants");
const models_1 = require("../../../models");
var GetRawSubscriptionByShortUuidCommand;
(function (GetRawSubscriptionByShortUuidCommand) {
    GetRawSubscriptionByShortUuidCommand.url = api_1.REST_API.SUBSCRIPTIONS.GET_BY.SHORT_UUID_RAW;
    GetRawSubscriptionByShortUuidCommand.TSQ_url = GetRawSubscriptionByShortUuidCommand.url(':shortUuid');
    GetRawSubscriptionByShortUuidCommand.endpointDetails = (0, constants_1.getEndpointDetails)(api_1.SUBSCRIPTIONS_ROUTES.GET_BY.SHORT_UUID_RAW(':shortUuid'), 'get', 'Get Raw Subscription by Short UUID', { scope: 'raw', kind: 'read' });
    GetRawSubscriptionByShortUuidCommand.RequestSchema = zod_1.z.object({
        shortUuid: zod_1.z.string(),
    });
    GetRawSubscriptionByShortUuidCommand.RequestQuerySchema = zod_1.z.object({
        withDisabledHosts: zod_1.z
            .string()
            .transform((str) => str === 'true')
            .optional()
            .default('false'),
    });
    GetRawSubscriptionByShortUuidCommand.ResponseSchema = zod_1.z.object({
        response: zod_1.z.object({
            user: models_1.ExtendedUsersSchema,
            convertedUserInfo: zod_1.z.object({
                daysLeft: zod_1.z.number(),
                trafficLimit: zod_1.z.string(),
                trafficUsed: zod_1.z.string(),
                lifetimeTrafficUsed: zod_1.z.string(),
                hwidCheckup: zod_1.z
                    .object({
                    subscriptionAllowed: zod_1.z.boolean(),
                    maxDeviceReached: zod_1.z.boolean(),
                    hwidNotSupported: zod_1.z.boolean(),
                    limitBypassed: zod_1.z.boolean(),
                })
                    .nullable(),
            }),
            headers: zod_1.z.record(zod_1.z.string(), zod_1.z.string().optional()),
            resolvedProxyConfigs: zod_1.z.array(models_1.ResolvedProxyConfigSchema),
        }),
    });
})(GetRawSubscriptionByShortUuidCommand || (exports.GetRawSubscriptionByShortUuidCommand = GetRawSubscriptionByShortUuidCommand = {}));
