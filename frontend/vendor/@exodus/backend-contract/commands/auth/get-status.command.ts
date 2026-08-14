import { z } from 'zod';

import { AUTH_ROUTES, REST_API } from '../../api';
import { getEndpointDetails, OAUTH2_PROVIDERS } from '../../constants';

export namespace GetStatusCommand {
    export const url = REST_API.AUTH.GET_STATUS;
    export const TSQ_url = url;

    export const endpointDetails = getEndpointDetails(
        AUTH_ROUTES.GET_STATUS,
        'get',
        'Get the status of the authentication',
        { scope: 'get-status', kind: 'read' },
    );

    export const ResponseSchema = z.object({
        response: z.object({
            isLoginAllowed: z.boolean(),
            isRegisterAllowed: z.boolean(),
            authentication: z.nullable(
                z.object({
                    passkey: z.object({
                        enabled: z.boolean(),
                    }),
                    oauth2: z.object({
                        providers: z.record(z.enum(OAUTH2_PROVIDERS), z.boolean()),
                    }),
                    password: z.object({
                        enabled: z.boolean(),
                    }),
                }),
            ),
            branding: z.object({
                title: z.nullable(z.string()),
                logoUrl: z.nullable(z.string()),
            }),
        }),
    });

    export type Response = z.infer<typeof ResponseSchema>;
}
