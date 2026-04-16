import { z } from 'zod';
export declare namespace GetStatsNodeUsersUsageCommand {
    const url: (uuid: string) => string;
    const TSQ_url: string;
    const endpointDetails: import("../../../constants").EndpointDetails;
    const RequestQuerySchema: z.ZodObject<{
        start: z.ZodString;
        end: z.ZodString;
        topUsersLimit: z.ZodDefault<z.ZodNumber>;
    }, "strip", z.ZodTypeAny, {
        start: string;
        end: string;
        topUsersLimit: number;
    }, {
        start: string;
        end: string;
        topUsersLimit?: number | undefined;
    }>;
    type RequestQuery = z.infer<typeof RequestQuerySchema>;
    const RequestSchema: z.ZodObject<{
        uuid: z.ZodString;
    }, "strip", z.ZodTypeAny, {
        uuid: string;
    }, {
        uuid: string;
    }>;
    type Request = z.infer<typeof RequestSchema>;
    const ResponseSchema: z.ZodObject<{
        response: z.ZodObject<{
            categories: z.ZodArray<z.ZodString, "many">;
            sparklineData: z.ZodArray<z.ZodNumber, "many">;
            topUsers: z.ZodArray<z.ZodObject<{
                color: z.ZodString;
                username: z.ZodString;
                total: z.ZodNumber;
            }, "strip", z.ZodTypeAny, {
                username: string;
                total: number;
                color: string;
            }, {
                username: string;
                total: number;
                color: string;
            }>, "many">;
        }, "strip", z.ZodTypeAny, {
            categories: string[];
            sparklineData: number[];
            topUsers: {
                username: string;
                total: number;
                color: string;
            }[];
        }, {
            categories: string[];
            sparklineData: number[];
            topUsers: {
                username: string;
                total: number;
                color: string;
            }[];
        }>;
    }, "strip", z.ZodTypeAny, {
        response: {
            categories: string[];
            sparklineData: number[];
            topUsers: {
                username: string;
                total: number;
                color: string;
            }[];
        };
    }, {
        response: {
            categories: string[];
            sparklineData: number[];
            topUsers: {
                username: string;
                total: number;
                color: string;
            }[];
        };
    }>;
    type Response = z.infer<typeof ResponseSchema>;
}
//# sourceMappingURL=get-stats-node-users-usage.command.d.ts.map