export type TEndpointScopeKind = 'read' | 'write';

export interface EndpointDetails {
    CONTROLLER_URL: string;
    REQUEST_METHOD: 'post' | 'get' | 'put' | 'delete' | 'patch';
    METHOD_DESCRIPTION: string;
    METHOD_LONG_DESCRIPTION?: string;
    SCOPE: string;
    SCOPE_KIND: TEndpointScopeKind;
}

export function getEndpointDetails(
    controllerUrl: string,
    requestMethod: 'post' | 'get' | 'put' | 'delete' | 'patch',
    methodDescription: string,
    scopeOptions: { scope: string; kind: TEndpointScopeKind },
    methodLongDescription?: string,
): EndpointDetails {
    return {
        CONTROLLER_URL: controllerUrl,
        REQUEST_METHOD: requestMethod,
        METHOD_DESCRIPTION: methodDescription,
        METHOD_LONG_DESCRIPTION: methodLongDescription,
        SCOPE: scopeOptions.scope,
        SCOPE_KIND: scopeOptions.kind,
    };
}
