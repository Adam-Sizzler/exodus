import { z } from 'zod';
export declare namespace GetConfigProfilesTagsCommand {
    const url: "/api/config-profiles/tags";
    const TSQ_url: "/api/config-profiles/tags";
    const endpointDetails: import("../../../constants").EndpointDetails;
    const ResponseSchema: z.ZodObject<{
        response: z.ZodObject<{
            tags: z.ZodArray<z.ZodString>;
        }, z.core.$strip>;
    }, z.core.$strip>;
    type Response = z.infer<typeof ResponseSchema>;
}
//# sourceMappingURL=get-config-profiles-tags.command.d.ts.map