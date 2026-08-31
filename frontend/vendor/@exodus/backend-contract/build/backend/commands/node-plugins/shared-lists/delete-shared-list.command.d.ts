import { z } from 'zod';
export declare namespace DeleteSharedListCommand {
    const url: "/api/node-plugins/shared-lists";
    const TSQ_url: "/api/node-plugins/shared-lists";
    const endpointDetails: import("../../../constants").EndpointDetails;
    const RequestBodySchema: z.ZodObject<{
        name: z.ZodString;
    }, z.core.$strip>;
    type RequestBody = z.infer<typeof RequestBodySchema>;
}
//# sourceMappingURL=delete-shared-list.command.d.ts.map