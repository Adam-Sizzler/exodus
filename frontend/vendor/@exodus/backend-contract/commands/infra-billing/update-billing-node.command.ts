import { z } from 'zod';

import { INFRA_BILLING_ROUTES, REST_API } from '../../api';
import { getEndpointDetails } from '../../constants';
import { InfraBillingAvailableNodeSchema, InfraBillingNodeSchema } from '../../models';

export namespace UpdateInfraBillingNodeCommand {
    export const url = REST_API.INFRA_BILLING.UPDATE_BILLING_NODE;
    export const TSQ_url = url;

    export const endpointDetails = getEndpointDetails(
        INFRA_BILLING_ROUTES.UPDATE_BILLING_NODE,
        'patch',
        'Update infra billing nodes',
        { scope: 'update-billing-node', kind: 'write' },
    );

    export const RequestBodySchema = z.object({
        uuids: z.array(z.uuid()),
        nextBillingAt: z.iso
            .datetime({ offset: true, local: true })
            .transform((str) => new Date(str)),
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
