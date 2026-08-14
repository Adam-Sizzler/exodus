import { z } from 'zod';

import { REST_API, USERS_ROUTES } from '../../../api';
import { getEndpointDetails } from '../../../constants';

export namespace BulkAllExtendExpirationDateCommand {
    export const url = REST_API.USERS.BULK.ALL.EXTEND_EXPIRATION_DATE;
    export const TSQ_url = url;

    export const endpointDetails = getEndpointDetails(
        USERS_ROUTES.BULK.ALL.EXTEND_EXPIRATION_DATE,
        'post',
        'Extend expiration date for all users by days',
        { scope: 'bulk-all-extend-expiration-date', kind: 'write' },
    );

    export const RequestBodySchema = z.object({
        extendDays: z.int().min(1),
    });

    export type RequestBody = z.infer<typeof RequestBodySchema>;
}
