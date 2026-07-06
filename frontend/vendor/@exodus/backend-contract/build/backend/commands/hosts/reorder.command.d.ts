import { z } from 'zod';
export declare namespace ReorderHostCommand {
    const url: "/api/hosts/actions/reorder";
    const TSQ_url: "/api/hosts/actions/reorder";
    const endpointDetails: import("../../constants").EndpointDetails;
    const RequestSchema: z.ZodObject<{
        hosts: z.ZodArray<z.ZodObject<Pick<{
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
        }, "uuid" | "viewPosition">, "strip", z.ZodTypeAny, {
            uuid: string;
            viewPosition: number;
        }, {
            uuid: string;
            viewPosition: number;
        }>, "many">;
    }, "strip", z.ZodTypeAny, {
        hosts: {
            uuid: string;
            viewPosition: number;
        }[];
    }, {
        hosts: {
            uuid: string;
            viewPosition: number;
        }[];
    }>;
    type Request = z.infer<typeof RequestSchema>;
    const ResponseSchema: z.ZodObject<{
        response: z.ZodObject<{
            isUpdated: z.ZodBoolean;
        }, "strip", z.ZodTypeAny, {
            isUpdated: boolean;
        }, {
            isUpdated: boolean;
        }>;
    }, "strip", z.ZodTypeAny, {
        response: {
            isUpdated: boolean;
        };
    }, {
        response: {
            isUpdated: boolean;
        };
    }>;
    type Response = z.infer<typeof ResponseSchema>;
}
//# sourceMappingURL=reorder.command.d.ts.map