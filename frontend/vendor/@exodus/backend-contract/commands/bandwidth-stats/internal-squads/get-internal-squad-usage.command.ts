import { z } from 'zod';

import { BANDWIDTH_STATS_ROUTES, REST_API } from '../../../api';
import { getEndpointDetails } from '../../../constants';

export namespace GetInternalSquadUsageCommand {
    export const url = REST_API.BANDWIDTH_STATS.INTERNAL_SQUADS.GET_USAGE;
    export const TSQ_url = url(':uuid');

    export const endpointDetails = getEndpointDetails(
        BANDWIDTH_STATS_ROUTES.INTERNAL_SQUADS.GET_USAGE(':uuid'),
        'get',
        'Get internal squad users traffic usage for a period',
        { scope: 'internal-squad-usage', kind: 'read' },
        'Returns users whose total usage over the period on the given nodes is >= minTotalBytes, scoped to the nodes reachable via the internal squad inbounds. Underlying usage data is flushed to the database roughly every 2 minutes.',
    );

    export const RequestParamSchema = z.object({
        uuid: z.uuid().describe('Internal squad UUID'),
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
        limit: z.coerce
            .number()
            .min(1)
            .max(1000)
            .optional()
            .default(250)
            .describe('Number of users to return, no more than 1000'),
        cursor: z.coerce
            .number()
            .optional()
            .describe('Pass the nextCursor from the previous response. Omit on the first request.'),
    });

    export const ResponseSchema = z.object({
        response: z.object({
            squadUuid: z.uuid(),
            users: z.array(
                z.object({
                    id: z.number(),
                    totalBytes: z.number().describe('Total used bytes over the period (raw bytes)'),
                }),
            ),
            nextCursor: z
                .string()
                .nullable()
                .describe('Cursor to fetch the next page, or null if there are no more results'),
            hasMore: z.boolean().describe('Whether there are more results to fetch'),
        }),
    });

    export type RequestParam = z.infer<typeof RequestParamSchema>;
    export type RequestQuery = z.infer<typeof RequestQuerySchema>;
    export type Response = z.infer<typeof ResponseSchema>;
}
