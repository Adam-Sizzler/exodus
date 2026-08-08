"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.GetInfraProviderCommand = void 0;
const zod_1 = require("zod");
const api_1 = require("../../api");
const constants_1 = require("../../constants");
const models_1 = require("../../models");
var GetInfraProviderCommand;
(function (GetInfraProviderCommand) {
    GetInfraProviderCommand.url = api_1.REST_API.INFRA_BILLING.GET_PROVIDER_BY_UUID;
    GetInfraProviderCommand.TSQ_url = GetInfraProviderCommand.url(':uuid');
    GetInfraProviderCommand.endpointDetails = (0, constants_1.getEndpointDetails)(api_1.INFRA_BILLING_ROUTES.GET_PROVIDER_BY_UUID(':uuid'), 'get', 'Get infra provider by uuid', { scope: 'get-provider', kind: 'read' });
    GetInfraProviderCommand.RequestParamSchema = zod_1.z.object({
        uuid: zod_1.z.uuid(),
    });
    GetInfraProviderCommand.ResponseSchema = zod_1.z.object({
        response: models_1.InfraProviderSchema,
    });

})(GetInfraProviderCommand || (exports.GetInfraProviderCommand = GetInfraProviderCommand = {}));
