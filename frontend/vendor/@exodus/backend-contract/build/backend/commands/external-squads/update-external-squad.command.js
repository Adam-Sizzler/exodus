"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.UpdateExternalSquadCommand = void 0;
const zod_1 = require("zod");
const api_1 = require("../../api");
const constants_1 = require("../../constants");
const models_1 = require("../../models");
var UpdateExternalSquadCommand;
(function (UpdateExternalSquadCommand) {
    UpdateExternalSquadCommand.url = api_1.REST_API.EXTERNAL_SQUADS.UPDATE;
    UpdateExternalSquadCommand.TSQ_url = UpdateExternalSquadCommand.url;
    UpdateExternalSquadCommand.endpointDetails = (0, constants_1.getEndpointDetails)(api_1.EXTERNAL_SQUADS_ROUTES.UPDATE, 'patch', 'Update external squad', { scope: 'update', kind: 'write' });
    UpdateExternalSquadCommand.RequestBodySchema = zod_1.z.object({
        uuid: zod_1.z.uuid().describe('UUID of the external squad'),
        name: zod_1.z
            .string()
            .min(2)
            .max(30)
            .regex(/^[A-Za-z0-9_\s-]+$/, 'Name can only contain letters, numbers, underscores, dashes and spaces')
            .optional(),
        templates: zod_1.z
            .array(zod_1.z.object({
            templateUuid: zod_1.z.uuid().describe('UUID of the subscription template'),
            templateType: zod_1.z
                .enum(constants_1.SUBSCRIPTION_TEMPLATE_TYPE)
                .describe('Type of the subscription template'),
        }))
            .optional(),
        subscriptionSettings: models_1.ExternalSquadSubscriptionSettingsSchema.optional(),
        hostOverrides: models_1.ExternalSquadHostOverridesSchema.optional(),
        responseHeadersAdd: models_1.ExternalSquadResponseHeadersAddSchema.optional(),
        responseHeadersRemove: models_1.ExternalSquadResponseHeadersRemoveSchema.optional(),
        hwidSettings: models_1.HwidSettingsSchema.nullish(),
        customRemarks: models_1.CustomRemarksSchema.nullish(),
        subpageConfigUuid: zod_1.z.uuid().nullish(),
    });
    UpdateExternalSquadCommand.ResponseSchema = zod_1.z.object({
        response: models_1.ExternalSquadSchema,
    });

})(UpdateExternalSquadCommand || (exports.UpdateExternalSquadCommand = UpdateExternalSquadCommand = {}));
