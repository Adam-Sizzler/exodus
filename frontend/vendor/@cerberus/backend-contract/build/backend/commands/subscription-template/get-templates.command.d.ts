import { z } from 'zod';
export declare namespace GetSubscriptionTemplatesCommand {
    const url: "/api/subscription-templates/";
    const TSQ_url: "/api/subscription-templates/";
    const endpointDetails: import("../../constants").EndpointDetails;
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
//# sourceMappingURL=get-templates.command.d.ts.map