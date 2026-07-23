import { z } from 'zod';
export declare namespace GetInfraBillingNodesCommand {
    const url: "/api/infra-billing/nodes";
    const TSQ_url: "/api/infra-billing/nodes";
    const endpointDetails: import("../../constants").EndpointDetails;
    const ResponseSchema: z.ZodObject<{
        response: z.ZodObject<{
            totalBillingNodes: z.ZodNumber;
            billingNodes: z.ZodArray<z.ZodObject<{
                uuid: z.ZodString;
                nodeUuid: z.ZodNullable<z.ZodString>;
                name: z.ZodNullable<z.ZodString>;
                providerUuid: z.ZodString;
                provider: z.ZodObject<Pick<{
                    uuid: z.ZodString;
                    name: z.ZodString;
                    faviconLink: z.ZodNullable<z.ZodString>;
                    loginUrl: z.ZodNullable<z.ZodString>;
                    createdAt: z.ZodEffects<z.ZodString, Date, string>;
                    updatedAt: z.ZodEffects<z.ZodString, Date, string>;
                }, "uuid" | "name" | "faviconLink" | "loginUrl">, "strip", z.ZodTypeAny, {
                    uuid: string;
                    name: string;
                    faviconLink: string | null;
                    loginUrl: string | null;
                }, {
                    uuid: string;
                    name: string;
                    faviconLink: string | null;
                    loginUrl: string | null;
                }>;
                node: z.ZodNullable<z.ZodObject<Pick<{
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
                }>>;
                nextBillingAt: z.ZodEffects<z.ZodString, Date, string>;
                createdAt: z.ZodEffects<z.ZodString, Date, string>;
                updatedAt: z.ZodEffects<z.ZodString, Date, string>;
            }, "strip", z.ZodTypeAny, {
                node: {
                    uuid: string;
                    name: string;
                    countryCode: string;
                } | null;
                uuid: string;
                name: string | null;
                createdAt: Date;
                updatedAt: Date;
                provider: {
                    uuid: string;
                    name: string;
                    faviconLink: string | null;
                    loginUrl: string | null;
                };
                nodeUuid: string | null;
                providerUuid: string;
                nextBillingAt: Date;
            }, {
                node: {
                    uuid: string;
                    name: string;
                    countryCode: string;
                } | null;
                uuid: string;
                name: string | null;
                createdAt: string;
                updatedAt: string;
                provider: {
                    uuid: string;
                    name: string;
                    faviconLink: string | null;
                    loginUrl: string | null;
                };
                nodeUuid: string | null;
                providerUuid: string;
                nextBillingAt: string;
            }>, "many">;
            availableBillingNodes: z.ZodArray<z.ZodObject<Pick<{
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
            }>, "many">;
            totalAvailableBillingNodes: z.ZodNumber;
            stats: z.ZodObject<{
                upcomingNodesCount: z.ZodNumber;
                currentMonthPayments: z.ZodNumber;
                totalSpent: z.ZodNumber;
            }, "strip", z.ZodTypeAny, {
                upcomingNodesCount: number;
                currentMonthPayments: number;
                totalSpent: number;
            }, {
                upcomingNodesCount: number;
                currentMonthPayments: number;
                totalSpent: number;
            }>;
        }, "strip", z.ZodTypeAny, {
            stats: {
                upcomingNodesCount: number;
                currentMonthPayments: number;
                totalSpent: number;
            };
            billingNodes: {
                node: {
                    uuid: string;
                    name: string;
                    countryCode: string;
                } | null;
                uuid: string;
                name: string | null;
                createdAt: Date;
                updatedAt: Date;
                provider: {
                    uuid: string;
                    name: string;
                    faviconLink: string | null;
                    loginUrl: string | null;
                };
                nodeUuid: string | null;
                providerUuid: string;
                nextBillingAt: Date;
            }[];
            totalBillingNodes: number;
            availableBillingNodes: {
                uuid: string;
                name: string;
                countryCode: string;
            }[];
            totalAvailableBillingNodes: number;
        }, {
            stats: {
                upcomingNodesCount: number;
                currentMonthPayments: number;
                totalSpent: number;
            };
            billingNodes: {
                node: {
                    uuid: string;
                    name: string;
                    countryCode: string;
                } | null;
                uuid: string;
                name: string | null;
                createdAt: string;
                updatedAt: string;
                provider: {
                    uuid: string;
                    name: string;
                    faviconLink: string | null;
                    loginUrl: string | null;
                };
                nodeUuid: string | null;
                providerUuid: string;
                nextBillingAt: string;
            }[];
            totalBillingNodes: number;
            availableBillingNodes: {
                uuid: string;
                name: string;
                countryCode: string;
            }[];
            totalAvailableBillingNodes: number;
        }>;
    }, "strip", z.ZodTypeAny, {
        response: {
            stats: {
                upcomingNodesCount: number;
                currentMonthPayments: number;
                totalSpent: number;
            };
            billingNodes: {
                node: {
                    uuid: string;
                    name: string;
                    countryCode: string;
                } | null;
                uuid: string;
                name: string | null;
                createdAt: Date;
                updatedAt: Date;
                provider: {
                    uuid: string;
                    name: string;
                    faviconLink: string | null;
                    loginUrl: string | null;
                };
                nodeUuid: string | null;
                providerUuid: string;
                nextBillingAt: Date;
            }[];
            totalBillingNodes: number;
            availableBillingNodes: {
                uuid: string;
                name: string;
                countryCode: string;
            }[];
            totalAvailableBillingNodes: number;
        };
    }, {
        response: {
            stats: {
                upcomingNodesCount: number;
                currentMonthPayments: number;
                totalSpent: number;
            };
            billingNodes: {
                node: {
                    uuid: string;
                    name: string;
                    countryCode: string;
                } | null;
                uuid: string;
                name: string | null;
                createdAt: string;
                updatedAt: string;
                provider: {
                    uuid: string;
                    name: string;
                    faviconLink: string | null;
                    loginUrl: string | null;
                };
                nodeUuid: string | null;
                providerUuid: string;
                nextBillingAt: string;
            }[];
            totalBillingNodes: number;
            availableBillingNodes: {
                uuid: string;
                name: string;
                countryCode: string;
            }[];
            totalAvailableBillingNodes: number;
        };
    }>;
    type Response = z.infer<typeof ResponseSchema>;
}
//# sourceMappingURL=get-billing-nodes.command.d.ts.map