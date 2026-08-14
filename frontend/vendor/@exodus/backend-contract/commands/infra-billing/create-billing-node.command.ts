import { z } from 'zod';

import { INFRA_BILLING_ROUTES, REST_API } from '../../api';
import { getEndpointDetails } from '../../constants';
import { InfraBillingAvailableNodeSchema, InfraBillingNodeSchema } from '../../models';

export namespace CreateInfraBillingNodeCommand {
    export const url = REST_API.INFRA_BILLING.CREATE_BILLING_NODE;
    export const TSQ_url = url;

    export const endpointDetails = getEndpointDetails(
        INFRA_BILLING_ROUTES.CREATE_BILLING_NODE,
        'post',
        'Create infra billing node',
        { scope: 'create-billing-node', kind: 'write' },
    );

    export const RequestBodySchema = z.object({
        providerUuid: z.uuid(),
        nodeUuid: z.uuid().nullable(),
        name: z.string().min(1).max(255).nullable(),
        nextBillingAt: z.iso
            .datetime({ offset: true, local: true })
            .transform((str) => new Date(str))
            .describe('Next billing date. Format: 2025-01-17T15:38:45.065Z'),
    });

    export const ResponseSchema = z.object({
        response: z.object({
            totalBillingNodes: z.number(),
            billingNodes: z.array(InfraBillingNodeSchema),
            availableBillingNodes: z.array(InfraBillingAvailableNodeSchema),
            totalAvailableBillingNodes: z.number(),
            stats: z.object({
                upcomingNodesCount: z.number(),
                currentMonthPayments: z.number(),
                totalSpent: z.number(),
            }),
        }),
    });

    export type RequestBody = z.infer<typeof RequestBodySchema>;
    export type Response = z.infer<typeof ResponseSchema>;
}
