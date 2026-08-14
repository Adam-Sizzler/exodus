import { z } from 'zod';

import { CONNECTIONS_ROUTES, REST_API } from '../../api';
import { getEndpointDetails } from '../../constants';

export namespace ConnectionsByNodeResultCommand {
    export const url = REST_API.CONNECTIONS.CONNECTIONS_BY_NODE_RESULT;
    export const TSQ_url = url(':jobId');

    export const endpointDetails = getEndpointDetails(
        CONNECTIONS_ROUTES.CONNECTIONS_BY_NODE_RESULT(':jobId'),
        'get',
        'Get Connections for Node by Job ID',
        { scope: 'by-node-result', kind: 'read' },
    );

    export const RequestParamSchema = z.object({
        jobId: z.string(),
    });

    export const ResponseSchema = z.object({
        response: z.object({
            isCompleted: z.boolean(),
            isFailed: z.boolean(),
            result: z
                .object({
                    success: z.boolean(),
                    nodeUuid: z.uuid(),
                    users: z.array(
                        z.object({
                            userId: z.number(),
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
