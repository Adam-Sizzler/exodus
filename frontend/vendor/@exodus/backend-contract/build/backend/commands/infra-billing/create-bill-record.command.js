"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.CreateInfraBillingRecordCommand = void 0;
const zod_1 = require("zod");
const api_1 = require("../../api");
const constants_1 = require("../../constants");
const models_1 = require("../../models");
var CreateInfraBillingRecordCommand;
(function (CreateInfraBillingRecordCommand) {
    CreateInfraBillingRecordCommand.url = api_1.REST_API.INFRA_BILLING.CREATE_BILLING_HISTORY;
    CreateInfraBillingRecordCommand.TSQ_url = CreateInfraBillingRecordCommand.url;
    CreateInfraBillingRecordCommand.endpointDetails = (0, constants_1.getEndpointDetails)(api_1.INFRA_BILLING_ROUTES.CREATE_BILLING_HISTORY, 'post', 'Create infra billing history', { scope: 'create-bill-record', kind: 'write' });
    CreateInfraBillingRecordCommand.RequestBodySchema = zod_1.z.object({
        providerUuid: zod_1.z.uuid(),
        amount: zod_1.z.number().min(0),
        billedAt: zod_1.z.iso
            .datetime({ offset: true, local: true })
            .transform((str) => new Date(str))
            .describe('Billing date. Format: 2025-01-17T15:38:45.065Z'),
    });
    CreateInfraBillingRecordCommand.ResponseSchema = zod_1.z.object({
        response: zod_1.z.object({
            records: zod_1.z.array(models_1.InfraBillingHistoryRecordSchema),
            total: zod_1.z.number(),
        }),
    });
})(CreateInfraBillingRecordCommand || (exports.CreateInfraBillingRecordCommand = CreateInfraBillingRecordCommand = {}));
