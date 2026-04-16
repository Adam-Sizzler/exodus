"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.DeleteInfraBillingHistoryRecordCommand = void 0;
const zod_1 = require("zod");
const models_1 = require("../../models");
const api_1 = require("../../api");
const constants_1 = require("../../constants");
var DeleteInfraBillingHistoryRecordCommand;
(function (DeleteInfraBillingHistoryRecordCommand) {
    DeleteInfraBillingHistoryRecordCommand.url = api_1.REST_API.INFRA_BILLING.DELETE_BILLING_HISTORY;
    DeleteInfraBillingHistoryRecordCommand.TSQ_url = DeleteInfraBillingHistoryRecordCommand.url(':uuid');
    DeleteInfraBillingHistoryRecordCommand.endpointDetails = (0, constants_1.getEndpointDetails)(api_1.INFRA_BILLING_ROUTES.DELETE_BILLING_HISTORY(':uuid'), 'delete', 'Delete infra billing history');
    DeleteInfraBillingHistoryRecordCommand.RequestSchema = zod_1.z.object({
        uuid: zod_1.z.string().uuid(),
    });
    DeleteInfraBillingHistoryRecordCommand.ResponseSchema = zod_1.z.object({
        response: zod_1.z.object({
            records: zod_1.z.array(models_1.InfraBillingHistoryRecordSchema),
            total: zod_1.z.number(),
        }),
    });
})(DeleteInfraBillingHistoryRecordCommand || (exports.DeleteInfraBillingHistoryRecordCommand = DeleteInfraBillingHistoryRecordCommand = {}));
