import { z } from 'zod';

import { CONNECTIONS_ROUTES, REST_API } from '../../api';
import { getEndpointDetails } from '../../constants';

export namespace ConnectionsByUserResultCommand {
    export const url = REST_API.CONNECTIONS.CONNECTIONS_BY_USER_RESULT;
    export const TSQ_url = url(':jobId');

    export const endpointDetails = getEndpointDetails(
        CONNECTIONS_ROUTES.CONNECTIONS_BY_USER_RESULT(':jobId'),
        'get',
        'Get Connections for User by Job ID',
        { scope: 'by-user-result', kind: 'read' },
    );

    export const RequestParamSchema = z.object({
        jobId: z.string(),
    });

    export const ResponseSchema = z.object({
        response: z.object({
            isCompleted: z.boolean(),
            isFailed: z.boolean(),
            progress: z.object({
                total: z.number(),
                completed: z.number(),
                percent: z.number(),
            }),
            result: z
                .object({
                    success: z.boolean(),
                    userId: z.number(),
                    nodes: z.array(
                        z.object({
                            nodeUuid: z.uuid(),
                            nodeName: z.string(),
                            countryCode: z.string(),
                            ips: z.array(
                                z.object({
                                    ip: z.string(),
                                    lastSeen: z.iso
                                        .datetime({
                                            local: true,
                                            offset: true,
                                        })
                                        .transform((str) => new Date(str)),
                                }),
                            ),
                        }),
                    ),
                })
                .nullable(),
        }),
    });

    export type RequestParam = z.infer<typeof RequestParamSchema>;
    export type Response = z.infer<typeof ResponseSchema>;
}
