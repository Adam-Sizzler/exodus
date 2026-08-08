import { z } from 'zod';
export declare namespace BulkEnableHostsCommand {
    const url: "/api/hosts/bulk/enable";
    const TSQ_url: "/api/hosts/bulk/enable";
    const endpointDetails: import("../../../constants").EndpointDetails;
    const RequestBodySchema: z.ZodObject<{
        uuids: z.ZodArray<z.ZodUUID>;
    }, z.core.$strip>;
    type RequestBody = z.infer<typeof RequestBodySchema>;
}
//# sourceMappingURL=enable-many-hosts.command.d.ts.map