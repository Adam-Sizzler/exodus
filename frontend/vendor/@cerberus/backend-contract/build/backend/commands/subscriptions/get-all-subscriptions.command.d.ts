import { z } from 'zod';
export declare namespace GetAllSubscriptionsCommand {
    const url: "/api/subscriptions/";
    const TSQ_url: "/api/subscriptions/";
    const endpointDetails: import("../../constants").EndpointDetails;
    const RequestQuerySchema: z.ZodObject<{
        start: z.ZodDefault<z.ZodNumber>;
        size: z.ZodDefault<z.ZodNumber>;
    }, "strip", z.ZodTypeAny, {
        start: number;
        size: number;
    }, {
        start?: number | undefined;
        size?: number | undefined;
    }>;
    type RequestQuery = z.infer<typeof RequestQuerySchema>;
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
                    expiresAt: z.ZodEffects<z.ZodString, Date, string>;
                    isActive: z.ZodBoolean;
                    userStatus: z.ZodNativeEnum<{
                        readonly ACTIVE: "ACTIVE";
                        readonly DISABLED: "DISABLED";
                        readonly LIMITED: "LIMITED";
                        readonly EXPIRED: "EXPIRED";
                    }>;
                    trafficLimitStrategy: z.ZodNativeEnum<{
                        readonly NO_RESET: "NO_RESET";
                        readonly DAY: "DAY";
                        readonly WEEK: "WEEK";
                        readonly MONTH: "MONTH";
                    }>;
                }, "strip", z.ZodTypeAny, {
                    username: string;
                    shortUuid: string;
                    trafficLimitBytes: string;
                    trafficLimitStrategy: "MONTH" | "NO_RESET" | "DAY" | "WEEK";
                    trafficUsedBytes: string;
                    daysLeft: number;
                    trafficUsed: string;
                    trafficLimit: string;
                    lifetimeTrafficUsed: string;
                    lifetimeTrafficUsedBytes: string;
                    expiresAt: Date;
                    isActive: boolean;
                    userStatus: "DISABLED" | "LIMITED" | "EXPIRED" | "ACTIVE";
                }, {
                    username: string;
                    shortUuid: string;
                    trafficLimitBytes: string;
                    trafficLimitStrategy: "MONTH" | "NO_RESET" | "DAY" | "WEEK";
                    trafficUsedBytes: string;
                    daysLeft: number;
                    trafficUsed: string;
                    trafficLimit: string;
                    lifetimeTrafficUsed: string;
                    lifetimeTrafficUsedBytes: string;
                    expiresAt: string;
                    isActive: boolean;
                    userStatus: "DISABLED" | "LIMITED" | "EXPIRED" | "ACTIVE";
                }>;
                links: z.ZodArray<z.ZodString, "many">;
                ssConfLinks: z.ZodRecord<z.ZodString, z.ZodString>;
                subscriptionUrl: z.ZodString;
            }, "strip", z.ZodTypeAny, {
                user: {
                    username: string;
                    shortUuid: string;
                    trafficLimitBytes: string;
                    trafficLimitStrategy: "MONTH" | "NO_RESET" | "DAY" | "WEEK";
                    trafficUsedBytes: string;
                    daysLeft: number;
                    trafficUsed: string;
                    trafficLimit: string;
                    lifetimeTrafficUsed: string;
                    lifetimeTrafficUsedBytes: string;
                    expiresAt: Date;
                    isActive: boolean;
                    userStatus: "DISABLED" | "LIMITED" | "EXPIRED" | "ACTIVE";
                };
                subscriptionUrl: string;
                isFound: boolean;
                links: string[];
                ssConfLinks: Record<string, string>;
            }, {
                user: {
                    username: string;
                    shortUuid: string;
                    trafficLimitBytes: string;
                    trafficLimitStrategy: "MONTH" | "NO_RESET" | "DAY" | "WEEK";
                    trafficUsedBytes: string;
                    daysLeft: number;
                    trafficUsed: string;
                    trafficLimit: string;
                    lifetimeTrafficUsed: string;
                    lifetimeTrafficUsedBytes: string;
                    expiresAt: string;
                    isActive: boolean;
                    userStatus: "DISABLED" | "LIMITED" | "EXPIRED" | "ACTIVE";
                };
                subscriptionUrl: string;
                isFound: boolean;
                links: string[];
                ssConfLinks: Record<string, string>;
            }>, "many">;
            total: z.ZodNumber;
        }, "strip", z.ZodTypeAny, {
            subscriptions: {
                user: {
                    username: string;
                    shortUuid: string;
                    trafficLimitBytes: string;
                    trafficLimitStrategy: "MONTH" | "NO_RESET" | "DAY" | "WEEK";
                    trafficUsedBytes: string;
                    daysLeft: number;
                    trafficUsed: string;
                    trafficLimit: string;
                    lifetimeTrafficUsed: string;
                    lifetimeTrafficUsedBytes: string;
                    expiresAt: Date;
                    isActive: boolean;
                    userStatus: "DISABLED" | "LIMITED" | "EXPIRED" | "ACTIVE";
                };
                subscriptionUrl: string;
                isFound: boolean;
                links: string[];
                ssConfLinks: Record<string, string>;
            }[];
            total: number;
        }, {
            subscriptions: {
                user: {
                    username: string;
                    shortUuid: string;
                    trafficLimitBytes: string;
                    trafficLimitStrategy: "MONTH" | "NO_RESET" | "DAY" | "WEEK";
                    trafficUsedBytes: string;
                    daysLeft: number;
                    trafficUsed: string;
                    trafficLimit: string;
                    lifetimeTrafficUsed: string;
                    lifetimeTrafficUsedBytes: string;
                    expiresAt: string;
                    isActive: boolean;
                    userStatus: "DISABLED" | "LIMITED" | "EXPIRED" | "ACTIVE";
                };
                subscriptionUrl: string;
                isFound: boolean;
                links: string[];
                ssConfLinks: Record<string, string>;
            }[];
            total: number;
        }>;
    }, "strip", z.ZodTypeAny, {
        response: {
            subscriptions: {
                user: {
                    username: string;
                    shortUuid: string;
                    trafficLimitBytes: string;
                    trafficLimitStrategy: "MONTH" | "NO_RESET" | "DAY" | "WEEK";
                    trafficUsedBytes: string;
                    daysLeft: number;
                    trafficUsed: string;
                    trafficLimit: string;
                    lifetimeTrafficUsed: string;
                    lifetimeTrafficUsedBytes: string;
                    expiresAt: Date;
                    isActive: boolean;
                    userStatus: "DISABLED" | "LIMITED" | "EXPIRED" | "ACTIVE";
                };
                subscriptionUrl: string;
                isFound: boolean;
                links: string[];
                ssConfLinks: Record<string, string>;
            }[];
            total: number;
        };
    }, {
        response: {
            subscriptions: {
                user: {
                    username: string;
                    shortUuid: string;
                    trafficLimitBytes: string;
                    trafficLimitStrategy: "MONTH" | "NO_RESET" | "DAY" | "WEEK";
                    trafficUsedBytes: string;
                    daysLeft: number;
                    trafficUsed: string;
                    trafficLimit: string;
                    lifetimeTrafficUsed: string;
                    lifetimeTrafficUsedBytes: string;
                    expiresAt: string;
                    isActive: boolean;
                    userStatus: "DISABLED" | "LIMITED" | "EXPIRED" | "ACTIVE";
                };
                subscriptionUrl: string;
                isFound: boolean;
                links: string[];
                ssConfLinks: Record<string, string>;
            }[];
            total: number;
        };
    }>;
    type Response = z.infer<typeof ResponseSchema>;
}
//# sourceMappingURL=get-all-subscriptions.command.d.ts.map