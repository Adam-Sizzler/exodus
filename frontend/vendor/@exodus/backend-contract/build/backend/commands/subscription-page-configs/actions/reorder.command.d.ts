import { z } from 'zod';
export declare namespace ReorderSubpageConfigsCommand {
    const url: "/api/subscription-page-configs/actions/reorder";
    const TSQ_url: "/api/subscription-page-configs/actions/reorder";
    const endpointDetails: import("../../../constants").EndpointDetails;
    const RequestBodySchema: z.ZodObject<{
        items: z.ZodArray<z.ZodObject<{
            uuid: z.ZodUUID;
            viewPosition: z.ZodNumber;
        }, z.core.$strip>>;
    }, z.core.$strip>;
    const ResponseSchema: z.ZodObject<{
        response: z.ZodObject<{
            total: z.ZodNumber;
            configs: z.ZodArray<z.ZodObject<{
                uuid: z.ZodUUID;
                viewPosition: z.ZodNumber;
                name: z.ZodString;
                config: z.ZodNullable<z.ZodUnknown>;
            }, z.core.$strip>>;
        }, z.core.$strip>;
    }, z.core.$strip>;
    type RequestBody = z.infer<typeof RequestBodySchema>;
    type Response = z.infer<typeof ResponseSchema>;
}
//# sourceMappingURL=reorder.command.d.ts.map