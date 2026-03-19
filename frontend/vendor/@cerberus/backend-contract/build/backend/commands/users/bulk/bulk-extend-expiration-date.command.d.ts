import { z } from 'zod';
export declare namespace BulkExtendExpirationDateCommand {
    const url: "/api/users/bulk/extend-expiration-date";
    const TSQ_url: "/api/users/bulk/extend-expiration-date";
    const endpointDetails: import("../../../constants").EndpointDetails;
    const RequestSchema: z.ZodObject<{
        uuids: z.ZodArray<z.ZodString, "many">;
        extendDays: z.ZodNumber;
    }, "strip", z.ZodTypeAny, {
        uuids: string[];
        extendDays: number;
    }, {
        uuids: string[];
        extendDays: number;
    }>;
    type Request = z.infer<typeof RequestSchema>;
    const ResponseSchema: z.ZodObject<{
        response: z.ZodObject<{
            affectedRows: z.ZodNumber;
        }, "strip", z.ZodTypeAny, {
            affectedRows: number;
        }, {
            affectedRows: number;
        }>;
    }, "strip", z.ZodTypeAny, {
        response: {
            affectedRows: number;
        };
    }, {
        response: {
            affectedRows: number;
        };
    }>;
    type Response = z.infer<typeof ResponseSchema>;
}
//# sourceMappingURL=bulk-extend-expiration-date.command.d.ts.map