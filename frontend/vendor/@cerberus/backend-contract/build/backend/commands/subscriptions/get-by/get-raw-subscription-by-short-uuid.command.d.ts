import { z } from 'zod';
export declare namespace GetRawSubscriptionByShortUuidCommand {
    const url: (shortUuid: string) => string;
    const TSQ_url: string;
    const endpointDetails: import("../../../constants").EndpointDetails;
    const RequestSchema: z.ZodObject<{
        shortUuid: z.ZodString;
    }, "strip", z.ZodTypeAny, {
        shortUuid: string;
    }, {
        shortUuid: string;
    }>;
    type Request = z.infer<typeof RequestSchema>;
    const RequestQuerySchema: z.ZodObject<{
        withDisabledHosts: z.ZodDefault<z.ZodOptional<z.ZodEffects<z.ZodString, boolean, string>>>;
    }, "strip", z.ZodTypeAny, {
        withDisabledHosts: boolean;
    }, {
        withDisabledHosts?: string | undefined;
    }>;
    type RequestQuery = z.infer<typeof RequestQuerySchema>;
    const ResponseSchema: z.ZodObject<{
        response: z.ZodObject<{
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
            convertedUserInfo: z.ZodObject<{
                daysLeft: z.ZodNumber;
                trafficLimit: z.ZodString;
                trafficUsed: z.ZodString;
                lifetimeTrafficUsed: z.ZodString;
                isHwidLimited: z.ZodBoolean;
            }, "strip", z.ZodTypeAny, {
                daysLeft: number;
                trafficUsed: string;
                trafficLimit: string;
                lifetimeTrafficUsed: string;
                isHwidLimited: boolean;
            }, {
                daysLeft: number;
                trafficUsed: string;
                trafficLimit: string;
                lifetimeTrafficUsed: string;
                isHwidLimited: boolean;
            }>;
            headers: z.ZodRecord<z.ZodString, z.ZodOptional<z.ZodString>>;
            rawHosts: z.ZodArray<z.ZodObject<{
                address: z.ZodOptional<z.ZodNullable<z.ZodString>>;
                alpn: z.ZodOptional<z.ZodNullable<z.ZodString>>;
                fingerprint: z.ZodOptional<z.ZodNullable<z.ZodString>>;
                host: z.ZodOptional<z.ZodNullable<z.ZodString>>;
                network: z.ZodOptional<z.ZodNullable<z.ZodString>>;
                password: z.ZodObject<{
                    ssPassword: z.ZodString;
                    trojanPassword: z.ZodString;
                    vlessPassword: z.ZodString;
                }, "strip", z.ZodTypeAny, {
                    trojanPassword: string;
                    ssPassword: string;
                    vlessPassword: string;
                }, {
                    trojanPassword: string;
                    ssPassword: string;
                    vlessPassword: string;
                }>;
                path: z.ZodOptional<z.ZodNullable<z.ZodString>>;
                publicKey: z.ZodOptional<z.ZodNullable<z.ZodString>>;
                port: z.ZodOptional<z.ZodNullable<z.ZodNumber>>;
                protocol: z.ZodOptional<z.ZodNullable<z.ZodString>>;
                remark: z.ZodOptional<z.ZodNullable<z.ZodString>>;
                shortId: z.ZodOptional<z.ZodNullable<z.ZodString>>;
                sni: z.ZodOptional<z.ZodNullable<z.ZodString>>;
                spiderX: z.ZodOptional<z.ZodNullable<z.ZodString>>;
                tls: z.ZodOptional<z.ZodNullable<z.ZodString>>;
                rawSettings: z.ZodOptional<z.ZodNullable<z.ZodObject<{
                    headerType: z.ZodOptional<z.ZodNullable<z.ZodString>>;
                    request: z.ZodOptional<z.ZodNullable<z.ZodObject<{}, "strip", z.ZodTypeAny, {}, {}>>>;
                }, "strip", z.ZodTypeAny, {
                    headerType?: string | null | undefined;
                    request?: {} | null | undefined;
                }, {
                    headerType?: string | null | undefined;
                    request?: {} | null | undefined;
                }>>>;
                additionalParams: z.ZodOptional<z.ZodNullable<z.ZodObject<{
                    mode: z.ZodOptional<z.ZodNullable<z.ZodString>>;
                    heartbeatPeriod: z.ZodOptional<z.ZodNullable<z.ZodNumber>>;
                }, "strip", z.ZodTypeAny, {
                    mode?: string | null | undefined;
                    heartbeatPeriod?: number | null | undefined;
                }, {
                    mode?: string | null | undefined;
                    heartbeatPeriod?: number | null | undefined;
                }>>>;
                xHttpExtraParams: z.ZodOptional<z.ZodNullable<z.ZodObject<{}, "strip", z.ZodTypeAny, {}, {}>>>;
                muxParams: z.ZodOptional<z.ZodNullable<z.ZodObject<{}, "strip", z.ZodTypeAny, {}, {}>>>;
                sockoptParams: z.ZodOptional<z.ZodNullable<z.ZodObject<{}, "strip", z.ZodTypeAny, {}, {}>>>;
                serverDescription: z.ZodOptional<z.ZodNullable<z.ZodString>>;
                flow: z.ZodOptional<z.ZodNullable<z.ZodString>>;
                allowInsecure: z.ZodOptional<z.ZodNullable<z.ZodBoolean>>;
                shuffleHost: z.ZodOptional<z.ZodNullable<z.ZodBoolean>>;
                mihomoX25519: z.ZodOptional<z.ZodNullable<z.ZodBoolean>>;
                mldsa65Verify: z.ZodOptional<z.ZodNullable<z.ZodString>>;
                encryption: z.ZodOptional<z.ZodNullable<z.ZodString>>;
                protocolOptions: z.ZodOptional<z.ZodNullable<z.ZodObject<{
                    ss: z.ZodOptional<z.ZodNullable<z.ZodObject<{
                        method: z.ZodOptional<z.ZodNullable<z.ZodString>>;
                    }, "strip", z.ZodTypeAny, {
                        method?: string | null | undefined;
                    }, {
                        method?: string | null | undefined;
                    }>>>;
                }, "strip", z.ZodTypeAny, {
                    ss?: {
                        method?: string | null | undefined;
                    } | null | undefined;
                }, {
                    ss?: {
                        method?: string | null | undefined;
                    } | null | undefined;
                }>>>;
                dbData: z.ZodOptional<z.ZodObject<{
                    rawInbound: z.ZodNullable<z.ZodObject<{}, "strip", z.ZodTypeAny, {}, {}>>;
                    inboundTag: z.ZodString;
                    uuid: z.ZodString;
                    configProfileUuid: z.ZodNullable<z.ZodString>;
                    configProfileInboundUuid: z.ZodNullable<z.ZodString>;
                    isDisabled: z.ZodBoolean;
                    viewPosition: z.ZodNumber;
                    remark: z.ZodString;
                    isHidden: z.ZodBoolean;
                    tag: z.ZodNullable<z.ZodString>;
                    vlessRouteId: z.ZodNullable<z.ZodNumber>;
                }, "strip", z.ZodTypeAny, {
                    uuid: string;
                    tag: string | null;
                    rawInbound: {} | null;
                    viewPosition: number;
                    remark: string;
                    isDisabled: boolean;
                    configProfileUuid: string | null;
                    configProfileInboundUuid: string | null;
                    isHidden: boolean;
                    vlessRouteId: number | null;
                    inboundTag: string;
                }, {
                    uuid: string;
                    tag: string | null;
                    rawInbound: {} | null;
                    viewPosition: number;
                    remark: string;
                    isDisabled: boolean;
                    configProfileUuid: string | null;
                    configProfileInboundUuid: string | null;
                    isHidden: boolean;
                    vlessRouteId: number | null;
                    inboundTag: string;
                }>>;
                xrayJsonTemplate: z.ZodOptional<z.ZodNullable<z.ZodObject<{}, "strip", z.ZodTypeAny, {}, {}>>>;
            }, "strip", z.ZodTypeAny, {
                password: {
                    trojanPassword: string;
                    ssPassword: string;
                    vlessPassword: string;
                };
                path?: string | null | undefined;
                network?: string | null | undefined;
                port?: number | null | undefined;
                remark?: string | null | undefined;
                address?: string | null | undefined;
                sni?: string | null | undefined;
                host?: string | null | undefined;
                alpn?: string | null | undefined;
                fingerprint?: string | null | undefined;
                xHttpExtraParams?: {} | null | undefined;
                muxParams?: {} | null | undefined;
                sockoptParams?: {} | null | undefined;
                serverDescription?: string | null | undefined;
                allowInsecure?: boolean | null | undefined;
                shuffleHost?: boolean | null | undefined;
                mihomoX25519?: boolean | null | undefined;
                publicKey?: string | null | undefined;
                protocol?: string | null | undefined;
                shortId?: string | null | undefined;
                spiderX?: string | null | undefined;
                tls?: string | null | undefined;
                rawSettings?: {
                    headerType?: string | null | undefined;
                    request?: {} | null | undefined;
                } | null | undefined;
                additionalParams?: {
                    mode?: string | null | undefined;
                    heartbeatPeriod?: number | null | undefined;
                } | null | undefined;
                flow?: string | null | undefined;
                mldsa65Verify?: string | null | undefined;
                encryption?: string | null | undefined;
                protocolOptions?: {
                    ss?: {
                        method?: string | null | undefined;
                    } | null | undefined;
                } | null | undefined;
                dbData?: {
                    uuid: string;
                    tag: string | null;
                    rawInbound: {} | null;
                    viewPosition: number;
                    remark: string;
                    isDisabled: boolean;
                    configProfileUuid: string | null;
                    configProfileInboundUuid: string | null;
                    isHidden: boolean;
                    vlessRouteId: number | null;
                    inboundTag: string;
                } | undefined;
                xrayJsonTemplate?: {} | null | undefined;
            }, {
                password: {
                    trojanPassword: string;
                    ssPassword: string;
                    vlessPassword: string;
                };
                path?: string | null | undefined;
                network?: string | null | undefined;
                port?: number | null | undefined;
                remark?: string | null | undefined;
                address?: string | null | undefined;
                sni?: string | null | undefined;
                host?: string | null | undefined;
                alpn?: string | null | undefined;
                fingerprint?: string | null | undefined;
                xHttpExtraParams?: {} | null | undefined;
                muxParams?: {} | null | undefined;
                sockoptParams?: {} | null | undefined;
                serverDescription?: string | null | undefined;
                allowInsecure?: boolean | null | undefined;
                shuffleHost?: boolean | null | undefined;
                mihomoX25519?: boolean | null | undefined;
                publicKey?: string | null | undefined;
                protocol?: string | null | undefined;
                shortId?: string | null | undefined;
                spiderX?: string | null | undefined;
                tls?: string | null | undefined;
                rawSettings?: {
                    headerType?: string | null | undefined;
                    request?: {} | null | undefined;
                } | null | undefined;
                additionalParams?: {
                    mode?: string | null | undefined;
                    heartbeatPeriod?: number | null | undefined;
                } | null | undefined;
                flow?: string | null | undefined;
                mldsa65Verify?: string | null | undefined;
                encryption?: string | null | undefined;
                protocolOptions?: {
                    ss?: {
                        method?: string | null | undefined;
                    } | null | undefined;
                } | null | undefined;
                dbData?: {
                    uuid: string;
                    tag: string | null;
                    rawInbound: {} | null;
                    viewPosition: number;
                    remark: string;
                    isDisabled: boolean;
                    configProfileUuid: string | null;
                    configProfileInboundUuid: string | null;
                    isHidden: boolean;
                    vlessRouteId: number | null;
                    inboundTag: string;
                } | undefined;
                xrayJsonTemplate?: {} | null | undefined;
            }>, "many">;
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
            headers: Record<string, string | undefined>;
            convertedUserInfo: {
                daysLeft: number;
                trafficUsed: string;
                trafficLimit: string;
                lifetimeTrafficUsed: string;
                isHwidLimited: boolean;
            };
            rawHosts: {
                password: {
                    trojanPassword: string;
                    ssPassword: string;
                    vlessPassword: string;
                };
                path?: string | null | undefined;
                network?: string | null | undefined;
                port?: number | null | undefined;
                remark?: string | null | undefined;
                address?: string | null | undefined;
                sni?: string | null | undefined;
                host?: string | null | undefined;
                alpn?: string | null | undefined;
                fingerprint?: string | null | undefined;
                xHttpExtraParams?: {} | null | undefined;
                muxParams?: {} | null | undefined;
                sockoptParams?: {} | null | undefined;
                serverDescription?: string | null | undefined;
                allowInsecure?: boolean | null | undefined;
                shuffleHost?: boolean | null | undefined;
                mihomoX25519?: boolean | null | undefined;
                publicKey?: string | null | undefined;
                protocol?: string | null | undefined;
                shortId?: string | null | undefined;
                spiderX?: string | null | undefined;
                tls?: string | null | undefined;
                rawSettings?: {
                    headerType?: string | null | undefined;
                    request?: {} | null | undefined;
                } | null | undefined;
                additionalParams?: {
                    mode?: string | null | undefined;
                    heartbeatPeriod?: number | null | undefined;
                } | null | undefined;
                flow?: string | null | undefined;
                mldsa65Verify?: string | null | undefined;
                encryption?: string | null | undefined;
                protocolOptions?: {
                    ss?: {
                        method?: string | null | undefined;
                    } | null | undefined;
                } | null | undefined;
                dbData?: {
                    uuid: string;
                    tag: string | null;
                    rawInbound: {} | null;
                    viewPosition: number;
                    remark: string;
                    isDisabled: boolean;
                    configProfileUuid: string | null;
                    configProfileInboundUuid: string | null;
                    isHidden: boolean;
                    vlessRouteId: number | null;
                    inboundTag: string;
                } | undefined;
                xrayJsonTemplate?: {} | null | undefined;
            }[];
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
            headers: Record<string, string | undefined>;
            convertedUserInfo: {
                daysLeft: number;
                trafficUsed: string;
                trafficLimit: string;
                lifetimeTrafficUsed: string;
                isHwidLimited: boolean;
            };
            rawHosts: {
                password: {
                    trojanPassword: string;
                    ssPassword: string;
                    vlessPassword: string;
                };
                path?: string | null | undefined;
                network?: string | null | undefined;
                port?: number | null | undefined;
                remark?: string | null | undefined;
                address?: string | null | undefined;
                sni?: string | null | undefined;
                host?: string | null | undefined;
                alpn?: string | null | undefined;
                fingerprint?: string | null | undefined;
                xHttpExtraParams?: {} | null | undefined;
                muxParams?: {} | null | undefined;
                sockoptParams?: {} | null | undefined;
                serverDescription?: string | null | undefined;
                allowInsecure?: boolean | null | undefined;
                shuffleHost?: boolean | null | undefined;
                mihomoX25519?: boolean | null | undefined;
                publicKey?: string | null | undefined;
                protocol?: string | null | undefined;
                shortId?: string | null | undefined;
                spiderX?: string | null | undefined;
                tls?: string | null | undefined;
                rawSettings?: {
                    headerType?: string | null | undefined;
                    request?: {} | null | undefined;
                } | null | undefined;
                additionalParams?: {
                    mode?: string | null | undefined;
                    heartbeatPeriod?: number | null | undefined;
                } | null | undefined;
                flow?: string | null | undefined;
                mldsa65Verify?: string | null | undefined;
                encryption?: string | null | undefined;
                protocolOptions?: {
                    ss?: {
                        method?: string | null | undefined;
                    } | null | undefined;
                } | null | undefined;
                dbData?: {
                    uuid: string;
                    tag: string | null;
                    rawInbound: {} | null;
                    viewPosition: number;
                    remark: string;
                    isDisabled: boolean;
                    configProfileUuid: string | null;
                    configProfileInboundUuid: string | null;
                    isHidden: boolean;
                    vlessRouteId: number | null;
                    inboundTag: string;
                } | undefined;
                xrayJsonTemplate?: {} | null | undefined;
            }[];
        }>;
    }, "strip", z.ZodTypeAny, {
        response: {
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
            headers: Record<string, string | undefined>;
            convertedUserInfo: {
                daysLeft: number;
                trafficUsed: string;
                trafficLimit: string;
                lifetimeTrafficUsed: string;
                isHwidLimited: boolean;
            };
            rawHosts: {
                password: {
                    trojanPassword: string;
                    ssPassword: string;
                    vlessPassword: string;
                };
                path?: string | null | undefined;
                network?: string | null | undefined;
                port?: number | null | undefined;
                remark?: string | null | undefined;
                address?: string | null | undefined;
                sni?: string | null | undefined;
                host?: string | null | undefined;
                alpn?: string | null | undefined;
                fingerprint?: string | null | undefined;
                xHttpExtraParams?: {} | null | undefined;
                muxParams?: {} | null | undefined;
                sockoptParams?: {} | null | undefined;
                serverDescription?: string | null | undefined;
                allowInsecure?: boolean | null | undefined;
                shuffleHost?: boolean | null | undefined;
                mihomoX25519?: boolean | null | undefined;
                publicKey?: string | null | undefined;
                protocol?: string | null | undefined;
                shortId?: string | null | undefined;
                spiderX?: string | null | undefined;
                tls?: string | null | undefined;
                rawSettings?: {
                    headerType?: string | null | undefined;
                    request?: {} | null | undefined;
                } | null | undefined;
                additionalParams?: {
                    mode?: string | null | undefined;
                    heartbeatPeriod?: number | null | undefined;
                } | null | undefined;
                flow?: string | null | undefined;
                mldsa65Verify?: string | null | undefined;
                encryption?: string | null | undefined;
                protocolOptions?: {
                    ss?: {
                        method?: string | null | undefined;
                    } | null | undefined;
                } | null | undefined;
                dbData?: {
                    uuid: string;
                    tag: string | null;
                    rawInbound: {} | null;
                    viewPosition: number;
                    remark: string;
                    isDisabled: boolean;
                    configProfileUuid: string | null;
                    configProfileInboundUuid: string | null;
                    isHidden: boolean;
                    vlessRouteId: number | null;
                    inboundTag: string;
                } | undefined;
                xrayJsonTemplate?: {} | null | undefined;
            }[];
        };
    }, {
        response: {
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
            headers: Record<string, string | undefined>;
            convertedUserInfo: {
                daysLeft: number;
                trafficUsed: string;
                trafficLimit: string;
                lifetimeTrafficUsed: string;
                isHwidLimited: boolean;
            };
            rawHosts: {
                password: {
                    trojanPassword: string;
                    ssPassword: string;
                    vlessPassword: string;
                };
                path?: string | null | undefined;
                network?: string | null | undefined;
                port?: number | null | undefined;
                remark?: string | null | undefined;
                address?: string | null | undefined;
                sni?: string | null | undefined;
                host?: string | null | undefined;
                alpn?: string | null | undefined;
                fingerprint?: string | null | undefined;
                xHttpExtraParams?: {} | null | undefined;
                muxParams?: {} | null | undefined;
                sockoptParams?: {} | null | undefined;
                serverDescription?: string | null | undefined;
                allowInsecure?: boolean | null | undefined;
                shuffleHost?: boolean | null | undefined;
                mihomoX25519?: boolean | null | undefined;
                publicKey?: string | null | undefined;
                protocol?: string | null | undefined;
                shortId?: string | null | undefined;
                spiderX?: string | null | undefined;
                tls?: string | null | undefined;
                rawSettings?: {
                    headerType?: string | null | undefined;
                    request?: {} | null | undefined;
                } | null | undefined;
                additionalParams?: {
                    mode?: string | null | undefined;
                    heartbeatPeriod?: number | null | undefined;
                } | null | undefined;
                flow?: string | null | undefined;
                mldsa65Verify?: string | null | undefined;
                encryption?: string | null | undefined;
                protocolOptions?: {
                    ss?: {
                        method?: string | null | undefined;
                    } | null | undefined;
                } | null | undefined;
                dbData?: {
                    uuid: string;
                    tag: string | null;
                    rawInbound: {} | null;
                    viewPosition: number;
                    remark: string;
                    isDisabled: boolean;
                    configProfileUuid: string | null;
                    configProfileInboundUuid: string | null;
                    isHidden: boolean;
                    vlessRouteId: number | null;
                    inboundTag: string;
                } | undefined;
                xrayJsonTemplate?: {} | null | undefined;
            }[];
        };
    }>;
    type Response = z.infer<typeof ResponseSchema>;
}
//# sourceMappingURL=get-raw-subscription-by-short-uuid.command.d.ts.map