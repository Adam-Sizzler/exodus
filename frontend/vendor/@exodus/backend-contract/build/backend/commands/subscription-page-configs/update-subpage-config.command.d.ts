import { z } from 'zod';
export declare namespace UpdateSubpageConfigCommand {
    const url: "/api/subscription-page-configs/";
    const TSQ_url: "/api/subscription-page-configs/";
    const endpointDetails: import("../../constants").EndpointDetails;
    const RequestBodySchema: z.ZodObject<{
        uuid: z.ZodUUID;
        name: z.ZodOptional<z.ZodString>;
        config: z.ZodOptional<z.ZodUnknown>;
    }, z.core.$strip>;
    const ResponseSchema: z.ZodObject<{
        response: z.ZodObject<{
            uuid: z.ZodUUID;
            viewPosition: z.ZodNumber;
            name: z.ZodString;
            config: z.ZodUnknown;
        }, z.core.$strip>;
    }, z.core.$strip>;
    type RequestBody = z.infer<typeof RequestBodySchema>;
    type Response = z.infer<typeof ResponseSchema>;
}
//# sourceMappingURL=update-subpage-config.command.d.ts.map