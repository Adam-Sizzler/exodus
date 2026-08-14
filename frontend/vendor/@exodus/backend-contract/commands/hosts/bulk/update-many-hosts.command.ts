import { z } from 'zod';

import { HOSTS_ROUTES, REST_API } from '../../../api';
import { getEndpointDetails } from '../../../constants';
import { UpdateHostCommand } from '../update.command';

export namespace UpdateManyHostsCommand {
    export const url = REST_API.HOSTS.BULK.UPDATE;
    export const TSQ_url = url;

    export const endpointDetails = getEndpointDetails(
        HOSTS_ROUTES.BULK.UPDATE,
        'patch',
        'Update many hosts',
        { scope: 'bulk-update', kind: 'write' },
    );

    export const RequestBodySchema = UpdateHostCommand.RequestBodySchema.omit({ uuid: true })
        .partial()
        .extend({
            uuids: z.array(z.uuid()).min(1),
        });

    export type RequestBody = z.infer<typeof RequestBodySchema>;
}
