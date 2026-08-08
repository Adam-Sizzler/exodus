import { z } from 'zod';
export declare const InfraBillingNodeSchema: z.ZodObject<{
    uuid: z.ZodUUID;
    nodeUuid: z.ZodNullable<z.ZodUUID>;
    name: z.ZodNullable<z.ZodString>;
    providerUuid: z.ZodUUID;
    provider: z.ZodObject<{
        uuid: z.ZodUUID;
        name: z.ZodString;
        faviconLink: z.ZodNullable<z.ZodString>;
        loginUrl: z.ZodNullable<z.ZodString>;
    }, z.core.$strip>;
    node: z.ZodNullable<z.ZodObject<{
        uuid: z.ZodUUID;
        name: z.ZodString;
        countryCode: z.ZodString;
    }, z.core.$strip>>;
    nextBillingAt: z.ZodPipe<z.ZodISODateTime, z.ZodTransform<Date, string>>;
    createdAt: z.ZodPipe<z.ZodISODateTime, z.ZodTransform<Date, string>>;
    updatedAt: z.ZodPipe<z.ZodISODateTime, z.ZodTransform<Date, string>>;
}, z.core.$strip>;
//# sourceMappingURL=infra-billing-node.schema.d.ts.map