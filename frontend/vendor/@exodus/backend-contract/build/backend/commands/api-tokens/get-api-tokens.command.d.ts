import { z } from 'zod';
export declare namespace GetApiTokensCommand {
    const url: "/api/tokens/";
    const TSQ_url: "/api/tokens/";
    const endpointDetails: import("../../constants").EndpointDetails;
    const ResponseSchema: z.ZodObject<{
        response: z.ZodObject<{
            tokens: z.ZodArray<z.ZodObject<{
                uuid: z.ZodUUID;
                name: z.ZodString;
                expireAt: z.ZodPipe<z.ZodISODateTime, z.ZodTransform<Date, string>>;
                scopes: z.ZodArray<z.ZodString>;
                createdAt: z.ZodPipe<z.ZodISODateTime, z.ZodTransform<Date, string>>;
                updatedAt: z.ZodPipe<z.ZodISODateTime, z.ZodTransform<Date, string>>;
            }, z.core.$strip>>;
        }, z.core.$strip>;
    }, z.core.$strip>;
    type Response = z.infer<typeof ResponseSchema>;
}
//# sourceMappingURL=get-api-tokens.command.d.ts.map