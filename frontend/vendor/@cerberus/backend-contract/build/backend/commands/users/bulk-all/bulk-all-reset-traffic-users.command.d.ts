import { z } from 'zod';
export declare namespace BulkAllResetTrafficUsersCommand {
    const url: "/api/users/bulk/all/reset-traffic";
    const TSQ_url: "/api/users/bulk/all/reset-traffic";
    const endpointDetails: import("../../../constants").EndpointDetails;
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
//# sourceMappingURL=bulk-all-reset-traffic-users.command.d.ts.map