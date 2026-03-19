import { z } from 'zod';
export declare namespace GetSubscriptionTemplateCommand {
    const url: (uuid: string) => string;
    const TSQ_url: string;
    const endpointDetails: import("../../constants").EndpointDetails;
    const RequestSchema: z.ZodObject<{
        uuid: z.ZodString;
    }, "strip", z.ZodTypeAny, {
        uuid: string;
    }, {
        uuid: string;
    }>;
    type Request = z.infer<typeof RequestSchema>;
    const ResponseSchema: z.ZodObject<{
        response: z.ZodObject<{
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
        }>;
    }, "strip", z.ZodTypeAny, {
        response: {
            uuid: string;
            name: string;
            viewPosition: number;
            templateType: "STASH" | "SINGBOX" | "MIHOMO" | "XRAY_JSON" | "CLASH" | "XRAY_BASE64";
            encodedTemplateYaml: string | null;
            templateJson?: unknown;
        };
    }, {
        response: {
            uuid: string;
            name: string;
            viewPosition: number;
            templateType: "STASH" | "SINGBOX" | "MIHOMO" | "XRAY_JSON" | "CLASH" | "XRAY_BASE64";
            encodedTemplateYaml: string | null;
            templateJson?: unknown;
        };
    }>;
    type Response = z.infer<typeof ResponseSchema>;
}
//# sourceMappingURL=get-template.command.d.ts.map