import { z } from 'zod';

import { REST_API, USERS_ROUTES } from '../../../api';
import { getEndpointDetails } from '../../../constants';

export namespace BulkResetTrafficUsersCommand {
    export const url = REST_API.USERS.BULK.RESET_TRAFFIC;
    export const TSQ_url = url;

    export const endpointDetails = getEndpointDetails(
        USERS_ROUTES.BULK.RESET_TRAFFIC,
        'post',
        'Bulk reset traffic users by User IDs',
        { scope: 'bulk-reset-traffic', kind: 'write' },
    );

    export const RequestBodySchema = z.object({
        userIds: z.array(z.number()).min(1).max(500),
    });

    export type RequestBody = z.infer<typeof RequestBodySchema>;
}
