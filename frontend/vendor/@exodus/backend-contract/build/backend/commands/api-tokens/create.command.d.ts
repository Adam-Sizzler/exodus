import { z } from 'zod';
export declare namespace CreateApiTokenCommand {
    const url: "/api/tokens/";
    const TSQ_url: "/api/tokens/";
    const endpointDetails: import("../../constants").EndpointDetails;
    const RequestSchema: z.ZodObject<{
        name: z.ZodString;
        expiresInDays: z.ZodNumber;
        scopes: z.ZodDefault<z.ZodOptional<z.ZodArray<z.ZodString, "many">>>;
    }, "strip", z.ZodTypeAny, {
        scopes: string[];
        name: string;
        expiresInDays: number;
    }, {
        name: string;
        expiresInDays: number;
        scopes?: string[] | undefined;
    }>;
    type Request = z.infer<typeof RequestSchema>;
    const ResponseSchema: z.ZodObject<{
        response: z.ZodObject<{
            uuid: z.ZodString;
            name: z.ZodString;
            expireAt: z.ZodEffects<z.ZodString, Date, string>;
            scopes: z.ZodArray<z.ZodString, "many">;
            createdAt: z.ZodEffects<z.ZodString, Date, string>;
            updatedAt: z.ZodEffects<z.ZodString, Date, string>;
        } & {
            token: z.ZodString;
        }, "strip", z.ZodTypeAny, {
            scopes: string[];
            uuid: string;
            name: string;
            expireAt: Date;
            createdAt: Date;
            updatedAt: Date;
            token: string;
        }, {
            scopes: string[];
            uuid: string;
            name: string;
            expireAt: string;
            createdAt: string;
            updatedAt: string;
            token: string;
        }>;
    }, "strip", z.ZodTypeAny, {
        response: {
            scopes: string[];
            uuid: string;
            name: string;
            expireAt: Date;
            createdAt: Date;
            updatedAt: Date;
            token: string;
        };
    }, {
        response: {
            scopes: string[];
            uuid: string;
            name: string;
            expireAt: string;
            createdAt: string;
            updatedAt: string;
            token: string;
        };
    }>;
    type Response = z.infer<typeof ResponseSchema>;
}
//# sourceMappingURL=create.command.d.ts.map