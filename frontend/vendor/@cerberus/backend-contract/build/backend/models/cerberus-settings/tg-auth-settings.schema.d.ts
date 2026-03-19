import z from 'zod';
export declare const TgAuthSettingsSchema: z.ZodObject<{
    enabled: z.ZodBoolean;
    botToken: z.ZodNullable<z.ZodString>;
    adminIds: z.ZodArray<z.ZodString, "many">;
}, "strip", z.ZodTypeAny, {
    enabled: boolean;
    botToken: string | null;
    adminIds: string[];
}, {
    enabled: boolean;
    botToken: string | null;
    adminIds: string[];
}>;
export type TTgAuthSettings = z.infer<typeof TgAuthSettingsSchema>;
//# sourceMappingURL=tg-auth-settings.schema.d.ts.map