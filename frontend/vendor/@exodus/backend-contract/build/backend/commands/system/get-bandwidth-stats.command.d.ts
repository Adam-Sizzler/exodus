import { z } from 'zod';
export declare namespace GetBandwidthStatsCommand {
    const url: "/api/system/stats/bandwidth";
    const TSQ_url: "/api/system/stats/bandwidth";
    const endpointDetails: import("../../constants").EndpointDetails;
    const RequestQuerySchema: z.ZodObject<{
        tz: z.ZodOptional<z.ZodString>;
    }, "strip", z.ZodTypeAny, {
        tz?: string | undefined;
    }, {
        tz?: string | undefined;
    }>;
    type Request = z.infer<typeof RequestQuerySchema>;
    const ResponseSchema: z.ZodObject<{
        response: z.ZodObject<{
            bandwidthLastTwoDays: z.ZodObject<{
                current: z.ZodString;
                previous: z.ZodString;
                difference: z.ZodString;
            }, "strip", z.ZodTypeAny, {
                current: string;
                previous: string;
                difference: string;
            }, {
                current: string;
                previous: string;
                difference: string;
            }>;
            bandwidthLastSevenDays: z.ZodObject<{
                current: z.ZodString;
                previous: z.ZodString;
                difference: z.ZodString;
            }, "strip", z.ZodTypeAny, {
                current: string;
                previous: string;
                difference: string;
            }, {
                current: string;
                previous: string;
                difference: string;
            }>;
            bandwidthLast30Days: z.ZodObject<{
                current: z.ZodString;
                previous: z.ZodString;
                difference: z.ZodString;
            }, "strip", z.ZodTypeAny, {
                current: string;
                previous: string;
                difference: string;
            }, {
                current: string;
                previous: string;
                difference: string;
            }>;
            bandwidthCalendarMonth: z.ZodObject<{
                current: z.ZodString;
                previous: z.ZodString;
                difference: z.ZodString;
            }, "strip", z.ZodTypeAny, {
                current: string;
                previous: string;
                difference: string;
            }, {
                current: string;
                previous: string;
                difference: string;
            }>;
            bandwidthCurrentYear: z.ZodObject<{
                current: z.ZodString;
                previous: z.ZodString;
                difference: z.ZodString;
            }, "strip", z.ZodTypeAny, {
                current: string;
                previous: string;
                difference: string;
            }, {
                current: string;
                previous: string;
                difference: string;
            }>;
        }, "strip", z.ZodTypeAny, {
            bandwidthLastTwoDays: {
                current: string;
                previous: string;
                difference: string;
            };
            bandwidthLastSevenDays: {
                current: string;
                previous: string;
                difference: string;
            };
            bandwidthLast30Days: {
                current: string;
                previous: string;
                difference: string;
            };
            bandwidthCalendarMonth: {
                current: string;
                previous: string;
                difference: string;
            };
            bandwidthCurrentYear: {
                current: string;
                previous: string;
                difference: string;
            };
        }, {
            bandwidthLastTwoDays: {
                current: string;
                previous: string;
                difference: string;
            };
            bandwidthLastSevenDays: {
                current: string;
                previous: string;
                difference: string;
            };
            bandwidthLast30Days: {
                current: string;
                previous: string;
                difference: string;
            };
            bandwidthCalendarMonth: {
                current: string;
                previous: string;
                difference: string;
            };
            bandwidthCurrentYear: {
                current: string;
                previous: string;
                difference: string;
            };
        }>;
    }, "strip", z.ZodTypeAny, {
        response: {
            bandwidthLastTwoDays: {
                current: string;
                previous: string;
                difference: string;
            };
            bandwidthLastSevenDays: {
                current: string;
                previous: string;
                difference: string;
            };
            bandwidthLast30Days: {
                current: string;
                previous: string;
                difference: string;
            };
            bandwidthCalendarMonth: {
                current: string;
                previous: string;
                difference: string;
            };
            bandwidthCurrentYear: {
                current: string;
                previous: string;
                difference: string;
            };
        };
    }, {
        response: {
            bandwidthLastTwoDays: {
                current: string;
                previous: string;
                difference: string;
            };
            bandwidthLastSevenDays: {
                current: string;
                previous: string;
                difference: string;
            };
            bandwidthLast30Days: {
                current: string;
                previous: string;
                difference: string;
            };
            bandwidthCalendarMonth: {
                current: string;
                previous: string;
                difference: string;
            };
            bandwidthCurrentYear: {
                current: string;
                previous: string;
                difference: string;
            };
        };
    }>;
    type Response = z.infer<typeof ResponseSchema>;
}
//# sourceMappingURL=get-bandwidth-stats.command.d.ts.map