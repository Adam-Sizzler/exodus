import z from 'zod';
export declare const Oauth2SettingsSchema: z.ZodObject<{
    github: z.ZodObject<{
        enabled: z.ZodBoolean;
        clientId: z.ZodNullable<z.ZodString>;
        clientSecret: z.ZodNullable<z.ZodString>;
        allowedEmails: z.ZodArray<z.ZodString, "many">;
    }, "strip", z.ZodTypeAny, {
        enabled: boolean;
        clientId: string | null;
        clientSecret: string | null;
        allowedEmails: string[];
    }, {
        enabled: boolean;
        clientId: string | null;
        clientSecret: string | null;
        allowedEmails: string[];
    }>;
    pocketid: z.ZodObject<{
        enabled: z.ZodBoolean;
        clientId: z.ZodNullable<z.ZodString>;
        clientSecret: z.ZodNullable<z.ZodString>;
        plainDomain: z.ZodNullable<z.ZodEffects<z.ZodString, string, string>>;
        allowedEmails: z.ZodArray<z.ZodString, "many">;
    }, "strip", z.ZodTypeAny, {
        enabled: boolean;
        clientId: string | null;
        clientSecret: string | null;
        allowedEmails: string[];
        plainDomain: string | null;
    }, {
        enabled: boolean;
        clientId: string | null;
        clientSecret: string | null;
        allowedEmails: string[];
        plainDomain: string | null;
    }>;
    yandex: z.ZodObject<{
        enabled: z.ZodBoolean;
        clientId: z.ZodNullable<z.ZodString>;
        clientSecret: z.ZodNullable<z.ZodString>;
        allowedEmails: z.ZodArray<z.ZodString, "many">;
    }, "strip", z.ZodTypeAny, {
        enabled: boolean;
        clientId: string | null;
        clientSecret: string | null;
        allowedEmails: string[];
    }, {
        enabled: boolean;
        clientId: string | null;
        clientSecret: string | null;
        allowedEmails: string[];
    }>;
    keycloak: z.ZodDefault<z.ZodObject<{
        enabled: z.ZodBoolean;
        realm: z.ZodNullable<z.ZodString>;
        clientId: z.ZodNullable<z.ZodString>;
        clientSecret: z.ZodNullable<z.ZodString>;
        frontendDomain: z.ZodNullable<z.ZodEffects<z.ZodString, string, string>>;
        keycloakDomain: z.ZodNullable<z.ZodEffects<z.ZodString, string, string>>;
        allowedEmails: z.ZodArray<z.ZodString, "many">;
    }, "strip", z.ZodTypeAny, {
        enabled: boolean;
        clientId: string | null;
        clientSecret: string | null;
        allowedEmails: string[];
        realm: string | null;
        frontendDomain: string | null;
        keycloakDomain: string | null;
    }, {
        enabled: boolean;
        clientId: string | null;
        clientSecret: string | null;
        allowedEmails: string[];
        realm: string | null;
        frontendDomain: string | null;
        keycloakDomain: string | null;
    }>>;
    generic: z.ZodDefault<z.ZodObject<{
        enabled: z.ZodBoolean;
        clientId: z.ZodNullable<z.ZodString>;
        clientSecret: z.ZodNullable<z.ZodString>;
        withPkce: z.ZodBoolean;
        authorizationUrl: z.ZodNullable<z.ZodString>;
        tokenUrl: z.ZodNullable<z.ZodString>;
        frontendDomain: z.ZodNullable<z.ZodEffects<z.ZodString, string, string>>;
        allowedEmails: z.ZodArray<z.ZodString, "many">;
    }, "strip", z.ZodTypeAny, {
        enabled: boolean;
        authorizationUrl: string | null;
        clientId: string | null;
        clientSecret: string | null;
        allowedEmails: string[];
        frontendDomain: string | null;
        withPkce: boolean;
        tokenUrl: string | null;
    }, {
        enabled: boolean;
        authorizationUrl: string | null;
        clientId: string | null;
        clientSecret: string | null;
        allowedEmails: string[];
        frontendDomain: string | null;
        withPkce: boolean;
        tokenUrl: string | null;
    }>>;
    telegram: z.ZodDefault<z.ZodObject<{
        enabled: z.ZodBoolean;
        clientId: z.ZodNullable<z.ZodString>;
        clientSecret: z.ZodNullable<z.ZodString>;
        allowedIds: z.ZodArray<z.ZodString, "many">;
        frontendDomain: z.ZodNullable<z.ZodEffects<z.ZodString, string, string>>;
    }, "strip", z.ZodTypeAny, {
        enabled: boolean;
        clientId: string | null;
        clientSecret: string | null;
        frontendDomain: string | null;
        allowedIds: string[];
    }, {
        enabled: boolean;
        clientId: string | null;
        clientSecret: string | null;
        frontendDomain: string | null;
        allowedIds: string[];
    }>>;
}, "strip", z.ZodTypeAny, {
    telegram: {
        enabled: boolean;
        clientId: string | null;
        clientSecret: string | null;
        frontendDomain: string | null;
        allowedIds: string[];
    };
    github: {
        enabled: boolean;
        clientId: string | null;
        clientSecret: string | null;
        allowedEmails: string[];
    };
    pocketid: {
        enabled: boolean;
        clientId: string | null;
        clientSecret: string | null;
        allowedEmails: string[];
        plainDomain: string | null;
    };
    yandex: {
        enabled: boolean;
        clientId: string | null;
        clientSecret: string | null;
        allowedEmails: string[];
    };
    keycloak: {
        enabled: boolean;
        clientId: string | null;
        clientSecret: string | null;
        allowedEmails: string[];
        realm: string | null;
        frontendDomain: string | null;
        keycloakDomain: string | null;
    };
    generic: {
        enabled: boolean;
        authorizationUrl: string | null;
        clientId: string | null;
        clientSecret: string | null;
        allowedEmails: string[];
        frontendDomain: string | null;
        withPkce: boolean;
        tokenUrl: string | null;
    };
}, {
    github: {
        enabled: boolean;
        clientId: string | null;
        clientSecret: string | null;
        allowedEmails: string[];
    };
    pocketid: {
        enabled: boolean;
        clientId: string | null;
        clientSecret: string | null;
        allowedEmails: string[];
        plainDomain: string | null;
    };
    yandex: {
        enabled: boolean;
        clientId: string | null;
        clientSecret: string | null;
        allowedEmails: string[];
    };
    telegram?: {
        enabled: boolean;
        clientId: string | null;
        clientSecret: string | null;
        frontendDomain: string | null;
        allowedIds: string[];
    } | undefined;
    keycloak?: {
        enabled: boolean;
        clientId: string | null;
        clientSecret: string | null;
        allowedEmails: string[];
        realm: string | null;
        frontendDomain: string | null;
        keycloakDomain: string | null;
    } | undefined;
    generic?: {
        enabled: boolean;
        authorizationUrl: string | null;
        clientId: string | null;
        clientSecret: string | null;
        allowedEmails: string[];
        frontendDomain: string | null;
        withPkce: boolean;
        tokenUrl: string | null;
    } | undefined;
}>;
export type TOauth2Settings = z.infer<typeof Oauth2SettingsSchema>;
//# sourceMappingURL=oauth2-settings.schema.d.ts.map