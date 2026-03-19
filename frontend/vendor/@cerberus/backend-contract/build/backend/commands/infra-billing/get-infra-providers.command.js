"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.GetInfraProvidersCommand = void 0;
const zod_1 = require("zod");
const api_1 = require("../../api");
const constants_1 = require("../../constants");
const models_1 = require("../../models");
var GetInfraProvidersCommand;
(function (GetInfraProvidersCommand) {
    GetInfraProvidersCommand.url = api_1.REST_API.INFRA_BILLING.GET_PROVIDERS;
    GetInfraProvidersCommand.TSQ_url = GetInfraProvidersCommand.url;
    GetInfraProvidersCommand.endpointDetails = (0, constants_1.getEndpointDetails)(api_1.INFRA_BILLING_ROUTES.GET_PROVIDERS, 'get', 'Get all infra providers');
    GetInfraProvidersCommand.ResponseSchema = zod_1.z.object({
        response: zod_1.z.object({
            total: zod_1.z.number(),
            providers: zod_1.z.array(models_1.InfraProviderSchema),
        }),
    });
})(GetInfraProvidersCommand || (exports.GetInfraProvidersCommand = GetInfraProvidersCommand = {}));
