import z from 'zod';
export declare const SubscriptionPageConfigSchema: z.ZodObject<{
    uuid: z.ZodString;
    viewPosition: z.ZodNumber;
    name: z.ZodString;
    config: z.ZodNullable<z.ZodUnknown>;
}, "strip", z.ZodTypeAny, {
    uuid: string;
    name: string;
    viewPosition: number;
    config?: unknown;
}, {
    uuid: string;
    name: string;
    viewPosition: number;
    config?: unknown;
}>;
//# sourceMappingURL=subscripion-page-config.schema.d.ts.map