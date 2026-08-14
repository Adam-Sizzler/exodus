import { z } from 'zod';

import { REST_API, USERS_ROUTES } from '../../api';
import { getEndpointDetails } from '../../constants';

export namespace ResolveUserCommand {
    export const url = REST_API.USERS.RESOLVE;
    export const TSQ_url = url;

    export const endpointDetails = getEndpointDetails(
        USERS_ROUTES.RESOLVE,
        'post',
        'Resolve a user',
        { scope: 'resolve', kind: 'read' },
        'Resolve a user by ID, Short UUID or username. Exactly one of the fields must be provided.',
    );

    export const RequestBodySchema = z
        .object({
            id: z.number().optional(),
            shortUuid: z.string().optional(),
            username: z.string().optional(),
        })
        .refine(
            (data) => {
                const provided = [data.id, data.shortUuid, data.username].filter(
                    (v) => v !== undefined,
                );
                return provided.length === 1;
            },
            {
                error: 'Exactly one of id, shortUuid, or username must be provided',
            },
        );

    export const ResponseSchema = z.object({
        response: z.object({
            id: z.number(),
            username: z.string(),
            shortUuid: z.string(),
        }),
    });

    export type RequestBody = z.infer<typeof RequestBodySchema>;
    export type Response = z.infer<typeof ResponseSchema>;
}
