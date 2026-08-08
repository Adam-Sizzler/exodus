import { z } from 'zod';
export declare namespace OAuth2AuthorizeCommand {
    const url: "/api/auth/oauth2/authorize";
    const TSQ_url: "/api/auth/oauth2/authorize";
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
    }, z.core.$strip>;
    const ResponseSchema: z.ZodObject<{
        response: z.ZodObject<{
            authorizationUrl: z.ZodNullable<z.ZodURL>;
        }, z.core.$strip>;
    }, z.core.$strip>;
    type RequestBody = z.infer<typeof RequestBodySchema>;
    type Response = z.infer<typeof ResponseSchema>;
}
//# sourceMappingURL=authorize.command.d.ts.map