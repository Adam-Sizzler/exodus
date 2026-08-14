import { z } from 'zod';

import { CONNECTIONS_ROUTES, REST_API } from '../../api';
import { getEndpointDetails } from '../../constants';
import { numberParamSchema } from '../../models';

export namespace ConnectionsByUserCommand {
    export const url = REST_API.CONNECTIONS.CONNECTIONS_BY_USER;
    export const TSQ_url = url(':userId');

    export const endpointDetails = getEndpointDetails(
        CONNECTIONS_ROUTES.CONNECTIONS_BY_USER(':userId'),
        'post',
        'Request Connections for User',
        { scope: 'by-user', kind: 'read' },
    );

    export const RequestParamSchema = z.object({
        userId: numberParamSchema,
    });

    export const ResponseSchema = z.object({
        response: z.object({
            jobId: z.string(),
        }),
    });

    export type RequestParam = z.infer<typeof RequestParamSchema>;
    export type Response = z.infer<typeof ResponseSchema>;
}
