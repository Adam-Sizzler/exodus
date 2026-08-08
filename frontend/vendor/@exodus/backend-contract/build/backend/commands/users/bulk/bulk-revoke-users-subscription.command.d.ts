import { z } from 'zod';
export declare namespace BulkRevokeUsersSubscriptionCommand {
    const url: "/api/users/bulk/revoke-subscription";
    const TSQ_url: "/api/users/bulk/revoke-subscription";
    const endpointDetails: import("../../../constants").EndpointDetails;
    const RequestBodySchema: z.ZodObject<{
        userIds: z.ZodArray<z.ZodNumber>;
    }, z.core.$strip>;
    type RequestBody = z.infer<typeof RequestBodySchema>;
}
//# sourceMappingURL=bulk-revoke-users-subscription.command.d.ts.map