import { z } from 'zod';
export declare const ApiTokensSchema: z.ZodObject<{
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
}>;
//# sourceMappingURL=api-tokens.schema.d.ts.map