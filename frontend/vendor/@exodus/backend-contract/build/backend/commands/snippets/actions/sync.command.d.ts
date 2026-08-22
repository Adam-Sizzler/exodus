import { z } from 'zod';
export declare namespace SyncSnippetCommand {
    const url: "/api/snippets/actions/sync";
    const TSQ_url: "/api/snippets/actions/sync";
    const endpointDetails: import("../../../constants").EndpointDetails;
    const RequestBodySchema: z.ZodObject<{
        name: z.ZodString;
    }, z.core.$strip>;
    type RequestBody = z.infer<typeof RequestBodySchema>;
}
//# sourceMappingURL=sync.command.d.ts.map