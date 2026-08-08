import { z } from 'zod';
export declare namespace GetHttpStatsCommand {
    const url: "/api/system/stats/http";
    const TSQ_url: "/api/system/stats/http";
    const endpointDetails: import("../../constants").EndpointDetails;
    const ResponseSchema: z.ZodObject<{
        response: z.ZodObject<{
            routes: z.ZodArray<z.ZodObject<{
                method: z.ZodString;
                route: z.ZodString;
                count: z.ZodInt32;
            }, z.core.$strip>>;
            total: z.ZodInt32;
        }, z.core.$strip>;
    }, z.core.$strip>;
    type Response = z.infer<typeof ResponseSchema>;
}
//# sourceMappingURL=get-http-stats.command.d.ts.map