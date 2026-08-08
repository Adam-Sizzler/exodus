import { z } from 'zod';
export declare namespace GetNodeUsageCommand {
    const url: "/api/bandwidth-stats/nodes/usage";
    const TSQ_url: "/api/bandwidth-stats/nodes/usage";
    const endpointDetails: import("../../../constants").EndpointDetails;
    const RequestBodySchema: z.ZodObject<{
        nodesUuids: z.ZodArray<z.ZodUUID>;
    }, z.core.$strip>;
    const RequestQuerySchema: z.ZodObject<{
        start: z.ZodISODate;
        end: z.ZodISODate;
        minTotalBytes: z.ZodDefault<z.ZodOptional<z.ZodCoercedNumber<unknown>>>;
    }, z.core.$strip>;
    const ResponseSchema: z.ZodObject<{
        response: z.ZodObject<{
            nodes: z.ZodArray<z.ZodObject<{
                uuid: z.ZodUUID;
                users: z.ZodArray<z.ZodObject<{
                    id: z.ZodNumber;
                    totalBytes: z.ZodNumber;
                }, z.core.$strip>>;
            }, z.core.$strip>>;
        }, z.core.$strip>;
    }, z.core.$strip>;
    type RequestBody = z.infer<typeof RequestBodySchema>;
    type RequestQuery = z.infer<typeof RequestQuerySchema>;
    type Response = z.infer<typeof ResponseSchema>;
}
//# sourceMappingURL=get-node-usage.command.d.ts.map