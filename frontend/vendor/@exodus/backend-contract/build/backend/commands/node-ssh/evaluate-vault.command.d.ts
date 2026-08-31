import { z } from 'zod';
export declare namespace EvaluateVaultCommand {
    const url: "/api/node-ssh/vault/evaluate";
    const TSQ_url: "/api/node-ssh/vault/evaluate";
    const endpointDetails: import("../../constants").EndpointDetails;
    const RequestBodySchema: z.ZodObject<{
        blinded: z.ZodBase64;
    }, z.core.$strip>;
    const ResponseSchema: z.ZodObject<{
        response: z.ZodObject<{
            evaluated: z.ZodBase64;
        }, z.core.$strip>;
    }, z.core.$strip>;
    type RequestBody = z.infer<typeof RequestBodySchema>;
    type Response = z.infer<typeof ResponseSchema>;
}
//# sourceMappingURL=evaluate-vault.command.d.ts.map