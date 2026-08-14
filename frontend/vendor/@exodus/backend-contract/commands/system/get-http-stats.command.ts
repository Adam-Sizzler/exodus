import { z } from 'zod';

import { REST_API, SYSTEM_ROUTES } from '../../api';
import { getEndpointDetails } from '../../constants';

export namespace GetHttpStatsCommand {
    export const url = REST_API.SYSTEM.STATS.HTTP;
    export const TSQ_url = url;

    export const endpointDetails = getEndpointDetails(
        SYSTEM_ROUTES.STATS.HTTP,
        'get',
        'Get HTTP Stats',
        { scope: 'http', kind: 'read' },
    );

    export const ResponseSchema = z.object({
        response: z.object({
            routes: z.array(
                z.object({
                    method: z.string(),
                    route: z.string(),
                    count: z.int32().nonnegative(),
                }),
            ),
            total: z.int32().nonnegative(),
        }),
    });

    export type Response = z.infer<typeof ResponseSchema>;
}
