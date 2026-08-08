import { z } from 'zod';
export declare namespace GetStatusCommand {
    const url: "/api/auth/status";
    const TSQ_url: "/api/auth/status";
    const endpointDetails: import("../../constants").EndpointDetails;
    const ResponseSchema: z.ZodObject<{
        response: z.ZodObject<{
            isLoginAllowed: z.ZodBoolean;
            isRegisterAllowed: z.ZodBoolean;
            authentication: z.ZodNullable<z.ZodObject<{
                passkey: z.ZodObject<{
                    enabled: z.ZodBoolean;
                }, z.core.$strip>;
                oauth2: z.ZodObject<{
                    providers: z.ZodRecord<z.ZodEnum<{
                        readonly TELEGRAM: "telegram";
                        readonly GITHUB: "github";
                        readonly POCKETID: "pocketid";
                        readonly YANDEX: "yandex";
                        readonly KEYCLOAK: "keycloak";
                        readonly GENERIC: "generic";
                    }>, z.ZodBoolean>;
                }, z.core.$strip>;
                password: z.ZodObject<{
                    enabled: z.ZodBoolean;
                }, z.core.$strip>;
            }, z.core.$strip>>;
            branding: z.ZodObject<{
                title: z.ZodNullable<z.ZodString>;
                logoUrl: z.ZodNullable<z.ZodString>;
            }, z.core.$strip>;
        }, z.core.$strip>;
    }, z.core.$strip>;
    type Response = z.infer<typeof ResponseSchema>;
}
//# sourceMappingURL=get-status.command.d.ts.map