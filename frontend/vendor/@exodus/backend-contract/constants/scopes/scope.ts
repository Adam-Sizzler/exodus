export const SCOPE_ACTION = {
    READ: 'read',
    WRITE: 'write',
} as const;

export type TScopeAction = (typeof SCOPE_ACTION)[keyof typeof SCOPE_ACTION];

export const SCOPE_WILDCARD = '*' as const;

export const normalizeControllerUrl = (controllerUrl: string): string => {
    const segments = controllerUrl
        .split('/')
        .filter(Boolean)
        .map((segment) => (segment.startsWith(':') ? `{${segment.slice(1)}}` : segment));

    return `/${segments.join('/')}`;
};

export const buildResourceScope = (resource: string): string => `${resource}:${SCOPE_WILDCARD}`;

export const buildActionScope = (resource: string, action: TScopeAction): string =>
    `${resource}:${action}`;

export const buildEndpointScope = (resource: string, scopeSlug: string): string =>
    `${resource}:${scopeSlug}`;
