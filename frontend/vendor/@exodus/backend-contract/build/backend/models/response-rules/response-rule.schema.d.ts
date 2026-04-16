import { z } from 'zod';
export declare const ResponseRuleSchemaBase: z.ZodObject<{
    name: z.ZodString;
    description: z.ZodOptional<z.ZodString>;
    enabled: z.ZodBoolean;
    operator: z.ZodNativeEnum<{
        readonly AND: "AND";
        readonly OR: "OR";
    }>;
    conditions: z.ZodArray<z.ZodObject<{
        headerName: z.ZodString;
        operator: z.ZodNativeEnum<{
            readonly EQUALS: "EQUALS";
            readonly NOT_EQUALS: "NOT_EQUALS";
            readonly CONTAINS: "CONTAINS";
            readonly NOT_CONTAINS: "NOT_CONTAINS";
            readonly STARTS_WITH: "STARTS_WITH";
            readonly NOT_STARTS_WITH: "NOT_STARTS_WITH";
            readonly ENDS_WITH: "ENDS_WITH";
            readonly NOT_ENDS_WITH: "NOT_ENDS_WITH";
            readonly REGEX: "REGEX";
            readonly NOT_REGEX: "NOT_REGEX";
        }>;
        value: z.ZodString;
        caseSensitive: z.ZodBoolean;
    }, "strip", z.ZodTypeAny, {
        value: string;
        headerName: string;
        operator: "EQUALS" | "NOT_EQUALS" | "CONTAINS" | "NOT_CONTAINS" | "STARTS_WITH" | "NOT_STARTS_WITH" | "ENDS_WITH" | "NOT_ENDS_WITH" | "REGEX" | "NOT_REGEX";
        caseSensitive: boolean;
    }, {
        value: string;
        headerName: string;
        operator: "EQUALS" | "NOT_EQUALS" | "CONTAINS" | "NOT_CONTAINS" | "STARTS_WITH" | "NOT_STARTS_WITH" | "ENDS_WITH" | "NOT_ENDS_WITH" | "REGEX" | "NOT_REGEX";
        caseSensitive: boolean;
    }>, "many">;
    responseType: z.ZodNativeEnum<{
        readonly BROWSER: "BROWSER";
        readonly BLOCK: "BLOCK";
        readonly STATUS_CODE_404: "STATUS_CODE_404";
        readonly STATUS_CODE_451: "STATUS_CODE_451";
        readonly SOCKET_DROP: "SOCKET_DROP";
        readonly XRAY_JSON: "XRAY_JSON";
        readonly XRAY_BASE64: "XRAY_BASE64";
        readonly MIHOMO: "MIHOMO";
        readonly STASH: "STASH";
        readonly CLASH: "CLASH";
        readonly SINGBOX: "SINGBOX";
    }>;
    responseModifications: z.ZodOptional<z.ZodObject<{
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
}, "strip", z.ZodTypeAny, {
    enabled: boolean;
    name: string;
    operator: "AND" | "OR";
    conditions: {
        value: string;
        headerName: string;
        operator: "EQUALS" | "NOT_EQUALS" | "CONTAINS" | "NOT_CONTAINS" | "STARTS_WITH" | "NOT_STARTS_WITH" | "ENDS_WITH" | "NOT_ENDS_WITH" | "REGEX" | "NOT_REGEX";
        caseSensitive: boolean;
    }[];
    responseType: "STASH" | "SINGBOX" | "MIHOMO" | "XRAY_JSON" | "CLASH" | "XRAY_BASE64" | "BROWSER" | "BLOCK" | "STATUS_CODE_404" | "STATUS_CODE_451" | "SOCKET_DROP";
    description?: string | undefined;
    responseModifications?: {
        headers?: {
            value: string;
            key: string;
        }[] | undefined;
        subscriptionTemplate?: string | undefined;
    } | undefined;
}, {
    enabled: boolean;
    name: string;
    operator: "AND" | "OR";
    conditions: {
        value: string;
        headerName: string;
        operator: "EQUALS" | "NOT_EQUALS" | "CONTAINS" | "NOT_CONTAINS" | "STARTS_WITH" | "NOT_STARTS_WITH" | "ENDS_WITH" | "NOT_ENDS_WITH" | "REGEX" | "NOT_REGEX";
        caseSensitive: boolean;
    }[];
    responseType: "STASH" | "SINGBOX" | "MIHOMO" | "XRAY_JSON" | "CLASH" | "XRAY_BASE64" | "BROWSER" | "BLOCK" | "STATUS_CODE_404" | "STATUS_CODE_451" | "SOCKET_DROP";
    description?: string | undefined;
    responseModifications?: {
        headers?: {
            value: string;
            key: string;
        }[] | undefined;
        subscriptionTemplate?: string | undefined;
    } | undefined;
}>;
export declare const ResponseRuleSchema: z.ZodObject<{
    name: z.ZodString;
    description: z.ZodOptional<z.ZodString>;
    enabled: z.ZodBoolean;
    operator: z.ZodNativeEnum<{
        readonly AND: "AND";
        readonly OR: "OR";
    }>;
    conditions: z.ZodArray<z.ZodObject<{
        headerName: z.ZodString;
        operator: z.ZodNativeEnum<{
            readonly EQUALS: "EQUALS";
            readonly NOT_EQUALS: "NOT_EQUALS";
            readonly CONTAINS: "CONTAINS";
            readonly NOT_CONTAINS: "NOT_CONTAINS";
            readonly STARTS_WITH: "STARTS_WITH";
            readonly NOT_STARTS_WITH: "NOT_STARTS_WITH";
            readonly ENDS_WITH: "ENDS_WITH";
            readonly NOT_ENDS_WITH: "NOT_ENDS_WITH";
            readonly REGEX: "REGEX";
            readonly NOT_REGEX: "NOT_REGEX";
        }>;
        value: z.ZodString;
        caseSensitive: z.ZodBoolean;
    }, "strip", z.ZodTypeAny, {
        value: string;
        headerName: string;
        operator: "EQUALS" | "NOT_EQUALS" | "CONTAINS" | "NOT_CONTAINS" | "STARTS_WITH" | "NOT_STARTS_WITH" | "ENDS_WITH" | "NOT_ENDS_WITH" | "REGEX" | "NOT_REGEX";
        caseSensitive: boolean;
    }, {
        value: string;
        headerName: string;
        operator: "EQUALS" | "NOT_EQUALS" | "CONTAINS" | "NOT_CONTAINS" | "STARTS_WITH" | "NOT_STARTS_WITH" | "ENDS_WITH" | "NOT_ENDS_WITH" | "REGEX" | "NOT_REGEX";
        caseSensitive: boolean;
    }>, "many">;
    responseType: z.ZodNativeEnum<{
        readonly BROWSER: "BROWSER";
        readonly BLOCK: "BLOCK";
        readonly STATUS_CODE_404: "STATUS_CODE_404";
        readonly STATUS_CODE_451: "STATUS_CODE_451";
        readonly SOCKET_DROP: "SOCKET_DROP";
        readonly XRAY_JSON: "XRAY_JSON";
        readonly XRAY_BASE64: "XRAY_BASE64";
        readonly MIHOMO: "MIHOMO";
        readonly STASH: "STASH";
        readonly CLASH: "CLASH";
        readonly SINGBOX: "SINGBOX";
    }>;
    responseModifications: z.ZodOptional<z.ZodObject<{
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
}, "strip", z.ZodTypeAny, {
    enabled: boolean;
    name: string;
    operator: "AND" | "OR";
    conditions: {
        value: string;
        headerName: string;
        operator: "EQUALS" | "NOT_EQUALS" | "CONTAINS" | "NOT_CONTAINS" | "STARTS_WITH" | "NOT_STARTS_WITH" | "ENDS_WITH" | "NOT_ENDS_WITH" | "REGEX" | "NOT_REGEX";
        caseSensitive: boolean;
    }[];
    responseType: "STASH" | "SINGBOX" | "MIHOMO" | "XRAY_JSON" | "CLASH" | "XRAY_BASE64" | "BROWSER" | "BLOCK" | "STATUS_CODE_404" | "STATUS_CODE_451" | "SOCKET_DROP";
    description?: string | undefined;
    responseModifications?: {
        headers?: {
            value: string;
            key: string;
        }[] | undefined;
        subscriptionTemplate?: string | undefined;
    } | undefined;
}, {
    enabled: boolean;
    name: string;
    operator: "AND" | "OR";
    conditions: {
        value: string;
        headerName: string;
        operator: "EQUALS" | "NOT_EQUALS" | "CONTAINS" | "NOT_CONTAINS" | "STARTS_WITH" | "NOT_STARTS_WITH" | "ENDS_WITH" | "NOT_ENDS_WITH" | "REGEX" | "NOT_REGEX";
        caseSensitive: boolean;
    }[];
    responseType: "STASH" | "SINGBOX" | "MIHOMO" | "XRAY_JSON" | "CLASH" | "XRAY_BASE64" | "BROWSER" | "BLOCK" | "STATUS_CODE_404" | "STATUS_CODE_451" | "SOCKET_DROP";
    description?: string | undefined;
    responseModifications?: {
        headers?: {
            value: string;
            key: string;
        }[] | undefined;
        subscriptionTemplate?: string | undefined;
    } | undefined;
}>;
//# sourceMappingURL=response-rule.schema.d.ts.map