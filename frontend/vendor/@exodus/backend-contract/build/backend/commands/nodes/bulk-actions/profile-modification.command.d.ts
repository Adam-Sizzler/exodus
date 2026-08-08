import { z } from 'zod';
export declare namespace BulkNodesProfileModificationCommand {
    const url: "/api/nodes/bulk-actions/profile-modification";
    const TSQ_url: "/api/nodes/bulk-actions/profile-modification";
    const endpointDetails: import("../../../constants").EndpointDetails;
    const RequestBodySchema: z.ZodObject<{
        uuids: z.ZodArray<z.ZodUUID>;
        configProfile: z.ZodObject<{
            activeConfigProfileUuid: z.ZodUUID;
            activeInbounds: z.ZodArray<z.ZodUUID>;
        }, z.core.$strip>;
    }, z.core.$strip>;
    type RequestBody = z.infer<typeof RequestBodySchema>;
}
//# sourceMappingURL=profile-modification.command.d.ts.map