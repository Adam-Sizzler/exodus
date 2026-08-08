import { z } from 'zod';
export declare namespace GetSubscriptionsCommand {
    const url: "/api/subscriptions/";
    const TSQ_url: "/api/subscriptions/";
    const endpointDetails: import("../../constants").EndpointDetails;
    const RequestQuerySchema: z.ZodObject<{
        start: z.ZodDefault<z.ZodCoercedNumber<unknown>>;
        size: z.ZodDefault<z.ZodCoercedNumber<unknown>>;
    }, z.core.$strip>;
    const ResponseSchema: z.ZodObject<{
        response: z.ZodObject<{
            subscriptions: z.ZodArray<z.ZodObject<{
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
            }, z.core.$strip>>;
            total: z.ZodNumber;
        }, z.core.$strip>;
    }, z.core.$strip>;
    type RequestQuery = z.infer<typeof RequestQuerySchema>;
    type Response = z.infer<typeof ResponseSchema>;
}
//# sourceMappingURL=get-subscriptions.command.d.ts.map