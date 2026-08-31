import { z } from 'zod';
export declare namespace GetInternalSquadsTagsCommand {
    const url: "/api/internal-squads/tags";
    const TSQ_url: "/api/internal-squads/tags";
    const endpointDetails: import("../../../constants").EndpointDetails;
    const ResponseSchema: z.ZodObject<{
        response: z.ZodObject<{
            tags: z.ZodArray<z.ZodString>;
        }, z.core.$strip>;
    }, z.core.$strip>;
    type Response = z.infer<typeof ResponseSchema>;
}
//# sourceMappingURL=get-internal-squads-tags.command.d.ts.map