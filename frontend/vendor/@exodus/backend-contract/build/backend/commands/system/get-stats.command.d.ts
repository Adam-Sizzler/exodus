import { z } from 'zod';
export declare namespace GetStatsCommand {
    const url: "/api/system/stats";
    const TSQ_url: "/api/system/stats";
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
            cpu: z.ZodObject<{
                cores: z.ZodNumber;
                physicalCores: z.ZodNumber;
            }, "strip", z.ZodTypeAny, {
                cores: number;
                physicalCores: number;
            }, {
                cores: number;
                physicalCores: number;
            }>;
            memory: z.ZodObject<{
                total: z.ZodNumber;
                free: z.ZodNumber;
                used: z.ZodNumber;
                active: z.ZodNumber;
                available: z.ZodNumber;
            }, "strip", z.ZodTypeAny, {
                total: number;
                free: number;
                used: number;
                active: number;
                available: number;
            }, {
                total: number;
                free: number;
                used: number;
                active: number;
                available: number;
            }>;
            uptime: z.ZodNumber;
            timestamp: z.ZodNumber;
            users: z.ZodObject<{
                statusCounts: z.ZodRecord<z.ZodEnum<[string, ...string[]]>, z.ZodNumber>;
                totalUsers: z.ZodNumber;
            }, "strip", z.ZodTypeAny, {
                statusCounts: Record<string, number>;
                totalUsers: number;
            }, {
                statusCounts: Record<string, number>;
                totalUsers: number;
            }>;
            onlineStats: z.ZodObject<{
                lastDay: z.ZodNumber;
                lastWeek: z.ZodNumber;
                neverOnline: z.ZodNumber;
                onlineNow: z.ZodNumber;
            }, "strip", z.ZodTypeAny, {
                lastDay: number;
                lastWeek: number;
                neverOnline: number;
                onlineNow: number;
            }, {
                lastDay: number;
                lastWeek: number;
                neverOnline: number;
                onlineNow: number;
            }>;
            nodes: z.ZodObject<{
                totalOnline: z.ZodNumber;
                totalBytesLifetime: z.ZodString;
            }, "strip", z.ZodTypeAny, {
                totalOnline: number;
                totalBytesLifetime: string;
            }, {
                totalOnline: number;
                totalBytesLifetime: string;
            }>;
        }, "strip", z.ZodTypeAny, {
            nodes: {
                totalOnline: number;
                totalBytesLifetime: string;
            };
            users: {
                statusCounts: Record<string, number>;
                totalUsers: number;
            };
            timestamp: number;
            memory: {
                total: number;
                free: number;
                used: number;
                active: number;
                available: number;
            };
            cpu: {
                cores: number;
                physicalCores: number;
            };
            uptime: number;
            onlineStats: {
                lastDay: number;
                lastWeek: number;
                neverOnline: number;
                onlineNow: number;
            };
        }, {
            nodes: {
                totalOnline: number;
                totalBytesLifetime: string;
            };
            users: {
                statusCounts: Record<string, number>;
                totalUsers: number;
            };
            timestamp: number;
            memory: {
                total: number;
                free: number;
                used: number;
                active: number;
                available: number;
            };
            cpu: {
                cores: number;
                physicalCores: number;
            };
            uptime: number;
            onlineStats: {
                lastDay: number;
                lastWeek: number;
                neverOnline: number;
                onlineNow: number;
            };
        }>;
    }, "strip", z.ZodTypeAny, {
        response: {
            nodes: {
                totalOnline: number;
                totalBytesLifetime: string;
            };
            users: {
                statusCounts: Record<string, number>;
                totalUsers: number;
            };
            timestamp: number;
            memory: {
                total: number;
                free: number;
                used: number;
                active: number;
                available: number;
            };
            cpu: {
                cores: number;
                physicalCores: number;
            };
            uptime: number;
            onlineStats: {
                lastDay: number;
                lastWeek: number;
                neverOnline: number;
                onlineNow: number;
            };
        };
    }, {
        response: {
            nodes: {
                totalOnline: number;
                totalBytesLifetime: string;
            };
            users: {
                statusCounts: Record<string, number>;
                totalUsers: number;
            };
            timestamp: number;
            memory: {
                total: number;
                free: number;
                used: number;
                active: number;
                available: number;
            };
            cpu: {
                cores: number;
                physicalCores: number;
            };
            uptime: number;
            onlineStats: {
                lastDay: number;
                lastWeek: number;
                neverOnline: number;
                onlineNow: number;
            };
        };
    }>;
    type Response = z.infer<typeof ResponseSchema>;
}
//# sourceMappingURL=get-stats.command.d.ts.map