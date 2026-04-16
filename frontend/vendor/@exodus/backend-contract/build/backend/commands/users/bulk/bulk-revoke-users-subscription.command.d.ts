import { z } from 'zod';
export declare namespace BulkRevokeUsersSubscriptionCommand {
    const url: "/api/users/bulk/revoke-subscription";
    const TSQ_url: "/api/users/bulk/revoke-subscription";
    const endpointDetails: import("../../../constants").EndpointDetails;
    const RequestSchema: z.ZodObject<{
        uuids: z.ZodArray<z.ZodString, "many">;
    }, "strip", z.ZodTypeAny, {
        uuids: string[];
    }, {
        uuids: string[];
    }>;
    type Request = z.infer<typeof RequestSchema>;
    const ResponseSchema: z.ZodObject<{
        response: z.ZodObject<{
            affectedRows: z.ZodNumber;
        }, "strip", z.ZodTypeAny, {
            affectedRows: number;
        }, {
            affectedRows: number;
        }>;
    }, "strip", z.ZodTypeAny, {
        response: {
            affectedRows: number;
        };
    }, {
        response: {
            affectedRows: number;
        };
    }>;
    type Response = z.infer<typeof ResponseSchema>;
}
//# sourceMappingURL=bulk-revoke-users-subscription.command.d.ts.map