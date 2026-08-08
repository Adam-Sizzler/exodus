import { z } from 'zod';
export declare namespace RestartNodeCommand {
    const url: (uuid: string) => string;
    const TSQ_url: string;
    const endpointDetails: import("../../../constants").EndpointDetails;
    const RequestParamSchema: z.ZodObject<{
        uuid: z.ZodUUID;
    }, z.core.$strip>;
    const RequestBodySchema: z.ZodObject<{
        forceRestart: z.ZodBoolean;
    }, z.core.$strip>;
    type RequestParam = z.infer<typeof RequestParamSchema>;
    type RequestBody = z.infer<typeof RequestBodySchema>;
}
//# sourceMappingURL=restart.command.d.ts.map