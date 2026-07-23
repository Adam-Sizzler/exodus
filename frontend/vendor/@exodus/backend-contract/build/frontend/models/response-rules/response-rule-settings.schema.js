"use strict";
var __importDefault = (this && this.__importDefault) || function (mod) {
    return (mod && mod.__esModule) ? mod : { "default": mod };
};
Object.defineProperty(exports, "__esModule", { value: true });
exports.ResponseRuleSettingsSchema = void 0;
const zod_1 = __importDefault(require("zod"));
exports.ResponseRuleSettingsSchema = zod_1.default
    .object({
    disableSubscriptionAccessByPath: zod_1.default
        .boolean()
        .optional()
        .describe(JSON.stringify({
        markdownDescription: "Usually, a user's subscription may also be available via additional paths such as **/json**, **/stash**, or **/mihomo**. If this flag is set to **true**, access via these additional paths will be disabled.",
    })),
})
    .optional()
    .describe(JSON.stringify({
    markdownDescription: 'Settings for the **response rules** config. Optional.',
}));
