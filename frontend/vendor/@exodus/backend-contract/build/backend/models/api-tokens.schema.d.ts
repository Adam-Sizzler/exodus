import { z } from 'zod';
export declare const ApiTokensSchema: z.ZodObject<{
    uuid: z.ZodUUID;
    name: z.ZodString;
    expireAt: z.ZodPipe<z.ZodISODateTime, z.ZodTransform<Date, string>>;
    scopes: z.ZodArray<z.ZodString>;
    createdAt: z.ZodPipe<z.ZodISODateTime, z.ZodTransform<Date, string>>;
    updatedAt: z.ZodPipe<z.ZodISODateTime, z.ZodTransform<Date, string>>;
}, z.core.$strip>;
//# sourceMappingURL=api-tokens.schema.d.ts.map