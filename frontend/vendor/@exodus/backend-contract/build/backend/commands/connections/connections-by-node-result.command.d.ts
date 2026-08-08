import { z } from 'zod';
export declare namespace ConnectionsByNodeResultCommand {
    const url: (jobId: string) => string;
    const TSQ_url: string;
    const endpointDetails: import("../../constants").EndpointDetails;
    const RequestParamSchema: z.ZodObject<{
        jobId: z.ZodString;
    }, z.core.$strip>;
    const ResponseSchema: z.ZodObject<{
        response: z.ZodObject<{
            isCompleted: z.ZodBoolean;
            isFailed: z.ZodBoolean;
            result: z.ZodNullable<z.ZodObject<{
                success: z.ZodBoolean;
                nodeUuid: z.ZodUUID;
                users: z.ZodArray<z.ZodObject<{
                    userId: z.ZodNumber;
                    ips: z.ZodArray<z.ZodObject<{
                        ip: z.ZodString;
                        lastSeen: z.ZodPipe<z.ZodISODateTime, z.ZodTransform<Date, string>>;
                    }, z.core.$strip>>;
                }, z.core.$strip>>;
            }, z.core.$strip>>;
        }, z.core.$strip>;
    }, z.core.$strip>;
    type RequestParam = z.infer<typeof RequestParamSchema>;
    type Response = z.infer<typeof ResponseSchema>;
}
//# sourceMappingURL=connections-by-node-result.command.d.ts.map