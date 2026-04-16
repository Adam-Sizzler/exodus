import { z } from 'zod';
export declare namespace BulkNodesActionsCommand {
    const url: "/api/nodes/bulk-actions";
    const TSQ_url: "/api/nodes/bulk-actions";
    const endpointDetails: import("../../../constants").EndpointDetails;
    const RequestSchema: z.ZodObject<{
        uuids: z.ZodArray<z.ZodString, "many">;
        action: z.ZodNativeEnum<{
            readonly ENABLE: "ENABLE";
            readonly DISABLE: "DISABLE";
            readonly RESTART: "RESTART";
            readonly RESET_TRAFFIC: "RESET_TRAFFIC";
        }>;
    }, "strip", z.ZodTypeAny, {
        action: "ENABLE" | "DISABLE" | "RESTART" | "RESET_TRAFFIC";
        uuids: string[];
    }, {
        action: "ENABLE" | "DISABLE" | "RESTART" | "RESET_TRAFFIC";
        uuids: string[];
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
//# sourceMappingURL=actions.command.d.ts.map