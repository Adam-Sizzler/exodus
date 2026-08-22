import { z } from 'zod';
export declare namespace SyncSharedListCommand {
    const url: "/api/node-plugins/shared-lists/actions/sync";
    const TSQ_url: "/api/node-plugins/shared-lists/actions/sync";
    const endpointDetails: import("../../../constants").EndpointDetails;
    const RequestBodySchema: z.ZodObject<{
        name: z.ZodString;
    }, z.core.$strip>;
    type RequestBody = z.infer<typeof RequestBodySchema>;
}
//# sourceMappingURL=sync-shared-list.command.d.ts.map