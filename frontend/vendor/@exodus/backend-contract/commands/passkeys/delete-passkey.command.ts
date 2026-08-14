import { z } from 'zod';

import { PASSKEYS_ROUTES, REST_API } from '../../api';
import { getEndpointDetails } from '../../constants';

export namespace DeletePasskeyCommand {
    export const url = REST_API.PASSKEYS.DELETE_PASSKEY;
    export const TSQ_url = url;

    export const endpointDetails = getEndpointDetails(
        PASSKEYS_ROUTES.DELETE_PASSKEY,
        'delete',
        'Delete a passkey by ID',
        { scope: 'delete', kind: 'write' },
    );

    export const RequestBodySchema = z.object({
        id: z.string(),
    });

    export type RequestBody = z.infer<typeof RequestBodySchema>;
}
