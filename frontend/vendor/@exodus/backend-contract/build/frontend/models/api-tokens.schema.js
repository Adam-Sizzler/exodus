"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.ApiTokensSchema = void 0;
const zod_1 = require("zod");
exports.ApiTokensSchema = zod_1.z.object({
    uuid: zod_1.z.uuid(),
    name: zod_1.z.string(),
    expireAt: zod_1.z.iso.datetime()
        .transform((str) => new Date(str)),
    scopes: zod_1.z.array(zod_1.z.string()),
    createdAt: zod_1.z.iso.datetime()
        .transform((str) => new Date(str)),
    updatedAt: zod_1.z.iso.datetime()
        .transform((str) => new Date(str)),
});
