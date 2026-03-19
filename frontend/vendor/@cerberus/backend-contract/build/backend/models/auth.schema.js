"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.AuthSchema = void 0;
const zod_1 = require("zod");
exports.AuthSchema = zod_1.z.object({
    hash: zod_1.z.string(),
});
