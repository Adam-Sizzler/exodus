import { z } from 'zod';
export declare namespace BulkDeleteUsersCommand {
    const url: "/api/users/bulk/delete";
    const TSQ_url: "/api/users/bulk/delete";
    const endpointDetails: import("../../../constants").EndpointDetails;
    const RequestBodySchema: z.ZodObject<{
        userIds: z.ZodArray<z.ZodNumber>;
    }, z.core.$strip>;
    type RequestBody = z.infer<typeof RequestBodySchema>;
}
//# sourceMappingURL=bulk-delete-users.command.d.ts.map