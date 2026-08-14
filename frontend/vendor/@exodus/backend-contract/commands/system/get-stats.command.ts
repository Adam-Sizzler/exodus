import { z } from 'zod';

import { REST_API, SYSTEM_ROUTES } from '../../api';
import { getEndpointDetails } from '../../constants';
import { USERS_STATUS } from '../../constants';

export namespace GetStatsCommand {
    export const url = REST_API.SYSTEM.STATS.SYSTEM_STATS;
    export const TSQ_url = url;

    export const endpointDetails = getEndpointDetails(
        SYSTEM_ROUTES.STATS.SYSTEM_STATS,
        'get',
        'Get Stats',
        { scope: 'stats', kind: 'read' },
    );

    export const RequestQuerySchema = z.object({
        tz: z.string().optional(),
    });

    export const ResponseSchema = z.object({
        response: z.object({
            cpu: z.object({
                cores: z.number(),
            }),
            memory: z.object({
                total: z.number(),
                free: z.number(),
                used: z.number(),
            }),
            uptime: z.number(),
            timestamp: z.number(),
            users: z.object({
                statusCounts: z.record(z.enum(USERS_STATUS), z.number()),
                totalUsers: z.number(),
            }),
            onlineStats: z.object({
                lastDay: z.number(),
                lastWeek: z.number(),
                neverOnline: z.number(),
                onlineNow: z.number(),
            }),
            nodes: z.object({
                totalOnline: z.number(),
                totalBytesLifetime: z.string(),
            }),
        }),
    });

    export type RequestQuery = z.infer<typeof RequestQuerySchema>;
    export type Response = z.infer<typeof ResponseSchema>;
}
