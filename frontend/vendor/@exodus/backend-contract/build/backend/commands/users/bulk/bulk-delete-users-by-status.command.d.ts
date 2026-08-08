import { z } from 'zod';
export declare namespace BulkDeleteUsersByStatusCommand {
    const url: "/api/users/bulk/delete-by-status";
    const TSQ_url: "/api/users/bulk/delete-by-status";
    const endpointDetails: import("../../../constants").EndpointDetails;
    const RequestBodySchema: z.ZodObject<{
        status: z.ZodEnum<{
            readonly ACTIVE: "ACTIVE";
            readonly DISABLED: "DISABLED";
            readonly LIMITED: "LIMITED";
            readonly EXPIRED: "EXPIRED";
        }>;
    }, z.core.$strip>;
    type RequestBody = z.infer<typeof RequestBodySchema>;
}
//# sourceMappingURL=bulk-delete-users-by-status.command.d.ts.map