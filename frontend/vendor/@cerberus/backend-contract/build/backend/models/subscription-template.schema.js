"use strict";
var __importDefault = (this && this.__importDefault) || function (mod) {
    return (mod && mod.__esModule) ? mod : { "default": mod };
};
Object.defineProperty(exports, "__esModule", { value: true });
exports.SubscriptionTemplateSchema = void 0;
const zod_1 = __importDefault(require("zod"));
const constants_1 = require("../constants");
exports.SubscriptionTemplateSchema = zod_1.default.object({
    uuid: zod_1.default.string().uuid(),
    viewPosition: zod_1.default.number().int(),
    name: zod_1.default.string(),
    templateType: zod_1.default.nativeEnum(constants_1.SUBSCRIPTION_TEMPLATE_TYPE),
    templateJson: zod_1.default.nullable(zod_1.default.unknown()),
    encodedTemplateYaml: zod_1.default.nullable(zod_1.default.string()),
});
