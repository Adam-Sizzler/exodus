import { z } from 'zod';

import { REST_API, USERS_ROUTES } from '../../api';
import { getEndpointDetails } from '../../constants';
import { numberParamSchema } from '../../models';
import { UserResponseSchema } from './user.response';

export namespace GetUserByIdCommand {
    export const url = REST_API.USERS.GET_BY_ID;
    export const TSQ_url = url(':userId');

    export const endpointDetails = getEndpointDetails(
        USERS_ROUTES.GET_BY_ID(':userId'),
        'get',
        'Get user by ID',
        { scope: 'by-id', kind: 'read' },
    );

    export const RequestParamSchema = z.object({
        userId: numberParamSchema,
    });

    export const ResponseSchema = UserResponseSchema;

    export type RequestParam = z.infer<typeof RequestParamSchema>;
    export type Response = z.infer<typeof ResponseSchema>;
}
