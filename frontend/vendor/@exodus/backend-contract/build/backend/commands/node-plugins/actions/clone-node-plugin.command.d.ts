import { z } from 'zod';
export declare namespace CloneNodePluginCommand {
    const url: "/api/node-plugins/actions/clone";
    const TSQ_url: "/api/node-plugins/actions/clone";
    const endpointDetails: import("../../../constants").EndpointDetails;
    const RequestBodySchema: z.ZodObject<{
        cloneFromUuid: z.ZodUUID;
    }, z.core.$strip>;
    const ResponseSchema: z.ZodObject<{
        response: z.ZodObject<{
            uuid: z.ZodUUID;
            viewPosition: z.ZodNumber;
            name: z.ZodString;
            tags: z.ZodDefault<z.ZodArray<z.ZodString>>;
            pluginConfig: z.ZodUnknown;
        }, z.core.$strip>;
    }, z.core.$strip>;
    type RequestBody = z.infer<typeof RequestBodySchema>;
    type Response = z.infer<typeof ResponseSchema>;
}
//# sourceMappingURL=clone-node-plugin.command.d.ts.map