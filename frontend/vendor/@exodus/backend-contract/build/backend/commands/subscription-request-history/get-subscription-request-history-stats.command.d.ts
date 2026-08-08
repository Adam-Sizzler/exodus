import { z } from 'zod';
export declare namespace GetSubscriptionRequestHistoryStatsCommand {
    const url: "/api/subscription-request-history/stats";
    const TSQ_url: "/api/subscription-request-history/stats";
    const endpointDetails: import("../../constants").EndpointDetails;
    const ResponseSchema: z.ZodObject<{
        response: z.ZodObject<{
            byParsedApp: z.ZodArray<z.ZodObject<{
                app: z.ZodString;
                count: z.ZodNumber;
            }, z.core.$strip>>;
            hourlyRequestStats: z.ZodArray<z.ZodObject<{
                dateTime: z.ZodPipe<z.ZodISODateTime, z.ZodTransform<Date, string>>;
                requestCount: z.ZodNumber;
            }, z.core.$strip>>;
        }, z.core.$strip>;
    }, z.core.$strip>;
    type Response = z.infer<typeof ResponseSchema>;
}
//# sourceMappingURL=get-subscription-request-history-stats.command.d.ts.map