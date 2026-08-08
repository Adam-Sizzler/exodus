"use strict";
var __importDefault = (this && this.__importDefault) || function (mod) {
    return (mod && mod.__esModule) ? mod : { "default": mod };
};
Object.defineProperty(exports, "__esModule", { value: true });
exports.ResponseRuleModificationsSchema = exports.ResponseRuleEncryptionSchema = void 0;
const zod_1 = require("zod");
exports.ResponseRuleEncryptionSchema = (zod_1.z || zod_1.default.z || zod_1.default).object({
    method: (zod_1.z || zod_1.default.z || zod_1.default).enum(['age1', 'age1pq1']),
    key: (zod_1.z || zod_1.default.z || zod_1.default).string(),
});
exports.ResponseRuleModificationsSchema = zod_1.default
    .object({
    headers: zod_1.default
        .array(zod_1.default
        .object({
        key: zod_1.default
            .string()
            .regex(/^[!#$%&'*+\-.0-9A-Z^_`a-z|~]+$/, 'Invalid header name. Only letters(a-z, A-Z), numbers(0-9), underscores(_) and hyphens(-) are allowed.')
            .meta({
            title: 'Header Key',
            markdownDescription: 'Key of the response header. Must comply with RFC 7230.',
        }),
        value: (zod_1.z || zod_1.default.z || zod_1.default).string().min(1).meta({
            title: 'Header Value',
            markdownDescription: 'Value of the response header.',
        }),
    })
        .meta({
        title: 'Headers',
        markdownDescription: '**Key** and **value** of the response header will be added to the response.',
    }))
        .meta({
        defaultSnippets: [
            {
                label: 'Examples: Add custom header',
                markdownDescription: 'Add a custom header to the response',
                body: [
                    {
                        key: 'X-Custom-Header',
                        value: 'CustomValue',
                    },
                ],
            },
        ],
        markdownDescription: 'Array of headers to be added when the rule is matched.',
    })
        .optional(),
    applyHeadersToEnd: (zod_1.z || zod_1.default.z || zod_1.default).boolean().optional().meta({
        title: 'Apply Headers to End',
        markdownDescription: 'By default, headers are added when forming the response. In some cases, headers set in SRR may be overridden by headers from other parts of the system. If you set this flag to **true**, headers from SRR will be added at the very end, just before the response is sent. In this case, SRR headers may override headers from other sections.',
    }),
    subscriptionTemplate: zod_1.default
        .string()
        .min(1, 'Subscription template name is required')
        .optional()
        .meta({
        title: 'Subscription Template',
        markdownDescription: 'Override the subscription template with the given name. If not provided, the default subscription template will be used. If the template name is not found, the default subscription template for this type will be used. **This modification have higher priority than settings from External Squads.**',
    }),
    ignoreHostXrayJsonTemplate: (zod_1.z || zod_1.default.z || zod_1.default).boolean().optional().meta({
        title: 'Ignore Host Xray Json Template',
        markdownDescription: "Each Host may have its own Xray Json Template. If you set this flag to **true**, the Xray Json Template defined by the SRR will be used. **The Host's Xray Json Template will be ignored.**",
    }),
    ignoreServeJsonAtBaseSubscription: (zod_1.z || zod_1.default.z || zod_1.default).boolean().optional().meta({
        title: 'Ignore Serve Json at Base Subscription',
        markdownDescription: 'If you set this flag to **true**, the **Serve JSON at Base Subscription** setting will be ignored (set to **false**).',
    }),
    additionalExtendedClientsRegex: zod_1.default
        .array((zod_1.z || zod_1.default.z || zod_1.default).string().min(1))
        .optional()
        .meta({
        markdownDescription: 'Additional regex patterns to match extended clients. Matched clients will receive `serverDescription` in the subscription response.\n\n' +
            '**Default Mihomo extended clients:**\n' +
            '- `^FlClash ?X/`\n' +
            '- `^Flowvy/`\n' +
            '- `^prizrak-box/`\n' +
            '- `^koala-clash/`\n\n' +
            '**Default Xray extended clients:**\n' +
            '- `^Happ/`\n' +
            '- `^INCY/`\n\n' +
            '**Example:** `["^MyClient/", "^CustomApp\\\\/v2"]`',
    }),
    disableHwidCheck: (zod_1.z || zod_1.default.z || zod_1.default).boolean().optional().meta({
        title: 'Disable HWID Check',
        markdownDescription: 'If you set this flag to **true**, the HWID check will be disabled. **This modification have higher priority than settings from Subscription Settings.**',
    }),
    encryption: exports.ResponseRuleEncryptionSchema.optional().meta({
        title: 'Encryption',
        markdownDescription: 'Encrypt response body with given parameters. Generate keypairs with Rescue CLI: `docker exec -it exodus cli`, select "Generate keypairs".',
    }),
    excludeHostsByTags: zod_1.default
        .array(zod_1.default
        .string()
        .regex(/^[A-Z0-9_:]+$/, 'Tag can only contain uppercase letters, numbers, underscores and colons')
        .max(36, 'Each tag must be less than 36 characters'))
        .min(1)
        .optional()
        .meta({
        title: 'Exclude Hosts by Tags',
        markdownDescription: 'Excludes hosts from the subscription output if at least one tag in the host matches the given tags.',
    }),
})
    .optional()
    .meta({
    title: 'Response Modifications',
    examples: [
        {
            headers: [
                {
                    key: 'X-Custom-Header',
                    value: 'CustomValue',
                },
            ],
        },
    ],
    markdownDescription: 'Response modifications to be applied when the rule is matched. Optional.',
});
