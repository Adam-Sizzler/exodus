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
                }, "strip", z.ZodTypeAny, {
                    enabled: boolean;
                }, {
                    enabled: boolean;
                }>;
                tgAuth: z.ZodObject<{
                    enabled: z.ZodBoolean;
                    botId: z.ZodNullable<z.ZodNumber>;
                }, "strip", z.ZodTypeAny, {
                    enabled: boolean;
                    botId: number | null;
                }, {
                    enabled: boolean;
                    botId: number | null;
                }>;
                oauth2: z.ZodObject<{
                    providers: z.ZodRecord<z.ZodNativeEnum<{
                        readonly GITHUB: "github";
                        readonly POCKETID: "pocketid";
                        readonly YANDEX: "yandex";
                        readonly KEYCLOAK: "keycloak";
                    }>, z.ZodBoolean>;
                }, "strip", z.ZodTypeAny, {
                    providers: Partial<Record<"github" | "pocketid" | "yandex" | "keycloak", boolean>>;
                }, {
                    providers: Partial<Record<"github" | "pocketid" | "yandex" | "keycloak", boolean>>;
                }>;
                password: z.ZodObject<{
                    enabled: z.ZodBoolean;
                }, "strip", z.ZodTypeAny, {
                    enabled: boolean;
                }, {
                    enabled: boolean;
                }>;
            }, "strip", z.ZodTypeAny, {
                passkey: {
                    enabled: boolean;
                };
                tgAuth: {
                    enabled: boolean;
                    botId: number | null;
                };
                oauth2: {
                    providers: Partial<Record<"github" | "pocketid" | "yandex" | "keycloak", boolean>>;
                };
                password: {
                    enabled: boolean;
                };
            }, {
                passkey: {
                    enabled: boolean;
                };
                tgAuth: {
                    enabled: boolean;
                    botId: number | null;
                };
                oauth2: {
                    providers: Partial<Record<"github" | "pocketid" | "yandex" | "keycloak", boolean>>;
                };
                password: {
                    enabled: boolean;
                };
            }>>;
            branding: z.ZodObject<{
                title: z.ZodNullable<z.ZodString>;
                logoUrl: z.ZodNullable<z.ZodString>;
            }, "strip", z.ZodTypeAny, {
                title: string | null;
                logoUrl: string | null;
            }, {
                title: string | null;
                logoUrl: string | null;
            }>;
        }, "strip", z.ZodTypeAny, {
            isLoginAllowed: boolean;
            isRegisterAllowed: boolean;
            authentication: {
                passkey: {
                    enabled: boolean;
                };
                tgAuth: {
                    enabled: boolean;
                    botId: number | null;
                };
                oauth2: {
                    providers: Partial<Record<"github" | "pocketid" | "yandex" | "keycloak", boolean>>;
                };
                password: {
                    enabled: boolean;
                };
            } | null;
            branding: {
                title: string | null;
                logoUrl: string | null;
            };
        }, {
            isLoginAllowed: boolean;
            isRegisterAllowed: boolean;
            authentication: {
                passkey: {
                    enabled: boolean;
                };
                tgAuth: {
                    enabled: boolean;
                    botId: number | null;
                };
                oauth2: {
                    providers: Partial<Record<"github" | "pocketid" | "yandex" | "keycloak", boolean>>;
                };
                password: {
                    enabled: boolean;
                };
            } | null;
            branding: {
                title: string | null;
                logoUrl: string | null;
            };
        }>;
    }, "strip", z.ZodTypeAny, {
        response: {
            isLoginAllowed: boolean;
            isRegisterAllowed: boolean;
            authentication: {
                passkey: {
                    enabled: boolean;
                };
                tgAuth: {
                    enabled: boolean;
                    botId: number | null;
                };
                oauth2: {
                    providers: Partial<Record<"github" | "pocketid" | "yandex" | "keycloak", boolean>>;
                };
                password: {
                    enabled: boolean;
                };
            } | null;
            branding: {
                title: string | null;
                logoUrl: string | null;
            };
        };
    }, {
        response: {
            isLoginAllowed: boolean;
            isRegisterAllowed: boolean;
            authentication: {
                passkey: {
                    enabled: boolean;
                };
                tgAuth: {
                    enabled: boolean;
                    botId: number | null;
                };
                oauth2: {
                    providers: Partial<Record<"github" | "pocketid" | "yandex" | "keycloak", boolean>>;
                };
                password: {
                    enabled: boolean;
                };
            } | null;
            branding: {
                title: string | null;
                logoUrl: string | null;
            };
        };
    }>;
    type Response = z.infer<typeof ResponseSchema>;
}
//# sourceMappingURL=get-status.command.d.ts.map