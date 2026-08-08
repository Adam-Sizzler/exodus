import z from 'zod';
export declare const PasskeySettingsSchema: z.ZodObject<{
    enabled: z.ZodBoolean;
    rpId: z.ZodNullable<z.ZodString>;
    origin: z.ZodNullable<z.ZodString>;
}, z.core.$strip>;
export type TExodusPasskeySettings = z.infer<typeof PasskeySettingsSchema>;
//# sourceMappingURL=passkey-settings.schema.d.ts.map