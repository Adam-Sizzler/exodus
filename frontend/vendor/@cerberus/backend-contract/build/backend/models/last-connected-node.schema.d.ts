import { z } from 'zod';
export declare const LastConnectedNodeSchema: z.ZodNullable<z.ZodObject<{
    connectedAt: z.ZodEffects<z.ZodString, Date, string>;
    nodeName: z.ZodString;
    countryCode: z.ZodString;
}, "strip", z.ZodTypeAny, {
    nodeName: string;
    countryCode: string;
    connectedAt: Date;
}, {
    nodeName: string;
    countryCode: string;
    connectedAt: string;
}>>;
//# sourceMappingURL=last-connected-node.schema.d.ts.map