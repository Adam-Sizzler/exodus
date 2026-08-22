import { z } from 'zod';
export declare namespace GetNodeIntegrationsCommand {
    const url: "/api/node-integrations/";
    const TSQ_url: "/api/node-integrations/";
    const endpointDetails: import("../../constants").EndpointDetails;
    const ResponseSchema: z.ZodObject<{
        response: z.ZodObject<{
            total: z.ZodNumber;
            nodeIntegrations: z.ZodArray<z.ZodObject<{
                uuid: z.ZodUUID;
                name: z.ZodString;
                description: z.ZodNullable<z.ZodString>;
                config: z.ZodRecord<z.ZodString, z.ZodUnknown>;
            }, z.core.$strip>>;
        }, z.core.$strip>;
    }, z.core.$strip>;
    type Response = z.infer<typeof ResponseSchema>;
}
//# sourceMappingURL=get-node-integrations.command.d.ts.map