import { z } from 'zod';
export declare const SubscriptionRequestHistorySchema: z.ZodObject<{
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
}>;
//# sourceMappingURL=subscription-request-history.schema.d.ts.map