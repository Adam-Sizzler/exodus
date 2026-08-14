import { z } from 'zod';

import { HOSTS_ROUTES, REST_API } from '../../../api';
import { getEndpointDetails } from '../../../constants';

export namespace BulkEnableHostsCommand {
    export const url = REST_API.HOSTS.BULK.ENABLE_HOSTS;
    export const TSQ_url = url;

    export const endpointDetails = getEndpointDetails(
        HOSTS_ROUTES.BULK.ENABLE_HOSTS,
        'post',
        'Enable hosts by UUIDs',
        { scope: 'bulk-enable', kind: 'write' },
    );

    export const RequestBodySchema = z.object({
        uuids: z.array(z.uuid()),
    });

    export type RequestBody = z.infer<typeof RequestBodySchema>;
}
