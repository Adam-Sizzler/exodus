import { z } from 'zod';

import { HWID_ROUTES, REST_API } from '../../api';
import { getEndpointDetails } from '../../constants';

export namespace GetTopUsersByHwidDevicesCommand {
    export const url = REST_API.HWID.TOP_USERS_BY_DEVICES;
    export const TSQ_url = url;

    export const endpointDetails = getEndpointDetails(
        HWID_ROUTES.TOP_USERS_BY_DEVICES,
        'get',
        'Get top users by HWID devices',
        { scope: 'top-users', kind: 'read' },
    );

    export const RequestQuerySchema = z.object({
        start: z.coerce
            .number()
            .default(0)
            .describe('Start index (offset) of the results to return, default is 0'),
        size: z.coerce
            .number()
            .min(1, 'Size (limit) must be greater than 0')
            .max(100, 'Size (limit) must be less than 100')
            .describe('Number of results to return, no more than 100')
            .default(5),
    });

    export const ResponseSchema = z.object({
        response: z.object({
            users: z.array(
                z.object({
                    id: z.number(),
                    username: z.string(),
                    devicesCount: z.number(),
                }),
            ),
            total: z.number(),
        }),
    });

    export type RequestQuery = z.infer<typeof RequestQuerySchema>;
    export type Response = z.infer<typeof ResponseSchema>;
}
