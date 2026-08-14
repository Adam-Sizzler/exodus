import { z } from 'zod';

import { API_TOKENS_ROUTES, REST_API } from '../../api';
import { getEndpointDetails } from '../../constants';

export namespace GetOttCommand {
    export const url = REST_API.API_TOKENS.OTT;
    export const TSQ_url = url;

    export const endpointDetails = getEndpointDetails(
        API_TOKENS_ROUTES.OTT,
        'post',
        'Get short-lived token for accessing backend tools (Swagger, Scalar, Bull Board)',
        { scope: 'ott', kind: 'write' },
    );

    export const ResponseSchema = z.object({
        response: z.object({
            ott: z.string(),
        }),
    });

    export type Response = z.infer<typeof ResponseSchema>;
}
