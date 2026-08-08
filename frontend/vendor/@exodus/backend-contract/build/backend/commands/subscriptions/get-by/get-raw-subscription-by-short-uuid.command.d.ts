import { z } from 'zod';
export declare namespace GetRawSubscriptionByShortUuidCommand {
    const url: (shortUuid: string) => string;
    const TSQ_url: string;
    const endpointDetails: import("../../../constants").EndpointDetails;
    const RequestParamSchema: z.ZodObject<{
        shortUuid: z.ZodString;
    }, z.core.$strip>;
    const RequestQuerySchema: z.ZodObject<{
        withDisabledHosts: z.ZodPrefault<z.ZodOptional<z.ZodPipe<z.ZodString, z.ZodTransform<boolean, string>>>>;
    }, z.core.$strip>;
    const ResponseSchema: z.ZodObject<{
        response: z.ZodObject<{
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
                }, z.core.$strip>>;
            }, z.core.$strip>;
            headers: z.ZodRecord<z.ZodString, z.ZodOptional<z.ZodString>>;
            resolvedProxyConfigs: z.ZodArray<z.ZodObject<{
                finalRemark: z.ZodString;
                address: z.ZodString;
                port: z.ZodInt;
                protocol: z.ZodEnum<{
                    vless: "vless";
                    trojan: "trojan";
                    shadowsocks: "shadowsocks";
                    hysteria: "hysteria";
                }>;
                protocolOptions: z.ZodUnion<readonly [z.ZodObject<{
                    encryption: z.ZodString;
                    id: z.ZodString;
                    flow: z.ZodEnum<{
                        "": "";
                        "xtls-rprx-vision": "xtls-rprx-vision";
                        "xtls-rprx-vision-udp443": "xtls-rprx-vision-udp443";
                    }>;
                }, z.core.$strip>, z.ZodObject<{
                    password: z.ZodString;
                }, z.core.$strip>, z.ZodObject<{
                    method: z.ZodString;
                    password: z.ZodString;
                    uot: z.ZodBoolean;
                    uotVersion: z.ZodInt;
                }, z.core.$strip>, z.ZodObject<{
                    version: z.ZodInt;
                }, z.core.$strip>]>;
                transport: z.ZodEnum<{
                    hysteria: "hysteria";
                    tcp: "tcp";
                    xhttp: "xhttp";
                    ws: "ws";
                    httpupgrade: "httpupgrade";
                    grpc: "grpc";
                    kcp: "kcp";
                }>;
                transportOptions: z.ZodUnion<readonly [z.ZodObject<{
                    header: z.ZodNullable<z.ZodDiscriminatedUnion<[z.ZodObject<{
                        type: z.ZodLiteral<"none">;
                    }, z.core.$strip>, z.ZodObject<{
                        type: z.ZodLiteral<"http">;
                        request: z.ZodOptional<z.ZodObject<{
                            version: z.ZodOptional<z.ZodString>;
                            method: z.ZodOptional<z.ZodString>;
                            path: z.ZodOptional<z.ZodArray<z.ZodString>>;
                            headers: z.ZodOptional<z.ZodRecord<z.ZodString, z.ZodUnknown>>;
                        }, z.core.$strip>>;
                        response: z.ZodOptional<z.ZodObject<{
                            version: z.ZodOptional<z.ZodString>;
                            status: z.ZodOptional<z.ZodString>;
                            reason: z.ZodOptional<z.ZodString>;
                            headers: z.ZodOptional<z.ZodRecord<z.ZodString, z.ZodUnknown>>;
                        }, z.core.$strip>>;
                    }, z.core.$strip>], "type">>;
                }, z.core.$strip>, z.ZodObject<{
                    path: z.ZodNullable<z.ZodString>;
                    host: z.ZodNullable<z.ZodString>;
                    mode: z.ZodEnum<{
                        auto: "auto";
                        "packet-up": "packet-up";
                        "stream-up": "stream-up";
                        "stream-one": "stream-one";
                    }>;
                    extra: z.ZodNullable<z.ZodRecord<z.ZodString, z.ZodUnknown>>;
                }, z.core.$strip>, z.ZodObject<{
                    path: z.ZodNullable<z.ZodString>;
                    host: z.ZodNullable<z.ZodString>;
                    headers: z.ZodNullable<z.ZodRecord<z.ZodString, z.ZodString>>;
                    heartbeatPeriod: z.ZodNullable<z.ZodNumber>;
                }, z.core.$strip>, z.ZodObject<{
                    path: z.ZodNullable<z.ZodString>;
                    host: z.ZodNullable<z.ZodString>;
                    headers: z.ZodNullable<z.ZodRecord<z.ZodString, z.ZodString>>;
                }, z.core.$strip>, z.ZodObject<{
                    authority: z.ZodNullable<z.ZodString>;
                    serviceName: z.ZodNullable<z.ZodString>;
                    multiMode: z.ZodBoolean;
                }, z.core.$strip>, z.ZodObject<{
                    clientMtu: z.ZodInt;
                    clientTti: z.ZodInt;
                    congestion: z.ZodBoolean;
                }, z.core.$strip>, z.ZodObject<{
                    version: z.ZodInt;
                    auth: z.ZodString;
                }, z.core.$strip>]>;
                security: z.ZodEnum<{
                    none: "none";
                    tls: "tls";
                    reality: "reality";
                }>;
                securityOptions: z.ZodOptional<z.ZodUnion<readonly [z.ZodObject<{
                    pinnedPeerCertSha256: z.ZodNullable<z.ZodString>;
                    verifyPeerCertByName: z.ZodNullable<z.ZodString>;
                    alpn: z.ZodNullable<z.ZodString>;
                    enableSessionResumption: z.ZodBoolean;
                    fingerprint: z.ZodNullable<z.ZodString>;
                    serverName: z.ZodNullable<z.ZodString>;
                    echConfigList: z.ZodNullable<z.ZodString>;
                    echForceQuery: z.ZodNullable<z.ZodString>;
                    echSockopt: z.ZodNullable<z.ZodUnknown>;
                }, z.core.$strip>, z.ZodObject<{
                    fingerprint: z.ZodString;
                    publicKey: z.ZodString;
                    shortId: z.ZodNullable<z.ZodString>;
                    serverName: z.ZodString;
                    spiderX: z.ZodNullable<z.ZodString>;
                    mldsa65Verify: z.ZodNullable<z.ZodString>;
                }, z.core.$strip>]>>;
                streamOverrides: z.ZodObject<{
                    finalMask: z.ZodNullable<z.ZodUnknown>;
                    sockopt: z.ZodNullable<z.ZodUnknown>;
                }, z.core.$strip>;
                mux: z.ZodNullable<z.ZodUnknown>;
                clientOverrides: z.ZodObject<{
                    shuffleHost: z.ZodBoolean;
                    mihomoX25519: z.ZodBoolean;
                    mihomoIpVersion: z.ZodNullable<z.ZodEnum<{
                        readonly DUAL: "dual";
                        readonly IPV4: "ipv4";
                        readonly IPV6: "ipv6";
                        readonly IPV4_PREFER: "ipv4-prefer";
                        readonly IPV6_PREFER: "ipv6-prefer";
                    }>>;
                    serverDescription: z.ZodNullable<z.ZodString>;
                    xrayJsonTemplate: z.ZodNullable<z.ZodUnknown>;
                }, z.core.$strip>;
                metadata: z.ZodObject<{
                    uuid: z.ZodUUID;
                    tags: z.ZodArray<z.ZodString>;
                    excludeFromSubscriptionTypes: z.ZodArray<z.ZodEnum<{
                        readonly XRAY_JSON: "XRAY_JSON";
                        readonly XRAY_BASE64: "XRAY_BASE64";
                        readonly MIHOMO: "MIHOMO";
                        readonly STASH: "STASH";
                        readonly CLASH: "CLASH";
                        readonly SINGBOX: "SINGBOX";
                    }>>;
                    inboundTag: z.ZodString;
                    configProfileUuid: z.ZodNullable<z.ZodUUID>;
                    configProfileInboundUuid: z.ZodNullable<z.ZodUUID>;
                    isDisabled: z.ZodBoolean;
                    isHidden: z.ZodBoolean;
                    viewPosition: z.ZodInt;
                    remark: z.ZodString;
                    vlessRouteId: z.ZodNullable<z.ZodInt>;
                    rawInbound: z.ZodNullable<z.ZodUnknown>;
                }, z.core.$strip>;
            }, z.core.$strip>>;
        }, z.core.$strip>;
    }, z.core.$strip>;
    type RequestParam = z.infer<typeof RequestParamSchema>;
    type RequestQuery = z.infer<typeof RequestQuerySchema>;
    type Response = z.infer<typeof ResponseSchema>;
}
//# sourceMappingURL=get-raw-subscription-by-short-uuid.command.d.ts.map