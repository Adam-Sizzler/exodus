export interface EndpointDetails {
    CONTROLLER_URL: string;
    REQUEST_METHOD: 'post' | 'get' | 'put' | 'delete' | 'patch';
    METHOD_DESCRIPTION: string;
    METHOD_LONG_DESCRIPTION?: string;
}
export declare function getEndpointDetails(controllerUrl: string, requestMethod: 'post' | 'get' | 'put' | 'delete' | 'patch', methodDescription: string, methodLongDescription?: string): EndpointDetails;
//# sourceMappingURL=endpoint-details.d.ts.map