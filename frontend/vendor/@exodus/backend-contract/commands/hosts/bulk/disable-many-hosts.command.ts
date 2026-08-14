import { z } from 'zod';

import { HOSTS_ROUTES, REST_API } from '../../../api';
import { getEndpointDetails } from '../../../constants';

export namespace BulkDisableHostsCommand {
    export const url = REST_API.HOSTS.BULK.DISABLE_HOSTS;
    export const TSQ_url = url;

    export const endpointDetails = getEndpointDetails(
        HOSTS_ROUTES.BULK.DISABLE_HOSTS,
        'post',
        'Disable hosts by UUIDs',
        { scope: 'bulk-disable', kind: 'write' },
    );

    export const RequestBodySchema = z.object({
        uuids: z.array(z.uuid()),
    });

    export type RequestBody = z.infer<typeof RequestBodySchema>;
}
