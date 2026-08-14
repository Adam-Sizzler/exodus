import { z } from 'zod';

import { INFRA_BILLING_ROUTES, REST_API } from '../../api';
import { getEndpointDetails } from '../../constants';
import { InfraProviderSchema } from '../../models';

export namespace GetInfraProviderCommand {
    export const url = REST_API.INFRA_BILLING.GET_PROVIDER_BY_UUID;
    export const TSQ_url = url(':uuid');

    export const endpointDetails = getEndpointDetails(
        INFRA_BILLING_ROUTES.GET_PROVIDER_BY_UUID(':uuid'),
        'get',
        'Get infra provider by uuid',
        { scope: 'get-provider', kind: 'read' },
    );

    export const RequestParamSchema = z.object({
        uuid: z.uuid(),
    });

    export const ResponseSchema = z.object({
        response: InfraProviderSchema,
    });

    export type RequestParam = z.infer<typeof RequestParamSchema>;
    export type Response = z.infer<typeof ResponseSchema>;
}
