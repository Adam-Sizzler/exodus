"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.DeleteInfraProviderByUuidCommand = void 0;
const zod_1 = require("zod");
const api_1 = require("../../api");
const constants_1 = require("../../constants");
var DeleteInfraProviderByUuidCommand;
(function (DeleteInfraProviderByUuidCommand) {
    DeleteInfraProviderByUuidCommand.url = api_1.REST_API.INFRA_BILLING.DELETE_PROVIDER;
    DeleteInfraProviderByUuidCommand.TSQ_url = DeleteInfraProviderByUuidCommand.url(':uuid');
    DeleteInfraProviderByUuidCommand.endpointDetails = (0, constants_1.getEndpointDetails)(api_1.INFRA_BILLING_ROUTES.DELETE_PROVIDER(':uuid'), 'delete', 'Delete infra provider by uuid');
    DeleteInfraProviderByUuidCommand.RequestSchema = zod_1.z.object({
        uuid: zod_1.z.string().uuid(),
    });
    DeleteInfraProviderByUuidCommand.ResponseSchema = zod_1.z.object({
        response: zod_1.z.object({
            isDeleted: zod_1.z.boolean(),
        }),
    });
})(DeleteInfraProviderByUuidCommand || (exports.DeleteInfraProviderByUuidCommand = DeleteInfraProviderByUuidCommand = {}));
