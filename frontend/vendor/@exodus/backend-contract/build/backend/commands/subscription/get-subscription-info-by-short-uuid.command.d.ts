import { z } from 'zod';
export declare namespace GetSubscriptionInfoByShortUuidCommand {
    const url: (shortUuid: string) => string;
    const TSQ_url: string;
    const endpointDetails: import("../../constants").EndpointDetails;
    const RequestParamSchema: z.ZodObject<{
        shortUuid: z.ZodString;
    }, z.core.$strip>;
    const ResponseSchema: z.ZodObject<{
        response: z.ZodObject<{
            isFound: z.ZodBoolean;
            user: z.ZodObject<{
                shortUuid: z.ZodString;
                daysLeft: z.ZodNumber;
                trafficUsed: z.ZodString;
                trafficLimit: z.ZodString;
                lifetimeTrafficUsed: z.ZodString;
                trafficUsedBytes: z.ZodString;
                trafficLimitBytes: z.ZodString;
                lifetimeTrafficUsedBytes: z.ZodString;
                username: z.ZodString;
                expiresAt: z.ZodPipe<z.ZodISODateTime, z.ZodTransform<Date, string>>;
                isActive: z.ZodBoolean;
                userStatus: z.ZodEnum<{
                    readonly ACTIVE: "ACTIVE";
                    readonly DISABLED: "DISABLED";
                    readonly LIMITED: "LIMITED";
                    readonly EXPIRED: "EXPIRED";
                }>;
                trafficLimitStrategy: z.ZodEnum<{
                    readonly NO_RESET: "NO_RESET";
                    readonly DAY: "DAY";
                    readonly WEEK: "WEEK";
                    readonly MONTH: "MONTH";
                    readonly MONTH_ROLLING: "MONTH_ROLLING";
                }>;
            }, z.core.$strip>;
            links: z.ZodArray<z.ZodString>;
            ssConfLinks: z.ZodRecord<z.ZodString, z.ZodString>;
            subscriptionUrl: z.ZodString;
        }, z.core.$strip>;
    }, z.core.$strip>;
    type RequestParam = z.infer<typeof RequestParamSchema>;
    type Response = z.infer<typeof ResponseSchema>;
}
//# sourceMappingURL=get-subscription-info-by-short-uuid.command.d.ts.map