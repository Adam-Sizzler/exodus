import { z } from 'zod';
export declare namespace UpdateSnippetCommand {
    const url: "/api/snippets/";
    const TSQ_url: "/api/snippets/";
    const endpointDetails: import("../../constants").EndpointDetails;
    const RequestBodySchema: z.ZodObject<{
        name: z.ZodString;
        snippet: z.ZodArray<z.ZodObject<{}, z.core.$loose>>;
    }, z.core.$strip>;
    const ResponseSchema: z.ZodObject<{
        response: z.ZodObject<{
            total: z.ZodNumber;
            snippets: z.ZodArray<z.ZodObject<{
                name: z.ZodString;
                snippet: z.ZodUnknown;
            }, z.core.$strip>>;
        }, z.core.$strip>;
    }, z.core.$strip>;
    type RequestBody = z.infer<typeof RequestBodySchema>;
    type Response = z.infer<typeof ResponseSchema>;
}
//# sourceMappingURL=update-snippet.command.d.ts.map