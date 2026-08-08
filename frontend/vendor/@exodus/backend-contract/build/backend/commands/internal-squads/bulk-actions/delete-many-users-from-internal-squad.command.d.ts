import { z } from 'zod';
export declare namespace DeleteManyUsersFromInternalSquadCommand {
    const url: (uuid: string) => string;
    const TSQ_url: string;
    const endpointDetails: import("../../../constants").EndpointDetails;
    const RequestParamSchema: z.ZodObject<{
        uuid: z.ZodUUID;
    }, z.core.$strip>;
    const RequestBodySchema: z.ZodObject<{
        userIds: z.ZodArray<z.ZodNumber>;
    }, z.core.$strip>;
    type RequestParam = z.infer<typeof RequestParamSchema>;
    type RequestBody = z.infer<typeof RequestBodySchema>;
}
//# sourceMappingURL=delete-many-users-from-internal-squad.command.d.ts.map