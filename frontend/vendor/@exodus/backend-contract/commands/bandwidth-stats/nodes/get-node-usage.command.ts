import { z } from 'zod';

import { BANDWIDTH_STATS_ROUTES, REST_API } from '../../../api';
import { getEndpointDetails } from '../../../constants';

export namespace GetNodeUsageCommand {
    export const url = REST_API.BANDWIDTH_STATS.NODES.GET_USAGE;
    export const TSQ_url = url;

    export const endpointDetails = getEndpointDetails(
        BANDWIDTH_STATS_ROUTES.NODES.GET_USAGE,
        'post',
        'Get users exceeding a traffic threshold on the given nodes for a period',
        { scope: 'node-usage', kind: 'read' },
        'Returns users whose total usage over the period on the given nodes is >= minTotalBytes. Underlying usage data is flushed to the database roughly every 2 minutes.',
    );

    export const RequestBodySchema = z.object({
        nodesUuids: z.array(z.uuid()).min(1).describe('Node UUIDs to include'),
    });

    export const RequestQuerySchema = z.object({
        start: z.iso.date().describe('Start date (YYYY-MM-DD)'),
        end: z.iso.date().describe('End date (YYYY-MM-DD)'),
        minTotalBytes: z.coerce
            .number()
            .min(0)
            .optional()
            .default(0)
            .describe('Only include users whose total usage over the period is >= this (bytes)'),
    });

    export const ResponseSchema = z.object({
        response: z.object({
            nodes: z.array(
                z.object({
                    uuid: z.uuid(),
                    users: z.array(
                        z.object({
                            id: z.number(),
                            totalBytes: z
                                .number()
                                .describe('Total used bytes over the period (raw bytes)'),
                        }),
                    ),
                }),
            ),
        }),
    });

    export type RequestBody = z.infer<typeof RequestBodySchema>;
    export type RequestQuery = z.infer<typeof RequestQuerySchema>;
    export type Response = z.infer<typeof ResponseSchema>;
}
