"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.CreateInfraBillingNodeCommand = void 0;
const zod_1 = require("zod");
const models_1 = require("../../models");
const api_1 = require("../../api");
const constants_1 = require("../../constants");
var CreateInfraBillingNodeCommand;
(function (CreateInfraBillingNodeCommand) {
    CreateInfraBillingNodeCommand.url = api_1.REST_API.INFRA_BILLING.CREATE_BILLING_NODE;
    CreateInfraBillingNodeCommand.TSQ_url = CreateInfraBillingNodeCommand.url;
    CreateInfraBillingNodeCommand.endpointDetails = (0, constants_1.getEndpointDetails)(api_1.INFRA_BILLING_ROUTES.CREATE_BILLING_NODE, 'post', 'Create infra billing node');
    CreateInfraBillingNodeCommand.RequestSchema = zod_1.z.object({
        providerUuid: zod_1.z.string().uuid(),
        nodeUuid: zod_1.z.string().uuid(),
        nextBillingAt: zod_1.z
            .string({
            invalid_type_error: 'Invalid date format',
        })
            .datetime({ message: 'Invalid date format', offset: true, local: true })
            .transform((str) => new Date(str))
            .optional()
            .describe('Next billing date. Format: 2025-01-17T15:38:45.065Z'),
    });
    CreateInfraBillingNodeCommand.ResponseSchema = zod_1.z.object({
        response: zod_1.z.object({
            totalBillingNodes: zod_1.z.number(),
            billingNodes: zod_1.z.array(models_1.InfraBillingNodeSchema),
            availableBillingNodes: zod_1.z.array(models_1.InfraBillingAvailableNodeSchema),
            totalAvailableBillingNodes: zod_1.z.number(),
            stats: zod_1.z.object({
                upcomingNodesCount: zod_1.z.number(),
                currentMonthPayments: zod_1.z.number(),
                totalSpent: zod_1.z.number(),
            }),
        }),
    });
})(CreateInfraBillingNodeCommand || (exports.CreateInfraBillingNodeCommand = CreateInfraBillingNodeCommand = {}));
