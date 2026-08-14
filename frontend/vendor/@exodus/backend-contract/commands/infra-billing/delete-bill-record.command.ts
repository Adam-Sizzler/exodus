import { z } from 'zod';

import { INFRA_BILLING_ROUTES, REST_API } from '../../api';
import { getEndpointDetails } from '../../constants';

export namespace DeleteInfraBillingRecordCommand {
    export const url = REST_API.INFRA_BILLING.DELETE_BILLING_HISTORY;
    export const TSQ_url = url(':uuid');

    export const endpointDetails = getEndpointDetails(
        INFRA_BILLING_ROUTES.DELETE_BILLING_HISTORY(':uuid'),
        'delete',
        'Delete infra billing history',
        { scope: 'delete-bill-record', kind: 'write' },
    );

    export const RequestParamSchema = z.object({
        uuid: z.uuid(),
    });

    export type RequestParam = z.infer<typeof RequestParamSchema>;
}
