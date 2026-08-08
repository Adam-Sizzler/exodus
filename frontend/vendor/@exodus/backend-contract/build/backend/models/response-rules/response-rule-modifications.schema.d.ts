import z from 'zod';
export declare const ResponseRuleEncryptionSchema: z.ZodObject<{
    method: z.ZodEnum<{
        age1: "age1";
        age1pq1: "age1pq1";
    }>;
    key: z.ZodString;
}, z.core.$strip>;
export declare const ResponseRuleModificationsSchema: z.ZodOptional<z.ZodObject<{
    headers: z.ZodOptional<z.ZodArray<z.ZodObject<{
        key: z.ZodString;
        value: z.ZodString;
    }, z.core.$strip>>>;
    applyHeadersToEnd: z.ZodOptional<z.ZodBoolean>;
    subscriptionTemplate: z.ZodOptional<z.ZodString>;
    ignoreHostXrayJsonTemplate: z.ZodOptional<z.ZodBoolean>;
    ignoreServeJsonAtBaseSubscription: z.ZodOptional<z.ZodBoolean>;
    additionalExtendedClientsRegex: z.ZodOptional<z.ZodArray<z.ZodString>>;
    disableHwidCheck: z.ZodOptional<z.ZodBoolean>;
    encryption: z.ZodOptional<z.ZodObject<{
        method: z.ZodEnum<{
            age1: "age1";
            age1pq1: "age1pq1";
        }>;
        key: z.ZodString;
    }, z.core.$strip>>;
    excludeHostsByTags: z.ZodOptional<z.ZodArray<z.ZodString>>;
}, z.core.$strip>>;
//# sourceMappingURL=response-rule-modifications.schema.d.ts.map