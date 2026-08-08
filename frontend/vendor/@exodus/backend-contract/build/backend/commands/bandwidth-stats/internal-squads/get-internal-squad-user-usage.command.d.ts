import { z } from 'zod';
export declare namespace GetInternalSquadUserUsageCommand {
    const url: (squadUuid: string, userId: string) => string;
    const TSQ_url: string;
    const endpointDetails: import("../../../constants").EndpointDetails;
    const RequestParamSchema: z.ZodObject<{
        squadUuid: z.ZodUUID;
        userId: z.ZodCoercedNumber<unknown>;
    }, z.core.$strip>;
    const RequestQuerySchema: z.ZodObject<{
        start: z.ZodISODate;
        end: z.ZodISODate;
    }, z.core.$strip>;
    const ResponseSchema: z.ZodObject<{
        response: z.ZodObject<{
            days: z.ZodArray<z.ZodObject<{
                date: z.ZodString;
                nodes: z.ZodArray<z.ZodObject<{
                    uuid: z.ZodUUID;
                    totalBytes: z.ZodNumber;
                }, z.core.$strip>>;
            }, z.core.$strip>>;
        }, z.core.$strip>;
    }, z.core.$strip>;
    type RequestParam = z.infer<typeof RequestParamSchema>;
    type RequestQuery = z.infer<typeof RequestQuerySchema>;
    type Response = z.infer<typeof ResponseSchema>;
}
//# sourceMappingURL=get-internal-squad-user-usage.command.d.ts.map