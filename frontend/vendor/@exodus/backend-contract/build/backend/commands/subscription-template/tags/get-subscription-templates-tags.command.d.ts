import { z } from 'zod';
export declare namespace GetSubscriptionTemplatesTagsCommand {
    const url: "/api/subscription-templates/tags";
    const TSQ_url: "/api/subscription-templates/tags";
    const endpointDetails: import("../../../constants").EndpointDetails;
    const ResponseSchema: z.ZodObject<{
        response: z.ZodObject<{
            tags: z.ZodArray<z.ZodString>;
        }, z.core.$strip>;
    }, z.core.$strip>;
    type Response = z.infer<typeof ResponseSchema>;
}
//# sourceMappingURL=get-subscription-templates-tags.command.d.ts.map