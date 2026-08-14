import { z } from 'zod';

import { PASSKEYS_ROUTES, REST_API } from '../../api';
import { getEndpointDetails } from '../../constants';

export namespace VerifyPasskeyRegistrationCommand {
    export const url = REST_API.PASSKEYS.VERIFY_REGISTRATION;
    export const TSQ_url = url;

    export const endpointDetails = getEndpointDetails(
        PASSKEYS_ROUTES.VERIFY_REGISTRATION,
        'post',
        'Verify registration for passkey',
        { scope: 'verify-registration', kind: 'write' },
    );

    export const RequestBodySchema = z.object({
        response: z.unknown(),
    });

    export const ResponseSchema = z.object({
        response: z.object({
            verified: z.boolean(),
        }),
    });

    export type RequestBody = z.infer<typeof RequestBodySchema>;
    export type Response = z.infer<typeof ResponseSchema>;
}
