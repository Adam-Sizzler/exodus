import z from 'zod';
export declare const NodePluginSchema: z.ZodObject<{
    uuid: z.ZodString;
    viewPosition: z.ZodNumber;
    name: z.ZodString;
    pluginConfig: z.ZodNullable<z.ZodUnknown>;
}, "strip", z.ZodTypeAny, {
    uuid: string;
    name: string;
    viewPosition: number;
    pluginConfig?: unknown;
}, {
    uuid: string;
    name: string;
    viewPosition: number;
    pluginConfig?: unknown;
}>;
//# sourceMappingURL=node-plugin.schema.d.ts.map