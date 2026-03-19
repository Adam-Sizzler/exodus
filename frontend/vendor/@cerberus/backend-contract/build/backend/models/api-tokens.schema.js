"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.ApiTokensSchema = void 0;
const zod_1 = require("zod");
exports.ApiTokensSchema = zod_1.z.object({
    uuid: zod_1.z.string().uuid(),
    token: zod_1.z.string(),
    tokenName: zod_1.z.string(),
    createdAt: zod_1.z
        .string()
        .datetime()
        .transform((str) => new Date(str)),
    updatedAt: zod_1.z
        .string()
        .datetime()
        .transform((str) => new Date(str)),
});
