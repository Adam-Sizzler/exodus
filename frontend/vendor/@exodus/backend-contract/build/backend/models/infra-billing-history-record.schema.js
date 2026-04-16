"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.InfraBillingHistoryRecordSchema = void 0;
const zod_1 = require("zod");
const infra_provider_schema_1 = require("./infra-provider.schema");
exports.InfraBillingHistoryRecordSchema = zod_1.z.object({
    uuid: zod_1.z.string().uuid(),
    providerUuid: zod_1.z.string().uuid(),
    amount: zod_1.z.number(),
    billedAt: zod_1.z
        .string()
        .datetime()
        .transform((str) => new Date(str)),
    provider: infra_provider_schema_1.PartialInfraProviderSchema.omit({
        createdAt: true,
        updatedAt: true,
        loginUrl: true,
    }),
});
