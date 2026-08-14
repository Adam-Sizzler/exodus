import { z } from 'zod';

import { REST_API, API_TOKENS_ROUTES } from '../../api';
import { getEndpointDetails } from '../../constants';

export namespace GetApiTokenScopesCommand {
    export const url = REST_API.API_TOKENS.GET_SCOPES;
    export const TSQ_url = url;

    export const endpointDetails = getEndpointDetails(
        API_TOKENS_ROUTES.GET_SCOPES,
        'get',
        'Get available API token scopes',
        { scope: 'list-scopes', kind: 'read' },
        'Returns the catalog of scopes that can be granted to an API token, grouped by resource. Forbidden via "API-key", admin JWT only.',
    );

    export const EndpointScopeSchema = z.object({
        key: z.string(),
        kind: z.enum(['read', 'write']),
        method: z.string(),
        path: z.string(),
        description: z.string(),
    });

    export const ResourceScopesSchema = z.object({
        resource: z.string(),
        resourceScopes: z.array(z.string()),
        endpoints: z.array(EndpointScopeSchema),
    });

    export const ResponseSchema = z.object({
        response: z.object({
            wildcard: z.string(),
            resources: z.array(ResourceScopesSchema),
        }),
    });

    export type Response = z.infer<typeof ResponseSchema>;
}
