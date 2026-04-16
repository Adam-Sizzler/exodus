import { z } from 'zod';
export declare namespace GetAllNodesTagsCommand {
    const url: "/api/nodes/tags";
    const TSQ_url: "/api/nodes/tags";
    const endpointDetails: import("../../../constants").EndpointDetails;
    const ResponseSchema: z.ZodObject<{
        response: z.ZodObject<{
            tags: z.ZodArray<z.ZodString, "many">;
        }, "strip", z.ZodTypeAny, {
            tags: string[];
        }, {
            tags: string[];
        }>;
    }, "strip", z.ZodTypeAny, {
        response: {
            tags: string[];
        };
    }, {
        response: {
            tags: string[];
        };
    }>;
    type Response = z.infer<typeof ResponseSchema>;
}
//# sourceMappingURL=get-all-nodes-tags.command.d.ts.map