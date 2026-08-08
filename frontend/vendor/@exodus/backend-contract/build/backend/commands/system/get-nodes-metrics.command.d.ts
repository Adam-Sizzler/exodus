import { z } from 'zod';
export declare namespace GetNodesMetricsCommand {
    const url: "/api/system/nodes/metrics";
    const TSQ_url: "/api/system/nodes/metrics";
    const endpointDetails: import("../../constants").EndpointDetails;
    const ResponseSchema: z.ZodObject<{
        response: z.ZodObject<{
            nodes: z.ZodArray<z.ZodObject<{
                nodeUuid: z.ZodString;
                nodeName: z.ZodString;
                countryEmoji: z.ZodString;
                providerName: z.ZodString;
                usersOnline: z.ZodNumber;
                inboundsStats: z.ZodArray<z.ZodObject<{
                    tag: z.ZodString;
                    upload: z.ZodString;
                    download: z.ZodString;
                }, z.core.$strip>>;
                outboundsStats: z.ZodArray<z.ZodObject<{
                    tag: z.ZodString;
                    upload: z.ZodString;
                    download: z.ZodString;
                }, z.core.$strip>>;
            }, z.core.$strip>>;
        }, z.core.$strip>;
    }, z.core.$strip>;
    type Response = z.infer<typeof ResponseSchema>;
}
//# sourceMappingURL=get-nodes-metrics.command.d.ts.map