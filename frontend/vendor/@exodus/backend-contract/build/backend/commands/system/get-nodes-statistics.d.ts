import { z } from 'zod';
export declare namespace GetNodesStatisticsCommand {
    const url: "/api/system/stats/nodes";
    const TSQ_url: "/api/system/stats/nodes";
    const endpointDetails: import("../../constants").EndpointDetails;
    const RequestQuerySchema: z.ZodObject<{
        tz: z.ZodOptional<z.ZodString>;
    }, z.core.$strip>;
    const ResponseSchema: z.ZodObject<{
        response: z.ZodObject<{
            lastSevenDays: z.ZodArray<z.ZodObject<{
                nodeName: z.ZodString;
                date: z.ZodString;
                totalBytes: z.ZodString;
            }, z.core.$strip>>;
        }, z.core.$strip>;
    }, z.core.$strip>;
    type RequestQuery = z.infer<typeof RequestQuerySchema>;
    type Response = z.infer<typeof ResponseSchema>;
}
//# sourceMappingURL=get-nodes-statistics.d.ts.map