import { z } from 'zod';
export declare namespace CreateApiTokenCommand {
    const url: "/api/tokens/";
    const TSQ_url: "/api/tokens/";
    const endpointDetails: import("../../constants").EndpointDetails;
    const RequestBodySchema: z.ZodObject<{
        name: z.ZodString;
        expiresInDays: z.ZodNumber;
        scopes: z.ZodDefault<z.ZodOptional<z.ZodArray<z.ZodString>>>;
    }, z.core.$strip>;
    const ResponseSchema: z.ZodObject<{
        response: z.ZodObject<{
            uuid: z.ZodUUID;
            name: z.ZodString;
            expireAt: z.ZodPipe<z.ZodISODateTime, z.ZodTransform<Date, string>>;
            scopes: z.ZodArray<z.ZodString>;
            createdAt: z.ZodPipe<z.ZodISODateTime, z.ZodTransform<Date, string>>;
            updatedAt: z.ZodPipe<z.ZodISODateTime, z.ZodTransform<Date, string>>;
            token: z.ZodString;
        }, z.core.$strip>;
    }, z.core.$strip>;
    type RequestBody = z.infer<typeof RequestBodySchema>;
    type Response = z.infer<typeof ResponseSchema>;
}
//# sourceMappingURL=create.command.d.ts.map