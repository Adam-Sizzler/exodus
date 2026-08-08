import { z } from 'zod';
export declare namespace GetUsersTagsCommand {
    const url: "/api/users/tags";
    const TSQ_url: "/api/users/tags";
    const endpointDetails: import("../../../constants").EndpointDetails;
    const ResponseSchema: z.ZodObject<{
        response: z.ZodObject<{
            tags: z.ZodArray<z.ZodString>;
        }, z.core.$strip>;
    }, z.core.$strip>;
    type Response = z.infer<typeof ResponseSchema>;
}
//# sourceMappingURL=get-users-tags.command.d.ts.map