"use strict";
var __importDefault = (this && this.__importDefault) || function (mod) {
    return (mod && mod.__esModule) ? mod : { "default": mod };
};
Object.defineProperty(exports, "__esModule", { value: true });
exports.SharedListPreviewSchema = exports.SharedListsSchema = exports.SharedListNameSchema = void 0;
const zod_1 = __importDefault(require("zod"));
exports.SharedListNameSchema = zod_1.default
    .string()
    .min(2, 'Name must be at least 2 characters')
    .max(255, 'Name must be less than 255 characters')
    .regex(/^[A-Za-z0-9_-]+$/, 'Name can only contain letters, numbers, underscores and dashes. The "ext:" prefix is added automatically');
exports.SharedListsSchema = zod_1.default.object({
    name: exports.SharedListNameSchema,
    config: zod_1.default.record(zod_1.default.string(), zod_1.default.unknown()),
});
exports.SharedListPreviewSchema = zod_1.default.object({
    name: exports.SharedListNameSchema,
    type: zod_1.default.string(),
    itemsCount: zod_1.default.number(),
});
