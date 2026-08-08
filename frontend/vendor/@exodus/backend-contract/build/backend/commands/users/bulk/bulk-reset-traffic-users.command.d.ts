import { z } from 'zod';
export declare namespace BulkResetTrafficUsersCommand {
    const url: "/api/users/bulk/reset-traffic";
    const TSQ_url: "/api/users/bulk/reset-traffic";
    const endpointDetails: import("../../../constants").EndpointDetails;
    const RequestBodySchema: z.ZodObject<{
        userIds: z.ZodArray<z.ZodNumber>;
    }, z.core.$strip>;
    type RequestBody = z.infer<typeof RequestBodySchema>;
}
//# sourceMappingURL=bulk-reset-traffic-users.command.d.ts.map