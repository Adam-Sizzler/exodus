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
    applyHeadersToEnd: z.ZodOptional<z.ZodOptional<z.ZodBoolean>>;
    subscriptionTemplate: z.ZodOptional<z.ZodString>;
    ignoreHostXrayJsonTemplate: z.ZodOptional<z.ZodBoolean>;
    ignoreServeJsonAtBaseSubscription: z.ZodOptional<z.ZodBoolean>;
}, "strip", z.ZodTypeAny, {
    headers?: {
        value: string;
        key: string;
    }[] | undefined;
    applyHeadersToEnd?: boolean | undefined;
    subscriptionTemplate?: string | undefined;
    ignoreHostXrayJsonTemplate?: boolean | undefined;
    ignoreServeJsonAtBaseSubscription?: boolean | undefined;
}, {
    headers?: {
        value: string;
        key: string;
    }[] | undefined;
    applyHeadersToEnd?: boolean | undefined;
    subscriptionTemplate?: string | undefined;
    ignoreHostXrayJsonTemplate?: boolean | undefined;
    ignoreServeJsonAtBaseSubscription?: boolean | undefined;
}>>;
//# sourceMappingURL=response-rule-modifications.schema.d.ts.map