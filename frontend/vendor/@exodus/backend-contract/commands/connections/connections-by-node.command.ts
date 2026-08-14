import { z } from 'zod';

import { CONNECTIONS_ROUTES, REST_API } from '../../api';
import { getEndpointDetails } from '../../constants';

export namespace ConnectionsByNodeCommand {
    export const url = REST_API.CONNECTIONS.CONNECTIONS_BY_NODE;
    export const TSQ_url = url(':nodeUuid');

    export const endpointDetails = getEndpointDetails(
        CONNECTIONS_ROUTES.CONNECTIONS_BY_NODE(':nodeUuid'),
        'post',
        'Request Connections for Node',
        { scope: 'by-node', kind: 'read' },
    );

    export const RequestParamSchema = z.object({
        nodeUuid: z.uuid().describe('Node UUID'),
    });

    export const ResponseSchema = z.object({
        response: z.object({
            jobId: z.string(),
        }),
    });

    export type RequestParam = z.infer<typeof RequestParamSchema>;
    export type Response = z.infer<typeof ResponseSchema>;
}
