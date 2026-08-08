import { z } from 'zod';
export declare namespace ConnectionsByUserResultCommand {
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
            progress: z.ZodObject<{
                total: z.ZodNumber;
                completed: z.ZodNumber;
                percent: z.ZodNumber;
            }, z.core.$strip>;
            result: z.ZodNullable<z.ZodObject<{
                success: z.ZodBoolean;
                userId: z.ZodNumber;
                nodes: z.ZodArray<z.ZodObject<{
                    nodeUuid: z.ZodUUID;
                    nodeName: z.ZodString;
                    countryCode: z.ZodString;
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
//# sourceMappingURL=connections-by-user-result.command.d.ts.map