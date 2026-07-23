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
            convertedUserInfo: z.ZodObject<{
                daysLeft: z.ZodNumber;
                trafficLimit: z.ZodString;
                trafficUsed: z.ZodString;
                lifetimeTrafficUsed: z.ZodString;
                hwidCheckup: z.ZodNullable<z.ZodObject<{
                    subscriptionAllowed: z.ZodBoolean;
                    maxDeviceReached: z.ZodBoolean;
                    hwidNotSupported: z.ZodBoolean;
                    limitBypassed: z.ZodBoolean;
                }, "strip", z.ZodTypeAny, {
                    subscriptionAllowed: boolean;
                    maxDeviceReached: boolean;
                    hwidNotSupported: boolean;
                    limitBypassed: boolean;
                }, {
                    subscriptionAllowed: boolean;
                    maxDeviceReached: boolean;
                    hwidNotSupported: boolean;
                    limitBypassed: boolean;
                }>>;
            }, "strip", z.ZodTypeAny, {
                daysLeft: number;
                trafficUsed: string;
                trafficLimit: string;
                lifetimeTrafficUsed: string;
                hwidCheckup: {
                    subscriptionAllowed: boolean;
                    maxDeviceReached: boolean;
                    hwidNotSupported: boolean;
                    limitBypassed: boolean;
                } | null;
            }, {
                daysLeft: number;
                trafficUsed: string;
                trafficLimit: string;
                lifetimeTrafficUsed: string;
                hwidCheckup: {
                    subscriptionAllowed: boolean;
                    maxDeviceReached: boolean;
                    hwidNotSupported: boolean;
                    limitBypassed: boolean;
                } | null;
            }>;
            headers: z.ZodRecord<z.ZodString, z.ZodOptional<z.ZodString>>;
            resolvedProxyConfigs: z.ZodArray<z.ZodObject<{
                finalRemark: z.ZodString;
                address: z.ZodString;
                port: z.ZodNumber;
                protocol: z.ZodEnum<["vless", "trojan", "shadowsocks", "hysteria"]>;
                protocolOptions: z.ZodUnion<[z.ZodObject<{
                    encryption: z.ZodString;
                    id: z.ZodString;
                    flow: z.ZodEnum<["", "xtls-rprx-vision", "xtls-rprx-vision-udp443"]>;
                }, "strip", z.ZodTypeAny, {
                    id: string;
                    encryption: string;
                    flow: "" | "xtls-rprx-vision" | "xtls-rprx-vision-udp443";
                }, {
                    id: string;
                    encryption: string;
                    flow: "" | "xtls-rprx-vision" | "xtls-rprx-vision-udp443";
                }>, z.ZodObject<{
                    password: z.ZodString;
                }, "strip", z.ZodTypeAny, {
                    password: string;
                }, {
                    password: string;
                }>, z.ZodObject<{
                    method: z.ZodString;
                    password: z.ZodString;
                    uot: z.ZodBoolean;
                    uotVersion: z.ZodNumber;
                }, "strip", z.ZodTypeAny, {
                    method: string;
                    password: string;
                    uot: boolean;
                    uotVersion: number;
                }, {
                    method: string;
                    password: string;
                    uot: boolean;
                    uotVersion: number;
                }>, z.ZodObject<{
                    version: z.ZodNumber;
                }, "strip", z.ZodTypeAny, {
                    version: number;
                }, {
                    version: number;
                }>]>;
                transport: z.ZodEnum<["tcp", "xhttp", "ws", "httpupgrade", "grpc", "kcp", "hysteria"]>;
                transportOptions: z.ZodUnion<[z.ZodObject<{
                    header: z.ZodNullable<z.ZodDiscriminatedUnion<"type", [z.ZodObject<{
                        type: z.ZodLiteral<"none">;
                    }, "strip", z.ZodTypeAny, {
                        type: "none";
                    }, {
                        type: "none";
                    }>, z.ZodObject<{
                        type: z.ZodLiteral<"http">;
                        request: z.ZodOptional<z.ZodObject<{
                            version: z.ZodOptional<z.ZodString>;
                            method: z.ZodOptional<z.ZodString>;
                            path: z.ZodOptional<z.ZodArray<z.ZodString, "many">>;
                            headers: z.ZodOptional<z.ZodRecord<z.ZodString, z.ZodUnknown>>;
                        }, "strip", z.ZodTypeAny, {
                            path?: string[] | undefined;
                            method?: string | undefined;
                            headers?: Record<string, unknown> | undefined;
                            version?: string | undefined;
                        }, {
                            path?: string[] | undefined;
                            method?: string | undefined;
                            headers?: Record<string, unknown> | undefined;
                            version?: string | undefined;
                        }>>;
                        response: z.ZodOptional<z.ZodObject<{
                            version: z.ZodOptional<z.ZodString>;
                            status: z.ZodOptional<z.ZodString>;
                            reason: z.ZodOptional<z.ZodString>;
                            headers: z.ZodOptional<z.ZodRecord<z.ZodString, z.ZodUnknown>>;
                        }, "strip", z.ZodTypeAny, {
                            status?: string | undefined;
                            headers?: Record<string, unknown> | undefined;
                            version?: string | undefined;
                            reason?: string | undefined;
                        }, {
                            status?: string | undefined;
                            headers?: Record<string, unknown> | undefined;
                            version?: string | undefined;
                            reason?: string | undefined;
                        }>>;
                    }, "strip", z.ZodTypeAny, {
                        type: "http";
                        response?: {
                            status?: string | undefined;
                            headers?: Record<string, unknown> | undefined;
                            version?: string | undefined;
                            reason?: string | undefined;
                        } | undefined;
                        request?: {
                            path?: string[] | undefined;
                            method?: string | undefined;
                            headers?: Record<string, unknown> | undefined;
                            version?: string | undefined;
                        } | undefined;
                    }, {
                        type: "http";
                        response?: {
                            status?: string | undefined;
                            headers?: Record<string, unknown> | undefined;
                            version?: string | undefined;
                            reason?: string | undefined;
                        } | undefined;
                        request?: {
                            path?: string[] | undefined;
                            method?: string | undefined;
                            headers?: Record<string, unknown> | undefined;
                            version?: string | undefined;
                        } | undefined;
                    }>]>>;
                }, "strip", z.ZodTypeAny, {
                    header: {
                        type: "none";
                    } | {
                        type: "http";
                        response?: {
                            status?: string | undefined;
                            headers?: Record<string, unknown> | undefined;
                            version?: string | undefined;
                            reason?: string | undefined;
                        } | undefined;
                        request?: {
                            path?: string[] | undefined;
                            method?: string | undefined;
                            headers?: Record<string, unknown> | undefined;
                            version?: string | undefined;
                        } | undefined;
                    } | null;
                }, {
                    header: {
                        type: "none";
                    } | {
                        type: "http";
                        response?: {
                            status?: string | undefined;
                            headers?: Record<string, unknown> | undefined;
                            version?: string | undefined;
                            reason?: string | undefined;
                        } | undefined;
                        request?: {
                            path?: string[] | undefined;
                            method?: string | undefined;
                            headers?: Record<string, unknown> | undefined;
                            version?: string | undefined;
                        } | undefined;
                    } | null;
                }>, z.ZodObject<{
                    path: z.ZodNullable<z.ZodString>;
                    host: z.ZodNullable<z.ZodString>;
                    mode: z.ZodEnum<["auto", "packet-up", "stream-up", "stream-one"]>;
                    extra: z.ZodNullable<z.ZodRecord<z.ZodString, z.ZodUnknown>>;
                }, "strip", z.ZodTypeAny, {
                    path: string | null;
                    host: string | null;
                    mode: "auto" | "packet-up" | "stream-up" | "stream-one";
                    extra: Record<string, unknown> | null;
                }, {
                    path: string | null;
                    host: string | null;
                    mode: "auto" | "packet-up" | "stream-up" | "stream-one";
                    extra: Record<string, unknown> | null;
                }>, z.ZodObject<{
                    path: z.ZodNullable<z.ZodString>;
                    host: z.ZodNullable<z.ZodString>;
                    headers: z.ZodNullable<z.ZodRecord<z.ZodString, z.ZodString>>;
                    heartbeatPeriod: z.ZodNullable<z.ZodNumber>;
                }, "strip", z.ZodTypeAny, {
                    path: string | null;
                    host: string | null;
                    headers: Record<string, string> | null;
                    heartbeatPeriod: number | null;
                }, {
                    path: string | null;
                    host: string | null;
                    headers: Record<string, string> | null;
                    heartbeatPeriod: number | null;
                }>, z.ZodObject<{
                    path: z.ZodNullable<z.ZodString>;
                    host: z.ZodNullable<z.ZodString>;
                    headers: z.ZodNullable<z.ZodRecord<z.ZodString, z.ZodString>>;
                }, "strip", z.ZodTypeAny, {
                    path: string | null;
                    host: string | null;
                    headers: Record<string, string> | null;
                }, {
                    path: string | null;
                    host: string | null;
                    headers: Record<string, string> | null;
                }>, z.ZodObject<{
                    authority: z.ZodNullable<z.ZodString>;
                    serviceName: z.ZodNullable<z.ZodString>;
                    multiMode: z.ZodBoolean;
                }, "strip", z.ZodTypeAny, {
                    authority: string | null;
                    serviceName: string | null;
                    multiMode: boolean;
                }, {
                    authority: string | null;
                    serviceName: string | null;
                    multiMode: boolean;
                }>, z.ZodObject<{
                    clientMtu: z.ZodNumber;
                    clientTti: z.ZodNumber;
                    congestion: z.ZodBoolean;
                }, "strip", z.ZodTypeAny, {
                    clientMtu: number;
                    clientTti: number;
                    congestion: boolean;
                }, {
                    clientMtu: number;
                    clientTti: number;
                    congestion: boolean;
                }>, z.ZodObject<{
                    version: z.ZodNumber;
                    auth: z.ZodString;
                }, "strip", z.ZodTypeAny, {
                    auth: string;
                    version: number;
                }, {
                    auth: string;
                    version: number;
                }>]>;
                security: z.ZodEnum<["tls", "reality", "none"]>;
                securityOptions: z.ZodOptional<z.ZodUnion<[z.ZodObject<{
                    pinnedPeerCertSha256: z.ZodNullable<z.ZodString>;
                    verifyPeerCertByName: z.ZodNullable<z.ZodString>;
                    alpn: z.ZodNullable<z.ZodString>;
                    enableSessionResumption: z.ZodBoolean;
                    fingerprint: z.ZodNullable<z.ZodString>;
                    serverName: z.ZodNullable<z.ZodString>;
                    echConfigList: z.ZodNullable<z.ZodString>;
                    echForceQuery: z.ZodNullable<z.ZodString>;
                }, "strip", z.ZodTypeAny, {
                    alpn: string | null;
                    fingerprint: string | null;
                    pinnedPeerCertSha256: string | null;
                    verifyPeerCertByName: string | null;
                    enableSessionResumption: boolean;
                    serverName: string | null;
                    echConfigList: string | null;
                    echForceQuery: string | null;
                }, {
                    alpn: string | null;
                    fingerprint: string | null;
                    pinnedPeerCertSha256: string | null;
                    verifyPeerCertByName: string | null;
                    enableSessionResumption: boolean;
                    serverName: string | null;
                    echConfigList: string | null;
                    echForceQuery: string | null;
                }>, z.ZodObject<{
                    fingerprint: z.ZodString;
                    publicKey: z.ZodString;
                    shortId: z.ZodNullable<z.ZodString>;
                    serverName: z.ZodString;
                    spiderX: z.ZodNullable<z.ZodString>;
                    mldsa65Verify: z.ZodNullable<z.ZodString>;
                }, "strip", z.ZodTypeAny, {
                    fingerprint: string;
                    serverName: string;
                    publicKey: string;
                    shortId: string | null;
                    spiderX: string | null;
                    mldsa65Verify: string | null;
                }, {
                    fingerprint: string;
                    serverName: string;
                    publicKey: string;
                    shortId: string | null;
                    spiderX: string | null;
                    mldsa65Verify: string | null;
                }>]>>;
                streamOverrides: z.ZodObject<{
                    finalMask: z.ZodNullable<z.ZodUnknown>;
                    sockopt: z.ZodNullable<z.ZodUnknown>;
                }, "strip", z.ZodTypeAny, {
                    finalMask?: unknown;
                    sockopt?: unknown;
                }, {
                    finalMask?: unknown;
                    sockopt?: unknown;
                }>;
                mux: z.ZodNullable<z.ZodUnknown>;
                clientOverrides: z.ZodObject<{
                    shuffleHost: z.ZodBoolean;
                    mihomoX25519: z.ZodBoolean;
                    mihomoIpVersion: z.ZodNullable<z.ZodNativeEnum<{
                        readonly DUAL: "dual";
                        readonly IPV4: "ipv4";
                        readonly IPV6: "ipv6";
                        readonly IPV4_PREFER: "ipv4-prefer";
                        readonly IPV6_PREFER: "ipv6-prefer";
                    }>>;
                    serverDescription: z.ZodNullable<z.ZodString>;
                    xrayJsonTemplate: z.ZodNullable<z.ZodUnknown>;
                }, "strip", z.ZodTypeAny, {
                    serverDescription: string | null;
                    shuffleHost: boolean;
                    mihomoX25519: boolean;
                    mihomoIpVersion: "dual" | "ipv4" | "ipv6" | "ipv4-prefer" | "ipv6-prefer" | null;
                    xrayJsonTemplate?: unknown;
                }, {
                    serverDescription: string | null;
                    shuffleHost: boolean;
                    mihomoX25519: boolean;
                    mihomoIpVersion: "dual" | "ipv4" | "ipv6" | "ipv4-prefer" | "ipv6-prefer" | null;
                    xrayJsonTemplate?: unknown;
                }>;
                metadata: z.ZodObject<{
                    uuid: z.ZodString;
                    tags: z.ZodArray<z.ZodString, "many">;
                    excludeFromSubscriptionTypes: z.ZodArray<z.ZodNativeEnum<{
                        readonly XRAY_JSON: "XRAY_JSON";
                        readonly XRAY_BASE64: "XRAY_BASE64";
                        readonly MIHOMO: "MIHOMO";
                        readonly STASH: "STASH";
                        readonly CLASH: "CLASH";
                        readonly SINGBOX: "SINGBOX";
                    }>, "many">;
                    inboundTag: z.ZodString;
                    configProfileUuid: z.ZodNullable<z.ZodString>;
                    configProfileInboundUuid: z.ZodNullable<z.ZodString>;
                    isDisabled: z.ZodBoolean;
                    isHidden: z.ZodBoolean;
                    viewPosition: z.ZodNumber;
                    remark: z.ZodString;
                    vlessRouteId: z.ZodNullable<z.ZodNumber>;
                    rawInbound: z.ZodNullable<z.ZodUnknown>;
                }, "strip", z.ZodTypeAny, {
                    tags: string[];
                    uuid: string;
                    viewPosition: number;
                    remark: string;
                    isDisabled: boolean;
                    configProfileUuid: string | null;
                    configProfileInboundUuid: string | null;
                    isHidden: boolean;
                    vlessRouteId: number | null;
                    excludeFromSubscriptionTypes: ("STASH" | "SINGBOX" | "MIHOMO" | "XRAY_JSON" | "CLASH" | "XRAY_BASE64")[];
                    inboundTag: string;
                    rawInbound?: unknown;
                }, {
                    tags: string[];
                    uuid: string;
                    viewPosition: number;
                    remark: string;
                    isDisabled: boolean;
                    configProfileUuid: string | null;
                    configProfileInboundUuid: string | null;
                    isHidden: boolean;
                    vlessRouteId: number | null;
                    excludeFromSubscriptionTypes: ("STASH" | "SINGBOX" | "MIHOMO" | "XRAY_JSON" | "CLASH" | "XRAY_BASE64")[];
                    inboundTag: string;
                    rawInbound?: unknown;
                }>;
            }, "strip", z.ZodTypeAny, {
                metadata: {
                    tags: string[];
                    uuid: string;
                    viewPosition: number;
                    remark: string;
                    isDisabled: boolean;
                    configProfileUuid: string | null;
                    configProfileInboundUuid: string | null;
                    isHidden: boolean;
                    vlessRouteId: number | null;
                    excludeFromSubscriptionTypes: ("STASH" | "SINGBOX" | "MIHOMO" | "XRAY_JSON" | "CLASH" | "XRAY_BASE64")[];
                    inboundTag: string;
                    rawInbound?: unknown;
                };
                security: "none" | "tls" | "reality";
                port: number;
                address: string;
                protocol: "vless" | "trojan" | "shadowsocks" | "hysteria";
                protocolOptions: {
                    id: string;
                    encryption: string;
                    flow: "" | "xtls-rprx-vision" | "xtls-rprx-vision-udp443";
                } | {
                    password: string;
                } | {
                    method: string;
                    password: string;
                    uot: boolean;
                    uotVersion: number;
                } | {
                    version: number;
                };
                transport: "hysteria" | "tcp" | "xhttp" | "ws" | "httpupgrade" | "grpc" | "kcp";
                transportOptions: {
                    header: {
                        type: "none";
                    } | {
                        type: "http";
                        response?: {
                            status?: string | undefined;
                            headers?: Record<string, unknown> | undefined;
                            version?: string | undefined;
                            reason?: string | undefined;
                        } | undefined;
                        request?: {
                            path?: string[] | undefined;
                            method?: string | undefined;
                            headers?: Record<string, unknown> | undefined;
                            version?: string | undefined;
                        } | undefined;
                    } | null;
                } | {
                    path: string | null;
                    host: string | null;
                    mode: "auto" | "packet-up" | "stream-up" | "stream-one";
                    extra: Record<string, unknown> | null;
                } | {
                    path: string | null;
                    host: string | null;
                    headers: Record<string, string> | null;
                    heartbeatPeriod: number | null;
                } | {
                    path: string | null;
                    host: string | null;
                    headers: Record<string, string> | null;
                } | {
                    authority: string | null;
                    serviceName: string | null;
                    multiMode: boolean;
                } | {
                    clientMtu: number;
                    clientTti: number;
                    congestion: boolean;
                } | {
                    auth: string;
                    version: number;
                };
                finalRemark: string;
                streamOverrides: {
                    finalMask?: unknown;
                    sockopt?: unknown;
                };
                clientOverrides: {
                    serverDescription: string | null;
                    shuffleHost: boolean;
                    mihomoX25519: boolean;
                    mihomoIpVersion: "dual" | "ipv4" | "ipv6" | "ipv4-prefer" | "ipv6-prefer" | null;
                    xrayJsonTemplate?: unknown;
                };
                securityOptions?: {
                    alpn: string | null;
                    fingerprint: string | null;
                    pinnedPeerCertSha256: string | null;
                    verifyPeerCertByName: string | null;
                    enableSessionResumption: boolean;
                    serverName: string | null;
                    echConfigList: string | null;
                    echForceQuery: string | null;
                } | {
                    fingerprint: string;
                    serverName: string;
                    publicKey: string;
                    shortId: string | null;
                    spiderX: string | null;
                    mldsa65Verify: string | null;
                } | undefined;
                mux?: unknown;
            }, {
                metadata: {
                    tags: string[];
                    uuid: string;
                    viewPosition: number;
                    remark: string;
                    isDisabled: boolean;
                    configProfileUuid: string | null;
                    configProfileInboundUuid: string | null;
                    isHidden: boolean;
                    vlessRouteId: number | null;
                    excludeFromSubscriptionTypes: ("STASH" | "SINGBOX" | "MIHOMO" | "XRAY_JSON" | "CLASH" | "XRAY_BASE64")[];
                    inboundTag: string;
                    rawInbound?: unknown;
                };
                security: "none" | "tls" | "reality";
                port: number;
                address: string;
                protocol: "vless" | "trojan" | "shadowsocks" | "hysteria";
                protocolOptions: {
                    id: string;
                    encryption: string;
                    flow: "" | "xtls-rprx-vision" | "xtls-rprx-vision-udp443";
                } | {
                    password: string;
                } | {
                    method: string;
                    password: string;
                    uot: boolean;
                    uotVersion: number;
                } | {
                    version: number;
                };
                transport: "hysteria" | "tcp" | "xhttp" | "ws" | "httpupgrade" | "grpc" | "kcp";
                transportOptions: {
                    header: {
                        type: "none";
                    } | {
                        type: "http";
                        response?: {
                            status?: string | undefined;
                            headers?: Record<string, unknown> | undefined;
                            version?: string | undefined;
                            reason?: string | undefined;
                        } | undefined;
                        request?: {
                            path?: string[] | undefined;
                            method?: string | undefined;
                            headers?: Record<string, unknown> | undefined;
                            version?: string | undefined;
                        } | undefined;
                    } | null;
                } | {
                    path: string | null;
                    host: string | null;
                    mode: "auto" | "packet-up" | "stream-up" | "stream-one";
                    extra: Record<string, unknown> | null;
                } | {
                    path: string | null;
                    host: string | null;
                    headers: Record<string, string> | null;
                    heartbeatPeriod: number | null;
                } | {
                    path: string | null;
                    host: string | null;
                    headers: Record<string, string> | null;
                } | {
                    authority: string | null;
                    serviceName: string | null;
                    multiMode: boolean;
                } | {
                    clientMtu: number;
                    clientTti: number;
                    congestion: boolean;
                } | {
                    auth: string;
                    version: number;
                };
                finalRemark: string;
                streamOverrides: {
                    finalMask?: unknown;
                    sockopt?: unknown;
                };
                clientOverrides: {
                    serverDescription: string | null;
                    shuffleHost: boolean;
                    mihomoX25519: boolean;
                    mihomoIpVersion: "dual" | "ipv4" | "ipv6" | "ipv4-prefer" | "ipv6-prefer" | null;
                    xrayJsonTemplate?: unknown;
                };
                securityOptions?: {
                    alpn: string | null;
                    fingerprint: string | null;
                    pinnedPeerCertSha256: string | null;
                    verifyPeerCertByName: string | null;
                    enableSessionResumption: boolean;
                    serverName: string | null;
                    echConfigList: string | null;
                    echForceQuery: string | null;
                } | {
                    fingerprint: string;
                    serverName: string;
                    publicKey: string;
                    shortId: string | null;
                    spiderX: string | null;
                    mldsa65Verify: string | null;
                } | undefined;
                mux?: unknown;
            }>, "many">;
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
            headers: Record<string, string | undefined>;
            convertedUserInfo: {
                daysLeft: number;
                trafficUsed: string;
                trafficLimit: string;
                lifetimeTrafficUsed: string;
                hwidCheckup: {
                    subscriptionAllowed: boolean;
                    maxDeviceReached: boolean;
                    hwidNotSupported: boolean;
                    limitBypassed: boolean;
                } | null;
            };
            resolvedProxyConfigs: {
                metadata: {
                    tags: string[];
                    uuid: string;
                    viewPosition: number;
                    remark: string;
                    isDisabled: boolean;
                    configProfileUuid: string | null;
                    configProfileInboundUuid: string | null;
                    isHidden: boolean;
                    vlessRouteId: number | null;
                    excludeFromSubscriptionTypes: ("STASH" | "SINGBOX" | "MIHOMO" | "XRAY_JSON" | "CLASH" | "XRAY_BASE64")[];
                    inboundTag: string;
                    rawInbound?: unknown;
                };
                security: "none" | "tls" | "reality";
                port: number;
                address: string;
                protocol: "vless" | "trojan" | "shadowsocks" | "hysteria";
                protocolOptions: {
                    id: string;
                    encryption: string;
                    flow: "" | "xtls-rprx-vision" | "xtls-rprx-vision-udp443";
                } | {
                    password: string;
                } | {
                    method: string;
                    password: string;
                    uot: boolean;
                    uotVersion: number;
                } | {
                    version: number;
                };
                transport: "hysteria" | "tcp" | "xhttp" | "ws" | "httpupgrade" | "grpc" | "kcp";
                transportOptions: {
                    header: {
                        type: "none";
                    } | {
                        type: "http";
                        response?: {
                            status?: string | undefined;
                            headers?: Record<string, unknown> | undefined;
                            version?: string | undefined;
                            reason?: string | undefined;
                        } | undefined;
                        request?: {
                            path?: string[] | undefined;
                            method?: string | undefined;
                            headers?: Record<string, unknown> | undefined;
                            version?: string | undefined;
                        } | undefined;
                    } | null;
                } | {
                    path: string | null;
                    host: string | null;
                    mode: "auto" | "packet-up" | "stream-up" | "stream-one";
                    extra: Record<string, unknown> | null;
                } | {
                    path: string | null;
                    host: string | null;
                    headers: Record<string, string> | null;
                    heartbeatPeriod: number | null;
                } | {
                    path: string | null;
                    host: string | null;
                    headers: Record<string, string> | null;
                } | {
                    authority: string | null;
                    serviceName: string | null;
                    multiMode: boolean;
                } | {
                    clientMtu: number;
                    clientTti: number;
                    congestion: boolean;
                } | {
                    auth: string;
                    version: number;
                };
                finalRemark: string;
                streamOverrides: {
                    finalMask?: unknown;
                    sockopt?: unknown;
                };
                clientOverrides: {
                    serverDescription: string | null;
                    shuffleHost: boolean;
                    mihomoX25519: boolean;
                    mihomoIpVersion: "dual" | "ipv4" | "ipv6" | "ipv4-prefer" | "ipv6-prefer" | null;
                    xrayJsonTemplate?: unknown;
                };
                securityOptions?: {
                    alpn: string | null;
                    fingerprint: string | null;
                    pinnedPeerCertSha256: string | null;
                    verifyPeerCertByName: string | null;
                    enableSessionResumption: boolean;
                    serverName: string | null;
                    echConfigList: string | null;
                    echForceQuery: string | null;
                } | {
                    fingerprint: string;
                    serverName: string;
                    publicKey: string;
                    shortId: string | null;
                    spiderX: string | null;
                    mldsa65Verify: string | null;
                } | undefined;
                mux?: unknown;
            }[];
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
            headers: Record<string, string | undefined>;
            convertedUserInfo: {
                daysLeft: number;
                trafficUsed: string;
                trafficLimit: string;
                lifetimeTrafficUsed: string;
                hwidCheckup: {
                    subscriptionAllowed: boolean;
                    maxDeviceReached: boolean;
                    hwidNotSupported: boolean;
                    limitBypassed: boolean;
                } | null;
            };
            resolvedProxyConfigs: {
                metadata: {
                    tags: string[];
                    uuid: string;
                    viewPosition: number;
                    remark: string;
                    isDisabled: boolean;
                    configProfileUuid: string | null;
                    configProfileInboundUuid: string | null;
                    isHidden: boolean;
                    vlessRouteId: number | null;
                    excludeFromSubscriptionTypes: ("STASH" | "SINGBOX" | "MIHOMO" | "XRAY_JSON" | "CLASH" | "XRAY_BASE64")[];
                    inboundTag: string;
                    rawInbound?: unknown;
                };
                security: "none" | "tls" | "reality";
                port: number;
                address: string;
                protocol: "vless" | "trojan" | "shadowsocks" | "hysteria";
                protocolOptions: {
                    id: string;
                    encryption: string;
                    flow: "" | "xtls-rprx-vision" | "xtls-rprx-vision-udp443";
                } | {
                    password: string;
                } | {
                    method: string;
                    password: string;
                    uot: boolean;
                    uotVersion: number;
                } | {
                    version: number;
                };
                transport: "hysteria" | "tcp" | "xhttp" | "ws" | "httpupgrade" | "grpc" | "kcp";
                transportOptions: {
                    header: {
                        type: "none";
                    } | {
                        type: "http";
                        response?: {
                            status?: string | undefined;
                            headers?: Record<string, unknown> | undefined;
                            version?: string | undefined;
                            reason?: string | undefined;
                        } | undefined;
                        request?: {
                            path?: string[] | undefined;
                            method?: string | undefined;
                            headers?: Record<string, unknown> | undefined;
                            version?: string | undefined;
                        } | undefined;
                    } | null;
                } | {
                    path: string | null;
                    host: string | null;
                    mode: "auto" | "packet-up" | "stream-up" | "stream-one";
                    extra: Record<string, unknown> | null;
                } | {
                    path: string | null;
                    host: string | null;
                    headers: Record<string, string> | null;
                    heartbeatPeriod: number | null;
                } | {
                    path: string | null;
                    host: string | null;
                    headers: Record<string, string> | null;
                } | {
                    authority: string | null;
                    serviceName: string | null;
                    multiMode: boolean;
                } | {
                    clientMtu: number;
                    clientTti: number;
                    congestion: boolean;
                } | {
                    auth: string;
                    version: number;
                };
                finalRemark: string;
                streamOverrides: {
                    finalMask?: unknown;
                    sockopt?: unknown;
                };
                clientOverrides: {
                    serverDescription: string | null;
                    shuffleHost: boolean;
                    mihomoX25519: boolean;
                    mihomoIpVersion: "dual" | "ipv4" | "ipv6" | "ipv4-prefer" | "ipv6-prefer" | null;
                    xrayJsonTemplate?: unknown;
                };
                securityOptions?: {
                    alpn: string | null;
                    fingerprint: string | null;
                    pinnedPeerCertSha256: string | null;
                    verifyPeerCertByName: string | null;
                    enableSessionResumption: boolean;
                    serverName: string | null;
                    echConfigList: string | null;
                    echForceQuery: string | null;
                } | {
                    fingerprint: string;
                    serverName: string;
                    publicKey: string;
                    shortId: string | null;
                    spiderX: string | null;
                    mldsa65Verify: string | null;
                } | undefined;
                mux?: unknown;
            }[];
        }>;
    }, "strip", z.ZodTypeAny, {
        response: {
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
            headers: Record<string, string | undefined>;
            convertedUserInfo: {
                daysLeft: number;
                trafficUsed: string;
                trafficLimit: string;
                lifetimeTrafficUsed: string;
                hwidCheckup: {
                    subscriptionAllowed: boolean;
                    maxDeviceReached: boolean;
                    hwidNotSupported: boolean;
                    limitBypassed: boolean;
                } | null;
            };
            resolvedProxyConfigs: {
                metadata: {
                    tags: string[];
                    uuid: string;
                    viewPosition: number;
                    remark: string;
                    isDisabled: boolean;
                    configProfileUuid: string | null;
                    configProfileInboundUuid: string | null;
                    isHidden: boolean;
                    vlessRouteId: number | null;
                    excludeFromSubscriptionTypes: ("STASH" | "SINGBOX" | "MIHOMO" | "XRAY_JSON" | "CLASH" | "XRAY_BASE64")[];
                    inboundTag: string;
                    rawInbound?: unknown;
                };
                security: "none" | "tls" | "reality";
                port: number;
                address: string;
                protocol: "vless" | "trojan" | "shadowsocks" | "hysteria";
                protocolOptions: {
                    id: string;
                    encryption: string;
                    flow: "" | "xtls-rprx-vision" | "xtls-rprx-vision-udp443";
                } | {
                    password: string;
                } | {
                    method: string;
                    password: string;
                    uot: boolean;
                    uotVersion: number;
                } | {
                    version: number;
                };
                transport: "hysteria" | "tcp" | "xhttp" | "ws" | "httpupgrade" | "grpc" | "kcp";
                transportOptions: {
                    header: {
                        type: "none";
                    } | {
                        type: "http";
                        response?: {
                            status?: string | undefined;
                            headers?: Record<string, unknown> | undefined;
                            version?: string | undefined;
                            reason?: string | undefined;
                        } | undefined;
                        request?: {
                            path?: string[] | undefined;
                            method?: string | undefined;
                            headers?: Record<string, unknown> | undefined;
                            version?: string | undefined;
                        } | undefined;
                    } | null;
                } | {
                    path: string | null;
                    host: string | null;
                    mode: "auto" | "packet-up" | "stream-up" | "stream-one";
                    extra: Record<string, unknown> | null;
                } | {
                    path: string | null;
                    host: string | null;
                    headers: Record<string, string> | null;
                    heartbeatPeriod: number | null;
                } | {
                    path: string | null;
                    host: string | null;
                    headers: Record<string, string> | null;
                } | {
                    authority: string | null;
                    serviceName: string | null;
                    multiMode: boolean;
                } | {
                    clientMtu: number;
                    clientTti: number;
                    congestion: boolean;
                } | {
                    auth: string;
                    version: number;
                };
                finalRemark: string;
                streamOverrides: {
                    finalMask?: unknown;
                    sockopt?: unknown;
                };
                clientOverrides: {
                    serverDescription: string | null;
                    shuffleHost: boolean;
                    mihomoX25519: boolean;
                    mihomoIpVersion: "dual" | "ipv4" | "ipv6" | "ipv4-prefer" | "ipv6-prefer" | null;
                    xrayJsonTemplate?: unknown;
                };
                securityOptions?: {
                    alpn: string | null;
                    fingerprint: string | null;
                    pinnedPeerCertSha256: string | null;
                    verifyPeerCertByName: string | null;
                    enableSessionResumption: boolean;
                    serverName: string | null;
                    echConfigList: string | null;
                    echForceQuery: string | null;
                } | {
                    fingerprint: string;
                    serverName: string;
                    publicKey: string;
                    shortId: string | null;
                    spiderX: string | null;
                    mldsa65Verify: string | null;
                } | undefined;
                mux?: unknown;
            }[];
        };
    }, {
        response: {
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
            headers: Record<string, string | undefined>;
            convertedUserInfo: {
                daysLeft: number;
                trafficUsed: string;
                trafficLimit: string;
                lifetimeTrafficUsed: string;
                hwidCheckup: {
                    subscriptionAllowed: boolean;
                    maxDeviceReached: boolean;
                    hwidNotSupported: boolean;
                    limitBypassed: boolean;
                } | null;
            };
            resolvedProxyConfigs: {
                metadata: {
                    tags: string[];
                    uuid: string;
                    viewPosition: number;
                    remark: string;
                    isDisabled: boolean;
                    configProfileUuid: string | null;
                    configProfileInboundUuid: string | null;
                    isHidden: boolean;
                    vlessRouteId: number | null;
                    excludeFromSubscriptionTypes: ("STASH" | "SINGBOX" | "MIHOMO" | "XRAY_JSON" | "CLASH" | "XRAY_BASE64")[];
                    inboundTag: string;
                    rawInbound?: unknown;
                };
                security: "none" | "tls" | "reality";
                port: number;
                address: string;
                protocol: "vless" | "trojan" | "shadowsocks" | "hysteria";
                protocolOptions: {
                    id: string;
                    encryption: string;
                    flow: "" | "xtls-rprx-vision" | "xtls-rprx-vision-udp443";
                } | {
                    password: string;
                } | {
                    method: string;
                    password: string;
                    uot: boolean;
                    uotVersion: number;
                } | {
                    version: number;
                };
                transport: "hysteria" | "tcp" | "xhttp" | "ws" | "httpupgrade" | "grpc" | "kcp";
                transportOptions: {
                    header: {
                        type: "none";
                    } | {
                        type: "http";
                        response?: {
                            status?: string | undefined;
                            headers?: Record<string, unknown> | undefined;
                            version?: string | undefined;
                            reason?: string | undefined;
                        } | undefined;
                        request?: {
                            path?: string[] | undefined;
                            method?: string | undefined;
                            headers?: Record<string, unknown> | undefined;
                            version?: string | undefined;
                        } | undefined;
                    } | null;
                } | {
                    path: string | null;
                    host: string | null;
                    mode: "auto" | "packet-up" | "stream-up" | "stream-one";
                    extra: Record<string, unknown> | null;
                } | {
                    path: string | null;
                    host: string | null;
                    headers: Record<string, string> | null;
                    heartbeatPeriod: number | null;
                } | {
                    path: string | null;
                    host: string | null;
                    headers: Record<string, string> | null;
                } | {
                    authority: string | null;
                    serviceName: string | null;
                    multiMode: boolean;
                } | {
                    clientMtu: number;
                    clientTti: number;
                    congestion: boolean;
                } | {
                    auth: string;
                    version: number;
                };
                finalRemark: string;
                streamOverrides: {
                    finalMask?: unknown;
                    sockopt?: unknown;
                };
                clientOverrides: {
                    serverDescription: string | null;
                    shuffleHost: boolean;
                    mihomoX25519: boolean;
                    mihomoIpVersion: "dual" | "ipv4" | "ipv6" | "ipv4-prefer" | "ipv6-prefer" | null;
                    xrayJsonTemplate?: unknown;
                };
                securityOptions?: {
                    alpn: string | null;
                    fingerprint: string | null;
                    pinnedPeerCertSha256: string | null;
                    verifyPeerCertByName: string | null;
                    enableSessionResumption: boolean;
                    serverName: string | null;
                    echConfigList: string | null;
                    echForceQuery: string | null;
                } | {
                    fingerprint: string;
                    serverName: string;
                    publicKey: string;
                    shortId: string | null;
                    spiderX: string | null;
                    mldsa65Verify: string | null;
                } | undefined;
                mux?: unknown;
            }[];
        };
    }>;
    type Response = z.infer<typeof ResponseSchema>;
}
//# sourceMappingURL=get-raw-subscription-by-short-uuid.command.d.ts.map