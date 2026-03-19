import { z } from 'zod';
export declare const InfraBillingHistoryRecordSchema: z.ZodObject<{
    uuid: z.ZodString;
    providerUuid: z.ZodString;
    amount: z.ZodNumber;
    billedAt: z.ZodEffects<z.ZodString, Date, string>;
    provider: z.ZodObject<Omit<{
        uuid: z.ZodString;
        name: z.ZodString;
        faviconLink: z.ZodNullable<z.ZodString>;
        loginUrl: z.ZodNullable<z.ZodString>;
        createdAt: z.ZodEffects<z.ZodString, Date, string>;
        updatedAt: z.ZodEffects<z.ZodString, Date, string>;
    }, "createdAt" | "updatedAt" | "loginUrl">, "strip", z.ZodTypeAny, {
        uuid: string;
        name: string;
        faviconLink: string | null;
    }, {
        uuid: string;
        name: string;
        faviconLink: string | null;
    }>;
}, "strip", z.ZodTypeAny, {
    uuid: string;
    provider: {
        uuid: string;
        name: string;
        faviconLink: string | null;
    };
    providerUuid: string;
    amount: number;
    billedAt: Date;
}, {
    uuid: string;
    provider: {
        uuid: string;
        name: string;
        faviconLink: string | null;
    };
    providerUuid: string;
    amount: number;
    billedAt: string;
}>;
//# sourceMappingURL=infra-billing-history-record.schema.d.ts.map