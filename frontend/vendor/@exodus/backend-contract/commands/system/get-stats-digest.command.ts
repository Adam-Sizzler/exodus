import { z } from 'zod';

import { REST_API, SYSTEM_ROUTES } from '../../api';
import { getEndpointDetails } from '../../constants';

export namespace GetStatsDigestCommand {
    export const url = REST_API.SYSTEM.STATS.DIGEST;
    export const TSQ_url = url;

    export const endpointDetails = getEndpointDetails(
        SYSTEM_ROUTES.STATS.DIGEST,
        'get',
        'Get Stats Digest',
        { scope: 'digest', kind: 'read' },
        'Aggregated statistics for a datetime range [start, end): created and expired users, total traffic, traffic spent by users created within the range and new HWID devices. Per-user traffic history is stored with daily granularity (UTC), so the "traffic by new users" metric snaps to whole days at the range edges.',
    );

    export const RequestQuerySchema = z.object({
        start: z.iso
            .datetime({ offset: true })
            .describe(
                'Start of the range, ISO 8601 datetime with timezone (e.g. 2026-07-15T00:00:00Z). Inclusive.',
            ),
        end: z.iso
            .datetime({ offset: true })
            .describe(
                'End of the range, ISO 8601 datetime with timezone (e.g. 2026-07-16T00:00:00Z). Exclusive.',
            ),
    });

    export const ResponseSchema = z.object({
        response: z.object({
            users: z.object({
                createdCount: z.number(),
                expiredCount: z.number(),
            }),
            traffic: z.object({
                totalBytes: z.string(),
                byUsersCreatedInRangeBytes: z.string(),
            }),
            hwidDevices: z.object({
                createdCount: z.number(),
            }),
        }),
    });

    export type RequestQuery = z.infer<typeof RequestQuerySchema>;
    export type Response = z.infer<typeof ResponseSchema>;
}
