import { z } from 'zod';
export declare namespace GetAllTagsCommand {
    const url: "/api/users/tags";
    const TSQ_url: "/api/users/tags";
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
//# sourceMappingURL=get-all-tags.command.d.ts.map