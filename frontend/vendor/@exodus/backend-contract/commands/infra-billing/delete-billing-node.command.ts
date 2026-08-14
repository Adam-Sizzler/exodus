import { z } from 'zod';

import { INFRA_BILLING_ROUTES, REST_API } from '../../api';
import { getEndpointDetails } from '../../constants';

export namespace DeleteInfraBillingNodeCommand {
    export const url = REST_API.INFRA_BILLING.DELETE_BILLING_NODE;
    export const TSQ_url = url(':uuid');

    export const endpointDetails = getEndpointDetails(
        INFRA_BILLING_ROUTES.DELETE_BILLING_NODE(':uuid'),
        'delete',
        'Delete infra billing node',
        { scope: 'delete-billing-node', kind: 'write' },
    );

    export const RequestParamSchema = z.object({
        uuid: z.uuid(),
    });

    export type RequestParam = z.infer<typeof RequestParamSchema>;
}
