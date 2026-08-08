import { z } from 'zod';
export declare namespace GetPasskeysCommand {
    const url: "/api/passkeys/";
    const TSQ_url: "/api/passkeys/";
    const endpointDetails: import("../../constants").EndpointDetails;
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
    type Response = z.infer<typeof ResponseSchema>;
}
//# sourceMappingURL=get-passkeys.command.d.ts.map