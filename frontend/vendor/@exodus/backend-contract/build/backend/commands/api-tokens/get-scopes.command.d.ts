import { z } from 'zod';
export declare namespace GetApiTokenScopesCommand {
    const url: "/api/tokens/scopes";
    const TSQ_url: "/api/tokens/scopes";
    const endpointDetails: import("../../constants").EndpointDetails;
    const EndpointScopeSchema: z.ZodObject<{
        key: z.ZodString;
        kind: z.ZodEnum<{
            read: "read";
            write: "write";
        }>;
        method: z.ZodString;
        path: z.ZodString;
        description: z.ZodString;
    }, z.core.$strip>;
    const ResourceScopesSchema: z.ZodObject<{
        resource: z.ZodString;
        resourceScopes: z.ZodArray<z.ZodString>;
        endpoints: z.ZodArray<z.ZodObject<{
            key: z.ZodString;
            kind: z.ZodEnum<{
                read: "read";
                write: "write";
            }>;
            method: z.ZodString;
            path: z.ZodString;
            description: z.ZodString;
        }, z.core.$strip>>;
    }, z.core.$strip>;
    const ResponseSchema: z.ZodObject<{
        response: z.ZodObject<{
            wildcard: z.ZodString;
            resources: z.ZodArray<z.ZodObject<{
                resource: z.ZodString;
                resourceScopes: z.ZodArray<z.ZodString>;
                endpoints: z.ZodArray<z.ZodObject<{
                    key: z.ZodString;
                    kind: z.ZodEnum<{
                        read: "read";
                        write: "write";
                    }>;
                    method: z.ZodString;
                    path: z.ZodString;
                    description: z.ZodString;
                }, z.core.$strip>>;
            }, z.core.$strip>>;
        }, z.core.$strip>;
    }, z.core.$strip>;
    type Response = z.infer<typeof ResponseSchema>;
}
//# sourceMappingURL=get-scopes.command.d.ts.map