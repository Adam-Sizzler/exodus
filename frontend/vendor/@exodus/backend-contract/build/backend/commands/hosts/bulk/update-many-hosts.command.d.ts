import { z } from 'zod';
export declare namespace UpdateManyHostsCommand {
    const url: "/api/hosts/bulk/update";
    const TSQ_url: "/api/hosts/bulk/update";
    const endpointDetails: import("../../../constants").EndpointDetails;
    const RequestBodySchema: z.ZodObject<{
        nodes: z.ZodOptional<z.ZodOptional<z.ZodArray<z.ZodUUID>>>;
        tags: z.ZodOptional<z.ZodOptional<z.ZodArray<z.ZodString>>>;
        path: z.ZodOptional<z.ZodOptional<z.ZodNullable<z.ZodString>>>;
        port: z.ZodOptional<z.ZodOptional<z.ZodInt>>;
        remark: z.ZodOptional<z.ZodOptional<z.ZodString>>;
        address: z.ZodOptional<z.ZodOptional<z.ZodString>>;
        sni: z.ZodOptional<z.ZodOptional<z.ZodNullable<z.ZodString>>>;
        host: z.ZodOptional<z.ZodOptional<z.ZodNullable<z.ZodString>>>;
        alpn: z.ZodOptional<z.ZodOptional<z.ZodNullable<z.ZodEnum<{
            readonly H3: "h3";
            readonly H2: "h2";
            readonly HTTP_1_1: "http/1.1";
            readonly H_COMBINED: "h2,http/1.1";
            readonly H3_H2_H1_COMBINED: "h3,h2,http/1.1";
            readonly H3_H2_COMBINED: "h3,h2";
        }>>>>;
        fingerprint: z.ZodOptional<z.ZodOptional<z.ZodNullable<z.ZodString>>>;
        isDisabled: z.ZodOptional<z.ZodDefault<z.ZodBoolean>>;
        securityLayer: z.ZodOptional<z.ZodOptional<z.ZodEnum<{
            readonly DEFAULT: "DEFAULT";
            readonly TLS: "TLS";
            readonly NONE: "NONE";
        }>>>;
        xhttpExtraParams: z.ZodOptional<z.ZodOptional<z.ZodNullable<z.ZodUnknown>>>;
        muxParams: z.ZodOptional<z.ZodOptional<z.ZodNullable<z.ZodUnknown>>>;
        singboxMuxParams: z.ZodOptional<z.ZodOptional<z.ZodNullable<z.ZodUnknown>>>;
        clashMuxParams: z.ZodOptional<z.ZodOptional<z.ZodNullable<z.ZodString>>>;
        singboxCustomParams: z.ZodOptional<z.ZodOptional<z.ZodNullable<z.ZodUnknown>>>;
        mihomoCustomParams: z.ZodOptional<z.ZodOptional<z.ZodNullable<z.ZodString>>>;
        sockoptParams: z.ZodOptional<z.ZodOptional<z.ZodNullable<z.ZodUnknown>>>;
        finalMask: z.ZodOptional<z.ZodOptional<z.ZodNullable<z.ZodUnknown>>>;
        inbound: z.ZodOptional<z.ZodOptional<z.ZodObject<{
            configProfileUuid: z.ZodUUID;
            configProfileInboundUuid: z.ZodUUID;
        }, z.core.$strip>>>;
        serverDescription: z.ZodOptional<z.ZodOptional<z.ZodNullable<z.ZodString>>>;
        isHidden: z.ZodOptional<z.ZodOptional<z.ZodBoolean>>;
        overrideSniFromAddress: z.ZodOptional<z.ZodOptional<z.ZodBoolean>>;
        keepSniBlank: z.ZodOptional<z.ZodOptional<z.ZodBoolean>>;
        overrideProtocolCredential: z.ZodOptional<z.ZodOptional<z.ZodBoolean>>;
        protocolCredential: z.ZodOptional<z.ZodOptional<z.ZodNullable<z.ZodString>>>;
        vlessRouteId: z.ZodOptional<z.ZodOptional<z.ZodNullable<z.ZodInt>>>;
        pinnedPeerCertSha256: z.ZodOptional<z.ZodOptional<z.ZodNullable<z.ZodString>>>;
        verifyPeerCertByName: z.ZodOptional<z.ZodOptional<z.ZodNullable<z.ZodString>>>;
        shuffleHost: z.ZodOptional<z.ZodOptional<z.ZodBoolean>>;
        mihomoX25519: z.ZodOptional<z.ZodOptional<z.ZodBoolean>>;
        mihomoIpVersion: z.ZodOptional<z.ZodOptional<z.ZodNullable<z.ZodEnum<{
            readonly DUAL: "dual";
            readonly IPV4: "ipv4";
            readonly IPV6: "ipv6";
            readonly IPV4_PREFER: "ipv4-prefer";
            readonly IPV6_PREFER: "ipv6-prefer";
        }>>>>;
        xrayJsonTemplateUuid: z.ZodOptional<z.ZodOptional<z.ZodNullable<z.ZodUUID>>>;
        excludedInternalSquads: z.ZodOptional<z.ZodOptional<z.ZodArray<z.ZodUUID>>>;
        excludeFromSubscriptionTypes: z.ZodOptional<z.ZodOptional<z.ZodArray<z.ZodEnum<{
            readonly XRAY_JSON: "XRAY_JSON";
            readonly XRAY_BASE64: "XRAY_BASE64";
            readonly MIHOMO: "MIHOMO";
            readonly STASH: "STASH";
            readonly CLASH: "CLASH";
            readonly SINGBOX: "SINGBOX";
        }>>>>;
        uuids: z.ZodArray<z.ZodUUID>;
    }, z.core.$strip>;
    type RequestBody = z.infer<typeof RequestBodySchema>;
}
//# sourceMappingURL=update-many-hosts.command.d.ts.map