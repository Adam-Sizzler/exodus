import { z } from 'zod';
export declare namespace BulkAllExtendExpirationDateCommand {
    const url: "/api/users/bulk/all/extend-expiration-date";
    const TSQ_url: "/api/users/bulk/all/extend-expiration-date";
    const endpointDetails: import("../../../constants").EndpointDetails;
    const RequestSchema: z.ZodObject<{
        extendDays: z.ZodNumber;
    }, "strip", z.ZodTypeAny, {
        extendDays: number;
    }, {
        extendDays: number;
    }>;
    type Request = z.infer<typeof RequestSchema>;
    const ResponseSchema: z.ZodObject<{
        response: z.ZodObject<{
            eventSent: z.ZodBoolean;
        }, "strip", z.ZodTypeAny, {
            eventSent: boolean;
        }, {
            eventSent: boolean;
        }>;
    }, "strip", z.ZodTypeAny, {
        response: {
            eventSent: boolean;
        };
    }, {
        response: {
            eventSent: boolean;
        };
    }>;
    type Response = z.infer<typeof ResponseSchema>;
}
//# sourceMappingURL=bulk-all-extend-expiration-date.command.d.ts.map