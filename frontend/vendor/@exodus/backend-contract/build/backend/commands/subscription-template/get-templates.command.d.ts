import { z } from 'zod';
export declare namespace GetSubscriptionTemplatesCommand {
    const url: "/api/subscription-templates/";
    const TSQ_url: "/api/subscription-templates/";
    const endpointDetails: import("../../constants").EndpointDetails;
    const ResponseSchema: z.ZodObject<{
        response: z.ZodObject<{
            total: z.ZodNumber;
            templates: z.ZodArray<z.ZodObject<{
                uuid: z.ZodUUID;
                viewPosition: z.ZodNumber;
                name: z.ZodString;
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
            }, z.core.$strip>>;
        }, z.core.$strip>;
    }, z.core.$strip>;
    type Response = z.infer<typeof ResponseSchema>;
}
//# sourceMappingURL=get-templates.command.d.ts.map