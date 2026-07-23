export declare const OAUTH2_PROVIDERS: {
    readonly TELEGRAM: "telegram";
    readonly GITHUB: "github";
    readonly POCKETID: "pocketid";
    readonly YANDEX: "yandex";
    readonly KEYCLOAK: "keycloak";
    readonly GENERIC: "generic";
};
export type TOAuth2Providers = [keyof typeof OAUTH2_PROVIDERS][number];
export type TOAuth2ProvidersKeys = (typeof OAUTH2_PROVIDERS)[keyof typeof OAUTH2_PROVIDERS];
export declare const OAUTH2_PROVIDERS_VALUES: ("telegram" | "github" | "pocketid" | "yandex" | "keycloak" | "generic")[];
//# sourceMappingURL=providers.contants.d.ts.map