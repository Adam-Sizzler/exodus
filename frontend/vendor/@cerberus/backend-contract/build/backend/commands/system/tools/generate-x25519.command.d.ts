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
            }, "strip", z.ZodTypeAny, {
                publicKey: string;
                privateKey: string;
            }, {
                publicKey: string;
                privateKey: string;
            }>, "many">;
        }, "strip", z.ZodTypeAny, {
            keypairs: {
                publicKey: string;
                privateKey: string;
            }[];
        }, {
            keypairs: {
                publicKey: string;
                privateKey: string;
            }[];
        }>;
    }, "strip", z.ZodTypeAny, {
        response: {
            keypairs: {
                publicKey: string;
                privateKey: string;
            }[];
        };
    }, {
        response: {
            keypairs: {
                publicKey: string;
                privateKey: string;
            }[];
        };
    }>;
    type Response = z.infer<typeof ResponseSchema>;
}
//# sourceMappingURL=generate-x25519.command.d.ts.map