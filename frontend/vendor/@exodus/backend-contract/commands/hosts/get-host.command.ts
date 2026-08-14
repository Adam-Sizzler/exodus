import { z } from 'zod';

import { HOSTS_ROUTES, REST_API } from '../../api';
import { getEndpointDetails } from '../../constants';
import { HostResponseSchema } from './host.response';

export namespace GetHostCommand {
    export const url = REST_API.HOSTS.GET_BY_UUID;
    export const TSQ_url = url(':uuid');

    export const endpointDetails = getEndpointDetails(
        HOSTS_ROUTES.GET_BY_UUID(':uuid'),
        'get',
        'Get a host by UUID',
        { scope: 'get', kind: 'read' },
    );

    export const RequestParamSchema = z.object({
        uuid: z.uuid(),
    });

    export const ResponseSchema = HostResponseSchema;

    export type RequestParam = z.infer<typeof RequestParamSchema>;
    export type Response = z.infer<typeof ResponseSchema>;
}
