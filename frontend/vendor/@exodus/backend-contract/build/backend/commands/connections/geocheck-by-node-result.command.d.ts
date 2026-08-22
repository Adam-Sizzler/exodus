import { z } from 'zod';
export declare namespace GeocheckByNodeResultCommand {
    const url: (jobId: string) => string;
    const TSQ_url: string;
    const endpointDetails: import("../../constants").EndpointDetails;
    const RequestParamSchema: z.ZodObject<{
        jobId: z.ZodString;
    }, z.core.$strip>;
    const GeocheckImageSchema: z.ZodObject<{
        format: z.ZodLiteral<"svg">;
        media_type: z.ZodLiteral<"image/svg+xml">;
        encoding: z.ZodLiteral<"base64">;
        data: z.ZodString;
    }, z.core.$strip>;
    const ResponseSchema: z.ZodObject<{
        response: z.ZodObject<{
            isCompleted: z.ZodBoolean;
            isFailed: z.ZodBoolean;
            result: z.ZodNullable<z.ZodObject<{
                success: z.ZodBoolean;
                nodeUuid: z.ZodUUID;
                image: z.ZodNullable<z.ZodObject<{
                    format: z.ZodLiteral<"svg">;
                    media_type: z.ZodLiteral<"image/svg+xml">;
                    encoding: z.ZodLiteral<"base64">;
                    data: z.ZodString;
                }, z.core.$strip>>;
                rawReport: z.ZodNullable<z.ZodRecord<z.ZodString, z.ZodUnknown>>;
                message: z.ZodNullable<z.ZodString>;
            }, z.core.$strip>>;
        }, z.core.$strip>;
    }, z.core.$strip>;
    type GeocheckImage = z.infer<typeof GeocheckImageSchema>;
    type RequestParam = z.infer<typeof RequestParamSchema>;
    type Response = z.infer<typeof ResponseSchema>;
}
//# sourceMappingURL=geocheck-by-node-result.command.d.ts.map