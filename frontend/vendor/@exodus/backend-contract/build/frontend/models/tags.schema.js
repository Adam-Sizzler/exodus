"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.TagsSchema = exports.TagSchema = void 0;
const zod_1 = require("zod");
exports.TagSchema = zod_1.z
    .string()
    .regex(/^[A-Z0-9_:]+$/, 'Tag can only contain uppercase letters, numbers, underscores and colons')
    .max(36, 'Each tag must be less than 36 characters');
exports.TagsSchema = zod_1.z.array(exports.TagSchema).max(10, 'Maximum 10 tags');
