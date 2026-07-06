import { z } from 'zod';
export declare const VlessProtocolOptionsSchema: z.ZodObject<{
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
}>;
export declare const TrojanProtocolOptionsSchema: z.ZodObject<{
    password: z.ZodString;
}, "strip", z.ZodTypeAny, {
    password: string;
}, {
    password: string;
}>;
export declare const ShadowsocksProtocolOptionsSchema: z.ZodObject<{
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
}>;
export declare const TcpTransportOptionsSchema: z.ZodObject<{
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
}>;
export declare const XhttpTransportOptionsSchema: z.ZodObject<{
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
}>;
export declare const WsTransportOptionsSchema: z.ZodObject<{
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
}>;
export declare const HttpUpgradeTransportOptionsSchema: z.ZodObject<{
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
}>;
export declare const GrpcTransportOptionsSchema: z.ZodObject<{
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
}>;
export declare const KcpTransportOptionsSchema: z.ZodObject<{
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
}>;
export declare const HysteriaProtocolOptionsSchema: z.ZodObject<{
    version: z.ZodNumber;
}, "strip", z.ZodTypeAny, {
    version: number;
}, {
    version: number;
}>;
export declare const HysteriaTransportOptionsSchema: z.ZodObject<{
    version: z.ZodNumber;
    auth: z.ZodString;
}, "strip", z.ZodTypeAny, {
    auth: string;
    version: number;
}, {
    auth: string;
    version: number;
}>;
export declare const TlsSecurityOptionsSchema: z.ZodObject<{
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
}>;
export declare const RealitySecurityOptionsSchema: z.ZodObject<{
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
}>;
export declare const ProtocolVariantSchema: z.ZodDiscriminatedUnion<"protocol", [z.ZodObject<{
    protocol: z.ZodLiteral<"vless">;
    protocolOptions: z.ZodObject<{
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
    }>;
}, "strip", z.ZodTypeAny, {
    protocol: "vless";
    protocolOptions: {
        id: string;
        encryption: string;
        flow: "" | "xtls-rprx-vision" | "xtls-rprx-vision-udp443";
    };
}, {
    protocol: "vless";
    protocolOptions: {
        id: string;
        encryption: string;
        flow: "" | "xtls-rprx-vision" | "xtls-rprx-vision-udp443";
    };
}>, z.ZodObject<{
    protocol: z.ZodLiteral<"trojan">;
    protocolOptions: z.ZodObject<{
        password: z.ZodString;
    }, "strip", z.ZodTypeAny, {
        password: string;
    }, {
        password: string;
    }>;
}, "strip", z.ZodTypeAny, {
    protocol: "trojan";
    protocolOptions: {
        password: string;
    };
}, {
    protocol: "trojan";
    protocolOptions: {
        password: string;
    };
}>, z.ZodObject<{
    protocol: z.ZodLiteral<"shadowsocks">;
    protocolOptions: z.ZodObject<{
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
    }>;
}, "strip", z.ZodTypeAny, {
    protocol: "shadowsocks";
    protocolOptions: {
        method: string;
        password: string;
        uot: boolean;
        uotVersion: number;
    };
}, {
    protocol: "shadowsocks";
    protocolOptions: {
        method: string;
        password: string;
        uot: boolean;
        uotVersion: number;
    };
}>, z.ZodObject<{
    protocol: z.ZodLiteral<"hysteria">;
    protocolOptions: z.ZodObject<{
        version: z.ZodNumber;
    }, "strip", z.ZodTypeAny, {
        version: number;
    }, {
        version: number;
    }>;
}, "strip", z.ZodTypeAny, {
    protocol: "hysteria";
    protocolOptions: {
        version: number;
    };
}, {
    protocol: "hysteria";
    protocolOptions: {
        version: number;
    };
}>]>;
export declare const TransportVariantSchema: z.ZodDiscriminatedUnion<"transport", [z.ZodObject<{
    transport: z.ZodLiteral<"tcp">;
    transportOptions: z.ZodObject<{
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
    }>;
}, "strip", z.ZodTypeAny, {
    transport: "tcp";
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
    };
}, {
    transport: "tcp";
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
    };
}>, z.ZodObject<{
    transport: z.ZodLiteral<"xhttp">;
    transportOptions: z.ZodObject<{
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
    }>;
}, "strip", z.ZodTypeAny, {
    transport: "xhttp";
    transportOptions: {
        path: string | null;
        host: string | null;
        mode: "auto" | "packet-up" | "stream-up" | "stream-one";
        extra: Record<string, unknown> | null;
    };
}, {
    transport: "xhttp";
    transportOptions: {
        path: string | null;
        host: string | null;
        mode: "auto" | "packet-up" | "stream-up" | "stream-one";
        extra: Record<string, unknown> | null;
    };
}>, z.ZodObject<{
    transport: z.ZodLiteral<"ws">;
    transportOptions: z.ZodObject<{
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
    }>;
}, "strip", z.ZodTypeAny, {
    transport: "ws";
    transportOptions: {
        path: string | null;
        host: string | null;
        headers: Record<string, string> | null;
        heartbeatPeriod: number | null;
    };
}, {
    transport: "ws";
    transportOptions: {
        path: string | null;
        host: string | null;
        headers: Record<string, string> | null;
        heartbeatPeriod: number | null;
    };
}>, z.ZodObject<{
    transport: z.ZodLiteral<"httpupgrade">;
    transportOptions: z.ZodObject<{
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
    }>;
}, "strip", z.ZodTypeAny, {
    transport: "httpupgrade";
    transportOptions: {
        path: string | null;
        host: string | null;
        headers: Record<string, string> | null;
    };
}, {
    transport: "httpupgrade";
    transportOptions: {
        path: string | null;
        host: string | null;
        headers: Record<string, string> | null;
    };
}>, z.ZodObject<{
    transport: z.ZodLiteral<"grpc">;
    transportOptions: z.ZodObject<{
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
    }>;
}, "strip", z.ZodTypeAny, {
    transport: "grpc";
    transportOptions: {
        authority: string | null;
        serviceName: string | null;
        multiMode: boolean;
    };
}, {
    transport: "grpc";
    transportOptions: {
        authority: string | null;
        serviceName: string | null;
        multiMode: boolean;
    };
}>, z.ZodObject<{
    transport: z.ZodLiteral<"kcp">;
    transportOptions: z.ZodObject<{
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
    }>;
}, "strip", z.ZodTypeAny, {
    transport: "kcp";
    transportOptions: {
        clientMtu: number;
        clientTti: number;
        congestion: boolean;
    };
}, {
    transport: "kcp";
    transportOptions: {
        clientMtu: number;
        clientTti: number;
        congestion: boolean;
    };
}>, z.ZodObject<{
    transport: z.ZodLiteral<"hysteria">;
    transportOptions: z.ZodObject<{
        version: z.ZodNumber;
        auth: z.ZodString;
    }, "strip", z.ZodTypeAny, {
        auth: string;
        version: number;
    }, {
        auth: string;
        version: number;
    }>;
}, "strip", z.ZodTypeAny, {
    transport: "hysteria";
    transportOptions: {
        auth: string;
        version: number;
    };
}, {
    transport: "hysteria";
    transportOptions: {
        auth: string;
        version: number;
    };
}>]>;
export declare const SecurityVariantSchema: z.ZodDiscriminatedUnion<"security", [z.ZodObject<{
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
    }>;
}, "strip", z.ZodTypeAny, {
    security: "tls";
    securityOptions: {
        alpn: string | null;
        fingerprint: string | null;
        pinnedPeerCertSha256: string | null;
        verifyPeerCertByName: string | null;
        enableSessionResumption: boolean;
        serverName: string | null;
        echConfigList: string | null;
        echForceQuery: string | null;
    };
}, {
    security: "tls";
    securityOptions: {
        alpn: string | null;
        fingerprint: string | null;
        pinnedPeerCertSha256: string | null;
        verifyPeerCertByName: string | null;
        enableSessionResumption: boolean;
        serverName: string | null;
        echConfigList: string | null;
        echForceQuery: string | null;
    };
}>, z.ZodObject<{
    security: z.ZodLiteral<"reality">;
    securityOptions: z.ZodObject<{
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
    }>;
}, "strip", z.ZodTypeAny, {
    security: "reality";
    securityOptions: {
        fingerprint: string;
        serverName: string;
        publicKey: string;
        shortId: string | null;
        spiderX: string | null;
        mldsa65Verify: string | null;
    };
}, {
    security: "reality";
    securityOptions: {
        fingerprint: string;
        serverName: string;
        publicKey: string;
        shortId: string | null;
        spiderX: string | null;
        mldsa65Verify: string | null;
    };
}>, z.ZodObject<{
    security: z.ZodLiteral<"none">;
}, "strip", z.ZodTypeAny, {
    security: "none";
}, {
    security: "none";
}>]>;
export declare const ProxyEntryMetadataSchema: z.ZodObject<{
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
export declare const ResolvedProxyConfigSchema: z.ZodObject<{
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
}>;
//# sourceMappingURL=resolved-proxy-config.schema.d.ts.map