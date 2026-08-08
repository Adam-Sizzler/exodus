import { z } from 'zod';
export declare namespace GetSubscriptionByShortUuidCommand {
    const url: (shortUuid: string) => string;
    const TSQ_url: string;
    const RequestParamSchema: z.ZodObject<{
        shortUuid: z.ZodString;
    }, z.core.$strip>;
    type RequestParam = z.infer<typeof RequestParamSchema>;
}
//# sourceMappingURL=get-subscription-by-short-uuid.command.d.ts.map