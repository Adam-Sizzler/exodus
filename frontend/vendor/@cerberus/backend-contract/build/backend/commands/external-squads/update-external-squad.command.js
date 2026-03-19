"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.UpdateExternalSquadCommand = void 0;
const zod_1 = require("zod");
const models_1 = require("../../models");
const constants_1 = require("../../constants");
const api_1 = require("../../api");
var UpdateExternalSquadCommand;
(function (UpdateExternalSquadCommand) {
    UpdateExternalSquadCommand.url = api_1.REST_API.EXTERNAL_SQUADS.UPDATE;
    UpdateExternalSquadCommand.TSQ_url = UpdateExternalSquadCommand.url;
    UpdateExternalSquadCommand.endpointDetails = (0, constants_1.getEndpointDetails)(api_1.EXTERNAL_SQUADS_ROUTES.UPDATE, 'patch', 'Update external squad');
    UpdateExternalSquadCommand.RequestSchema = zod_1.z.object({
        uuid: zod_1.z.string().uuid(),
        name: zod_1.z
            .string()
            .min(2, 'Name must be at least 2 characters')
            .max(30, 'Name must be less than 30 characters')
            .regex(/^[A-Za-z0-9_\s-]+$/, 'Name can only contain letters, numbers, underscores, dashes and spaces')
            .optional(),
        templates: zod_1.z
            .array(zod_1.z.object({
            templateUuid: zod_1.z.string().uuid(),
            templateType: zod_1.z.nativeEnum(constants_1.SUBSCRIPTION_TEMPLATE_TYPE),
        }))
            .optional(),
        subscriptionSettings: models_1.ExternalSquadSubscriptionSettingsSchema.optional(),
        hostOverrides: models_1.ExternalSquadHostOverridesSchema.optional(),
        responseHeaders: models_1.ExternalSquadResponseHeadersSchema.optional(),
        hwidSettings: zod_1.z.optional(zod_1.z.nullable(models_1.HwidSettingsSchema)),
        customRemarks: zod_1.z.optional(zod_1.z.nullable(models_1.CustomRemarksSchema)),
        subpageConfigUuid: zod_1.z.optional(zod_1.z.nullable(zod_1.z.string().uuid())),
    });
    UpdateExternalSquadCommand.ResponseSchema = zod_1.z.object({
        response: models_1.ExternalSquadSchema,
    });
})(UpdateExternalSquadCommand || (exports.UpdateExternalSquadCommand = UpdateExternalSquadCommand = {}));
