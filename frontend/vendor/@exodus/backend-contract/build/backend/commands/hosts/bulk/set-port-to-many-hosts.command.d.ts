import { z } from 'zod';
export declare namespace SetPortToManyHostsCommand {
    const url: "/api/hosts/bulk/set-port";
    const TSQ_url: "/api/hosts/bulk/set-port";
    const endpointDetails: import("../../../constants").EndpointDetails;
    const RequestSchema: z.ZodObject<{
        uuids: z.ZodArray<z.ZodString, "many">;
        port: z.ZodNumber;
    }, "strip", z.ZodTypeAny, {
        port: number;
        uuids: string[];
    }, {
        port: number;
        uuids: string[];
    }>;
    type Request = z.infer<typeof RequestSchema>;
    const ResponseSchema: z.ZodObject<{
        response: z.ZodArray<z.ZodObject<{
            uuid: z.ZodString;
            viewPosition: z.ZodNumber;
            remark: z.ZodString;
            address: z.ZodString;
            port: z.ZodNumber;
            path: z.ZodNullable<z.ZodString>;
            sni: z.ZodNullable<z.ZodString>;
            host: z.ZodNullable<z.ZodString>;
            alpn: z.ZodNullable<z.ZodString>;
            fingerprint: z.ZodNullable<z.ZodString>;
            isDisabled: z.ZodDefault<z.ZodBoolean>;
            securityLayer: z.ZodDefault<z.ZodNativeEnum<{
                readonly DEFAULT: "DEFAULT";
                readonly TLS: "TLS";
                readonly NONE: "NONE";
            }>>;
            muxParams: z.ZodNullable<z.ZodUnknown>;
            sockoptParams: z.ZodNullable<z.ZodUnknown>;
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
            tag: z.ZodNullable<z.ZodString>;
            isHidden: z.ZodDefault<z.ZodBoolean>;
            overrideSniFromAddress: z.ZodDefault<z.ZodBoolean>;
            keepSniBlank: z.ZodDefault<z.ZodBoolean>;
            vlessRouteId: z.ZodNullable<z.ZodNumber>;
            allowInsecure: z.ZodDefault<z.ZodBoolean>;
            shuffleHost: z.ZodBoolean;
            mihomoX25519: z.ZodBoolean;
            nodes: z.ZodArray<z.ZodString, "many">;
            xrayJsonTemplateUuid: z.ZodNullable<z.ZodString>;
            excludedInternalSquads: z.ZodArray<z.ZodString, "many">;
            excludeFromSubscriptionTypes: z.ZodOptional<z.ZodArray<z.ZodNativeEnum<{
                readonly XRAY_JSON: "XRAY_JSON";
                readonly XRAY_BASE64: "XRAY_BASE64";
                readonly MIHOMO: "MIHOMO";
                readonly STASH: "STASH";
                readonly CLASH: "CLASH";
                readonly SINGBOX: "SINGBOX";
            }>, "many">>;
        }, "strip", z.ZodTypeAny, {
            nodes: string[];
            path: string | null;
            uuid: string;
            tag: string | null;
            port: number;
            viewPosition: number;
            remark: string;
            address: string;
            sni: string | null;
            host: string | null;
            alpn: string | null;
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
            allowInsecure: boolean;
            shuffleHost: boolean;
            mihomoX25519: boolean;
            xrayJsonTemplateUuid: string | null;
            excludedInternalSquads: string[];
            muxParams?: unknown;
            sockoptParams?: unknown;
            excludeFromSubscriptionTypes?: ("STASH" | "SINGBOX" | "MIHOMO" | "XRAY_JSON" | "CLASH" | "XRAY_BASE64")[] | undefined;
        }, {
            nodes: string[];
            path: string | null;
            uuid: string;
            tag: string | null;
            port: number;
            viewPosition: number;
            remark: string;
            address: string;
            sni: string | null;
            host: string | null;
            alpn: string | null;
            fingerprint: string | null;
            inbound: {
                configProfileUuid: string | null;
                configProfileInboundUuid: string | null;
            };
            serverDescription: string | null;
            vlessRouteId: number | null;
            shuffleHost: boolean;
            mihomoX25519: boolean;
            xrayJsonTemplateUuid: string | null;
            excludedInternalSquads: string[];
            isDisabled?: boolean | undefined;
            securityLayer?: "DEFAULT" | "TLS" | "NONE" | undefined;
            muxParams?: unknown;
            sockoptParams?: unknown;
            isHidden?: boolean | undefined;
            overrideSniFromAddress?: boolean | undefined;
            keepSniBlank?: boolean | undefined;
            allowInsecure?: boolean | undefined;
            excludeFromSubscriptionTypes?: ("STASH" | "SINGBOX" | "MIHOMO" | "XRAY_JSON" | "CLASH" | "XRAY_BASE64")[] | undefined;
        }>, "many">;
    }, "strip", z.ZodTypeAny, {
        response: {
            nodes: string[];
            path: string | null;
            uuid: string;
            tag: string | null;
            port: number;
            viewPosition: number;
            remark: string;
            address: string;
            sni: string | null;
            host: string | null;
            alpn: string | null;
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
            allowInsecure: boolean;
            shuffleHost: boolean;
            mihomoX25519: boolean;
            xrayJsonTemplateUuid: string | null;
            excludedInternalSquads: string[];
            muxParams?: unknown;
            sockoptParams?: unknown;
            excludeFromSubscriptionTypes?: ("STASH" | "SINGBOX" | "MIHOMO" | "XRAY_JSON" | "CLASH" | "XRAY_BASE64")[] | undefined;
        }[];
    }, {
        response: {
            nodes: string[];
            path: string | null;
            uuid: string;
            tag: string | null;
            port: number;
            viewPosition: number;
            remark: string;
            address: string;
            sni: string | null;
            host: string | null;
            alpn: string | null;
            fingerprint: string | null;
            inbound: {
                configProfileUuid: string | null;
                configProfileInboundUuid: string | null;
            };
            serverDescription: string | null;
            vlessRouteId: number | null;
            shuffleHost: boolean;
            mihomoX25519: boolean;
            xrayJsonTemplateUuid: string | null;
            excludedInternalSquads: string[];
            isDisabled?: boolean | undefined;
            securityLayer?: "DEFAULT" | "TLS" | "NONE" | undefined;
            muxParams?: unknown;
            sockoptParams?: unknown;
            isHidden?: boolean | undefined;
            overrideSniFromAddress?: boolean | undefined;
            keepSniBlank?: boolean | undefined;
            allowInsecure?: boolean | undefined;
            excludeFromSubscriptionTypes?: ("STASH" | "SINGBOX" | "MIHOMO" | "XRAY_JSON" | "CLASH" | "XRAY_BASE64")[] | undefined;
        }[];
    }>;
    type Response = z.infer<typeof ResponseSchema>;
}
//# sourceMappingURL=set-port-to-many-hosts.command.d.ts.map