import { z } from 'zod';

import { REST_API, API_TOKENS_ROUTES } from '../../api';
import { getEndpointDetails } from '../../constants';

export namespace DeleteApiTokenCommand {
    export const url = REST_API.API_TOKENS.DELETE;
    export const TSQ_url = url(':uuid');

    export const endpointDetails = getEndpointDetails(
        API_TOKENS_ROUTES.DELETE(':uuid'),
        'delete',
        'Delete API token',
        { scope: 'delete', kind: 'write' },
        'This endpoint is forbidden to use via "API-key". It can be used only with an admin JWT-token.',
    );

    export const RequestParamSchema = z.object({
        uuid: z.uuid().describe('UUID of the API token'),
    });

    export type RequestParam = z.infer<typeof RequestParamSchema>;
}
