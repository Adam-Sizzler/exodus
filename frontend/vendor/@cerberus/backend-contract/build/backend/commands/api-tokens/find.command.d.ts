import { z } from 'zod';
export declare namespace FindAllApiTokensCommand {
    const url: "/api/tokens/";
    const TSQ_url: "/api/tokens/";
    const endpointDetails: import("../../constants").EndpointDetails;
    const ResponseSchema: z.ZodObject<{
        response: z.ZodObject<{
            apiKeys: z.ZodArray<z.ZodObject<{
                uuid: z.ZodString;
                token: z.ZodString;
                tokenName: z.ZodString;
                createdAt: z.ZodEffects<z.ZodString, Date, string>;
                updatedAt: z.ZodEffects<z.ZodString, Date, string>;
            }, "strip", z.ZodTypeAny, {
                tokenName: string;
                token: string;
                uuid: string;
                createdAt: Date;
                updatedAt: Date;
            }, {
                tokenName: string;
                token: string;
                uuid: string;
                createdAt: string;
                updatedAt: string;
            }>, "many">;
            docs: z.ZodObject<{
                isDocsEnabled: z.ZodBoolean;
                scalarPath: z.ZodNullable<z.ZodString>;
                swaggerPath: z.ZodNullable<z.ZodString>;
            }, "strip", z.ZodTypeAny, {
                isDocsEnabled: boolean;
                scalarPath: string | null;
                swaggerPath: string | null;
            }, {
                isDocsEnabled: boolean;
                scalarPath: string | null;
                swaggerPath: string | null;
            }>;
        }, "strip", z.ZodTypeAny, {
            apiKeys: {
                tokenName: string;
                token: string;
                uuid: string;
                createdAt: Date;
                updatedAt: Date;
            }[];
            docs: {
                isDocsEnabled: boolean;
                scalarPath: string | null;
                swaggerPath: string | null;
            };
        }, {
            apiKeys: {
                tokenName: string;
                token: string;
                uuid: string;
                createdAt: string;
                updatedAt: string;
            }[];
            docs: {
                isDocsEnabled: boolean;
                scalarPath: string | null;
                swaggerPath: string | null;
            };
        }>;
    }, "strip", z.ZodTypeAny, {
        response: {
            apiKeys: {
                tokenName: string;
                token: string;
                uuid: string;
                createdAt: Date;
                updatedAt: Date;
            }[];
            docs: {
                isDocsEnabled: boolean;
                scalarPath: string | null;
                swaggerPath: string | null;
            };
        };
    }, {
        response: {
            apiKeys: {
                tokenName: string;
                token: string;
                uuid: string;
                createdAt: string;
                updatedAt: string;
            }[];
            docs: {
                isDocsEnabled: boolean;
                scalarPath: string | null;
                swaggerPath: string | null;
            };
        };
    }>;
    type Response = z.infer<typeof ResponseSchema>;
}
//# sourceMappingURL=find.command.d.ts.map