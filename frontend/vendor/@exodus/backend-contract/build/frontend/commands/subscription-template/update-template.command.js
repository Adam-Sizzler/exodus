"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.UpdateSubscriptionTemplateCommand = void 0;
const zod_1 = require("zod");
const api_1 = require("../../api");
const constants_1 = require("../../constants");
const models_1 = require("../../models");
var UpdateSubscriptionTemplateCommand;
(function (UpdateSubscriptionTemplateCommand) {
    UpdateSubscriptionTemplateCommand.url = api_1.REST_API.SUBSCRIPTION_TEMPLATE.UPDATE;
    UpdateSubscriptionTemplateCommand.TSQ_url = UpdateSubscriptionTemplateCommand.url;
    UpdateSubscriptionTemplateCommand.endpointDetails = (0, constants_1.getEndpointDetails)(api_1.SUBSCRIPTION_TEMPLATE_ROUTES.UPDATE, 'patch', 'Update subscription template', { scope: 'update', kind: 'write' });
    UpdateSubscriptionTemplateCommand.RequestBodySchema = zod_1.z.object({
        uuid: zod_1.z.uuid(),
        name: zod_1.z
            .string()
            .min(2)
            .max(255)
            .regex(/^[A-Za-z0-9_\s-]+$/, 'Name can only contain letters, numbers, underscores, dashes and spaces')
            .optional(),
        templateJson: zod_1.z.optional(zod_1.z.looseObject({})),
        encodedTemplateYaml: zod_1.z.optional(zod_1.z.string()),
    });
    UpdateSubscriptionTemplateCommand.ResponseSchema = zod_1.z.object({
        response: models_1.SubscriptionTemplateSchema,
    });

})(UpdateSubscriptionTemplateCommand || (exports.UpdateSubscriptionTemplateCommand = UpdateSubscriptionTemplateCommand = {}));
