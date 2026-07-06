import z from 'zod';
export declare const ExodusWebhookUserEvents: z.ZodObject<{
    scope: z.ZodLiteral<"user">;
    event: z.ZodEnum<["user.created" | "user.modified" | "user.deleted" | "user.revoked" | "user.disabled" | "user.enabled" | "user.limited" | "user.expired" | "user.traffic_reset" | "user.first_connected" | "user.bandwidth_usage_threshold_reached" | "user.not_connected" | "user.expiration", ...("user.created" | "user.modified" | "user.deleted" | "user.revoked" | "user.disabled" | "user.enabled" | "user.limited" | "user.expired" | "user.traffic_reset" | "user.first_connected" | "user.bandwidth_usage_threshold_reached" | "user.not_connected" | "user.expiration")[]]>;
    timestamp: z.ZodEffects<z.ZodString, Date, string>;
    data: z.ZodObject<{
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
    }>;
    meta: z.ZodNullable<z.ZodObject<{
        notConnectedAfterHours: z.ZodOptional<z.ZodNullable<z.ZodNumber>>;
        expiration: z.ZodOptional<z.ZodNullable<z.ZodNumber>>;
    }, "strip", z.ZodTypeAny, {
        notConnectedAfterHours?: number | null | undefined;
        expiration?: number | null | undefined;
    }, {
        notConnectedAfterHours?: number | null | undefined;
        expiration?: number | null | undefined;
    }>>;
}, "strip", z.ZodTypeAny, {
    data: {
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
    };
    scope: "user";
    event: "user.created" | "user.modified" | "user.deleted" | "user.revoked" | "user.disabled" | "user.enabled" | "user.limited" | "user.expired" | "user.traffic_reset" | "user.first_connected" | "user.bandwidth_usage_threshold_reached" | "user.not_connected" | "user.expiration";
    timestamp: Date;
    meta: {
        notConnectedAfterHours?: number | null | undefined;
        expiration?: number | null | undefined;
    } | null;
}, {
    data: {
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
    };
    scope: "user";
    event: "user.created" | "user.modified" | "user.deleted" | "user.revoked" | "user.disabled" | "user.enabled" | "user.limited" | "user.expired" | "user.traffic_reset" | "user.first_connected" | "user.bandwidth_usage_threshold_reached" | "user.not_connected" | "user.expiration";
    timestamp: string;
    meta: {
        notConnectedAfterHours?: number | null | undefined;
        expiration?: number | null | undefined;
    } | null;
}>;
export declare const ExodusWebhookUserHwidDevicesEvents: z.ZodObject<{
    scope: z.ZodLiteral<"user_hwid_devices">;
    event: z.ZodEnum<["user_hwid_devices.added" | "user_hwid_devices.deleted", ...("user_hwid_devices.added" | "user_hwid_devices.deleted")[]]>;
    timestamp: z.ZodEffects<z.ZodString, Date, string>;
    data: z.ZodObject<{
        user: z.ZodObject<{
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
        }>;
        hwidUserDevice: z.ZodObject<{
            hwid: z.ZodString;
            userId: z.ZodNumber;
            platform: z.ZodNullable<z.ZodString>;
            osVersion: z.ZodNullable<z.ZodString>;
            deviceModel: z.ZodNullable<z.ZodString>;
            userAgent: z.ZodNullable<z.ZodString>;
            requestIp: z.ZodNullable<z.ZodString>;
            createdAt: z.ZodEffects<z.ZodString, Date, string>;
            updatedAt: z.ZodEffects<z.ZodString, Date, string>;
        }, "strip", z.ZodTypeAny, {
            hwid: string;
            createdAt: Date;
            updatedAt: Date;
            userId: number;
            platform: string | null;
            osVersion: string | null;
            deviceModel: string | null;
            userAgent: string | null;
            requestIp: string | null;
        }, {
            hwid: string;
            createdAt: string;
            updatedAt: string;
            userId: number;
            platform: string | null;
            osVersion: string | null;
            deviceModel: string | null;
            userAgent: string | null;
            requestIp: string | null;
        }>;
    }, "strip", z.ZodTypeAny, {
        user: {
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
        };
        hwidUserDevice: {
            hwid: string;
            createdAt: Date;
            updatedAt: Date;
            userId: number;
            platform: string | null;
            osVersion: string | null;
            deviceModel: string | null;
            userAgent: string | null;
            requestIp: string | null;
        };
    }, {
        user: {
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
        };
        hwidUserDevice: {
            hwid: string;
            createdAt: string;
            updatedAt: string;
            userId: number;
            platform: string | null;
            osVersion: string | null;
            deviceModel: string | null;
            userAgent: string | null;
            requestIp: string | null;
        };
    }>;
}, "strip", z.ZodTypeAny, {
    data: {
        user: {
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
        };
        hwidUserDevice: {
            hwid: string;
            createdAt: Date;
            updatedAt: Date;
            userId: number;
            platform: string | null;
            osVersion: string | null;
            deviceModel: string | null;
            userAgent: string | null;
            requestIp: string | null;
        };
    };
    scope: "user_hwid_devices";
    event: "user_hwid_devices.added" | "user_hwid_devices.deleted";
    timestamp: Date;
}, {
    data: {
        user: {
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
        };
        hwidUserDevice: {
            hwid: string;
            createdAt: string;
            updatedAt: string;
            userId: number;
            platform: string | null;
            osVersion: string | null;
            deviceModel: string | null;
            userAgent: string | null;
            requestIp: string | null;
        };
    };
    scope: "user_hwid_devices";
    event: "user_hwid_devices.added" | "user_hwid_devices.deleted";
    timestamp: string;
}>;
export declare const ExodusWebhookNodeEvents: z.ZodObject<{
    scope: z.ZodLiteral<"node">;
    event: z.ZodEnum<["node.created" | "node.modified" | "node.disabled" | "node.enabled" | "node.deleted" | "node.connection_lost" | "node.connection_restored" | "node.traffic_notify", ...("node.created" | "node.modified" | "node.disabled" | "node.enabled" | "node.deleted" | "node.connection_lost" | "node.connection_restored" | "node.traffic_notify")[]]>;
    timestamp: z.ZodEffects<z.ZodString, Date, string>;
    data: z.ZodObject<{
        uuid: z.ZodString;
        name: z.ZodString;
        address: z.ZodString;
        port: z.ZodNullable<z.ZodNumber>;
        proxyUrl: z.ZodNullable<z.ZodString>;
        isConnected: z.ZodBoolean;
        isDisabled: z.ZodBoolean;
        isConnecting: z.ZodBoolean;
        lastStatusChange: z.ZodNullable<z.ZodEffects<z.ZodString, Date, string>>;
            lastStatusMessage: z.ZodNullable<z.ZodString>;

            singboxVersion: z.ZodNullable<z.ZodString>;

            nodeVersion: z.ZodNullable<z.ZodString>;
        isTrafficTrackingActive: z.ZodBoolean;
        trafficResetDay: z.ZodNullable<z.ZodNumber>;
        trafficLimitBytes: z.ZodNullable<z.ZodNumber>;
        trafficUsedBytes: z.ZodNullable<z.ZodNumber>;
        notifyPercent: z.ZodNullable<z.ZodNumber>;
        viewPosition: z.ZodNumber;
        countryCode: z.ZodString;
        consumptionMultiplier: z.ZodNumber;
        nodeConsumptionMultiplier: z.ZodNumber;
        tags: z.ZodArray<z.ZodString, "many">;
        createdAt: z.ZodEffects<z.ZodString, Date, string>;
        updatedAt: z.ZodEffects<z.ZodString, Date, string>;
        configProfile: z.ZodObject<{
            activeConfigProfileUuid: z.ZodNullable<z.ZodString>;
            activeInbounds: z.ZodArray<z.ZodObject<{
                uuid: z.ZodString;
                profileUuid: z.ZodString;
                tag: z.ZodString;
                type: z.ZodString;
                network: z.ZodNullable<z.ZodString>;
                security: z.ZodNullable<z.ZodString>;
                port: z.ZodNullable<z.ZodNumber>;
                rawInbound: z.ZodNullable<z.ZodUnknown>;
            }, "strip", z.ZodTypeAny, {
                uuid: string;
                type: string;
                profileUuid: string;
                tag: string;
                network: string | null;
                security: string | null;
                port: number | null;
                rawInbound?: unknown;
            }, {
                uuid: string;
                type: string;
                profileUuid: string;
                tag: string;
                network: string | null;
                security: string | null;
                port: number | null;
                rawInbound?: unknown;
            }>, "many">;
        }, "strip", z.ZodTypeAny, {
            activeConfigProfileUuid: string | null;
            activeInbounds: {
                uuid: string;
                type: string;
                profileUuid: string;
                tag: string;
                network: string | null;
                security: string | null;
                port: number | null;
                rawInbound?: unknown;
            }[];
        }, {
            activeConfigProfileUuid: string | null;
            activeInbounds: {
                uuid: string;
                type: string;
                profileUuid: string;
                tag: string;
                network: string | null;
                security: string | null;
                port: number | null;
                rawInbound?: unknown;
            }[];
        }>;
        providerUuid: z.ZodNullable<z.ZodString>;
        provider: z.ZodNullable<z.ZodObject<{
            uuid: z.ZodString;
            name: z.ZodString;
            faviconLink: z.ZodNullable<z.ZodString>;
            loginUrl: z.ZodNullable<z.ZodString>;
            createdAt: z.ZodEffects<z.ZodString, Date, string>;
            updatedAt: z.ZodEffects<z.ZodString, Date, string>;
        }, "strip", z.ZodTypeAny, {
            uuid: string;
            name: string;
            createdAt: Date;
            updatedAt: Date;
            faviconLink: string | null;
            loginUrl: string | null;
        }, {
            uuid: string;
            name: string;
            createdAt: string;
            updatedAt: string;
            faviconLink: string | null;
            loginUrl: string | null;
        }>>;
        activePluginUuid: z.ZodNullable<z.ZodString>;
        system: z.ZodNullable<z.ZodObject<{
            info: z.ZodObject<{
                arch: z.ZodString;
                cpus: z.ZodNumber;
                cpuModel: z.ZodString;
                memoryTotal: z.ZodNumber;
                hostname: z.ZodString;
                platform: z.ZodString;
                release: z.ZodString;
                type: z.ZodString;
                version: z.ZodString;
                networkInterfaces: z.ZodArray<z.ZodString, "many">;
            }, "strip", z.ZodTypeAny, {
                type: string;
                version: string;
                platform: string;
                arch: string;
                cpus: number;
                cpuModel: string;
                memoryTotal: number;
                hostname: string;
                release: string;
                networkInterfaces: string[];
            }, {
                type: string;
                version: string;
                platform: string;
                arch: string;
                cpus: number;
                cpuModel: string;
                memoryTotal: number;
                hostname: string;
                release: string;
                networkInterfaces: string[];
            }>;
            stats: z.ZodObject<{
                memoryFree: z.ZodNumber;
                memoryUsed: z.ZodNumber;
                uptime: z.ZodNumber;
                loadAvg: z.ZodArray<z.ZodNumber, "many">;
                interface: z.ZodNullable<z.ZodObject<{
                    interface: z.ZodString;
                    rxBytesPerSec: z.ZodNumber;
                    txBytesPerSec: z.ZodNumber;
                    rxTotal: z.ZodNumber;
                    txTotal: z.ZodNumber;
                }, "strip", z.ZodTypeAny, {
                    interface: string;
                    rxBytesPerSec: number;
                    txBytesPerSec: number;
                    rxTotal: number;
                    txTotal: number;
                }, {
                    interface: string;
                    rxBytesPerSec: number;
                    txBytesPerSec: number;
                    rxTotal: number;
                    txTotal: number;
                }>>;
            }, "strip", z.ZodTypeAny, {
                interface: {
                    interface: string;
                    rxBytesPerSec: number;
                    txBytesPerSec: number;
                    rxTotal: number;
                    txTotal: number;
                } | null;
                memoryFree: number;
                memoryUsed: number;
                uptime: number;
                loadAvg: number[];
            }, {
                interface: {
                    interface: string;
                    rxBytesPerSec: number;
                    txBytesPerSec: number;
                    rxTotal: number;
                    txTotal: number;
                } | null;
                memoryFree: number;
                memoryUsed: number;
                uptime: number;
                loadAvg: number[];
            }>;
        }, "strip", z.ZodTypeAny, {
            stats: {
                interface: {
                    interface: string;
                    rxBytesPerSec: number;
                    txBytesPerSec: number;
                    rxTotal: number;
                    txTotal: number;
                } | null;
                memoryFree: number;
                memoryUsed: number;
                uptime: number;
                loadAvg: number[];
            };
            info: {
                type: string;
                version: string;
                platform: string;
                arch: string;
                cpus: number;
                cpuModel: string;
                memoryTotal: number;
                hostname: string;
                release: string;
                networkInterfaces: string[];
            };
        }, {
            stats: {
                interface: {
                    interface: string;
                    rxBytesPerSec: number;
                    txBytesPerSec: number;
                    rxTotal: number;
                    txTotal: number;
                } | null;
                memoryFree: number;
                memoryUsed: number;
                uptime: number;
                loadAvg: number[];
            };
            info: {
                type: string;
                version: string;
                platform: string;
                arch: string;
                cpus: number;
                cpuModel: string;
                memoryTotal: number;
                hostname: string;
                release: string;
                networkInterfaces: string[];
            };
        }>>;
        versions: z.ZodNullable<z.ZodObject<{
            singbox: z.ZodString;
            node: z.ZodString;
        }, "strip", z.ZodTypeAny, {
            node: string;
            singbox: string;
        }, {
            node: string;
            singbox: string;
        }>>;
        singboxUptime: z.ZodNumber;
        usersOnline: z.ZodNumber;
        note: z.ZodNullable<z.ZodString>;
    }, "strip", z.ZodTypeAny, {
        tags: string[];
        system: {
            stats: {
                interface: {
                    interface: string;
                    rxBytesPerSec: number;
                    txBytesPerSec: number;
                    rxTotal: number;
                    txTotal: number;
                } | null;
                memoryFree: number;
                memoryUsed: number;
                uptime: number;
                loadAvg: number[];
            };
            info: {
                type: string;
                version: string;
                platform: string;
                arch: string;
                cpus: number;
                cpuModel: string;
                memoryTotal: number;
                hostname: string;
                release: string;
                networkInterfaces: string[];
            };
        } | null;
        uuid: string;
        name: string;
        createdAt: Date;
        updatedAt: Date;
        provider: {
            uuid: string;
            name: string;
            createdAt: Date;
            updatedAt: Date;
            faviconLink: string | null;
            loginUrl: string | null;
        } | null;
        countryCode: string;
        port: number | null;
        viewPosition: number;
        trafficLimitBytes: number | null;
        address: string;
        isDisabled: boolean;
        proxyUrl: string | null;
        isConnected: boolean;
        isConnecting: boolean;
        lastStatusChange: Date | null;
            lastStatusMessage: string | null;

            singboxVersion: string | null;

            nodeVersion: string | null;
        isTrafficTrackingActive: boolean;
        trafficResetDay: number | null;
        trafficUsedBytes: number | null;
        notifyPercent: number | null;
        consumptionMultiplier: number;
        nodeConsumptionMultiplier: number;
        configProfile: {
            activeConfigProfileUuid: string | null;
            activeInbounds: {
                uuid: string;
                type: string;
                profileUuid: string;
                tag: string;
                network: string | null;
                security: string | null;
                port: number | null;
                rawInbound?: unknown;
            }[];
        };
        providerUuid: string | null;
        activePluginUuid: string | null;
        versions: {
            node: string;
            singbox: string;
        } | null;
        singboxUptime: number;
        usersOnline: number;
        note: string | null;
    }, {
        tags: string[];
        system: {
            stats: {
                interface: {
                    interface: string;
                    rxBytesPerSec: number;
                    txBytesPerSec: number;
                    rxTotal: number;
                    txTotal: number;
                } | null;
                memoryFree: number;
                memoryUsed: number;
                uptime: number;
                loadAvg: number[];
            };
            info: {
                type: string;
                version: string;
                platform: string;
                arch: string;
                cpus: number;
                cpuModel: string;
                memoryTotal: number;
                hostname: string;
                release: string;
                networkInterfaces: string[];
            };
        } | null;
        uuid: string;
        name: string;
        createdAt: string;
        updatedAt: string;
        provider: {
            uuid: string;
            name: string;
            createdAt: string;
            updatedAt: string;
            faviconLink: string | null;
            loginUrl: string | null;
        } | null;
        countryCode: string;
        port: number | null;
        viewPosition: number;
        trafficLimitBytes: number | null;
        address: string;
        isDisabled: boolean;
        proxyUrl: string | null;
        isConnected: boolean;
        isConnecting: boolean;
        lastStatusChange: string | null;
            lastStatusMessage: string | null;

            singboxVersion: string | null;

            nodeVersion: string | null;
        isTrafficTrackingActive: boolean;
        trafficResetDay: number | null;
        trafficUsedBytes: number | null;
        notifyPercent: number | null;
        consumptionMultiplier: number;
        nodeConsumptionMultiplier: number;
        configProfile: {
            activeConfigProfileUuid: string | null;
            activeInbounds: {
                uuid: string;
                type: string;
                profileUuid: string;
                tag: string;
                network: string | null;
                security: string | null;
                port: number | null;
                rawInbound?: unknown;
            }[];
        };
        providerUuid: string | null;
        activePluginUuid: string | null;
        versions: {
            node: string;
            singbox: string;
        } | null;
        singboxUptime: number;
        usersOnline: number;
        note: string | null;
    }>;
}, "strip", z.ZodTypeAny, {
    data: {
        tags: string[];
        system: {
            stats: {
                interface: {
                    interface: string;
                    rxBytesPerSec: number;
                    txBytesPerSec: number;
                    rxTotal: number;
                    txTotal: number;
                } | null;
                memoryFree: number;
                memoryUsed: number;
                uptime: number;
                loadAvg: number[];
            };
            info: {
                type: string;
                version: string;
                platform: string;
                arch: string;
                cpus: number;
                cpuModel: string;
                memoryTotal: number;
                hostname: string;
                release: string;
                networkInterfaces: string[];
            };
        } | null;
        uuid: string;
        name: string;
        createdAt: Date;
        updatedAt: Date;
        provider: {
            uuid: string;
            name: string;
            createdAt: Date;
            updatedAt: Date;
            faviconLink: string | null;
            loginUrl: string | null;
        } | null;
        countryCode: string;
        port: number | null;
        viewPosition: number;
        trafficLimitBytes: number | null;
        address: string;
        isDisabled: boolean;
        proxyUrl: string | null;
        isConnected: boolean;
        isConnecting: boolean;
        lastStatusChange: Date | null;
            lastStatusMessage: string | null;

            singboxVersion: string | null;

            nodeVersion: string | null;
        isTrafficTrackingActive: boolean;
        trafficResetDay: number | null;
        trafficUsedBytes: number | null;
        notifyPercent: number | null;
        consumptionMultiplier: number;
        nodeConsumptionMultiplier: number;
        configProfile: {
            activeConfigProfileUuid: string | null;
            activeInbounds: {
                uuid: string;
                type: string;
                profileUuid: string;
                tag: string;
                network: string | null;
                security: string | null;
                port: number | null;
                rawInbound?: unknown;
            }[];
        };
        providerUuid: string | null;
        activePluginUuid: string | null;
        versions: {
            node: string;
            singbox: string;
        } | null;
        singboxUptime: number;
        usersOnline: number;
        note: string | null;
    };
    scope: "node";
    event: "node.created" | "node.modified" | "node.disabled" | "node.enabled" | "node.deleted" | "node.connection_lost" | "node.connection_restored" | "node.traffic_notify";
    timestamp: Date;
}, {
    data: {
        tags: string[];
        system: {
            stats: {
                interface: {
                    interface: string;
                    rxBytesPerSec: number;
                    txBytesPerSec: number;
                    rxTotal: number;
                    txTotal: number;
                } | null;
                memoryFree: number;
                memoryUsed: number;
                uptime: number;
                loadAvg: number[];
            };
            info: {
                type: string;
                version: string;
                platform: string;
                arch: string;
                cpus: number;
                cpuModel: string;
                memoryTotal: number;
                hostname: string;
                release: string;
                networkInterfaces: string[];
            };
        } | null;
        uuid: string;
        name: string;
        createdAt: string;
        updatedAt: string;
        provider: {
            uuid: string;
            name: string;
            createdAt: string;
            updatedAt: string;
            faviconLink: string | null;
            loginUrl: string | null;
        } | null;
        countryCode: string;
        port: number | null;
        viewPosition: number;
        trafficLimitBytes: number | null;
        address: string;
        isDisabled: boolean;
        proxyUrl: string | null;
        isConnected: boolean;
        isConnecting: boolean;
        lastStatusChange: string | null;
            lastStatusMessage: string | null;

            singboxVersion: string | null;

            nodeVersion: string | null;
        isTrafficTrackingActive: boolean;
        trafficResetDay: number | null;
        trafficUsedBytes: number | null;
        notifyPercent: number | null;
        consumptionMultiplier: number;
        nodeConsumptionMultiplier: number;
        configProfile: {
            activeConfigProfileUuid: string | null;
            activeInbounds: {
                uuid: string;
                type: string;
                profileUuid: string;
                tag: string;
                network: string | null;
                security: string | null;
                port: number | null;
                rawInbound?: unknown;
            }[];
        };
        providerUuid: string | null;
        activePluginUuid: string | null;
        versions: {
            node: string;
            singbox: string;
        } | null;
        singboxUptime: number;
        usersOnline: number;
        note: string | null;
    };
    scope: "node";
    event: "node.created" | "node.modified" | "node.disabled" | "node.enabled" | "node.deleted" | "node.connection_lost" | "node.connection_restored" | "node.traffic_notify";
    timestamp: string;
}>;
export declare const ExodusWebhookServiceEvents: z.ZodObject<{
    scope: z.ZodLiteral<"service">;
    event: z.ZodEnum<["service.panel_started" | "service.login_attempt_failed" | "service.login_attempt_success" | "service.subpage_config_changed" | "service.api_token_created" | "service.api_token_deleted", ...("service.panel_started" | "service.login_attempt_failed" | "service.login_attempt_success" | "service.subpage_config_changed" | "service.api_token_created" | "service.api_token_deleted")[]]>;
    timestamp: z.ZodEffects<z.ZodString, Date, string>;
    data: z.ZodObject<{
        loginAttempt: z.ZodOptional<z.ZodObject<{
            username: z.ZodString;
            ip: z.ZodString;
            userAgent: z.ZodString;
            description: z.ZodOptional<z.ZodString>;
            password: z.ZodOptional<z.ZodString>;
        }, "strip", z.ZodTypeAny, {
            username: string;
            ip: string;
            userAgent: string;
            description?: string | undefined;
            password?: string | undefined;
        }, {
            username: string;
            ip: string;
            userAgent: string;
            description?: string | undefined;
            password?: string | undefined;
        }>>;
        panelVersion: z.ZodOptional<z.ZodString>;
        subpageConfig: z.ZodOptional<z.ZodObject<{
            action: z.ZodEnum<["CREATED" | "UPDATED" | "DELETED", ...("CREATED" | "UPDATED" | "DELETED")[]]>;
            uuid: z.ZodString;
        }, "strip", z.ZodTypeAny, {
            uuid: string;
            action: "CREATED" | "UPDATED" | "DELETED";
        }, {
            uuid: string;
            action: "CREATED" | "UPDATED" | "DELETED";
        }>>;
        apiToken: z.ZodOptional<z.ZodObject<{
            name: z.ZodString;
            uuid: z.ZodString;
            expireAt: z.ZodEffects<z.ZodString, Date, string>;
            scopes: z.ZodArray<z.ZodString, "many">;
        }, "strip", z.ZodTypeAny, {
            scopes: string[];
            uuid: string;
            name: string;
            expireAt: Date;
        }, {
            scopes: string[];
            uuid: string;
            name: string;
            expireAt: string;
        }>>;
    }, "strip", z.ZodTypeAny, {
        loginAttempt?: {
            username: string;
            ip: string;
            userAgent: string;
            description?: string | undefined;
            password?: string | undefined;
        } | undefined;
        panelVersion?: string | undefined;
        subpageConfig?: {
            uuid: string;
            action: "CREATED" | "UPDATED" | "DELETED";
        } | undefined;
        apiToken?: {
            scopes: string[];
            uuid: string;
            name: string;
            expireAt: Date;
        } | undefined;
    }, {
        loginAttempt?: {
            username: string;
            ip: string;
            userAgent: string;
            description?: string | undefined;
            password?: string | undefined;
        } | undefined;
        panelVersion?: string | undefined;
        subpageConfig?: {
            uuid: string;
            action: "CREATED" | "UPDATED" | "DELETED";
        } | undefined;
        apiToken?: {
            scopes: string[];
            uuid: string;
            name: string;
            expireAt: string;
        } | undefined;
    }>;
}, "strip", z.ZodTypeAny, {
    data: {
        loginAttempt?: {
            username: string;
            ip: string;
            userAgent: string;
            description?: string | undefined;
            password?: string | undefined;
        } | undefined;
        panelVersion?: string | undefined;
        subpageConfig?: {
            uuid: string;
            action: "CREATED" | "UPDATED" | "DELETED";
        } | undefined;
        apiToken?: {
            scopes: string[];
            uuid: string;
            name: string;
            expireAt: Date;
        } | undefined;
    };
    scope: "service";
    event: "service.panel_started" | "service.login_attempt_failed" | "service.login_attempt_success" | "service.subpage_config_changed" | "service.api_token_created" | "service.api_token_deleted";
    timestamp: Date;
}, {
    data: {
        loginAttempt?: {
            username: string;
            ip: string;
            userAgent: string;
            description?: string | undefined;
            password?: string | undefined;
        } | undefined;
        panelVersion?: string | undefined;
        subpageConfig?: {
            uuid: string;
            action: "CREATED" | "UPDATED" | "DELETED";
        } | undefined;
        apiToken?: {
            scopes: string[];
            uuid: string;
            name: string;
            expireAt: string;
        } | undefined;
    };
    scope: "service";
    event: "service.panel_started" | "service.login_attempt_failed" | "service.login_attempt_success" | "service.subpage_config_changed" | "service.api_token_created" | "service.api_token_deleted";
    timestamp: string;
}>;
export declare const ExodusWebhookErrorsEvents: z.ZodObject<{
    scope: z.ZodLiteral<"errors">;
    event: z.ZodEnum<["errors.bandwidth_usage_threshold_reached_max_notifications", ..."errors.bandwidth_usage_threshold_reached_max_notifications"[]]>;
    timestamp: z.ZodEffects<z.ZodString, Date, string>;
    data: z.ZodObject<{
        description: z.ZodString;
    }, "strip", z.ZodTypeAny, {
        description: string;
    }, {
        description: string;
    }>;
}, "strip", z.ZodTypeAny, {
    data: {
        description: string;
    };
    scope: "errors";
    event: "errors.bandwidth_usage_threshold_reached_max_notifications";
    timestamp: Date;
}, {
    data: {
        description: string;
    };
    scope: "errors";
    event: "errors.bandwidth_usage_threshold_reached_max_notifications";
    timestamp: string;
}>;
export declare const ExodusWebhookCrmEvents: z.ZodObject<{
    scope: z.ZodLiteral<"crm">;
    event: z.ZodEnum<["crm.infra_billing_node_payment_in_7_days" | "crm.infra_billing_node_payment_in_48hrs" | "crm.infra_billing_node_payment_in_24hrs" | "crm.infra_billing_node_payment_due_today" | "crm.infra_billing_node_payment_overdue_24hrs" | "crm.infra_billing_node_payment_overdue_48hrs" | "crm.infra_billing_node_payment_overdue_7_days", ...("crm.infra_billing_node_payment_in_7_days" | "crm.infra_billing_node_payment_in_48hrs" | "crm.infra_billing_node_payment_in_24hrs" | "crm.infra_billing_node_payment_due_today" | "crm.infra_billing_node_payment_overdue_24hrs" | "crm.infra_billing_node_payment_overdue_48hrs" | "crm.infra_billing_node_payment_overdue_7_days")[]]>;
    timestamp: z.ZodEffects<z.ZodString, Date, string>;
    data: z.ZodObject<{
        providerName: z.ZodString;
        nodeName: z.ZodString;
        nextBillingAt: z.ZodEffects<z.ZodString, Date, string>;
        loginUrl: z.ZodString;
    }, "strip", z.ZodTypeAny, {
        nodeName: string;
        loginUrl: string;
        nextBillingAt: Date;
        providerName: string;
    }, {
        nodeName: string;
        loginUrl: string;
        nextBillingAt: string;
        providerName: string;
    }>;
}, "strip", z.ZodTypeAny, {
    data: {
        nodeName: string;
        loginUrl: string;
        nextBillingAt: Date;
        providerName: string;
    };
    scope: "crm";
    event: "crm.infra_billing_node_payment_in_7_days" | "crm.infra_billing_node_payment_in_48hrs" | "crm.infra_billing_node_payment_in_24hrs" | "crm.infra_billing_node_payment_due_today" | "crm.infra_billing_node_payment_overdue_24hrs" | "crm.infra_billing_node_payment_overdue_48hrs" | "crm.infra_billing_node_payment_overdue_7_days";
    timestamp: Date;
}, {
    data: {
        nodeName: string;
        loginUrl: string;
        nextBillingAt: string;
        providerName: string;
    };
    scope: "crm";
    event: "crm.infra_billing_node_payment_in_7_days" | "crm.infra_billing_node_payment_in_48hrs" | "crm.infra_billing_node_payment_in_24hrs" | "crm.infra_billing_node_payment_due_today" | "crm.infra_billing_node_payment_overdue_24hrs" | "crm.infra_billing_node_payment_overdue_48hrs" | "crm.infra_billing_node_payment_overdue_7_days";
    timestamp: string;
}>;
export declare const ExodusWebhookTorrentBlockerEvents: z.ZodObject<{
    scope: z.ZodLiteral<"torrent_blocker">;
    event: z.ZodEnum<["torrent_blocker.report", ..."torrent_blocker.report"[]]>;
    timestamp: z.ZodEffects<z.ZodString, Date, string>;
    data: z.ZodObject<{
        node: z.ZodObject<{
            uuid: z.ZodString;
            name: z.ZodString;
            address: z.ZodString;
            port: z.ZodNullable<z.ZodNumber>;
            proxyUrl: z.ZodNullable<z.ZodString>;
            isConnected: z.ZodBoolean;
            isDisabled: z.ZodBoolean;
            isConnecting: z.ZodBoolean;
            lastStatusChange: z.ZodNullable<z.ZodEffects<z.ZodString, Date, string>>;
            lastStatusMessage: z.ZodNullable<z.ZodString>;

            singboxVersion: z.ZodNullable<z.ZodString>;

            nodeVersion: z.ZodNullable<z.ZodString>;
            isTrafficTrackingActive: z.ZodBoolean;
            trafficResetDay: z.ZodNullable<z.ZodNumber>;
            trafficLimitBytes: z.ZodNullable<z.ZodNumber>;
            trafficUsedBytes: z.ZodNullable<z.ZodNumber>;
            notifyPercent: z.ZodNullable<z.ZodNumber>;
            viewPosition: z.ZodNumber;
            countryCode: z.ZodString;
            consumptionMultiplier: z.ZodNumber;
            nodeConsumptionMultiplier: z.ZodNumber;
            tags: z.ZodArray<z.ZodString, "many">;
            createdAt: z.ZodEffects<z.ZodString, Date, string>;
            updatedAt: z.ZodEffects<z.ZodString, Date, string>;
            configProfile: z.ZodObject<{
                activeConfigProfileUuid: z.ZodNullable<z.ZodString>;
                activeInbounds: z.ZodArray<z.ZodObject<{
                    uuid: z.ZodString;
                    profileUuid: z.ZodString;
                    tag: z.ZodString;
                    type: z.ZodString;
                    network: z.ZodNullable<z.ZodString>;
                    security: z.ZodNullable<z.ZodString>;
                    port: z.ZodNullable<z.ZodNumber>;
                    rawInbound: z.ZodNullable<z.ZodUnknown>;
                }, "strip", z.ZodTypeAny, {
                    uuid: string;
                    type: string;
                    profileUuid: string;
                    tag: string;
                    network: string | null;
                    security: string | null;
                    port: number | null;
                    rawInbound?: unknown;
                }, {
                    uuid: string;
                    type: string;
                    profileUuid: string;
                    tag: string;
                    network: string | null;
                    security: string | null;
                    port: number | null;
                    rawInbound?: unknown;
                }>, "many">;
            }, "strip", z.ZodTypeAny, {
                activeConfigProfileUuid: string | null;
                activeInbounds: {
                    uuid: string;
                    type: string;
                    profileUuid: string;
                    tag: string;
                    network: string | null;
                    security: string | null;
                    port: number | null;
                    rawInbound?: unknown;
                }[];
            }, {
                activeConfigProfileUuid: string | null;
                activeInbounds: {
                    uuid: string;
                    type: string;
                    profileUuid: string;
                    tag: string;
                    network: string | null;
                    security: string | null;
                    port: number | null;
                    rawInbound?: unknown;
                }[];
            }>;
            providerUuid: z.ZodNullable<z.ZodString>;
            provider: z.ZodNullable<z.ZodObject<{
                uuid: z.ZodString;
                name: z.ZodString;
                faviconLink: z.ZodNullable<z.ZodString>;
                loginUrl: z.ZodNullable<z.ZodString>;
                createdAt: z.ZodEffects<z.ZodString, Date, string>;
                updatedAt: z.ZodEffects<z.ZodString, Date, string>;
            }, "strip", z.ZodTypeAny, {
                uuid: string;
                name: string;
                createdAt: Date;
                updatedAt: Date;
                faviconLink: string | null;
                loginUrl: string | null;
            }, {
                uuid: string;
                name: string;
                createdAt: string;
                updatedAt: string;
                faviconLink: string | null;
                loginUrl: string | null;
            }>>;
            activePluginUuid: z.ZodNullable<z.ZodString>;
            system: z.ZodNullable<z.ZodObject<{
                info: z.ZodObject<{
                    arch: z.ZodString;
                    cpus: z.ZodNumber;
                    cpuModel: z.ZodString;
                    memoryTotal: z.ZodNumber;
                    hostname: z.ZodString;
                    platform: z.ZodString;
                    release: z.ZodString;
                    type: z.ZodString;
                    version: z.ZodString;
                    networkInterfaces: z.ZodArray<z.ZodString, "many">;
                }, "strip", z.ZodTypeAny, {
                    type: string;
                    version: string;
                    platform: string;
                    arch: string;
                    cpus: number;
                    cpuModel: string;
                    memoryTotal: number;
                    hostname: string;
                    release: string;
                    networkInterfaces: string[];
                }, {
                    type: string;
                    version: string;
                    platform: string;
                    arch: string;
                    cpus: number;
                    cpuModel: string;
                    memoryTotal: number;
                    hostname: string;
                    release: string;
                    networkInterfaces: string[];
                }>;
                stats: z.ZodObject<{
                    memoryFree: z.ZodNumber;
                    memoryUsed: z.ZodNumber;
                    uptime: z.ZodNumber;
                    loadAvg: z.ZodArray<z.ZodNumber, "many">;
                    interface: z.ZodNullable<z.ZodObject<{
                        interface: z.ZodString;
                        rxBytesPerSec: z.ZodNumber;
                        txBytesPerSec: z.ZodNumber;
                        rxTotal: z.ZodNumber;
                        txTotal: z.ZodNumber;
                    }, "strip", z.ZodTypeAny, {
                        interface: string;
                        rxBytesPerSec: number;
                        txBytesPerSec: number;
                        rxTotal: number;
                        txTotal: number;
                    }, {
                        interface: string;
                        rxBytesPerSec: number;
                        txBytesPerSec: number;
                        rxTotal: number;
                        txTotal: number;
                    }>>;
                }, "strip", z.ZodTypeAny, {
                    interface: {
                        interface: string;
                        rxBytesPerSec: number;
                        txBytesPerSec: number;
                        rxTotal: number;
                        txTotal: number;
                    } | null;
                    memoryFree: number;
                    memoryUsed: number;
                    uptime: number;
                    loadAvg: number[];
                }, {
                    interface: {
                        interface: string;
                        rxBytesPerSec: number;
                        txBytesPerSec: number;
                        rxTotal: number;
                        txTotal: number;
                    } | null;
                    memoryFree: number;
                    memoryUsed: number;
                    uptime: number;
                    loadAvg: number[];
                }>;
            }, "strip", z.ZodTypeAny, {
                stats: {
                    interface: {
                        interface: string;
                        rxBytesPerSec: number;
                        txBytesPerSec: number;
                        rxTotal: number;
                        txTotal: number;
                    } | null;
                    memoryFree: number;
                    memoryUsed: number;
                    uptime: number;
                    loadAvg: number[];
                };
                info: {
                    type: string;
                    version: string;
                    platform: string;
                    arch: string;
                    cpus: number;
                    cpuModel: string;
                    memoryTotal: number;
                    hostname: string;
                    release: string;
                    networkInterfaces: string[];
                };
            }, {
                stats: {
                    interface: {
                        interface: string;
                        rxBytesPerSec: number;
                        txBytesPerSec: number;
                        rxTotal: number;
                        txTotal: number;
                    } | null;
                    memoryFree: number;
                    memoryUsed: number;
                    uptime: number;
                    loadAvg: number[];
                };
                info: {
                    type: string;
                    version: string;
                    platform: string;
                    arch: string;
                    cpus: number;
                    cpuModel: string;
                    memoryTotal: number;
                    hostname: string;
                    release: string;
                    networkInterfaces: string[];
                };
            }>>;
            versions: z.ZodNullable<z.ZodObject<{
                singbox: z.ZodString;
                node: z.ZodString;
            }, "strip", z.ZodTypeAny, {
                node: string;
                singbox: string;
            }, {
                node: string;
                singbox: string;
            }>>;
            singboxUptime: z.ZodNumber;
            usersOnline: z.ZodNumber;
            note: z.ZodNullable<z.ZodString>;
        }, "strip", z.ZodTypeAny, {
            tags: string[];
            system: {
                stats: {
                    interface: {
                        interface: string;
                        rxBytesPerSec: number;
                        txBytesPerSec: number;
                        rxTotal: number;
                        txTotal: number;
                    } | null;
                    memoryFree: number;
                    memoryUsed: number;
                    uptime: number;
                    loadAvg: number[];
                };
                info: {
                    type: string;
                    version: string;
                    platform: string;
                    arch: string;
                    cpus: number;
                    cpuModel: string;
                    memoryTotal: number;
                    hostname: string;
                    release: string;
                    networkInterfaces: string[];
                };
            } | null;
            uuid: string;
            name: string;
            createdAt: Date;
            updatedAt: Date;
            provider: {
                uuid: string;
                name: string;
                createdAt: Date;
                updatedAt: Date;
                faviconLink: string | null;
                loginUrl: string | null;
            } | null;
            countryCode: string;
            port: number | null;
            viewPosition: number;
            trafficLimitBytes: number | null;
            address: string;
            isDisabled: boolean;
            proxyUrl: string | null;
            isConnected: boolean;
            isConnecting: boolean;
            lastStatusChange: Date | null;
            lastStatusMessage: string | null;

            singboxVersion: string | null;

            nodeVersion: string | null;
            isTrafficTrackingActive: boolean;
            trafficResetDay: number | null;
            trafficUsedBytes: number | null;
            notifyPercent: number | null;
            consumptionMultiplier: number;
            nodeConsumptionMultiplier: number;
            configProfile: {
                activeConfigProfileUuid: string | null;
                activeInbounds: {
                    uuid: string;
                    type: string;
                    profileUuid: string;
                    tag: string;
                    network: string | null;
                    security: string | null;
                    port: number | null;
                    rawInbound?: unknown;
                }[];
            };
            providerUuid: string | null;
            activePluginUuid: string | null;
            versions: {
                node: string;
                singbox: string;
            } | null;
            singboxUptime: number;
            usersOnline: number;
            note: string | null;
        }, {
            tags: string[];
            system: {
                stats: {
                    interface: {
                        interface: string;
                        rxBytesPerSec: number;
                        txBytesPerSec: number;
                        rxTotal: number;
                        txTotal: number;
                    } | null;
                    memoryFree: number;
                    memoryUsed: number;
                    uptime: number;
                    loadAvg: number[];
                };
                info: {
                    type: string;
                    version: string;
                    platform: string;
                    arch: string;
                    cpus: number;
                    cpuModel: string;
                    memoryTotal: number;
                    hostname: string;
                    release: string;
                    networkInterfaces: string[];
                };
            } | null;
            uuid: string;
            name: string;
            createdAt: string;
            updatedAt: string;
            provider: {
                uuid: string;
                name: string;
                createdAt: string;
                updatedAt: string;
                faviconLink: string | null;
                loginUrl: string | null;
            } | null;
            countryCode: string;
            port: number | null;
            viewPosition: number;
            trafficLimitBytes: number | null;
            address: string;
            isDisabled: boolean;
            proxyUrl: string | null;
            isConnected: boolean;
            isConnecting: boolean;
            lastStatusChange: string | null;
            lastStatusMessage: string | null;

            singboxVersion: string | null;

            nodeVersion: string | null;
            isTrafficTrackingActive: boolean;
            trafficResetDay: number | null;
            trafficUsedBytes: number | null;
            notifyPercent: number | null;
            consumptionMultiplier: number;
            nodeConsumptionMultiplier: number;
            configProfile: {
                activeConfigProfileUuid: string | null;
                activeInbounds: {
                    uuid: string;
                    type: string;
                    profileUuid: string;
                    tag: string;
                    network: string | null;
                    security: string | null;
                    port: number | null;
                    rawInbound?: unknown;
                }[];
            };
            providerUuid: string | null;
            activePluginUuid: string | null;
            versions: {
                node: string;
                singbox: string;
            } | null;
            singboxUptime: number;
            usersOnline: number;
            note: string | null;
        }>;
        user: z.ZodObject<{
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
        }>;
        report: z.ZodObject<{
            actionReport: z.ZodObject<{
                blocked: z.ZodBoolean;
                ip: z.ZodString;
                blockDuration: z.ZodNumber;
                willUnblockAt: z.ZodEffects<z.ZodString, Date, string>;
                userId: z.ZodString;
                processedAt: z.ZodEffects<z.ZodString, Date, string>;
            }, "strip", z.ZodTypeAny, {
                ip: string;
                userId: string;
                blocked: boolean;
                blockDuration: number;
                willUnblockAt: Date;
                processedAt: Date;
            }, {
                ip: string;
                userId: string;
                blocked: boolean;
                blockDuration: number;
                willUnblockAt: string;
                processedAt: string;
            }>;
            xrayReport: z.ZodObject<{
                email: z.ZodNullable<z.ZodString>;
                level: z.ZodNullable<z.ZodNumber>;
                protocol: z.ZodNullable<z.ZodString>;
                network: z.ZodString;
                source: z.ZodNullable<z.ZodString>;
                destination: z.ZodString;
                routeTarget: z.ZodNullable<z.ZodString>;
                originalTarget: z.ZodNullable<z.ZodString>;
                inboundTag: z.ZodNullable<z.ZodString>;
                inboundName: z.ZodNullable<z.ZodString>;
                inboundLocal: z.ZodNullable<z.ZodString>;
                outboundTag: z.ZodNullable<z.ZodString>;
                ts: z.ZodNumber;
            }, "strip", z.ZodTypeAny, {
                network: string;
                email: string | null;
                protocol: string | null;
                inboundTag: string | null;
                level: number | null;
                source: string | null;
                destination: string;
                routeTarget: string | null;
                originalTarget: string | null;
                inboundName: string | null;
                inboundLocal: string | null;
                outboundTag: string | null;
                ts: number;
            }, {
                network: string;
                email: string | null;
                protocol: string | null;
                inboundTag: string | null;
                level: number | null;
                source: string | null;
                destination: string;
                routeTarget: string | null;
                originalTarget: string | null;
                inboundName: string | null;
                inboundLocal: string | null;
                outboundTag: string | null;
                ts: number;
            }>;
        }, "strip", z.ZodTypeAny, {
            actionReport: {
                ip: string;
                userId: string;
                blocked: boolean;
                blockDuration: number;
                willUnblockAt: Date;
                processedAt: Date;
            };
            xrayReport: {
                network: string;
                email: string | null;
                protocol: string | null;
                inboundTag: string | null;
                level: number | null;
                source: string | null;
                destination: string;
                routeTarget: string | null;
                originalTarget: string | null;
                inboundName: string | null;
                inboundLocal: string | null;
                outboundTag: string | null;
                ts: number;
            };
        }, {
            actionReport: {
                ip: string;
                userId: string;
                blocked: boolean;
                blockDuration: number;
                willUnblockAt: string;
                processedAt: string;
            };
            xrayReport: {
                network: string;
                email: string | null;
                protocol: string | null;
                inboundTag: string | null;
                level: number | null;
                source: string | null;
                destination: string;
                routeTarget: string | null;
                originalTarget: string | null;
                inboundName: string | null;
                inboundLocal: string | null;
                outboundTag: string | null;
                ts: number;
            };
        }>;
    }, "strip", z.ZodTypeAny, {
        user: {
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
        };
        node: {
            tags: string[];
            system: {
                stats: {
                    interface: {
                        interface: string;
                        rxBytesPerSec: number;
                        txBytesPerSec: number;
                        rxTotal: number;
                        txTotal: number;
                    } | null;
                    memoryFree: number;
                    memoryUsed: number;
                    uptime: number;
                    loadAvg: number[];
                };
                info: {
                    type: string;
                    version: string;
                    platform: string;
                    arch: string;
                    cpus: number;
                    cpuModel: string;
                    memoryTotal: number;
                    hostname: string;
                    release: string;
                    networkInterfaces: string[];
                };
            } | null;
            uuid: string;
            name: string;
            createdAt: Date;
            updatedAt: Date;
            provider: {
                uuid: string;
                name: string;
                createdAt: Date;
                updatedAt: Date;
                faviconLink: string | null;
                loginUrl: string | null;
            } | null;
            countryCode: string;
            port: number | null;
            viewPosition: number;
            trafficLimitBytes: number | null;
            address: string;
            isDisabled: boolean;
            proxyUrl: string | null;
            isConnected: boolean;
            isConnecting: boolean;
            lastStatusChange: Date | null;
            lastStatusMessage: string | null;

            singboxVersion: string | null;

            nodeVersion: string | null;
            isTrafficTrackingActive: boolean;
            trafficResetDay: number | null;
            trafficUsedBytes: number | null;
            notifyPercent: number | null;
            consumptionMultiplier: number;
            nodeConsumptionMultiplier: number;
            configProfile: {
                activeConfigProfileUuid: string | null;
                activeInbounds: {
                    uuid: string;
                    type: string;
                    profileUuid: string;
                    tag: string;
                    network: string | null;
                    security: string | null;
                    port: number | null;
                    rawInbound?: unknown;
                }[];
            };
            providerUuid: string | null;
            activePluginUuid: string | null;
            versions: {
                node: string;
                singbox: string;
            } | null;
            singboxUptime: number;
            usersOnline: number;
            note: string | null;
        };
        report: {
            actionReport: {
                ip: string;
                userId: string;
                blocked: boolean;
                blockDuration: number;
                willUnblockAt: Date;
                processedAt: Date;
            };
            xrayReport: {
                network: string;
                email: string | null;
                protocol: string | null;
                inboundTag: string | null;
                level: number | null;
                source: string | null;
                destination: string;
                routeTarget: string | null;
                originalTarget: string | null;
                inboundName: string | null;
                inboundLocal: string | null;
                outboundTag: string | null;
                ts: number;
            };
        };
    }, {
        user: {
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
        };
        node: {
            tags: string[];
            system: {
                stats: {
                    interface: {
                        interface: string;
                        rxBytesPerSec: number;
                        txBytesPerSec: number;
                        rxTotal: number;
                        txTotal: number;
                    } | null;
                    memoryFree: number;
                    memoryUsed: number;
                    uptime: number;
                    loadAvg: number[];
                };
                info: {
                    type: string;
                    version: string;
                    platform: string;
                    arch: string;
                    cpus: number;
                    cpuModel: string;
                    memoryTotal: number;
                    hostname: string;
                    release: string;
                    networkInterfaces: string[];
                };
            } | null;
            uuid: string;
            name: string;
            createdAt: string;
            updatedAt: string;
            provider: {
                uuid: string;
                name: string;
                createdAt: string;
                updatedAt: string;
                faviconLink: string | null;
                loginUrl: string | null;
            } | null;
            countryCode: string;
            port: number | null;
            viewPosition: number;
            trafficLimitBytes: number | null;
            address: string;
            isDisabled: boolean;
            proxyUrl: string | null;
            isConnected: boolean;
            isConnecting: boolean;
            lastStatusChange: string | null;
            lastStatusMessage: string | null;

            singboxVersion: string | null;

            nodeVersion: string | null;
            isTrafficTrackingActive: boolean;
            trafficResetDay: number | null;
            trafficUsedBytes: number | null;
            notifyPercent: number | null;
            consumptionMultiplier: number;
            nodeConsumptionMultiplier: number;
            configProfile: {
                activeConfigProfileUuid: string | null;
                activeInbounds: {
                    uuid: string;
                    type: string;
                    profileUuid: string;
                    tag: string;
                    network: string | null;
                    security: string | null;
                    port: number | null;
                    rawInbound?: unknown;
                }[];
            };
            providerUuid: string | null;
            activePluginUuid: string | null;
            versions: {
                node: string;
                singbox: string;
            } | null;
            singboxUptime: number;
            usersOnline: number;
            note: string | null;
        };
        report: {
            actionReport: {
                ip: string;
                userId: string;
                blocked: boolean;
                blockDuration: number;
                willUnblockAt: string;
                processedAt: string;
            };
            xrayReport: {
                network: string;
                email: string | null;
                protocol: string | null;
                inboundTag: string | null;
                level: number | null;
                source: string | null;
                destination: string;
                routeTarget: string | null;
                originalTarget: string | null;
                inboundName: string | null;
                inboundLocal: string | null;
                outboundTag: string | null;
                ts: number;
            };
        };
    }>;
}, "strip", z.ZodTypeAny, {
    data: {
        user: {
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
        };
        node: {
            tags: string[];
            system: {
                stats: {
                    interface: {
                        interface: string;
                        rxBytesPerSec: number;
                        txBytesPerSec: number;
                        rxTotal: number;
                        txTotal: number;
                    } | null;
                    memoryFree: number;
                    memoryUsed: number;
                    uptime: number;
                    loadAvg: number[];
                };
                info: {
                    type: string;
                    version: string;
                    platform: string;
                    arch: string;
                    cpus: number;
                    cpuModel: string;
                    memoryTotal: number;
                    hostname: string;
                    release: string;
                    networkInterfaces: string[];
                };
            } | null;
            uuid: string;
            name: string;
            createdAt: Date;
            updatedAt: Date;
            provider: {
                uuid: string;
                name: string;
                createdAt: Date;
                updatedAt: Date;
                faviconLink: string | null;
                loginUrl: string | null;
            } | null;
            countryCode: string;
            port: number | null;
            viewPosition: number;
            trafficLimitBytes: number | null;
            address: string;
            isDisabled: boolean;
            proxyUrl: string | null;
            isConnected: boolean;
            isConnecting: boolean;
            lastStatusChange: Date | null;
            lastStatusMessage: string | null;

            singboxVersion: string | null;

            nodeVersion: string | null;
            isTrafficTrackingActive: boolean;
            trafficResetDay: number | null;
            trafficUsedBytes: number | null;
            notifyPercent: number | null;
            consumptionMultiplier: number;
            nodeConsumptionMultiplier: number;
            configProfile: {
                activeConfigProfileUuid: string | null;
                activeInbounds: {
                    uuid: string;
                    type: string;
                    profileUuid: string;
                    tag: string;
                    network: string | null;
                    security: string | null;
                    port: number | null;
                    rawInbound?: unknown;
                }[];
            };
            providerUuid: string | null;
            activePluginUuid: string | null;
            versions: {
                node: string;
                singbox: string;
            } | null;
            singboxUptime: number;
            usersOnline: number;
            note: string | null;
        };
        report: {
            actionReport: {
                ip: string;
                userId: string;
                blocked: boolean;
                blockDuration: number;
                willUnblockAt: Date;
                processedAt: Date;
            };
            xrayReport: {
                network: string;
                email: string | null;
                protocol: string | null;
                inboundTag: string | null;
                level: number | null;
                source: string | null;
                destination: string;
                routeTarget: string | null;
                originalTarget: string | null;
                inboundName: string | null;
                inboundLocal: string | null;
                outboundTag: string | null;
                ts: number;
            };
        };
    };
    scope: "torrent_blocker";
    event: "torrent_blocker.report";
    timestamp: Date;
}, {
    data: {
        user: {
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
        };
        node: {
            tags: string[];
            system: {
                stats: {
                    interface: {
                        interface: string;
                        rxBytesPerSec: number;
                        txBytesPerSec: number;
                        rxTotal: number;
                        txTotal: number;
                    } | null;
                    memoryFree: number;
                    memoryUsed: number;
                    uptime: number;
                    loadAvg: number[];
                };
                info: {
                    type: string;
                    version: string;
                    platform: string;
                    arch: string;
                    cpus: number;
                    cpuModel: string;
                    memoryTotal: number;
                    hostname: string;
                    release: string;
                    networkInterfaces: string[];
                };
            } | null;
            uuid: string;
            name: string;
            createdAt: string;
            updatedAt: string;
            provider: {
                uuid: string;
                name: string;
                createdAt: string;
                updatedAt: string;
                faviconLink: string | null;
                loginUrl: string | null;
            } | null;
            countryCode: string;
            port: number | null;
            viewPosition: number;
            trafficLimitBytes: number | null;
            address: string;
            isDisabled: boolean;
            proxyUrl: string | null;
            isConnected: boolean;
            isConnecting: boolean;
            lastStatusChange: string | null;
            lastStatusMessage: string | null;

            singboxVersion: string | null;

            nodeVersion: string | null;
            isTrafficTrackingActive: boolean;
            trafficResetDay: number | null;
            trafficUsedBytes: number | null;
            notifyPercent: number | null;
            consumptionMultiplier: number;
            nodeConsumptionMultiplier: number;
            configProfile: {
                activeConfigProfileUuid: string | null;
                activeInbounds: {
                    uuid: string;
                    type: string;
                    profileUuid: string;
                    tag: string;
                    network: string | null;
                    security: string | null;
                    port: number | null;
                    rawInbound?: unknown;
                }[];
            };
            providerUuid: string | null;
            activePluginUuid: string | null;
            versions: {
                node: string;
                singbox: string;
            } | null;
            singboxUptime: number;
            usersOnline: number;
            note: string | null;
        };
        report: {
            actionReport: {
                ip: string;
                userId: string;
                blocked: boolean;
                blockDuration: number;
                willUnblockAt: string;
                processedAt: string;
            };
            xrayReport: {
                network: string;
                email: string | null;
                protocol: string | null;
                inboundTag: string | null;
                level: number | null;
                source: string | null;
                destination: string;
                routeTarget: string | null;
                originalTarget: string | null;
                inboundName: string | null;
                inboundLocal: string | null;
                outboundTag: string | null;
                ts: number;
            };
        };
    };
    scope: "torrent_blocker";
    event: "torrent_blocker.report";
    timestamp: string;
}>;
export declare const ExodusWebhookEventSchema: z.ZodDiscriminatedUnion<"scope", [z.ZodObject<{
    scope: z.ZodLiteral<"user">;
    event: z.ZodEnum<["user.created" | "user.modified" | "user.deleted" | "user.revoked" | "user.disabled" | "user.enabled" | "user.limited" | "user.expired" | "user.traffic_reset" | "user.first_connected" | "user.bandwidth_usage_threshold_reached" | "user.not_connected" | "user.expiration", ...("user.created" | "user.modified" | "user.deleted" | "user.revoked" | "user.disabled" | "user.enabled" | "user.limited" | "user.expired" | "user.traffic_reset" | "user.first_connected" | "user.bandwidth_usage_threshold_reached" | "user.not_connected" | "user.expiration")[]]>;
    timestamp: z.ZodEffects<z.ZodString, Date, string>;
    data: z.ZodObject<{
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
    }>;
    meta: z.ZodNullable<z.ZodObject<{
        notConnectedAfterHours: z.ZodOptional<z.ZodNullable<z.ZodNumber>>;
        expiration: z.ZodOptional<z.ZodNullable<z.ZodNumber>>;
    }, "strip", z.ZodTypeAny, {
        notConnectedAfterHours?: number | null | undefined;
        expiration?: number | null | undefined;
    }, {
        notConnectedAfterHours?: number | null | undefined;
        expiration?: number | null | undefined;
    }>>;
}, "strip", z.ZodTypeAny, {
    data: {
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
    };
    scope: "user";
    event: "user.created" | "user.modified" | "user.deleted" | "user.revoked" | "user.disabled" | "user.enabled" | "user.limited" | "user.expired" | "user.traffic_reset" | "user.first_connected" | "user.bandwidth_usage_threshold_reached" | "user.not_connected" | "user.expiration";
    timestamp: Date;
    meta: {
        notConnectedAfterHours?: number | null | undefined;
        expiration?: number | null | undefined;
    } | null;
}, {
    data: {
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
    };
    scope: "user";
    event: "user.created" | "user.modified" | "user.deleted" | "user.revoked" | "user.disabled" | "user.enabled" | "user.limited" | "user.expired" | "user.traffic_reset" | "user.first_connected" | "user.bandwidth_usage_threshold_reached" | "user.not_connected" | "user.expiration";
    timestamp: string;
    meta: {
        notConnectedAfterHours?: number | null | undefined;
        expiration?: number | null | undefined;
    } | null;
}>, z.ZodObject<{
    scope: z.ZodLiteral<"user_hwid_devices">;
    event: z.ZodEnum<["user_hwid_devices.added" | "user_hwid_devices.deleted", ...("user_hwid_devices.added" | "user_hwid_devices.deleted")[]]>;
    timestamp: z.ZodEffects<z.ZodString, Date, string>;
    data: z.ZodObject<{
        user: z.ZodObject<{
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
        }>;
        hwidUserDevice: z.ZodObject<{
            hwid: z.ZodString;
            userId: z.ZodNumber;
            platform: z.ZodNullable<z.ZodString>;
            osVersion: z.ZodNullable<z.ZodString>;
            deviceModel: z.ZodNullable<z.ZodString>;
            userAgent: z.ZodNullable<z.ZodString>;
            requestIp: z.ZodNullable<z.ZodString>;
            createdAt: z.ZodEffects<z.ZodString, Date, string>;
            updatedAt: z.ZodEffects<z.ZodString, Date, string>;
        }, "strip", z.ZodTypeAny, {
            hwid: string;
            createdAt: Date;
            updatedAt: Date;
            userId: number;
            platform: string | null;
            osVersion: string | null;
            deviceModel: string | null;
            userAgent: string | null;
            requestIp: string | null;
        }, {
            hwid: string;
            createdAt: string;
            updatedAt: string;
            userId: number;
            platform: string | null;
            osVersion: string | null;
            deviceModel: string | null;
            userAgent: string | null;
            requestIp: string | null;
        }>;
    }, "strip", z.ZodTypeAny, {
        user: {
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
        };
        hwidUserDevice: {
            hwid: string;
            createdAt: Date;
            updatedAt: Date;
            userId: number;
            platform: string | null;
            osVersion: string | null;
            deviceModel: string | null;
            userAgent: string | null;
            requestIp: string | null;
        };
    }, {
        user: {
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
        };
        hwidUserDevice: {
            hwid: string;
            createdAt: string;
            updatedAt: string;
            userId: number;
            platform: string | null;
            osVersion: string | null;
            deviceModel: string | null;
            userAgent: string | null;
            requestIp: string | null;
        };
    }>;
}, "strip", z.ZodTypeAny, {
    data: {
        user: {
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
        };
        hwidUserDevice: {
            hwid: string;
            createdAt: Date;
            updatedAt: Date;
            userId: number;
            platform: string | null;
            osVersion: string | null;
            deviceModel: string | null;
            userAgent: string | null;
            requestIp: string | null;
        };
    };
    scope: "user_hwid_devices";
    event: "user_hwid_devices.added" | "user_hwid_devices.deleted";
    timestamp: Date;
}, {
    data: {
        user: {
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
        };
        hwidUserDevice: {
            hwid: string;
            createdAt: string;
            updatedAt: string;
            userId: number;
            platform: string | null;
            osVersion: string | null;
            deviceModel: string | null;
            userAgent: string | null;
            requestIp: string | null;
        };
    };
    scope: "user_hwid_devices";
    event: "user_hwid_devices.added" | "user_hwid_devices.deleted";
    timestamp: string;
}>, z.ZodObject<{
    scope: z.ZodLiteral<"node">;
    event: z.ZodEnum<["node.created" | "node.modified" | "node.disabled" | "node.enabled" | "node.deleted" | "node.connection_lost" | "node.connection_restored" | "node.traffic_notify", ...("node.created" | "node.modified" | "node.disabled" | "node.enabled" | "node.deleted" | "node.connection_lost" | "node.connection_restored" | "node.traffic_notify")[]]>;
    timestamp: z.ZodEffects<z.ZodString, Date, string>;
    data: z.ZodObject<{
        uuid: z.ZodString;
        name: z.ZodString;
        address: z.ZodString;
        port: z.ZodNullable<z.ZodNumber>;
        proxyUrl: z.ZodNullable<z.ZodString>;
        isConnected: z.ZodBoolean;
        isDisabled: z.ZodBoolean;
        isConnecting: z.ZodBoolean;
        lastStatusChange: z.ZodNullable<z.ZodEffects<z.ZodString, Date, string>>;
            lastStatusMessage: z.ZodNullable<z.ZodString>;

            singboxVersion: z.ZodNullable<z.ZodString>;

            nodeVersion: z.ZodNullable<z.ZodString>;
        isTrafficTrackingActive: z.ZodBoolean;
        trafficResetDay: z.ZodNullable<z.ZodNumber>;
        trafficLimitBytes: z.ZodNullable<z.ZodNumber>;
        trafficUsedBytes: z.ZodNullable<z.ZodNumber>;
        notifyPercent: z.ZodNullable<z.ZodNumber>;
        viewPosition: z.ZodNumber;
        countryCode: z.ZodString;
        consumptionMultiplier: z.ZodNumber;
        nodeConsumptionMultiplier: z.ZodNumber;
        tags: z.ZodArray<z.ZodString, "many">;
        createdAt: z.ZodEffects<z.ZodString, Date, string>;
        updatedAt: z.ZodEffects<z.ZodString, Date, string>;
        configProfile: z.ZodObject<{
            activeConfigProfileUuid: z.ZodNullable<z.ZodString>;
            activeInbounds: z.ZodArray<z.ZodObject<{
                uuid: z.ZodString;
                profileUuid: z.ZodString;
                tag: z.ZodString;
                type: z.ZodString;
                network: z.ZodNullable<z.ZodString>;
                security: z.ZodNullable<z.ZodString>;
                port: z.ZodNullable<z.ZodNumber>;
                rawInbound: z.ZodNullable<z.ZodUnknown>;
            }, "strip", z.ZodTypeAny, {
                uuid: string;
                type: string;
                profileUuid: string;
                tag: string;
                network: string | null;
                security: string | null;
                port: number | null;
                rawInbound?: unknown;
            }, {
                uuid: string;
                type: string;
                profileUuid: string;
                tag: string;
                network: string | null;
                security: string | null;
                port: number | null;
                rawInbound?: unknown;
            }>, "many">;
        }, "strip", z.ZodTypeAny, {
            activeConfigProfileUuid: string | null;
            activeInbounds: {
                uuid: string;
                type: string;
                profileUuid: string;
                tag: string;
                network: string | null;
                security: string | null;
                port: number | null;
                rawInbound?: unknown;
            }[];
        }, {
            activeConfigProfileUuid: string | null;
            activeInbounds: {
                uuid: string;
                type: string;
                profileUuid: string;
                tag: string;
                network: string | null;
                security: string | null;
                port: number | null;
                rawInbound?: unknown;
            }[];
        }>;
        providerUuid: z.ZodNullable<z.ZodString>;
        provider: z.ZodNullable<z.ZodObject<{
            uuid: z.ZodString;
            name: z.ZodString;
            faviconLink: z.ZodNullable<z.ZodString>;
            loginUrl: z.ZodNullable<z.ZodString>;
            createdAt: z.ZodEffects<z.ZodString, Date, string>;
            updatedAt: z.ZodEffects<z.ZodString, Date, string>;
        }, "strip", z.ZodTypeAny, {
            uuid: string;
            name: string;
            createdAt: Date;
            updatedAt: Date;
            faviconLink: string | null;
            loginUrl: string | null;
        }, {
            uuid: string;
            name: string;
            createdAt: string;
            updatedAt: string;
            faviconLink: string | null;
            loginUrl: string | null;
        }>>;
        activePluginUuid: z.ZodNullable<z.ZodString>;
        system: z.ZodNullable<z.ZodObject<{
            info: z.ZodObject<{
                arch: z.ZodString;
                cpus: z.ZodNumber;
                cpuModel: z.ZodString;
                memoryTotal: z.ZodNumber;
                hostname: z.ZodString;
                platform: z.ZodString;
                release: z.ZodString;
                type: z.ZodString;
                version: z.ZodString;
                networkInterfaces: z.ZodArray<z.ZodString, "many">;
            }, "strip", z.ZodTypeAny, {
                type: string;
                version: string;
                platform: string;
                arch: string;
                cpus: number;
                cpuModel: string;
                memoryTotal: number;
                hostname: string;
                release: string;
                networkInterfaces: string[];
            }, {
                type: string;
                version: string;
                platform: string;
                arch: string;
                cpus: number;
                cpuModel: string;
                memoryTotal: number;
                hostname: string;
                release: string;
                networkInterfaces: string[];
            }>;
            stats: z.ZodObject<{
                memoryFree: z.ZodNumber;
                memoryUsed: z.ZodNumber;
                uptime: z.ZodNumber;
                loadAvg: z.ZodArray<z.ZodNumber, "many">;
                interface: z.ZodNullable<z.ZodObject<{
                    interface: z.ZodString;
                    rxBytesPerSec: z.ZodNumber;
                    txBytesPerSec: z.ZodNumber;
                    rxTotal: z.ZodNumber;
                    txTotal: z.ZodNumber;
                }, "strip", z.ZodTypeAny, {
                    interface: string;
                    rxBytesPerSec: number;
                    txBytesPerSec: number;
                    rxTotal: number;
                    txTotal: number;
                }, {
                    interface: string;
                    rxBytesPerSec: number;
                    txBytesPerSec: number;
                    rxTotal: number;
                    txTotal: number;
                }>>;
            }, "strip", z.ZodTypeAny, {
                interface: {
                    interface: string;
                    rxBytesPerSec: number;
                    txBytesPerSec: number;
                    rxTotal: number;
                    txTotal: number;
                } | null;
                memoryFree: number;
                memoryUsed: number;
                uptime: number;
                loadAvg: number[];
            }, {
                interface: {
                    interface: string;
                    rxBytesPerSec: number;
                    txBytesPerSec: number;
                    rxTotal: number;
                    txTotal: number;
                } | null;
                memoryFree: number;
                memoryUsed: number;
                uptime: number;
                loadAvg: number[];
            }>;
        }, "strip", z.ZodTypeAny, {
            stats: {
                interface: {
                    interface: string;
                    rxBytesPerSec: number;
                    txBytesPerSec: number;
                    rxTotal: number;
                    txTotal: number;
                } | null;
                memoryFree: number;
                memoryUsed: number;
                uptime: number;
                loadAvg: number[];
            };
            info: {
                type: string;
                version: string;
                platform: string;
                arch: string;
                cpus: number;
                cpuModel: string;
                memoryTotal: number;
                hostname: string;
                release: string;
                networkInterfaces: string[];
            };
        }, {
            stats: {
                interface: {
                    interface: string;
                    rxBytesPerSec: number;
                    txBytesPerSec: number;
                    rxTotal: number;
                    txTotal: number;
                } | null;
                memoryFree: number;
                memoryUsed: number;
                uptime: number;
                loadAvg: number[];
            };
            info: {
                type: string;
                version: string;
                platform: string;
                arch: string;
                cpus: number;
                cpuModel: string;
                memoryTotal: number;
                hostname: string;
                release: string;
                networkInterfaces: string[];
            };
        }>>;
        versions: z.ZodNullable<z.ZodObject<{
            singbox: z.ZodString;
            node: z.ZodString;
        }, "strip", z.ZodTypeAny, {
            node: string;
            singbox: string;
        }, {
            node: string;
            singbox: string;
        }>>;
        singboxUptime: z.ZodNumber;
        usersOnline: z.ZodNumber;
        note: z.ZodNullable<z.ZodString>;
    }, "strip", z.ZodTypeAny, {
        tags: string[];
        system: {
            stats: {
                interface: {
                    interface: string;
                    rxBytesPerSec: number;
                    txBytesPerSec: number;
                    rxTotal: number;
                    txTotal: number;
                } | null;
                memoryFree: number;
                memoryUsed: number;
                uptime: number;
                loadAvg: number[];
            };
            info: {
                type: string;
                version: string;
                platform: string;
                arch: string;
                cpus: number;
                cpuModel: string;
                memoryTotal: number;
                hostname: string;
                release: string;
                networkInterfaces: string[];
            };
        } | null;
        uuid: string;
        name: string;
        createdAt: Date;
        updatedAt: Date;
        provider: {
            uuid: string;
            name: string;
            createdAt: Date;
            updatedAt: Date;
            faviconLink: string | null;
            loginUrl: string | null;
        } | null;
        countryCode: string;
        port: number | null;
        viewPosition: number;
        trafficLimitBytes: number | null;
        address: string;
        isDisabled: boolean;
        proxyUrl: string | null;
        isConnected: boolean;
        isConnecting: boolean;
        lastStatusChange: Date | null;
            lastStatusMessage: string | null;

            singboxVersion: string | null;

            nodeVersion: string | null;
        isTrafficTrackingActive: boolean;
        trafficResetDay: number | null;
        trafficUsedBytes: number | null;
        notifyPercent: number | null;
        consumptionMultiplier: number;
        nodeConsumptionMultiplier: number;
        configProfile: {
            activeConfigProfileUuid: string | null;
            activeInbounds: {
                uuid: string;
                type: string;
                profileUuid: string;
                tag: string;
                network: string | null;
                security: string | null;
                port: number | null;
                rawInbound?: unknown;
            }[];
        };
        providerUuid: string | null;
        activePluginUuid: string | null;
        versions: {
            node: string;
            singbox: string;
        } | null;
        singboxUptime: number;
        usersOnline: number;
        note: string | null;
    }, {
        tags: string[];
        system: {
            stats: {
                interface: {
                    interface: string;
                    rxBytesPerSec: number;
                    txBytesPerSec: number;
                    rxTotal: number;
                    txTotal: number;
                } | null;
                memoryFree: number;
                memoryUsed: number;
                uptime: number;
                loadAvg: number[];
            };
            info: {
                type: string;
                version: string;
                platform: string;
                arch: string;
                cpus: number;
                cpuModel: string;
                memoryTotal: number;
                hostname: string;
                release: string;
                networkInterfaces: string[];
            };
        } | null;
        uuid: string;
        name: string;
        createdAt: string;
        updatedAt: string;
        provider: {
            uuid: string;
            name: string;
            createdAt: string;
            updatedAt: string;
            faviconLink: string | null;
            loginUrl: string | null;
        } | null;
        countryCode: string;
        port: number | null;
        viewPosition: number;
        trafficLimitBytes: number | null;
        address: string;
        isDisabled: boolean;
        proxyUrl: string | null;
        isConnected: boolean;
        isConnecting: boolean;
        lastStatusChange: string | null;
            lastStatusMessage: string | null;

            singboxVersion: string | null;

            nodeVersion: string | null;
        isTrafficTrackingActive: boolean;
        trafficResetDay: number | null;
        trafficUsedBytes: number | null;
        notifyPercent: number | null;
        consumptionMultiplier: number;
        nodeConsumptionMultiplier: number;
        configProfile: {
            activeConfigProfileUuid: string | null;
            activeInbounds: {
                uuid: string;
                type: string;
                profileUuid: string;
                tag: string;
                network: string | null;
                security: string | null;
                port: number | null;
                rawInbound?: unknown;
            }[];
        };
        providerUuid: string | null;
        activePluginUuid: string | null;
        versions: {
            node: string;
            singbox: string;
        } | null;
        singboxUptime: number;
        usersOnline: number;
        note: string | null;
    }>;
}, "strip", z.ZodTypeAny, {
    data: {
        tags: string[];
        system: {
            stats: {
                interface: {
                    interface: string;
                    rxBytesPerSec: number;
                    txBytesPerSec: number;
                    rxTotal: number;
                    txTotal: number;
                } | null;
                memoryFree: number;
                memoryUsed: number;
                uptime: number;
                loadAvg: number[];
            };
            info: {
                type: string;
                version: string;
                platform: string;
                arch: string;
                cpus: number;
                cpuModel: string;
                memoryTotal: number;
                hostname: string;
                release: string;
                networkInterfaces: string[];
            };
        } | null;
        uuid: string;
        name: string;
        createdAt: Date;
        updatedAt: Date;
        provider: {
            uuid: string;
            name: string;
            createdAt: Date;
            updatedAt: Date;
            faviconLink: string | null;
            loginUrl: string | null;
        } | null;
        countryCode: string;
        port: number | null;
        viewPosition: number;
        trafficLimitBytes: number | null;
        address: string;
        isDisabled: boolean;
        proxyUrl: string | null;
        isConnected: boolean;
        isConnecting: boolean;
        lastStatusChange: Date | null;
            lastStatusMessage: string | null;

            singboxVersion: string | null;

            nodeVersion: string | null;
        isTrafficTrackingActive: boolean;
        trafficResetDay: number | null;
        trafficUsedBytes: number | null;
        notifyPercent: number | null;
        consumptionMultiplier: number;
        nodeConsumptionMultiplier: number;
        configProfile: {
            activeConfigProfileUuid: string | null;
            activeInbounds: {
                uuid: string;
                type: string;
                profileUuid: string;
                tag: string;
                network: string | null;
                security: string | null;
                port: number | null;
                rawInbound?: unknown;
            }[];
        };
        providerUuid: string | null;
        activePluginUuid: string | null;
        versions: {
            node: string;
            singbox: string;
        } | null;
        singboxUptime: number;
        usersOnline: number;
        note: string | null;
    };
    scope: "node";
    event: "node.created" | "node.modified" | "node.disabled" | "node.enabled" | "node.deleted" | "node.connection_lost" | "node.connection_restored" | "node.traffic_notify";
    timestamp: Date;
}, {
    data: {
        tags: string[];
        system: {
            stats: {
                interface: {
                    interface: string;
                    rxBytesPerSec: number;
                    txBytesPerSec: number;
                    rxTotal: number;
                    txTotal: number;
                } | null;
                memoryFree: number;
                memoryUsed: number;
                uptime: number;
                loadAvg: number[];
            };
            info: {
                type: string;
                version: string;
                platform: string;
                arch: string;
                cpus: number;
                cpuModel: string;
                memoryTotal: number;
                hostname: string;
                release: string;
                networkInterfaces: string[];
            };
        } | null;
        uuid: string;
        name: string;
        createdAt: string;
        updatedAt: string;
        provider: {
            uuid: string;
            name: string;
            createdAt: string;
            updatedAt: string;
            faviconLink: string | null;
            loginUrl: string | null;
        } | null;
        countryCode: string;
        port: number | null;
        viewPosition: number;
        trafficLimitBytes: number | null;
        address: string;
        isDisabled: boolean;
        proxyUrl: string | null;
        isConnected: boolean;
        isConnecting: boolean;
        lastStatusChange: string | null;
            lastStatusMessage: string | null;

            singboxVersion: string | null;

            nodeVersion: string | null;
        isTrafficTrackingActive: boolean;
        trafficResetDay: number | null;
        trafficUsedBytes: number | null;
        notifyPercent: number | null;
        consumptionMultiplier: number;
        nodeConsumptionMultiplier: number;
        configProfile: {
            activeConfigProfileUuid: string | null;
            activeInbounds: {
                uuid: string;
                type: string;
                profileUuid: string;
                tag: string;
                network: string | null;
                security: string | null;
                port: number | null;
                rawInbound?: unknown;
            }[];
        };
        providerUuid: string | null;
        activePluginUuid: string | null;
        versions: {
            node: string;
            singbox: string;
        } | null;
        singboxUptime: number;
        usersOnline: number;
        note: string | null;
    };
    scope: "node";
    event: "node.created" | "node.modified" | "node.disabled" | "node.enabled" | "node.deleted" | "node.connection_lost" | "node.connection_restored" | "node.traffic_notify";
    timestamp: string;
}>, z.ZodObject<{
    scope: z.ZodLiteral<"service">;
    event: z.ZodEnum<["service.panel_started" | "service.login_attempt_failed" | "service.login_attempt_success" | "service.subpage_config_changed" | "service.api_token_created" | "service.api_token_deleted", ...("service.panel_started" | "service.login_attempt_failed" | "service.login_attempt_success" | "service.subpage_config_changed" | "service.api_token_created" | "service.api_token_deleted")[]]>;
    timestamp: z.ZodEffects<z.ZodString, Date, string>;
    data: z.ZodObject<{
        loginAttempt: z.ZodOptional<z.ZodObject<{
            username: z.ZodString;
            ip: z.ZodString;
            userAgent: z.ZodString;
            description: z.ZodOptional<z.ZodString>;
            password: z.ZodOptional<z.ZodString>;
        }, "strip", z.ZodTypeAny, {
            username: string;
            ip: string;
            userAgent: string;
            description?: string | undefined;
            password?: string | undefined;
        }, {
            username: string;
            ip: string;
            userAgent: string;
            description?: string | undefined;
            password?: string | undefined;
        }>>;
        panelVersion: z.ZodOptional<z.ZodString>;
        subpageConfig: z.ZodOptional<z.ZodObject<{
            action: z.ZodEnum<["CREATED" | "UPDATED" | "DELETED", ...("CREATED" | "UPDATED" | "DELETED")[]]>;
            uuid: z.ZodString;
        }, "strip", z.ZodTypeAny, {
            uuid: string;
            action: "CREATED" | "UPDATED" | "DELETED";
        }, {
            uuid: string;
            action: "CREATED" | "UPDATED" | "DELETED";
        }>>;
        apiToken: z.ZodOptional<z.ZodObject<{
            name: z.ZodString;
            uuid: z.ZodString;
            expireAt: z.ZodEffects<z.ZodString, Date, string>;
            scopes: z.ZodArray<z.ZodString, "many">;
        }, "strip", z.ZodTypeAny, {
            scopes: string[];
            uuid: string;
            name: string;
            expireAt: Date;
        }, {
            scopes: string[];
            uuid: string;
            name: string;
            expireAt: string;
        }>>;
    }, "strip", z.ZodTypeAny, {
        loginAttempt?: {
            username: string;
            ip: string;
            userAgent: string;
            description?: string | undefined;
            password?: string | undefined;
        } | undefined;
        panelVersion?: string | undefined;
        subpageConfig?: {
            uuid: string;
            action: "CREATED" | "UPDATED" | "DELETED";
        } | undefined;
        apiToken?: {
            scopes: string[];
            uuid: string;
            name: string;
            expireAt: Date;
        } | undefined;
    }, {
        loginAttempt?: {
            username: string;
            ip: string;
            userAgent: string;
            description?: string | undefined;
            password?: string | undefined;
        } | undefined;
        panelVersion?: string | undefined;
        subpageConfig?: {
            uuid: string;
            action: "CREATED" | "UPDATED" | "DELETED";
        } | undefined;
        apiToken?: {
            scopes: string[];
            uuid: string;
            name: string;
            expireAt: string;
        } | undefined;
    }>;
}, "strip", z.ZodTypeAny, {
    data: {
        loginAttempt?: {
            username: string;
            ip: string;
            userAgent: string;
            description?: string | undefined;
            password?: string | undefined;
        } | undefined;
        panelVersion?: string | undefined;
        subpageConfig?: {
            uuid: string;
            action: "CREATED" | "UPDATED" | "DELETED";
        } | undefined;
        apiToken?: {
            scopes: string[];
            uuid: string;
            name: string;
            expireAt: Date;
        } | undefined;
    };
    scope: "service";
    event: "service.panel_started" | "service.login_attempt_failed" | "service.login_attempt_success" | "service.subpage_config_changed" | "service.api_token_created" | "service.api_token_deleted";
    timestamp: Date;
}, {
    data: {
        loginAttempt?: {
            username: string;
            ip: string;
            userAgent: string;
            description?: string | undefined;
            password?: string | undefined;
        } | undefined;
        panelVersion?: string | undefined;
        subpageConfig?: {
            uuid: string;
            action: "CREATED" | "UPDATED" | "DELETED";
        } | undefined;
        apiToken?: {
            scopes: string[];
            uuid: string;
            name: string;
            expireAt: string;
        } | undefined;
    };
    scope: "service";
    event: "service.panel_started" | "service.login_attempt_failed" | "service.login_attempt_success" | "service.subpage_config_changed" | "service.api_token_created" | "service.api_token_deleted";
    timestamp: string;
}>, z.ZodObject<{
    scope: z.ZodLiteral<"errors">;
    event: z.ZodEnum<["errors.bandwidth_usage_threshold_reached_max_notifications", ..."errors.bandwidth_usage_threshold_reached_max_notifications"[]]>;
    timestamp: z.ZodEffects<z.ZodString, Date, string>;
    data: z.ZodObject<{
        description: z.ZodString;
    }, "strip", z.ZodTypeAny, {
        description: string;
    }, {
        description: string;
    }>;
}, "strip", z.ZodTypeAny, {
    data: {
        description: string;
    };
    scope: "errors";
    event: "errors.bandwidth_usage_threshold_reached_max_notifications";
    timestamp: Date;
}, {
    data: {
        description: string;
    };
    scope: "errors";
    event: "errors.bandwidth_usage_threshold_reached_max_notifications";
    timestamp: string;
}>, z.ZodObject<{
    scope: z.ZodLiteral<"crm">;
    event: z.ZodEnum<["crm.infra_billing_node_payment_in_7_days" | "crm.infra_billing_node_payment_in_48hrs" | "crm.infra_billing_node_payment_in_24hrs" | "crm.infra_billing_node_payment_due_today" | "crm.infra_billing_node_payment_overdue_24hrs" | "crm.infra_billing_node_payment_overdue_48hrs" | "crm.infra_billing_node_payment_overdue_7_days", ...("crm.infra_billing_node_payment_in_7_days" | "crm.infra_billing_node_payment_in_48hrs" | "crm.infra_billing_node_payment_in_24hrs" | "crm.infra_billing_node_payment_due_today" | "crm.infra_billing_node_payment_overdue_24hrs" | "crm.infra_billing_node_payment_overdue_48hrs" | "crm.infra_billing_node_payment_overdue_7_days")[]]>;
    timestamp: z.ZodEffects<z.ZodString, Date, string>;
    data: z.ZodObject<{
        providerName: z.ZodString;
        nodeName: z.ZodString;
        nextBillingAt: z.ZodEffects<z.ZodString, Date, string>;
        loginUrl: z.ZodString;
    }, "strip", z.ZodTypeAny, {
        nodeName: string;
        loginUrl: string;
        nextBillingAt: Date;
        providerName: string;
    }, {
        nodeName: string;
        loginUrl: string;
        nextBillingAt: string;
        providerName: string;
    }>;
}, "strip", z.ZodTypeAny, {
    data: {
        nodeName: string;
        loginUrl: string;
        nextBillingAt: Date;
        providerName: string;
    };
    scope: "crm";
    event: "crm.infra_billing_node_payment_in_7_days" | "crm.infra_billing_node_payment_in_48hrs" | "crm.infra_billing_node_payment_in_24hrs" | "crm.infra_billing_node_payment_due_today" | "crm.infra_billing_node_payment_overdue_24hrs" | "crm.infra_billing_node_payment_overdue_48hrs" | "crm.infra_billing_node_payment_overdue_7_days";
    timestamp: Date;
}, {
    data: {
        nodeName: string;
        loginUrl: string;
        nextBillingAt: string;
        providerName: string;
    };
    scope: "crm";
    event: "crm.infra_billing_node_payment_in_7_days" | "crm.infra_billing_node_payment_in_48hrs" | "crm.infra_billing_node_payment_in_24hrs" | "crm.infra_billing_node_payment_due_today" | "crm.infra_billing_node_payment_overdue_24hrs" | "crm.infra_billing_node_payment_overdue_48hrs" | "crm.infra_billing_node_payment_overdue_7_days";
    timestamp: string;
}>, z.ZodObject<{
    scope: z.ZodLiteral<"torrent_blocker">;
    event: z.ZodEnum<["torrent_blocker.report", ..."torrent_blocker.report"[]]>;
    timestamp: z.ZodEffects<z.ZodString, Date, string>;
    data: z.ZodObject<{
        node: z.ZodObject<{
            uuid: z.ZodString;
            name: z.ZodString;
            address: z.ZodString;
            port: z.ZodNullable<z.ZodNumber>;
            proxyUrl: z.ZodNullable<z.ZodString>;
            isConnected: z.ZodBoolean;
            isDisabled: z.ZodBoolean;
            isConnecting: z.ZodBoolean;
            lastStatusChange: z.ZodNullable<z.ZodEffects<z.ZodString, Date, string>>;
            lastStatusMessage: z.ZodNullable<z.ZodString>;

            singboxVersion: z.ZodNullable<z.ZodString>;

            nodeVersion: z.ZodNullable<z.ZodString>;
            isTrafficTrackingActive: z.ZodBoolean;
            trafficResetDay: z.ZodNullable<z.ZodNumber>;
            trafficLimitBytes: z.ZodNullable<z.ZodNumber>;
            trafficUsedBytes: z.ZodNullable<z.ZodNumber>;
            notifyPercent: z.ZodNullable<z.ZodNumber>;
            viewPosition: z.ZodNumber;
            countryCode: z.ZodString;
            consumptionMultiplier: z.ZodNumber;
            nodeConsumptionMultiplier: z.ZodNumber;
            tags: z.ZodArray<z.ZodString, "many">;
            createdAt: z.ZodEffects<z.ZodString, Date, string>;
            updatedAt: z.ZodEffects<z.ZodString, Date, string>;
            configProfile: z.ZodObject<{
                activeConfigProfileUuid: z.ZodNullable<z.ZodString>;
                activeInbounds: z.ZodArray<z.ZodObject<{
                    uuid: z.ZodString;
                    profileUuid: z.ZodString;
                    tag: z.ZodString;
                    type: z.ZodString;
                    network: z.ZodNullable<z.ZodString>;
                    security: z.ZodNullable<z.ZodString>;
                    port: z.ZodNullable<z.ZodNumber>;
                    rawInbound: z.ZodNullable<z.ZodUnknown>;
                }, "strip", z.ZodTypeAny, {
                    uuid: string;
                    type: string;
                    profileUuid: string;
                    tag: string;
                    network: string | null;
                    security: string | null;
                    port: number | null;
                    rawInbound?: unknown;
                }, {
                    uuid: string;
                    type: string;
                    profileUuid: string;
                    tag: string;
                    network: string | null;
                    security: string | null;
                    port: number | null;
                    rawInbound?: unknown;
                }>, "many">;
            }, "strip", z.ZodTypeAny, {
                activeConfigProfileUuid: string | null;
                activeInbounds: {
                    uuid: string;
                    type: string;
                    profileUuid: string;
                    tag: string;
                    network: string | null;
                    security: string | null;
                    port: number | null;
                    rawInbound?: unknown;
                }[];
            }, {
                activeConfigProfileUuid: string | null;
                activeInbounds: {
                    uuid: string;
                    type: string;
                    profileUuid: string;
                    tag: string;
                    network: string | null;
                    security: string | null;
                    port: number | null;
                    rawInbound?: unknown;
                }[];
            }>;
            providerUuid: z.ZodNullable<z.ZodString>;
            provider: z.ZodNullable<z.ZodObject<{
                uuid: z.ZodString;
                name: z.ZodString;
                faviconLink: z.ZodNullable<z.ZodString>;
                loginUrl: z.ZodNullable<z.ZodString>;
                createdAt: z.ZodEffects<z.ZodString, Date, string>;
                updatedAt: z.ZodEffects<z.ZodString, Date, string>;
            }, "strip", z.ZodTypeAny, {
                uuid: string;
                name: string;
                createdAt: Date;
                updatedAt: Date;
                faviconLink: string | null;
                loginUrl: string | null;
            }, {
                uuid: string;
                name: string;
                createdAt: string;
                updatedAt: string;
                faviconLink: string | null;
                loginUrl: string | null;
            }>>;
            activePluginUuid: z.ZodNullable<z.ZodString>;
            system: z.ZodNullable<z.ZodObject<{
                info: z.ZodObject<{
                    arch: z.ZodString;
                    cpus: z.ZodNumber;
                    cpuModel: z.ZodString;
                    memoryTotal: z.ZodNumber;
                    hostname: z.ZodString;
                    platform: z.ZodString;
                    release: z.ZodString;
                    type: z.ZodString;
                    version: z.ZodString;
                    networkInterfaces: z.ZodArray<z.ZodString, "many">;
                }, "strip", z.ZodTypeAny, {
                    type: string;
                    version: string;
                    platform: string;
                    arch: string;
                    cpus: number;
                    cpuModel: string;
                    memoryTotal: number;
                    hostname: string;
                    release: string;
                    networkInterfaces: string[];
                }, {
                    type: string;
                    version: string;
                    platform: string;
                    arch: string;
                    cpus: number;
                    cpuModel: string;
                    memoryTotal: number;
                    hostname: string;
                    release: string;
                    networkInterfaces: string[];
                }>;
                stats: z.ZodObject<{
                    memoryFree: z.ZodNumber;
                    memoryUsed: z.ZodNumber;
                    uptime: z.ZodNumber;
                    loadAvg: z.ZodArray<z.ZodNumber, "many">;
                    interface: z.ZodNullable<z.ZodObject<{
                        interface: z.ZodString;
                        rxBytesPerSec: z.ZodNumber;
                        txBytesPerSec: z.ZodNumber;
                        rxTotal: z.ZodNumber;
                        txTotal: z.ZodNumber;
                    }, "strip", z.ZodTypeAny, {
                        interface: string;
                        rxBytesPerSec: number;
                        txBytesPerSec: number;
                        rxTotal: number;
                        txTotal: number;
                    }, {
                        interface: string;
                        rxBytesPerSec: number;
                        txBytesPerSec: number;
                        rxTotal: number;
                        txTotal: number;
                    }>>;
                }, "strip", z.ZodTypeAny, {
                    interface: {
                        interface: string;
                        rxBytesPerSec: number;
                        txBytesPerSec: number;
                        rxTotal: number;
                        txTotal: number;
                    } | null;
                    memoryFree: number;
                    memoryUsed: number;
                    uptime: number;
                    loadAvg: number[];
                }, {
                    interface: {
                        interface: string;
                        rxBytesPerSec: number;
                        txBytesPerSec: number;
                        rxTotal: number;
                        txTotal: number;
                    } | null;
                    memoryFree: number;
                    memoryUsed: number;
                    uptime: number;
                    loadAvg: number[];
                }>;
            }, "strip", z.ZodTypeAny, {
                stats: {
                    interface: {
                        interface: string;
                        rxBytesPerSec: number;
                        txBytesPerSec: number;
                        rxTotal: number;
                        txTotal: number;
                    } | null;
                    memoryFree: number;
                    memoryUsed: number;
                    uptime: number;
                    loadAvg: number[];
                };
                info: {
                    type: string;
                    version: string;
                    platform: string;
                    arch: string;
                    cpus: number;
                    cpuModel: string;
                    memoryTotal: number;
                    hostname: string;
                    release: string;
                    networkInterfaces: string[];
                };
            }, {
                stats: {
                    interface: {
                        interface: string;
                        rxBytesPerSec: number;
                        txBytesPerSec: number;
                        rxTotal: number;
                        txTotal: number;
                    } | null;
                    memoryFree: number;
                    memoryUsed: number;
                    uptime: number;
                    loadAvg: number[];
                };
                info: {
                    type: string;
                    version: string;
                    platform: string;
                    arch: string;
                    cpus: number;
                    cpuModel: string;
                    memoryTotal: number;
                    hostname: string;
                    release: string;
                    networkInterfaces: string[];
                };
            }>>;
            versions: z.ZodNullable<z.ZodObject<{
                singbox: z.ZodString;
                node: z.ZodString;
            }, "strip", z.ZodTypeAny, {
                node: string;
                singbox: string;
            }, {
                node: string;
                singbox: string;
            }>>;
            singboxUptime: z.ZodNumber;
            usersOnline: z.ZodNumber;
            note: z.ZodNullable<z.ZodString>;
        }, "strip", z.ZodTypeAny, {
            tags: string[];
            system: {
                stats: {
                    interface: {
                        interface: string;
                        rxBytesPerSec: number;
                        txBytesPerSec: number;
                        rxTotal: number;
                        txTotal: number;
                    } | null;
                    memoryFree: number;
                    memoryUsed: number;
                    uptime: number;
                    loadAvg: number[];
                };
                info: {
                    type: string;
                    version: string;
                    platform: string;
                    arch: string;
                    cpus: number;
                    cpuModel: string;
                    memoryTotal: number;
                    hostname: string;
                    release: string;
                    networkInterfaces: string[];
                };
            } | null;
            uuid: string;
            name: string;
            createdAt: Date;
            updatedAt: Date;
            provider: {
                uuid: string;
                name: string;
                createdAt: Date;
                updatedAt: Date;
                faviconLink: string | null;
                loginUrl: string | null;
            } | null;
            countryCode: string;
            port: number | null;
            viewPosition: number;
            trafficLimitBytes: number | null;
            address: string;
            isDisabled: boolean;
            proxyUrl: string | null;
            isConnected: boolean;
            isConnecting: boolean;
            lastStatusChange: Date | null;
            lastStatusMessage: string | null;

            singboxVersion: string | null;

            nodeVersion: string | null;
            isTrafficTrackingActive: boolean;
            trafficResetDay: number | null;
            trafficUsedBytes: number | null;
            notifyPercent: number | null;
            consumptionMultiplier: number;
            nodeConsumptionMultiplier: number;
            configProfile: {
                activeConfigProfileUuid: string | null;
                activeInbounds: {
                    uuid: string;
                    type: string;
                    profileUuid: string;
                    tag: string;
                    network: string | null;
                    security: string | null;
                    port: number | null;
                    rawInbound?: unknown;
                }[];
            };
            providerUuid: string | null;
            activePluginUuid: string | null;
            versions: {
                node: string;
                singbox: string;
            } | null;
            singboxUptime: number;
            usersOnline: number;
            note: string | null;
        }, {
            tags: string[];
            system: {
                stats: {
                    interface: {
                        interface: string;
                        rxBytesPerSec: number;
                        txBytesPerSec: number;
                        rxTotal: number;
                        txTotal: number;
                    } | null;
                    memoryFree: number;
                    memoryUsed: number;
                    uptime: number;
                    loadAvg: number[];
                };
                info: {
                    type: string;
                    version: string;
                    platform: string;
                    arch: string;
                    cpus: number;
                    cpuModel: string;
                    memoryTotal: number;
                    hostname: string;
                    release: string;
                    networkInterfaces: string[];
                };
            } | null;
            uuid: string;
            name: string;
            createdAt: string;
            updatedAt: string;
            provider: {
                uuid: string;
                name: string;
                createdAt: string;
                updatedAt: string;
                faviconLink: string | null;
                loginUrl: string | null;
            } | null;
            countryCode: string;
            port: number | null;
            viewPosition: number;
            trafficLimitBytes: number | null;
            address: string;
            isDisabled: boolean;
            proxyUrl: string | null;
            isConnected: boolean;
            isConnecting: boolean;
            lastStatusChange: string | null;
            lastStatusMessage: string | null;

            singboxVersion: string | null;

            nodeVersion: string | null;
            isTrafficTrackingActive: boolean;
            trafficResetDay: number | null;
            trafficUsedBytes: number | null;
            notifyPercent: number | null;
            consumptionMultiplier: number;
            nodeConsumptionMultiplier: number;
            configProfile: {
                activeConfigProfileUuid: string | null;
                activeInbounds: {
                    uuid: string;
                    type: string;
                    profileUuid: string;
                    tag: string;
                    network: string | null;
                    security: string | null;
                    port: number | null;
                    rawInbound?: unknown;
                }[];
            };
            providerUuid: string | null;
            activePluginUuid: string | null;
            versions: {
                node: string;
                singbox: string;
            } | null;
            singboxUptime: number;
            usersOnline: number;
            note: string | null;
        }>;
        user: z.ZodObject<{
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
        }>;
        report: z.ZodObject<{
            actionReport: z.ZodObject<{
                blocked: z.ZodBoolean;
                ip: z.ZodString;
                blockDuration: z.ZodNumber;
                willUnblockAt: z.ZodEffects<z.ZodString, Date, string>;
                userId: z.ZodString;
                processedAt: z.ZodEffects<z.ZodString, Date, string>;
            }, "strip", z.ZodTypeAny, {
                ip: string;
                userId: string;
                blocked: boolean;
                blockDuration: number;
                willUnblockAt: Date;
                processedAt: Date;
            }, {
                ip: string;
                userId: string;
                blocked: boolean;
                blockDuration: number;
                willUnblockAt: string;
                processedAt: string;
            }>;
            xrayReport: z.ZodObject<{
                email: z.ZodNullable<z.ZodString>;
                level: z.ZodNullable<z.ZodNumber>;
                protocol: z.ZodNullable<z.ZodString>;
                network: z.ZodString;
                source: z.ZodNullable<z.ZodString>;
                destination: z.ZodString;
                routeTarget: z.ZodNullable<z.ZodString>;
                originalTarget: z.ZodNullable<z.ZodString>;
                inboundTag: z.ZodNullable<z.ZodString>;
                inboundName: z.ZodNullable<z.ZodString>;
                inboundLocal: z.ZodNullable<z.ZodString>;
                outboundTag: z.ZodNullable<z.ZodString>;
                ts: z.ZodNumber;
            }, "strip", z.ZodTypeAny, {
                network: string;
                email: string | null;
                protocol: string | null;
                inboundTag: string | null;
                level: number | null;
                source: string | null;
                destination: string;
                routeTarget: string | null;
                originalTarget: string | null;
                inboundName: string | null;
                inboundLocal: string | null;
                outboundTag: string | null;
                ts: number;
            }, {
                network: string;
                email: string | null;
                protocol: string | null;
                inboundTag: string | null;
                level: number | null;
                source: string | null;
                destination: string;
                routeTarget: string | null;
                originalTarget: string | null;
                inboundName: string | null;
                inboundLocal: string | null;
                outboundTag: string | null;
                ts: number;
            }>;
        }, "strip", z.ZodTypeAny, {
            actionReport: {
                ip: string;
                userId: string;
                blocked: boolean;
                blockDuration: number;
                willUnblockAt: Date;
                processedAt: Date;
            };
            xrayReport: {
                network: string;
                email: string | null;
                protocol: string | null;
                inboundTag: string | null;
                level: number | null;
                source: string | null;
                destination: string;
                routeTarget: string | null;
                originalTarget: string | null;
                inboundName: string | null;
                inboundLocal: string | null;
                outboundTag: string | null;
                ts: number;
            };
        }, {
            actionReport: {
                ip: string;
                userId: string;
                blocked: boolean;
                blockDuration: number;
                willUnblockAt: string;
                processedAt: string;
            };
            xrayReport: {
                network: string;
                email: string | null;
                protocol: string | null;
                inboundTag: string | null;
                level: number | null;
                source: string | null;
                destination: string;
                routeTarget: string | null;
                originalTarget: string | null;
                inboundName: string | null;
                inboundLocal: string | null;
                outboundTag: string | null;
                ts: number;
            };
        }>;
    }, "strip", z.ZodTypeAny, {
        user: {
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
        };
        node: {
            tags: string[];
            system: {
                stats: {
                    interface: {
                        interface: string;
                        rxBytesPerSec: number;
                        txBytesPerSec: number;
                        rxTotal: number;
                        txTotal: number;
                    } | null;
                    memoryFree: number;
                    memoryUsed: number;
                    uptime: number;
                    loadAvg: number[];
                };
                info: {
                    type: string;
                    version: string;
                    platform: string;
                    arch: string;
                    cpus: number;
                    cpuModel: string;
                    memoryTotal: number;
                    hostname: string;
                    release: string;
                    networkInterfaces: string[];
                };
            } | null;
            uuid: string;
            name: string;
            createdAt: Date;
            updatedAt: Date;
            provider: {
                uuid: string;
                name: string;
                createdAt: Date;
                updatedAt: Date;
                faviconLink: string | null;
                loginUrl: string | null;
            } | null;
            countryCode: string;
            port: number | null;
            viewPosition: number;
            trafficLimitBytes: number | null;
            address: string;
            isDisabled: boolean;
            proxyUrl: string | null;
            isConnected: boolean;
            isConnecting: boolean;
            lastStatusChange: Date | null;
            lastStatusMessage: string | null;

            singboxVersion: string | null;

            nodeVersion: string | null;
            isTrafficTrackingActive: boolean;
            trafficResetDay: number | null;
            trafficUsedBytes: number | null;
            notifyPercent: number | null;
            consumptionMultiplier: number;
            nodeConsumptionMultiplier: number;
            configProfile: {
                activeConfigProfileUuid: string | null;
                activeInbounds: {
                    uuid: string;
                    type: string;
                    profileUuid: string;
                    tag: string;
                    network: string | null;
                    security: string | null;
                    port: number | null;
                    rawInbound?: unknown;
                }[];
            };
            providerUuid: string | null;
            activePluginUuid: string | null;
            versions: {
                node: string;
                singbox: string;
            } | null;
            singboxUptime: number;
            usersOnline: number;
            note: string | null;
        };
        report: {
            actionReport: {
                ip: string;
                userId: string;
                blocked: boolean;
                blockDuration: number;
                willUnblockAt: Date;
                processedAt: Date;
            };
            xrayReport: {
                network: string;
                email: string | null;
                protocol: string | null;
                inboundTag: string | null;
                level: number | null;
                source: string | null;
                destination: string;
                routeTarget: string | null;
                originalTarget: string | null;
                inboundName: string | null;
                inboundLocal: string | null;
                outboundTag: string | null;
                ts: number;
            };
        };
    }, {
        user: {
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
        };
        node: {
            tags: string[];
            system: {
                stats: {
                    interface: {
                        interface: string;
                        rxBytesPerSec: number;
                        txBytesPerSec: number;
                        rxTotal: number;
                        txTotal: number;
                    } | null;
                    memoryFree: number;
                    memoryUsed: number;
                    uptime: number;
                    loadAvg: number[];
                };
                info: {
                    type: string;
                    version: string;
                    platform: string;
                    arch: string;
                    cpus: number;
                    cpuModel: string;
                    memoryTotal: number;
                    hostname: string;
                    release: string;
                    networkInterfaces: string[];
                };
            } | null;
            uuid: string;
            name: string;
            createdAt: string;
            updatedAt: string;
            provider: {
                uuid: string;
                name: string;
                createdAt: string;
                updatedAt: string;
                faviconLink: string | null;
                loginUrl: string | null;
            } | null;
            countryCode: string;
            port: number | null;
            viewPosition: number;
            trafficLimitBytes: number | null;
            address: string;
            isDisabled: boolean;
            proxyUrl: string | null;
            isConnected: boolean;
            isConnecting: boolean;
            lastStatusChange: string | null;
            lastStatusMessage: string | null;

            singboxVersion: string | null;

            nodeVersion: string | null;
            isTrafficTrackingActive: boolean;
            trafficResetDay: number | null;
            trafficUsedBytes: number | null;
            notifyPercent: number | null;
            consumptionMultiplier: number;
            nodeConsumptionMultiplier: number;
            configProfile: {
                activeConfigProfileUuid: string | null;
                activeInbounds: {
                    uuid: string;
                    type: string;
                    profileUuid: string;
                    tag: string;
                    network: string | null;
                    security: string | null;
                    port: number | null;
                    rawInbound?: unknown;
                }[];
            };
            providerUuid: string | null;
            activePluginUuid: string | null;
            versions: {
                node: string;
                singbox: string;
            } | null;
            singboxUptime: number;
            usersOnline: number;
            note: string | null;
        };
        report: {
            actionReport: {
                ip: string;
                userId: string;
                blocked: boolean;
                blockDuration: number;
                willUnblockAt: string;
                processedAt: string;
            };
            xrayReport: {
                network: string;
                email: string | null;
                protocol: string | null;
                inboundTag: string | null;
                level: number | null;
                source: string | null;
                destination: string;
                routeTarget: string | null;
                originalTarget: string | null;
                inboundName: string | null;
                inboundLocal: string | null;
                outboundTag: string | null;
                ts: number;
            };
        };
    }>;
}, "strip", z.ZodTypeAny, {
    data: {
        user: {
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
        };
        node: {
            tags: string[];
            system: {
                stats: {
                    interface: {
                        interface: string;
                        rxBytesPerSec: number;
                        txBytesPerSec: number;
                        rxTotal: number;
                        txTotal: number;
                    } | null;
                    memoryFree: number;
                    memoryUsed: number;
                    uptime: number;
                    loadAvg: number[];
                };
                info: {
                    type: string;
                    version: string;
                    platform: string;
                    arch: string;
                    cpus: number;
                    cpuModel: string;
                    memoryTotal: number;
                    hostname: string;
                    release: string;
                    networkInterfaces: string[];
                };
            } | null;
            uuid: string;
            name: string;
            createdAt: Date;
            updatedAt: Date;
            provider: {
                uuid: string;
                name: string;
                createdAt: Date;
                updatedAt: Date;
                faviconLink: string | null;
                loginUrl: string | null;
            } | null;
            countryCode: string;
            port: number | null;
            viewPosition: number;
            trafficLimitBytes: number | null;
            address: string;
            isDisabled: boolean;
            proxyUrl: string | null;
            isConnected: boolean;
            isConnecting: boolean;
            lastStatusChange: Date | null;
            lastStatusMessage: string | null;

            singboxVersion: string | null;

            nodeVersion: string | null;
            isTrafficTrackingActive: boolean;
            trafficResetDay: number | null;
            trafficUsedBytes: number | null;
            notifyPercent: number | null;
            consumptionMultiplier: number;
            nodeConsumptionMultiplier: number;
            configProfile: {
                activeConfigProfileUuid: string | null;
                activeInbounds: {
                    uuid: string;
                    type: string;
                    profileUuid: string;
                    tag: string;
                    network: string | null;
                    security: string | null;
                    port: number | null;
                    rawInbound?: unknown;
                }[];
            };
            providerUuid: string | null;
            activePluginUuid: string | null;
            versions: {
                node: string;
                singbox: string;
            } | null;
            singboxUptime: number;
            usersOnline: number;
            note: string | null;
        };
        report: {
            actionReport: {
                ip: string;
                userId: string;
                blocked: boolean;
                blockDuration: number;
                willUnblockAt: Date;
                processedAt: Date;
            };
            xrayReport: {
                network: string;
                email: string | null;
                protocol: string | null;
                inboundTag: string | null;
                level: number | null;
                source: string | null;
                destination: string;
                routeTarget: string | null;
                originalTarget: string | null;
                inboundName: string | null;
                inboundLocal: string | null;
                outboundTag: string | null;
                ts: number;
            };
        };
    };
    scope: "torrent_blocker";
    event: "torrent_blocker.report";
    timestamp: Date;
}, {
    data: {
        user: {
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
        };
        node: {
            tags: string[];
            system: {
                stats: {
                    interface: {
                        interface: string;
                        rxBytesPerSec: number;
                        txBytesPerSec: number;
                        rxTotal: number;
                        txTotal: number;
                    } | null;
                    memoryFree: number;
                    memoryUsed: number;
                    uptime: number;
                    loadAvg: number[];
                };
                info: {
                    type: string;
                    version: string;
                    platform: string;
                    arch: string;
                    cpus: number;
                    cpuModel: string;
                    memoryTotal: number;
                    hostname: string;
                    release: string;
                    networkInterfaces: string[];
                };
            } | null;
            uuid: string;
            name: string;
            createdAt: string;
            updatedAt: string;
            provider: {
                uuid: string;
                name: string;
                createdAt: string;
                updatedAt: string;
                faviconLink: string | null;
                loginUrl: string | null;
            } | null;
            countryCode: string;
            port: number | null;
            viewPosition: number;
            trafficLimitBytes: number | null;
            address: string;
            isDisabled: boolean;
            proxyUrl: string | null;
            isConnected: boolean;
            isConnecting: boolean;
            lastStatusChange: string | null;
            lastStatusMessage: string | null;

            singboxVersion: string | null;

            nodeVersion: string | null;
            isTrafficTrackingActive: boolean;
            trafficResetDay: number | null;
            trafficUsedBytes: number | null;
            notifyPercent: number | null;
            consumptionMultiplier: number;
            nodeConsumptionMultiplier: number;
            configProfile: {
                activeConfigProfileUuid: string | null;
                activeInbounds: {
                    uuid: string;
                    type: string;
                    profileUuid: string;
                    tag: string;
                    network: string | null;
                    security: string | null;
                    port: number | null;
                    rawInbound?: unknown;
                }[];
            };
            providerUuid: string | null;
            activePluginUuid: string | null;
            versions: {
                node: string;
                singbox: string;
            } | null;
            singboxUptime: number;
            usersOnline: number;
            note: string | null;
        };
        report: {
            actionReport: {
                ip: string;
                userId: string;
                blocked: boolean;
                blockDuration: number;
                willUnblockAt: string;
                processedAt: string;
            };
            xrayReport: {
                network: string;
                email: string | null;
                protocol: string | null;
                inboundTag: string | null;
                level: number | null;
                source: string | null;
                destination: string;
                routeTarget: string | null;
                originalTarget: string | null;
                inboundName: string | null;
                inboundLocal: string | null;
                outboundTag: string | null;
                ts: number;
            };
        };
    };
    scope: "torrent_blocker";
    event: "torrent_blocker.report";
    timestamp: string;
}>]>;
export type TExodusWebhookEvent = z.infer<typeof ExodusWebhookEventSchema>;
export type TExodusWebhookUserEvent = z.infer<typeof ExodusWebhookUserEvents>;
export type TExodusWebhookNodeEvent = z.infer<typeof ExodusWebhookNodeEvents>;
export type TExodusWebhookServiceEvent = z.infer<typeof ExodusWebhookServiceEvents>;
export type TExodusWebhookErrorsEvent = z.infer<typeof ExodusWebhookErrorsEvents>;
export type TExodusWebhookCrmEvent = z.infer<typeof ExodusWebhookCrmEvents>;
export type TExodusWebhookUserHwidDevicesEvent = z.infer<typeof ExodusWebhookUserHwidDevicesEvents>;
export type TExodusWebhookTorrentBlockerEvent = z.infer<typeof ExodusWebhookTorrentBlockerEvents>;
//# sourceMappingURL=webhook.schema.d.ts.map