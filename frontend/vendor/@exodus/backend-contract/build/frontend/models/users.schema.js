"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.UsersSchema = void 0;
const zod_1 = require("zod");
const constants_1 = require("../constants");
exports.UsersSchema = zod_1.z.object({
    id: zod_1.z.number(),
    shortUuid: zod_1.z.string(),
    username: zod_1.z.string(),
    status: zod_1.z.enum(constants_1.USERS_STATUS),
    trafficLimitBytes: zod_1.z.number(),
    trafficLimitStrategy: zod_1.z.enum(constants_1.RESET_PERIODS).describe('Available reset periods'),
    expireAt: zod_1.z.iso.datetime().transform((str) => new Date(str)),
    telegramId: zod_1.z.nullable(zod_1.z.number()),
    email: zod_1.z.nullable(zod_1.z.email()),
    description: zod_1.z.nullable(zod_1.z.string()),
    tag: zod_1.z.nullable(zod_1.z.string()),
    hwidDeviceLimit: zod_1.z.nullable(zod_1.z.int()),
    externalSquadUuid: zod_1.z.nullable(zod_1.z.uuid()),
    trojanPassword: zod_1.z.string(),
    vlessUuid: zod_1.z.uuid(),
    ssPassword: zod_1.z.string(),
    naivePassword: zod_1.z.string(),
    shadowtlsPassword: zod_1.z.string(),
    hysteria2Password: zod_1.z.string(),
    anytlsPassword: zod_1.z.string(),
    lastTriggeredThreshold: zod_1.z.int(),
    subRevokedAt: zod_1.z.nullable(zod_1.z.iso.datetime().transform((str) => new Date(str))),
    lastTrafficResetAt: zod_1.z.nullable(zod_1.z.iso.datetime().transform((str) => new Date(str))),
    createdAt: zod_1.z.iso.datetime().transform((str) => new Date(str)),
    updatedAt: zod_1.z.iso.datetime().transform((str) => new Date(str)),
});
