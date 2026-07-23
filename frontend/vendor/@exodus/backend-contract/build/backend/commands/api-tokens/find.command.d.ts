import { z } from 'zod';
export declare namespace FindAllApiTokensCommand {
    const url: "/api/tokens/";
    const TSQ_url: "/api/tokens/";
    const endpointDetails: import("../../constants").EndpointDetails;
    const ResponseSchema: z.ZodObject<{
        response: z.ZodObject<{
            tokens: z.ZodArray<z.ZodObject<{
                uuid: z.ZodString;
                name: z.ZodString;
                expireAt: z.ZodEffects<z.ZodString, Date, string>;
                scopes: z.ZodArray<z.ZodString, "many">;
                createdAt: z.ZodEffects<z.ZodString, Date, string>;
                updatedAt: z.ZodEffects<z.ZodString, Date, string>;
            }, "strip", z.ZodTypeAny, {
                scopes: string[];
                uuid: string;
                name: string;
                expireAt: Date;
                createdAt: Date;
                updatedAt: Date;
            }, {
                scopes: string[];
                uuid: string;
                name: string;
                expireAt: string;
                createdAt: string;
                updatedAt: string;
            }>, "many">;
            docs: z.ZodObject<{
                enabled: z.ZodBoolean;
                scalarPath: z.ZodNullable<z.ZodString>;
                swaggerPath: z.ZodNullable<z.ZodString>;
            }, "strip", z.ZodTypeAny, {
                enabled: boolean;
                scalarPath: string | null;
                swaggerPath: string | null;
            }, {
                enabled: boolean;
                scalarPath: string | null;
                swaggerPath: string | null;
            }>;
        }, "strip", z.ZodTypeAny, {
            tokens: {
                scopes: string[];
                uuid: string;
                name: string;
                expireAt: Date;
                createdAt: Date;
                updatedAt: Date;
            }[];
            docs: {
                enabled: boolean;
                scalarPath: string | null;
                swaggerPath: string | null;
            };
        }, {
            tokens: {
                scopes: string[];
                uuid: string;
                name: string;
                expireAt: string;
                createdAt: string;
                updatedAt: string;
            }[];
            docs: {
                enabled: boolean;
                scalarPath: string | null;
                swaggerPath: string | null;
            };
        }>;
    }, "strip", z.ZodTypeAny, {
        response: {
            tokens: {
                scopes: string[];
                uuid: string;
                name: string;
                expireAt: Date;
                createdAt: Date;
                updatedAt: Date;
            }[];
            docs: {
                enabled: boolean;
                scalarPath: string | null;
                swaggerPath: string | null;
            };
        };
    }, {
        response: {
            tokens: {
                scopes: string[];
                uuid: string;
                name: string;
                expireAt: string;
                createdAt: string;
                updatedAt: string;
            }[];
            docs: {
                enabled: boolean;
                scalarPath: string | null;
                swaggerPath: string | null;
            };
        };
    }>;
    type Response = z.infer<typeof ResponseSchema>;
}
//# sourceMappingURL=find.command.d.ts.map