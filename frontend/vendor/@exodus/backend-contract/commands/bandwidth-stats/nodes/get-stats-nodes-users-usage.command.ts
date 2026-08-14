import { z } from 'zod';

import { BANDWIDTH_STATS_ROUTES, REST_API } from '../../../api';
import { getEndpointDetails } from '../../../constants';

export namespace GetStatsNodesUsersUsageCommand {
    export const url = REST_API.BANDWIDTH_STATS.NODES.GET_USERS_BY_NODES;
    export const TSQ_url = url;

    export const endpointDetails = getEndpointDetails(
        BANDWIDTH_STATS_ROUTES.NODES.GET_USERS_BY_NODES,
        'post',
        'Get Nodes Users Usage by Nodes UUIDs',
        { scope: 'nodes-users-usage', kind: 'read' },
    );

    export const RequestBodySchema = z.object({
        nodesUuids: z.array(z.uuid()).min(1),
    });

    export const RequestQuerySchema = z.object({
        start: z.iso.date().describe('Start date (YYYY-MM-DD)'),
        end: z.iso.date().describe('End date (YYYY-MM-DD)'),
        topUsersLimit: z.coerce.number().min(1).default(100),
    });

    export const ResponseSchema = z.object({
        response: z.object({
            categories: z.array(z.string()),
            sparklineData: z.array(z.number()),
            topUsers: z.array(
                z.object({
                    color: z.string(),
                    username: z.string(),
                    total: z.number(),
                }),
            ),
        }),
    });

    export type RequestBody = z.infer<typeof RequestBodySchema>;
    export type RequestQuery = z.infer<typeof RequestQuerySchema>;
    export type Response = z.infer<typeof ResponseSchema>;
}
