import { z } from 'zod';
export declare namespace UpsertNodeMetadataCommand {
    const url: (uuid: string) => string;
    const TSQ_url: string;
    const endpointDetails: import("../../../constants").EndpointDetails;
    const RequestParamsSchema: z.ZodObject<{
        uuid: z.ZodUUID;
    }, z.core.$strip>;
    const RequestBodySchema: z.ZodObject<{
        metadata: z.ZodObject<{}, z.core.$loose>;
    }, z.core.$strip>;
    const ResponseSchema: z.ZodObject<{
        response: z.ZodObject<{
            metadata: z.ZodObject<{}, z.core.$loose>;
        }, z.core.$strip>;
    }, z.core.$strip>;
    type RequestParams = z.infer<typeof RequestParamsSchema>;
    type RequestBody = z.infer<typeof RequestBodySchema>;
    type Response = z.infer<typeof ResponseSchema>;
}
//# sourceMappingURL=upsert-node-metadata.command.d.ts.map