import { z } from 'zod';
export declare namespace RestartAllNodesCommand {
    const url: "/api/nodes/actions/restart-all";
    const TSQ_url: "/api/nodes/actions/restart-all";
    const endpointDetails: import("../../../constants").EndpointDetails;
    const RequestBodySchema: z.ZodObject<{
        forceRestart: z.ZodBoolean;
    }, z.core.$strip>;
    type RequestBody = z.infer<typeof RequestBodySchema>;
}
//# sourceMappingURL=restart-all.command.d.ts.map