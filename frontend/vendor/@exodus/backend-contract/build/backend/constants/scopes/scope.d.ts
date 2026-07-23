export declare const SCOPE_ACTION: {
    readonly READ: "read";
    readonly WRITE: "write";
};
export type TScopeAction = (typeof SCOPE_ACTION)[keyof typeof SCOPE_ACTION];
export declare const SCOPE_WILDCARD: "*";
export declare const normalizeControllerUrl: (controllerUrl: string) => string;
export declare const buildResourceScope: (resource: string) => string;
export declare const buildActionScope: (resource: string, action: TScopeAction) => string;
export declare const buildEndpointScope: (resource: string, scopeSlug: string) => string;
//# sourceMappingURL=scope.d.ts.map