import { z } from 'zod';
export declare namespace BulkNodesActionsCommand {
    const url: "/api/nodes/bulk-actions";
    const TSQ_url: "/api/nodes/bulk-actions";
    const endpointDetails: import("../../../constants").EndpointDetails;
    const RequestBodySchema: z.ZodObject<{
        uuids: z.ZodArray<z.ZodUUID>;
        action: z.ZodEnum<{
            readonly ENABLE: "ENABLE";
            readonly DISABLE: "DISABLE";
            readonly RESTART: "RESTART";
            readonly RESET_TRAFFIC: "RESET_TRAFFIC";
        }>;
    }, z.core.$strip>;
    type RequestBody = z.infer<typeof RequestBodySchema>;
}
//# sourceMappingURL=actions.command.d.ts.map