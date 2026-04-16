import { z } from 'zod';
export declare namespace DeletePasskeyCommand {
    const url: "/api/passkeys/";
    const TSQ_url: "/api/passkeys/";
    const endpointDetails: import("../../constants").EndpointDetails;
    const RequestSchema: z.ZodObject<{
        id: z.ZodString;
    }, "strip", z.ZodTypeAny, {
        id: string;
    }, {
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
                createdAt: Date;
                id: string;
                name: string;
                lastUsedAt: Date;
            }, {
                createdAt: string;
                id: string;
                name: string;
                lastUsedAt: string;
            }>, "many">;
        }, "strip", z.ZodTypeAny, {
            passkeys: {
                createdAt: Date;
                id: string;
                name: string;
                lastUsedAt: Date;
            }[];
        }, {
            passkeys: {
                createdAt: string;
                id: string;
                name: string;
                lastUsedAt: string;
            }[];
        }>;
    }, "strip", z.ZodTypeAny, {
        response: {
            passkeys: {
                createdAt: Date;
                id: string;
                name: string;
                lastUsedAt: Date;
            }[];
        };
    }, {
        response: {
            passkeys: {
                createdAt: string;
                id: string;
                name: string;
                lastUsedAt: string;
            }[];
        };
    }>;
    type Response = z.infer<typeof ResponseSchema>;
}
//# sourceMappingURL=delete-passkeys.command.d.ts.map