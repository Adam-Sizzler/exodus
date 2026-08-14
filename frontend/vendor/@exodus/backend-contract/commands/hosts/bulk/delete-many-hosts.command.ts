import { z } from 'zod';

import { HOSTS_ROUTES, REST_API } from '../../../api';
import { getEndpointDetails } from '../../../constants';

export namespace BulkDeleteHostsCommand {
    export const url = REST_API.HOSTS.BULK.DELETE_HOSTS;
    export const TSQ_url = url;

    export const endpointDetails = getEndpointDetails(
        HOSTS_ROUTES.BULK.DELETE_HOSTS,
        'post',
        'Delete hosts by UUIDs',
        { scope: 'bulk-delete', kind: 'write' },
    );

    export const RequestBodySchema = z.object({
        uuids: z.array(z.uuid()),
    });

    export type RequestBody = z.infer<typeof RequestBodySchema>;
}
