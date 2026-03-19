"use strict";
var __importDefault = (this && this.__importDefault) || function (mod) {
    return (mod && mod.__esModule) ? mod : { "default": mod };
};
Object.defineProperty(exports, "__esModule", { value: true });
exports.HwidSettingsSchema = void 0;
const zod_1 = __importDefault(require("zod"));
exports.HwidSettingsSchema = zod_1.default.object({
    enabled: zod_1.default.boolean(),
    fallbackDeviceLimit: zod_1.default.number(),
    maxDevicesAnnounce: zod_1.default.nullable(zod_1.default.string().max(200, { message: 'Announce must be less than 200 characters' })),
});
