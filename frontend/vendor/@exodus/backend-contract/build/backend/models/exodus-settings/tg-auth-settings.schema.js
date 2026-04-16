"use strict";
var __importDefault = (this && this.__importDefault) || function (mod) {
    return (mod && mod.__esModule) ? mod : { "default": mod };
};
Object.defineProperty(exports, "__esModule", { value: true });
exports.TgAuthSettingsSchema = void 0;
const zod_1 = __importDefault(require("zod"));
exports.TgAuthSettingsSchema = zod_1.default.object({
    enabled: zod_1.default.boolean(),
    botToken: zod_1.default.nullable(zod_1.default.string()),
    adminIds: zod_1.default.array(zod_1.default.string()),
});
