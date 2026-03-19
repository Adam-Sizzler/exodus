"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.PartialInfraProviderSchema = exports.InfraProviderSchema = void 0;
const zod_1 = require("zod");
exports.InfraProviderSchema = zod_1.z.object({
    uuid: zod_1.z.string().uuid(),
    name: zod_1.z.string(),
    faviconLink: zod_1.z.nullable(zod_1.z.string()),
    loginUrl: zod_1.z.nullable(zod_1.z.string()),
    createdAt: zod_1.z
        .string()
        .datetime()
        .transform((str) => new Date(str)),
    updatedAt: zod_1.z
        .string()
        .datetime()
        .transform((str) => new Date(str)),
    billingHistory: zod_1.z.object({
        totalAmount: zod_1.z.number(),
        totalBills: zod_1.z.number(),
    }),
    billingNodes: zod_1.z.array(zod_1.z.object({
        nodeUuid: zod_1.z.string().uuid(),
        name: zod_1.z.string(),
        countryCode: zod_1.z.string(),
    })),
});
exports.PartialInfraProviderSchema = zod_1.z.object({
    uuid: zod_1.z.string().uuid(),
    name: zod_1.z.string(),
    faviconLink: zod_1.z.nullable(zod_1.z.string()),
    loginUrl: zod_1.z.nullable(zod_1.z.string()),
    createdAt: zod_1.z
        .string()
        .datetime()
        .transform((str) => new Date(str)),
    updatedAt: zod_1.z
        .string()
        .datetime()
        .transform((str) => new Date(str)),
});
