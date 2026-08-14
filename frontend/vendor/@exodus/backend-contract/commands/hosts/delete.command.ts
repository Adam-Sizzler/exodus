import { z } from 'zod';

import { HOSTS_ROUTES, REST_API } from '../../api';
import { getEndpointDetails } from '../../constants';

export namespace DeleteHostCommand {
    export const url = REST_API.HOSTS.DELETE;
    export const TSQ_url = url(':uuid');

    export const endpointDetails = getEndpointDetails(
        HOSTS_ROUTES.DELETE(':uuid'),
        'delete',
        'Delete a host by UUID',
        { scope: 'delete', kind: 'write' },
    );

    export const RequestParamSchema = z.object({
        uuid: z.uuid(),
    });

    export type RequestParam = z.infer<typeof RequestParamSchema>;
}
