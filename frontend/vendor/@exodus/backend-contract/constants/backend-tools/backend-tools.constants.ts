export const BACKEND_TOOLS_AUTH_COOKIE_NAME = 'ex-tools';
export const BACKEND_TOOLS_JWT_ISSUER = 'Exodus';
export const BACKEND_TOOLS_JWT_LIFETIME_HOURS = 2;
export const BACKEND_TOOLS_JWT_SCOPES = {
    ACCESS: 'access',
    OTT: 'ott',
} as const;
export type TBackendToolsJwtScope =
    (typeof BACKEND_TOOLS_JWT_SCOPES)[keyof typeof BACKEND_TOOLS_JWT_SCOPES];
