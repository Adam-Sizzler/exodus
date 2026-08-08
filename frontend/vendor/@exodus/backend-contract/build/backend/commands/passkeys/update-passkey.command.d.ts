import { z } from 'zod';
export declare namespace UpdatePasskeyCommand {
    const url: "/api/passkeys/";
    const TSQ_url: "/api/passkeys/";
    const endpointDetails: import("../../constants").EndpointDetails;
    const RequestBodySchema: z.ZodObject<{
        id: z.ZodString;
        name: z.ZodString;
    }, z.core.$strip>;
    const ResponseSchema: z.ZodObject<{
        response: z.ZodObject<{
            passkeys: z.ZodArray<z.ZodObject<{
                id: z.ZodString;
                name: z.ZodString;
                createdAt: z.ZodPipe<z.ZodISODateTime, z.ZodTransform<Date, string>>;
                lastUsedAt: z.ZodPipe<z.ZodISODateTime, z.ZodTransform<Date, string>>;
            }, z.core.$strip>>;
        }, z.core.$strip>;
    }, z.core.$strip>;
    type RequestBody = z.infer<typeof RequestBodySchema>;
    type Response = z.infer<typeof ResponseSchema>;
}
//# sourceMappingURL=update-passkey.command.d.ts.map