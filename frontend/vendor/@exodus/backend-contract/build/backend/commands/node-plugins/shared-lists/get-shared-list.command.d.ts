import { z } from 'zod';
export declare namespace GetSharedListCommand {
    const url: "/api/node-plugins/shared-lists/by-name";
    const TSQ_url: "/api/node-plugins/shared-lists/by-name";
    const endpointDetails: import("../../../constants").EndpointDetails;
    const RequestQuerySchema: z.ZodObject<{
        name: z.ZodString;
    }, z.core.$strip>;
    const ResponseSchema: z.ZodObject<{
        response: z.ZodObject<{
            name: z.ZodString;
            config: z.ZodRecord<z.ZodString, z.ZodUnknown>;
        }, z.core.$strip>;
    }, z.core.$strip>;
    type RequestQuery = z.infer<typeof RequestQuerySchema>;
    type Response = z.infer<typeof ResponseSchema>;
}
//# sourceMappingURL=get-shared-list.command.d.ts.map