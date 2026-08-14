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
    UpdateInfraProviderCommand.endpointDetails = (0, constants_1.getEndpointDetails)(api_1.INFRA_BILLING_ROUTES.UPDATE_PROVIDER, 'patch', 'Update infra provider', { scope: 'update-provider', kind: 'write' });
    UpdateInfraProviderCommand.RequestBodySchema = zod_1.z.object({
        uuid: zod_1.z.uuid(),
        name: zod_1.z.string().min(2).max(30).optional(),
        faviconLink: zod_1.z.url().nullish(),
        loginUrl: zod_1.z.url().nullish(),
    });
    UpdateInfraProviderCommand.ResponseSchema = zod_1.z.object({
        response: models_1.InfraProviderSchema,
    });
})(UpdateInfraProviderCommand || (exports.UpdateInfraProviderCommand = UpdateInfraProviderCommand = {}));
