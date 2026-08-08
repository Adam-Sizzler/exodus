import z from 'zod';
export declare const BrandingSettingsSchema: z.ZodObject<{
    title: z.ZodNullable<z.ZodString>;
    logoUrl: z.ZodNullable<z.ZodString>;
}, z.core.$strip>;
export type TBrandingSettings = z.infer<typeof BrandingSettingsSchema>;
//# sourceMappingURL=branding-settings.schema.d.ts.map