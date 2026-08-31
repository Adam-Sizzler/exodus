import { z } from 'zod';
export declare namespace SetSubpageConfigTagsCommand {
    const url: "/api/subscription-page-configs/tags";
    const TSQ_url: "/api/subscription-page-configs/tags";
    const endpointDetails: import("../../../constants").EndpointDetails;
    const RequestBodySchema: z.ZodObject<{
        uuid: z.ZodUUID;
        tags: z.ZodArray<z.ZodString>;
    }, z.core.$strip>;
    const ResponseSchema: z.ZodObject<{
        response: z.ZodObject<{
            uuid: z.ZodUUID;
            tags: z.ZodArray<z.ZodString>;
        }, z.core.$strip>;
    }, z.core.$strip>;
    type RequestBody = z.infer<typeof RequestBodySchema>;
    type Response = z.infer<typeof ResponseSchema>;
}
//# sourceMappingURL=set-subpage-config-tags.command.d.ts.map