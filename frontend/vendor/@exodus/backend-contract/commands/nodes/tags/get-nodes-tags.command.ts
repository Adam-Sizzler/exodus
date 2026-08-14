import { z } from 'zod';

import { REST_API, NODES_ROUTES } from '../../../api';
import { getEndpointDetails } from '../../../constants';

export namespace GetNodesTagsCommand {
    export const url = REST_API.NODES.TAGS.GET;
    export const TSQ_url = url;

    export const endpointDetails = getEndpointDetails(
        NODES_ROUTES.TAGS.GET,
        'get',
        'Get nodes tags',
        { scope: 'list-tags', kind: 'read' },
    );

    export const ResponseSchema = z.object({
        response: z.object({
            tags: z.array(z.string()),
        }),
    });

    export type Response = z.infer<typeof ResponseSchema>;
}
