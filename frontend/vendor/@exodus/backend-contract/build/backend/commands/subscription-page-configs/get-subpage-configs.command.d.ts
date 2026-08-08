import { z } from 'zod';
export declare namespace GetSubpageConfigsCommand {
    const url: "/api/subscription-page-configs/";
    const TSQ_url: "/api/subscription-page-configs/";
    const endpointDetails: import("../../constants").EndpointDetails;
    const ResponseSchema: z.ZodObject<{
        response: z.ZodObject<{
            total: z.ZodNumber;
            configs: z.ZodArray<z.ZodObject<{
                uuid: z.ZodUUID;
                viewPosition: z.ZodNumber;
                name: z.ZodString;
                config: z.ZodNullable<z.ZodUnknown>;
            }, z.core.$strip>>;
        }, z.core.$strip>;
    }, z.core.$strip>;
    type Response = z.infer<typeof ResponseSchema>;
}
//# sourceMappingURL=get-subpage-configs.command.d.ts.map