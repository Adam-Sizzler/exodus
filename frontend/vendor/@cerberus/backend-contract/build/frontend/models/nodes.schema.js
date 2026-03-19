"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.NodesSchema = void 0;
const zod_1 = require("zod");
const config_profile_inbounds_schema_1 = require("./config-profile-inbounds.schema");
const infra_provider_schema_1 = require("./infra-provider.schema");
exports.NodesSchema = zod_1.z.object({
    uuid: zod_1.z.string().uuid(),
    name: zod_1.z.string(),
    address: zod_1.z.string(),
    port: zod_1.z.nullable(zod_1.z.number().int()),
    isConnected: zod_1.z.boolean(),
    isDisabled: zod_1.z.boolean(),
    isConnecting: zod_1.z.boolean(),
    lastStatusChange: zod_1.z.nullable(zod_1.z
        .string()
        .datetime()
        .transform((str) => new Date(str))),
    lastStatusMessage: zod_1.z.nullable(zod_1.z.string()),
    xrayVersion: zod_1.z.nullable(zod_1.z.string()),
    nodeVersion: zod_1.z.nullable(zod_1.z.string()),
    xrayUptime: zod_1.z.string(),
    isTrafficTrackingActive: zod_1.z.boolean(),
    trafficResetDay: zod_1.z.nullable(zod_1.z.number().int()),
    trafficLimitBytes: zod_1.z.nullable(zod_1.z.number()),
    trafficUsedBytes: zod_1.z.nullable(zod_1.z.number()),
    notifyPercent: zod_1.z.nullable(zod_1.z.number().int()),
    usersOnline: zod_1.z.nullable(zod_1.z.number().int()),
    viewPosition: zod_1.z.number().int(),
    countryCode: zod_1.z.string(),
    consumptionMultiplier: zod_1.z.number(),
    tags: zod_1.z.array(zod_1.z.string()),
    cpuCount: zod_1.z.nullable(zod_1.z.number().int()),
    cpuModel: zod_1.z.nullable(zod_1.z.string()),
    totalRam: zod_1.z.nullable(zod_1.z.string()),
    createdAt: zod_1.z
        .string()
        .datetime()
        .transform((str) => new Date(str)),
    updatedAt: zod_1.z
        .string()
        .datetime()
        .transform((str) => new Date(str)),
    configProfile: zod_1.z.object({
        activeConfigProfileUuid: zod_1.z.nullable(zod_1.z.string().uuid()),
        activeInbounds: zod_1.z.array(config_profile_inbounds_schema_1.ConfigProfileInboundsSchema),
    }),
    providerUuid: zod_1.z.nullable(zod_1.z.string().uuid()),
    provider: zod_1.z.nullable(infra_provider_schema_1.PartialInfraProviderSchema),
});
