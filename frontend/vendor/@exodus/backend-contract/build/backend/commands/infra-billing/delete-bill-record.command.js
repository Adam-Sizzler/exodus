"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.DeleteInfraBillingRecordCommand = void 0;
const zod_1 = require("zod");
const api_1 = require("../../api");
const constants_1 = require("../../constants");
var DeleteInfraBillingRecordCommand;
(function (DeleteInfraBillingRecordCommand) {
    DeleteInfraBillingRecordCommand.url = api_1.REST_API.INFRA_BILLING.DELETE_BILLING_HISTORY;
    DeleteInfraBillingRecordCommand.TSQ_url = DeleteInfraBillingRecordCommand.url(':uuid');
    DeleteInfraBillingRecordCommand.endpointDetails = (0, constants_1.getEndpointDetails)(api_1.INFRA_BILLING_ROUTES.DELETE_BILLING_HISTORY(':uuid'), 'delete', 'Delete infra billing history', { scope: 'delete-bill-record', kind: 'write' });
    DeleteInfraBillingRecordCommand.RequestParamSchema = zod_1.z.object({
        uuid: zod_1.z.uuid(),
    });

})(DeleteInfraBillingRecordCommand || (exports.DeleteInfraBillingRecordCommand = DeleteInfraBillingRecordCommand = {}));
