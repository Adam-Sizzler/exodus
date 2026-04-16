import z from 'zod';
export declare const HwidSettingsSchema: z.ZodObject<{
    enabled: z.ZodBoolean;
    fallbackDeviceLimit: z.ZodNumber;
    maxDevicesAnnounce: z.ZodNullable<z.ZodString>;
}, "strip", z.ZodTypeAny, {
    enabled: boolean;
    fallbackDeviceLimit: number;
    maxDevicesAnnounce: string | null;
}, {
    enabled: boolean;
    fallbackDeviceLimit: number;
    maxDevicesAnnounce: string | null;
}>;
export type THwidSettings = z.infer<typeof HwidSettingsSchema>;
//# sourceMappingURL=hwid-settings.schema.d.ts.map