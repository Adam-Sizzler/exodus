import { z } from 'zod';
export declare namespace DropConnectionsCommand {
    const url: "/api/ip-control/drop-connections";
    const TSQ_url: "/api/ip-control/drop-connections";
    const endpointDetails: import("../../constants").EndpointDetails;
    const DropBySchema: z.ZodDiscriminatedUnion<"by", [z.ZodObject<{
        by: z.ZodLiteral<"userUuids">;
        userUuids: z.ZodArray<z.ZodString, "many">;
    }, "strip", z.ZodTypeAny, {
        by: "userUuids";
        userUuids: string[];
    }, {
        by: "userUuids";
        userUuids: string[];
    }>, z.ZodObject<{
        by: z.ZodLiteral<"ipAddresses">;
        ipAddresses: z.ZodArray<z.ZodString, "many">;
    }, "strip", z.ZodTypeAny, {
        by: "ipAddresses";
        ipAddresses: string[];
    }, {
        by: "ipAddresses";
        ipAddresses: string[];
    }>]>;
    const TargetNodesSchema: z.ZodDiscriminatedUnion<"target", [z.ZodObject<{
        target: z.ZodLiteral<"allNodes">;
    }, "strip", z.ZodTypeAny, {
        target: "allNodes";
    }, {
        target: "allNodes";
    }>, z.ZodObject<{
        target: z.ZodLiteral<"specificNodes">;
        nodeUuids: z.ZodArray<z.ZodString, "many">;
    }, "strip", z.ZodTypeAny, {
        target: "specificNodes";
        nodeUuids: string[];
    }, {
        target: "specificNodes";
        nodeUuids: string[];
    }>]>;
    const RequestSchema: z.ZodObject<{
        dropBy: z.ZodDiscriminatedUnion<"by", [z.ZodObject<{
            by: z.ZodLiteral<"userUuids">;
            userUuids: z.ZodArray<z.ZodString, "many">;
        }, "strip", z.ZodTypeAny, {
            by: "userUuids";
            userUuids: string[];
        }, {
            by: "userUuids";
            userUuids: string[];
        }>, z.ZodObject<{
            by: z.ZodLiteral<"ipAddresses">;
            ipAddresses: z.ZodArray<z.ZodString, "many">;
        }, "strip", z.ZodTypeAny, {
            by: "ipAddresses";
            ipAddresses: string[];
        }, {
            by: "ipAddresses";
            ipAddresses: string[];
        }>]>;
        targetNodes: z.ZodDiscriminatedUnion<"target", [z.ZodObject<{
            target: z.ZodLiteral<"allNodes">;
        }, "strip", z.ZodTypeAny, {
            target: "allNodes";
        }, {
            target: "allNodes";
        }>, z.ZodObject<{
            target: z.ZodLiteral<"specificNodes">;
            nodeUuids: z.ZodArray<z.ZodString, "many">;
        }, "strip", z.ZodTypeAny, {
            target: "specificNodes";
            nodeUuids: string[];
        }, {
            target: "specificNodes";
            nodeUuids: string[];
        }>]>;
    }, "strip", z.ZodTypeAny, {
        dropBy: {
            by: "userUuids";
            userUuids: string[];
        } | {
            by: "ipAddresses";
            ipAddresses: string[];
        };
        targetNodes: {
            target: "allNodes";
        } | {
            target: "specificNodes";
            nodeUuids: string[];
        };
    }, {
        dropBy: {
            by: "userUuids";
            userUuids: string[];
        } | {
            by: "ipAddresses";
            ipAddresses: string[];
        };
        targetNodes: {
            target: "allNodes";
        } | {
            target: "specificNodes";
            nodeUuids: string[];
        };
    }>;
    type Request = z.infer<typeof RequestSchema>;
    const ResponseSchema: z.ZodObject<{
        response: z.ZodObject<{
            eventSent: z.ZodBoolean;
        }, "strip", z.ZodTypeAny, {
            eventSent: boolean;
        }, {
            eventSent: boolean;
        }>;
    }, "strip", z.ZodTypeAny, {
        response: {
            eventSent: boolean;
        };
    }, {
        response: {
            eventSent: boolean;
        };
    }>;
    type Response = z.infer<typeof ResponseSchema>;
}
//# sourceMappingURL=drop-connections.command.d.ts.map