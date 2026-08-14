import { z } from 'zod';
export declare namespace GetNodesTagsCommand {
    const url: "/api/nodes/tags";
    const TSQ_url: "/api/nodes/tags";
    const endpointDetails: import("../../../constants").EndpointDetails;
    const ResponseSchema: z.ZodObject<{
        response: z.ZodObject<{
            tags: z.ZodArray<z.ZodString>;
        }, z.core.$strip>;
    }, z.core.$strip>;
    type Response = z.infer<typeof ResponseSchema>;
}
//# sourceMappingURL=get-nodes-tags.command.d.ts.map