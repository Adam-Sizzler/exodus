"use strict";
var __importDefault = (this && this.__importDefault) || function (mod) {
    return (mod && mod.__esModule) ? mod : { "default": mod };
};
Object.defineProperty(exports, "__esModule", { value: true });
exports.BrandingSettingsSchema = void 0;
const zod_1 = __importDefault(require("zod"));
exports.BrandingSettingsSchema = zod_1.default.object({
    title: zod_1.default.nullable(zod_1.default.string()),
    logoUrl: zod_1.default.nullable(zod_1.default.string().url()),
});
