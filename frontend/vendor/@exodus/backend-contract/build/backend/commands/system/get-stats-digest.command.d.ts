import { z } from 'zod';
export declare namespace GetStatsDigestCommand {
    const url: "/api/system/stats/digest";
    const TSQ_url: "/api/system/stats/digest";
    const endpointDetails: import("../../constants").EndpointDetails;
    const RequestQuerySchema: z.ZodObject<{
        start: z.ZodISODateTime;
        end: z.ZodISODateTime;
    }, z.core.$strip>;
    const ResponseSchema: z.ZodObject<{
        response: z.ZodObject<{
            users: z.ZodObject<{
                createdCount: z.ZodNumber;
                expiredCount: z.ZodNumber;
            }, z.core.$strip>;
            traffic: z.ZodObject<{
                totalBytes: z.ZodString;
                byUsersCreatedInRangeBytes: z.ZodString;
            }, z.core.$strip>;
            hwidDevices: z.ZodObject<{
                createdCount: z.ZodNumber;
            }, z.core.$strip>;
        }, z.core.$strip>;
    }, z.core.$strip>;
    type RequestQuery = z.infer<typeof RequestQuerySchema>;
    type Response = z.infer<typeof ResponseSchema>;
}
//# sourceMappingURL=get-stats-digest.command.d.ts.map