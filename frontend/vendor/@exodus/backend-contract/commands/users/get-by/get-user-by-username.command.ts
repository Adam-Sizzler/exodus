import { z } from 'zod';

import { REST_API, USERS_ROUTES } from '../../../api';
import { getEndpointDetails } from '../../../constants';
import { UserResponseSchema } from '../user.response';

export namespace GetUserByUsernameCommand {
    export const url = REST_API.USERS.GET_BY.USERNAME;
    export const TSQ_url = url(':username');

    export const endpointDetails = getEndpointDetails(
        USERS_ROUTES.GET_BY.USERNAME(':username'),
        'get',
        'Get user by username',
        { scope: 'by-username', kind: 'read' },
    );

    export const RequestParamSchema = z.object({
        username: z.string(),
    });

    export const ResponseSchema = UserResponseSchema;

    export type RequestParam = z.infer<typeof RequestParamSchema>;
    export type Response = z.infer<typeof ResponseSchema>;
}
