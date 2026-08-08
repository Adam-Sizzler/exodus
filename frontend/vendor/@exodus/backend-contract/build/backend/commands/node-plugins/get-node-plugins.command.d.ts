import { z } from 'zod';
export declare namespace GetNodePluginsCommand {
    const url: "/api/node-plugins/";
    const TSQ_url: "/api/node-plugins/";
    const endpointDetails: import("../../constants").EndpointDetails;
    const ResponseSchema: z.ZodObject<{
        response: z.ZodObject<{
            total: z.ZodNumber;
            nodePlugins: z.ZodArray<z.ZodObject<{
                uuid: z.ZodUUID;
                viewPosition: z.ZodNumber;
                name: z.ZodString;
                pluginConfig: z.ZodNullable<z.ZodUnknown>;
            }, z.core.$strip>>;
        }, z.core.$strip>;
    }, z.core.$strip>;
    type Response = z.infer<typeof ResponseSchema>;
}
//# sourceMappingURL=get-node-plugins.command.d.ts.map