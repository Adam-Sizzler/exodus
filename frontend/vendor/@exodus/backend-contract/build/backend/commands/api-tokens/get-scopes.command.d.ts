import { z } from 'zod';
export declare namespace GetApiTokenScopesCommand {
    const url: "/api/tokens/scopes";
    const TSQ_url: "/api/tokens/scopes";
    const endpointDetails: import("../../constants").EndpointDetails;
    const EndpointScopeSchema: z.ZodObject<{
        key: z.ZodString;
        kind: z.ZodEnum<["read", "write"]>;
        method: z.ZodString;
        path: z.ZodString;
        description: z.ZodString;
    }, "strip", z.ZodTypeAny, {
        path: string;
        key: string;
        kind: "read" | "write";
        method: string;
        description: string;
    }, {
        path: string;
        key: string;
        kind: "read" | "write";
        method: string;
        description: string;
    }>;
    const ResourceScopesSchema: z.ZodObject<{
        resource: z.ZodString;
        resourceScopes: z.ZodArray<z.ZodString, "many">;
        endpoints: z.ZodArray<z.ZodObject<{
            key: z.ZodString;
            kind: z.ZodEnum<["read", "write"]>;
            method: z.ZodString;
            path: z.ZodString;
            description: z.ZodString;
        }, "strip", z.ZodTypeAny, {
            path: string;
            key: string;
            kind: "read" | "write";
            method: string;
            description: string;
        }, {
            path: string;
            key: string;
            kind: "read" | "write";
            method: string;
            description: string;
        }>, "many">;
    }, "strip", z.ZodTypeAny, {
        resource: string;
        resourceScopes: string[];
        endpoints: {
            path: string;
            key: string;
            kind: "read" | "write";
            method: string;
            description: string;
        }[];
    }, {
        resource: string;
        resourceScopes: string[];
        endpoints: {
            path: string;
            key: string;
            kind: "read" | "write";
            method: string;
            description: string;
        }[];
    }>;
    const ResponseSchema: z.ZodObject<{
        response: z.ZodObject<{
            wildcard: z.ZodString;
            resources: z.ZodArray<z.ZodObject<{
                resource: z.ZodString;
                resourceScopes: z.ZodArray<z.ZodString, "many">;
                endpoints: z.ZodArray<z.ZodObject<{
                    key: z.ZodString;
                    kind: z.ZodEnum<["read", "write"]>;
                    method: z.ZodString;
                    path: z.ZodString;
                    description: z.ZodString;
                }, "strip", z.ZodTypeAny, {
                    path: string;
                    key: string;
                    kind: "read" | "write";
                    method: string;
                    description: string;
                }, {
                    path: string;
                    key: string;
                    kind: "read" | "write";
                    method: string;
                    description: string;
                }>, "many">;
            }, "strip", z.ZodTypeAny, {
                resource: string;
                resourceScopes: string[];
                endpoints: {
                    path: string;
                    key: string;
                    kind: "read" | "write";
                    method: string;
                    description: string;
                }[];
            }, {
                resource: string;
                resourceScopes: string[];
                endpoints: {
                    path: string;
                    key: string;
                    kind: "read" | "write";
                    method: string;
                    description: string;
                }[];
            }>, "many">;
        }, "strip", z.ZodTypeAny, {
            wildcard: string;
            resources: {
                resource: string;
                resourceScopes: string[];
                endpoints: {
                    path: string;
                    key: string;
                    kind: "read" | "write";
                    method: string;
                    description: string;
                }[];
            }[];
        }, {
            wildcard: string;
            resources: {
                resource: string;
                resourceScopes: string[];
                endpoints: {
                    path: string;
                    key: string;
                    kind: "read" | "write";
                    method: string;
                    description: string;
                }[];
            }[];
        }>;
    }, "strip", z.ZodTypeAny, {
        response: {
            wildcard: string;
            resources: {
                resource: string;
                resourceScopes: string[];
                endpoints: {
                    path: string;
                    key: string;
                    kind: "read" | "write";
                    method: string;
                    description: string;
                }[];
            }[];
        };
    }, {
        response: {
            wildcard: string;
            resources: {
                resource: string;
                resourceScopes: string[];
                endpoints: {
                    path: string;
                    key: string;
                    kind: "read" | "write";
                    method: string;
                    description: string;
                }[];
            }[];
        };
    }>;
    type Response = z.infer<typeof ResponseSchema>;
}
//# sourceMappingURL=get-scopes.command.d.ts.map