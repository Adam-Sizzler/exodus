import z from 'zod';
export declare const SubscriptionPageConfigSchema: z.ZodObject<{
    uuid: z.ZodUUID;
    viewPosition: z.ZodNumber;
    name: z.ZodString;
    tags: z.ZodDefault<z.ZodArray<z.ZodString>>;
    config: z.ZodNullable<z.ZodUnknown>;
}, z.core.$strip>;
//# sourceMappingURL=subscripion-page-config.schema.d.ts.map