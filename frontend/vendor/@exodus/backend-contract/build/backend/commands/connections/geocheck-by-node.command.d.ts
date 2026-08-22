import { z } from 'zod';
export declare namespace GeocheckByNodeCommand {
    const url: (uuid: string) => string;
    const TSQ_url: string;
    const endpointDetails: import("../../constants").EndpointDetails;
    const RequestParamSchema: z.ZodObject<{
        nodeUuid: z.ZodUUID;
    }, z.core.$strip>;
    const RequestBodySchema: z.ZodObject<{
        ip: z.ZodOptional<z.ZodString>;
        interface: z.ZodOptional<z.ZodString>;
    }, z.core.$strip>;
    const ResponseSchema: z.ZodObject<{
        response: z.ZodObject<{
            jobId: z.ZodString;
        }, z.core.$strip>;
    }, z.core.$strip>;
    type RequestBody = z.infer<typeof RequestBodySchema>;
    type RequestParam = z.infer<typeof RequestParamSchema>;
    type Response = z.infer<typeof ResponseSchema>;
}
//# sourceMappingURL=geocheck-by-node.command.d.ts.map