import { z } from 'zod';
export declare namespace UpdateNodeIntegrationCommand {
    const url: "/api/node-integrations/";
    const TSQ_url: "/api/node-integrations/";
    const endpointDetails: import("../../constants").EndpointDetails;
    const RequestBodySchema: z.ZodObject<{
        uuid: z.ZodUUID;
        name: z.ZodOptional<z.ZodString>;
        description: z.ZodOptional<z.ZodNullable<z.ZodString>>;
        config: z.ZodOptional<z.ZodRecord<z.ZodString, z.ZodUnknown>>;
        restartNodes: z.ZodOptional<z.ZodBoolean>;
    }, z.core.$strip>;
    const ResponseSchema: z.ZodObject<{
        response: z.ZodObject<{
            uuid: z.ZodUUID;
            name: z.ZodString;
            description: z.ZodNullable<z.ZodString>;
            config: z.ZodRecord<z.ZodString, z.ZodUnknown>;
        }, z.core.$strip>;
    }, z.core.$strip>;
    type RequestBody = z.infer<typeof RequestBodySchema>;
    type Response = z.infer<typeof ResponseSchema>;
}
//# sourceMappingURL=update-node-integration.command.d.ts.map