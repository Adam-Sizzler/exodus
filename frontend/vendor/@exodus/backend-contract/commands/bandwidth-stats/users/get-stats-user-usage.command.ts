import { z } from 'zod';

import { BANDWIDTH_STATS_ROUTES, REST_API } from '../../../api';
import { getEndpointDetails } from '../../../constants';
import { numberParamSchema } from '../../../models';

export namespace GetStatsUserUsageCommand {
    export const url = REST_API.BANDWIDTH_STATS.USERS.GET_BY_ID;
    export const TSQ_url = url(':userId');

    export const endpointDetails = getEndpointDetails(
        BANDWIDTH_STATS_ROUTES.USERS.GET_BY_ID(':userId'),
        'get',
        'Get User Usage by Range',
        { scope: 'user-usage', kind: 'read' },
    );

    export const RequestParamSchema = z.object({
        userId: numberParamSchema.describe('ID of the user'),
    });

    export const RequestQuerySchema = z.object({
        start: z.iso.date().describe('Start date (YYYY-MM-DD)'),
        end: z.iso.date().describe('End date (YYYY-MM-DD)'),
        topNodesLimit: z.coerce.number().min(1).default(20),
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

    export type RequestParam = z.infer<typeof RequestParamSchema>;
    export type RequestQuery = z.infer<typeof RequestQuerySchema>;
    export type Response = z.infer<typeof ResponseSchema>;
}
