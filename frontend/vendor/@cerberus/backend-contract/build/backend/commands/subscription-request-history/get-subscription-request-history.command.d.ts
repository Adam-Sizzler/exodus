import { z } from 'zod';
export declare namespace GetSubscriptionRequestHistoryCommand {
    const url: "/api/subscription-request-history/";
    const TSQ_url: "/api/subscription-request-history/";
    const endpointDetails: import("../../constants").EndpointDetails;
    const RequestQuerySchema: z.ZodObject<{
        start: z.ZodDefault<z.ZodNumber>;
        size: z.ZodDefault<z.ZodNumber>;
        filters: z.ZodOptional<z.ZodEffects<z.ZodArray<z.ZodObject<{
            id: z.ZodString;
            value: z.ZodUnknown;
        }, "strip", z.ZodTypeAny, {
            id: string;
            value?: unknown;
        }, {
            id: string;
            value?: unknown;
        }>, "many">, {
            id: string;
            value?: unknown;
        }[], unknown>>;
        filterModes: z.ZodOptional<z.ZodEffects<z.ZodRecord<z.ZodString, z.ZodString>, Record<string, string>, unknown>>;
        globalFilterMode: z.ZodOptional<z.ZodString>;
        sorting: z.ZodOptional<z.ZodEffects<z.ZodArray<z.ZodObject<{
            id: z.ZodString;
            desc: z.ZodBoolean;
        }, "strip", z.ZodTypeAny, {
            id: string;
            desc: boolean;
        }, {
            id: string;
            desc: boolean;
        }>, "many">, {
            id: string;
            desc: boolean;
        }[], unknown>>;
    }, "strip", z.ZodTypeAny, {
        start: number;
        size: number;
        filters?: {
            id: string;
            value?: unknown;
        }[] | undefined;
        filterModes?: Record<string, string> | undefined;
        globalFilterMode?: string | undefined;
        sorting?: {
            id: string;
            desc: boolean;
        }[] | undefined;
    }, {
        start?: number | undefined;
        size?: number | undefined;
        filters?: unknown;
        filterModes?: unknown;
        globalFilterMode?: string | undefined;
        sorting?: unknown;
    }>;
    type RequestQuery = z.infer<typeof RequestQuerySchema>;
    const ResponseSchema: z.ZodObject<{
        response: z.ZodObject<{
            records: z.ZodArray<z.ZodObject<{
                id: z.ZodNumber;
                userUuid: z.ZodString;
                requestIp: z.ZodNullable<z.ZodString>;
                userAgent: z.ZodNullable<z.ZodString>;
                requestAt: z.ZodEffects<z.ZodString, Date, string>;
            }, "strip", z.ZodTypeAny, {
                id: number;
                userUuid: string;
                userAgent: string | null;
                requestIp: string | null;
                requestAt: Date;
            }, {
                id: number;
                userUuid: string;
                userAgent: string | null;
                requestIp: string | null;
                requestAt: string;
            }>, "many">;
            total: z.ZodNumber;
        }, "strip", z.ZodTypeAny, {
            total: number;
            records: {
                id: number;
                userUuid: string;
                userAgent: string | null;
                requestIp: string | null;
                requestAt: Date;
            }[];
        }, {
            total: number;
            records: {
                id: number;
                userUuid: string;
                userAgent: string | null;
                requestIp: string | null;
                requestAt: string;
            }[];
        }>;
    }, "strip", z.ZodTypeAny, {
        response: {
            total: number;
            records: {
                id: number;
                userUuid: string;
                userAgent: string | null;
                requestIp: string | null;
                requestAt: Date;
            }[];
        };
    }, {
        response: {
            total: number;
            records: {
                id: number;
                userUuid: string;
                userAgent: string | null;
                requestIp: string | null;
                requestAt: string;
            }[];
        };
    }>;
    type Response = z.infer<typeof ResponseSchema>;
}
//# sourceMappingURL=get-subscription-request-history.command.d.ts.map