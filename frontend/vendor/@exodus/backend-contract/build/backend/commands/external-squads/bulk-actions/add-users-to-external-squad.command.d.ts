import { z } from 'zod';
export declare namespace AddUsersToExternalSquadCommand {
    const url: (uuid: string) => string;
    const TSQ_url: string;
    const endpointDetails: import("../../../constants").EndpointDetails;
    const RequestParamSchema: z.ZodObject<{
        uuid: z.ZodUUID;
    }, z.core.$strip>;
    type RequestParam = z.infer<typeof RequestParamSchema>;
}
//# sourceMappingURL=add-users-to-external-squad.command.d.ts.map