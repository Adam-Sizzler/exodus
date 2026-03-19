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