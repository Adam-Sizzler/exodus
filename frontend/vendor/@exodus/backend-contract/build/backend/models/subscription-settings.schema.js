"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.SubscriptionSettingsSchema = void 0;
const zod_1 = require("zod");
const custom_remarks_schema_1 = require("./subscription-settings/custom-remarks.schema");
const hwid_settings_schema_1 = require("./subscription-settings/hwid-settings.schema");
const response_rules_1 = require("./response-rules");
exports.SubscriptionSettingsSchema = zod_1.z.object({
    uuid: zod_1.z.string().uuid(),
    profileTitle: zod_1.z.string(),
    supportLink: zod_1.z.string(),
    profileUpdateInterval: zod_1.z
        .number()
        .int()
        .min(1, 'Profile update interval must be greater than 0'),
    isProfileWebpageUrlEnabled: zod_1.z.boolean(),
    serveJsonAtBaseSubscription: zod_1.z.boolean(),
    isShowCustomRemarks: zod_1.z.boolean(),
    customRemarks: custom_remarks_schema_1.CustomRemarksSchema,
    happAnnounce: zod_1.z.string().nullable(),
    happRouting: zod_1.z.string().nullable(),
    customResponseHeaders: zod_1.z.nullable(zod_1.z.record(zod_1.z.string(), zod_1.z.string())),
    randomizeHosts: zod_1.z.boolean(),
    responseRules: zod_1.z.nullable(response_rules_1.ResponseRulesConfigSchema),
    hwidSettings: zod_1.z.nullable(hwid_settings_schema_1.HwidSettingsSchema),
    createdAt: zod_1.z
        .string()
        .datetime()
        .transform((str) => new Date(str)),
    updatedAt: zod_1.z
        .string()
        .datetime()
        .transform((str) => new Date(str)),
});
