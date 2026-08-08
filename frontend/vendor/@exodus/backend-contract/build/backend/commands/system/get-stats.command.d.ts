import { z } from 'zod';
export declare namespace GetStatsCommand {
    const url: "/api/system/stats";
    const TSQ_url: "/api/system/stats";
    const endpointDetails: import("../../constants").EndpointDetails;
    const RequestQuerySchema: z.ZodObject<{
        tz: z.ZodOptional<z.ZodString>;
    }, z.core.$strip>;
    const ResponseSchema: z.ZodObject<{
        response: z.ZodObject<{
            cpu: z.ZodObject<{
                cores: z.ZodNumber;
            }, z.core.$strip>;
            memory: z.ZodObject<{
                total: z.ZodNumber;
                free: z.ZodNumber;
                used: z.ZodNumber;
            }, z.core.$strip>;
            uptime: z.ZodNumber;
            timestamp: z.ZodNumber;
            users: z.ZodObject<{
                statusCounts: z.ZodRecord<z.ZodEnum<{
                    readonly ACTIVE: "ACTIVE";
                    readonly DISABLED: "DISABLED";
                    readonly LIMITED: "LIMITED";
                    readonly EXPIRED: "EXPIRED";
                }>, z.ZodNumber>;
                totalUsers: z.ZodNumber;
            }, z.core.$strip>;
            onlineStats: z.ZodObject<{
                lastDay: z.ZodNumber;
                lastWeek: z.ZodNumber;
                neverOnline: z.ZodNumber;
                onlineNow: z.ZodNumber;
            }, z.core.$strip>;
            nodes: z.ZodObject<{
                totalOnline: z.ZodNumber;
                totalBytesLifetime: z.ZodString;
            }, z.core.$strip>;
        }, z.core.$strip>;
    }, z.core.$strip>;
    type RequestQuery = z.infer<typeof RequestQuerySchema>;
    type Response = z.infer<typeof ResponseSchema>;
}
//# sourceMappingURL=get-stats.command.d.ts.map