import z from 'zod';
export declare const SubscriptionTemplateSchema: z.ZodObject<{
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
//# sourceMappingURL=subscription-template.schema.d.ts.map