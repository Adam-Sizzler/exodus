import { z } from 'zod';

import { REST_API, USERS_ROUTES } from '../../../api';
import { getEndpointDetails } from '../../../constants';

export namespace BulkDeleteUsersCommand {
    export const url = REST_API.USERS.BULK.DELETE;
    export const TSQ_url = url;

    export const endpointDetails = getEndpointDetails(
        USERS_ROUTES.BULK.DELETE,
        'post',
        'Bulk delete users by User IDs',
        { scope: 'bulk-delete', kind: 'write' },
    );

    export const RequestBodySchema = z.object({
        userIds: z.array(z.number()).min(1).max(500),
    });

    export type RequestBody = z.infer<typeof RequestBodySchema>;
}
