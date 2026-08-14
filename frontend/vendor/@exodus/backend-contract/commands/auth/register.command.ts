import { z } from 'zod';

import { AUTH_ROUTES, REST_API } from '../../api';
import { getEndpointDetails } from '../../constants';

export namespace RegisterCommand {
    export const url = REST_API.AUTH.REGISTER;
    export const TSQ_url = url;

    export const endpointDetails = getEndpointDetails(
        AUTH_ROUTES.REGISTER,
        'post',
        'Register as superadmin',
        { scope: 'register', kind: 'write' },
    );

    export const RequestBodySchema = z.object({
        username: z.string().describe('Username of the user'),
        password: z
            .string()
            .min(24, 'Password must contain at least 24 characters')
            .regex(
                /^(?=.*?[A-Z])(?=.*?[a-z])(?=.*?[0-9]).{24,}$/,
                'Password must contain uppercase and lowercase letters and numbers, and be at least 24 characters long.',
            ),
    });

    export const ResponseSchema = z.object({
        response: z.object({
            accessToken: z.string(),
        }),
    });

    export type RequestBody = z.infer<typeof RequestBodySchema>;
    export type Response = z.infer<typeof ResponseSchema>;
}
