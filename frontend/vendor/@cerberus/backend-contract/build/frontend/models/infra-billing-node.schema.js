"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.InfraBillingNodeSchema = void 0;
const zod_1 = require("zod");
const infra_provider_schema_1 = require("./infra-provider.schema");
const nodes_schema_1 = require("./nodes.schema");
exports.InfraBillingNodeSchema = zod_1.z.object({
    uuid: zod_1.z.string().uuid(),
    nodeUuid: zod_1.z.string().uuid(),
    providerUuid: zod_1.z.string().uuid(),
    provider: infra_provider_schema_1.PartialInfraProviderSchema.pick({
        uuid: true,
        name: true,
        loginUrl: true,
        faviconLink: true,
    }),
    node: nodes_schema_1.NodesSchema.pick({
        uuid: true,
        name: true,
        countryCode: true,
    }),
    nextBillingAt: zod_1.z
        .string()
        .datetime()
        .transform((str) => new Date(str)),
    createdAt: zod_1.z
        .string()
        .datetime()
        .transform((str) => new Date(str)),
    updatedAt: zod_1.z
        .string()
        .datetime()
        .transform((str) => new Date(str)),
});
