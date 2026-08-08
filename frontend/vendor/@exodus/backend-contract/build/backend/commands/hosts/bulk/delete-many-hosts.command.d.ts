import { z } from 'zod';
export declare namespace BulkDeleteHostsCommand {
    const url: "/api/hosts/bulk/delete";
    const TSQ_url: "/api/hosts/bulk/delete";
    const endpointDetails: import("../../../constants").EndpointDetails;
    const RequestBodySchema: z.ZodObject<{
        uuids: z.ZodArray<z.ZodUUID>;
    }, z.core.$strip>;
    type RequestBody = z.infer<typeof RequestBodySchema>;
}
//# sourceMappingURL=delete-many-hosts.command.d.ts.map