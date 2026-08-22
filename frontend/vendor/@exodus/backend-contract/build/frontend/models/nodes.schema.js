"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.NodesSchema = void 0;
const zod_1 = require("zod");
const config_profile_inbounds_schema_1 = require("./config-profile-inbounds.schema");
const infra_provider_schema_1 = require("./infra-provider.schema");
const node_ips_schema_1 = require("./node-ips.schema");
const node_system_schema_1 = require("./node-system.schema");
exports.NodesSchema = zod_1.z.object({
    uuid: zod_1.z.uuid(),
    id: zod_1.z.number(),
    name: zod_1.z.string(),
    address: zod_1.z.string(),
    port: zod_1.z.nullable(zod_1.z.int()),
    proxyUrl: zod_1.z.nullable(zod_1.z.string()),
    apiSchema: zod_1.z.enum(['mtls', 'tls']).default('mtls'),
    apiPath: zod_1.z.nullable(zod_1.z.string()).default('/'),
    grpcAuthToken: zod_1.z.nullable(zod_1.z.string()),
    isConnected: zod_1.z.boolean(),
    isDisabled: zod_1.z.boolean(),
    isConnecting: zod_1.z.boolean(),
    lastStatusChange: zod_1.z.nullable(zod_1.z.iso.datetime().transform((str) => new Date(str))),
    lastStatusMessage: zod_1.z.nullable(zod_1.z.string()),
    isTrafficTrackingActive: zod_1.z.boolean(),
    trafficResetDay: zod_1.z.nullable(zod_1.z.int()),
    trafficLimitBytes: zod_1.z.nullable(zod_1.z.number()),
    trafficUsedBytes: zod_1.z.nullable(zod_1.z.number()),
    notifyPercent: zod_1.z.nullable(zod_1.z.int()),
    viewPosition: zod_1.z.int(),
    countryCode: zod_1.z.string(),
    consumptionMultiplier: zod_1.z.number(),
    nodeConsumptionMultiplier: zod_1.z.number(),
    tags: zod_1.z.array(zod_1.z.string()),
    integrationUuids: zod_1.z.array(zod_1.z.uuid()),
    ips: node_ips_schema_1.NodeIpsSchema,
    createdAt: zod_1.z.iso.datetime().transform((str) => new Date(str)),
    updatedAt: zod_1.z.iso.datetime().transform((str) => new Date(str)),
    configProfile: zod_1.z.object({
        activeConfigProfileUuid: zod_1.z.nullable(zod_1.z.uuid()),
        activeInbounds: zod_1.z.array(config_profile_inbounds_schema_1.ConfigProfileInboundsSchema),
    }),
    providerUuid: zod_1.z.nullable(zod_1.z.uuid()),
    provider: zod_1.z.nullable(infra_provider_schema_1.PartialInfraProviderSchema),
    activePluginUuid: zod_1.z.nullable(zod_1.z.uuid()),
    system: zod_1.z.nullable(node_system_schema_1.NodeSystemSchema),
    versions: zod_1.z.nullable(zod_1.z.object({
        singbox: zod_1.z.string(),
        node: zod_1.z.string(),
    })),
    singboxUptime: zod_1.z.number(),
    usersOnline: zod_1.z.number(),
    note: zod_1.z.nullable(zod_1.z.string()),
});
