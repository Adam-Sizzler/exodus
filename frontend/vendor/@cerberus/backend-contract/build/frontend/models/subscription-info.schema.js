"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.SubscriptionInfoSchema = void 0;
const zod_1 = require("zod");
const constants_1 = require("../constants");
exports.SubscriptionInfoSchema = zod_1.z.object({
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
        expiresAt: zod_1.z
            .string()
            .datetime()
            .transform((str) => new Date(str)),
        isActive: zod_1.z.boolean(),
        userStatus: zod_1.z.nativeEnum(constants_1.USERS_STATUS),
        trafficLimitStrategy: zod_1.z.nativeEnum(constants_1.RESET_PERIODS),
    }),
    links: zod_1.z.array(zod_1.z.string()),
    ssConfLinks: zod_1.z.record(zod_1.z.string(), zod_1.z.string()),
    subscriptionUrl: zod_1.z.string(),
});
