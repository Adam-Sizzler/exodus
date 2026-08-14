import { z } from 'zod';

import { HOSTS_ROUTES, REST_API } from '../../api';
import { getEndpointDetails } from '../../constants';
import { HostsSchema } from '../../models';

export namespace ReorderHostsCommand {
    export const url = REST_API.HOSTS.ACTIONS.REORDER;
    export const TSQ_url = url;

    export const endpointDetails = getEndpointDetails(
        HOSTS_ROUTES.ACTIONS.REORDER,
        'post',
        'Reorder hosts',
        { scope: 'reorder', kind: 'write' },
    );

    export const RequestBodySchema = z.object({
        hosts: z.array(
            HostsSchema.pick({
                viewPosition: true,
                uuid: true,
            }),
        ),
    });

    export const ResponseSchema = z.object({
        response: z.object({
            isUpdated: z.boolean(),
        }),
    });

    export type RequestBody = z.infer<typeof RequestBodySchema>;
    export type Response = z.infer<typeof ResponseSchema>;
}
