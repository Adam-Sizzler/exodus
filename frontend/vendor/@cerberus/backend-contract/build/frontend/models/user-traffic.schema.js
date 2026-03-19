"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.UserTrafficSchema = void 0;
const zod_1 = require("zod");
exports.UserTrafficSchema = zod_1.z.object({
    usedTrafficBytes: zod_1.z.number(),
    lifetimeUsedTrafficBytes: zod_1.z.number(),
    onlineAt: zod_1.z.nullable(zod_1.z
        .string()
        .datetime()
        .transform((str) => new Date(str))),
    firstConnectedAt: zod_1.z.nullable(zod_1.z
        .string()
        .datetime()
        .transform((str) => new Date(str))),
    lastConnectedNodeUuid: zod_1.z.nullable(zod_1.z.string().uuid()),
});
