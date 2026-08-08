import { z } from 'zod';
export declare namespace PluginExecutorCommand {
    const url: "/api/node-plugins/executor";
    const TSQ_url: "/api/node-plugins/executor";
    const endpointDetails: import("../../constants").EndpointDetails;
    const CommandSchema: z.ZodDiscriminatedUnion<[z.ZodObject<{
        command: z.ZodLiteral<"blockIps">;
        ips: z.ZodArray<z.ZodObject<{
            ip: z.ZodUnion<readonly [z.ZodIPv4, z.ZodIPv6]>;
            timeout: z.ZodNumber;
        }, z.core.$strip>>;
    }, z.core.$strip>, z.ZodObject<{
        command: z.ZodLiteral<"unblockIps">;
        ips: z.ZodArray<z.ZodUnion<readonly [z.ZodIPv4, z.ZodIPv6]>>;
    }, z.core.$strip>, z.ZodObject<{
        command: z.ZodLiteral<"recreateTables">;
    }, z.core.$strip>], "command">;
    const TargetNodesSchema: z.ZodDiscriminatedUnion<[z.ZodObject<{
        target: z.ZodLiteral<"allNodes">;
    }, z.core.$strip>, z.ZodObject<{
        target: z.ZodLiteral<"specificNodes">;
        nodeUuids: z.ZodArray<z.ZodUUID>;
    }, z.core.$strip>], "target">;
    const RequestBodySchema: z.ZodObject<{
        command: z.ZodDiscriminatedUnion<[z.ZodObject<{
            command: z.ZodLiteral<"blockIps">;
            ips: z.ZodArray<z.ZodObject<{
                ip: z.ZodUnion<readonly [z.ZodIPv4, z.ZodIPv6]>;
                timeout: z.ZodNumber;
            }, z.core.$strip>>;
        }, z.core.$strip>, z.ZodObject<{
            command: z.ZodLiteral<"unblockIps">;
            ips: z.ZodArray<z.ZodUnion<readonly [z.ZodIPv4, z.ZodIPv6]>>;
        }, z.core.$strip>, z.ZodObject<{
            command: z.ZodLiteral<"recreateTables">;
        }, z.core.$strip>], "command">;
        targetNodes: z.ZodDiscriminatedUnion<[z.ZodObject<{
            target: z.ZodLiteral<"allNodes">;
        }, z.core.$strip>, z.ZodObject<{
            target: z.ZodLiteral<"specificNodes">;
            nodeUuids: z.ZodArray<z.ZodUUID>;
        }, z.core.$strip>], "target">;
    }, z.core.$strip>;
    type RequestBody = z.infer<typeof RequestBodySchema>;
}
//# sourceMappingURL=executor.command.d.ts.map