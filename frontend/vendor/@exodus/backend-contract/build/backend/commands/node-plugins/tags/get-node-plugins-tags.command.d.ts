import { z } from 'zod';
export declare namespace GetNodePluginsTagsCommand {
    const url: "/api/node-plugins/tags";
    const TSQ_url: "/api/node-plugins/tags";
    const endpointDetails: import("../../../constants").EndpointDetails;
    const ResponseSchema: z.ZodObject<{
        response: z.ZodObject<{
            tags: z.ZodArray<z.ZodString>;
        }, z.core.$strip>;
    }, z.core.$strip>;
    type Response = z.infer<typeof ResponseSchema>;
}
//# sourceMappingURL=get-node-plugins-tags.command.d.ts.map