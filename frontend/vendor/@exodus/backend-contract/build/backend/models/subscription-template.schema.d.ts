import z from 'zod';
export declare const SubscriptionTemplateSchema: z.ZodObject<{
    uuid: z.ZodUUID;
    viewPosition: z.ZodNumber;
    name: z.ZodString;
    tags: z.ZodDefault<z.ZodArray<z.ZodString>>;
    templateType: z.ZodEnum<{
        readonly XRAY_JSON: "XRAY_JSON";
        readonly XRAY_BASE64: "XRAY_BASE64";
        readonly MIHOMO: "MIHOMO";
        readonly STASH: "STASH";
        readonly CLASH: "CLASH";
        readonly SINGBOX: "SINGBOX";
    }>;
    templateJson: z.ZodNullable<z.ZodUnknown>;
    encodedTemplateYaml: z.ZodNullable<z.ZodString>;
}, z.core.$strip>;
//# sourceMappingURL=subscription-template.schema.d.ts.map