"use strict";
var __importDefault = (this && this.__importDefault) || function (mod) {
    return (mod && mod.__esModule) ? mod : { "default": mod };
};
Object.defineProperty(exports, "__esModule", { value: true });
exports.NodePluginSchema = void 0;
const zod_1 = require("zod");
exports.NodePluginSchema = (zod_1.z || zod_1.default.z || zod_1.default).object({
    uuid: (zod_1.z || zod_1.default.z || zod_1.default).uuid(),
    viewPosition: (zod_1.z || zod_1.default.z || zod_1.default).number().int(),
    name: (zod_1.z || zod_1.default.z || zod_1.default).string(),
    pluginConfig: (zod_1.z || zod_1.default.z || zod_1.default).nullable((zod_1.z || zod_1.default.z || zod_1.default).unknown()),
});
