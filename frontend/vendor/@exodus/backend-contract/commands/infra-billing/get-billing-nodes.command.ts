import { z } from 'zod';

import { INFRA_BILLING_ROUTES, REST_API } from '../../api';
import { getEndpointDetails } from '../../constants';
import { InfraBillingAvailableNodeSchema, InfraBillingNodeSchema } from '../../models';

export namespace GetInfraBillingNodesCommand {
    export const url = REST_API.INFRA_BILLING.GET_BILLING_NODES;
    export const TSQ_url = url;

    export const endpointDetails = getEndpointDetails(
        INFRA_BILLING_ROUTES.GET_BILLING_NODES,
        'get',
        'Get infra billing nodes',
        { scope: 'list-billing-nodes', kind: 'read' },
    );

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

    export type Response = z.infer<typeof ResponseSchema>;
}
