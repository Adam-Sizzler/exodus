"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.SubscriptionSettingsSchema = void 0;
const zod_1 = require("zod");
const response_rules_1 = require("./response-rules");
const custom_remarks_schema_1 = require("./subscription-settings/custom-remarks.schema");
const hwid_settings_schema_1 = require("./subscription-settings/hwid-settings.schema");
exports.SubscriptionSettingsSchema = zod_1.z.object({
    uuid: zod_1.z.uuid(),
    serveJsonAtBaseSubscription: zod_1.z.boolean(),
    isShowCustomRemarks: zod_1.z.boolean(),
    customRemarks: custom_remarks_schema_1.CustomRemarksSchema,
    customResponseHeaders: zod_1.z.nullable(zod_1.z.record(zod_1.z.string(), zod_1.z.string())),
    randomizeHosts: zod_1.z.boolean(),
    responseRules: zod_1.z.nullable(response_rules_1.ResponseRulesConfigSchema),
    hwidSettings: zod_1.z.nullable(hwid_settings_schema_1.HwidSettingsSchema),
    createdAt: zod_1.z.iso.datetime().transform((str) => new Date(str)),
    updatedAt: zod_1.z.iso.datetime().transform((str) => new Date(str)),
});
