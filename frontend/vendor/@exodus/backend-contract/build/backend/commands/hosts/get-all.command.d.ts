import { z } from 'zod';
export declare namespace GetAllHostsCommand {
    const url: "/api/hosts/";
    const TSQ_url: "/api/hosts/";
    const endpointDetails: import("../../constants").EndpointDetails;
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
            xHttpExtraParams: z.ZodNullable<z.ZodUnknown>;
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
            xHttpExtraParams?: unknown;
            muxParams?: unknown;
            sockoptParams?: unknown;
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
            xHttpExtraParams?: unknown;
            muxParams?: unknown;
            sockoptParams?: unknown;
            isHidden?: boolean | undefined;
            overrideSniFromAddress?: boolean | undefined;
            keepSniBlank?: boolean | undefined;
            allowInsecure?: boolean | undefined;
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
            xHttpExtraParams?: unknown;
            muxParams?: unknown;
            sockoptParams?: unknown;
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
            xHttpExtraParams?: unknown;
            muxParams?: unknown;
            sockoptParams?: unknown;
            isHidden?: boolean | undefined;
            overrideSniFromAddress?: boolean | undefined;
            keepSniBlank?: boolean | undefined;
            allowInsecure?: boolean | undefined;
        }[];
    }>;
    type Response = z.infer<typeof ResponseSchema>;
}
//# sourceMappingURL=get-all.command.d.ts.map