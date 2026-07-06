import { z } from 'zod';
export declare namespace GetUsersStreamCommand {
    const url: "/api/users/stream";
    const TSQ_url: "/api/users/stream";
    const endpointDetails: import("../../constants").EndpointDetails;
    const RequestQuerySchema: z.ZodObject<{
        cursor: z.ZodOptional<z.ZodString>;
        size: z.ZodDefault<z.ZodNumber>;
    }, "strip", z.ZodTypeAny, {
        size: number;
        cursor?: string | undefined;
    }, {
        size?: number | undefined;
        cursor?: string | undefined;
    }>;
    type RequestQuery = z.infer<typeof RequestQuerySchema>;
    const ResponseSchema: z.ZodObject<{
        response: z.ZodObject<{
            users: z.ZodArray<z.ZodObject<{
                uuid: z.ZodString;
                id: z.ZodNumber;
                shortUuid: z.ZodString;
                username: z.ZodString;
                status: z.ZodDefault<z.ZodNativeEnum<{
                    readonly ACTIVE: "ACTIVE";
                    readonly DISABLED: "DISABLED";
                    readonly LIMITED: "LIMITED";
                    readonly EXPIRED: "EXPIRED";
                }>>;
                trafficLimitBytes: z.ZodDefault<z.ZodNumber>;
                trafficLimitStrategy: z.ZodDefault<z.ZodNativeEnum<{
                    readonly NO_RESET: "NO_RESET";
                    readonly DAY: "DAY";
                    readonly WEEK: "WEEK";
                    readonly MONTH: "MONTH";
                    readonly MONTH_ROLLING: "MONTH_ROLLING";
                }>>;
                expireAt: z.ZodEffects<z.ZodString, Date, string>;
                telegramId: z.ZodNullable<z.ZodNumber>;
                email: z.ZodNullable<z.ZodString>;
                description: z.ZodNullable<z.ZodString>;
                tag: z.ZodNullable<z.ZodString>;
                hwidDeviceLimit: z.ZodNullable<z.ZodNumber>;
                externalSquadUuid: z.ZodNullable<z.ZodString>;
                trojanPassword: z.ZodString;
                vlessUuid: z.ZodString;
                ssPassword: z.ZodString;
                lastTriggeredThreshold: z.ZodDefault<z.ZodNumber>;
                subRevokedAt: z.ZodNullable<z.ZodEffects<z.ZodString, Date, string>>;
                lastTrafficResetAt: z.ZodNullable<z.ZodEffects<z.ZodString, Date, string>>;
                createdAt: z.ZodEffects<z.ZodString, Date, string>;
                updatedAt: z.ZodEffects<z.ZodString, Date, string>;
            } & {
                subscriptionUrl: z.ZodString;
                activeInternalSquads: z.ZodArray<z.ZodObject<{
                    uuid: z.ZodString;
                    name: z.ZodString;
                }, "strip", z.ZodTypeAny, {
                    uuid: string;
                    name: string;
                }, {
                    uuid: string;
                    name: string;
                }>, "many">;
                userTraffic: z.ZodObject<{
                    usedTrafficBytes: z.ZodNumber;
                    lifetimeUsedTrafficBytes: z.ZodNumber;
                    onlineAt: z.ZodNullable<z.ZodEffects<z.ZodString, Date, string>>;
                    firstConnectedAt: z.ZodNullable<z.ZodEffects<z.ZodString, Date, string>>;
                    lastConnectedNodeUuid: z.ZodNullable<z.ZodString>;
                }, "strip", z.ZodTypeAny, {
                    usedTrafficBytes: number;
                    lifetimeUsedTrafficBytes: number;
                    onlineAt: Date | null;
                    firstConnectedAt: Date | null;
                    lastConnectedNodeUuid: string | null;
                }, {
                    usedTrafficBytes: number;
                    lifetimeUsedTrafficBytes: number;
                    onlineAt: string | null;
                    firstConnectedAt: string | null;
                    lastConnectedNodeUuid: string | null;
                }>;
            }, "strip", z.ZodTypeAny, {
                status: "DISABLED" | "LIMITED" | "EXPIRED" | "ACTIVE";
                uuid: string;
                expireAt: Date;
                createdAt: Date;
                updatedAt: Date;
                description: string | null;
                username: string;
                tag: string | null;
                id: number;
                shortUuid: string;
                trafficLimitBytes: number;
                trafficLimitStrategy: "MONTH" | "NO_RESET" | "DAY" | "WEEK" | "MONTH_ROLLING";
                telegramId: number | null;
                email: string | null;
                hwidDeviceLimit: number | null;
                externalSquadUuid: string | null;
                trojanPassword: string;
                vlessUuid: string;
                ssPassword: string;
                lastTriggeredThreshold: number;
                subRevokedAt: Date | null;
                lastTrafficResetAt: Date | null;
                subscriptionUrl: string;
                activeInternalSquads: {
                    uuid: string;
                    name: string;
                }[];
                userTraffic: {
                    usedTrafficBytes: number;
                    lifetimeUsedTrafficBytes: number;
                    onlineAt: Date | null;
                    firstConnectedAt: Date | null;
                    lastConnectedNodeUuid: string | null;
                };
            }, {
                uuid: string;
                expireAt: string;
                createdAt: string;
                updatedAt: string;
                description: string | null;
                username: string;
                tag: string | null;
                id: number;
                shortUuid: string;
                telegramId: number | null;
                email: string | null;
                hwidDeviceLimit: number | null;
                externalSquadUuid: string | null;
                trojanPassword: string;
                vlessUuid: string;
                ssPassword: string;
                subRevokedAt: string | null;
                lastTrafficResetAt: string | null;
                subscriptionUrl: string;
                activeInternalSquads: {
                    uuid: string;
                    name: string;
                }[];
                userTraffic: {
                    usedTrafficBytes: number;
                    lifetimeUsedTrafficBytes: number;
                    onlineAt: string | null;
                    firstConnectedAt: string | null;
                    lastConnectedNodeUuid: string | null;
                };
                status?: "DISABLED" | "LIMITED" | "EXPIRED" | "ACTIVE" | undefined;
                trafficLimitBytes?: number | undefined;
                trafficLimitStrategy?: "MONTH" | "NO_RESET" | "DAY" | "WEEK" | "MONTH_ROLLING" | undefined;
                lastTriggeredThreshold?: number | undefined;
            }>, "many">;
            nextCursor: z.ZodNullable<z.ZodString>;
            hasMore: z.ZodBoolean;
        }, "strip", z.ZodTypeAny, {
            users: {
                status: "DISABLED" | "LIMITED" | "EXPIRED" | "ACTIVE";
                uuid: string;
                expireAt: Date;
                createdAt: Date;
                updatedAt: Date;
                description: string | null;
                username: string;
                tag: string | null;
                id: number;
                shortUuid: string;
                trafficLimitBytes: number;
                trafficLimitStrategy: "MONTH" | "NO_RESET" | "DAY" | "WEEK" | "MONTH_ROLLING";
                telegramId: number | null;
                email: string | null;
                hwidDeviceLimit: number | null;
                externalSquadUuid: string | null;
                trojanPassword: string;
                vlessUuid: string;
                ssPassword: string;
                lastTriggeredThreshold: number;
                subRevokedAt: Date | null;
                lastTrafficResetAt: Date | null;
                subscriptionUrl: string;
                activeInternalSquads: {
                    uuid: string;
                    name: string;
                }[];
                userTraffic: {
                    usedTrafficBytes: number;
                    lifetimeUsedTrafficBytes: number;
                    onlineAt: Date | null;
                    firstConnectedAt: Date | null;
                    lastConnectedNodeUuid: string | null;
                };
            }[];
            nextCursor: string | null;
            hasMore: boolean;
        }, {
            users: {
                uuid: string;
                expireAt: string;
                createdAt: string;
                updatedAt: string;
                description: string | null;
                username: string;
                tag: string | null;
                id: number;
                shortUuid: string;
                telegramId: number | null;
                email: string | null;
                hwidDeviceLimit: number | null;
                externalSquadUuid: string | null;
                trojanPassword: string;
                vlessUuid: string;
                ssPassword: string;
                subRevokedAt: string | null;
                lastTrafficResetAt: string | null;
                subscriptionUrl: string;
                activeInternalSquads: {
                    uuid: string;
                    name: string;
                }[];
                userTraffic: {
                    usedTrafficBytes: number;
                    lifetimeUsedTrafficBytes: number;
                    onlineAt: string | null;
                    firstConnectedAt: string | null;
                    lastConnectedNodeUuid: string | null;
                };
                status?: "DISABLED" | "LIMITED" | "EXPIRED" | "ACTIVE" | undefined;
                trafficLimitBytes?: number | undefined;
                trafficLimitStrategy?: "MONTH" | "NO_RESET" | "DAY" | "WEEK" | "MONTH_ROLLING" | undefined;
                lastTriggeredThreshold?: number | undefined;
            }[];
            nextCursor: string | null;
            hasMore: boolean;
        }>;
    }, "strip", z.ZodTypeAny, {
        response: {
            users: {
                status: "DISABLED" | "LIMITED" | "EXPIRED" | "ACTIVE";
                uuid: string;
                expireAt: Date;
                createdAt: Date;
                updatedAt: Date;
                description: string | null;
                username: string;
                tag: string | null;
                id: number;
                shortUuid: string;
                trafficLimitBytes: number;
                trafficLimitStrategy: "MONTH" | "NO_RESET" | "DAY" | "WEEK" | "MONTH_ROLLING";
                telegramId: number | null;
                email: string | null;
                hwidDeviceLimit: number | null;
                externalSquadUuid: string | null;
                trojanPassword: string;
                vlessUuid: string;
                ssPassword: string;
                lastTriggeredThreshold: number;
                subRevokedAt: Date | null;
                lastTrafficResetAt: Date | null;
                subscriptionUrl: string;
                activeInternalSquads: {
                    uuid: string;
                    name: string;
                }[];
                userTraffic: {
                    usedTrafficBytes: number;
                    lifetimeUsedTrafficBytes: number;
                    onlineAt: Date | null;
                    firstConnectedAt: Date | null;
                    lastConnectedNodeUuid: string | null;
                };
            }[];
            nextCursor: string | null;
            hasMore: boolean;
        };
    }, {
        response: {
            users: {
                uuid: string;
                expireAt: string;
                createdAt: string;
                updatedAt: string;
                description: string | null;
                username: string;
                tag: string | null;
                id: number;
                shortUuid: string;
                telegramId: number | null;
                email: string | null;
                hwidDeviceLimit: number | null;
                externalSquadUuid: string | null;
                trojanPassword: string;
                vlessUuid: string;
                ssPassword: string;
                subRevokedAt: string | null;
                lastTrafficResetAt: string | null;
                subscriptionUrl: string;
                activeInternalSquads: {
                    uuid: string;
                    name: string;
                }[];
                userTraffic: {
                    usedTrafficBytes: number;
                    lifetimeUsedTrafficBytes: number;
                    onlineAt: string | null;
                    firstConnectedAt: string | null;
                    lastConnectedNodeUuid: string | null;
                };
                status?: "DISABLED" | "LIMITED" | "EXPIRED" | "ACTIVE" | undefined;
                trafficLimitBytes?: number | undefined;
                trafficLimitStrategy?: "MONTH" | "NO_RESET" | "DAY" | "WEEK" | "MONTH_ROLLING" | undefined;
                lastTriggeredThreshold?: number | undefined;
            }[];
            nextCursor: string | null;
            hasMore: boolean;
        };
    }>;
    type Response = z.infer<typeof ResponseSchema>;
}
//# sourceMappingURL=get-users-stream.command.d.ts.map