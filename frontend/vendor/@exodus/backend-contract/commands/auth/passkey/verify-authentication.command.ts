import { z } from 'zod';

import { AUTH_ROUTES, REST_API } from '../../../api';
import { getEndpointDetails } from '../../../constants';

export namespace VerifyPasskeyAuthenticationCommand {
    export const url = REST_API.AUTH.PASSKEY.VERIFY_AUTHENTICATION;
    export const TSQ_url = url;

    export const endpointDetails = getEndpointDetails(
        AUTH_ROUTES.PASSKEY.VERIFY_AUTHENTICATION,
        'post',
        'Verify the authentication for passkey',
        { scope: 'verify-authentication', kind: 'write' },
    );

    export const RequestBodySchema = z.object({
        response: z.unknown(),
    });

    export const ResponseSchema = z.object({
        response: z.object({
            accessToken: z.string(),
        }),
    });

    export type RequestBody = z.infer<typeof RequestBodySchema>;
    export type Response = z.infer<typeof ResponseSchema>;
}
