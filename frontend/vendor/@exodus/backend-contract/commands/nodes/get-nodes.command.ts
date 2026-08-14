import { z } from 'zod';

import { NODES_ROUTES, REST_API } from '../../api';
import { getEndpointDetails } from '../../constants';
import { NodesSchema } from '../../models';

export namespace GetNodesCommand {
    export const url = REST_API.NODES.GET;
    export const TSQ_url = url;

    export const endpointDetails = getEndpointDetails(NODES_ROUTES.GET, 'get', 'Get nodes', {
        scope: 'list',
        kind: 'read',
    });

    export const ResponseSchema = z.object({
        response: z.array(NodesSchema),
    });

    export type Response = z.infer<typeof ResponseSchema>;
}
