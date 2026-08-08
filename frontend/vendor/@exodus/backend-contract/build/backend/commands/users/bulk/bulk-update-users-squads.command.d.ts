import { z } from 'zod';
export declare namespace BulkUpdateUsersSquadsCommand {
    const url: "/api/users/bulk/update-squads";
    const TSQ_url: "/api/users/bulk/update-squads";
    const endpointDetails: import("../../../constants").EndpointDetails;
    const RequestBodySchema: z.ZodObject<{
        userIds: z.ZodArray<z.ZodNumber>;
        activeInternalSquads: z.ZodArray<z.ZodUUID>;
    }, z.core.$strip>;
    type RequestBody = z.infer<typeof RequestBodySchema>;
}
//# sourceMappingURL=bulk-update-users-squads.command.d.ts.map