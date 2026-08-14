import { z } from 'zod';

import { NODE_PLUGINS_ROUTES, REST_API } from '../../api';
import { getEndpointDetails } from '../../constants';

export namespace PluginExecutorCommand {
    export const url = REST_API.NODE_PLUGINS.EXECUTOR;
    export const TSQ_url = url;

    export const endpointDetails = getEndpointDetails(
        NODE_PLUGINS_ROUTES.EXECUTOR,
        'post',
        'Execute command on node plugins',
        { scope: 'executor', kind: 'write' },
    );

    export const CommandSchema = z.discriminatedUnion('command', [
        z
            .object({
                command: z.literal('blockIps'),
                ips: z
                    .array(
                        z.object({
                            ip: z.union([z.ipv4(), z.ipv6()]),
                            timeout: z.number(),
                        }),
                    )
                    .min(1),
            })
            .describe('Block IPs'),
        z
            .object({
                command: z.literal('unblockIps'),
                ips: z.array(z.union([z.ipv4(), z.ipv6()])).min(1),
            })
            .describe('Unblock IPs'),
        z
            .object({
                command: z.literal('recreateTables'),
            })
            .describe('Recreate tables'),
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
        command: CommandSchema,
        targetNodes: TargetNodesSchema,
    });

    export type RequestBody = z.infer<typeof RequestBodySchema>;
}
