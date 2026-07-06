import { z } from 'zod';
export declare namespace GetUserSubscriptionRequestHistoryCommand {
    const url: (uuid: string) => string;
    const TSQ_url: string;
    const endpointDetails: import("../../constants").EndpointDetails;
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
            total: z.ZodNumber;
            records: z.ZodArray<z.ZodObject<{
                id: z.ZodNumber;
                userId: z.ZodNumber;
                requestAt: z.ZodEffects<z.ZodString, Date, string>;
                requestIp: z.ZodNullable<z.ZodOptional<z.ZodString>>;
                userAgent: z.ZodNullable<z.ZodOptional<z.ZodString>>;
            }, "strip", z.ZodTypeAny, {
                id: number;
                userId: number;
                requestAt: Date;
                userAgent?: string | null | undefined;
                requestIp?: string | null | undefined;
            }, {
                id: number;
                userId: number;
                requestAt: string;
                userAgent?: string | null | undefined;
                requestIp?: string | null | undefined;
            }>, "many">;
        }, "strip", z.ZodTypeAny, {
            total: number;
            records: {
                id: number;
                userId: number;
                requestAt: Date;
                userAgent?: string | null | undefined;
                requestIp?: string | null | undefined;
            }[];
        }, {
            total: number;
            records: {
                id: number;
                userId: number;
                requestAt: string;
                userAgent?: string | null | undefined;
                requestIp?: string | null | undefined;
            }[];
        }>;
    }, "strip", z.ZodTypeAny, {
        response: {
            total: number;
            records: {
                id: number;
                userId: number;
                requestAt: Date;
                userAgent?: string | null | undefined;
                requestIp?: string | null | undefined;
            }[];
        };
    }, {
        response: {
            total: number;
            records: {
                id: number;
                userId: number;
                requestAt: string;
                userAgent?: string | null | undefined;
                requestIp?: string | null | undefined;
            }[];
        };
    }>;
    type Response = z.infer<typeof ResponseSchema>;
}
//# sourceMappingURL=get-user-subscription-request-history.command.d.ts.map