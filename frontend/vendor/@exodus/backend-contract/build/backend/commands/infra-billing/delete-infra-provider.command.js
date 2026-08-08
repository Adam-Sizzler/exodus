"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.DeleteInfraProviderCommand = void 0;
const zod_1 = require("zod");
const api_1 = require("../../api");
const constants_1 = require("../../constants");
var DeleteInfraProviderCommand;
(function (DeleteInfraProviderCommand) {
    DeleteInfraProviderCommand.url = api_1.REST_API.INFRA_BILLING.DELETE_PROVIDER;
    DeleteInfraProviderCommand.TSQ_url = DeleteInfraProviderCommand.url(':uuid');
    DeleteInfraProviderCommand.endpointDetails = (0, constants_1.getEndpointDetails)(api_1.INFRA_BILLING_ROUTES.DELETE_PROVIDER(':uuid'), 'delete', 'Delete infra provider by uuid', { scope: 'delete-provider', kind: 'write' });
    DeleteInfraProviderCommand.RequestParamSchema = zod_1.z.object({
        uuid: zod_1.z.uuid(),
    });

})(DeleteInfraProviderCommand || (exports.DeleteInfraProviderCommand = DeleteInfraProviderCommand = {}));
