"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.ResponseRulesConfigSchema = void 0;
const zod_1 = require("zod");
const response_rules_examples_1 = require("./response-rules-examples");
const constants_1 = require("../../constants");
const response_rule_schema_1 = require("./response-rule.schema");
exports.ResponseRulesConfigSchema = zod_1.z.object({
    version: zod_1.z.nativeEnum(constants_1.RESPONSE_RULES_CONFIG_VERSION).describe(JSON.stringify({
        title: 'Response Rules Config Version',
        markdownDescription: 'Version of the **response rules** config. Currently supported version is **1**.',
    })),
    rules: zod_1.z.array(response_rule_schema_1.ResponseRuleSchema).describe(JSON.stringify({
        title: 'Response Rules',
        markdownDescription: `Array of **response rules**. Rules are evaluated in order and the first rule that matches is applied. If no rule matches, request will be blocked by default.\n\n**Example:**\n\`\`\`json\n${JSON.stringify([response_rules_examples_1.EXAMPLES_SRR_BLANK_RULE], null, 2)}\n\`\`\``,
        defaultSnippets: [],
    })),
});
