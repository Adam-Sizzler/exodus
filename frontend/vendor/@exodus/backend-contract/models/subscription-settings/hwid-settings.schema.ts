import z from 'zod';

export const HwidSettingsSchema = z.object({
    enabled: z.boolean(),
    fallbackDeviceLimit: z.number(),
    maxDevicesAnnounce: z.nullable(z.string().max(200)),
});

export type THwidSettings = z.infer<typeof HwidSettingsSchema>;
