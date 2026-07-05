import { z } from 'zod';
export declare namespace FetchUsersIpsResultCommand {
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
            result: z.ZodNullable<z.ZodObject<{
                success: z.ZodBoolean;
                nodeUuid: z.ZodString;
                users: z.ZodArray<z.ZodObject<{
                    userId: z.ZodString;
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
                    userId: string;
                    ips: {
                        ip: string;
                        lastSeen: Date;
                    }[];
                }, {
                    userId: string;
                    ips: {
                        ip: string;
                        lastSeen: string;
                    }[];
                }>, "many">;
            }, "strip", z.ZodTypeAny, {
                users: {
                    userId: string;
                    ips: {
                        ip: string;
                        lastSeen: Date;
                    }[];
                }[];
                nodeUuid: string;
                success: boolean;
            }, {
                users: {
                    userId: string;
                    ips: {
                        ip: string;
                        lastSeen: string;
                    }[];
                }[];
                nodeUuid: string;
                success: boolean;
            }>>;
        }, "strip", z.ZodTypeAny, {
            isCompleted: boolean;
            isFailed: boolean;
            result: {
                users: {
                    userId: string;
                    ips: {
                        ip: string;
                        lastSeen: Date;
                    }[];
                }[];
                nodeUuid: string;
                success: boolean;
            } | null;
        }, {
            isCompleted: boolean;
            isFailed: boolean;
            result: {
                users: {
                    userId: string;
                    ips: {
                        ip: string;
                        lastSeen: string;
                    }[];
                }[];
                nodeUuid: string;
                success: boolean;
            } | null;
        }>;
    }, "strip", z.ZodTypeAny, {
        response: {
            isCompleted: boolean;
            isFailed: boolean;
            result: {
                users: {
                    userId: string;
                    ips: {
                        ip: string;
                        lastSeen: Date;
                    }[];
                }[];
                nodeUuid: string;
                success: boolean;
            } | null;
        };
    }, {
        response: {
            isCompleted: boolean;
            isFailed: boolean;
            result: {
                users: {
                    userId: string;
                    ips: {
                        ip: string;
                        lastSeen: string;
                    }[];
                }[];
                nodeUuid: string;
                success: boolean;
            } | null;
        };
    }>;
    type Response = z.infer<typeof ResponseSchema>;
}
//# sourceMappingURL=fetch-users-ips-result.command.d.ts.map