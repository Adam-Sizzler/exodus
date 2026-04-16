"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.ConfigProfileInboundsSchema = void 0;
const zod_1 = require("zod");
exports.ConfigProfileInboundsSchema = zod_1.z.object({
    uuid: zod_1.z.string().uuid(),
    profileUuid: zod_1.z.string().uuid(),
    tag: zod_1.z.string(),
    type: zod_1.z.string(),
    network: zod_1.z.nullable(zod_1.z.string()),
    security: zod_1.z.nullable(zod_1.z.string()),
    port: zod_1.z.nullable(zod_1.z.number()),
    rawInbound: zod_1.z.nullable(zod_1.z.unknown()),
});
