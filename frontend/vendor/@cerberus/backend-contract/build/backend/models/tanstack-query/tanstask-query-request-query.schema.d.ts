import { z } from 'zod';
export declare const TanstackQueryRequestQuerySchema: z.ZodObject<{
    start: z.ZodDefault<z.ZodNumber>;
    size: z.ZodDefault<z.ZodNumber>;
    filters: z.ZodOptional<z.ZodEffects<z.ZodArray<z.ZodObject<{
        id: z.ZodString;
        value: z.ZodUnknown;
    }, "strip", z.ZodTypeAny, {
        id: string;
        value?: unknown;
    }, {
        id: string;
        value?: unknown;
    }>, "many">, {
        id: string;
        value?: unknown;
    }[], unknown>>;
    filterModes: z.ZodOptional<z.ZodEffects<z.ZodRecord<z.ZodString, z.ZodString>, Record<string, string>, unknown>>;
    globalFilterMode: z.ZodOptional<z.ZodString>;
    sorting: z.ZodOptional<z.ZodEffects<z.ZodArray<z.ZodObject<{
        id: z.ZodString;
        desc: z.ZodBoolean;
    }, "strip", z.ZodTypeAny, {
        id: string;
        desc: boolean;
    }, {
        id: string;
        desc: boolean;
    }>, "many">, {
        id: string;
        desc: boolean;
    }[], unknown>>;
}, "strip", z.ZodTypeAny, {
    start: number;
    size: number;
    filters?: {
        id: string;
        value?: unknown;
    }[] | undefined;
    filterModes?: Record<string, string> | undefined;
    globalFilterMode?: string | undefined;
    sorting?: {
        id: string;
        desc: boolean;
    }[] | undefined;
}, {
    start?: number | undefined;
    size?: number | undefined;
    filters?: unknown;
    filterModes?: unknown;
    globalFilterMode?: string | undefined;
    sorting?: unknown;
}>;
//# sourceMappingURL=tanstask-query-request-query.schema.d.ts.map