import { z } from 'zod';
export declare namespace GetAllUsersCommand {
    const url: "/api/users/";
    const TSQ_url: "/api/users/";
    const endpointDetails: import("../../constants").EndpointDetails;
    const RequestQuerySchema: z.ZodObject<{
        start: z.ZodDefault<z.ZodNumber>;
        size: z.ZodDefault<z.ZodNumber>;
        filters: z.ZodOptional<z.ZodEffects<z.ZodArray<z.ZodObject<{
            id: z.ZodString;
            value: z.ZodUnknown;
        }, "strip", z.ZodTypeAny, {
            id: string;
            value?: unknown;
        }, {
            id: string;
            value?: unknown;
        }>, "many">, {
            id: string;
            value?: unknown;
        }[], unknown>>;
        filterModes: z.ZodOptional<z.ZodEffects<z.ZodRecord<z.ZodString, z.ZodString>, Record<string, string>, unknown>>;
        globalFilterMode: z.ZodOptional<z.ZodString>;
        sorting: z.ZodOptional<z.ZodEffects<z.ZodArray<z.ZodObject<{
            id: z.ZodString;
            desc: z.ZodBoolean;
        }, "strip", z.ZodTypeAny, {
            id: string;
            desc: boolean;
        }, {
            id: string;
            desc: boolean;
        }>, "many">, {
            id: string;
            desc: boolean;
        }[], unknown>>;
    }, "strip", z.ZodTypeAny, {
        start: number;
        size: number;
        filters?: {
            id: string;
            value?: unknown;
        }[] | undefined;
        filterModes?: Record<string, string> | undefined;
        globalFilterMode?: string | undefined;
        sorting?: {
            id: string;
            desc: boolean;
        }[] | undefined;
    }, {
        start?: number | undefined;
        size?: number | undefined;
        filters?: unknown;
        filterModes?: unknown;
        globalFilterMode?: string | undefined;
        sorting?: unknown;
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
                subLastUserAgent: z.ZodNullable<z.ZodString>;
                subLastOpenedAt: z.ZodNullable<z.ZodEffects<z.ZodString, Date, string>>;
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
                createdAt: Date;
                updatedAt: Date;
                username: string;
                id: number;
                tag: string | null;
                shortUuid: string;
                trafficLimitBytes: number;
                description: string | null;
                trafficLimitStrategy: "MONTH" | "NO_RESET" | "DAY" | "WEEK";
                expireAt: Date;
                telegramId: number | null;
                email: string | null;
                hwidDeviceLimit: number | null;
                externalSquadUuid: string | null;
                trojanPassword: string;
                vlessUuid: string;
                ssPassword: string;
                lastTriggeredThreshold: number;
                subRevokedAt: Date | null;
                subLastUserAgent: string | null;
                subLastOpenedAt: Date | null;
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
                createdAt: string;
                updatedAt: string;
                username: string;
                id: number;
                tag: string | null;
                shortUuid: string;
                description: string | null;
                expireAt: string;
                telegramId: number | null;
                email: string | null;
                hwidDeviceLimit: number | null;
                externalSquadUuid: string | null;
                trojanPassword: string;
                vlessUuid: string;
                ssPassword: string;
                subRevokedAt: string | null;
                subLastUserAgent: string | null;
                subLastOpenedAt: string | null;
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
                trafficLimitStrategy?: "MONTH" | "NO_RESET" | "DAY" | "WEEK" | undefined;
                lastTriggeredThreshold?: number | undefined;
            }>, "many">;
            total: z.ZodNumber;
        }, "strip", z.ZodTypeAny, {
            users: {
                status: "DISABLED" | "LIMITED" | "EXPIRED" | "ACTIVE";
                uuid: string;
                createdAt: Date;
                updatedAt: Date;
                username: string;
                id: number;
                tag: string | null;
                shortUuid: string;
                trafficLimitBytes: number;
                description: string | null;
                trafficLimitStrategy: "MONTH" | "NO_RESET" | "DAY" | "WEEK";
                expireAt: Date;
                telegramId: number | null;
                email: string | null;
                hwidDeviceLimit: number | null;
                externalSquadUuid: string | null;
                trojanPassword: string;
                vlessUuid: string;
                ssPassword: string;
                lastTriggeredThreshold: number;
                subRevokedAt: Date | null;
                subLastUserAgent: string | null;
                subLastOpenedAt: Date | null;
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
            total: number;
        }, {
            users: {
                uuid: string;
                createdAt: string;
                updatedAt: string;
                username: string;
                id: number;
                tag: string | null;
                shortUuid: string;
                description: string | null;
                expireAt: string;
                telegramId: number | null;
                email: string | null;
                hwidDeviceLimit: number | null;
                externalSquadUuid: string | null;
                trojanPassword: string;
                vlessUuid: string;
                ssPassword: string;
                subRevokedAt: string | null;
                subLastUserAgent: string | null;
                subLastOpenedAt: string | null;
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
                trafficLimitStrategy?: "MONTH" | "NO_RESET" | "DAY" | "WEEK" | undefined;
                lastTriggeredThreshold?: number | undefined;
            }[];
            total: number;
        }>;
    }, "strip", z.ZodTypeAny, {
        response: {
            users: {
                status: "DISABLED" | "LIMITED" | "EXPIRED" | "ACTIVE";
                uuid: string;
                createdAt: Date;
                updatedAt: Date;
                username: string;
                id: number;
                tag: string | null;
                shortUuid: string;
                trafficLimitBytes: number;
                description: string | null;
                trafficLimitStrategy: "MONTH" | "NO_RESET" | "DAY" | "WEEK";
                expireAt: Date;
                telegramId: number | null;
                email: string | null;
                hwidDeviceLimit: number | null;
                externalSquadUuid: string | null;
                trojanPassword: string;
                vlessUuid: string;
                ssPassword: string;
                lastTriggeredThreshold: number;
                subRevokedAt: Date | null;
                subLastUserAgent: string | null;
                subLastOpenedAt: Date | null;
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
            total: number;
        };
    }, {
        response: {
            users: {
                uuid: string;
                createdAt: string;
                updatedAt: string;
                username: string;
                id: number;
                tag: string | null;
                shortUuid: string;
                description: string | null;
                expireAt: string;
                telegramId: number | null;
                email: string | null;
                hwidDeviceLimit: number | null;
                externalSquadUuid: string | null;
                trojanPassword: string;
                vlessUuid: string;
                ssPassword: string;
                subRevokedAt: string | null;
                subLastUserAgent: string | null;
                subLastOpenedAt: string | null;
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
                trafficLimitStrategy?: "MONTH" | "NO_RESET" | "DAY" | "WEEK" | undefined;
                lastTriggeredThreshold?: number | undefined;
            }[];
            total: number;
        };
    }>;
    type Response = z.infer<typeof ResponseSchema>;
}
//# sourceMappingURL=get-all-users.command.d.ts.map