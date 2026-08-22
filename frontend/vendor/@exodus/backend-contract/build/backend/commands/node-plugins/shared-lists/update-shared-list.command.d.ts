import { z } from 'zod';
export declare namespace UpdateSharedListCommand {
    const url: "/api/node-plugins/shared-lists";
    const TSQ_url: "/api/node-plugins/shared-lists";
    const endpointDetails: import("../../../constants").EndpointDetails;
    const RequestBodySchema: z.ZodObject<{
        name: z.ZodString;
        config: z.ZodRecord<z.ZodString, z.ZodUnknown>;
    }, z.core.$strip>;
    const ResponseSchema: z.ZodObject<{
        response: z.ZodObject<{
            name: z.ZodString;
            config: z.ZodRecord<z.ZodString, z.ZodUnknown>;
        }, z.core.$strip>;
    }, z.core.$strip>;
    type RequestBody = z.infer<typeof RequestBodySchema>;
    type Response = z.infer<typeof ResponseSchema>;
}
//# sourceMappingURL=update-shared-list.command.d.ts.map