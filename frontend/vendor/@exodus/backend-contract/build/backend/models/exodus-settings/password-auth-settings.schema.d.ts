import z from 'zod';
export declare const PasswordAuthSettingsSchema: z.ZodObject<{
    enabled: z.ZodBoolean;
}, z.core.$strip>;
export type TPasswordAuthSettings = z.infer<typeof PasswordAuthSettingsSchema>;
//# sourceMappingURL=password-auth-settings.schema.d.ts.map