import { z } from 'zod';

import { NODES_ROUTES, REST_API } from '../../../api';
import { getEndpointDetails } from '../../../constants';
import { NodeResponseSchema } from '../node.response';

export namespace DisableNodeCommand {
    export const url = REST_API.NODES.ACTIONS.DISABLE;
    export const TSQ_url = url(':uuid');

    export const endpointDetails = getEndpointDetails(
        NODES_ROUTES.ACTIONS.DISABLE(':uuid'),
        'post',
        'Disable a node',
        { scope: 'disable', kind: 'write' },
    );

    export const RequestParamSchema = z.object({
        uuid: z.uuid(),
    });

    export const ResponseSchema = NodeResponseSchema;

    export type RequestParam = z.infer<typeof RequestParamSchema>;
    export type Response = z.infer<typeof ResponseSchema>;
}
