import { z } from 'zod';

import { CONNECTIONS_ROUTES, REST_API } from '../../api';
import { getEndpointDetails } from '../../constants';

export namespace DropConnectionsCommand {
    export const url = REST_API.CONNECTIONS.DROP_CONNECTIONS;
    export const TSQ_url = url;

    export const endpointDetails = getEndpointDetails(
        CONNECTIONS_ROUTES.DROP_CONNECTIONS,
        'post',
        'Drop Connections for Users or IPs',
        { scope: 'drop', kind: 'write' },
    );

    export const DropBySchema = z.discriminatedUnion('by', [
        z
            .object({
                by: z.literal('userIds'),
                userIds: z.array(z.number()).min(1),
            })
            .describe('Drop by user IDs'),
        z
            .object({
                by: z.literal('ipAddresses'),
                ipAddresses: z.array(z.union([z.ipv4(), z.ipv6()])).min(1),
            })
            .describe('Drop by IP addresses'),
    ]);

    export const TargetNodesSchema = z.discriminatedUnion('target', [
        z
            .object({
                target: z.literal('allNodes'),
            })
            .describe('Target all connected nodes'),
        z
            .object({
                target: z.literal('specificNodes'),
                nodeUuids: z.array(z.uuid()).min(1),
            })
            .describe('Target specific nodes'),
    ]);

    export const RequestBodySchema = z.object({
        dropBy: DropBySchema,
        targetNodes: TargetNodesSchema,
    });

    export type RequestBody = z.infer<typeof RequestBodySchema>;
}
