import z from 'zod';
export declare const PasskeySettingsSchema: z.ZodObject<{
    enabled: z.ZodBoolean;
    rpId: z.ZodNullable<z.ZodEffects<z.ZodString, string, string>>;
    origin: z.ZodNullable<z.ZodEffects<z.ZodString, string, string>>;
}, "strip", z.ZodTypeAny, {
    enabled: boolean;
    rpId: string | null;
    origin: string | null;
}, {
    enabled: boolean;
    rpId: string | null;
    origin: string | null;
}>;
export type TRemnawavePasskeySettings = z.infer<typeof PasskeySettingsSchema>;
//# sourceMappingURL=passkey-settings.schema.d.ts.map