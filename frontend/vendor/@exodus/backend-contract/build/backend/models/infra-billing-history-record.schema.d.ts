import { z } from 'zod';
export declare const InfraBillingHistoryRecordSchema: z.ZodObject<{
    uuid: z.ZodUUID;
    providerUuid: z.ZodUUID;
    amount: z.ZodNumber;
    billedAt: z.ZodPipe<z.ZodISODateTime, z.ZodTransform<Date, string>>;
    provider: z.ZodObject<{
        uuid: z.ZodUUID;
        name: z.ZodString;
        faviconLink: z.ZodNullable<z.ZodString>;
    }, z.core.$strip>;
}, z.core.$strip>;
//# sourceMappingURL=infra-billing-history-record.schema.d.ts.map