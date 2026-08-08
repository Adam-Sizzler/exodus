import { z } from 'zod';
export declare namespace DropConnectionsCommand {
    const url: "/api/connections/drop";
    const TSQ_url: "/api/connections/drop";
    const endpointDetails: import("../../constants").EndpointDetails;
    const DropBySchema: z.ZodDiscriminatedUnion<[z.ZodObject<{
        by: z.ZodLiteral<"userIds">;
        userIds: z.ZodArray<z.ZodNumber>;
    }, z.core.$strip>, z.ZodObject<{
        by: z.ZodLiteral<"ipAddresses">;
        ipAddresses: z.ZodArray<z.ZodUnion<readonly [z.ZodIPv4, z.ZodIPv6]>>;
    }, z.core.$strip>], "by">;
    const TargetNodesSchema: z.ZodDiscriminatedUnion<[z.ZodObject<{
        target: z.ZodLiteral<"allNodes">;
    }, z.core.$strip>, z.ZodObject<{
        target: z.ZodLiteral<"specificNodes">;
        nodeUuids: z.ZodArray<z.ZodUUID>;
    }, z.core.$strip>], "target">;
    const RequestBodySchema: z.ZodObject<{
        dropBy: z.ZodDiscriminatedUnion<[z.ZodObject<{
            by: z.ZodLiteral<"userIds">;
            userIds: z.ZodArray<z.ZodNumber>;
        }, z.core.$strip>, z.ZodObject<{
            by: z.ZodLiteral<"ipAddresses">;
            ipAddresses: z.ZodArray<z.ZodUnion<readonly [z.ZodIPv4, z.ZodIPv6]>>;
        }, z.core.$strip>], "by">;
        targetNodes: z.ZodDiscriminatedUnion<[z.ZodObject<{
            target: z.ZodLiteral<"allNodes">;
        }, z.core.$strip>, z.ZodObject<{
            target: z.ZodLiteral<"specificNodes">;
            nodeUuids: z.ZodArray<z.ZodUUID>;
        }, z.core.$strip>], "target">;
    }, z.core.$strip>;
    type RequestBody = z.infer<typeof RequestBodySchema>;
}
//# sourceMappingURL=drop-connections.command.d.ts.map