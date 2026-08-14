import { z } from 'zod';

import { REST_API, USERS_ROUTES } from '../../../api';
import { getEndpointDetails } from '../../../constants';

export namespace BulkExtendExpirationDateCommand {
    export const url = REST_API.USERS.BULK.EXTEND_EXPIRATION_DATE;
    export const TSQ_url = url;

    export const endpointDetails = getEndpointDetails(
        USERS_ROUTES.BULK.EXTEND_EXPIRATION_DATE,
        'post',
        'Extend expiration date for specified users by days',
        { scope: 'bulk-extend-expiration-date', kind: 'write' },
    );

    export const RequestBodySchema = z.object({
        userIds: z.array(z.number()).min(1).max(500),
        extendDays: z.int().min(1).max(9999),
    });

    export type RequestBody = z.infer<typeof RequestBodySchema>;
}
