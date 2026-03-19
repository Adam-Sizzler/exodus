import { z } from 'zod';
export declare namespace UpdateCerberusSettingsCommand {
    const url: "/api/cerberus-settings/";
    const TSQ_url: "/api/cerberus-settings/";
    const endpointDetails: import("../../constants").EndpointDetails;
    const RequestSchema: z.ZodObject<{
        passkeySettings: z.ZodOptional<z.ZodObject<{
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
        }>>;
        oauth2Settings: z.ZodOptional<z.ZodObject<{
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
        }, "strip", z.ZodTypeAny, {
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
        }>>;
        tgAuthSettings: z.ZodOptional<z.ZodObject<{
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
        }>>;
        passwordSettings: z.ZodOptional<z.ZodObject<{
            enabled: z.ZodBoolean;
        }, "strip", z.ZodTypeAny, {
            enabled: boolean;
        }, {
            enabled: boolean;
        }>>;
        brandingSettings: z.ZodOptional<z.ZodObject<{
            title: z.ZodNullable<z.ZodString>;
            logoUrl: z.ZodNullable<z.ZodString>;
        }, "strip", z.ZodTypeAny, {
            title: string | null;
            logoUrl: string | null;
        }, {
            title: string | null;
            logoUrl: string | null;
        }>>;
    }, "strip", z.ZodTypeAny, {
        passkeySettings?: {
            enabled: boolean;
            rpId: string | null;
            origin: string | null;
        } | undefined;
        oauth2Settings?: {
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
        } | undefined;
        tgAuthSettings?: {
            enabled: boolean;
            botToken: string | null;
            adminIds: string[];
        } | undefined;
        passwordSettings?: {
            enabled: boolean;
        } | undefined;
        brandingSettings?: {
            title: string | null;
            logoUrl: string | null;
        } | undefined;
    }, {
        passkeySettings?: {
            enabled: boolean;
            rpId: string | null;
            origin: string | null;
        } | undefined;
        oauth2Settings?: {
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
        } | undefined;
        tgAuthSettings?: {
            enabled: boolean;
            botToken: string | null;
            adminIds: string[];
        } | undefined;
        passwordSettings?: {
            enabled: boolean;
        } | undefined;
        brandingSettings?: {
            title: string | null;
            logoUrl: string | null;
        } | undefined;
    }>;
    type Request = z.infer<typeof RequestSchema>;
    const ResponseSchema: z.ZodObject<{
        response: z.ZodObject<{
            passkeySettings: z.ZodNullable<z.ZodObject<{
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
            }>>;
            oauth2Settings: z.ZodNullable<z.ZodObject<{
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
            }, "strip", z.ZodTypeAny, {
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
            }>>;
            tgAuthSettings: z.ZodNullable<z.ZodObject<{
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
            }>>;
            passwordSettings: z.ZodNullable<z.ZodObject<{
                enabled: z.ZodBoolean;
            }, "strip", z.ZodTypeAny, {
                enabled: boolean;
            }, {
                enabled: boolean;
            }>>;
            brandingSettings: z.ZodNullable<z.ZodObject<{
                title: z.ZodNullable<z.ZodString>;
                logoUrl: z.ZodNullable<z.ZodString>;
            }, "strip", z.ZodTypeAny, {
                title: string | null;
                logoUrl: string | null;
            }, {
                title: string | null;
                logoUrl: string | null;
            }>>;
        }, "strip", z.ZodTypeAny, {
            passkeySettings: {
                enabled: boolean;
                rpId: string | null;
                origin: string | null;
            } | null;
            oauth2Settings: {
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
            } | null;
            tgAuthSettings: {
                enabled: boolean;
                botToken: string | null;
                adminIds: string[];
            } | null;
            passwordSettings: {
                enabled: boolean;
            } | null;
            brandingSettings: {
                title: string | null;
                logoUrl: string | null;
            } | null;
        }, {
            passkeySettings: {
                enabled: boolean;
                rpId: string | null;
                origin: string | null;
            } | null;
            oauth2Settings: {
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
            } | null;
            tgAuthSettings: {
                enabled: boolean;
                botToken: string | null;
                adminIds: string[];
            } | null;
            passwordSettings: {
                enabled: boolean;
            } | null;
            brandingSettings: {
                title: string | null;
                logoUrl: string | null;
            } | null;
        }>;
    }, "strip", z.ZodTypeAny, {
        response: {
            passkeySettings: {
                enabled: boolean;
                rpId: string | null;
                origin: string | null;
            } | null;
            oauth2Settings: {
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
            } | null;
            tgAuthSettings: {
                enabled: boolean;
                botToken: string | null;
                adminIds: string[];
            } | null;
            passwordSettings: {
                enabled: boolean;
            } | null;
            brandingSettings: {
                title: string | null;
                logoUrl: string | null;
            } | null;
        };
    }, {
        response: {
            passkeySettings: {
                enabled: boolean;
                rpId: string | null;
                origin: string | null;
            } | null;
            oauth2Settings: {
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
            } | null;
            tgAuthSettings: {
                enabled: boolean;
                botToken: string | null;
                adminIds: string[];
            } | null;
            passwordSettings: {
                enabled: boolean;
            } | null;
            brandingSettings: {
                title: string | null;
                logoUrl: string | null;
            } | null;
        };
    }>;
    type Response = z.infer<typeof ResponseSchema>;
}
//# sourceMappingURL=update-cerberus-settings.command.d.ts.map