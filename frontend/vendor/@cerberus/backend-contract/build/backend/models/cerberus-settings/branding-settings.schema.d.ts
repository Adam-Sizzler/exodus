import z from 'zod';
export declare const BrandingSettingsSchema: z.ZodObject<{
    title: z.ZodNullable<z.ZodString>;
    logoUrl: z.ZodNullable<z.ZodString>;
}, "strip", z.ZodTypeAny, {
    title: string | null;
    logoUrl: string | null;
}, {
    title: string | null;
    logoUrl: string | null;
}>;
export type TBrandingSettings = z.infer<typeof BrandingSettingsSchema>;
//# sourceMappingURL=branding-settings.schema.d.ts.map