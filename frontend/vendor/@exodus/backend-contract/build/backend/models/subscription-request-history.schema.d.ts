import { z } from 'zod';
export declare const SubscriptionRequestHistorySchema: z.ZodObject<{
    id: z.ZodNumber;
    userId: z.ZodNumber;
    srrResponseType: z.ZodString;
    srrRuleName: z.ZodNullable<z.ZodString>;
    requestIp: z.ZodNullable<z.ZodString>;
    userAgent: z.ZodNullable<z.ZodString>;
    requestAt: z.ZodPipe<z.ZodISODateTime, z.ZodTransform<Date, string>>;
}, z.core.$strip>;
//# sourceMappingURL=subscription-request-history.schema.d.ts.map