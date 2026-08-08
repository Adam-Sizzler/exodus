import { z } from 'zod';
export declare namespace GetSnippetsCommand {
    const url: "/api/snippets/";
    const TSQ_url: "/api/snippets/";
    const endpointDetails: import("../../constants").EndpointDetails;
    const ResponseSchema: z.ZodObject<{
        response: z.ZodObject<{
            total: z.ZodNumber;
            snippets: z.ZodArray<z.ZodObject<{
                name: z.ZodString;
                snippet: z.ZodUnknown;
            }, z.core.$strip>>;
        }, z.core.$strip>;
    }, z.core.$strip>;
    type Response = z.infer<typeof ResponseSchema>;
}
//# sourceMappingURL=get-snippets.command.d.ts.map