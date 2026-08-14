export const SECURITY_LAYERS = {
    DEFAULT: 'DEFAULT',
    TLS: 'TLS',
    NONE: 'NONE',
} as const;

export type TSecurityLayers = [keyof typeof SECURITY_LAYERS][number];
export const SECURITY_LAYERS_VALUES = Object.values(SECURITY_LAYERS);
