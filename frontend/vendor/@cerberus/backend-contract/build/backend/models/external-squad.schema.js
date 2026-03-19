"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.ExternalSquadSchema = void 0;
const zod_1 = require("zod");
const external_squads_1 = require("./external-squads");
const subscription_settings_1 = require("./subscription-settings");
const constants_1 = require("../constants");
exports.ExternalSquadSchema = zod_1.z.object({
    uuid: zod_1.z.string().uuid(),
    viewPosition: zod_1.z.number().int(),
    name: zod_1.z.string(),
    info: zod_1.z.object({
        membersCount: zod_1.z.number(),
    }),
    templates: zod_1.z.array(zod_1.z.object({
        templateUuid: zod_1.z.string().uuid(),
        templateType: zod_1.z.nativeEnum(constants_1.SUBSCRIPTION_TEMPLATE_TYPE),
    })),
    subscriptionSettings: zod_1.z.nullable(external_squads_1.ExternalSquadSubscriptionSettingsSchema),
    hostOverrides: zod_1.z.nullable(external_squads_1.ExternalSquadHostOverridesSchema),
    responseHeaders: external_squads_1.ExternalSquadResponseHeadersSchema,
    hwidSettings: zod_1.z.nullable(subscription_settings_1.HwidSettingsSchema),
    customRemarks: zod_1.z.nullable(subscription_settings_1.CustomRemarksSchema),
    subpageConfigUuid: zod_1.z.nullable(zod_1.z.string().uuid()),
    createdAt: zod_1.z
        .string()
        .datetime()
        .transform((str) => new Date(str)),
    updatedAt: zod_1.z
        .string()
        .datetime()
        .transform((str) => new Date(str)),
});
