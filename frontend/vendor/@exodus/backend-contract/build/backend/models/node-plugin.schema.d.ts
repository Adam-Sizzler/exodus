import z from 'zod';
export declare const NodePluginSchema: z.ZodObject<{
    uuid: z.ZodUUID;
    viewPosition: z.ZodNumber;
    name: z.ZodString;
    pluginConfig: z.ZodNullable<z.ZodUnknown>;
}, z.core.$strip>;
//# sourceMappingURL=node-plugin.schema.d.ts.map