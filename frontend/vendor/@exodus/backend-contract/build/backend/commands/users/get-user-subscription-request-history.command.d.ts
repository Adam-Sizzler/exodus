import { z } from 'zod';
export declare namespace GetUserSubscriptionRequestHistoryCommand {
    const url: (userId: string) => string;
    const TSQ_url: string;
    const endpointDetails: import("../../constants").EndpointDetails;
    const RequestParamSchema: z.ZodObject<{
        userId: z.ZodCoercedNumber<unknown>;
    }, z.core.$strip>;
    const ResponseSchema: z.ZodObject<{
        response: z.ZodObject<{
            total: z.ZodNumber;
            records: z.ZodArray<z.ZodObject<{
                id: z.ZodNumber;
                userId: z.ZodNumber;
                requestAt: z.ZodPipe<z.ZodISODateTime, z.ZodTransform<Date, string>>;
                srrResponseType: z.ZodString;
                requestIp: z.ZodNullable<z.ZodOptional<z.ZodString>>;
                userAgent: z.ZodNullable<z.ZodOptional<z.ZodString>>;
                srrRuleName: z.ZodNullable<z.ZodOptional<z.ZodString>>;
            }, z.core.$strip>>;
        }, z.core.$strip>;
    }, z.core.$strip>;
    type RequestParam = z.infer<typeof RequestParamSchema>;
    type Response = z.infer<typeof ResponseSchema>;
}
//# sourceMappingURL=get-user-subscription-request-history.command.d.ts.map