"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.GetInfraBillingRecordsCommand = void 0;
const zod_1 = require("zod");
const api_1 = require("../../api");
const constants_1 = require("../../constants");
const models_1 = require("../../models");
var GetInfraBillingRecordsCommand;
(function (GetInfraBillingRecordsCommand) {
    GetInfraBillingRecordsCommand.url = api_1.REST_API.INFRA_BILLING.GET_BILLING_HISTORY;
    GetInfraBillingRecordsCommand.TSQ_url = GetInfraBillingRecordsCommand.url;
    GetInfraBillingRecordsCommand.endpointDetails = (0, constants_1.getEndpointDetails)(api_1.INFRA_BILLING_ROUTES.GET_BILLING_HISTORY, 'get', 'Get infra billing history', { scope: 'list-bill-records', kind: 'read' });
    GetInfraBillingRecordsCommand.RequestQuerySchema = zod_1.z.object({
        start: zod_1.z.coerce
            .number()
            .default(0)
            .describe('Start index (offset) of the billing history records to return, default is 0'),
        size: zod_1.z.coerce
            .number()
            .min(1)
            .max(500)
            .describe('Number of billing records to return, no more than 500')
            .default(50),
    });
    GetInfraBillingRecordsCommand.ResponseSchema = zod_1.z.object({
        response: zod_1.z.object({
            records: zod_1.z.array(models_1.InfraBillingHistoryRecordSchema),
            total: zod_1.z.number(),
        }),
    });
})(GetInfraBillingRecordsCommand || (exports.GetInfraBillingRecordsCommand = GetInfraBillingRecordsCommand = {}));
