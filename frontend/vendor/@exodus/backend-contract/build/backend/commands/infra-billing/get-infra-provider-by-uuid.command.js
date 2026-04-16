"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.GetInfraProviderByUuidCommand = void 0;
const zod_1 = require("zod");
const api_1 = require("../../api");
const constants_1 = require("../../constants");
const models_1 = require("../../models");
var GetInfraProviderByUuidCommand;
(function (GetInfraProviderByUuidCommand) {
    GetInfraProviderByUuidCommand.url = api_1.REST_API.INFRA_BILLING.GET_PROVIDER_BY_UUID;
    GetInfraProviderByUuidCommand.TSQ_url = GetInfraProviderByUuidCommand.url(':uuid');
    GetInfraProviderByUuidCommand.endpointDetails = (0, constants_1.getEndpointDetails)(api_1.INFRA_BILLING_ROUTES.GET_PROVIDER_BY_UUID(':uuid'), 'get', 'Get infra provider by uuid');
    GetInfraProviderByUuidCommand.RequestSchema = zod_1.z.object({
        uuid: zod_1.z.string().uuid(),
    });
    GetInfraProviderByUuidCommand.ResponseSchema = zod_1.z.object({
        response: models_1.InfraProviderSchema,
    });
})(GetInfraProviderByUuidCommand || (exports.GetInfraProviderByUuidCommand = GetInfraProviderByUuidCommand = {}));
