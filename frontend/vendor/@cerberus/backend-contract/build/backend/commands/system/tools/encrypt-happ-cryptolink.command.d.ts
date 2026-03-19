import { z } from 'zod';
export declare namespace EncryptHappCryptoLinkCommand {
    const url: "/api/system/tools/happ/encrypt";
    const TSQ_url: "/api/system/tools/happ/encrypt";
    const endpointDetails: import("../../../constants").EndpointDetails;
    const RequestSchema: z.ZodObject<{
        linkToEncrypt: z.ZodString;
    }, "strip", z.ZodTypeAny, {
        linkToEncrypt: string;
    }, {
        linkToEncrypt: string;
    }>;
    type Request = z.infer<typeof RequestSchema>;
    const ResponseSchema: z.ZodObject<{
        response: z.ZodObject<{
            encryptedLink: z.ZodString;
        }, "strip", z.ZodTypeAny, {
            encryptedLink: string;
        }, {
            encryptedLink: string;
        }>;
    }, "strip", z.ZodTypeAny, {
        response: {
            encryptedLink: string;
        };
    }, {
        response: {
            encryptedLink: string;
        };
    }>;
    type Response = z.infer<typeof ResponseSchema>;
}
//# sourceMappingURL=encrypt-happ-cryptolink.command.d.ts.map