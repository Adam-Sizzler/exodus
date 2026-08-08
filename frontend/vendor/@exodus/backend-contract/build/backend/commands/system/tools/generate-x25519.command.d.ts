import { z } from 'zod';
export declare namespace GenerateX25519Command {
    const url: "/api/system/tools/x25519/generate";
    const TSQ_url: "/api/system/tools/x25519/generate";
    const endpointDetails: import("../../../constants").EndpointDetails;
    const ResponseSchema: z.ZodObject<{
        response: z.ZodObject<{
            keypairs: z.ZodArray<z.ZodObject<{
                publicKey: z.ZodString;
                privateKey: z.ZodString;
            }, z.core.$strip>>;
        }, z.core.$strip>;
    }, z.core.$strip>;
    type Response = z.infer<typeof ResponseSchema>;
}
//# sourceMappingURL=generate-x25519.command.d.ts.map