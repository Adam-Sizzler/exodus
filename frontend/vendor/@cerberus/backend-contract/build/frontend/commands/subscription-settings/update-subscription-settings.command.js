"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.UpdateSubscriptionSettingsCommand = void 0;
const zod_1 = require("zod");
const models_1 = require("../../models");
const api_1 = require("../../api");
const constants_1 = require("../../constants");
var UpdateSubscriptionSettingsCommand;
(function (UpdateSubscriptionSettingsCommand) {
    UpdateSubscriptionSettingsCommand.url = api_1.REST_API.SUBSCRIPTION_SETTINGS.UPDATE;
    UpdateSubscriptionSettingsCommand.TSQ_url = UpdateSubscriptionSettingsCommand.url;
    UpdateSubscriptionSettingsCommand.endpointDetails = (0, constants_1.getEndpointDetails)(api_1.SUBSCRIPTION_SETTINGS_ROUTES.UPDATE, 'patch', 'Update subscription settings');
    UpdateSubscriptionSettingsCommand.RequestSchema = zod_1.z.object({
        uuid: zod_1.z.string().uuid(),
        profileTitle: zod_1.z.optional(zod_1.z.string()),
        supportLink: zod_1.z.optional(zod_1.z.string()),
        profileUpdateInterval: zod_1.z.optional(zod_1.z.number().int()),
        isProfileWebpageUrlEnabled: zod_1.z.optional(zod_1.z.boolean()),
        serveJsonAtBaseSubscription: zod_1.z.optional(zod_1.z.boolean()),
        happAnnounce: zod_1.z.optional(zod_1.z
            .string()
            .max(200, { message: 'Announce must be less than 200 characters' })
            .nullable()),
        happRouting: zod_1.z.optional(zod_1.z.string().nullable()),
        isShowCustomRemarks: zod_1.z.optional(zod_1.z.boolean()),
        customRemarks: zod_1.z.optional(models_1.CustomRemarksSchema),
        customResponseHeaders: zod_1.z.optional(zod_1.z.record(zod_1.z
            .string()
            .regex(/^[a-zA-Z0-9_-]+$/, 'Invalid header name. Only letters(a-z, A-Z), numbers(0-9), underscores(_) and hyphens(-) are allowed.'), zod_1.z.string())),
        randomizeHosts: zod_1.z.optional(zod_1.z.boolean()),
        responseRules: zod_1.z.optional(models_1.ResponseRulesConfigSchema),
        hwidSettings: zod_1.z.optional(models_1.HwidSettingsSchema),
    });
    UpdateSubscriptionSettingsCommand.ResponseSchema = zod_1.z.object({
        response: models_1.SubscriptionSettingsSchema,
    });
})(UpdateSubscriptionSettingsCommand || (exports.UpdateSubscriptionSettingsCommand = UpdateSubscriptionSettingsCommand = {}));
