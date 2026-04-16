export declare const OAUTH2_PROVIDERS: {
    readonly GITHUB: "github";
    readonly POCKETID: "pocketid";
    readonly YANDEX: "yandex";
    readonly KEYCLOAK: "keycloak";
};
export type TOAuth2Providers = [keyof typeof OAUTH2_PROVIDERS][number];
export type TOAuth2ProvidersKeys = (typeof OAUTH2_PROVIDERS)[keyof typeof OAUTH2_PROVIDERS];
export declare const OAUTH2_PROVIDERS_VALUES: ("github" | "pocketid" | "yandex" | "keycloak")[];
//# sourceMappingURL=providers.contants.d.ts.map