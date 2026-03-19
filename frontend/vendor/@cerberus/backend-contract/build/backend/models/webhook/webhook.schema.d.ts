import z from 'zod';
export declare const CerberusWebhookUserEvents: z.ZodObject<{
    scope: z.ZodLiteral<"user">;
    event: z.ZodEnum<["user.created" | "user.modified" | "user.deleted" | "user.revoked" | "user.disabled" | "user.enabled" | "user.limited" | "user.expired" | "user.traffic_reset" | "user.expires_in_72_hours" | "user.expires_in_48_hours" | "user.expires_in_24_hours" | "user.expired_24_hours_ago" | "user.first_connected" | "user.bandwidth_usage_threshold_reached" | "user.not_connected", ...("user.created" | "user.modified" | "user.deleted" | "user.revoked" | "user.disabled" | "user.enabled" | "user.limited" | "user.expired" | "user.traffic_reset" | "user.expires_in_72_hours" | "user.expires_in_48_hours" | "user.expires_in_24_hours" | "user.expired_24_hours_ago" | "user.first_connected" | "user.bandwidth_usage_threshold_reached" | "user.not_connected")[]]>;
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
    }>;
    meta: z.ZodNullable<z.ZodObject<{
        notConnectedAfterHours: z.ZodOptional<z.ZodNullable<z.ZodNumber>>;
    }, "strip", z.ZodTypeAny, {
        notConnectedAfterHours?: number | null | undefined;
    }, {
        notConnectedAfterHours?: number | null | undefined;
    }>>;
}, "strip", z.ZodTypeAny, {
    data: {
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
    };
    scope: "user";
    event: "user.created" | "user.modified" | "user.deleted" | "user.revoked" | "user.disabled" | "user.enabled" | "user.limited" | "user.expired" | "user.traffic_reset" | "user.expires_in_72_hours" | "user.expires_in_48_hours" | "user.expires_in_24_hours" | "user.expired_24_hours_ago" | "user.first_connected" | "user.bandwidth_usage_threshold_reached" | "user.not_connected";
    timestamp: Date;
    meta: {
        notConnectedAfterHours?: number | null | undefined;
    } | null;
}, {
    data: {
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
    };
    scope: "user";
    event: "user.created" | "user.modified" | "user.deleted" | "user.revoked" | "user.disabled" | "user.enabled" | "user.limited" | "user.expired" | "user.traffic_reset" | "user.expires_in_72_hours" | "user.expires_in_48_hours" | "user.expires_in_24_hours" | "user.expired_24_hours_ago" | "user.first_connected" | "user.bandwidth_usage_threshold_reached" | "user.not_connected";
    timestamp: string;
    meta: {
        notConnectedAfterHours?: number | null | undefined;
    } | null;
}>;
export declare const CerberusWebhookUserHwidDevicesEvents: z.ZodObject<{
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
        }>;
        hwidUserDevice: z.ZodObject<{
            hwid: z.ZodString;
            userUuid: z.ZodString;
            platform: z.ZodNullable<z.ZodString>;
            osVersion: z.ZodNullable<z.ZodString>;
            deviceModel: z.ZodNullable<z.ZodString>;
            userAgent: z.ZodNullable<z.ZodString>;
            createdAt: z.ZodEffects<z.ZodString, Date, string>;
            updatedAt: z.ZodEffects<z.ZodString, Date, string>;
        }, "strip", z.ZodTypeAny, {
            hwid: string;
            createdAt: Date;
            updatedAt: Date;
            userUuid: string;
            platform: string | null;
            osVersion: string | null;
            deviceModel: string | null;
            userAgent: string | null;
        }, {
            hwid: string;
            createdAt: string;
            updatedAt: string;
            userUuid: string;
            platform: string | null;
            osVersion: string | null;
            deviceModel: string | null;
            userAgent: string | null;
        }>;
    }, "strip", z.ZodTypeAny, {
        user: {
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
        };
        hwidUserDevice: {
            hwid: string;
            createdAt: Date;
            updatedAt: Date;
            userUuid: string;
            platform: string | null;
            osVersion: string | null;
            deviceModel: string | null;
            userAgent: string | null;
        };
    }, {
        user: {
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
        };
        hwidUserDevice: {
            hwid: string;
            createdAt: string;
            updatedAt: string;
            userUuid: string;
            platform: string | null;
            osVersion: string | null;
            deviceModel: string | null;
            userAgent: string | null;
        };
    }>;
}, "strip", z.ZodTypeAny, {
    data: {
        user: {
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
        };
        hwidUserDevice: {
            hwid: string;
            createdAt: Date;
            updatedAt: Date;
            userUuid: string;
            platform: string | null;
            osVersion: string | null;
            deviceModel: string | null;
            userAgent: string | null;
        };
    };
    scope: "user_hwid_devices";
    event: "user_hwid_devices.added" | "user_hwid_devices.deleted";
    timestamp: Date;
}, {
    data: {
        user: {
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
        };
        hwidUserDevice: {
            hwid: string;
            createdAt: string;
            updatedAt: string;
            userUuid: string;
            platform: string | null;
            osVersion: string | null;
            deviceModel: string | null;
            userAgent: string | null;
        };
    };
    scope: "user_hwid_devices";
    event: "user_hwid_devices.added" | "user_hwid_devices.deleted";
    timestamp: string;
}>;
export declare const CerberusWebhookNodeEvents: z.ZodObject<{
    scope: z.ZodLiteral<"node">;
    event: z.ZodEnum<["node.created" | "node.modified" | "node.disabled" | "node.enabled" | "node.deleted" | "node.connection_lost" | "node.connection_restored" | "node.traffic_notify", ...("node.created" | "node.modified" | "node.disabled" | "node.enabled" | "node.deleted" | "node.connection_lost" | "node.connection_restored" | "node.traffic_notify")[]]>;
    timestamp: z.ZodEffects<z.ZodString, Date, string>;
    data: z.ZodObject<{
        uuid: z.ZodString;
        name: z.ZodString;
        address: z.ZodString;
        port: z.ZodNullable<z.ZodNumber>;
        isConnected: z.ZodBoolean;
        isDisabled: z.ZodBoolean;
        isConnecting: z.ZodBoolean;
        lastStatusChange: z.ZodNullable<z.ZodEffects<z.ZodString, Date, string>>;
        lastStatusMessage: z.ZodNullable<z.ZodString>;
        xrayVersion: z.ZodNullable<z.ZodString>;
        nodeVersion: z.ZodNullable<z.ZodString>;
        xrayUptime: z.ZodString;
        isTrafficTrackingActive: z.ZodBoolean;
        trafficResetDay: z.ZodNullable<z.ZodNumber>;
        trafficLimitBytes: z.ZodNullable<z.ZodNumber>;
        trafficUsedBytes: z.ZodNullable<z.ZodNumber>;
        notifyPercent: z.ZodNullable<z.ZodNumber>;
        usersOnline: z.ZodNullable<z.ZodNumber>;
        viewPosition: z.ZodNumber;
        countryCode: z.ZodString;
        consumptionMultiplier: z.ZodNumber;
        tags: z.ZodArray<z.ZodString, "many">;
        cpuCount: z.ZodNullable<z.ZodNumber>;
        cpuModel: z.ZodNullable<z.ZodString>;
        totalRam: z.ZodNullable<z.ZodString>;
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
                type: string;
                uuid: string;
                profileUuid: string;
                tag: string;
                network: string | null;
                security: string | null;
                port: number | null;
                rawInbound?: unknown;
            }, {
                type: string;
                uuid: string;
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
                type: string;
                uuid: string;
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
                type: string;
                uuid: string;
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
            createdAt: Date;
            updatedAt: Date;
            name: string;
            faviconLink: string | null;
            loginUrl: string | null;
        }, {
            uuid: string;
            createdAt: string;
            updatedAt: string;
            name: string;
            faviconLink: string | null;
            loginUrl: string | null;
        }>>;
    }, "strip", z.ZodTypeAny, {
        tags: string[];
        uuid: string;
        createdAt: Date;
        updatedAt: Date;
        provider: {
            uuid: string;
            createdAt: Date;
            updatedAt: Date;
            name: string;
            faviconLink: string | null;
            loginUrl: string | null;
        } | null;
        countryCode: string;
        name: string;
        port: number | null;
        viewPosition: number;
        trafficLimitBytes: number | null;
        address: string;
        isDisabled: boolean;
        isConnected: boolean;
        isConnecting: boolean;
        lastStatusChange: Date | null;
        lastStatusMessage: string | null;
        xrayVersion: string | null;
        nodeVersion: string | null;
        xrayUptime: string;
        isTrafficTrackingActive: boolean;
        trafficResetDay: number | null;
        trafficUsedBytes: number | null;
        notifyPercent: number | null;
        usersOnline: number | null;
        consumptionMultiplier: number;
        cpuCount: number | null;
        cpuModel: string | null;
        totalRam: string | null;
        configProfile: {
            activeConfigProfileUuid: string | null;
            activeInbounds: {
                type: string;
                uuid: string;
                profileUuid: string;
                tag: string;
                network: string | null;
                security: string | null;
                port: number | null;
                rawInbound?: unknown;
            }[];
        };
        providerUuid: string | null;
    }, {
        tags: string[];
        uuid: string;
        createdAt: string;
        updatedAt: string;
        provider: {
            uuid: string;
            createdAt: string;
            updatedAt: string;
            name: string;
            faviconLink: string | null;
            loginUrl: string | null;
        } | null;
        countryCode: string;
        name: string;
        port: number | null;
        viewPosition: number;
        trafficLimitBytes: number | null;
        address: string;
        isDisabled: boolean;
        isConnected: boolean;
        isConnecting: boolean;
        lastStatusChange: string | null;
        lastStatusMessage: string | null;
        xrayVersion: string | null;
        nodeVersion: string | null;
        xrayUptime: string;
        isTrafficTrackingActive: boolean;
        trafficResetDay: number | null;
        trafficUsedBytes: number | null;
        notifyPercent: number | null;
        usersOnline: number | null;
        consumptionMultiplier: number;
        cpuCount: number | null;
        cpuModel: string | null;
        totalRam: string | null;
        configProfile: {
            activeConfigProfileUuid: string | null;
            activeInbounds: {
                type: string;
                uuid: string;
                profileUuid: string;
                tag: string;
                network: string | null;
                security: string | null;
                port: number | null;
                rawInbound?: unknown;
            }[];
        };
        providerUuid: string | null;
    }>;
}, "strip", z.ZodTypeAny, {
    data: {
        tags: string[];
        uuid: string;
        createdAt: Date;
        updatedAt: Date;
        provider: {
            uuid: string;
            createdAt: Date;
            updatedAt: Date;
            name: string;
            faviconLink: string | null;
            loginUrl: string | null;
        } | null;
        countryCode: string;
        name: string;
        port: number | null;
        viewPosition: number;
        trafficLimitBytes: number | null;
        address: string;
        isDisabled: boolean;
        isConnected: boolean;
        isConnecting: boolean;
        lastStatusChange: Date | null;
        lastStatusMessage: string | null;
        xrayVersion: string | null;
        nodeVersion: string | null;
        xrayUptime: string;
        isTrafficTrackingActive: boolean;
        trafficResetDay: number | null;
        trafficUsedBytes: number | null;
        notifyPercent: number | null;
        usersOnline: number | null;
        consumptionMultiplier: number;
        cpuCount: number | null;
        cpuModel: string | null;
        totalRam: string | null;
        configProfile: {
            activeConfigProfileUuid: string | null;
            activeInbounds: {
                type: string;
                uuid: string;
                profileUuid: string;
                tag: string;
                network: string | null;
                security: string | null;
                port: number | null;
                rawInbound?: unknown;
            }[];
        };
        providerUuid: string | null;
    };
    scope: "node";
    event: "node.created" | "node.modified" | "node.disabled" | "node.enabled" | "node.deleted" | "node.connection_lost" | "node.connection_restored" | "node.traffic_notify";
    timestamp: Date;
}, {
    data: {
        tags: string[];
        uuid: string;
        createdAt: string;
        updatedAt: string;
        provider: {
            uuid: string;
            createdAt: string;
            updatedAt: string;
            name: string;
            faviconLink: string | null;
            loginUrl: string | null;
        } | null;
        countryCode: string;
        name: string;
        port: number | null;
        viewPosition: number;
        trafficLimitBytes: number | null;
        address: string;
        isDisabled: boolean;
        isConnected: boolean;
        isConnecting: boolean;
        lastStatusChange: string | null;
        lastStatusMessage: string | null;
        xrayVersion: string | null;
        nodeVersion: string | null;
        xrayUptime: string;
        isTrafficTrackingActive: boolean;
        trafficResetDay: number | null;
        trafficUsedBytes: number | null;
        notifyPercent: number | null;
        usersOnline: number | null;
        consumptionMultiplier: number;
        cpuCount: number | null;
        cpuModel: string | null;
        totalRam: string | null;
        configProfile: {
            activeConfigProfileUuid: string | null;
            activeInbounds: {
                type: string;
                uuid: string;
                profileUuid: string;
                tag: string;
                network: string | null;
                security: string | null;
                port: number | null;
                rawInbound?: unknown;
            }[];
        };
        providerUuid: string | null;
    };
    scope: "node";
    event: "node.created" | "node.modified" | "node.disabled" | "node.enabled" | "node.deleted" | "node.connection_lost" | "node.connection_restored" | "node.traffic_notify";
    timestamp: string;
}>;
export declare const CerberusWebhookServiceEvents: z.ZodObject<{
    scope: z.ZodLiteral<"service">;
    event: z.ZodEnum<["service.panel_started" | "service.login_attempt_failed" | "service.login_attempt_success", ...("service.panel_started" | "service.login_attempt_failed" | "service.login_attempt_success")[]]>;
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
            password?: string | undefined;
            description?: string | undefined;
        }, {
            username: string;
            ip: string;
            userAgent: string;
            password?: string | undefined;
            description?: string | undefined;
        }>>;
        panelVersion: z.ZodOptional<z.ZodString>;
    }, "strip", z.ZodTypeAny, {
        loginAttempt?: {
            username: string;
            ip: string;
            userAgent: string;
            password?: string | undefined;
            description?: string | undefined;
        } | undefined;
        panelVersion?: string | undefined;
    }, {
        loginAttempt?: {
            username: string;
            ip: string;
            userAgent: string;
            password?: string | undefined;
            description?: string | undefined;
        } | undefined;
        panelVersion?: string | undefined;
    }>;
}, "strip", z.ZodTypeAny, {
    data: {
        loginAttempt?: {
            username: string;
            ip: string;
            userAgent: string;
            password?: string | undefined;
            description?: string | undefined;
        } | undefined;
        panelVersion?: string | undefined;
    };
    scope: "service";
    event: "service.panel_started" | "service.login_attempt_failed" | "service.login_attempt_success";
    timestamp: Date;
}, {
    data: {
        loginAttempt?: {
            username: string;
            ip: string;
            userAgent: string;
            password?: string | undefined;
            description?: string | undefined;
        } | undefined;
        panelVersion?: string | undefined;
    };
    scope: "service";
    event: "service.panel_started" | "service.login_attempt_failed" | "service.login_attempt_success";
    timestamp: string;
}>;
export declare const CerberusWebhookErrorsEvents: z.ZodObject<{
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
export declare const CerberusWebhookCrmEvents: z.ZodObject<{
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
export declare const CerberusWebhookEventSchema: z.ZodDiscriminatedUnion<"scope", [z.ZodObject<{
    scope: z.ZodLiteral<"user">;
    event: z.ZodEnum<["user.created" | "user.modified" | "user.deleted" | "user.revoked" | "user.disabled" | "user.enabled" | "user.limited" | "user.expired" | "user.traffic_reset" | "user.expires_in_72_hours" | "user.expires_in_48_hours" | "user.expires_in_24_hours" | "user.expired_24_hours_ago" | "user.first_connected" | "user.bandwidth_usage_threshold_reached" | "user.not_connected", ...("user.created" | "user.modified" | "user.deleted" | "user.revoked" | "user.disabled" | "user.enabled" | "user.limited" | "user.expired" | "user.traffic_reset" | "user.expires_in_72_hours" | "user.expires_in_48_hours" | "user.expires_in_24_hours" | "user.expired_24_hours_ago" | "user.first_connected" | "user.bandwidth_usage_threshold_reached" | "user.not_connected")[]]>;
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
    }>;
    meta: z.ZodNullable<z.ZodObject<{
        notConnectedAfterHours: z.ZodOptional<z.ZodNullable<z.ZodNumber>>;
    }, "strip", z.ZodTypeAny, {
        notConnectedAfterHours?: number | null | undefined;
    }, {
        notConnectedAfterHours?: number | null | undefined;
    }>>;
}, "strip", z.ZodTypeAny, {
    data: {
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
    };
    scope: "user";
    event: "user.created" | "user.modified" | "user.deleted" | "user.revoked" | "user.disabled" | "user.enabled" | "user.limited" | "user.expired" | "user.traffic_reset" | "user.expires_in_72_hours" | "user.expires_in_48_hours" | "user.expires_in_24_hours" | "user.expired_24_hours_ago" | "user.first_connected" | "user.bandwidth_usage_threshold_reached" | "user.not_connected";
    timestamp: Date;
    meta: {
        notConnectedAfterHours?: number | null | undefined;
    } | null;
}, {
    data: {
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
    };
    scope: "user";
    event: "user.created" | "user.modified" | "user.deleted" | "user.revoked" | "user.disabled" | "user.enabled" | "user.limited" | "user.expired" | "user.traffic_reset" | "user.expires_in_72_hours" | "user.expires_in_48_hours" | "user.expires_in_24_hours" | "user.expired_24_hours_ago" | "user.first_connected" | "user.bandwidth_usage_threshold_reached" | "user.not_connected";
    timestamp: string;
    meta: {
        notConnectedAfterHours?: number | null | undefined;
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
        }>;
        hwidUserDevice: z.ZodObject<{
            hwid: z.ZodString;
            userUuid: z.ZodString;
            platform: z.ZodNullable<z.ZodString>;
            osVersion: z.ZodNullable<z.ZodString>;
            deviceModel: z.ZodNullable<z.ZodString>;
            userAgent: z.ZodNullable<z.ZodString>;
            createdAt: z.ZodEffects<z.ZodString, Date, string>;
            updatedAt: z.ZodEffects<z.ZodString, Date, string>;
        }, "strip", z.ZodTypeAny, {
            hwid: string;
            createdAt: Date;
            updatedAt: Date;
            userUuid: string;
            platform: string | null;
            osVersion: string | null;
            deviceModel: string | null;
            userAgent: string | null;
        }, {
            hwid: string;
            createdAt: string;
            updatedAt: string;
            userUuid: string;
            platform: string | null;
            osVersion: string | null;
            deviceModel: string | null;
            userAgent: string | null;
        }>;
    }, "strip", z.ZodTypeAny, {
        user: {
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
        };
        hwidUserDevice: {
            hwid: string;
            createdAt: Date;
            updatedAt: Date;
            userUuid: string;
            platform: string | null;
            osVersion: string | null;
            deviceModel: string | null;
            userAgent: string | null;
        };
    }, {
        user: {
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
        };
        hwidUserDevice: {
            hwid: string;
            createdAt: string;
            updatedAt: string;
            userUuid: string;
            platform: string | null;
            osVersion: string | null;
            deviceModel: string | null;
            userAgent: string | null;
        };
    }>;
}, "strip", z.ZodTypeAny, {
    data: {
        user: {
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
        };
        hwidUserDevice: {
            hwid: string;
            createdAt: Date;
            updatedAt: Date;
            userUuid: string;
            platform: string | null;
            osVersion: string | null;
            deviceModel: string | null;
            userAgent: string | null;
        };
    };
    scope: "user_hwid_devices";
    event: "user_hwid_devices.added" | "user_hwid_devices.deleted";
    timestamp: Date;
}, {
    data: {
        user: {
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
        };
        hwidUserDevice: {
            hwid: string;
            createdAt: string;
            updatedAt: string;
            userUuid: string;
            platform: string | null;
            osVersion: string | null;
            deviceModel: string | null;
            userAgent: string | null;
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
        isConnected: z.ZodBoolean;
        isDisabled: z.ZodBoolean;
        isConnecting: z.ZodBoolean;
        lastStatusChange: z.ZodNullable<z.ZodEffects<z.ZodString, Date, string>>;
        lastStatusMessage: z.ZodNullable<z.ZodString>;
        xrayVersion: z.ZodNullable<z.ZodString>;
        nodeVersion: z.ZodNullable<z.ZodString>;
        xrayUptime: z.ZodString;
        isTrafficTrackingActive: z.ZodBoolean;
        trafficResetDay: z.ZodNullable<z.ZodNumber>;
        trafficLimitBytes: z.ZodNullable<z.ZodNumber>;
        trafficUsedBytes: z.ZodNullable<z.ZodNumber>;
        notifyPercent: z.ZodNullable<z.ZodNumber>;
        usersOnline: z.ZodNullable<z.ZodNumber>;
        viewPosition: z.ZodNumber;
        countryCode: z.ZodString;
        consumptionMultiplier: z.ZodNumber;
        tags: z.ZodArray<z.ZodString, "many">;
        cpuCount: z.ZodNullable<z.ZodNumber>;
        cpuModel: z.ZodNullable<z.ZodString>;
        totalRam: z.ZodNullable<z.ZodString>;
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
                type: string;
                uuid: string;
                profileUuid: string;
                tag: string;
                network: string | null;
                security: string | null;
                port: number | null;
                rawInbound?: unknown;
            }, {
                type: string;
                uuid: string;
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
                type: string;
                uuid: string;
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
                type: string;
                uuid: string;
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
            createdAt: Date;
            updatedAt: Date;
            name: string;
            faviconLink: string | null;
            loginUrl: string | null;
        }, {
            uuid: string;
            createdAt: string;
            updatedAt: string;
            name: string;
            faviconLink: string | null;
            loginUrl: string | null;
        }>>;
    }, "strip", z.ZodTypeAny, {
        tags: string[];
        uuid: string;
        createdAt: Date;
        updatedAt: Date;
        provider: {
            uuid: string;
            createdAt: Date;
            updatedAt: Date;
            name: string;
            faviconLink: string | null;
            loginUrl: string | null;
        } | null;
        countryCode: string;
        name: string;
        port: number | null;
        viewPosition: number;
        trafficLimitBytes: number | null;
        address: string;
        isDisabled: boolean;
        isConnected: boolean;
        isConnecting: boolean;
        lastStatusChange: Date | null;
        lastStatusMessage: string | null;
        xrayVersion: string | null;
        nodeVersion: string | null;
        xrayUptime: string;
        isTrafficTrackingActive: boolean;
        trafficResetDay: number | null;
        trafficUsedBytes: number | null;
        notifyPercent: number | null;
        usersOnline: number | null;
        consumptionMultiplier: number;
        cpuCount: number | null;
        cpuModel: string | null;
        totalRam: string | null;
        configProfile: {
            activeConfigProfileUuid: string | null;
            activeInbounds: {
                type: string;
                uuid: string;
                profileUuid: string;
                tag: string;
                network: string | null;
                security: string | null;
                port: number | null;
                rawInbound?: unknown;
            }[];
        };
        providerUuid: string | null;
    }, {
        tags: string[];
        uuid: string;
        createdAt: string;
        updatedAt: string;
        provider: {
            uuid: string;
            createdAt: string;
            updatedAt: string;
            name: string;
            faviconLink: string | null;
            loginUrl: string | null;
        } | null;
        countryCode: string;
        name: string;
        port: number | null;
        viewPosition: number;
        trafficLimitBytes: number | null;
        address: string;
        isDisabled: boolean;
        isConnected: boolean;
        isConnecting: boolean;
        lastStatusChange: string | null;
        lastStatusMessage: string | null;
        xrayVersion: string | null;
        nodeVersion: string | null;
        xrayUptime: string;
        isTrafficTrackingActive: boolean;
        trafficResetDay: number | null;
        trafficUsedBytes: number | null;
        notifyPercent: number | null;
        usersOnline: number | null;
        consumptionMultiplier: number;
        cpuCount: number | null;
        cpuModel: string | null;
        totalRam: string | null;
        configProfile: {
            activeConfigProfileUuid: string | null;
            activeInbounds: {
                type: string;
                uuid: string;
                profileUuid: string;
                tag: string;
                network: string | null;
                security: string | null;
                port: number | null;
                rawInbound?: unknown;
            }[];
        };
        providerUuid: string | null;
    }>;
}, "strip", z.ZodTypeAny, {
    data: {
        tags: string[];
        uuid: string;
        createdAt: Date;
        updatedAt: Date;
        provider: {
            uuid: string;
            createdAt: Date;
            updatedAt: Date;
            name: string;
            faviconLink: string | null;
            loginUrl: string | null;
        } | null;
        countryCode: string;
        name: string;
        port: number | null;
        viewPosition: number;
        trafficLimitBytes: number | null;
        address: string;
        isDisabled: boolean;
        isConnected: boolean;
        isConnecting: boolean;
        lastStatusChange: Date | null;
        lastStatusMessage: string | null;
        xrayVersion: string | null;
        nodeVersion: string | null;
        xrayUptime: string;
        isTrafficTrackingActive: boolean;
        trafficResetDay: number | null;
        trafficUsedBytes: number | null;
        notifyPercent: number | null;
        usersOnline: number | null;
        consumptionMultiplier: number;
        cpuCount: number | null;
        cpuModel: string | null;
        totalRam: string | null;
        configProfile: {
            activeConfigProfileUuid: string | null;
            activeInbounds: {
                type: string;
                uuid: string;
                profileUuid: string;
                tag: string;
                network: string | null;
                security: string | null;
                port: number | null;
                rawInbound?: unknown;
            }[];
        };
        providerUuid: string | null;
    };
    scope: "node";
    event: "node.created" | "node.modified" | "node.disabled" | "node.enabled" | "node.deleted" | "node.connection_lost" | "node.connection_restored" | "node.traffic_notify";
    timestamp: Date;
}, {
    data: {
        tags: string[];
        uuid: string;
        createdAt: string;
        updatedAt: string;
        provider: {
            uuid: string;
            createdAt: string;
            updatedAt: string;
            name: string;
            faviconLink: string | null;
            loginUrl: string | null;
        } | null;
        countryCode: string;
        name: string;
        port: number | null;
        viewPosition: number;
        trafficLimitBytes: number | null;
        address: string;
        isDisabled: boolean;
        isConnected: boolean;
        isConnecting: boolean;
        lastStatusChange: string | null;
        lastStatusMessage: string | null;
        xrayVersion: string | null;
        nodeVersion: string | null;
        xrayUptime: string;
        isTrafficTrackingActive: boolean;
        trafficResetDay: number | null;
        trafficUsedBytes: number | null;
        notifyPercent: number | null;
        usersOnline: number | null;
        consumptionMultiplier: number;
        cpuCount: number | null;
        cpuModel: string | null;
        totalRam: string | null;
        configProfile: {
            activeConfigProfileUuid: string | null;
            activeInbounds: {
                type: string;
                uuid: string;
                profileUuid: string;
                tag: string;
                network: string | null;
                security: string | null;
                port: number | null;
                rawInbound?: unknown;
            }[];
        };
        providerUuid: string | null;
    };
    scope: "node";
    event: "node.created" | "node.modified" | "node.disabled" | "node.enabled" | "node.deleted" | "node.connection_lost" | "node.connection_restored" | "node.traffic_notify";
    timestamp: string;
}>, z.ZodObject<{
    scope: z.ZodLiteral<"service">;
    event: z.ZodEnum<["service.panel_started" | "service.login_attempt_failed" | "service.login_attempt_success", ...("service.panel_started" | "service.login_attempt_failed" | "service.login_attempt_success")[]]>;
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
            password?: string | undefined;
            description?: string | undefined;
        }, {
            username: string;
            ip: string;
            userAgent: string;
            password?: string | undefined;
            description?: string | undefined;
        }>>;
        panelVersion: z.ZodOptional<z.ZodString>;
    }, "strip", z.ZodTypeAny, {
        loginAttempt?: {
            username: string;
            ip: string;
            userAgent: string;
            password?: string | undefined;
            description?: string | undefined;
        } | undefined;
        panelVersion?: string | undefined;
    }, {
        loginAttempt?: {
            username: string;
            ip: string;
            userAgent: string;
            password?: string | undefined;
            description?: string | undefined;
        } | undefined;
        panelVersion?: string | undefined;
    }>;
}, "strip", z.ZodTypeAny, {
    data: {
        loginAttempt?: {
            username: string;
            ip: string;
            userAgent: string;
            password?: string | undefined;
            description?: string | undefined;
        } | undefined;
        panelVersion?: string | undefined;
    };
    scope: "service";
    event: "service.panel_started" | "service.login_attempt_failed" | "service.login_attempt_success";
    timestamp: Date;
}, {
    data: {
        loginAttempt?: {
            username: string;
            ip: string;
            userAgent: string;
            password?: string | undefined;
            description?: string | undefined;
        } | undefined;
        panelVersion?: string | undefined;
    };
    scope: "service";
    event: "service.panel_started" | "service.login_attempt_failed" | "service.login_attempt_success";
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
}>]>;
export type TCerberusWebhookEvent = z.infer<typeof CerberusWebhookEventSchema>;
export type TCerberusWebhookUserEvent = z.infer<typeof CerberusWebhookUserEvents>;
export type TCerberusWebhookNodeEvent = z.infer<typeof CerberusWebhookNodeEvents>;
export type TCerberusWebhookServiceEvent = z.infer<typeof CerberusWebhookServiceEvents>;
export type TCerberusWebhookErrorsEvent = z.infer<typeof CerberusWebhookErrorsEvents>;
export type TCerberusWebhookCrmEvent = z.infer<typeof CerberusWebhookCrmEvents>;
export type TCerberusWebhookUserHwidDevicesEvent = z.infer<typeof CerberusWebhookUserHwidDevicesEvents>;
//# sourceMappingURL=webhook.schema.d.ts.map