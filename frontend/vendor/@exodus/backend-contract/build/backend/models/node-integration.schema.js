"use strict";
var __importDefault = (this && this.__importDefault) || function (mod) {
    return (mod && mod.__esModule) ? mod : { "default": mod };
};
Object.defineProperty(exports, "__esModule", { value: true });
exports.NodeIntegrationSchema = void 0;
const zod_1 = __importDefault(require("zod"));
exports.NodeIntegrationSchema = zod_1.default.object({
    uuid: zod_1.default.uuid(),
    name: zod_1.default.string(),
    description: zod_1.default.nullable(zod_1.default.string()),
    config: zod_1.default.record(zod_1.default.string(), zod_1.default.unknown()),
});
