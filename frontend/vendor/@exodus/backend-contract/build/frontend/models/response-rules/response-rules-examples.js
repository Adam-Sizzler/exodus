"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.EXAMPLES_SRR_BLOCK_LEGACY_CLIENTS_RULE = exports.EXAMPLES_SRR_BLANK_RULE = void 0;
exports.generateResponseRuleDescription = generateResponseRuleDescription;
const constants_1 = require("../../constants");
exports.EXAMPLES_SRR_BLANK_RULE = {
    name: 'Blank rule',
    description: 'Blank rule',
    operator: 'AND',
    enabled: true,
    conditions: [],
    responseType: 'BLOCK',
    responseModifications: {
        headers: [],
    },
};
exports.EXAMPLES_SRR_BLOCK_LEGACY_CLIENTS_RULE = {
    name: 'Block Legacy Clients',
    description: 'Block requests from legacy clients',
    enabled: true,
    operator: constants_1.RESPONSE_RULES_OPERATORS.OR,
    conditions: [
        {
            headerName: 'user-agent',
            operator: constants_1.RESPONSE_RULES_CONDITION_OPERATORS.CONTAINS,
            value: 'Hiddify',
            caseSensitive: true,
        },
        {
            headerName: 'user-agent',
            operator: constants_1.RESPONSE_RULES_CONDITION_OPERATORS.CONTAINS,
            value: 'FoxRay',
            caseSensitive: true,
        },
    ],
    responseType: constants_1.RESPONSE_RULES_RESPONSE_TYPES.BLOCK,
};
function generateResponseRuleDescription(schema) {
    const fields = Object.entries(schema.shape).map(([key, value]) => {
        const desc = value.description ? JSON.parse(value.description) : {};
        return {
            name: key,
            description: desc.markdownDescription || desc.title || 'No description',
        };
    });
    const fieldsText = fields
        .map((field) => `- **${field.name}**: ${field.description}`)
        .join('\n');
    return `Response rule configuration.\n\n**Fields:**\n${fieldsText}\n\n**Example:**\n\`\`\`json\n${JSON.stringify(exports.EXAMPLES_SRR_BLOCK_LEGACY_CLIENTS_RULE, null, 2)}\n\`\`\``;
}
