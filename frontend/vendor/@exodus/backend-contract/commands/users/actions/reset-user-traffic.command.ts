import { z } from 'zod';

import { REST_API, USERS_ROUTES } from '../../../api';
import { getEndpointDetails } from '../../../constants';
import { numberParamSchema } from '../../../models';
import { UserResponseSchema } from '../user.response';

export namespace ResetUserTrafficCommand {
    export const url = REST_API.USERS.ACTIONS.RESET_TRAFFIC;
    export const TSQ_url = url(':userId');

    export const endpointDetails = getEndpointDetails(
        USERS_ROUTES.ACTIONS.RESET_TRAFFIC(':userId'),
        'post',
        'Reset user traffic',
        { scope: 'reset-traffic', kind: 'write' },
    );

    export const RequestParamSchema = z.object({
        userId: numberParamSchema,
    });

    export const ResponseSchema = UserResponseSchema;

    export type RequestParam = z.infer<typeof RequestParamSchema>;
    export type Response = z.infer<typeof ResponseSchema>;
}
