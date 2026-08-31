import { z } from 'zod';
export declare namespace CreateSubpageConfigCommand {
    const url: "/api/subscription-page-configs/";
    const TSQ_url: "/api/subscription-page-configs/";
    const endpointDetails: import("../../constants").EndpointDetails;
    const RequestBodySchema: z.ZodObject<{
        name: z.ZodString;
    }, z.core.$strip>;
    const ResponseSchema: z.ZodObject<{
        response: z.ZodObject<{
            uuid: z.ZodUUID;
            viewPosition: z.ZodNumber;
            name: z.ZodString;
            tags: z.ZodDefault<z.ZodArray<z.ZodString>>;
            config: z.ZodNullable<z.ZodUnknown>;
        }, z.core.$strip>;
    }, z.core.$strip>;
    type RequestBody = z.infer<typeof RequestBodySchema>;
    type Response = z.infer<typeof ResponseSchema>;
}
//# sourceMappingURL=create-subpage-config.command.d.ts.map