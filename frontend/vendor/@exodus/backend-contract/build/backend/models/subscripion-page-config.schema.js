"use strict";
var __importDefault = (this && this.__importDefault) || function (mod) {
    return (mod && mod.__esModule) ? mod : { "default": mod };
};
Object.defineProperty(exports, "__esModule", { value: true });
exports.SubscriptionPageConfigSchema = void 0;
const zod_1 = __importDefault(require("zod"));
exports.SubscriptionPageConfigSchema = zod_1.default.object({
    uuid: zod_1.default.string().uuid(),
    viewPosition: zod_1.default.number().int(),
    name: zod_1.default.string(),
    config: zod_1.default.nullable(zod_1.default.unknown()),
});
