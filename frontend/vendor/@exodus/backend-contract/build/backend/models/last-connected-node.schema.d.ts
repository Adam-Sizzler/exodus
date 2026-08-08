import { z } from 'zod';
export declare const LastConnectedNodeSchema: z.ZodNullable<z.ZodObject<{
    connectedAt: z.ZodPipe<z.ZodISODateTime, z.ZodTransform<Date, string>>;
    nodeName: z.ZodString;
    countryCode: z.ZodString;
}, z.core.$strip>>;
//# sourceMappingURL=last-connected-node.schema.d.ts.map