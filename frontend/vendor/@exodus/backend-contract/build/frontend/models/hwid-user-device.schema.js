"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.HwidUserDeviceSchema = void 0;
const zod_1 = require("zod");
exports.HwidUserDeviceSchema = zod_1.z.object({
    hwid: zod_1.z.string(),
    userId: zod_1.z.number(),
    platform: zod_1.z.nullable(zod_1.z.string()),
    osVersion: zod_1.z.nullable(zod_1.z.string()),
    deviceModel: zod_1.z.nullable(zod_1.z.string()),
    userAgent: zod_1.z.nullable(zod_1.z.string()),
    requestIp: zod_1.z.nullable(zod_1.z.string()),
    createdAt: zod_1.z.iso.datetime()
        .transform((str) => new Date(str)),
    updatedAt: zod_1.z.iso.datetime()
        .transform((str) => new Date(str)),
});
