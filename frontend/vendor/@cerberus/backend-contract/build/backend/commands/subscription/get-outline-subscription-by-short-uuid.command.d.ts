import { z } from 'zod';
export declare namespace GetOutlineSubscriptionByShortUuidCommand {
    const url: (shortUuid: string) => string;
    const TSQ_url: string;
    const RequestSchema: z.ZodObject<{
        shortUuid: z.ZodString;
        type: z.ZodOptional<z.ZodString>;
        encodedTag: z.ZodOptional<z.ZodString>;
    }, "strip", z.ZodTypeAny, {
        shortUuid: string;
        type?: string | undefined;
        encodedTag?: string | undefined;
    }, {
        shortUuid: string;
        type?: string | undefined;
        encodedTag?: string | undefined;
    }>;
    type Request = z.infer<typeof RequestSchema>;
}
//# sourceMappingURL=get-outline-subscription-by-short-uuid.command.d.ts.map