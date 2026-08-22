import { z } from 'zod';
export declare namespace GetHostsCommand {
    const url: "/api/hosts/";
    const TSQ_url: "/api/hosts/";
    const endpointDetails: import("../../constants").EndpointDetails;
    const ResponseSchema: z.ZodObject<{
        response: z.ZodArray<z.ZodObject<{
            uuid: z.ZodUUID;
            viewPosition: z.ZodInt;
            remark: z.ZodString;
            address: z.ZodString;
            port: z.ZodInt;
            path: z.ZodNullable<z.ZodString>;
            sni: z.ZodNullable<z.ZodString>;
            host: z.ZodNullable<z.ZodString>;
            alpn: z.ZodNullable<z.ZodEnum<{
                readonly H3: "h3";
                readonly H2: "h2";
                readonly HTTP_1_1: "http/1.1";
                readonly H_COMBINED: "h2,http/1.1";
                readonly H3_H2_H1_COMBINED: "h3,h2,http/1.1";
                readonly H3_H2_COMBINED: "h3,h2";
            }>>;
            fingerprint: z.ZodNullable<z.ZodString>;
            isDisabled: z.ZodBoolean;
            securityLayer: z.ZodDefault<z.ZodEnum<{
                readonly DEFAULT: "DEFAULT";
                readonly TLS: "TLS";
                readonly NONE: "NONE";
            }>>;
            xhttpExtraParams: z.ZodNullable<z.ZodUnknown>;
            muxParams: z.ZodNullable<z.ZodUnknown>;
            sockoptParams: z.ZodNullable<z.ZodUnknown>;
            finalMask: z.ZodNullable<z.ZodUnknown>;
            inbound: z.ZodObject<{
                configProfileUuid: z.ZodNullable<z.ZodUUID>;
                configProfileInboundUuid: z.ZodNullable<z.ZodUUID>;
            }, z.core.$strip>;
            serverDescription: z.ZodNullable<z.ZodString>;
            tags: z.ZodDefault<z.ZodArray<z.ZodString>>;
            isHidden: z.ZodDefault<z.ZodBoolean>;
            overrideSniFromAddress: z.ZodDefault<z.ZodBoolean>;
            keepSniBlank: z.ZodDefault<z.ZodBoolean>;
            vlessRouteId: z.ZodNullable<z.ZodInt>;
            pinnedPeerCertSha256: z.ZodNullable<z.ZodString>;
            verifyPeerCertByName: z.ZodNullable<z.ZodString>;
            shuffleHost: z.ZodBoolean;
            mihomoX25519: z.ZodBoolean;
            mihomoIpVersion: z.ZodNullable<z.ZodEnum<{
                readonly DUAL: "dual";
                readonly IPV4: "ipv4";
                readonly IPV6: "ipv6";
                readonly IPV4_PREFER: "ipv4-prefer";
                readonly IPV6_PREFER: "ipv6-prefer";
            }>>;
            nodes: z.ZodArray<z.ZodUUID>;
            xrayJsonTemplateUuid: z.ZodNullable<z.ZodUUID>;
            excludedInternalSquads: z.ZodArray<z.ZodUUID>;
            excludeFromSubscriptionTypes: z.ZodArray<z.ZodEnum<{
                readonly XRAY_JSON: "XRAY_JSON";
                readonly XRAY_BASE64: "XRAY_BASE64";
                readonly MIHOMO: "MIHOMO";
                readonly STASH: "STASH";
                readonly CLASH: "CLASH";
                readonly SINGBOX: "SINGBOX";
            }>>;
            mapper: z.ZodPipe<z.ZodOptional<z.ZodNullable<z.ZodObject<{
                xrayJson: z.ZodOptional<z.ZodArray<z.ZodDiscriminatedUnion<[z.ZodObject<{
                    op: z.ZodLiteral<"copy">;
                    from: z.ZodString;
                    to: z.ZodString;
                }, z.core.$strip>, z.ZodObject<{
                    op: z.ZodLiteral<"set">;
                    value: z.ZodUnion<readonly [z.ZodString, z.ZodNumber, z.ZodBoolean, z.ZodArray<z.ZodJSONSchema>, z.ZodRecord<z.ZodString, z.ZodJSONSchema>]>;
                    to: z.ZodString;
                }, z.core.$strip>, z.ZodObject<{
                    op: z.ZodLiteral<"unset">;
                    to: z.ZodString;
                }, z.core.$strip>], "op">>>;
                mihomo: z.ZodOptional<z.ZodArray<z.ZodDiscriminatedUnion<[z.ZodObject<{
                    op: z.ZodLiteral<"copy">;
                    from: z.ZodString;
                    to: z.ZodString;
                }, z.core.$strip>, z.ZodObject<{
                    op: z.ZodLiteral<"set">;
                    value: z.ZodUnion<readonly [z.ZodString, z.ZodNumber, z.ZodBoolean, z.ZodArray<z.ZodJSONSchema>, z.ZodRecord<z.ZodString, z.ZodJSONSchema>]>;
                    to: z.ZodString;
                }, z.core.$strip>, z.ZodObject<{
                    op: z.ZodLiteral<"unset">;
                    to: z.ZodString;
                }, z.core.$strip>], "op">>>;
                base64: z.ZodOptional<z.ZodArray<z.ZodDiscriminatedUnion<[z.ZodObject<{
                    op: z.ZodLiteral<"copy">;
                    from: z.ZodString;
                    to: z.ZodString;
                }, z.core.$strip>, z.ZodObject<{
                    op: z.ZodLiteral<"set">;
                    value: z.ZodUnion<readonly [z.ZodString, z.ZodNumber, z.ZodBoolean, z.ZodArray<z.ZodJSONSchema>, z.ZodRecord<z.ZodString, z.ZodJSONSchema>]>;
                    to: z.ZodString;
                }, z.core.$strip>, z.ZodObject<{
                    op: z.ZodLiteral<"unset">;
                    to: z.ZodString;
                }, z.core.$strip>], "op">>>;
                singbox: z.ZodOptional<z.ZodArray<z.ZodDiscriminatedUnion<[z.ZodObject<{
                    op: z.ZodLiteral<"copy">;
                    from: z.ZodString;
                    to: z.ZodString;
                }, z.core.$strip>, z.ZodObject<{
                    op: z.ZodLiteral<"set">;
                    value: z.ZodUnion<readonly [z.ZodString, z.ZodNumber, z.ZodBoolean, z.ZodArray<z.ZodJSONSchema>, z.ZodRecord<z.ZodString, z.ZodJSONSchema>]>;
                    to: z.ZodString;
                }, z.core.$strip>, z.ZodObject<{
                    op: z.ZodLiteral<"unset">;
                    to: z.ZodString;
                }, z.core.$strip>], "op">>>;
            }, z.core.$strip>>>, z.ZodTransform<{
                xrayJson?: ({
                    op: "copy";
                    from: string;
                    to: string;
                } | {
                    op: "set";
                    value: string | number | boolean | Record<string, z.core.util.JSONType> | z.core.util.JSONType[];
                    to: string;
                } | {
                    op: "unset";
                    to: string;
                })[] | undefined;
                mihomo?: ({
                    op: "copy";
                    from: string;
                    to: string;
                } | {
                    op: "set";
                    value: string | number | boolean | Record<string, z.core.util.JSONType> | z.core.util.JSONType[];
                    to: string;
                } | {
                    op: "unset";
                    to: string;
                })[] | undefined;
                base64?: ({
                    op: "copy";
                    from: string;
                    to: string;
                } | {
                    op: "set";
                    value: string | number | boolean | Record<string, z.core.util.JSONType> | z.core.util.JSONType[];
                    to: string;
                } | {
                    op: "unset";
                    to: string;
                })[] | undefined;
                singbox?: ({
                    op: "copy";
                    from: string;
                    to: string;
                } | {
                    op: "set";
                    value: string | number | boolean | Record<string, z.core.util.JSONType> | z.core.util.JSONType[];
                    to: string;
                } | {
                    op: "unset";
                    to: string;
                })[] | undefined;
            }, {
                xrayJson?: ({
                    op: "copy";
                    from: string;
                    to: string;
                } | {
                    op: "set";
                    value: string | number | boolean | Record<string, z.core.util.JSONType> | z.core.util.JSONType[];
                    to: string;
                } | {
                    op: "unset";
                    to: string;
                })[] | undefined;
                mihomo?: ({
                    op: "copy";
                    from: string;
                    to: string;
                } | {
                    op: "set";
                    value: string | number | boolean | Record<string, z.core.util.JSONType> | z.core.util.JSONType[];
                    to: string;
                } | {
                    op: "unset";
                    to: string;
                })[] | undefined;
                base64?: ({
                    op: "copy";
                    from: string;
                    to: string;
                } | {
                    op: "set";
                    value: string | number | boolean | Record<string, z.core.util.JSONType> | z.core.util.JSONType[];
                    to: string;
                } | {
                    op: "unset";
                    to: string;
                })[] | undefined;
                singbox?: ({
                    op: "copy";
                    from: string;
                    to: string;
                } | {
                    op: "set";
                    value: string | number | boolean | Record<string, z.core.util.JSONType> | z.core.util.JSONType[];
                    to: string;
                } | {
                    op: "unset";
                    to: string;
                })[] | undefined;
            } | null | undefined>>;
        }, z.core.$strip>>;
    }, z.core.$strip>;
    type Response = z.infer<typeof ResponseSchema>;
}
//# sourceMappingURL=get-hosts.command.d.ts.map