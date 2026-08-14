export declare const BACKEND_TOOLS_AUTH_COOKIE_NAME = "ex-tools";
export declare const BACKEND_TOOLS_JWT_ISSUER = "Exodus";
export declare const BACKEND_TOOLS_JWT_LIFETIME_HOURS = 2;
export declare const BACKEND_TOOLS_JWT_SCOPES: {
    readonly ACCESS: "access";
    readonly OTT: "ott";
};
export type TBackendToolsJwtScope = (typeof BACKEND_TOOLS_JWT_SCOPES)[keyof typeof BACKEND_TOOLS_JWT_SCOPES];
//# sourceMappingURL=backend-tools.constants.d.ts.map