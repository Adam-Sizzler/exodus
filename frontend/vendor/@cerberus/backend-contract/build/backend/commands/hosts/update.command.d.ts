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
        path: z.ZodOptional<z.ZodString>;
        sni: z.ZodOptional<z.ZodString>;
        host: z.ZodOptional<z.ZodString>;
        alpn: z.ZodOptional<z.ZodNullable<z.ZodNativeEnum<{
            readonly H3: "h3";
            readonly H2: "h2";
            readonly HTTP_1_1: "http/1.1";
            readonly H_COMBINED: "h2,http/1.1";
            readonly H3_H2_H1_COMBINED: "h3,h2,http/1.1";
            readonly H3_H2_COMBINED: "h3,h2";
        }>>>;
        fingerprint: z.ZodOptional<z.ZodNullable<z.ZodNativeEnum<{
            readonly CHROME: "chrome";
            readonly FIREFOX: "firefox";
            readonly SAFARI: "safari";
            readonly IOS: "ios";
            readonly ANDROID: "android";
            readonly EDGE: "edge";
            readonly QQ: "qq";
            readonly RANDOM: "random";
            readonly RANDOMIZED: "randomized";
        }>>>;
        isDisabled: z.ZodOptional<z.ZodBoolean>;
        securityLayer: z.ZodOptional<z.ZodNativeEnum<{
            readonly DEFAULT: "DEFAULT";
            readonly TLS: "TLS";
            readonly NONE: "NONE";
        }>>;
        xHttpExtraParams: z.ZodOptional<z.ZodNullable<z.ZodUnknown>>;
        muxParams: z.ZodOptional<z.ZodNullable<z.ZodUnknown>>;
        sockoptParams: z.ZodOptional<z.ZodNullable<z.ZodUnknown>>;
        serverDescription: z.ZodOptional<z.ZodNullable<z.ZodString>>;
        tag: z.ZodOptional<z.ZodNullable<z.ZodString>>;
        isHidden: z.ZodOptional<z.ZodBoolean>;
        overrideSniFromAddress: z.ZodOptional<z.ZodBoolean>;
        keepSniBlank: z.ZodOptional<z.ZodBoolean>;
        vlessRouteId: z.ZodOptional<z.ZodNullable<z.ZodNumber>>;
        allowInsecure: z.ZodOptional<z.ZodBoolean>;
        shuffleHost: z.ZodOptional<z.ZodBoolean>;
        mihomoX25519: z.ZodOptional<z.ZodBoolean>;
        nodes: z.ZodOptional<z.ZodArray<z.ZodString, "many">>;
        xrayJsonTemplateUuid: z.ZodOptional<z.ZodNullable<z.ZodString>>;
        excludedInternalSquads: z.ZodOptional<z.ZodArray<z.ZodString, "many">>;
    }, "strip", z.ZodTypeAny, {
        uuid: string;
        nodes?: string[] | undefined;
        path?: string | undefined;
        tag?: string | null | undefined;
        port?: number | undefined;
        remark?: string | undefined;
        address?: string | undefined;
        sni?: string | undefined;
        host?: string | undefined;
        alpn?: "h3" | "h2" | "http/1.1" | "h2,http/1.1" | "h3,h2,http/1.1" | "h3,h2" | null | undefined;
        fingerprint?: "chrome" | "firefox" | "safari" | "ios" | "android" | "edge" | "qq" | "random" | "randomized" | null | undefined;
        isDisabled?: boolean | undefined;
        securityLayer?: "DEFAULT" | "TLS" | "NONE" | undefined;
        xHttpExtraParams?: unknown;
        muxParams?: unknown;
        sockoptParams?: unknown;
        inbound?: {
            configProfileUuid: string;
            configProfileInboundUuid: string;
        } | undefined;
        serverDescription?: string | null | undefined;
        isHidden?: boolean | undefined;
        overrideSniFromAddress?: boolean | undefined;
        keepSniBlank?: boolean | undefined;
        vlessRouteId?: number | null | undefined;
        allowInsecure?: boolean | undefined;
        shuffleHost?: boolean | undefined;
        mihomoX25519?: boolean | undefined;
        xrayJsonTemplateUuid?: string | null | undefined;
        excludedInternalSquads?: string[] | undefined;
    }, {
        uuid: string;
        nodes?: string[] | undefined;
        path?: string | undefined;
        tag?: string | null | undefined;
        port?: number | undefined;
        remark?: string | undefined;
        address?: string | undefined;
        sni?: string | undefined;
        host?: string | undefined;
        alpn?: "h3" | "h2" | "http/1.1" | "h2,http/1.1" | "h3,h2,http/1.1" | "h3,h2" | null | undefined;
        fingerprint?: "chrome" | "firefox" | "safari" | "ios" | "android" | "edge" | "qq" | "random" | "randomized" | null | undefined;
        isDisabled?: boolean | undefined;
        securityLayer?: "DEFAULT" | "TLS" | "NONE" | undefined;
        xHttpExtraParams?: unknown;
        muxParams?: unknown;
        sockoptParams?: unknown;
        inbound?: {
            configProfileUuid: string;
            configProfileInboundUuid: string;
        } | undefined;
        serverDescription?: string | null | undefined;
        isHidden?: boolean | undefined;
        overrideSniFromAddress?: boolean | undefined;
        keepSniBlank?: boolean | undefined;
        vlessRouteId?: number | null | undefined;
        allowInsecure?: boolean | undefined;
        shuffleHost?: boolean | undefined;
        mihomoX25519?: boolean | undefined;
        xrayJsonTemplateUuid?: string | null | undefined;
        excludedInternalSquads?: string[] | undefined;
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
        }>;
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
        };
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
        };
    }>;
    type Response = z.infer<typeof ResponseSchema>;
}
//# sourceMappingURL=update.command.d.ts.map