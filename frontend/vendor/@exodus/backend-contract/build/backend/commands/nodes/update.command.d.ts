import { z } from 'zod';
export declare namespace UpdateNodeCommand {
    const url: "/api/nodes/";
    const TSQ_url: "/api/nodes/";
    const endpointDetails: import("../../constants").EndpointDetails;
    const RequestSchema: z.ZodObject<Pick<{
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
            xray: z.ZodString;
            node: z.ZodString;
        }, "strip", z.ZodTypeAny, {
            node: string;
            xray: string;
        }, {
            node: string;
            xray: string;
        }>>;
        xrayUptime: z.ZodNumber;
        usersOnline: z.ZodNumber;
        note: z.ZodNullable<z.ZodString>;
    }, "uuid"> & {
        name: z.ZodOptional<z.ZodString>;
        address: z.ZodOptional<z.ZodString>;
        port: z.ZodOptional<z.ZodNumber>;
        proxyUrl: z.ZodOptional<z.ZodNullable<z.ZodString>>;
        isTrafficTrackingActive: z.ZodOptional<z.ZodBoolean>;
        trafficLimitBytes: z.ZodOptional<z.ZodNumber>;
        notifyPercent: z.ZodOptional<z.ZodNumber>;
        trafficResetDay: z.ZodOptional<z.ZodNumber>;
        countryCode: z.ZodOptional<z.ZodString>;
        consumptionMultiplier: z.ZodOptional<z.ZodEffects<z.ZodNumber, number, number>>;
        nodeConsumptionMultiplier: z.ZodOptional<z.ZodEffects<z.ZodNumber, number, number>>;
        configProfile: z.ZodOptional<z.ZodObject<{
            activeConfigProfileUuid: z.ZodString;
            activeInbounds: z.ZodArray<z.ZodString, "many">;
        }, "strip", z.ZodTypeAny, {
            activeConfigProfileUuid: string;
            activeInbounds: string[];
        }, {
            activeConfigProfileUuid: string;
            activeInbounds: string[];
        }>>;
        providerUuid: z.ZodOptional<z.ZodNullable<z.ZodString>>;
        tags: z.ZodOptional<z.ZodArray<z.ZodString, "many">>;
        activePluginUuid: z.ZodOptional<z.ZodNullable<z.ZodString>>;
        note: z.ZodOptional<z.ZodNullable<z.ZodString>>;
    }, "strip", z.ZodTypeAny, {
        uuid: string;
        tags?: string[] | undefined;
        name?: string | undefined;
        countryCode?: string | undefined;
        port?: number | undefined;
        trafficLimitBytes?: number | undefined;
        address?: string | undefined;
        proxyUrl?: string | null | undefined;
        isTrafficTrackingActive?: boolean | undefined;
        trafficResetDay?: number | undefined;
        notifyPercent?: number | undefined;
        consumptionMultiplier?: number | undefined;
        nodeConsumptionMultiplier?: number | undefined;
        configProfile?: {
            activeConfigProfileUuid: string;
            activeInbounds: string[];
        } | undefined;
        providerUuid?: string | null | undefined;
        activePluginUuid?: string | null | undefined;
        note?: string | null | undefined;
    }, {
        uuid: string;
        tags?: string[] | undefined;
        name?: string | undefined;
        countryCode?: string | undefined;
        port?: number | undefined;
        trafficLimitBytes?: number | undefined;
        address?: string | undefined;
        proxyUrl?: string | null | undefined;
        isTrafficTrackingActive?: boolean | undefined;
        trafficResetDay?: number | undefined;
        notifyPercent?: number | undefined;
        consumptionMultiplier?: number | undefined;
        nodeConsumptionMultiplier?: number | undefined;
        configProfile?: {
            activeConfigProfileUuid: string;
            activeInbounds: string[];
        } | undefined;
        providerUuid?: string | null | undefined;
        activePluginUuid?: string | null | undefined;
        note?: string | null | undefined;
    }>;
    type Request = z.infer<typeof RequestSchema>;
    const ResponseSchema: z.ZodObject<{
        response: z.ZodObject<{
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
                xray: z.ZodString;
                node: z.ZodString;
            }, "strip", z.ZodTypeAny, {
                node: string;
                xray: string;
            }, {
                node: string;
                xray: string;
            }>>;
            xrayUptime: z.ZodNumber;
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
                xray: string;
            } | null;
            xrayUptime: number;
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
                xray: string;
            } | null;
            xrayUptime: number;
            usersOnline: number;
            note: string | null;
        }>;
    }, "strip", z.ZodTypeAny, {
        response: {
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
                xray: string;
            } | null;
            xrayUptime: number;
            usersOnline: number;
            note: string | null;
        };
    }, {
        response: {
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
                xray: string;
            } | null;
            xrayUptime: number;
            usersOnline: number;
            note: string | null;
        };
    }>;
    type Response = z.infer<typeof ResponseSchema>;
}
//# sourceMappingURL=update.command.d.ts.map