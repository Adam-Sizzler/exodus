import { z } from 'zod';

import { REST_API, USERS_ROUTES } from '../../../api';
import { getEndpointDetails } from '../../../constants';

export namespace BulkUpdateUsersSquadsCommand {
    export const url = REST_API.USERS.BULK.UPDATE_SQUADS;
    export const TSQ_url = url;

    export const endpointDetails = getEndpointDetails(
        USERS_ROUTES.BULK.UPDATE_SQUADS,
        'post',
        'Bulk update users internal squads by User IDs',
        { scope: 'bulk-update-squads', kind: 'write' },
    );

    export const RequestBodySchema = z.object({
        userIds: z.array(z.number()).min(1).max(500),
        activeInternalSquads: z.array(z.uuid()),
    });

    export type RequestBody = z.infer<typeof RequestBodySchema>;
}
