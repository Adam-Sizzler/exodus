"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.ResponseRuleSchema = exports.ResponseRuleSchemaBase = void 0;
const zod_1 = require("zod");
const response_rules_examples_1 = require("./response-rules-examples");
const constants_1 = require("../../constants");
const response_rule_modifications_schema_1 = require("./response-rule-modifications.schema");
const response_rule_condition_schema_1 = require("./response-rule-condition.schema");
exports.ResponseRuleSchemaBase = zod_1.z.object({
    name: zod_1.z
        .string()
        .min(1, 'Name is required')
        .max(50, 'Name must be less than 50 characters')
        .describe(JSON.stringify({
        markdownDescription: 'Name of the response rule.',
    })),
    description: zod_1.z
        .string()
        .min(1, 'Description is required')
        .max(250, 'Description must be less than 250 characters')
        .optional()
        .describe(JSON.stringify({
        markdownDescription: 'Description of the response rule. Optional.',
    })),
    enabled: zod_1.z.boolean().describe(JSON.stringify({
        markdownDescription: 'Control whether the response rule is enabled or disabled. \n\n - `true` the rule will be applied. \n\n - `false` the rule will be always ignored.',
    })),
    operator: zod_1.z.nativeEnum(constants_1.RESPONSE_RULES_OPERATORS).describe(JSON.stringify({
        markdownDescription: 'Operator to use for combining conditions in the rule.',
    })),
    conditions: zod_1.z.array(response_rule_condition_schema_1.ResponseRuleConditionSchema).describe(JSON.stringify({
        markdownDescription: 'Array of conditions to check against the request headers. Conditions are applied with **operator**. If conditions are empty, the rule will be matched.',
    })),
    responseType: zod_1.z.nativeEnum(constants_1.RESPONSE_RULES_RESPONSE_TYPES).describe(JSON.stringify({
        errorMessage: 'Invalid response type. Please select a valid response type.',
        markdownDescription: `Type of the response. Determines the type of **response** to be returned when the rule is matched.`,
        markdownEnumDescriptions: Object.entries(constants_1.RESPONSE_RULES_RESPONSE_TYPES_DESCRIPTION).map(([_key, description]) => description),
    })),
    responseModifications: response_rule_modifications_schema_1.ResponseRuleModificationsSchema,
});
exports.ResponseRuleSchema = exports.ResponseRuleSchemaBase.describe(JSON.stringify({
    defaultSnippets: [
        {
            label: 'Examples: Blank rule',
            markdownDescription: `Simple blank rule with no conditions or modifications.\n\`\`\`json\n${JSON.stringify(response_rules_examples_1.EXAMPLES_SRR_BLANK_RULE, null, 2)}\n\`\`\``,
            body: {
                ...response_rules_examples_1.EXAMPLES_SRR_BLANK_RULE,
            },
        },
        {
            label: 'Examples: Block Legacy Clients',
            markdownDescription: `Block requests from legacy clients\n\`\`\`json\n${JSON.stringify(response_rules_examples_1.EXAMPLES_SRR_BLOCK_LEGACY_CLIENTS_RULE, null, 2)}\n\`\`\``,
            body: {
                ...response_rules_examples_1.EXAMPLES_SRR_BLOCK_LEGACY_CLIENTS_RULE,
            },
        },
    ],
    title: 'Response Rule',
    markdownDescription: (0, response_rules_examples_1.generateResponseRuleDescription)(exports.ResponseRuleSchemaBase),
}));
