import { z } from 'zod';
export declare namespace PluginExecutorCommand {
    const url: "/api/node-plugins/executor";
    const TSQ_url: "/api/node-plugins/executor";
    const endpointDetails: import("../../constants").EndpointDetails;
    const CommandSchema: z.ZodDiscriminatedUnion<"command", [z.ZodObject<{
        command: z.ZodLiteral<"blockIps">;
        ips: z.ZodArray<z.ZodObject<{
            ip: z.ZodString;
            timeout: z.ZodNumber;
        }, "strip", z.ZodTypeAny, {
            ip: string;
            timeout: number;
        }, {
            ip: string;
            timeout: number;
        }>, "many">;
    }, "strip", z.ZodTypeAny, {
        ips: {
            ip: string;
            timeout: number;
        }[];
        command: "blockIps";
    }, {
        ips: {
            ip: string;
            timeout: number;
        }[];
        command: "blockIps";
    }>, z.ZodObject<{
        command: z.ZodLiteral<"unblockIps">;
        ips: z.ZodArray<z.ZodString, "many">;
    }, "strip", z.ZodTypeAny, {
        ips: string[];
        command: "unblockIps";
    }, {
        ips: string[];
        command: "unblockIps";
    }>, z.ZodObject<{
        command: z.ZodLiteral<"recreateTables">;
    }, "strip", z.ZodTypeAny, {
        command: "recreateTables";
    }, {
        command: "recreateTables";
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
        command: z.ZodDiscriminatedUnion<"command", [z.ZodObject<{
            command: z.ZodLiteral<"blockIps">;
            ips: z.ZodArray<z.ZodObject<{
                ip: z.ZodString;
                timeout: z.ZodNumber;
            }, "strip", z.ZodTypeAny, {
                ip: string;
                timeout: number;
            }, {
                ip: string;
                timeout: number;
            }>, "many">;
        }, "strip", z.ZodTypeAny, {
            ips: {
                ip: string;
                timeout: number;
            }[];
            command: "blockIps";
        }, {
            ips: {
                ip: string;
                timeout: number;
            }[];
            command: "blockIps";
        }>, z.ZodObject<{
            command: z.ZodLiteral<"unblockIps">;
            ips: z.ZodArray<z.ZodString, "many">;
        }, "strip", z.ZodTypeAny, {
            ips: string[];
            command: "unblockIps";
        }, {
            ips: string[];
            command: "unblockIps";
        }>, z.ZodObject<{
            command: z.ZodLiteral<"recreateTables">;
        }, "strip", z.ZodTypeAny, {
            command: "recreateTables";
        }, {
            command: "recreateTables";
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
        targetNodes: {
            target: "allNodes";
        } | {
            target: "specificNodes";
            nodeUuids: string[];
        };
        command: {
            ips: {
                ip: string;
                timeout: number;
            }[];
            command: "blockIps";
        } | {
            ips: string[];
            command: "unblockIps";
        } | {
            command: "recreateTables";
        };
    }, {
        targetNodes: {
            target: "allNodes";
        } | {
            target: "specificNodes";
            nodeUuids: string[];
        };
        command: {
            ips: {
                ip: string;
                timeout: number;
            }[];
            command: "blockIps";
        } | {
            ips: string[];
            command: "unblockIps";
        } | {
            command: "recreateTables";
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
//# sourceMappingURL=executor.command.d.ts.map