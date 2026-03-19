"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.InternalSquadSchema = void 0;
const zod_1 = require("zod");
const config_profile_inbounds_schema_1 = require("./config-profile-inbounds.schema");
exports.InternalSquadSchema = zod_1.z.object({
    uuid: zod_1.z.string().uuid(),
    viewPosition: zod_1.z.number().int(),
    name: zod_1.z.string(),
    info: zod_1.z.object({
        membersCount: zod_1.z.number(),
        inboundsCount: zod_1.z.number(),
    }),
    inbounds: zod_1.z.array(config_profile_inbounds_schema_1.ConfigProfileInboundsSchema),
    createdAt: zod_1.z
        .string()
        .datetime()
        .transform((str) => new Date(str)),
    updatedAt: zod_1.z
        .string()
        .datetime()
        .transform((str) => new Date(str)),
});
