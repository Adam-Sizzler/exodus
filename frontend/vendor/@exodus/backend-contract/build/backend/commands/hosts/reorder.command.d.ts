import { z } from 'zod';
export declare namespace ReorderHostsCommand {
    const url: "/api/hosts/actions/reorder";
    const TSQ_url: "/api/hosts/actions/reorder";
    const endpointDetails: import("../../constants").EndpointDetails;
    const RequestBodySchema: z.ZodObject<{
        hosts: z.ZodArray<z.ZodObject<{
            uuid: z.ZodUUID;
            viewPosition: z.ZodInt;
        }, z.core.$strip>>;
    }, z.core.$strip>;
    const ResponseSchema: z.ZodObject<{
        response: z.ZodObject<{
            isUpdated: z.ZodBoolean;
        }, z.core.$strip>;
    }, z.core.$strip>;
    type RequestBody = z.infer<typeof RequestBodySchema>;
    type Response = z.infer<typeof ResponseSchema>;
}
//# sourceMappingURL=reorder.command.d.ts.map