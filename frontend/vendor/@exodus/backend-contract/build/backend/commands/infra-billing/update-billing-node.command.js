"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.UpdateInfraBillingNodeCommand = void 0;
const zod_1 = require("zod");
const models_1 = require("../../models");
const api_1 = require("../../api");
const constants_1 = require("../../constants");
var UpdateInfraBillingNodeCommand;
(function (UpdateInfraBillingNodeCommand) {
    UpdateInfraBillingNodeCommand.url = api_1.REST_API.INFRA_BILLING.UPDATE_BILLING_NODE;
    UpdateInfraBillingNodeCommand.TSQ_url = UpdateInfraBillingNodeCommand.url;
    UpdateInfraBillingNodeCommand.endpointDetails = (0, constants_1.getEndpointDetails)(api_1.INFRA_BILLING_ROUTES.UPDATE_BILLING_NODE, 'patch', 'Update infra billing nodes');
    UpdateInfraBillingNodeCommand.RequestSchema = zod_1.z.object({
        uuids: zod_1.z.array(zod_1.z.string().uuid()),
        nextBillingAt: zod_1.z
            .string({
            invalid_type_error: 'Invalid date format',
        })
            .datetime({ message: 'Invalid date format', offset: true, local: true })
            .transform((str) => new Date(str)),
    });
    UpdateInfraBillingNodeCommand.ResponseSchema = zod_1.z.object({
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
})(UpdateInfraBillingNodeCommand || (exports.UpdateInfraBillingNodeCommand = UpdateInfraBillingNodeCommand = {}));
