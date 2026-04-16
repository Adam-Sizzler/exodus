import { z } from 'zod';
export declare namespace ReorderSubscriptionTemplateCommand {
    const url: "/api/subscription-templates/actions/reorder";
    const TSQ_url: "/api/subscription-templates/actions/reorder";
    const endpointDetails: import("../../../constants").EndpointDetails;
    const RequestSchema: z.ZodObject<{
        items: z.ZodArray<z.ZodObject<Pick<{
            uuid: z.ZodString;
            viewPosition: z.ZodNumber;
            name: z.ZodString;
            templateType: z.ZodNativeEnum<{
                readonly XRAY_JSON: "XRAY_JSON";
                readonly XRAY_BASE64: "XRAY_BASE64";
                readonly MIHOMO: "MIHOMO";
                readonly STASH: "STASH";
                readonly CLASH: "CLASH";
                readonly SINGBOX: "SINGBOX";
            }>;
            templateJson: z.ZodNullable<z.ZodUnknown>;
            encodedTemplateYaml: z.ZodNullable<z.ZodString>;
        }, "uuid" | "viewPosition">, "strip", z.ZodTypeAny, {
            uuid: string;
            viewPosition: number;
        }, {
            uuid: string;
            viewPosition: number;
        }>, "many">;
    }, "strip", z.ZodTypeAny, {
        items: {
            uuid: string;
            viewPosition: number;
        }[];
    }, {
        items: {
            uuid: string;
            viewPosition: number;
        }[];
    }>;
    type Request = z.infer<typeof RequestSchema>;
    const ResponseSchema: z.ZodObject<{
        response: z.ZodObject<{
            total: z.ZodNumber;
            templates: z.ZodArray<z.ZodObject<{
                uuid: z.ZodString;
                viewPosition: z.ZodNumber;
                name: z.ZodString;
                templateType: z.ZodNativeEnum<{
                    readonly XRAY_JSON: "XRAY_JSON";
                    readonly XRAY_BASE64: "XRAY_BASE64";
                    readonly MIHOMO: "MIHOMO";
                    readonly STASH: "STASH";
                    readonly CLASH: "CLASH";
                    readonly SINGBOX: "SINGBOX";
                }>;
                templateJson: z.ZodNullable<z.ZodUnknown>;
                encodedTemplateYaml: z.ZodNullable<z.ZodString>;
            }, "strip", z.ZodTypeAny, {
                uuid: string;
                name: string;
                viewPosition: number;
                templateType: "STASH" | "SINGBOX" | "MIHOMO" | "XRAY_JSON" | "CLASH" | "XRAY_BASE64";
                encodedTemplateYaml: string | null;
                templateJson?: unknown;
            }, {
                uuid: string;
                name: string;
                viewPosition: number;
                templateType: "STASH" | "SINGBOX" | "MIHOMO" | "XRAY_JSON" | "CLASH" | "XRAY_BASE64";
                encodedTemplateYaml: string | null;
                templateJson?: unknown;
            }>, "many">;
        }, "strip", z.ZodTypeAny, {
            total: number;
            templates: {
                uuid: string;
                name: string;
                viewPosition: number;
                templateType: "STASH" | "SINGBOX" | "MIHOMO" | "XRAY_JSON" | "CLASH" | "XRAY_BASE64";
                encodedTemplateYaml: string | null;
                templateJson?: unknown;
            }[];
        }, {
            total: number;
            templates: {
                uuid: string;
                name: string;
                viewPosition: number;
                templateType: "STASH" | "SINGBOX" | "MIHOMO" | "XRAY_JSON" | "CLASH" | "XRAY_BASE64";
                encodedTemplateYaml: string | null;
                templateJson?: unknown;
            }[];
        }>;
    }, "strip", z.ZodTypeAny, {
        response: {
            total: number;
            templates: {
                uuid: string;
                name: string;
                viewPosition: number;
                templateType: "STASH" | "SINGBOX" | "MIHOMO" | "XRAY_JSON" | "CLASH" | "XRAY_BASE64";
                encodedTemplateYaml: string | null;
                templateJson?: unknown;
            }[];
        };
    }, {
        response: {
            total: number;
            templates: {
                uuid: string;
                name: string;
                viewPosition: number;
                templateType: "STASH" | "SINGBOX" | "MIHOMO" | "XRAY_JSON" | "CLASH" | "XRAY_BASE64";
                encodedTemplateYaml: string | null;
                templateJson?: unknown;
            }[];
        };
    }>;
    type Response = z.infer<typeof ResponseSchema>;
}
//# sourceMappingURL=reorder.command.d.ts.map