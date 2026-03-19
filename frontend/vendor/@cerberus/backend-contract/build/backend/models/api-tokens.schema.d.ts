import { z } from 'zod';
export declare const ApiTokensSchema: z.ZodObject<{
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
}>;
//# sourceMappingURL=api-tokens.schema.d.ts.map