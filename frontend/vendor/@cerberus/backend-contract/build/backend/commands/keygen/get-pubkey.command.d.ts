import { z } from 'zod';
export declare namespace GetPubKeyCommand {
    const url: "/api/keygen/";
    const TSQ_url: "/api/keygen/";
    const endpointDetails: import("../../constants").EndpointDetails;
    const ResponseSchema: z.ZodObject<{
        response: z.ZodObject<{
            pubKey: z.ZodString;
        }, "strip", z.ZodTypeAny, {
            pubKey: string;
        }, {
            pubKey: string;
        }>;
    }, "strip", z.ZodTypeAny, {
        response: {
            pubKey: string;
        };
    }, {
        response: {
            pubKey: string;
        };
    }>;
    type Response = z.infer<typeof ResponseSchema>;
}
//# sourceMappingURL=get-pubkey.command.d.ts.map