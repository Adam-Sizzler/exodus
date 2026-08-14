import { z } from 'zod';

import { NODES_ROUTES, REST_API } from '../../../api';
import { getEndpointDetails } from '../../../constants';

export namespace ResetNodeTrafficCommand {
    export const url = REST_API.NODES.ACTIONS.RESET_TRAFFIC;
    export const TSQ_url = url(':uuid');

    export const endpointDetails = getEndpointDetails(
        NODES_ROUTES.ACTIONS.RESET_TRAFFIC(':uuid'),
        'post',
        'Reset Node Traffic',
        { scope: 'reset-traffic', kind: 'write' },
    );

    export const RequestParamSchema = z.object({
        uuid: z.uuid(),
    });

    export type RequestParam = z.infer<typeof RequestParamSchema>;
}
