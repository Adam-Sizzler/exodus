"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.ExternalSquadSchema = void 0;
const zod_1 = require("zod");
const constants_1 = require("../constants");
const external_squads_1 = require("./external-squads");
const subscription_settings_1 = require("./subscription-settings");
exports.ExternalSquadSchema = zod_1.z.object({
    uuid: zod_1.z.uuid(),
    viewPosition: zod_1.z.int(),
    name: zod_1.z.string(),
    tags: zod_1.z.array(zod_1.z.string()).default([]),
    info: zod_1.z.object({
        membersCount: zod_1.z.number(),
    }),
    templates: zod_1.z.array(zod_1.z.object({
        templateUuid: zod_1.z.uuid(),
        templateType: zod_1.z.enum(constants_1.SUBSCRIPTION_TEMPLATE_TYPE),
    })),
    subscriptionSettings: zod_1.z.nullable(external_squads_1.ExternalSquadSubscriptionSettingsSchema),
    hostOverrides: zod_1.z.nullable(external_squads_1.ExternalSquadHostOverridesSchema),
    responseHeadersAdd: external_squads_1.ExternalSquadResponseHeadersAddSchema,
    responseHeadersRemove: external_squads_1.ExternalSquadResponseHeadersRemoveSchema,
    hwidSettings: zod_1.z.nullable(subscription_settings_1.HwidSettingsSchema),
    customRemarks: zod_1.z.nullable(subscription_settings_1.CustomRemarksSchema),
    subpageConfigUuid: zod_1.z.nullable(zod_1.z.uuid()),
    createdAt: zod_1.z.iso.datetime().transform((str) => new Date(str)),
    updatedAt: zod_1.z.iso.datetime().transform((str) => new Date(str)),
});
