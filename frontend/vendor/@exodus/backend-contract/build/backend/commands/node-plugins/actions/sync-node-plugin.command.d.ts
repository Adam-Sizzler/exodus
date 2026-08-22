import { z } from 'zod';
export declare namespace SyncNodePluginCommand {
    const url: "/api/node-plugins/actions/sync";
    const TSQ_url: "/api/node-plugins/actions/sync";
    const endpointDetails: import("../../../constants").EndpointDetails;
    const RequestBodySchema: z.ZodObject<{
        uuid: z.ZodUUID;
    }, z.core.$strip>;
    type RequestBody = z.infer<typeof RequestBodySchema>;
}
//# sourceMappingURL=sync-node-plugin.command.d.ts.map