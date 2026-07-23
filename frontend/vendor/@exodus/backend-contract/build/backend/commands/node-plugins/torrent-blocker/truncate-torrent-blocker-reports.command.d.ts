import { z } from 'zod';
export declare namespace TruncateTorrentBlockerReportsCommand {
    const url: "/api/node-plugins/torrent-blocker/truncate";
    const TSQ_url: "/api/node-plugins/torrent-blocker/truncate";
    const endpointDetails: import("../../../constants").EndpointDetails;
    const ResponseSchema: z.ZodObject<{
        response: z.ZodObject<{
            records: z.ZodArray<z.ZodObject<{
                id: z.ZodNumber;
                userId: z.ZodNumber;
                nodeId: z.ZodNumber;
                user: z.ZodObject<Pick<{
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
                }, "uuid" | "username">, "strip", z.ZodTypeAny, {
                    uuid: string;
                    username: string;
                }, {
                    uuid: string;
                    username: string;
                }>;
                node: z.ZodObject<Pick<{
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
                }, "uuid" | "name" | "countryCode">, "strip", z.ZodTypeAny, {
                    uuid: string;
                    name: string;
                    countryCode: string;
                }, {
                    uuid: string;
                    name: string;
                    countryCode: string;
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
                createdAt: z.ZodEffects<z.ZodString, Date, string>;
            }, "strip", z.ZodTypeAny, {
                user: {
                    uuid: string;
                    username: string;
                };
                node: {
                    uuid: string;
                    name: string;
                    countryCode: string;
                };
                createdAt: Date;
                id: number;
                userId: number;
                nodeId: number;
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
                    username: string;
                };
                node: {
                    uuid: string;
                    name: string;
                    countryCode: string;
                };
                createdAt: string;
                id: number;
                userId: number;
                nodeId: number;
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
            }>, "many">;
            total: z.ZodNumber;
        }, "strip", z.ZodTypeAny, {
            total: number;
            records: {
                user: {
                    uuid: string;
                    username: string;
                };
                node: {
                    uuid: string;
                    name: string;
                    countryCode: string;
                };
                createdAt: Date;
                id: number;
                userId: number;
                nodeId: number;
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
            }[];
        }, {
            total: number;
            records: {
                user: {
                    uuid: string;
                    username: string;
                };
                node: {
                    uuid: string;
                    name: string;
                    countryCode: string;
                };
                createdAt: string;
                id: number;
                userId: number;
                nodeId: number;
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
            }[];
        }>;
    }, "strip", z.ZodTypeAny, {
        response: {
            total: number;
            records: {
                user: {
                    uuid: string;
                    username: string;
                };
                node: {
                    uuid: string;
                    name: string;
                    countryCode: string;
                };
                createdAt: Date;
                id: number;
                userId: number;
                nodeId: number;
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
            }[];
        };
    }, {
        response: {
            total: number;
            records: {
                user: {
                    uuid: string;
                    username: string;
                };
                node: {
                    uuid: string;
                    name: string;
                    countryCode: string;
                };
                createdAt: string;
                id: number;
                userId: number;
                nodeId: number;
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
            }[];
        };
    }>;
    type Response = z.infer<typeof ResponseSchema>;
}
//# sourceMappingURL=truncate-torrent-blocker-reports.command.d.ts.map