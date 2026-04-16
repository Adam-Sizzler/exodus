"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.ResponseRuleConditionSchema = void 0;
const zod_1 = require("zod");
const constants_1 = require("../../constants");
exports.ResponseRuleConditionSchema = zod_1.z
    .object({
    headerName: zod_1.z
        .string()
        .regex(/^[!#$%&'*+\-.0-9A-Z^_`a-z|~]+$/, 'Invalid header name. Only letters(a-z, A-Z), numbers(0-9), underscores(_) and hyphens(-) are allowed.')
        .describe(JSON.stringify({
        markdownDescription: '**Name** of the HTTP header to check. Must comply with RFC 7230.',
    })),
    operator: zod_1.z.nativeEnum(constants_1.RESPONSE_RULES_CONDITION_OPERATORS).describe(JSON.stringify({
        errorMessage: 'Invalid operator. Please select a valid operator.',
        markdownDescription: 'Operator to use for comparing the `headerName` with `value`.',
        markdownEnumDescriptions: Object.entries(constants_1.RESPONSE_RULES_CONDITION_OPERATORS_DESCRIPTION).map(([_key, description]) => description),
    })),
    value: zod_1.z
        .string()
        .min(1, 'Value is required')
        .max(255, 'Value must be less than 255 characters')
        .describe(JSON.stringify({
        markdownDescription: `**Value** to check against the **headerName**.`,
    })),
    caseSensitive: zod_1.z.boolean().describe(JSON.stringify({
        markdownDescription: 'Whether the value is **case sensitive**. \n\n - `true`: the value will be compared as is. \n\n - `false`: the value will be lowercased **before** comparison.',
    })),
})
    .describe(JSON.stringify({
    markdownDescription: 'Condition to check against the **headerName**.',
    defaultSnippets: [
        {
            label: 'Examples: Check if header contains "text/html"',
            markdownDescription: 'Condition to check if **headerName** contains "text/html"',
            body: {
                headerName: 'accept',
                operator: 'CONTAINS',
                value: 'text/html',
                caseSensitive: true,
            },
        },
    ],
}));
