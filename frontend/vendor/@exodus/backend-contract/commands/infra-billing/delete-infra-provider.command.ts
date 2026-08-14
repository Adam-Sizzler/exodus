import { z } from 'zod';

import { INFRA_BILLING_ROUTES, REST_API } from '../../api';
import { getEndpointDetails } from '../../constants';

export namespace DeleteInfraProviderCommand {
    export const url = REST_API.INFRA_BILLING.DELETE_PROVIDER;
    export const TSQ_url = url(':uuid');

    export const endpointDetails = getEndpointDetails(
        INFRA_BILLING_ROUTES.DELETE_PROVIDER(':uuid'),
        'delete',
        'Delete infra provider by uuid',
        { scope: 'delete-provider', kind: 'write' },
    );

    export const RequestParamSchema = z.object({
        uuid: z.uuid(),
    });

    export type RequestParam = z.infer<typeof RequestParamSchema>;
}
