import { z } from 'zod';

import {
    RESPONSE_RULES_CONDITION_OPERATORS,
    RESPONSE_RULES_CONDITION_OPERATORS_DESCRIPTION,
} from '../../constants';

export const ResponseRuleConditionSchema = z
    .object({
        headerName: z
            .string()
            .regex(
                /^[!#$%&'*+\-.0-9A-Z^_`a-z|~]+$/,
                'Invalid header name. Only letters(a-z, A-Z), numbers(0-9), underscores(_) and hyphens(-) are allowed.',
            )
            .meta({
                title: 'Header Name',
                markdownDescription:
                    '**Name** of the HTTP header to check. Must comply with RFC 7230.',
            }),
        operator: z.enum(RESPONSE_RULES_CONDITION_OPERATORS).meta({
            errorMessage: 'Invalid operator. Please select a valid operator.',
            markdownDescription: 'Operator to use for comparing the `headerName` with `value`.',
            markdownEnumDescriptions: Object.entries(
                RESPONSE_RULES_CONDITION_OPERATORS_DESCRIPTION,
            ).map(([_key, description]) => description),
        }),
        value: z
            .string()
            .min(1, 'Value is required')
            .max(255, 'Value must be less than 255 characters')
            .meta({
                markdownDescription: '**Value** to check against the **headerName**.',
            }),
        caseSensitive: z.boolean().meta({
            markdownDescription:
                'Whether the value is **case sensitive**. \n\n - `true`: the value will be compared as is. \n\n - `false`: the value will be lowercased **before** comparison.',
        }),
    })
    .meta({
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
    });
