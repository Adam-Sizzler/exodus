import z from 'zod';
import { ResponseRuleSchema, ResponseRuleSchemaBase } from './response-rule.schema';
export declare const EXAMPLES_SRR_BLANK_RULE: z.infer<typeof ResponseRuleSchema>;
export declare const EXAMPLES_SRR_BLOCK_LEGACY_CLIENTS_RULE: z.infer<typeof ResponseRuleSchema>;
export declare function generateResponseRuleDescription(schema: typeof ResponseRuleSchemaBase): string;
//# sourceMappingURL=response-rules-examples.d.ts.map