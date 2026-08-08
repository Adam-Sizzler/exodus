import { z } from 'zod';
export declare namespace GetBandwidthStatsCommand {
    const url: "/api/system/stats/bandwidth";
    const TSQ_url: "/api/system/stats/bandwidth";
    const endpointDetails: import("../../constants").EndpointDetails;
    const RequestQuerySchema: z.ZodObject<{
        tz: z.ZodOptional<z.ZodString>;
    }, z.core.$strip>;
    const ResponseSchema: z.ZodObject<{
        response: z.ZodObject<{
            bandwidthLastTwoDays: z.ZodObject<{
                current: z.ZodString;
                previous: z.ZodString;
                difference: z.ZodString;
            }, z.core.$strip>;
            bandwidthLastSevenDays: z.ZodObject<{
                current: z.ZodString;
                previous: z.ZodString;
                difference: z.ZodString;
            }, z.core.$strip>;
            bandwidthLast30Days: z.ZodObject<{
                current: z.ZodString;
                previous: z.ZodString;
                difference: z.ZodString;
            }, z.core.$strip>;
            bandwidthCalendarMonth: z.ZodObject<{
                current: z.ZodString;
                previous: z.ZodString;
                difference: z.ZodString;
            }, z.core.$strip>;
            bandwidthCurrentYear: z.ZodObject<{
                current: z.ZodString;
                previous: z.ZodString;
                difference: z.ZodString;
            }, z.core.$strip>;
        }, z.core.$strip>;
    }, z.core.$strip>;
    type RequestQuery = z.infer<typeof RequestQuerySchema>;
    type Response = z.infer<typeof ResponseSchema>;
}
//# sourceMappingURL=get-bandwidth-stats.command.d.ts.map