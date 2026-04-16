import { z } from 'zod';
export declare namespace FetchIpsResultCommand {
    const url: (jobId: string) => string;
    const TSQ_url: string;
    const endpointDetails: import("../../constants").EndpointDetails;
    const RequestSchema: z.ZodObject<{
        jobId: z.ZodString;
    }, "strip", z.ZodTypeAny, {
        jobId: string;
    }, {
        jobId: string;
    }>;
    type Request = z.infer<typeof RequestSchema>;
    const ResponseSchema: z.ZodObject<{
        response: z.ZodObject<{
            isCompleted: z.ZodBoolean;
            isFailed: z.ZodBoolean;
            progress: z.ZodObject<{
                total: z.ZodNumber;
                completed: z.ZodNumber;
                percent: z.ZodNumber;
            }, "strip", z.ZodTypeAny, {
                total: number;
                completed: number;
                percent: number;
            }, {
                total: number;
                completed: number;
                percent: number;
            }>;
            result: z.ZodNullable<z.ZodObject<{
                success: z.ZodBoolean;
                userUuid: z.ZodString;
                userId: z.ZodString;
                nodes: z.ZodArray<z.ZodObject<{
                    nodeUuid: z.ZodString;
                    nodeName: z.ZodString;
                    countryCode: z.ZodString;
                    ips: z.ZodArray<z.ZodString, "many">;
                }, "strip", z.ZodTypeAny, {
                    nodeUuid: string;
                    nodeName: string;
                    countryCode: string;
                    ips: string[];
                }, {
                    nodeUuid: string;
                    nodeName: string;
                    countryCode: string;
                    ips: string[];
                }>, "many">;
            }, "strip", z.ZodTypeAny, {
                nodes: {
                    nodeUuid: string;
                    nodeName: string;
                    countryCode: string;
                    ips: string[];
                }[];
                userUuid: string;
                success: boolean;
                userId: string;
            }, {
                nodes: {
                    nodeUuid: string;
                    nodeName: string;
                    countryCode: string;
                    ips: string[];
                }[];
                userUuid: string;
                success: boolean;
                userId: string;
            }>>;
        }, "strip", z.ZodTypeAny, {
            isCompleted: boolean;
            isFailed: boolean;
            progress: {
                total: number;
                completed: number;
                percent: number;
            };
            result: {
                nodes: {
                    nodeUuid: string;
                    nodeName: string;
                    countryCode: string;
                    ips: string[];
                }[];
                userUuid: string;
                success: boolean;
                userId: string;
            } | null;
        }, {
            isCompleted: boolean;
            isFailed: boolean;
            progress: {
                total: number;
                completed: number;
                percent: number;
            };
            result: {
                nodes: {
                    nodeUuid: string;
                    nodeName: string;
                    countryCode: string;
                    ips: string[];
                }[];
                userUuid: string;
                success: boolean;
                userId: string;
            } | null;
        }>;
    }, "strip", z.ZodTypeAny, {
        response: {
            isCompleted: boolean;
            isFailed: boolean;
            progress: {
                total: number;
                completed: number;
                percent: number;
            };
            result: {
                nodes: {
                    nodeUuid: string;
                    nodeName: string;
                    countryCode: string;
                    ips: string[];
                }[];
                userUuid: string;
                success: boolean;
                userId: string;
            } | null;
        };
    }, {
        response: {
            isCompleted: boolean;
            isFailed: boolean;
            progress: {
                total: number;
                completed: number;
                percent: number;
            };
            result: {
                nodes: {
                    nodeUuid: string;
                    nodeName: string;
                    countryCode: string;
                    ips: string[];
                }[];
                userUuid: string;
                success: boolean;
                userId: string;
            } | null;
        };
    }>;
    type Response = z.infer<typeof ResponseSchema>;
}
//# sourceMappingURL=fetch-ips-result.command.d.ts.map