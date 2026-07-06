import { z } from 'zod';
export declare namespace RestartAllNodesCommand {
    const url: "/api/nodes/actions/restart-all";
    const TSQ_url: "/api/nodes/actions/restart-all";
    const endpointDetails: import("../../../constants").EndpointDetails;
    const RequestBodySchema: z.ZodObject<{
        forceRestart: z.ZodBoolean;
    }, "strip", z.ZodTypeAny, {
        forceRestart: boolean;
    }, {
        forceRestart: boolean;
    }>;
    type RequestBody = z.infer<typeof RequestBodySchema>;
    const ResponseSchema: z.ZodObject<{
        response: z.ZodObject<{
            eventSent: z.ZodBoolean;
        }, "strip", z.ZodTypeAny, {
            eventSent: boolean;
        }, {
            eventSent: boolean;
        }>;
    }, "strip", z.ZodTypeAny, {
        response: {
            eventSent: boolean;
        };
    }, {
        response: {
            eventSent: boolean;
        };
    }>;
    type Response = z.infer<typeof ResponseSchema>;
}
//# sourceMappingURL=restart-all.command.d.ts.map