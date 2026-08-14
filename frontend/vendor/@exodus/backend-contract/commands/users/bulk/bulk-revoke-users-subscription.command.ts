import { z } from 'zod';

import { REST_API, USERS_ROUTES } from '../../../api';
import { getEndpointDetails } from '../../../constants';

export namespace BulkRevokeUsersSubscriptionCommand {
    export const url = REST_API.USERS.BULK.REVOKE_SUBSCRIPTION;
    export const TSQ_url = url;

    export const endpointDetails = getEndpointDetails(
        USERS_ROUTES.BULK.REVOKE_SUBSCRIPTION,
        'post',
        'Revoke users subscription by User IDs',
        { scope: 'bulk-revoke-subscription', kind: 'write' },
    );

    export const RequestBodySchema = z.object({
        userIds: z.array(z.number()).min(1).max(500),
    });

    export type RequestBody = z.infer<typeof RequestBodySchema>;
}
