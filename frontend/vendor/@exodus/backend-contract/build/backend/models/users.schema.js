"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.UsersSchema = void 0;
const zod_1 = require("zod");
const constants_1 = require("../constants");
exports.UsersSchema = zod_1.z.object({
    uuid: zod_1.z.string().uuid(),
    id: zod_1.z.number(),
    shortUuid: zod_1.z.string(),
    username: zod_1.z.string(),
    status: zod_1.z.nativeEnum(constants_1.USERS_STATUS).default(constants_1.USERS_STATUS.ACTIVE),
    trafficLimitBytes: zod_1.z.number().int().default(0),
    trafficLimitStrategy: zod_1.z
        .nativeEnum(constants_1.RESET_PERIODS, {
        description: 'Available reset periods',
    })
        .default(constants_1.RESET_PERIODS.NO_RESET),
    expireAt: zod_1.z
        .string()
        .datetime()
        .transform((str) => new Date(str)),
    telegramId: zod_1.z.nullable(zod_1.z.number().int()),
    email: zod_1.z.nullable(zod_1.z.string().email()),
    description: zod_1.z.nullable(zod_1.z.string()),
    tag: zod_1.z.nullable(zod_1.z.string()),
    hwidDeviceLimit: zod_1.z.nullable(zod_1.z.number().int()),
    externalSquadUuid: zod_1.z.nullable(zod_1.z.string().uuid()),
    trojanPassword: zod_1.z.string(),
    vlessUuid: zod_1.z.string().uuid(),
    ssPassword: zod_1.z.string(),
    lastTriggeredThreshold: zod_1.z.number().int().default(0),
    subRevokedAt: zod_1.z.nullable(zod_1.z
        .string()
        .datetime()
        .transform((str) => new Date(str))),
    subLastUserAgent: zod_1.z.nullable(zod_1.z.string()),
    subLastOpenedAt: zod_1.z.nullable(zod_1.z
        .string()
        .datetime()
        .transform((str) => new Date(str))),
    lastTrafficResetAt: zod_1.z.nullable(zod_1.z
        .string()
        .datetime()
        .transform((str) => new Date(str))),
    createdAt: zod_1.z
        .string()
        .datetime()
        .transform((str) => new Date(str)),
    updatedAt: zod_1.z
        .string()
        .datetime()
        .transform((str) => new Date(str)),
});
