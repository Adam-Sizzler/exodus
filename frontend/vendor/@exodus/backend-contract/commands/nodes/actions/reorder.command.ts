import { z } from 'zod';

import { NODES_ROUTES, REST_API } from '../../../api';
import { getEndpointDetails } from '../../../constants';
import { NodesSchema } from '../../../models';

export namespace ReorderNodesCommand {
    export const url = REST_API.NODES.ACTIONS.REORDER;
    export const TSQ_url = url;

    export const endpointDetails = getEndpointDetails(
        NODES_ROUTES.ACTIONS.REORDER,
        'post',
        'Reorder nodes',
        { scope: 'reorder', kind: 'write' },
    );

    export const RequestBodySchema = z.object({
        nodes: z.array(
            NodesSchema.pick({
                viewPosition: true,
                uuid: true,
            }),
        ),
    });

    export const ResponseSchema = z.object({
        response: z.array(NodesSchema),
    });

    export type RequestBody = z.infer<typeof RequestBodySchema>;
    export type Response = z.infer<typeof ResponseSchema>;
}
