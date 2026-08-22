import { z } from 'zod';
export declare namespace UpdateHostCommand {
    const url: "/api/hosts/";
    const TSQ_url: "/api/hosts/";
    const endpointDetails: import("../../constants").EndpointDetails;
    const RequestBodySchema: z.ZodObject<{
        uuid: z.ZodUUID;
        inbound: z.ZodOptional<z.ZodObject<{
            configProfileUuid: z.ZodUUID;
            configProfileInboundUuid: z.ZodUUID;
        }, z.core.$strip>>;
        remark: z.ZodOptional<z.ZodString>;
        address: z.ZodOptional<z.ZodString>;
        port: z.ZodOptional<z.ZodInt>;
        path: z.ZodOptional<z.ZodNullable<z.ZodString>>;
        sni: z.ZodOptional<z.ZodNullable<z.ZodString>>;
        host: z.ZodOptional<z.ZodNullable<z.ZodString>>;
        alpn: z.ZodOptional<z.ZodNullable<z.ZodEnum<{
            readonly H3: "h3";
            readonly H2: "h2";
            readonly HTTP_1_1: "http/1.1";
            readonly H_COMBINED: "h2,http/1.1";
            readonly H3_H2_H1_COMBINED: "h3,h2,http/1.1";
            readonly H3_H2_COMBINED: "h3,h2";
        }>>>;
        fingerprint: z.ZodOptional<z.ZodNullable<z.ZodString>>;
        isDisabled: z.ZodDefault<z.ZodBoolean>;
        securityLayer: z.ZodOptional<z.ZodEnum<{
            readonly DEFAULT: "DEFAULT";
            readonly TLS: "TLS";
            readonly NONE: "NONE";
        }>>;
        xhttpExtraParams: z.ZodOptional<z.ZodNullable<z.ZodUnknown>>;
        muxParams: z.ZodOptional<z.ZodNullable<z.ZodUnknown>>;
        sockoptParams: z.ZodOptional<z.ZodNullable<z.ZodUnknown>>;
        finalMask: z.ZodOptional<z.ZodNullable<z.ZodUnknown>>;
        serverDescription: z.ZodOptional<z.ZodNullable<z.ZodString>>;
        tags: z.ZodOptional<z.ZodArray<z.ZodString>>;
        isHidden: z.ZodOptional<z.ZodBoolean>;
        overrideSniFromAddress: z.ZodOptional<z.ZodBoolean>;
        keepSniBlank: z.ZodOptional<z.ZodBoolean>;
        vlessRouteId: z.ZodOptional<z.ZodNullable<z.ZodInt>>;
        pinnedPeerCertSha256: z.ZodOptional<z.ZodNullable<z.ZodString>>;
        verifyPeerCertByName: z.ZodOptional<z.ZodNullable<z.ZodString>>;
        shuffleHost: z.ZodOptional<z.ZodBoolean>;
        mihomoX25519: z.ZodOptional<z.ZodBoolean>;
        mihomoIpVersion: z.ZodOptional<z.ZodNullable<z.ZodEnum<{
            readonly DUAL: "dual";
            readonly IPV4: "ipv4";
            readonly IPV6: "ipv6";
            readonly IPV4_PREFER: "ipv4-prefer";
            readonly IPV6_PREFER: "ipv6-prefer";
        }>>>;
        nodes: z.ZodOptional<z.ZodArray<z.ZodUUID>>;
        xrayJsonTemplateUuid: z.ZodOptional<z.ZodNullable<z.ZodUUID>>;
        excludedInternalSquads: z.ZodOptional<z.ZodArray<z.ZodUUID>>;
        excludeFromSubscriptionTypes: z.ZodOptional<z.ZodArray<z.ZodEnum<{
            readonly XRAY_JSON: "XRAY_JSON";
            readonly XRAY_BASE64: "XRAY_BASE64";
            readonly MIHOMO: "MIHOMO";
            readonly STASH: "STASH";
            readonly CLASH: "CLASH";
            readonly SINGBOX: "SINGBOX";
        }>>>;
        mapper: z.ZodOptional<z.ZodObject<{
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
        }, z.core.$strip>>;
    }, z.core.$strip>;
    const ResponseSchema: z.ZodObject<{
        response: z.ZodObject<{
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
            mapper: z.ZodObject<{
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
            }, z.core.$strip>;
        }, z.core.$strip>;
    }, z.core.$strip>;
    type RequestBody = z.infer<typeof RequestBodySchema>;
    type Request = RequestBody;
    type Response = z.infer<typeof ResponseSchema>;
}
//# sourceMappingURL=update.command.d.ts.map