"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.DeleteInfraBillingNodeByUuidCommand = void 0;
const zod_1 = require("zod");
const models_1 = require("../../models");
const api_1 = require("../../api");
const constants_1 = require("../../constants");
var DeleteInfraBillingNodeByUuidCommand;
(function (DeleteInfraBillingNodeByUuidCommand) {
    DeleteInfraBillingNodeByUuidCommand.url = api_1.REST_API.INFRA_BILLING.DELETE_BILLING_NODE;
    DeleteInfraBillingNodeByUuidCommand.TSQ_url = DeleteInfraBillingNodeByUuidCommand.url(':uuid');
    DeleteInfraBillingNodeByUuidCommand.endpointDetails = (0, constants_1.getEndpointDetails)(api_1.INFRA_BILLING_ROUTES.DELETE_BILLING_NODE(':uuid'), 'delete', 'Delete infra billing node');
    DeleteInfraBillingNodeByUuidCommand.RequestSchema = zod_1.z.object({
        uuid: zod_1.z.string().uuid(),
    });
    DeleteInfraBillingNodeByUuidCommand.ResponseSchema = zod_1.z.object({
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
})(DeleteInfraBillingNodeByUuidCommand || (exports.DeleteInfraBillingNodeByUuidCommand = DeleteInfraBillingNodeByUuidCommand = {}));
