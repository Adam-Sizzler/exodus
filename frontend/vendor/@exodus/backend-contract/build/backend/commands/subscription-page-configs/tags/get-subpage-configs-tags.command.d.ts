import { z } from 'zod';
export declare namespace GetSubpageConfigsTagsCommand {
    const url: "/api/subscription-page-configs/tags";
    const TSQ_url: "/api/subscription-page-configs/tags";
    const endpointDetails: import("../../../constants").EndpointDetails;
    const ResponseSchema: z.ZodObject<{
        response: z.ZodObject<{
            tags: z.ZodArray<z.ZodString>;
        }, z.core.$strip>;
    }, z.core.$strip>;
    type Response = z.infer<typeof ResponseSchema>;
}
//# sourceMappingURL=get-subpage-configs-tags.command.d.ts.map