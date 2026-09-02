import { z } from 'zod';
export declare const TanstackQueryRequestQuerySchema: z.ZodObject<{
    start: z.ZodDefault<z.ZodCoercedNumber<unknown>>;
    size: z.ZodDefault<z.ZodCoercedNumber<unknown>>;
    filters: z.ZodOptional<z.ZodPreprocess<z.ZodArray<z.ZodObject<{
        id: z.ZodString;
        value: z.ZodUnknown;
    }, z.core.$strip>>, unknown>>;
    filterModes: z.ZodOptional<z.ZodPreprocess<z.ZodRecord<z.ZodString, z.ZodString>, unknown>>;
    globalFilterMode: z.ZodOptional<z.ZodString>;
    sorting: z.ZodOptional<z.ZodPreprocess<z.ZodArray<z.ZodObject<{
        id: z.ZodString;
        desc: z.ZodBoolean;
    }, z.core.$strip>>, unknown>>;
}, z.core.$strip>;
//# sourceMappingURL=tanstask-query-request-query.schema.d.ts.map