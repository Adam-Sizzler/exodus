"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.LastConnectedNodeSchema = void 0;
const zod_1 = require("zod");
exports.LastConnectedNodeSchema = zod_1.z
    .object({
    connectedAt: zod_1.z
        .string()
        .datetime()
        .transform((str) => new Date(str)),
    nodeName: zod_1.z.string(),
    countryCode: zod_1.z.string(),
})
    .nullable();
