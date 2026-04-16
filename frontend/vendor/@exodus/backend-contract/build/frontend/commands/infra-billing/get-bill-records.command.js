"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.GetInfraBillingHistoryRecordsCommand = void 0;
const zod_1 = require("zod");
const models_1 = require("../../models");
const api_1 = require("../../api");
const constants_1 = require("../../constants");
var GetInfraBillingHistoryRecordsCommand;
(function (GetInfraBillingHistoryRecordsCommand) {
    GetInfraBillingHistoryRecordsCommand.url = api_1.REST_API.INFRA_BILLING.GET_BILLING_HISTORY;
    GetInfraBillingHistoryRecordsCommand.TSQ_url = GetInfraBillingHistoryRecordsCommand.url;
    GetInfraBillingHistoryRecordsCommand.endpointDetails = (0, constants_1.getEndpointDetails)(api_1.INFRA_BILLING_ROUTES.GET_BILLING_HISTORY, 'get', 'Get infra billing history');
    GetInfraBillingHistoryRecordsCommand.RequestQuerySchema = zod_1.z.object({
        start: zod_1.z.coerce
            .number()
            .default(0)
            .describe('Start index (offset) of the billing history records to return, default is 0'),
        size: zod_1.z.coerce
            .number()
            .min(1, 'Size (limit) must be greater than 0')
            .max(500, 'Size (limit) must be less than 500')
            .describe('Number of billing records to return, no more than 500')
            .default(50),
    });
    GetInfraBillingHistoryRecordsCommand.ResponseSchema = zod_1.z.object({
        response: zod_1.z.object({
            records: zod_1.z.array(models_1.InfraBillingHistoryRecordSchema),
            total: zod_1.z.number(),
        }),
    });
})(GetInfraBillingHistoryRecordsCommand || (exports.GetInfraBillingHistoryRecordsCommand = GetInfraBillingHistoryRecordsCommand = {}));
