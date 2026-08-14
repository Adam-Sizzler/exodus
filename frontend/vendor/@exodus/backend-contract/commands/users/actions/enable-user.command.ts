import { z } from 'zod';

import { REST_API, USERS_ROUTES } from '../../../api';
import { getEndpointDetails } from '../../../constants';
import { numberParamSchema } from '../../../models';
import { UserResponseSchema } from '../user.response';

export namespace EnableUserCommand {
    export const url = REST_API.USERS.ACTIONS.ENABLE;
    export const TSQ_url = url(':userId');

    export const endpointDetails = getEndpointDetails(
        USERS_ROUTES.ACTIONS.ENABLE(':userId'),
        'post',
        'Enable user',
        { scope: 'enable', kind: 'write' },
    );

    export const RequestParamSchema = z.object({
        userId: numberParamSchema,
    });

    export const ResponseSchema = UserResponseSchema;

    export type RequestParam = z.infer<typeof RequestParamSchema>;
    export type Response = z.infer<typeof ResponseSchema>;
}
