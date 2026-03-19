import { z } from 'zod';
export declare namespace OAuth2CallbackCommand {
    const url: "/api/auth/oauth2/callback";
    const TSQ_url: "/api/auth/oauth2/callback";
    const endpointDetails: import("../../../constants").EndpointDetails;
    const RequestSchema: z.ZodObject<{
        provider: z.ZodNativeEnum<{
            readonly GITHUB: "github";
            readonly POCKETID: "pocketid";
            readonly YANDEX: "yandex";
            readonly KEYCLOAK: "keycloak";
        }>;
        code: z.ZodString;
        state: z.ZodString;
    }, "strip", z.ZodTypeAny, {
        code: string;
        provider: "github" | "pocketid" | "yandex" | "keycloak";
        state: string;
    }, {
        code: string;
        provider: "github" | "pocketid" | "yandex" | "keycloak";
        state: string;
    }>;
    type Request = z.infer<typeof RequestSchema>;
    const ResponseSchema: z.ZodObject<{
        response: z.ZodObject<{
            accessToken: z.ZodString;
        }, "strip", z.ZodTypeAny, {
            accessToken: string;
        }, {
            accessToken: string;
        }>;
    }, "strip", z.ZodTypeAny, {
        response: {
            accessToken: string;
        };
    }, {
        response: {
            accessToken: string;
        };
    }>;
    type Response = z.infer<typeof ResponseSchema>;
}
//# sourceMappingURL=callback.command.d.ts.map