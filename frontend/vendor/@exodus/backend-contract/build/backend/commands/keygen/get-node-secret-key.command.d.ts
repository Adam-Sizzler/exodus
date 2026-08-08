import { z } from 'zod';
export declare namespace GetNodeSecretKeyCommand {
    const url: "/api/keygen/";
    const TSQ_url: "/api/keygen/";
    const endpointDetails: import("../../constants").EndpointDetails;
    const ResponseSchema: z.ZodObject<{
        response: z.ZodObject<{
            secretKey: z.ZodString;
        }, z.core.$strip>;
    }, z.core.$strip>;
    type Response = z.infer<typeof ResponseSchema>;
}
//# sourceMappingURL=get-node-secret-key.command.d.ts.map