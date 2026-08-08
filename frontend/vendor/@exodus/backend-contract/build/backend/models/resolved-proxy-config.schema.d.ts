import { z } from 'zod';
export declare const VlessProtocolOptionsSchema: z.ZodObject<{
    encryption: z.ZodString;
    id: z.ZodString;
    flow: z.ZodEnum<{
        "": "";
        "xtls-rprx-vision": "xtls-rprx-vision";
        "xtls-rprx-vision-udp443": "xtls-rprx-vision-udp443";
    }>;
}, z.core.$strip>;
export declare const TrojanProtocolOptionsSchema: z.ZodObject<{
    password: z.ZodString;
}, z.core.$strip>;
export declare const ShadowsocksProtocolOptionsSchema: z.ZodObject<{
    method: z.ZodString;
    password: z.ZodString;
    uot: z.ZodBoolean;
    uotVersion: z.ZodInt;
}, z.core.$strip>;
export declare const TcpTransportOptionsSchema: z.ZodObject<{
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
}, z.core.$strip>;
export declare const XhttpTransportOptionsSchema: z.ZodObject<{
    path: z.ZodNullable<z.ZodString>;
    host: z.ZodNullable<z.ZodString>;
    mode: z.ZodEnum<{
        auto: "auto";
        "packet-up": "packet-up";
        "stream-up": "stream-up";
        "stream-one": "stream-one";
    }>;
    extra: z.ZodNullable<z.ZodRecord<z.ZodString, z.ZodUnknown>>;
}, z.core.$strip>;
export declare const WsTransportOptionsSchema: z.ZodObject<{
    path: z.ZodNullable<z.ZodString>;
    host: z.ZodNullable<z.ZodString>;
    headers: z.ZodNullable<z.ZodRecord<z.ZodString, z.ZodString>>;
    heartbeatPeriod: z.ZodNullable<z.ZodNumber>;
}, z.core.$strip>;
export declare const HttpUpgradeTransportOptionsSchema: z.ZodObject<{
    path: z.ZodNullable<z.ZodString>;
    host: z.ZodNullable<z.ZodString>;
    headers: z.ZodNullable<z.ZodRecord<z.ZodString, z.ZodString>>;
}, z.core.$strip>;
export declare const GrpcTransportOptionsSchema: z.ZodObject<{
    authority: z.ZodNullable<z.ZodString>;
    serviceName: z.ZodNullable<z.ZodString>;
    multiMode: z.ZodBoolean;
}, z.core.$strip>;
export declare const KcpTransportOptionsSchema: z.ZodObject<{
    clientMtu: z.ZodInt;
    clientTti: z.ZodInt;
    congestion: z.ZodBoolean;
}, z.core.$strip>;
export declare const HysteriaProtocolOptionsSchema: z.ZodObject<{
    version: z.ZodInt;
}, z.core.$strip>;
export declare const HysteriaTransportOptionsSchema: z.ZodObject<{
    version: z.ZodInt;
    auth: z.ZodString;
}, z.core.$strip>;
export declare const TlsSecurityOptionsSchema: z.ZodObject<{
    pinnedPeerCertSha256: z.ZodNullable<z.ZodString>;
    verifyPeerCertByName: z.ZodNullable<z.ZodString>;
    alpn: z.ZodNullable<z.ZodString>;
    enableSessionResumption: z.ZodBoolean;
    fingerprint: z.ZodNullable<z.ZodString>;
    serverName: z.ZodNullable<z.ZodString>;
    echConfigList: z.ZodNullable<z.ZodString>;
    echForceQuery: z.ZodNullable<z.ZodString>;
    echSockopt: z.ZodNullable<z.ZodUnknown>;
}, z.core.$strip>;
export declare const RealitySecurityOptionsSchema: z.ZodObject<{
    fingerprint: z.ZodString;
    publicKey: z.ZodString;
    shortId: z.ZodNullable<z.ZodString>;
    serverName: z.ZodString;
    spiderX: z.ZodNullable<z.ZodString>;
    mldsa65Verify: z.ZodNullable<z.ZodString>;
}, z.core.$strip>;
export declare const ProtocolVariantSchema: z.ZodDiscriminatedUnion<[z.ZodObject<{
    protocol: z.ZodLiteral<"vless">;
    protocolOptions: z.ZodObject<{
        encryption: z.ZodString;
        id: z.ZodString;
        flow: z.ZodEnum<{
            "": "";
            "xtls-rprx-vision": "xtls-rprx-vision";
            "xtls-rprx-vision-udp443": "xtls-rprx-vision-udp443";
        }>;
    }, z.core.$strip>;
}, z.core.$strip>, z.ZodObject<{
    protocol: z.ZodLiteral<"trojan">;
    protocolOptions: z.ZodObject<{
        password: z.ZodString;
    }, z.core.$strip>;
}, z.core.$strip>, z.ZodObject<{
    protocol: z.ZodLiteral<"shadowsocks">;
    protocolOptions: z.ZodObject<{
        method: z.ZodString;
        password: z.ZodString;
        uot: z.ZodBoolean;
        uotVersion: z.ZodInt;
    }, z.core.$strip>;
}, z.core.$strip>, z.ZodObject<{
    protocol: z.ZodLiteral<"hysteria">;
    protocolOptions: z.ZodObject<{
        version: z.ZodInt;
    }, z.core.$strip>;
}, z.core.$strip>], "protocol">;
export declare const TransportVariantSchema: z.ZodDiscriminatedUnion<[z.ZodObject<{
    transport: z.ZodLiteral<"tcp">;
    transportOptions: z.ZodObject<{
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
    }, z.core.$strip>;
}, z.core.$strip>, z.ZodObject<{
    transport: z.ZodLiteral<"xhttp">;
    transportOptions: z.ZodObject<{
        path: z.ZodNullable<z.ZodString>;
        host: z.ZodNullable<z.ZodString>;
        mode: z.ZodEnum<{
            auto: "auto";
            "packet-up": "packet-up";
            "stream-up": "stream-up";
            "stream-one": "stream-one";
        }>;
        extra: z.ZodNullable<z.ZodRecord<z.ZodString, z.ZodUnknown>>;
    }, z.core.$strip>;
}, z.core.$strip>, z.ZodObject<{
    transport: z.ZodLiteral<"ws">;
    transportOptions: z.ZodObject<{
        path: z.ZodNullable<z.ZodString>;
        host: z.ZodNullable<z.ZodString>;
        headers: z.ZodNullable<z.ZodRecord<z.ZodString, z.ZodString>>;
        heartbeatPeriod: z.ZodNullable<z.ZodNumber>;
    }, z.core.$strip>;
}, z.core.$strip>, z.ZodObject<{
    transport: z.ZodLiteral<"httpupgrade">;
    transportOptions: z.ZodObject<{
        path: z.ZodNullable<z.ZodString>;
        host: z.ZodNullable<z.ZodString>;
        headers: z.ZodNullable<z.ZodRecord<z.ZodString, z.ZodString>>;
    }, z.core.$strip>;
}, z.core.$strip>, z.ZodObject<{
    transport: z.ZodLiteral<"grpc">;
    transportOptions: z.ZodObject<{
        authority: z.ZodNullable<z.ZodString>;
        serviceName: z.ZodNullable<z.ZodString>;
        multiMode: z.ZodBoolean;
    }, z.core.$strip>;
}, z.core.$strip>, z.ZodObject<{
    transport: z.ZodLiteral<"kcp">;
    transportOptions: z.ZodObject<{
        clientMtu: z.ZodInt;
        clientTti: z.ZodInt;
        congestion: z.ZodBoolean;
    }, z.core.$strip>;
}, z.core.$strip>, z.ZodObject<{
    transport: z.ZodLiteral<"hysteria">;
    transportOptions: z.ZodObject<{
        version: z.ZodInt;
        auth: z.ZodString;
    }, z.core.$strip>;
}, z.core.$strip>], "transport">;
export declare const SecurityVariantSchema: z.ZodDiscriminatedUnion<[z.ZodObject<{
    security: z.ZodLiteral<"tls">;
    securityOptions: z.ZodObject<{
        pinnedPeerCertSha256: z.ZodNullable<z.ZodString>;
        verifyPeerCertByName: z.ZodNullable<z.ZodString>;
        alpn: z.ZodNullable<z.ZodString>;
        enableSessionResumption: z.ZodBoolean;
        fingerprint: z.ZodNullable<z.ZodString>;
        serverName: z.ZodNullable<z.ZodString>;
        echConfigList: z.ZodNullable<z.ZodString>;
        echForceQuery: z.ZodNullable<z.ZodString>;
        echSockopt: z.ZodNullable<z.ZodUnknown>;
    }, z.core.$strip>;
}, z.core.$strip>, z.ZodObject<{
    security: z.ZodLiteral<"reality">;
    securityOptions: z.ZodObject<{
        fingerprint: z.ZodString;
        publicKey: z.ZodString;
        shortId: z.ZodNullable<z.ZodString>;
        serverName: z.ZodString;
        spiderX: z.ZodNullable<z.ZodString>;
        mldsa65Verify: z.ZodNullable<z.ZodString>;
    }, z.core.$strip>;
}, z.core.$strip>, z.ZodObject<{
    security: z.ZodLiteral<"none">;
}, z.core.$strip>], "security">;
export declare const ProxyEntryMetadataSchema: z.ZodObject<{
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
export declare const ResolvedProxyConfigSchema: z.ZodObject<{
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
}, z.core.$strip>;
//# sourceMappingURL=resolved-proxy-config.schema.d.ts.map