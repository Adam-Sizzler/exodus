import { z } from 'zod';

import { RESPONSE_RULES_CONFIG_VERSION } from '../../constants';
import { ResponseRuleSettingsSchema } from './response-rule-settings.schema';
import { ResponseRuleSchema } from './response-rule.schema';
import { EXAMPLES_SRR_BLANK_RULE } from './response-rules-examples';

export const ResponseRulesConfigSchema = z.object({
    version: z.enum(RESPONSE_RULES_CONFIG_VERSION).meta({
        title: 'Response Rules Config Version',
        markdownDescription:
            'Version of the **response rules** config. Currently supported version is **1**.',
    }),
    settings: ResponseRuleSettingsSchema,
    rules: z.array(ResponseRuleSchema).meta({
        title: 'Response Rules',
        markdownDescription: `Array of **response rules**. Rules are evaluated in order and the first rule that matches is applied. If no rule matches, request will be blocked by default.\n\n**Example:**\n\`\`\`json\n${JSON.stringify([EXAMPLES_SRR_BLANK_RULE], null, 2)}\n\`\`\``,
        defaultSnippets: [],
    }),
});
