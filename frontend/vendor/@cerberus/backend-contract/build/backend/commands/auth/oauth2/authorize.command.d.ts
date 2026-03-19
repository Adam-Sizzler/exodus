import { z } from 'zod';
export declare namespace OAuth2AuthorizeCommand {
    const url: "/api/auth/oauth2/authorize";
    const TSQ_url: "/api/auth/oauth2/authorize";
    const endpointDetails: import("../../../constants").EndpointDetails;
    const RequestSchema: z.ZodObject<{
        provider: z.ZodNativeEnum<{
            readonly GITHUB: "github";
            readonly POCKETID: "pocketid";
            readonly YANDEX: "yandex";
            readonly KEYCLOAK: "keycloak";
        }>;
    }, "strip", z.ZodTypeAny, {
        provider: "github" | "pocketid" | "yandex" | "keycloak";
    }, {
        provider: "github" | "pocketid" | "yandex" | "keycloak";
    }>;
    const ResponseSchema: z.ZodObject<{
        response: z.ZodObject<{
            authorizationUrl: z.ZodNullable<z.ZodString>;
        }, "strip", z.ZodTypeAny, {
            authorizationUrl: string | null;
        }, {
            authorizationUrl: string | null;
        }>;
    }, "strip", z.ZodTypeAny, {
        response: {
            authorizationUrl: string | null;
        };
    }, {
        response: {
            authorizationUrl: string | null;
        };
    }>;
    type Response = z.infer<typeof ResponseSchema>;
}
//# sourceMappingURL=authorize.command.d.ts.map