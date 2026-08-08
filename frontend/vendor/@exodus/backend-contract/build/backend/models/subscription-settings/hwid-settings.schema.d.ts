import z from 'zod';
export declare const HwidSettingsSchema: z.ZodObject<{
    enabled: z.ZodBoolean;
    fallbackDeviceLimit: z.ZodNumber;
    maxDevicesAnnounce: z.ZodNullable<z.ZodString>;
}, z.core.$strip>;
export type THwidSettings = z.infer<typeof HwidSettingsSchema>;
//# sourceMappingURL=hwid-settings.schema.d.ts.map