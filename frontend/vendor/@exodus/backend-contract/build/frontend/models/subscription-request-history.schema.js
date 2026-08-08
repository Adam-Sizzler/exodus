"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.SubscriptionRequestHistorySchema = void 0;
const zod_1 = require("zod");
exports.SubscriptionRequestHistorySchema = zod_1.z.object({
    id: zod_1.z.number(),
    userId: zod_1.z.number(),
    requestIp: zod_1.z.nullable(zod_1.z.string()),
    userAgent: zod_1.z.nullable(zod_1.z.string()),
    requestAt: zod_1.z.iso.datetime().transform((str) => new Date(str)),
});
