"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.UpdateInfraProviderCommand = void 0;
const zod_1 = require("zod");
const api_1 = require("../../api");
const constants_1 = require("../../constants");
const models_1 = require("../../models");
var UpdateInfraProviderCommand;
(function (UpdateInfraProviderCommand) {
    UpdateInfraProviderCommand.url = api_1.REST_API.INFRA_BILLING.UPDATE_PROVIDER;
    UpdateInfraProviderCommand.TSQ_url = UpdateInfraProviderCommand.url;
    UpdateInfraProviderCommand.endpointDetails = (0, constants_1.getEndpointDetails)(api_1.INFRA_BILLING_ROUTES.UPDATE_PROVIDER, 'patch', 'Update infra provider');
    UpdateInfraProviderCommand.RequestSchema = zod_1.z.object({
        uuid: zod_1.z.string().uuid(),
        name: zod_1.z
            .string()
            .min(2, 'Name must be at least 2 characters')
            .max(30, 'Name must be less than 30 characters')
            .optional(),
        faviconLink: zod_1.z.optional(zod_1.z.nullable(zod_1.z.string().url())),
        loginUrl: zod_1.z.optional(zod_1.z.nullable(zod_1.z.string().url())),
    });
    UpdateInfraProviderCommand.ResponseSchema = zod_1.z.object({
        response: models_1.InfraProviderSchema,
    });
})(UpdateInfraProviderCommand || (exports.UpdateInfraProviderCommand = UpdateInfraProviderCommand = {}));
