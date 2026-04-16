"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.ConfigProfileSchema = void 0;
const zod_1 = require("zod");
const config_profile_inbounds_schema_1 = require("./config-profile-inbounds.schema");
exports.ConfigProfileSchema = zod_1.z.object({
    uuid: zod_1.z.string().uuid(),
    viewPosition: zod_1.z.number().int(),
    name: zod_1.z.string(),
    config: zod_1.z.unknown(),
    inbounds: zod_1.z.array(config_profile_inbounds_schema_1.ConfigProfileInboundsSchema),
    nodes: zod_1.z.array(zod_1.z.object({
        uuid: zod_1.z.string().uuid(),
        name: zod_1.z.string(),
        countryCode: zod_1.z.string(),
    })),
    createdAt: zod_1.z
        .string()
        .datetime()
        .transform((str) => new Date(str)),
    updatedAt: zod_1.z
        .string()
        .datetime()
        .transform((str) => new Date(str)),
});
