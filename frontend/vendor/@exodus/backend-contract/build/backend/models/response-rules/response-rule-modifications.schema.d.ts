import z from 'zod';
export declare const ResponseRuleEncryptionSchema: z.ZodObject<{
    method: z.ZodEnum<["age1", "age1pq1"]>;
    key: z.ZodString;
}, "strip", z.ZodTypeAny, {
    key: string;
    method: "age1" | "age1pq1";
}, {
    key: string;
    method: "age1" | "age1pq1";
}>;
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
    additionalExtendedClientsRegex: z.ZodOptional<z.ZodArray<z.ZodString, "many">>;
    disableHwidCheck: z.ZodOptional<z.ZodBoolean>;
    encryption: z.ZodOptional<z.ZodObject<{
        method: z.ZodEnum<["age1", "age1pq1"]>;
        key: z.ZodString;
    }, "strip", z.ZodTypeAny, {
        key: string;
        method: "age1" | "age1pq1";
    }, {
        key: string;
        method: "age1" | "age1pq1";
    }>>;
    excludeHostsByTags: z.ZodOptional<z.ZodArray<z.ZodString, "many">>;
}, "strip", z.ZodTypeAny, {
    headers?: {
        value: string;
        key: string;
    }[] | undefined;
    applyHeadersToEnd?: boolean | undefined;
    subscriptionTemplate?: string | undefined;
    ignoreHostXrayJsonTemplate?: boolean | undefined;
    ignoreServeJsonAtBaseSubscription?: boolean | undefined;
    additionalExtendedClientsRegex?: string[] | undefined;
    disableHwidCheck?: boolean | undefined;
    encryption?: {
        key: string;
        method: "age1" | "age1pq1";
    } | undefined;
    excludeHostsByTags?: string[] | undefined;
}, {
    headers?: {
        value: string;
        key: string;
    }[] | undefined;
    applyHeadersToEnd?: boolean | undefined;
    subscriptionTemplate?: string | undefined;
    ignoreHostXrayJsonTemplate?: boolean | undefined;
    ignoreServeJsonAtBaseSubscription?: boolean | undefined;
    additionalExtendedClientsRegex?: string[] | undefined;
    disableHwidCheck?: boolean | undefined;
    encryption?: {
        key: string;
        method: "age1" | "age1pq1";
    } | undefined;
    excludeHostsByTags?: string[] | undefined;
}>>;
//# sourceMappingURL=response-rule-modifications.schema.d.ts.map