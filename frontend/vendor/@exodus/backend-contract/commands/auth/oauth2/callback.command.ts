import { z } from 'zod';

import { AUTH_ROUTES, REST_API } from '../../../api';
import { getEndpointDetails, OAUTH2_PROVIDERS } from '../../../constants';

export namespace OAuth2CallbackCommand {
    export const url = REST_API.AUTH.OAUTH2.CALLBACK;
    export const TSQ_url = url;

    export const endpointDetails = getEndpointDetails(
        AUTH_ROUTES.OAUTH2.CALLBACK,
        'post',
        'Callback from OAuth2',
        { scope: 'callback', kind: 'write' },
    );

    export const RequestBodySchema = z.object({
        provider: z.enum(OAUTH2_PROVIDERS),
        code: z.string(),
        state: z.string(),
    });

    export const ResponseSchema = z.object({
        response: z.object({
            accessToken: z.string(),
        }),
    });

    export type RequestBody = z.infer<typeof RequestBodySchema>;
    export type Response = z.infer<typeof ResponseSchema>;
}
