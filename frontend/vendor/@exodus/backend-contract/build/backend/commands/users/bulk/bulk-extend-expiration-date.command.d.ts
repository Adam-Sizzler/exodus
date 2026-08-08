import { z } from 'zod';
export declare namespace BulkExtendExpirationDateCommand {
    const url: "/api/users/bulk/extend-expiration-date";
    const TSQ_url: "/api/users/bulk/extend-expiration-date";
    const endpointDetails: import("../../../constants").EndpointDetails;
    const RequestBodySchema: z.ZodObject<{
        userIds: z.ZodArray<z.ZodNumber>;
        extendDays: z.ZodInt;
    }, z.core.$strip>;
    type RequestBody = z.infer<typeof RequestBodySchema>;
}
//# sourceMappingURL=bulk-extend-expiration-date.command.d.ts.map