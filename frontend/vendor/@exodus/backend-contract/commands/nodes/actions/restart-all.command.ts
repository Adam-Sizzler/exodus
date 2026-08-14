import { z } from 'zod';

import { NODES_ROUTES, REST_API } from '../../../api';
import { getEndpointDetails } from '../../../constants';

export namespace RestartAllNodesCommand {
    export const url = REST_API.NODES.ACTIONS.RESTART_ALL;
    export const TSQ_url = url;

    export const endpointDetails = getEndpointDetails(
        NODES_ROUTES.ACTIONS.RESTART_ALL,
        'post',
        'Restart all nodes',
        { scope: 'restart-all', kind: 'write' },
    );

    export const RequestBodySchema = z.object({
        forceRestart: z.boolean(),
    });

    export type RequestBody = z.infer<typeof RequestBodySchema>;
}
