import z from 'zod';
export declare const ExodusWebhookUserEvents: z.ZodObject<{
    scope: z.ZodLiteral<"user">;
    event: z.ZodEnum<{
        "user.created": "user.created";
        "user.modified": "user.modified";
        "user.deleted": "user.deleted";
        "user.revoked": "user.revoked";
        "user.disabled": "user.disabled";
        "user.enabled": "user.enabled";
        "user.limited": "user.limited";
        "user.expired": "user.expired";
        "user.traffic_reset": "user.traffic_reset";
        "user.first_connected": "user.first_connected";
        "user.bandwidth_usage_threshold_reached": "user.bandwidth_usage_threshold_reached";
        "user.not_connected": "user.not_connected";
        "user.expiration": "user.expiration";
    }>;
    timestamp: z.ZodPipe<z.ZodString, z.ZodTransform<Date, string>>;
    data: z.ZodObject<{
        id: z.ZodNumber;
        shortUuid: z.ZodString;
        username: z.ZodString;
        status: z.ZodEnum<{
            readonly ACTIVE: "ACTIVE";
            readonly DISABLED: "DISABLED";
            readonly LIMITED: "LIMITED";
            readonly EXPIRED: "EXPIRED";
        }>;
        trafficLimitBytes: z.ZodNumber;
        trafficLimitStrategy: z.ZodEnum<{
            readonly NO_RESET: "NO_RESET";
            readonly DAY: "DAY";
            readonly WEEK: "WEEK";
            readonly MONTH: "MONTH";
            readonly MONTH_ROLLING: "MONTH_ROLLING";
        }>;
        expireAt: z.ZodPipe<z.ZodISODateTime, z.ZodTransform<Date, string>>;
        telegramId: z.ZodNullable<z.ZodNumber>;
        email: z.ZodNullable<z.ZodEmail>;
        description: z.ZodNullable<z.ZodString>;
        tag: z.ZodNullable<z.ZodString>;
        hwidDeviceLimit: z.ZodNullable<z.ZodInt>;
        externalSquadUuid: z.ZodNullable<z.ZodUUID>;
        trojanPassword: z.ZodString;
        vlessUuid: z.ZodUUID;
        ssPassword: z.ZodString;
        naivePassword: z.ZodString;
        shadowtlsPassword: z.ZodString;
        hysteria2Password: z.ZodString;
        anytlsPassword: z.ZodString;
        lastTriggeredThreshold: z.ZodInt;
        subRevokedAt: z.ZodNullable<z.ZodPipe<z.ZodISODateTime, z.ZodTransform<Date, string>>>;
        lastTrafficResetAt: z.ZodNullable<z.ZodPipe<z.ZodISODateTime, z.ZodTransform<Date, string>>>;
        createdAt: z.ZodPipe<z.ZodISODateTime, z.ZodTransform<Date, string>>;
        updatedAt: z.ZodPipe<z.ZodISODateTime, z.ZodTransform<Date, string>>;
        subscriptionUrl: z.ZodString;
        activeInternalSquads: z.ZodArray<z.ZodObject<{
            uuid: z.ZodUUID;
            name: z.ZodString;
        }, z.core.$strip>>;
        userTraffic: z.ZodObject<{
            usedTrafficBytes: z.ZodNumber;
            lifetimeUsedTrafficBytes: z.ZodNumber;
            onlineAt: z.ZodNullable<z.ZodPipe<z.ZodISODateTime, z.ZodTransform<Date, string>>>;
            firstConnectedAt: z.ZodNullable<z.ZodPipe<z.ZodISODateTime, z.ZodTransform<Date, string>>>;
            lastConnectedNodeUuid: z.ZodNullable<z.ZodUUID>;
        }, z.core.$strip>;
    }, z.core.$strip>;
    meta: z.ZodNullable<z.ZodObject<{
        notConnectedAfterHours: z.ZodOptional<z.ZodNullable<z.ZodNumber>>;
        expiration: z.ZodOptional<z.ZodNullable<z.ZodNumber>>;
    }, z.core.$strip>>;
}, z.core.$strip>;
export declare const ExodusWebhookUserHwidDevicesEvents: z.ZodObject<{
    scope: z.ZodLiteral<"user_hwid_devices">;
    event: z.ZodEnum<{
        "user_hwid_devices.added": "user_hwid_devices.added";
        "user_hwid_devices.deleted": "user_hwid_devices.deleted";
    }>;
    timestamp: z.ZodPipe<z.ZodString, z.ZodTransform<Date, string>>;
    data: z.ZodObject<{
        user: z.ZodObject<{
            id: z.ZodNumber;
            shortUuid: z.ZodString;
            username: z.ZodString;
            status: z.ZodEnum<{
                readonly ACTIVE: "ACTIVE";
                readonly DISABLED: "DISABLED";
                readonly LIMITED: "LIMITED";
                readonly EXPIRED: "EXPIRED";
            }>;
            trafficLimitBytes: z.ZodNumber;
            trafficLimitStrategy: z.ZodEnum<{
                readonly NO_RESET: "NO_RESET";
                readonly DAY: "DAY";
                readonly WEEK: "WEEK";
                readonly MONTH: "MONTH";
                readonly MONTH_ROLLING: "MONTH_ROLLING";
            }>;
            expireAt: z.ZodPipe<z.ZodISODateTime, z.ZodTransform<Date, string>>;
            telegramId: z.ZodNullable<z.ZodNumber>;
            email: z.ZodNullable<z.ZodEmail>;
            description: z.ZodNullable<z.ZodString>;
            tag: z.ZodNullable<z.ZodString>;
            hwidDeviceLimit: z.ZodNullable<z.ZodInt>;
            externalSquadUuid: z.ZodNullable<z.ZodUUID>;
            trojanPassword: z.ZodString;
            vlessUuid: z.ZodUUID;
            ssPassword: z.ZodString;
            naivePassword: z.ZodString;
            shadowtlsPassword: z.ZodString;
            hysteria2Password: z.ZodString;
            anytlsPassword: z.ZodString;
            lastTriggeredThreshold: z.ZodInt;
            subRevokedAt: z.ZodNullable<z.ZodPipe<z.ZodISODateTime, z.ZodTransform<Date, string>>>;
            lastTrafficResetAt: z.ZodNullable<z.ZodPipe<z.ZodISODateTime, z.ZodTransform<Date, string>>>;
            createdAt: z.ZodPipe<z.ZodISODateTime, z.ZodTransform<Date, string>>;
            updatedAt: z.ZodPipe<z.ZodISODateTime, z.ZodTransform<Date, string>>;
            subscriptionUrl: z.ZodString;
            activeInternalSquads: z.ZodArray<z.ZodObject<{
                uuid: z.ZodUUID;
                name: z.ZodString;
            }, z.core.$strip>>;
            userTraffic: z.ZodObject<{
                usedTrafficBytes: z.ZodNumber;
                lifetimeUsedTrafficBytes: z.ZodNumber;
                onlineAt: z.ZodNullable<z.ZodPipe<z.ZodISODateTime, z.ZodTransform<Date, string>>>;
                firstConnectedAt: z.ZodNullable<z.ZodPipe<z.ZodISODateTime, z.ZodTransform<Date, string>>>;
                lastConnectedNodeUuid: z.ZodNullable<z.ZodUUID>;
            }, z.core.$strip>;
        }, z.core.$strip>;
        hwidUserDevice: z.ZodObject<{
            hwid: z.ZodString;
            userId: z.ZodNumber;
            platform: z.ZodNullable<z.ZodString>;
            osVersion: z.ZodNullable<z.ZodString>;
            deviceModel: z.ZodNullable<z.ZodString>;
            userAgent: z.ZodNullable<z.ZodString>;
            requestIp: z.ZodNullable<z.ZodString>;
            createdAt: z.ZodPipe<z.ZodISODateTime, z.ZodTransform<Date, string>>;
            updatedAt: z.ZodPipe<z.ZodISODateTime, z.ZodTransform<Date, string>>;
        }, z.core.$strip>;
    }, z.core.$strip>;
}, z.core.$strip>;
export declare const ExodusWebhookNodeEvents: z.ZodObject<{
    scope: z.ZodLiteral<"node">;
    event: z.ZodEnum<{
        "node.created": "node.created";
        "node.modified": "node.modified";
        "node.disabled": "node.disabled";
        "node.enabled": "node.enabled";
        "node.deleted": "node.deleted";
        "node.connection_lost": "node.connection_lost";
        "node.connection_restored": "node.connection_restored";
        "node.traffic_notify": "node.traffic_notify";
    }>;
    timestamp: z.ZodPipe<z.ZodString, z.ZodTransform<Date, string>>;
    data: z.ZodObject<{
        uuid: z.ZodUUID;
        id: z.ZodNumber;
        name: z.ZodString;
        address: z.ZodString;
        port: z.ZodNullable<z.ZodInt>;
        proxyUrl: z.ZodNullable<z.ZodString>;
        apiSchema: z.ZodDefault<z.ZodEnum<{
            mtls: "mtls";
            tls: "tls";
        }>>;
        apiPath: z.ZodDefault<z.ZodNullable<z.ZodString>>;
        grpcAuthToken: z.ZodNullable<z.ZodString>;
        isConnected: z.ZodBoolean;
        isDisabled: z.ZodBoolean;
        isConnecting: z.ZodBoolean;
        lastStatusChange: z.ZodNullable<z.ZodPipe<z.ZodISODateTime, z.ZodTransform<Date, string>>>;
        lastStatusMessage: z.ZodNullable<z.ZodString>;
        isTrafficTrackingActive: z.ZodBoolean;
        trafficResetDay: z.ZodNullable<z.ZodInt>;
        trafficLimitBytes: z.ZodNullable<z.ZodNumber>;
        trafficUsedBytes: z.ZodNullable<z.ZodNumber>;
        notifyPercent: z.ZodNullable<z.ZodInt>;
        viewPosition: z.ZodInt;
        countryCode: z.ZodString;
        consumptionMultiplier: z.ZodNumber;
        nodeConsumptionMultiplier: z.ZodNumber;
        tags: z.ZodArray<z.ZodString>;
        createdAt: z.ZodPipe<z.ZodISODateTime, z.ZodTransform<Date, string>>;
        updatedAt: z.ZodPipe<z.ZodISODateTime, z.ZodTransform<Date, string>>;
        configProfile: z.ZodObject<{
            activeConfigProfileUuid: z.ZodNullable<z.ZodUUID>;
            activeInbounds: z.ZodArray<z.ZodObject<{
                uuid: z.ZodUUID;
                profileUuid: z.ZodUUID;
                tag: z.ZodString;
                type: z.ZodString;
                network: z.ZodNullable<z.ZodString>;
                security: z.ZodNullable<z.ZodString>;
                port: z.ZodNullable<z.ZodNumber>;
                rawInbound: z.ZodNullable<z.ZodUnknown>;
            }, z.core.$strip>>;
        }, z.core.$strip>;
        providerUuid: z.ZodNullable<z.ZodUUID>;
        provider: z.ZodNullable<z.ZodObject<{
            uuid: z.ZodUUID;
            name: z.ZodString;
            faviconLink: z.ZodNullable<z.ZodString>;
            loginUrl: z.ZodNullable<z.ZodString>;
            createdAt: z.ZodPipe<z.ZodISODateTime, z.ZodTransform<Date, string>>;
            updatedAt: z.ZodPipe<z.ZodISODateTime, z.ZodTransform<Date, string>>;
        }, z.core.$strip>>;
        activePluginUuid: z.ZodNullable<z.ZodUUID>;
        system: z.ZodNullable<z.ZodObject<{
            info: z.ZodObject<{
                arch: z.ZodString;
                cpus: z.ZodInt;
                cpuModel: z.ZodString;
                memoryTotal: z.ZodNumber;
                hostname: z.ZodString;
                platform: z.ZodString;
                release: z.ZodString;
                type: z.ZodString;
                version: z.ZodString;
                networkInterfaces: z.ZodArray<z.ZodString>;
            }, z.core.$strip>;
            stats: z.ZodObject<{
                memoryFree: z.ZodNumber;
                memoryUsed: z.ZodNumber;
                uptime: z.ZodNumber;
                loadAvg: z.ZodArray<z.ZodNumber>;
                interface: z.ZodNullable<z.ZodObject<{
                    interface: z.ZodString;
                    rxBytesPerSec: z.ZodNumber;
                    txBytesPerSec: z.ZodNumber;
                    rxTotal: z.ZodNumber;
                    txTotal: z.ZodNumber;
                }, z.core.$strip>>;
            }, z.core.$strip>;
        }, z.core.$strip>>;
        versions: z.ZodNullable<z.ZodObject<{
            singbox: z.ZodString;
            node: z.ZodString;
        }, z.core.$strip>>;
        singboxUptime: z.ZodNumber;
        usersOnline: z.ZodNumber;
        note: z.ZodNullable<z.ZodString>;
    }, z.core.$strip>;
}, z.core.$strip>;
export declare const ExodusWebhookServiceEvents: z.ZodObject<{
    scope: z.ZodLiteral<"service">;
    event: z.ZodEnum<{
        "service.panel_started": "service.panel_started";
        "service.login_attempt_failed": "service.login_attempt_failed";
        "service.login_attempt_success": "service.login_attempt_success";
        "service.subpage_config_changed": "service.subpage_config_changed";
        "service.api_token_created": "service.api_token_created";
        "service.api_token_deleted": "service.api_token_deleted";
    }>;
    timestamp: z.ZodPipe<z.ZodString, z.ZodTransform<Date, string>>;
    data: z.ZodObject<{
        loginAttempt: z.ZodOptional<z.ZodObject<{
            username: z.ZodString;
            ip: z.ZodString;
            userAgent: z.ZodString;
            description: z.ZodOptional<z.ZodString>;
            password: z.ZodOptional<z.ZodString>;
        }, z.core.$strip>>;
        panelVersion: z.ZodOptional<z.ZodString>;
        subpageConfig: z.ZodOptional<z.ZodObject<{
            action: z.ZodEnum<{
                CREATED: "CREATED";
                UPDATED: "UPDATED";
                DELETED: "DELETED";
            }>;
            uuid: z.ZodUUID;
        }, z.core.$strip>>;
        apiToken: z.ZodOptional<z.ZodObject<{
            name: z.ZodString;
            uuid: z.ZodUUID;
            expireAt: z.ZodPipe<z.ZodString, z.ZodTransform<Date, string>>;
            scopes: z.ZodArray<z.ZodString>;
        }, z.core.$strip>>;
    }, z.core.$strip>;
}, z.core.$strip>;
export declare const ExodusWebhookErrorsEvents: z.ZodObject<{
    scope: z.ZodLiteral<"errors">;
    event: z.ZodEnum<{
        "errors.bandwidth_usage_threshold_reached_max_notifications": "errors.bandwidth_usage_threshold_reached_max_notifications";
    }>;
    timestamp: z.ZodPipe<z.ZodString, z.ZodTransform<Date, string>>;
    data: z.ZodObject<{
        description: z.ZodString;
    }, z.core.$strip>;
}, z.core.$strip>;
export declare const ExodusWebhookCrmEvents: z.ZodObject<{
    scope: z.ZodLiteral<"crm">;
    event: z.ZodEnum<{
        "crm.infra_billing_node_payment_in_7_days": "crm.infra_billing_node_payment_in_7_days";
        "crm.infra_billing_node_payment_in_48hrs": "crm.infra_billing_node_payment_in_48hrs";
        "crm.infra_billing_node_payment_in_24hrs": "crm.infra_billing_node_payment_in_24hrs";
        "crm.infra_billing_node_payment_due_today": "crm.infra_billing_node_payment_due_today";
        "crm.infra_billing_node_payment_overdue_24hrs": "crm.infra_billing_node_payment_overdue_24hrs";
        "crm.infra_billing_node_payment_overdue_48hrs": "crm.infra_billing_node_payment_overdue_48hrs";
        "crm.infra_billing_node_payment_overdue_7_days": "crm.infra_billing_node_payment_overdue_7_days";
    }>;
    timestamp: z.ZodPipe<z.ZodString, z.ZodTransform<Date, string>>;
    data: z.ZodObject<{
        providerName: z.ZodString;
        nodeName: z.ZodString;
        nextBillingAt: z.ZodPipe<z.ZodString, z.ZodTransform<Date, string>>;
        loginUrl: z.ZodString;
    }, z.core.$strip>;
}, z.core.$strip>;
export declare const ExodusWebhookTorrentBlockerEvents: z.ZodObject<{
    scope: z.ZodLiteral<"torrent_blocker">;
    event: z.ZodEnum<{
        "torrent_blocker.report": "torrent_blocker.report";
    }>;
    timestamp: z.ZodPipe<z.ZodString, z.ZodTransform<Date, string>>;
    data: z.ZodObject<{
        node: z.ZodObject<{
            uuid: z.ZodUUID;
            id: z.ZodNumber;
            name: z.ZodString;
            address: z.ZodString;
            port: z.ZodNullable<z.ZodInt>;
            proxyUrl: z.ZodNullable<z.ZodString>;
            apiSchema: z.ZodDefault<z.ZodEnum<{
                mtls: "mtls";
                tls: "tls";
            }>>;
            apiPath: z.ZodDefault<z.ZodNullable<z.ZodString>>;
            grpcAuthToken: z.ZodNullable<z.ZodString>;
            isConnected: z.ZodBoolean;
            isDisabled: z.ZodBoolean;
            isConnecting: z.ZodBoolean;
            lastStatusChange: z.ZodNullable<z.ZodPipe<z.ZodISODateTime, z.ZodTransform<Date, string>>>;
            lastStatusMessage: z.ZodNullable<z.ZodString>;
            isTrafficTrackingActive: z.ZodBoolean;
            trafficResetDay: z.ZodNullable<z.ZodInt>;
            trafficLimitBytes: z.ZodNullable<z.ZodNumber>;
            trafficUsedBytes: z.ZodNullable<z.ZodNumber>;
            notifyPercent: z.ZodNullable<z.ZodInt>;
            viewPosition: z.ZodInt;
            countryCode: z.ZodString;
            consumptionMultiplier: z.ZodNumber;
            nodeConsumptionMultiplier: z.ZodNumber;
            tags: z.ZodArray<z.ZodString>;
            createdAt: z.ZodPipe<z.ZodISODateTime, z.ZodTransform<Date, string>>;
            updatedAt: z.ZodPipe<z.ZodISODateTime, z.ZodTransform<Date, string>>;
            configProfile: z.ZodObject<{
                activeConfigProfileUuid: z.ZodNullable<z.ZodUUID>;
                activeInbounds: z.ZodArray<z.ZodObject<{
                    uuid: z.ZodUUID;
                    profileUuid: z.ZodUUID;
                    tag: z.ZodString;
                    type: z.ZodString;
                    network: z.ZodNullable<z.ZodString>;
                    security: z.ZodNullable<z.ZodString>;
                    port: z.ZodNullable<z.ZodNumber>;
                    rawInbound: z.ZodNullable<z.ZodUnknown>;
                }, z.core.$strip>>;
            }, z.core.$strip>;
            providerUuid: z.ZodNullable<z.ZodUUID>;
            provider: z.ZodNullable<z.ZodObject<{
                uuid: z.ZodUUID;
                name: z.ZodString;
                faviconLink: z.ZodNullable<z.ZodString>;
                loginUrl: z.ZodNullable<z.ZodString>;
                createdAt: z.ZodPipe<z.ZodISODateTime, z.ZodTransform<Date, string>>;
                updatedAt: z.ZodPipe<z.ZodISODateTime, z.ZodTransform<Date, string>>;
            }, z.core.$strip>>;
            activePluginUuid: z.ZodNullable<z.ZodUUID>;
            system: z.ZodNullable<z.ZodObject<{
                info: z.ZodObject<{
                    arch: z.ZodString;
                    cpus: z.ZodInt;
                    cpuModel: z.ZodString;
                    memoryTotal: z.ZodNumber;
                    hostname: z.ZodString;
                    platform: z.ZodString;
                    release: z.ZodString;
                    type: z.ZodString;
                    version: z.ZodString;
                    networkInterfaces: z.ZodArray<z.ZodString>;
                }, z.core.$strip>;
                stats: z.ZodObject<{
                    memoryFree: z.ZodNumber;
                    memoryUsed: z.ZodNumber;
                    uptime: z.ZodNumber;
                    loadAvg: z.ZodArray<z.ZodNumber>;
                    interface: z.ZodNullable<z.ZodObject<{
                        interface: z.ZodString;
                        rxBytesPerSec: z.ZodNumber;
                        txBytesPerSec: z.ZodNumber;
                        rxTotal: z.ZodNumber;
                        txTotal: z.ZodNumber;
                    }, z.core.$strip>>;
                }, z.core.$strip>;
            }, z.core.$strip>>;
            versions: z.ZodNullable<z.ZodObject<{
                singbox: z.ZodString;
                node: z.ZodString;
            }, z.core.$strip>>;
            singboxUptime: z.ZodNumber;
            usersOnline: z.ZodNumber;
            note: z.ZodNullable<z.ZodString>;
        }, z.core.$strip>;
        user: z.ZodObject<{
            id: z.ZodNumber;
            shortUuid: z.ZodString;
            username: z.ZodString;
            status: z.ZodEnum<{
                readonly ACTIVE: "ACTIVE";
                readonly DISABLED: "DISABLED";
                readonly LIMITED: "LIMITED";
                readonly EXPIRED: "EXPIRED";
            }>;
            trafficLimitBytes: z.ZodNumber;
            trafficLimitStrategy: z.ZodEnum<{
                readonly NO_RESET: "NO_RESET";
                readonly DAY: "DAY";
                readonly WEEK: "WEEK";
                readonly MONTH: "MONTH";
                readonly MONTH_ROLLING: "MONTH_ROLLING";
            }>;
            expireAt: z.ZodPipe<z.ZodISODateTime, z.ZodTransform<Date, string>>;
            telegramId: z.ZodNullable<z.ZodNumber>;
            email: z.ZodNullable<z.ZodEmail>;
            description: z.ZodNullable<z.ZodString>;
            tag: z.ZodNullable<z.ZodString>;
            hwidDeviceLimit: z.ZodNullable<z.ZodInt>;
            externalSquadUuid: z.ZodNullable<z.ZodUUID>;
            trojanPassword: z.ZodString;
            vlessUuid: z.ZodUUID;
            ssPassword: z.ZodString;
            naivePassword: z.ZodString;
            shadowtlsPassword: z.ZodString;
            hysteria2Password: z.ZodString;
            anytlsPassword: z.ZodString;
            lastTriggeredThreshold: z.ZodInt;
            subRevokedAt: z.ZodNullable<z.ZodPipe<z.ZodISODateTime, z.ZodTransform<Date, string>>>;
            lastTrafficResetAt: z.ZodNullable<z.ZodPipe<z.ZodISODateTime, z.ZodTransform<Date, string>>>;
            createdAt: z.ZodPipe<z.ZodISODateTime, z.ZodTransform<Date, string>>;
            updatedAt: z.ZodPipe<z.ZodISODateTime, z.ZodTransform<Date, string>>;
            subscriptionUrl: z.ZodString;
            activeInternalSquads: z.ZodArray<z.ZodObject<{
                uuid: z.ZodUUID;
                name: z.ZodString;
            }, z.core.$strip>>;
            userTraffic: z.ZodObject<{
                usedTrafficBytes: z.ZodNumber;
                lifetimeUsedTrafficBytes: z.ZodNumber;
                onlineAt: z.ZodNullable<z.ZodPipe<z.ZodISODateTime, z.ZodTransform<Date, string>>>;
                firstConnectedAt: z.ZodNullable<z.ZodPipe<z.ZodISODateTime, z.ZodTransform<Date, string>>>;
                lastConnectedNodeUuid: z.ZodNullable<z.ZodUUID>;
            }, z.core.$strip>;
        }, z.core.$strip>;
        report: z.ZodObject<{
            actionReport: z.ZodObject<{
                blocked: z.ZodBoolean;
                ip: z.ZodString;
                blockDuration: z.ZodNumber;
                willUnblockAt: z.ZodPipe<z.ZodString, z.ZodTransform<Date, string>>;
                userId: z.ZodString;
                processedAt: z.ZodPipe<z.ZodString, z.ZodTransform<Date, string>>;
            }, z.core.$strip>;
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
            }, z.core.$strip>;
        }, z.core.$strip>;
    }, z.core.$strip>;
}, z.core.$strip>;
export declare const ExodusWebhookEventSchema: z.ZodDiscriminatedUnion<[z.ZodObject<{
    scope: z.ZodLiteral<"user">;
    event: z.ZodEnum<{
        "user.created": "user.created";
        "user.modified": "user.modified";
        "user.deleted": "user.deleted";
        "user.revoked": "user.revoked";
        "user.disabled": "user.disabled";
        "user.enabled": "user.enabled";
        "user.limited": "user.limited";
        "user.expired": "user.expired";
        "user.traffic_reset": "user.traffic_reset";
        "user.first_connected": "user.first_connected";
        "user.bandwidth_usage_threshold_reached": "user.bandwidth_usage_threshold_reached";
        "user.not_connected": "user.not_connected";
        "user.expiration": "user.expiration";
    }>;
    timestamp: z.ZodPipe<z.ZodString, z.ZodTransform<Date, string>>;
    data: z.ZodObject<{
        id: z.ZodNumber;
        shortUuid: z.ZodString;
        username: z.ZodString;
        status: z.ZodEnum<{
            readonly ACTIVE: "ACTIVE";
            readonly DISABLED: "DISABLED";
            readonly LIMITED: "LIMITED";
            readonly EXPIRED: "EXPIRED";
        }>;
        trafficLimitBytes: z.ZodNumber;
        trafficLimitStrategy: z.ZodEnum<{
            readonly NO_RESET: "NO_RESET";
            readonly DAY: "DAY";
            readonly WEEK: "WEEK";
            readonly MONTH: "MONTH";
            readonly MONTH_ROLLING: "MONTH_ROLLING";
        }>;
        expireAt: z.ZodPipe<z.ZodISODateTime, z.ZodTransform<Date, string>>;
        telegramId: z.ZodNullable<z.ZodNumber>;
        email: z.ZodNullable<z.ZodEmail>;
        description: z.ZodNullable<z.ZodString>;
        tag: z.ZodNullable<z.ZodString>;
        hwidDeviceLimit: z.ZodNullable<z.ZodInt>;
        externalSquadUuid: z.ZodNullable<z.ZodUUID>;
        trojanPassword: z.ZodString;
        vlessUuid: z.ZodUUID;
        ssPassword: z.ZodString;
        naivePassword: z.ZodString;
        shadowtlsPassword: z.ZodString;
        hysteria2Password: z.ZodString;
        anytlsPassword: z.ZodString;
        lastTriggeredThreshold: z.ZodInt;
        subRevokedAt: z.ZodNullable<z.ZodPipe<z.ZodISODateTime, z.ZodTransform<Date, string>>>;
        lastTrafficResetAt: z.ZodNullable<z.ZodPipe<z.ZodISODateTime, z.ZodTransform<Date, string>>>;
        createdAt: z.ZodPipe<z.ZodISODateTime, z.ZodTransform<Date, string>>;
        updatedAt: z.ZodPipe<z.ZodISODateTime, z.ZodTransform<Date, string>>;
        subscriptionUrl: z.ZodString;
        activeInternalSquads: z.ZodArray<z.ZodObject<{
            uuid: z.ZodUUID;
            name: z.ZodString;
        }, z.core.$strip>>;
        userTraffic: z.ZodObject<{
            usedTrafficBytes: z.ZodNumber;
            lifetimeUsedTrafficBytes: z.ZodNumber;
            onlineAt: z.ZodNullable<z.ZodPipe<z.ZodISODateTime, z.ZodTransform<Date, string>>>;
            firstConnectedAt: z.ZodNullable<z.ZodPipe<z.ZodISODateTime, z.ZodTransform<Date, string>>>;
            lastConnectedNodeUuid: z.ZodNullable<z.ZodUUID>;
        }, z.core.$strip>;
    }, z.core.$strip>;
    meta: z.ZodNullable<z.ZodObject<{
        notConnectedAfterHours: z.ZodOptional<z.ZodNullable<z.ZodNumber>>;
        expiration: z.ZodOptional<z.ZodNullable<z.ZodNumber>>;
    }, z.core.$strip>>;
}, z.core.$strip>, z.ZodObject<{
    scope: z.ZodLiteral<"user_hwid_devices">;
    event: z.ZodEnum<{
        "user_hwid_devices.added": "user_hwid_devices.added";
        "user_hwid_devices.deleted": "user_hwid_devices.deleted";
    }>;
    timestamp: z.ZodPipe<z.ZodString, z.ZodTransform<Date, string>>;
    data: z.ZodObject<{
        user: z.ZodObject<{
            id: z.ZodNumber;
            shortUuid: z.ZodString;
            username: z.ZodString;
            status: z.ZodEnum<{
                readonly ACTIVE: "ACTIVE";
                readonly DISABLED: "DISABLED";
                readonly LIMITED: "LIMITED";
                readonly EXPIRED: "EXPIRED";
            }>;
            trafficLimitBytes: z.ZodNumber;
            trafficLimitStrategy: z.ZodEnum<{
                readonly NO_RESET: "NO_RESET";
                readonly DAY: "DAY";
                readonly WEEK: "WEEK";
                readonly MONTH: "MONTH";
                readonly MONTH_ROLLING: "MONTH_ROLLING";
            }>;
            expireAt: z.ZodPipe<z.ZodISODateTime, z.ZodTransform<Date, string>>;
            telegramId: z.ZodNullable<z.ZodNumber>;
            email: z.ZodNullable<z.ZodEmail>;
            description: z.ZodNullable<z.ZodString>;
            tag: z.ZodNullable<z.ZodString>;
            hwidDeviceLimit: z.ZodNullable<z.ZodInt>;
            externalSquadUuid: z.ZodNullable<z.ZodUUID>;
            trojanPassword: z.ZodString;
            vlessUuid: z.ZodUUID;
            ssPassword: z.ZodString;
            naivePassword: z.ZodString;
            shadowtlsPassword: z.ZodString;
            hysteria2Password: z.ZodString;
            anytlsPassword: z.ZodString;
            lastTriggeredThreshold: z.ZodInt;
            subRevokedAt: z.ZodNullable<z.ZodPipe<z.ZodISODateTime, z.ZodTransform<Date, string>>>;
            lastTrafficResetAt: z.ZodNullable<z.ZodPipe<z.ZodISODateTime, z.ZodTransform<Date, string>>>;
            createdAt: z.ZodPipe<z.ZodISODateTime, z.ZodTransform<Date, string>>;
            updatedAt: z.ZodPipe<z.ZodISODateTime, z.ZodTransform<Date, string>>;
            subscriptionUrl: z.ZodString;
            activeInternalSquads: z.ZodArray<z.ZodObject<{
                uuid: z.ZodUUID;
                name: z.ZodString;
            }, z.core.$strip>>;
            userTraffic: z.ZodObject<{
                usedTrafficBytes: z.ZodNumber;
                lifetimeUsedTrafficBytes: z.ZodNumber;
                onlineAt: z.ZodNullable<z.ZodPipe<z.ZodISODateTime, z.ZodTransform<Date, string>>>;
                firstConnectedAt: z.ZodNullable<z.ZodPipe<z.ZodISODateTime, z.ZodTransform<Date, string>>>;
                lastConnectedNodeUuid: z.ZodNullable<z.ZodUUID>;
            }, z.core.$strip>;
        }, z.core.$strip>;
        hwidUserDevice: z.ZodObject<{
            hwid: z.ZodString;
            userId: z.ZodNumber;
            platform: z.ZodNullable<z.ZodString>;
            osVersion: z.ZodNullable<z.ZodString>;
            deviceModel: z.ZodNullable<z.ZodString>;
            userAgent: z.ZodNullable<z.ZodString>;
            requestIp: z.ZodNullable<z.ZodString>;
            createdAt: z.ZodPipe<z.ZodISODateTime, z.ZodTransform<Date, string>>;
            updatedAt: z.ZodPipe<z.ZodISODateTime, z.ZodTransform<Date, string>>;
        }, z.core.$strip>;
    }, z.core.$strip>;
}, z.core.$strip>, z.ZodObject<{
    scope: z.ZodLiteral<"node">;
    event: z.ZodEnum<{
        "node.created": "node.created";
        "node.modified": "node.modified";
        "node.disabled": "node.disabled";
        "node.enabled": "node.enabled";
        "node.deleted": "node.deleted";
        "node.connection_lost": "node.connection_lost";
        "node.connection_restored": "node.connection_restored";
        "node.traffic_notify": "node.traffic_notify";
    }>;
    timestamp: z.ZodPipe<z.ZodString, z.ZodTransform<Date, string>>;
    data: z.ZodObject<{
        uuid: z.ZodUUID;
        id: z.ZodNumber;
        name: z.ZodString;
        address: z.ZodString;
        port: z.ZodNullable<z.ZodInt>;
        proxyUrl: z.ZodNullable<z.ZodString>;
        apiSchema: z.ZodDefault<z.ZodEnum<{
            mtls: "mtls";
            tls: "tls";
        }>>;
        apiPath: z.ZodDefault<z.ZodNullable<z.ZodString>>;
        grpcAuthToken: z.ZodNullable<z.ZodString>;
        isConnected: z.ZodBoolean;
        isDisabled: z.ZodBoolean;
        isConnecting: z.ZodBoolean;
        lastStatusChange: z.ZodNullable<z.ZodPipe<z.ZodISODateTime, z.ZodTransform<Date, string>>>;
        lastStatusMessage: z.ZodNullable<z.ZodString>;
        isTrafficTrackingActive: z.ZodBoolean;
        trafficResetDay: z.ZodNullable<z.ZodInt>;
        trafficLimitBytes: z.ZodNullable<z.ZodNumber>;
        trafficUsedBytes: z.ZodNullable<z.ZodNumber>;
        notifyPercent: z.ZodNullable<z.ZodInt>;
        viewPosition: z.ZodInt;
        countryCode: z.ZodString;
        consumptionMultiplier: z.ZodNumber;
        nodeConsumptionMultiplier: z.ZodNumber;
        tags: z.ZodArray<z.ZodString>;
        createdAt: z.ZodPipe<z.ZodISODateTime, z.ZodTransform<Date, string>>;
        updatedAt: z.ZodPipe<z.ZodISODateTime, z.ZodTransform<Date, string>>;
        configProfile: z.ZodObject<{
            activeConfigProfileUuid: z.ZodNullable<z.ZodUUID>;
            activeInbounds: z.ZodArray<z.ZodObject<{
                uuid: z.ZodUUID;
                profileUuid: z.ZodUUID;
                tag: z.ZodString;
                type: z.ZodString;
                network: z.ZodNullable<z.ZodString>;
                security: z.ZodNullable<z.ZodString>;
                port: z.ZodNullable<z.ZodNumber>;
                rawInbound: z.ZodNullable<z.ZodUnknown>;
            }, z.core.$strip>>;
        }, z.core.$strip>;
        providerUuid: z.ZodNullable<z.ZodUUID>;
        provider: z.ZodNullable<z.ZodObject<{
            uuid: z.ZodUUID;
            name: z.ZodString;
            faviconLink: z.ZodNullable<z.ZodString>;
            loginUrl: z.ZodNullable<z.ZodString>;
            createdAt: z.ZodPipe<z.ZodISODateTime, z.ZodTransform<Date, string>>;
            updatedAt: z.ZodPipe<z.ZodISODateTime, z.ZodTransform<Date, string>>;
        }, z.core.$strip>>;
        activePluginUuid: z.ZodNullable<z.ZodUUID>;
        system: z.ZodNullable<z.ZodObject<{
            info: z.ZodObject<{
                arch: z.ZodString;
                cpus: z.ZodInt;
                cpuModel: z.ZodString;
                memoryTotal: z.ZodNumber;
                hostname: z.ZodString;
                platform: z.ZodString;
                release: z.ZodString;
                type: z.ZodString;
                version: z.ZodString;
                networkInterfaces: z.ZodArray<z.ZodString>;
            }, z.core.$strip>;
            stats: z.ZodObject<{
                memoryFree: z.ZodNumber;
                memoryUsed: z.ZodNumber;
                uptime: z.ZodNumber;
                loadAvg: z.ZodArray<z.ZodNumber>;
                interface: z.ZodNullable<z.ZodObject<{
                    interface: z.ZodString;
                    rxBytesPerSec: z.ZodNumber;
                    txBytesPerSec: z.ZodNumber;
                    rxTotal: z.ZodNumber;
                    txTotal: z.ZodNumber;
                }, z.core.$strip>>;
            }, z.core.$strip>;
        }, z.core.$strip>>;
        versions: z.ZodNullable<z.ZodObject<{
            singbox: z.ZodString;
            node: z.ZodString;
        }, z.core.$strip>>;
        singboxUptime: z.ZodNumber;
        usersOnline: z.ZodNumber;
        note: z.ZodNullable<z.ZodString>;
    }, z.core.$strip>;
}, z.core.$strip>, z.ZodObject<{
    scope: z.ZodLiteral<"service">;
    event: z.ZodEnum<{
        "service.panel_started": "service.panel_started";
        "service.login_attempt_failed": "service.login_attempt_failed";
        "service.login_attempt_success": "service.login_attempt_success";
        "service.subpage_config_changed": "service.subpage_config_changed";
        "service.api_token_created": "service.api_token_created";
        "service.api_token_deleted": "service.api_token_deleted";
    }>;
    timestamp: z.ZodPipe<z.ZodString, z.ZodTransform<Date, string>>;
    data: z.ZodObject<{
        loginAttempt: z.ZodOptional<z.ZodObject<{
            username: z.ZodString;
            ip: z.ZodString;
            userAgent: z.ZodString;
            description: z.ZodOptional<z.ZodString>;
            password: z.ZodOptional<z.ZodString>;
        }, z.core.$strip>>;
        panelVersion: z.ZodOptional<z.ZodString>;
        subpageConfig: z.ZodOptional<z.ZodObject<{
            action: z.ZodEnum<{
                CREATED: "CREATED";
                UPDATED: "UPDATED";
                DELETED: "DELETED";
            }>;
            uuid: z.ZodUUID;
        }, z.core.$strip>>;
        apiToken: z.ZodOptional<z.ZodObject<{
            name: z.ZodString;
            uuid: z.ZodUUID;
            expireAt: z.ZodPipe<z.ZodString, z.ZodTransform<Date, string>>;
            scopes: z.ZodArray<z.ZodString>;
        }, z.core.$strip>>;
    }, z.core.$strip>;
}, z.core.$strip>, z.ZodObject<{
    scope: z.ZodLiteral<"errors">;
    event: z.ZodEnum<{
        "errors.bandwidth_usage_threshold_reached_max_notifications": "errors.bandwidth_usage_threshold_reached_max_notifications";
    }>;
    timestamp: z.ZodPipe<z.ZodString, z.ZodTransform<Date, string>>;
    data: z.ZodObject<{
        description: z.ZodString;
    }, z.core.$strip>;
}, z.core.$strip>, z.ZodObject<{
    scope: z.ZodLiteral<"crm">;
    event: z.ZodEnum<{
        "crm.infra_billing_node_payment_in_7_days": "crm.infra_billing_node_payment_in_7_days";
        "crm.infra_billing_node_payment_in_48hrs": "crm.infra_billing_node_payment_in_48hrs";
        "crm.infra_billing_node_payment_in_24hrs": "crm.infra_billing_node_payment_in_24hrs";
        "crm.infra_billing_node_payment_due_today": "crm.infra_billing_node_payment_due_today";
        "crm.infra_billing_node_payment_overdue_24hrs": "crm.infra_billing_node_payment_overdue_24hrs";
        "crm.infra_billing_node_payment_overdue_48hrs": "crm.infra_billing_node_payment_overdue_48hrs";
        "crm.infra_billing_node_payment_overdue_7_days": "crm.infra_billing_node_payment_overdue_7_days";
    }>;
    timestamp: z.ZodPipe<z.ZodString, z.ZodTransform<Date, string>>;
    data: z.ZodObject<{
        providerName: z.ZodString;
        nodeName: z.ZodString;
        nextBillingAt: z.ZodPipe<z.ZodString, z.ZodTransform<Date, string>>;
        loginUrl: z.ZodString;
    }, z.core.$strip>;
}, z.core.$strip>, z.ZodObject<{
    scope: z.ZodLiteral<"torrent_blocker">;
    event: z.ZodEnum<{
        "torrent_blocker.report": "torrent_blocker.report";
    }>;
    timestamp: z.ZodPipe<z.ZodString, z.ZodTransform<Date, string>>;
    data: z.ZodObject<{
        node: z.ZodObject<{
            uuid: z.ZodUUID;
            id: z.ZodNumber;
            name: z.ZodString;
            address: z.ZodString;
            port: z.ZodNullable<z.ZodInt>;
            proxyUrl: z.ZodNullable<z.ZodString>;
            apiSchema: z.ZodDefault<z.ZodEnum<{
                mtls: "mtls";
                tls: "tls";
            }>>;
            apiPath: z.ZodDefault<z.ZodNullable<z.ZodString>>;
            grpcAuthToken: z.ZodNullable<z.ZodString>;
            isConnected: z.ZodBoolean;
            isDisabled: z.ZodBoolean;
            isConnecting: z.ZodBoolean;
            lastStatusChange: z.ZodNullable<z.ZodPipe<z.ZodISODateTime, z.ZodTransform<Date, string>>>;
            lastStatusMessage: z.ZodNullable<z.ZodString>;
            isTrafficTrackingActive: z.ZodBoolean;
            trafficResetDay: z.ZodNullable<z.ZodInt>;
            trafficLimitBytes: z.ZodNullable<z.ZodNumber>;
            trafficUsedBytes: z.ZodNullable<z.ZodNumber>;
            notifyPercent: z.ZodNullable<z.ZodInt>;
            viewPosition: z.ZodInt;
            countryCode: z.ZodString;
            consumptionMultiplier: z.ZodNumber;
            nodeConsumptionMultiplier: z.ZodNumber;
            tags: z.ZodArray<z.ZodString>;
            createdAt: z.ZodPipe<z.ZodISODateTime, z.ZodTransform<Date, string>>;
            updatedAt: z.ZodPipe<z.ZodISODateTime, z.ZodTransform<Date, string>>;
            configProfile: z.ZodObject<{
                activeConfigProfileUuid: z.ZodNullable<z.ZodUUID>;
                activeInbounds: z.ZodArray<z.ZodObject<{
                    uuid: z.ZodUUID;
                    profileUuid: z.ZodUUID;
                    tag: z.ZodString;
                    type: z.ZodString;
                    network: z.ZodNullable<z.ZodString>;
                    security: z.ZodNullable<z.ZodString>;
                    port: z.ZodNullable<z.ZodNumber>;
                    rawInbound: z.ZodNullable<z.ZodUnknown>;
                }, z.core.$strip>>;
            }, z.core.$strip>;
            providerUuid: z.ZodNullable<z.ZodUUID>;
            provider: z.ZodNullable<z.ZodObject<{
                uuid: z.ZodUUID;
                name: z.ZodString;
                faviconLink: z.ZodNullable<z.ZodString>;
                loginUrl: z.ZodNullable<z.ZodString>;
                createdAt: z.ZodPipe<z.ZodISODateTime, z.ZodTransform<Date, string>>;
                updatedAt: z.ZodPipe<z.ZodISODateTime, z.ZodTransform<Date, string>>;
            }, z.core.$strip>>;
            activePluginUuid: z.ZodNullable<z.ZodUUID>;
            system: z.ZodNullable<z.ZodObject<{
                info: z.ZodObject<{
                    arch: z.ZodString;
                    cpus: z.ZodInt;
                    cpuModel: z.ZodString;
                    memoryTotal: z.ZodNumber;
                    hostname: z.ZodString;
                    platform: z.ZodString;
                    release: z.ZodString;
                    type: z.ZodString;
                    version: z.ZodString;
                    networkInterfaces: z.ZodArray<z.ZodString>;
                }, z.core.$strip>;
                stats: z.ZodObject<{
                    memoryFree: z.ZodNumber;
                    memoryUsed: z.ZodNumber;
                    uptime: z.ZodNumber;
                    loadAvg: z.ZodArray<z.ZodNumber>;
                    interface: z.ZodNullable<z.ZodObject<{
                        interface: z.ZodString;
                        rxBytesPerSec: z.ZodNumber;
                        txBytesPerSec: z.ZodNumber;
                        rxTotal: z.ZodNumber;
                        txTotal: z.ZodNumber;
                    }, z.core.$strip>>;
                }, z.core.$strip>;
            }, z.core.$strip>>;
            versions: z.ZodNullable<z.ZodObject<{
                singbox: z.ZodString;
                node: z.ZodString;
            }, z.core.$strip>>;
            singboxUptime: z.ZodNumber;
            usersOnline: z.ZodNumber;
            note: z.ZodNullable<z.ZodString>;
        }, z.core.$strip>;
        user: z.ZodObject<{
            id: z.ZodNumber;
            shortUuid: z.ZodString;
            username: z.ZodString;
            status: z.ZodEnum<{
                readonly ACTIVE: "ACTIVE";
                readonly DISABLED: "DISABLED";
                readonly LIMITED: "LIMITED";
                readonly EXPIRED: "EXPIRED";
            }>;
            trafficLimitBytes: z.ZodNumber;
            trafficLimitStrategy: z.ZodEnum<{
                readonly NO_RESET: "NO_RESET";
                readonly DAY: "DAY";
                readonly WEEK: "WEEK";
                readonly MONTH: "MONTH";
                readonly MONTH_ROLLING: "MONTH_ROLLING";
            }>;
            expireAt: z.ZodPipe<z.ZodISODateTime, z.ZodTransform<Date, string>>;
            telegramId: z.ZodNullable<z.ZodNumber>;
            email: z.ZodNullable<z.ZodEmail>;
            description: z.ZodNullable<z.ZodString>;
            tag: z.ZodNullable<z.ZodString>;
            hwidDeviceLimit: z.ZodNullable<z.ZodInt>;
            externalSquadUuid: z.ZodNullable<z.ZodUUID>;
            trojanPassword: z.ZodString;
            vlessUuid: z.ZodUUID;
            ssPassword: z.ZodString;
            naivePassword: z.ZodString;
            shadowtlsPassword: z.ZodString;
            hysteria2Password: z.ZodString;
            anytlsPassword: z.ZodString;
            lastTriggeredThreshold: z.ZodInt;
            subRevokedAt: z.ZodNullable<z.ZodPipe<z.ZodISODateTime, z.ZodTransform<Date, string>>>;
            lastTrafficResetAt: z.ZodNullable<z.ZodPipe<z.ZodISODateTime, z.ZodTransform<Date, string>>>;
            createdAt: z.ZodPipe<z.ZodISODateTime, z.ZodTransform<Date, string>>;
            updatedAt: z.ZodPipe<z.ZodISODateTime, z.ZodTransform<Date, string>>;
            subscriptionUrl: z.ZodString;
            activeInternalSquads: z.ZodArray<z.ZodObject<{
                uuid: z.ZodUUID;
                name: z.ZodString;
            }, z.core.$strip>>;
            userTraffic: z.ZodObject<{
                usedTrafficBytes: z.ZodNumber;
                lifetimeUsedTrafficBytes: z.ZodNumber;
                onlineAt: z.ZodNullable<z.ZodPipe<z.ZodISODateTime, z.ZodTransform<Date, string>>>;
                firstConnectedAt: z.ZodNullable<z.ZodPipe<z.ZodISODateTime, z.ZodTransform<Date, string>>>;
                lastConnectedNodeUuid: z.ZodNullable<z.ZodUUID>;
            }, z.core.$strip>;
        }, z.core.$strip>;
        report: z.ZodObject<{
            actionReport: z.ZodObject<{
                blocked: z.ZodBoolean;
                ip: z.ZodString;
                blockDuration: z.ZodNumber;
                willUnblockAt: z.ZodPipe<z.ZodString, z.ZodTransform<Date, string>>;
                userId: z.ZodString;
                processedAt: z.ZodPipe<z.ZodString, z.ZodTransform<Date, string>>;
            }, z.core.$strip>;
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
            }, z.core.$strip>;
        }, z.core.$strip>;
    }, z.core.$strip>;
}, z.core.$strip>], "scope">;
export type TExodusWebhookEvent = z.infer<typeof ExodusWebhookEventSchema>;
export type TExodusWebhookUserEvent = z.infer<typeof ExodusWebhookUserEvents>;
export type TExodusWebhookNodeEvent = z.infer<typeof ExodusWebhookNodeEvents>;
export type TExodusWebhookServiceEvent = z.infer<typeof ExodusWebhookServiceEvents>;
export type TExodusWebhookErrorsEvent = z.infer<typeof ExodusWebhookErrorsEvents>;
export type TExodusWebhookCrmEvent = z.infer<typeof ExodusWebhookCrmEvents>;
export type TExodusWebhookUserHwidDevicesEvent = z.infer<typeof ExodusWebhookUserHwidDevicesEvents>;
export type TExodusWebhookTorrentBlockerEvent = z.infer<typeof ExodusWebhookTorrentBlockerEvents>;
//# sourceMappingURL=webhook.schema.d.ts.map