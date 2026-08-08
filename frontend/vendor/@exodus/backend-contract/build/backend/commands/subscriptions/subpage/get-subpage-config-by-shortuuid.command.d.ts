import { z } from 'zod';
export declare namespace GetSubpageConfigByShortUuidCommand {
    const url: (shortUuid: string) => string;
    const TSQ_url: string;
    const endpointDetails: import("../../../constants").EndpointDetails;
    const RequestParamSchema: z.ZodObject<{
        shortUuid: z.ZodString;
    }, z.core.$strip>;
    const RequestBodySchema: z.ZodObject<{
        requestHeaders: z.ZodRecord<z.ZodString, z.ZodString>;
    }, z.core.$strip>;
    type RequestBody = z.infer<typeof RequestBodySchema>;
    const ResponseSchema: z.ZodObject<{
        response: z.ZodObject<{
            subpageConfigUuid: z.ZodNullable<z.ZodUUID>;
            webpageAllowed: z.ZodBoolean;
        }, z.core.$strip>;
    }, z.core.$strip>;
    type RequestParam = z.infer<typeof RequestParamSchema>;
    type Response = z.infer<typeof ResponseSchema>;
}
//# sourceMappingURL=get-subpage-config-by-shortuuid.command.d.ts.map