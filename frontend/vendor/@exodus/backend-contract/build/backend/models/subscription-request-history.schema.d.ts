import { z } from 'zod';
export declare const SubscriptionRequestHistorySchema: z.ZodObject<{
    id: z.ZodNumber;
    userId: z.ZodNumber;
    requestIp: z.ZodNullable<z.ZodString>;
    userAgent: z.ZodNullable<z.ZodString>;
    requestAt: z.ZodEffects<z.ZodString, Date, string>;
}, "strip", z.ZodTypeAny, {
    id: number;
    userId: number;
    userAgent: string | null;
    requestIp: string | null;
    requestAt: Date;
}, {
    id: number;
    userId: number;
    userAgent: string | null;
    requestIp: string | null;
    requestAt: string;
}>;
//# sourceMappingURL=subscription-request-history.schema.d.ts.map