import { z } from 'zod';

import { REST_API, USERS_ROUTES } from '../../api';
import { getEndpointDetails } from '../../constants';
import { numberParamSchema } from '../../models';

export namespace DeleteUserCommand {
    export const url = REST_API.USERS.DELETE;
    export const TSQ_url = url(':userId');

    export const endpointDetails = getEndpointDetails(
        USERS_ROUTES.DELETE(':userId'),
        'delete',
        'Delete user',
        { scope: 'delete', kind: 'write' },
    );

    export const RequestParamSchema = z.object({
        userId: numberParamSchema,
    });

    export type RequestParam = z.infer<typeof RequestParamSchema>;
}
