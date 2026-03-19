import { z } from 'zod';
export declare namespace CreateInfraBillingNodeCommand {
    const url: "/api/infra-billing/nodes";
    const TSQ_url: "/api/infra-billing/nodes";
    const endpointDetails: import("../../constants").EndpointDetails;
    const RequestSchema: z.ZodObject<{
        providerUuid: z.ZodString;
        nodeUuid: z.ZodString;
        nextBillingAt: z.ZodOptional<z.ZodEffects<z.ZodString, Date, string>>;
    }, "strip", z.ZodTypeAny, {
        nodeUuid: string;
        providerUuid: string;
        nextBillingAt?: Date | undefined;
    }, {
        nodeUuid: string;
        providerUuid: string;
        nextBillingAt?: string | undefined;
    }>;
    type Request = z.infer<typeof RequestSchema>;
    const ResponseSchema: z.ZodObject<{
        response: z.ZodObject<{
            totalBillingNodes: z.ZodNumber;
            billingNodes: z.ZodArray<z.ZodObject<{
                uuid: z.ZodString;
                nodeUuid: z.ZodString;
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
                node: z.ZodObject<Pick<{
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
                }, "uuid" | "countryCode" | "name">, "strip", z.ZodTypeAny, {
                    uuid: string;
                    countryCode: string;
                    name: string;
                }, {
                    uuid: string;
                    countryCode: string;
                    name: string;
                }>;
                nextBillingAt: z.ZodEffects<z.ZodString, Date, string>;
                createdAt: z.ZodEffects<z.ZodString, Date, string>;
                updatedAt: z.ZodEffects<z.ZodString, Date, string>;
            }, "strip", z.ZodTypeAny, {
                node: {
                    uuid: string;
                    countryCode: string;
                    name: string;
                };
                uuid: string;
                createdAt: Date;
                updatedAt: Date;
                provider: {
                    uuid: string;
                    name: string;
                    faviconLink: string | null;
                    loginUrl: string | null;
                };
                nodeUuid: string;
                providerUuid: string;
                nextBillingAt: Date;
            }, {
                node: {
                    uuid: string;
                    countryCode: string;
                    name: string;
                };
                uuid: string;
                createdAt: string;
                updatedAt: string;
                provider: {
                    uuid: string;
                    name: string;
                    faviconLink: string | null;
                    loginUrl: string | null;
                };
                nodeUuid: string;
                providerUuid: string;
                nextBillingAt: string;
            }>, "many">;
            availableBillingNodes: z.ZodArray<z.ZodObject<Pick<{
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
            }, "uuid" | "countryCode" | "name">, "strip", z.ZodTypeAny, {
                uuid: string;
                countryCode: string;
                name: string;
            }, {
                uuid: string;
                countryCode: string;
                name: string;
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
                    countryCode: string;
                    name: string;
                };
                uuid: string;
                createdAt: Date;
                updatedAt: Date;
                provider: {
                    uuid: string;
                    name: string;
                    faviconLink: string | null;
                    loginUrl: string | null;
                };
                nodeUuid: string;
                providerUuid: string;
                nextBillingAt: Date;
            }[];
            totalBillingNodes: number;
            availableBillingNodes: {
                uuid: string;
                countryCode: string;
                name: string;
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
                    countryCode: string;
                    name: string;
                };
                uuid: string;
                createdAt: string;
                updatedAt: string;
                provider: {
                    uuid: string;
                    name: string;
                    faviconLink: string | null;
                    loginUrl: string | null;
                };
                nodeUuid: string;
                providerUuid: string;
                nextBillingAt: string;
            }[];
            totalBillingNodes: number;
            availableBillingNodes: {
                uuid: string;
                countryCode: string;
                name: string;
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
                    countryCode: string;
                    name: string;
                };
                uuid: string;
                createdAt: Date;
                updatedAt: Date;
                provider: {
                    uuid: string;
                    name: string;
                    faviconLink: string | null;
                    loginUrl: string | null;
                };
                nodeUuid: string;
                providerUuid: string;
                nextBillingAt: Date;
            }[];
            totalBillingNodes: number;
            availableBillingNodes: {
                uuid: string;
                countryCode: string;
                name: string;
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
                    countryCode: string;
                    name: string;
                };
                uuid: string;
                createdAt: string;
                updatedAt: string;
                provider: {
                    uuid: string;
                    name: string;
                    faviconLink: string | null;
                    loginUrl: string | null;
                };
                nodeUuid: string;
                providerUuid: string;
                nextBillingAt: string;
            }[];
            totalBillingNodes: number;
            availableBillingNodes: {
                uuid: string;
                countryCode: string;
                name: string;
            }[];
            totalAvailableBillingNodes: number;
        };
    }>;
    type Response = z.infer<typeof ResponseSchema>;
}
//# sourceMappingURL=create-billing-node.command.d.ts.map