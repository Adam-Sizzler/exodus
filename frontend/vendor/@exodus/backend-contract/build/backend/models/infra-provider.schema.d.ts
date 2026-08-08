import { z } from 'zod';
export declare const InfraProviderSchema: z.ZodObject<{
    uuid: z.ZodUUID;
    name: z.ZodString;
    faviconLink: z.ZodNullable<z.ZodString>;
    loginUrl: z.ZodNullable<z.ZodString>;
    createdAt: z.ZodPipe<z.ZodISODateTime, z.ZodTransform<Date, string>>;
    updatedAt: z.ZodPipe<z.ZodISODateTime, z.ZodTransform<Date, string>>;
    billingHistory: z.ZodObject<{
        totalAmount: z.ZodNumber;
        totalBills: z.ZodNumber;
    }, z.core.$strip>;
    billingNodes: z.ZodArray<z.ZodObject<{
        name: z.ZodString;
        details: z.ZodNullable<z.ZodObject<{
            nodeUuid: z.ZodUUID;
            countryCode: z.ZodString;
        }, z.core.$strip>>;
    }, z.core.$strip>>;
}, z.core.$strip>;
export declare const PartialInfraProviderSchema: z.ZodObject<{
    uuid: z.ZodUUID;
    name: z.ZodString;
    faviconLink: z.ZodNullable<z.ZodString>;
    loginUrl: z.ZodNullable<z.ZodString>;
    createdAt: z.ZodPipe<z.ZodISODateTime, z.ZodTransform<Date, string>>;
    updatedAt: z.ZodPipe<z.ZodISODateTime, z.ZodTransform<Date, string>>;
}, z.core.$strip>;
//# sourceMappingURL=infra-provider.schema.d.ts.map