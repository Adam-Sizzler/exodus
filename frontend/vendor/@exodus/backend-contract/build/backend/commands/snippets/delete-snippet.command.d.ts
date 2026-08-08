import { z } from 'zod';
export declare namespace DeleteSnippetCommand {
    const url: "/api/snippets/";
    const TSQ_url: "/api/snippets/";
    const endpointDetails: import("../../constants").EndpointDetails;
    const RequestBodySchema: z.ZodObject<{
        name: z.ZodString;
    }, z.core.$strip>;
    type RequestBody = z.infer<typeof RequestBodySchema>;
}
//# sourceMappingURL=delete-snippet.command.d.ts.map