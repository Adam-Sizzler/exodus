"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.DeleteInfraBillingNodeCommand = void 0;
const zod_1 = require("zod");
const api_1 = require("../../api");
const constants_1 = require("../../constants");
var DeleteInfraBillingNodeCommand;
(function (DeleteInfraBillingNodeCommand) {
    DeleteInfraBillingNodeCommand.url = api_1.REST_API.INFRA_BILLING.DELETE_BILLING_NODE;
    DeleteInfraBillingNodeCommand.TSQ_url = DeleteInfraBillingNodeCommand.url(':uuid');
    DeleteInfraBillingNodeCommand.endpointDetails = (0, constants_1.getEndpointDetails)(api_1.INFRA_BILLING_ROUTES.DELETE_BILLING_NODE(':uuid'), 'delete', 'Delete infra billing node', { scope: 'delete-billing-node', kind: 'write' });
    DeleteInfraBillingNodeCommand.RequestParamSchema = zod_1.z.object({
        uuid: zod_1.z.uuid(),
    });

})(DeleteInfraBillingNodeCommand || (exports.DeleteInfraBillingNodeCommand = DeleteInfraBillingNodeCommand = {}));
