import { z } from 'zod';

import { BANDWIDTH_STATS_ROUTES, REST_API } from '../../../api';
import { getEndpointDetails } from '../../../constants';

export namespace GetStatsNodesUsageCommand {
    export const url = REST_API.BANDWIDTH_STATS.NODES.GET;
    export const TSQ_url = url;

    export const endpointDetails = getEndpointDetails(
        BANDWIDTH_STATS_ROUTES.NODES.GET,
        'get',
        'Get Nodes Usage by Range',
        { scope: 'nodes-usage', kind: 'read' },
    );

    export const RequestQuerySchema = z.object({
        start: z.iso.date().describe('Start date (YYYY-MM-DD)'),
        end: z.iso.date().describe('End date (YYYY-MM-DD)'),
        topNodesLimit: z.coerce
            .number()
            .min(1)
            .default(20)
            .describe('Limit of top nodes to return'),
    });

    export const ResponseSchema = z.object({
        response: z.object({
            categories: z.array(z.string()),
            sparklineData: z.array(z.number()),
            topNodes: z.array(
                z.object({
                    uuid: z.uuid(),
                    color: z.string(),
                    name: z.string(),
                    countryCode: z.string(),
                    total: z.number(),
                }),
            ),
            series: z.array(
                z.object({
                    uuid: z.uuid(),
                    name: z.string(),
                    color: z.string(),
                    countryCode: z.string(),
                    total: z.number(),
                    data: z.array(z.number()),
                }),
            ),
        }),
    });

    export type RequestQuery = z.infer<typeof RequestQuerySchema>;
    export type Response = z.infer<typeof ResponseSchema>;
}
