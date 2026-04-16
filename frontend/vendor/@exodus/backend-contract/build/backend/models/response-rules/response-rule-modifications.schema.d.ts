import z from 'zod';
export declare const ResponseRuleModificationsSchema: z.ZodOptional<z.ZodObject<{
    headers: z.ZodOptional<z.ZodArray<z.ZodObject<{
        key: z.ZodString;
        value: z.ZodString;
    }, "strip", z.ZodTypeAny, {
        value: string;
        key: string;
    }, {
        value: string;
        key: string;
    }>, "many">>;
    subscriptionTemplate: z.ZodOptional<z.ZodString>;
}, "strip", z.ZodTypeAny, {
    headers?: {
        value: string;
        key: string;
    }[] | undefined;
    subscriptionTemplate?: string | undefined;
}, {
    headers?: {
        value: string;
        key: string;
    }[] | undefined;
    subscriptionTemplate?: string | undefined;
}>>;
//# sourceMappingURL=response-rule-modifications.schema.d.ts.map