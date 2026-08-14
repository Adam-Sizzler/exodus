import { z } from 'zod';

import { REST_API, API_TOKENS_ROUTES } from '../../api';
import { getEndpointDetails } from '../../constants';
import { ApiTokensSchema } from '../../models/api-tokens.schema';

export namespace GetApiTokensCommand {
    export const url = REST_API.API_TOKENS.GET;
    export const TSQ_url = url;

    export const endpointDetails = getEndpointDetails(
        API_TOKENS_ROUTES.GET,
        'get',
        'Get all API tokens',
        { scope: 'list', kind: 'read' },
        'This endpoint is forbidden to use via "API-key". It can only be used with admin JWT-token.',
    );

    export const ResponseSchema = z.object({
        response: z.object({
            tokens: z.array(ApiTokensSchema),
        }),
    });

    export type Response = z.infer<typeof ResponseSchema>;
}
