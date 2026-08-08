import { z } from 'zod';
export declare namespace GetExodusHealthCommand {
    const url: "/api/system/health";
    const TSQ_url: "/api/system/health";
    const endpointDetails: import("../../constants").EndpointDetails;
    const ResponseSchema: z.ZodObject<{
        response: z.ZodObject<{
            runtimeMetrics: z.ZodArray<z.ZodObject<{
                rss: z.ZodNumber;
                heapUsed: z.ZodNumber;
                heapTotal: z.ZodNumber;
                external: z.ZodNumber;
                arrayBuffers: z.ZodNumber;
                eventLoopDelayMs: z.ZodNumber;
                eventLoopP99Ms: z.ZodNumber;
                activeHandles: z.ZodNumber;
                uptime: z.ZodNumber;
                pid: z.ZodNumber;
                timestamp: z.ZodNumber;
                instanceId: z.ZodString;
                instanceType: z.ZodString;
            }, z.core.$strip>>;
        }, z.core.$strip>;
    }, z.core.$strip>;
    type Response = z.infer<typeof ResponseSchema>;
}
//# sourceMappingURL=get-exodus-health.command.d.ts.map
export import GetExodusHealthCommand = GetExodusHealthCommand;
