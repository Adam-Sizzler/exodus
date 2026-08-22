import z from 'zod';
export declare const NodeIntegrationSchema: z.ZodObject<{
    uuid: z.ZodUUID;
    name: z.ZodString;
    description: z.ZodNullable<z.ZodString>;
    config: z.ZodRecord<z.ZodString, z.ZodUnknown>;
}, z.core.$strip>;
//# sourceMappingURL=node-integration.schema.d.ts.map