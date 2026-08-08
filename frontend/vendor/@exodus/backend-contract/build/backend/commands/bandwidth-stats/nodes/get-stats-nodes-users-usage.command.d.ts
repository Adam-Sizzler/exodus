import { z } from 'zod';
export declare namespace GetStatsNodesUsersUsageCommand {
    const url: "/api/bandwidth-stats/nodes/users";
    const TSQ_url: "/api/bandwidth-stats/nodes/users";
    const endpointDetails: import("../../../constants").EndpointDetails;
    const RequestBodySchema: z.ZodObject<{
        nodesUuids: z.ZodArray<z.ZodUUID>;
    }, z.core.$strip>;
    const RequestQuerySchema: z.ZodObject<{
        start: z.ZodISODate;
        end: z.ZodISODate;
        topUsersLimit: z.ZodDefault<z.ZodCoercedNumber<unknown>>;
    }, z.core.$strip>;
    const ResponseSchema: z.ZodObject<{
        response: z.ZodObject<{
            categories: z.ZodArray<z.ZodString>;
            sparklineData: z.ZodArray<z.ZodNumber>;
            topUsers: z.ZodArray<z.ZodObject<{
                color: z.ZodString;
                username: z.ZodString;
                total: z.ZodNumber;
            }, z.core.$strip>>;
        }, z.core.$strip>;
    }, z.core.$strip>;
    type RequestBody = z.infer<typeof RequestBodySchema>;
    type RequestQuery = z.infer<typeof RequestQuerySchema>;
    type Response = z.infer<typeof ResponseSchema>;
}
//# sourceMappingURL=get-stats-nodes-users-usage.command.d.ts.map