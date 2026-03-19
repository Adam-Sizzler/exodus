import { z } from 'zod';
export declare namespace GetNodesStatisticsCommand {
    const url: "/api/system/stats/nodes";
    const TSQ_url: "/api/system/stats/nodes";
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
            lastSevenDays: z.ZodArray<z.ZodObject<{
                nodeName: z.ZodString;
                date: z.ZodString;
                totalBytes: z.ZodString;
            }, "strip", z.ZodTypeAny, {
                date: string;
                nodeName: string;
                totalBytes: string;
            }, {
                date: string;
                nodeName: string;
                totalBytes: string;
            }>, "many">;
        }, "strip", z.ZodTypeAny, {
            lastSevenDays: {
                date: string;
                nodeName: string;
                totalBytes: string;
            }[];
        }, {
            lastSevenDays: {
                date: string;
                nodeName: string;
                totalBytes: string;
            }[];
        }>;
    }, "strip", z.ZodTypeAny, {
        response: {
            lastSevenDays: {
                date: string;
                nodeName: string;
                totalBytes: string;
            }[];
        };
    }, {
        response: {
            lastSevenDays: {
                date: string;
                nodeName: string;
                totalBytes: string;
            }[];
        };
    }>;
    type Response = z.infer<typeof ResponseSchema>;
}
//# sourceMappingURL=get-nodes-statistics.d.ts.map