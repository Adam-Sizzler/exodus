import { z } from 'zod';
export declare namespace UpdateSubscriptionTemplateCommand {
    const url: "/api/subscription-templates/";
    const TSQ_url: "/api/subscription-templates/";
    const endpointDetails: import("../../constants").EndpointDetails;
    const RequestBodySchema: z.ZodObject<{
        uuid: z.ZodUUID;
        name: z.ZodOptional<z.ZodString>;
        templateJson: z.ZodOptional<z.ZodObject<{}, z.core.$loose>>;
        encodedTemplateYaml: z.ZodOptional<z.ZodString>;
    }, z.core.$strip>;
    const ResponseSchema: z.ZodObject<{
        response: z.ZodObject<{
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
        }, z.core.$strip>;
    }, z.core.$strip>;
    type RequestBody = z.infer<typeof RequestBodySchema>;
    type Response = z.infer<typeof ResponseSchema>;
}
//# sourceMappingURL=update-template.command.d.ts.map