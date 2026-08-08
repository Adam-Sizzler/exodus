import { z } from 'zod';
export declare namespace GetHostsTagsCommand {
    const url: "/api/hosts/tags";
    const TSQ_url: "/api/hosts/tags";
    const endpointDetails: import("../../../constants").EndpointDetails;
    const ResponseSchema: z.ZodObject<{
        response: z.ZodObject<{
            tags: z.ZodArray<z.ZodString>;
        }, z.core.$strip>;
    }, z.core.$strip>;
    type Response = z.infer<typeof ResponseSchema>;
}
//# sourceMappingURL=get-hosts-tags.command.d.ts.map