import { z } from 'zod';
export declare namespace ConnectionsByNodeCommand {
    const url: (uuid: string) => string;
    const TSQ_url: string;
    const endpointDetails: import("../../constants").EndpointDetails;
    const RequestParamSchema: z.ZodObject<{
        nodeUuid: z.ZodUUID;
    }, z.core.$strip>;
    const ResponseSchema: z.ZodObject<{
        response: z.ZodObject<{
            jobId: z.ZodString;
        }, z.core.$strip>;
    }, z.core.$strip>;
    type RequestParam = z.infer<typeof RequestParamSchema>;
    type Response = z.infer<typeof ResponseSchema>;
}
//# sourceMappingURL=connections-by-node.command.d.ts.map