import { z } from 'zod';
export declare namespace UpdatePasskeyCommand {
    const url: "/api/passkeys/";
    const TSQ_url: "/api/passkeys/";
    const endpointDetails: import("../../constants").EndpointDetails;
    const RequestSchema: z.ZodObject<{
        id: z.ZodString;
        name: z.ZodString;
    }, "strip", z.ZodTypeAny, {
        name: string;
        id: string;
    }, {
        name: string;
        id: string;
    }>;
    type Request = z.infer<typeof RequestSchema>;
    const ResponseSchema: z.ZodObject<{
        response: z.ZodObject<{
            passkeys: z.ZodArray<z.ZodObject<{
                id: z.ZodString;
                name: z.ZodString;
                createdAt: z.ZodEffects<z.ZodString, Date, string>;
                lastUsedAt: z.ZodEffects<z.ZodString, Date, string>;
            }, "strip", z.ZodTypeAny, {
                name: string;
                createdAt: Date;
                id: string;
                lastUsedAt: Date;
            }, {
                name: string;
                createdAt: string;
                id: string;
                lastUsedAt: string;
            }>, "many">;
        }, "strip", z.ZodTypeAny, {
            passkeys: {
                name: string;
                createdAt: Date;
                id: string;
                lastUsedAt: Date;
            }[];
        }, {
            passkeys: {
                name: string;
                createdAt: string;
                id: string;
                lastUsedAt: string;
            }[];
        }>;
    }, "strip", z.ZodTypeAny, {
        response: {
            passkeys: {
                name: string;
                createdAt: Date;
                id: string;
                lastUsedAt: Date;
            }[];
        };
    }, {
        response: {
            passkeys: {
                name: string;
                createdAt: string;
                id: string;
                lastUsedAt: string;
            }[];
        };
    }>;
    type Response = z.infer<typeof ResponseSchema>;
}
//# sourceMappingURL=update-passkey.command.d.ts.map