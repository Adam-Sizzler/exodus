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
                    ips: z.ZodArray<z.ZodObject<{
                        ip: z.ZodString;
                        lastSeen: z.ZodEffects<z.ZodString, Date, string>;
                    }, "strip", z.ZodTypeAny, {
                        ip: string;
                        lastSeen: Date;
                    }, {
                        ip: string;
                        lastSeen: string;
                    }>, "many">;
                }, "strip", z.ZodTypeAny, {
                    nodeUuid: string;
                    nodeName: string;
                    countryCode: string;
                    ips: {
                        ip: string;
                        lastSeen: Date;
                    }[];
                }, {
                    nodeUuid: string;
                    nodeName: string;
                    countryCode: string;
                    ips: {
                        ip: string;
                        lastSeen: string;
                    }[];
                }>, "many">;
            }, "strip", z.ZodTypeAny, {
                nodes: {
                    nodeUuid: string;
                    nodeName: string;
                    countryCode: string;
                    ips: {
                        ip: string;
                        lastSeen: Date;
                    }[];
                }[];
                userUuid: string;
                userId: string;
                success: boolean;
            }, {
                nodes: {
                    nodeUuid: string;
                    nodeName: string;
                    countryCode: string;
                    ips: {
                        ip: string;
                        lastSeen: string;
                    }[];
                }[];
                userUuid: string;
                userId: string;
                success: boolean;
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
                    ips: {
                        ip: string;
                        lastSeen: Date;
                    }[];
                }[];
                userUuid: string;
                userId: string;
                success: boolean;
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
                    ips: {
                        ip: string;
                        lastSeen: string;
                    }[];
                }[];
                userUuid: string;
                userId: string;
                success: boolean;
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
                    ips: {
                        ip: string;
                        lastSeen: Date;
                    }[];
                }[];
                userUuid: string;
                userId: string;
                success: boolean;
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
                    ips: {
                        ip: string;
                        lastSeen: string;
                    }[];
                }[];
                userUuid: string;
                userId: string;
                success: boolean;
            } | null;
        };
    }>;
    type Response = z.infer<typeof ResponseSchema>;
}
//# sourceMappingURL=fetch-ips-result.command.d.ts.map