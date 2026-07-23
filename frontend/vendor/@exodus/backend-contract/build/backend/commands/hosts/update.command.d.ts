import { z } from 'zod';
export declare namespace UpdateHostCommand {
    const url: "/api/hosts/";
    const TSQ_url: "/api/hosts/";
    const endpointDetails: import("../../constants").EndpointDetails;
    const RequestSchema: z.ZodObject<Pick<{
        uuid: z.ZodString;
        viewPosition: z.ZodNumber;
        remark: z.ZodString;
        address: z.ZodString;
        port: z.ZodNumber;
        path: z.ZodNullable<z.ZodString>;
        sni: z.ZodNullable<z.ZodString>;
        host: z.ZodNullable<z.ZodString>;
        alpn: z.ZodNullable<z.ZodNativeEnum<{
            readonly H3: "h3";
            readonly H2: "h2";
            readonly HTTP_1_1: "http/1.1";
            readonly H_COMBINED: "h2,http/1.1";
            readonly H3_H2_H1_COMBINED: "h3,h2,http/1.1";
            readonly H3_H2_COMBINED: "h3,h2";
        }>>;
        fingerprint: z.ZodNullable<z.ZodString>;
        isDisabled: z.ZodBoolean;
        securityLayer: z.ZodDefault<z.ZodNativeEnum<{
            readonly DEFAULT: "DEFAULT";
            readonly TLS: "TLS";
            readonly NONE: "NONE";
        }>>;
        xhttpExtraParams: z.ZodNullable<z.ZodUnknown>;
        muxParams: z.ZodNullable<z.ZodUnknown>;
        sockoptParams: z.ZodNullable<z.ZodUnknown>;
        finalMask: z.ZodNullable<z.ZodUnknown>;
        inbound: z.ZodObject<{
            configProfileUuid: z.ZodNullable<z.ZodString>;
            configProfileInboundUuid: z.ZodNullable<z.ZodString>;
        }, "strip", z.ZodTypeAny, {
            configProfileUuid: string | null;
            configProfileInboundUuid: string | null;
        }, {
            configProfileUuid: string | null;
            configProfileInboundUuid: string | null;
        }>;
        serverDescription: z.ZodNullable<z.ZodString>;
        tags: z.ZodDefault<z.ZodArray<z.ZodString, "many">>;
        isHidden: z.ZodDefault<z.ZodBoolean>;
        overrideSniFromAddress: z.ZodDefault<z.ZodBoolean>;
        keepSniBlank: z.ZodDefault<z.ZodBoolean>;
        vlessRouteId: z.ZodNullable<z.ZodNumber>;
        pinnedPeerCertSha256: z.ZodNullable<z.ZodString>;
        verifyPeerCertByName: z.ZodNullable<z.ZodString>;
        shuffleHost: z.ZodBoolean;
        mihomoX25519: z.ZodBoolean;
        mihomoIpVersion: z.ZodNullable<z.ZodNativeEnum<{
            readonly DUAL: "dual";
            readonly IPV4: "ipv4";
            readonly IPV6: "ipv6";
            readonly IPV4_PREFER: "ipv4-prefer";
            readonly IPV6_PREFER: "ipv6-prefer";
        }>>;
        nodes: z.ZodArray<z.ZodString, "many">;
        xrayJsonTemplateUuid: z.ZodNullable<z.ZodString>;
        excludedInternalSquads: z.ZodArray<z.ZodString, "many">;
        excludeFromSubscriptionTypes: z.ZodArray<z.ZodNativeEnum<{
            readonly XRAY_JSON: "XRAY_JSON";
            readonly XRAY_BASE64: "XRAY_BASE64";
            readonly MIHOMO: "MIHOMO";
            readonly STASH: "STASH";
            readonly CLASH: "CLASH";
            readonly SINGBOX: "SINGBOX";
        }>, "many">;
    }, "uuid"> & {
        inbound: z.ZodOptional<z.ZodObject<{
            configProfileUuid: z.ZodString;
            configProfileInboundUuid: z.ZodString;
        }, "strip", z.ZodTypeAny, {
            configProfileUuid: string;
            configProfileInboundUuid: string;
        }, {
            configProfileUuid: string;
            configProfileInboundUuid: string;
        }>>;
        remark: z.ZodOptional<z.ZodString>;
        address: z.ZodOptional<z.ZodString>;
        port: z.ZodOptional<z.ZodNumber>;
        path: z.ZodOptional<z.ZodNullable<z.ZodString>>;
        sni: z.ZodOptional<z.ZodNullable<z.ZodString>>;
        host: z.ZodOptional<z.ZodNullable<z.ZodString>>;
        alpn: z.ZodOptional<z.ZodNullable<z.ZodNativeEnum<{
            readonly H3: "h3";
            readonly H2: "h2";
            readonly HTTP_1_1: "http/1.1";
            readonly H_COMBINED: "h2,http/1.1";
            readonly H3_H2_H1_COMBINED: "h3,h2,http/1.1";
            readonly H3_H2_COMBINED: "h3,h2";
        }>>>;
        fingerprint: z.ZodOptional<z.ZodNullable<z.ZodString>>;
        isDisabled: z.ZodDefault<z.ZodBoolean>;
        securityLayer: z.ZodOptional<z.ZodNativeEnum<{
            readonly DEFAULT: "DEFAULT";
            readonly TLS: "TLS";
            readonly NONE: "NONE";
        }>>;
        xhttpExtraParams: z.ZodOptional<z.ZodNullable<z.ZodUnknown>>;
        muxParams: z.ZodOptional<z.ZodNullable<z.ZodUnknown>>;
        sockoptParams: z.ZodOptional<z.ZodNullable<z.ZodUnknown>>;
        finalMask: z.ZodOptional<z.ZodNullable<z.ZodUnknown>>;
        serverDescription: z.ZodOptional<z.ZodNullable<z.ZodString>>;
        tags: z.ZodOptional<z.ZodArray<z.ZodString, "many">>;
        isHidden: z.ZodOptional<z.ZodBoolean>;
        overrideSniFromAddress: z.ZodOptional<z.ZodBoolean>;
        keepSniBlank: z.ZodOptional<z.ZodBoolean>;
        vlessRouteId: z.ZodOptional<z.ZodNullable<z.ZodNumber>>;
        pinnedPeerCertSha256: z.ZodOptional<z.ZodNullable<z.ZodString>>;
        verifyPeerCertByName: z.ZodOptional<z.ZodNullable<z.ZodString>>;
        shuffleHost: z.ZodOptional<z.ZodBoolean>;
        mihomoX25519: z.ZodOptional<z.ZodBoolean>;
        mihomoIpVersion: z.ZodOptional<z.ZodNullable<z.ZodNativeEnum<{
            readonly DUAL: "dual";
            readonly IPV4: "ipv4";
            readonly IPV6: "ipv6";
            readonly IPV4_PREFER: "ipv4-prefer";
            readonly IPV6_PREFER: "ipv6-prefer";
        }>>>;
        nodes: z.ZodOptional<z.ZodArray<z.ZodString, "many">>;
        xrayJsonTemplateUuid: z.ZodOptional<z.ZodNullable<z.ZodString>>;
        excludedInternalSquads: z.ZodOptional<z.ZodArray<z.ZodString, "many">>;
        excludeFromSubscriptionTypes: z.ZodOptional<z.ZodArray<z.ZodNativeEnum<{
            readonly XRAY_JSON: "XRAY_JSON";
            readonly XRAY_BASE64: "XRAY_BASE64";
            readonly MIHOMO: "MIHOMO";
            readonly STASH: "STASH";
            readonly CLASH: "CLASH";
            readonly SINGBOX: "SINGBOX";
        }>, "many">>;
    }, "strip", z.ZodTypeAny, {
        uuid: string;
        isDisabled: boolean;
        nodes?: string[] | undefined;
        tags?: string[] | undefined;
        path?: string | null | undefined;
        port?: number | undefined;
        remark?: string | undefined;
        address?: string | undefined;
        sni?: string | null | undefined;
        host?: string | null | undefined;
        alpn?: "h3" | "h2" | "http/1.1" | "h2,http/1.1" | "h3,h2,http/1.1" | "h3,h2" | null | undefined;
        fingerprint?: string | null | undefined;
        securityLayer?: "DEFAULT" | "TLS" | "NONE" | undefined;
        xhttpExtraParams?: unknown;
        muxParams?: unknown;
        sockoptParams?: unknown;
        finalMask?: unknown;
        inbound?: {
            configProfileUuid: string;
            configProfileInboundUuid: string;
        } | undefined;
        serverDescription?: string | null | undefined;
        isHidden?: boolean | undefined;
        overrideSniFromAddress?: boolean | undefined;
        keepSniBlank?: boolean | undefined;
        vlessRouteId?: number | null | undefined;
        pinnedPeerCertSha256?: string | null | undefined;
        verifyPeerCertByName?: string | null | undefined;
        shuffleHost?: boolean | undefined;
        mihomoX25519?: boolean | undefined;
        mihomoIpVersion?: "dual" | "ipv4" | "ipv6" | "ipv4-prefer" | "ipv6-prefer" | null | undefined;
        xrayJsonTemplateUuid?: string | null | undefined;
        excludedInternalSquads?: string[] | undefined;
        excludeFromSubscriptionTypes?: ("STASH" | "SINGBOX" | "MIHOMO" | "XRAY_JSON" | "CLASH" | "XRAY_BASE64")[] | undefined;
    }, {
        uuid: string;
        nodes?: string[] | undefined;
        tags?: string[] | undefined;
        path?: string | null | undefined;
        port?: number | undefined;
        remark?: string | undefined;
        address?: string | undefined;
        sni?: string | null | undefined;
        host?: string | null | undefined;
        alpn?: "h3" | "h2" | "http/1.1" | "h2,http/1.1" | "h3,h2,http/1.1" | "h3,h2" | null | undefined;
        fingerprint?: string | null | undefined;
        isDisabled?: boolean | undefined;
        securityLayer?: "DEFAULT" | "TLS" | "NONE" | undefined;
        xhttpExtraParams?: unknown;
        muxParams?: unknown;
        sockoptParams?: unknown;
        finalMask?: unknown;
        inbound?: {
            configProfileUuid: string;
            configProfileInboundUuid: string;
        } | undefined;
        serverDescription?: string | null | undefined;
        isHidden?: boolean | undefined;
        overrideSniFromAddress?: boolean | undefined;
        keepSniBlank?: boolean | undefined;
        vlessRouteId?: number | null | undefined;
        pinnedPeerCertSha256?: string | null | undefined;
        verifyPeerCertByName?: string | null | undefined;
        shuffleHost?: boolean | undefined;
        mihomoX25519?: boolean | undefined;
        mihomoIpVersion?: "dual" | "ipv4" | "ipv6" | "ipv4-prefer" | "ipv6-prefer" | null | undefined;
        xrayJsonTemplateUuid?: string | null | undefined;
        excludedInternalSquads?: string[] | undefined;
        excludeFromSubscriptionTypes?: ("STASH" | "SINGBOX" | "MIHOMO" | "XRAY_JSON" | "CLASH" | "XRAY_BASE64")[] | undefined;
    }>;
    type Request = z.infer<typeof RequestSchema>;
    const ResponseSchema: z.ZodObject<{
        response: z.ZodObject<{
            uuid: z.ZodString;
            viewPosition: z.ZodNumber;
            remark: z.ZodString;
            address: z.ZodString;
            port: z.ZodNumber;
            path: z.ZodNullable<z.ZodString>;
            sni: z.ZodNullable<z.ZodString>;
            host: z.ZodNullable<z.ZodString>;
            alpn: z.ZodNullable<z.ZodNativeEnum<{
                readonly H3: "h3";
                readonly H2: "h2";
                readonly HTTP_1_1: "http/1.1";
                readonly H_COMBINED: "h2,http/1.1";
                readonly H3_H2_H1_COMBINED: "h3,h2,http/1.1";
                readonly H3_H2_COMBINED: "h3,h2";
            }>>;
            fingerprint: z.ZodNullable<z.ZodString>;
            isDisabled: z.ZodBoolean;
            securityLayer: z.ZodDefault<z.ZodNativeEnum<{
                readonly DEFAULT: "DEFAULT";
                readonly TLS: "TLS";
                readonly NONE: "NONE";
            }>>;
            xhttpExtraParams: z.ZodNullable<z.ZodUnknown>;
            muxParams: z.ZodNullable<z.ZodUnknown>;
            sockoptParams: z.ZodNullable<z.ZodUnknown>;
            finalMask: z.ZodNullable<z.ZodUnknown>;
            inbound: z.ZodObject<{
                configProfileUuid: z.ZodNullable<z.ZodString>;
                configProfileInboundUuid: z.ZodNullable<z.ZodString>;
            }, "strip", z.ZodTypeAny, {
                configProfileUuid: string | null;
                configProfileInboundUuid: string | null;
            }, {
                configProfileUuid: string | null;
                configProfileInboundUuid: string | null;
            }>;
            serverDescription: z.ZodNullable<z.ZodString>;
            tags: z.ZodDefault<z.ZodArray<z.ZodString, "many">>;
            isHidden: z.ZodDefault<z.ZodBoolean>;
            overrideSniFromAddress: z.ZodDefault<z.ZodBoolean>;
            keepSniBlank: z.ZodDefault<z.ZodBoolean>;
            vlessRouteId: z.ZodNullable<z.ZodNumber>;
            pinnedPeerCertSha256: z.ZodNullable<z.ZodString>;
            verifyPeerCertByName: z.ZodNullable<z.ZodString>;
            shuffleHost: z.ZodBoolean;
            mihomoX25519: z.ZodBoolean;
            mihomoIpVersion: z.ZodNullable<z.ZodNativeEnum<{
                readonly DUAL: "dual";
                readonly IPV4: "ipv4";
                readonly IPV6: "ipv6";
                readonly IPV4_PREFER: "ipv4-prefer";
                readonly IPV6_PREFER: "ipv6-prefer";
            }>>;
            nodes: z.ZodArray<z.ZodString, "many">;
            xrayJsonTemplateUuid: z.ZodNullable<z.ZodString>;
            excludedInternalSquads: z.ZodArray<z.ZodString, "many">;
            excludeFromSubscriptionTypes: z.ZodArray<z.ZodNativeEnum<{
                readonly XRAY_JSON: "XRAY_JSON";
                readonly XRAY_BASE64: "XRAY_BASE64";
                readonly MIHOMO: "MIHOMO";
                readonly STASH: "STASH";
                readonly CLASH: "CLASH";
                readonly SINGBOX: "SINGBOX";
            }>, "many">;
        }, "strip", z.ZodTypeAny, {
            nodes: string[];
            tags: string[];
            uuid: string;
            path: string | null;
            port: number;
            viewPosition: number;
            remark: string;
            address: string;
            sni: string | null;
            host: string | null;
            alpn: "h3" | "h2" | "http/1.1" | "h2,http/1.1" | "h3,h2,http/1.1" | "h3,h2" | null;
            fingerprint: string | null;
            isDisabled: boolean;
            securityLayer: "DEFAULT" | "TLS" | "NONE";
            inbound: {
                configProfileUuid: string | null;
                configProfileInboundUuid: string | null;
            };
            serverDescription: string | null;
            isHidden: boolean;
            overrideSniFromAddress: boolean;
            keepSniBlank: boolean;
            vlessRouteId: number | null;
            pinnedPeerCertSha256: string | null;
            verifyPeerCertByName: string | null;
            shuffleHost: boolean;
            mihomoX25519: boolean;
            mihomoIpVersion: "dual" | "ipv4" | "ipv6" | "ipv4-prefer" | "ipv6-prefer" | null;
            xrayJsonTemplateUuid: string | null;
            excludedInternalSquads: string[];
            excludeFromSubscriptionTypes: ("STASH" | "SINGBOX" | "MIHOMO" | "XRAY_JSON" | "CLASH" | "XRAY_BASE64")[];
            xhttpExtraParams?: unknown;
            muxParams?: unknown;
            sockoptParams?: unknown;
            finalMask?: unknown;
        }, {
            nodes: string[];
            uuid: string;
            path: string | null;
            port: number;
            viewPosition: number;
            remark: string;
            address: string;
            sni: string | null;
            host: string | null;
            alpn: "h3" | "h2" | "http/1.1" | "h2,http/1.1" | "h3,h2,http/1.1" | "h3,h2" | null;
            fingerprint: string | null;
            isDisabled: boolean;
            inbound: {
                configProfileUuid: string | null;
                configProfileInboundUuid: string | null;
            };
            serverDescription: string | null;
            vlessRouteId: number | null;
            pinnedPeerCertSha256: string | null;
            verifyPeerCertByName: string | null;
            shuffleHost: boolean;
            mihomoX25519: boolean;
            mihomoIpVersion: "dual" | "ipv4" | "ipv6" | "ipv4-prefer" | "ipv6-prefer" | null;
            xrayJsonTemplateUuid: string | null;
            excludedInternalSquads: string[];
            excludeFromSubscriptionTypes: ("STASH" | "SINGBOX" | "MIHOMO" | "XRAY_JSON" | "CLASH" | "XRAY_BASE64")[];
            tags?: string[] | undefined;
            securityLayer?: "DEFAULT" | "TLS" | "NONE" | undefined;
            xhttpExtraParams?: unknown;
            muxParams?: unknown;
            sockoptParams?: unknown;
            finalMask?: unknown;
            isHidden?: boolean | undefined;
            overrideSniFromAddress?: boolean | undefined;
            keepSniBlank?: boolean | undefined;
        }>;
    }, "strip", z.ZodTypeAny, {
        response: {
            nodes: string[];
            tags: string[];
            uuid: string;
            path: string | null;
            port: number;
            viewPosition: number;
            remark: string;
            address: string;
            sni: string | null;
            host: string | null;
            alpn: "h3" | "h2" | "http/1.1" | "h2,http/1.1" | "h3,h2,http/1.1" | "h3,h2" | null;
            fingerprint: string | null;
            isDisabled: boolean;
            securityLayer: "DEFAULT" | "TLS" | "NONE";
            inbound: {
                configProfileUuid: string | null;
                configProfileInboundUuid: string | null;
            };
            serverDescription: string | null;
            isHidden: boolean;
            overrideSniFromAddress: boolean;
            keepSniBlank: boolean;
            vlessRouteId: number | null;
            pinnedPeerCertSha256: string | null;
            verifyPeerCertByName: string | null;
            shuffleHost: boolean;
            mihomoX25519: boolean;
            mihomoIpVersion: "dual" | "ipv4" | "ipv6" | "ipv4-prefer" | "ipv6-prefer" | null;
            xrayJsonTemplateUuid: string | null;
            excludedInternalSquads: string[];
            excludeFromSubscriptionTypes: ("STASH" | "SINGBOX" | "MIHOMO" | "XRAY_JSON" | "CLASH" | "XRAY_BASE64")[];
            xhttpExtraParams?: unknown;
            muxParams?: unknown;
            sockoptParams?: unknown;
            finalMask?: unknown;
        };
    }, {
        response: {
            nodes: string[];
            uuid: string;
            path: string | null;
            port: number;
            viewPosition: number;
            remark: string;
            address: string;
            sni: string | null;
            host: string | null;
            alpn: "h3" | "h2" | "http/1.1" | "h2,http/1.1" | "h3,h2,http/1.1" | "h3,h2" | null;
            fingerprint: string | null;
            isDisabled: boolean;
            inbound: {
                configProfileUuid: string | null;
                configProfileInboundUuid: string | null;
            };
            serverDescription: string | null;
            vlessRouteId: number | null;
            pinnedPeerCertSha256: string | null;
            verifyPeerCertByName: string | null;
            shuffleHost: boolean;
            mihomoX25519: boolean;
            mihomoIpVersion: "dual" | "ipv4" | "ipv6" | "ipv4-prefer" | "ipv6-prefer" | null;
            xrayJsonTemplateUuid: string | null;
            excludedInternalSquads: string[];
            excludeFromSubscriptionTypes: ("STASH" | "SINGBOX" | "MIHOMO" | "XRAY_JSON" | "CLASH" | "XRAY_BASE64")[];
            tags?: string[] | undefined;
            securityLayer?: "DEFAULT" | "TLS" | "NONE" | undefined;
            xhttpExtraParams?: unknown;
            muxParams?: unknown;
            sockoptParams?: unknown;
            finalMask?: unknown;
            isHidden?: boolean | undefined;
            overrideSniFromAddress?: boolean | undefined;
            keepSniBlank?: boolean | undefined;
        };
    }>;
    type Response = z.infer<typeof ResponseSchema>;
}
//# sourceMappingURL=update.command.d.ts.map