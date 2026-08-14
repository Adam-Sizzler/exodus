import { z } from 'zod';

import { INFRA_BILLING_ROUTES, REST_API } from '../../api';
import { getEndpointDetails } from '../../constants';
import { InfraProviderSchema } from '../../models';

export namespace GetInfraProvidersCommand {
    export const url = REST_API.INFRA_BILLING.GET_PROVIDERS;
    export const TSQ_url = url;

    export const endpointDetails = getEndpointDetails(
        INFRA_BILLING_ROUTES.GET_PROVIDERS,
        'get',
        'Get all infra providers',
        { scope: 'list-providers', kind: 'read' },
    );

    export const ResponseSchema = z.object({
        response: z.object({
            total: z.number(),
            providers: z.array(InfraProviderSchema),
        }),
    });

    export type Response = z.infer<typeof ResponseSchema>;
}
