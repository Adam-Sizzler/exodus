"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.SnippetsSchema = void 0;
const zod_1 = require("zod");
exports.SnippetsSchema = zod_1.z.object({
    name: zod_1.z.string(),
    snippet: zod_1.z.unknown(),
});
