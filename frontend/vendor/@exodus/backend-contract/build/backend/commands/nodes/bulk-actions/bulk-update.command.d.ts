import { z } from 'zod';
export declare namespace BulkNodesUpdateCommand {
    const url: "/api/nodes/bulk-actions/update";
    const TSQ_url: "/api/nodes/bulk-actions/update";
    const endpointDetails: import("../../../constants").EndpointDetails;
    const RequestBodySchema: z.ZodObject<{
        uuids: z.ZodArray<z.ZodUUID>;
        fields: z.ZodObject<{
            countryCode: z.ZodOptional<z.ZodString>;
            consumptionMultiplier: z.ZodOptional<z.ZodPipe<z.ZodNumber, z.ZodTransform<number, number>>>;
            nodeConsumptionMultiplier: z.ZodOptional<z.ZodPipe<z.ZodNumber, z.ZodTransform<number, number>>>;
            providerUuid: z.ZodOptional<z.ZodNullable<z.ZodUUID>>;
            tags: z.ZodOptional<z.ZodArray<z.ZodString>>;
            activePluginUuid: z.ZodOptional<z.ZodNullable<z.ZodUUID>>;
            integrationUuids: z.ZodOptional<z.ZodArray<z.ZodUUID>>;
            note: z.ZodOptional<z.ZodNullable<z.ZodString>>;
        }, z.core.$strip>;
    }, z.core.$strip>;
    type RequestBody = z.infer<typeof RequestBodySchema>;
}
//# sourceMappingURL=bulk-update.command.d.ts.map