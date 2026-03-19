"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.CreateInfraProviderCommand = void 0;
const zod_1 = require("zod");
const api_1 = require("../../api");
const constants_1 = require("../../constants");
const models_1 = require("../../models");
var CreateInfraProviderCommand;
(function (CreateInfraProviderCommand) {
    CreateInfraProviderCommand.url = api_1.REST_API.INFRA_BILLING.CREATE_PROVIDER;
    CreateInfraProviderCommand.TSQ_url = CreateInfraProviderCommand.url;
    CreateInfraProviderCommand.endpointDetails = (0, constants_1.getEndpointDetails)(api_1.INFRA_BILLING_ROUTES.CREATE_PROVIDER, 'post', 'Create infra provider');
    CreateInfraProviderCommand.RequestSchema = zod_1.z.object({
        name: zod_1.z
            .string()
            .min(2, 'Name must be at least 2 characters')
            .max(30, 'Name must be less than 30 characters'),
        faviconLink: zod_1.z.string().url().optional(),
        loginUrl: zod_1.z.string().url().optional(),
    });
    CreateInfraProviderCommand.ResponseSchema = zod_1.z.object({
        response: models_1.InfraProviderSchema,
    });
})(CreateInfraProviderCommand || (exports.CreateInfraProviderCommand = CreateInfraProviderCommand = {}));
