import { z } from 'zod';
export declare namespace BulkDisableHostsCommand {
    const url: "/api/hosts/bulk/disable";
    const TSQ_url: "/api/hosts/bulk/disable";
    const endpointDetails: import("../../../constants").EndpointDetails;
    const RequestBodySchema: z.ZodObject<{
        uuids: z.ZodArray<z.ZodUUID>;
    }, z.core.$strip>;
    type RequestBody = z.infer<typeof RequestBodySchema>;
}
//# sourceMappingURL=disable-many-hosts.command.d.ts.map