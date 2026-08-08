import { z } from 'zod';
export declare namespace BulkAllExtendExpirationDateCommand {
    const url: "/api/users/bulk/all/extend-expiration-date";
    const TSQ_url: "/api/users/bulk/all/extend-expiration-date";
    const endpointDetails: import("../../../constants").EndpointDetails;
    const RequestBodySchema: z.ZodObject<{
        extendDays: z.ZodInt;
    }, z.core.$strip>;
    type RequestBody = z.infer<typeof RequestBodySchema>;
}
//# sourceMappingURL=bulk-all-extend-expiration-date.command.d.ts.map