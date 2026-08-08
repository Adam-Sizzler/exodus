import { z } from 'zod';
export declare namespace GetStatsUserUsageCommand {
    const url: (userId: string) => string;
    const TSQ_url: string;
    const endpointDetails: import("../../../constants").EndpointDetails;
    const RequestParamSchema: z.ZodObject<{
        userId: z.ZodCoercedNumber<unknown>;
    }, z.core.$strip>;
    const RequestQuerySchema: z.ZodObject<{
        start: z.ZodISODate;
        end: z.ZodISODate;
        topNodesLimit: z.ZodDefault<z.ZodCoercedNumber<unknown>>;
    }, z.core.$strip>;
    const ResponseSchema: z.ZodObject<{
        response: z.ZodObject<{
            categories: z.ZodArray<z.ZodString>;
            sparklineData: z.ZodArray<z.ZodNumber>;
            topNodes: z.ZodArray<z.ZodObject<{
                uuid: z.ZodUUID;
                color: z.ZodString;
                name: z.ZodString;
                countryCode: z.ZodString;
                total: z.ZodNumber;
            }, z.core.$strip>>;
            series: z.ZodArray<z.ZodObject<{
                uuid: z.ZodUUID;
                name: z.ZodString;
                color: z.ZodString;
                countryCode: z.ZodString;
                total: z.ZodNumber;
                data: z.ZodArray<z.ZodNumber>;
            }, z.core.$strip>>;
        }, z.core.$strip>;
    }, z.core.$strip>;
    type RequestParam = z.infer<typeof RequestParamSchema>;
    type RequestQuery = z.infer<typeof RequestQuerySchema>;
    type Response = z.infer<typeof ResponseSchema>;
}
//# sourceMappingURL=get-stats-user-usage.command.d.ts.map