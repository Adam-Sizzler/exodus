import { z } from 'zod';
export declare namespace GetExternalSquadsTagsCommand {
    const url: "/api/external-squads/tags";
    const TSQ_url: "/api/external-squads/tags";
    const endpointDetails: import("../../../constants").EndpointDetails;
    const ResponseSchema: z.ZodObject<{
        response: z.ZodObject<{
            tags: z.ZodArray<z.ZodString>;
        }, z.core.$strip>;
    }, z.core.$strip>;
    type Response = z.infer<typeof ResponseSchema>;
}
//# sourceMappingURL=get-external-squads-tags.command.d.ts.map