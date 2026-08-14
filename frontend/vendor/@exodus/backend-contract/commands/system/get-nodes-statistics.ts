import { z } from 'zod';

import { REST_API, SYSTEM_ROUTES } from '../../api';
import { getEndpointDetails } from '../../constants';

export namespace GetNodesStatisticsCommand {
    export const url = REST_API.SYSTEM.STATS.NODES_STATS;
    export const TSQ_url = url;

    export const endpointDetails = getEndpointDetails(
        SYSTEM_ROUTES.STATS.NODES_STATS,
        'get',
        'Get Nodes Statistics',
        { scope: 'nodes-statistics', kind: 'read' },
    );

    export const RequestQuerySchema = z.object({
        tz: z.string().optional(),
    });

    export const ResponseSchema = z.object({
        response: z.object({
            lastSevenDays: z.array(
                z.object({
                    nodeName: z.string(),
                    date: z.string(),
                    totalBytes: z.string(),
                }),
            ),
        }),
    });

    export type RequestQuery = z.infer<typeof RequestQuerySchema>;
    export type Response = z.infer<typeof ResponseSchema>;
}
