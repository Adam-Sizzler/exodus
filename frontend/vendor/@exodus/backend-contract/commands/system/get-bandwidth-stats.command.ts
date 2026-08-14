import { z } from 'zod';

import { REST_API, SYSTEM_ROUTES } from '../../api';
import { getEndpointDetails } from '../../constants';
import { BaseStatSchema } from '../../models/base-stat.schema';

export namespace GetBandwidthStatsCommand {
    export const url = REST_API.SYSTEM.STATS.BANDWIDTH_STATS;
    export const TSQ_url = url;

    export const endpointDetails = getEndpointDetails(
        SYSTEM_ROUTES.STATS.BANDWIDTH_STATS,
        'get',
        'Get Bandwidth Stats',
        { scope: 'bandwidth-stats', kind: 'read' },
    );

    export const RequestQuerySchema = z.object({
        tz: z.string().optional(),
    });

    export const ResponseSchema = z.object({
        response: z.object({
            bandwidthLastTwoDays: BaseStatSchema,
            bandwidthLastSevenDays: BaseStatSchema,
            bandwidthLast30Days: BaseStatSchema,
            bandwidthCalendarMonth: BaseStatSchema,
            bandwidthCurrentYear: BaseStatSchema,
        }),
    });

    export type RequestQuery = z.infer<typeof RequestQuerySchema>;
    export type Response = z.infer<typeof ResponseSchema>;
}
