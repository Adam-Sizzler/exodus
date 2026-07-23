import { z } from 'zod';
export declare namespace GetAllPasskeysCommand {
    const url: "/api/passkeys/";
    const TSQ_url: "/api/passkeys/";
    const endpointDetails: import("../../constants").EndpointDetails;
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
//# sourceMappingURL=get-active-passkeys.command.d.ts.map