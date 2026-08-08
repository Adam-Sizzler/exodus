import { z } from 'zod';
export declare namespace GetInternalSquadUsageCommand {
    const url: (uuid: string) => string;
    const TSQ_url: string;
    const endpointDetails: import("../../../constants").EndpointDetails;
    const RequestParamSchema: z.ZodObject<{
        uuid: z.ZodUUID;
    }, z.core.$strip>;
    const RequestQuerySchema: z.ZodObject<{
        start: z.ZodISODate;
        end: z.ZodISODate;
        minTotalBytes: z.ZodDefault<z.ZodOptional<z.ZodCoercedNumber<unknown>>>;
        limit: z.ZodDefault<z.ZodOptional<z.ZodCoercedNumber<unknown>>>;
        cursor: z.ZodOptional<z.ZodCoercedNumber<unknown>>;
    }, z.core.$strip>;
    const ResponseSchema: z.ZodObject<{
        response: z.ZodObject<{
            squadUuid: z.ZodUUID;
            users: z.ZodArray<z.ZodObject<{
                id: z.ZodNumber;
                totalBytes: z.ZodNumber;
            }, z.core.$strip>>;
            nextCursor: z.ZodNullable<z.ZodString>;
            hasMore: z.ZodBoolean;
        }, z.core.$strip>;
    }, z.core.$strip>;
    type RequestParam = z.infer<typeof RequestParamSchema>;
    type RequestQuery = z.infer<typeof RequestQuerySchema>;
    type Response = z.infer<typeof ResponseSchema>;
}
//# sourceMappingURL=get-internal-squad-usage.command.d.ts.map