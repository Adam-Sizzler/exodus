import { z } from 'zod';
export declare namespace ReorderSubscriptionTemplateCommand {
    const url: "/api/subscription-templates/actions/reorder";
    const TSQ_url: "/api/subscription-templates/actions/reorder";
    const endpointDetails: import("../../../constants").EndpointDetails;
    const RequestBodySchema: z.ZodObject<{
        items: z.ZodArray<z.ZodObject<{
            uuid: z.ZodUUID;
            viewPosition: z.ZodNumber;
        }, z.core.$strip>>;
    }, z.core.$strip>;
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
    type RequestBody = z.infer<typeof RequestBodySchema>;
    type Response = z.infer<typeof ResponseSchema>;
}
//# sourceMappingURL=reorder.command.d.ts.map
