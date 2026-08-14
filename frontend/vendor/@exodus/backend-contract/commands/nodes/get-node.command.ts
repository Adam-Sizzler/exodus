import { z } from 'zod';

import { NODES_ROUTES, REST_API } from '../../api';
import { getEndpointDetails } from '../../constants';
import { NodeResponseSchema } from './node.response';

export namespace GetNodeCommand {
    export const url = REST_API.NODES.GET_BY_UUID;
    export const TSQ_url = url(':uuid');

    export const endpointDetails = getEndpointDetails(
        NODES_ROUTES.GET_BY_UUID(':uuid'),
        'get',
        'Get node by UUID',
        { scope: 'get', kind: 'read' },
    );

    export const RequestParamSchema = z.object({
        uuid: z.uuid(),
    });

    export const ResponseSchema = NodeResponseSchema;

    export type RequestParam = z.infer<typeof RequestParamSchema>;
    export type Response = z.infer<typeof ResponseSchema>;
}
