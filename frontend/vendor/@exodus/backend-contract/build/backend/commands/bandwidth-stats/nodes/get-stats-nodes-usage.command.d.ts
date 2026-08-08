import { z } from 'zod';
export declare namespace GetStatsNodesUsageCommand {
    const url: "/api/bandwidth-stats/nodes/";
    const TSQ_url: "/api/bandwidth-stats/nodes/";
    const endpointDetails: import("../../../constants").EndpointDetails;
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
    type RequestQuery = z.infer<typeof RequestQuerySchema>;
    type Response = z.infer<typeof ResponseSchema>;
}
//# sourceMappingURL=get-stats-nodes-usage.command.d.ts.map