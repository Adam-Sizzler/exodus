import { z } from 'zod';
export declare namespace GetSubscriptionByShortUuidCommand {
    const url: (shortUuid: string) => string;
    const TSQ_url: string;
    const RequestSchema: z.ZodObject<{
        shortUuid: z.ZodString;
    }, "strip", z.ZodTypeAny, {
        shortUuid: string;
    }, {
        shortUuid: string;
    }>;
    type Request = z.infer<typeof RequestSchema>;
}
//# sourceMappingURL=get-subscription-by-short-uuid.command.d.ts.map