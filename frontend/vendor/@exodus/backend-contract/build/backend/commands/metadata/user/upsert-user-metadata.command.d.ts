import { z } from 'zod';
export declare namespace UpsertUserMetadataCommand {
    const url: (userId: string) => string;
    const TSQ_url: string;
    const endpointDetails: import("../../../constants").EndpointDetails;
    const RequestParamsSchema: z.ZodObject<{
        userId: z.ZodCoercedNumber<unknown>;
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
//# sourceMappingURL=upsert-user-metadata.command.d.ts.map