import { z } from 'zod';
export declare namespace OAuth2CallbackCommand {
    const url: "/api/auth/oauth2/callback";
    const TSQ_url: "/api/auth/oauth2/callback";
    const endpointDetails: import("../../../constants").EndpointDetails;
    const RequestBodySchema: z.ZodObject<{
        provider: z.ZodEnum<{
            readonly TELEGRAM: "telegram";
            readonly GITHUB: "github";
            readonly POCKETID: "pocketid";
            readonly YANDEX: "yandex";
            readonly KEYCLOAK: "keycloak";
            readonly GENERIC: "generic";
        }>;
        code: z.ZodString;
        state: z.ZodString;
    }, z.core.$strip>;
    const ResponseSchema: z.ZodObject<{
        response: z.ZodObject<{
            accessToken: z.ZodString;
        }, z.core.$strip>;
    }, z.core.$strip>;
    type RequestBody = z.infer<typeof RequestBodySchema>;
    type Response = z.infer<typeof ResponseSchema>;
}
//# sourceMappingURL=callback.command.d.ts.map