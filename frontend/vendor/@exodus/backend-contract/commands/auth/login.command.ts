import { z } from 'zod';

import { AUTH_ROUTES, REST_API } from '../../api';
import { getEndpointDetails } from '../../constants';

export namespace LoginCommand {
    export const url = REST_API.AUTH.LOGIN;
    export const TSQ_url = url;

    export const endpointDetails = getEndpointDetails(
        AUTH_ROUTES.LOGIN,
        'post',
        'Login as superadmin',
        { scope: 'login', kind: 'write' },
    );

    export const RequestBodySchema = z.object({
        username: z.string().describe('Username of the user'),
        password: z.string().describe('Password of the user'),
    });

    export const ResponseSchema = z.object({
        response: z.object({
            accessToken: z.string(),
        }),
    });

    export type RequestBody = z.infer<typeof RequestBodySchema>;
    export type Response = z.infer<typeof ResponseSchema>;
}
