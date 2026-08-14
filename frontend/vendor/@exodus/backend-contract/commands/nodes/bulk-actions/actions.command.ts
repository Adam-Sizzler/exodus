import { z } from 'zod';

import { NODES_ROUTES, REST_API } from '../../../api';
import { getEndpointDetails, NODES_BULK_ACTIONS } from '../../../constants';

export namespace BulkNodesActionsCommand {
    export const url = REST_API.NODES.BULK_ACTIONS.ACTIONS;
    export const TSQ_url = url;

    export const endpointDetails = getEndpointDetails(
        NODES_ROUTES.BULK_ACTIONS.ACTIONS,
        'post',
        'Perform actions for many nodes',
        { scope: 'bulk-actions', kind: 'write' },
    );

    export const RequestBodySchema = z.object({
        uuids: z.array(z.uuid()).min(1),
        action: z.enum(NODES_BULK_ACTIONS),
    });

    export type RequestBody = z.infer<typeof RequestBodySchema>;
}
