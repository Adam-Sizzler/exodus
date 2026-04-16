import { z } from 'zod';
export declare const ResponseRuleConditionSchema: z.ZodObject<{
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
}>;
//# sourceMappingURL=response-rule-condition.schema.d.ts.map