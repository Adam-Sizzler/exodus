import { z } from 'zod';
export declare namespace GetPubKeyCommand {
    const url: "/api/keygen/";
    const TSQ_url: "/api/keygen/";
    const endpointDetails: import("../../constants").EndpointDetails;
    const ResponseSchema: z.ZodObject<{
        response: z.ZodObject<{
            pubKey: z.ZodString;
            grpcToken: z.ZodOptional<z.ZodString>;
        }, "strip", z.ZodTypeAny, {
            pubKey: string;
            grpcToken?: string | undefined;
        }, {
            pubKey: string;
            grpcToken?: string | undefined;
        }>;
    }, "strip", z.ZodTypeAny, {
        response: {
            pubKey: string;
            grpcToken?: string | undefined;
        };
    }, {
        response: {
            pubKey: string;
            grpcToken?: string | undefined;
        };
    }>;
    type Response = z.infer<typeof ResponseSchema>;
}
//# sourceMappingURL=get-pubkey.command.d.ts.map