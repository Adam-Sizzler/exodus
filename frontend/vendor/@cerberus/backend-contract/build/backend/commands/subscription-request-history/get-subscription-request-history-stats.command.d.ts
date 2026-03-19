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
            }, "strip", z.ZodTypeAny, {
                count: number;
                app: string;
            }, {
                count: number;
                app: string;
            }>, "many">;
            hourlyRequestStats: z.ZodArray<z.ZodObject<{
                dateTime: z.ZodEffects<z.ZodString, Date, string>;
                requestCount: z.ZodNumber;
            }, "strip", z.ZodTypeAny, {
                dateTime: Date;
                requestCount: number;
            }, {
                dateTime: string;
                requestCount: number;
            }>, "many">;
        }, "strip", z.ZodTypeAny, {
            byParsedApp: {
                count: number;
                app: string;
            }[];
            hourlyRequestStats: {
                dateTime: Date;
                requestCount: number;
            }[];
        }, {
            byParsedApp: {
                count: number;
                app: string;
            }[];
            hourlyRequestStats: {
                dateTime: string;
                requestCount: number;
            }[];
        }>;
    }, "strip", z.ZodTypeAny, {
        response: {
            byParsedApp: {
                count: number;
                app: string;
            }[];
            hourlyRequestStats: {
                dateTime: Date;
                requestCount: number;
            }[];
        };
    }, {
        response: {
            byParsedApp: {
                count: number;
                app: string;
            }[];
            hourlyRequestStats: {
                dateTime: string;
                requestCount: number;
            }[];
        };
    }>;
    type Response = z.infer<typeof ResponseSchema>;
}
//# sourceMappingURL=get-subscription-request-history-stats.command.d.ts.map