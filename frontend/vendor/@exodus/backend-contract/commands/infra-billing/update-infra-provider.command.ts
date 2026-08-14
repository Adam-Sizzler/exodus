import { z } from 'zod';

import { INFRA_BILLING_ROUTES, REST_API } from '../../api';
import { getEndpointDetails } from '../../constants';
import { InfraProviderSchema } from '../../models';

export namespace UpdateInfraProviderCommand {
    export const url = REST_API.INFRA_BILLING.UPDATE_PROVIDER;
    export const TSQ_url = url;

    export const endpointDetails = getEndpointDetails(
        INFRA_BILLING_ROUTES.UPDATE_PROVIDER,
        'patch',
        'Update infra provider',
        { scope: 'update-provider', kind: 'write' },
    );

    export const RequestBodySchema = z.object({
        uuid: z.uuid(),
        name: z.string().min(2).max(30).optional(),
        faviconLink: z.url().nullish(),
        loginUrl: z.url().nullish(),
    });

    export const ResponseSchema = z.object({
        response: InfraProviderSchema,
    });

    export type RequestBody = z.infer<typeof RequestBodySchema>;
    export type Response = z.infer<typeof ResponseSchema>;
}
