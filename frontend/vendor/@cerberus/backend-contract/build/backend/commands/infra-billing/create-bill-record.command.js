"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.CreateInfraBillingHistoryRecordCommand = void 0;
const zod_1 = require("zod");
const models_1 = require("../../models");
const api_1 = require("../../api");
const constants_1 = require("../../constants");
var CreateInfraBillingHistoryRecordCommand;
(function (CreateInfraBillingHistoryRecordCommand) {
    CreateInfraBillingHistoryRecordCommand.url = api_1.REST_API.INFRA_BILLING.CREATE_BILLING_HISTORY;
    CreateInfraBillingHistoryRecordCommand.TSQ_url = CreateInfraBillingHistoryRecordCommand.url;
    CreateInfraBillingHistoryRecordCommand.endpointDetails = (0, constants_1.getEndpointDetails)(api_1.INFRA_BILLING_ROUTES.CREATE_BILLING_HISTORY, 'post', 'Create infra billing history');
    CreateInfraBillingHistoryRecordCommand.RequestSchema = zod_1.z.object({
        providerUuid: zod_1.z.string().uuid(),
        amount: zod_1.z.number().min(0, 'Amount must be greater than 0'),
        billedAt: zod_1.z
            .string({
            invalid_type_error: 'Invalid date format',
        })
            .datetime({ message: 'Invalid date format', offset: true, local: true })
            .transform((str) => new Date(str))
            .describe('Billing date. Format: 2025-01-17T15:38:45.065Z'),
    });
    CreateInfraBillingHistoryRecordCommand.ResponseSchema = zod_1.z.object({
        response: zod_1.z.object({
            records: zod_1.z.array(models_1.InfraBillingHistoryRecordSchema),
            total: zod_1.z.number(),
        }),
    });
})(CreateInfraBillingHistoryRecordCommand || (exports.CreateInfraBillingHistoryRecordCommand = CreateInfraBillingHistoryRecordCommand = {}));
