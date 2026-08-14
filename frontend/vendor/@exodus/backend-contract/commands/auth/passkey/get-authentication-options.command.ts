import { z } from 'zod';

import { AUTH_ROUTES, REST_API } from '../../../api';
import { getEndpointDetails } from '../../../constants';

export namespace GetPasskeyAuthenticationOptionsCommand {
    export const url = REST_API.AUTH.PASSKEY.GET_AUTHENTICATION_OPTIONS;
    export const TSQ_url = url;

    export const endpointDetails = getEndpointDetails(
        AUTH_ROUTES.PASSKEY.GET_AUTHENTICATION_OPTIONS,
        'get',
        'Get the authentication options for passkey',
        { scope: 'get-authentication-options', kind: 'read' },
    );

    export const ResponseSchema = z.object({
        response: z.unknown(),
    });

    export type Response = z.infer<typeof ResponseSchema>;
}
