import { z } from 'zod';

import { REST_API, USERS_ROUTES } from '../../api';
import { getEndpointDetails, RESET_PERIODS, USERS_STATUS } from '../../constants';
import { ExtendedUsersSchema } from '../../models';

export namespace GetUsersStreamCommand {
    export const url = REST_API.USERS.STREAM;
    export const TSQ_url = url;

    export const endpointDetails = getEndpointDetails(
        USERS_ROUTES.STREAM,
        'get',
        'Get all users using cursor-based (keyset) pagination with filtering options',
        { scope: 'stream', kind: 'read' },
    );

    export const RequestQuerySchema = z.object({
        cursor: z.coerce
            .number()
            .optional()
            .describe(
                'Cursor for pagination — pass the nextCursor from the previous response. Omit on the first request.',
            ),
        size: z.coerce
            .number()
            .min(1)
            .max(1000)
            .describe('Number of results to return, no more than 1000')
            .optional()
            .default(250),

        // Filtering
        status: z.enum(USERS_STATUS).optional().describe('Status to filter users by'),
        trafficLimitStrategy: z
            .enum(RESET_PERIODS)
            .optional()
            .describe('Traffic limit strategy to filter users by'),
        telegramId: z
            .string()
            .transform(Number)
            .pipe(z.number().nonnegative())
            .optional()
            .describe('Telegram ID to filter users by'),
        email: z.email().optional().describe('Email to filter users by'),
        tag: z.string().optional().describe('Tag to filter users by'),
        externalSquadUuid: z.uuid().optional().describe('External squad UUID to filter users by'),
    });

    export const ResponseSchema = z.object({
        response: z.object({
            users: z.array(ExtendedUsersSchema),
            nextCursor: z
                .string()
                .nullable()
                .describe('Cursor to fetch the next page, or null if there are no more results'),
            hasMore: z.boolean().describe('Whether there are more results to fetch'),
        }),
    });

    export type RequestQuery = z.infer<typeof RequestQuerySchema>;
    export type Response = z.infer<typeof ResponseSchema>;
}
