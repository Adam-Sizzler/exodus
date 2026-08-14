import { z } from 'zod';

import { BANDWIDTH_STATS_ROUTES, REST_API } from '../../../api';
import { getEndpointDetails } from '../../../constants';
import { numberParamSchema } from '../../../models';

export namespace GetInternalSquadUserUsageCommand {
    export const url = REST_API.BANDWIDTH_STATS.INTERNAL_SQUADS.USER_USAGE;
    export const TSQ_url = url(':squadUuid', ':userId');

    export const endpointDetails = getEndpointDetails(
        BANDWIDTH_STATS_ROUTES.INTERNAL_SQUADS.USER_USAGE(':squadUuid', ':userId'),
        'get',
        'Get a single user daily traffic usage on the internal squad nodes for a period',
        { scope: 'internal-squad-user-usage', kind: 'read' },
        'Returns users whose total usage over the period on the given nodes is >= minTotalBytes, scoped to the nodes reachable via the Internal Squad inbounds. Every day in the range is present (zero-filled). Underlying usage data is flushed to the database roughly every 2 minutes.',
    );

    export const RequestParamSchema = z.object({
        squadUuid: z.uuid().describe('Internal squad UUID'),
        userId: numberParamSchema,
    });

    export const RequestQuerySchema = z.object({
        start: z.iso.date().describe('Start date (YYYY-MM-DD)'),
        end: z.iso.date().describe('End date (YYYY-MM-DD)'),
    });

    export const ResponseSchema = z.object({
        response: z.object({
            days: z.array(
                z.object({
                    date: z.string().describe('Day (YYYY-MM-DD)'),
                    nodes: z.array(
                        z.object({
                            uuid: z.uuid(),
                            totalBytes: z
                                .number()
                                .describe('Used bytes on this node that day (raw bytes)'),
                        }),
                    ),
                }),
            ),
        }),
    });

    export type RequestParam = z.infer<typeof RequestParamSchema>;
    export type RequestQuery = z.infer<typeof RequestQuerySchema>;
    export type Response = z.infer<typeof ResponseSchema>;
}
