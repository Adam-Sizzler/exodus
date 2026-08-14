import { z } from 'zod';

import { REST_API, USERS_ROUTES } from '../../../api';
import { getEndpointDetails } from '../../../constants';
import { numberParamSchema } from '../../../models';
import { UserResponseSchema } from '../user.response';

export namespace ExtendUserCommand {
    export const url = REST_API.USERS.ACTIONS.EXTEND_EXPIRATION_DATE;
    export const TSQ_url = url(':userId');

    export const endpointDetails = getEndpointDetails(
        USERS_ROUTES.ACTIONS.EXTEND_EXPIRATION_DATE(':userId'),
        'post',
        'Extend user expiration date',
        { scope: 'extend', kind: 'write' },
        'If user status is EXPIRED, the new expiration date is calculated from the current date and the user becomes ACTIVE. If user status is ACTIVE, the given number of days is added to the existing expiration date. DISABLED and LIMITED users will be extended, but their status will not change.',
    );

    export const RequestParamSchema = z.object({
        userId: numberParamSchema,
    });

    export const RequestBodySchema = z.object({
        days: z.number().min(1).describe('The number of days to extend the expiration date.'),
    });

    export const ResponseSchema = UserResponseSchema;

    export type RequestParam = z.infer<typeof RequestParamSchema>;
    export type RequestBody = z.infer<typeof RequestBodySchema>;
    export type Response = z.infer<typeof ResponseSchema>;
}
