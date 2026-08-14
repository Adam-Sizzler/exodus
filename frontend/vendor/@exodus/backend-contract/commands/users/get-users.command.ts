import { z } from 'zod';

import { REST_API, USERS_ROUTES } from '../../api';
import { getEndpointDetails } from '../../constants';
import { ExtendedUsersSchema, TanstackQueryRequestQuerySchema } from '../../models';

export namespace GetUsersCommand {
    export const url = REST_API.USERS.GET;
    export const TSQ_url = url;

    export const endpointDetails = getEndpointDetails(
        USERS_ROUTES.GET,
        'get',
        'Get all users using offset-based pagination',
        {
            scope: 'list',
            kind: 'read',
        },
        'Please note that the filters here are primarily intended for use by the frontend and rely on expensive operators such as LIKE under the hood. Misusing these filters may negatively impact the performance of your database.',
    );

    export const RequestQuerySchema = TanstackQueryRequestQuerySchema;

    export const ResponseSchema = z.object({
        response: z.object({
            users: z.array(ExtendedUsersSchema),
            total: z.number(),
        }),
    });

    export type RequestQuery = z.infer<typeof RequestQuerySchema>;
    export type Response = z.infer<typeof ResponseSchema>;
}
