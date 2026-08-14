import { z } from 'zod';

import { AUTH_ROUTES, REST_API } from '../../../api';
import { getEndpointDetails, OAUTH2_PROVIDERS } from '../../../constants';

export namespace OAuth2AuthorizeCommand {
    export const url = REST_API.AUTH.OAUTH2.AUTHORIZE;
    export const TSQ_url = url;

    export const endpointDetails = getEndpointDetails(
        AUTH_ROUTES.OAUTH2.AUTHORIZE,
        'post',
        'Initiate OAuth2 authorization',
        { scope: 'authorize', kind: 'read' },
    );

    export const RequestBodySchema = z.object({
        provider: z.enum(OAUTH2_PROVIDERS),
    });

    export const ResponseSchema = z.object({
        response: z.object({
            authorizationUrl: z.nullable(z.url()),
        }),
    });

    export type RequestBody = z.infer<typeof RequestBodySchema>;
    export type Response = z.infer<typeof ResponseSchema>;
}
