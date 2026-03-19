import { z } from 'zod';
export declare namespace BulkDeleteUsersByStatusCommand {
    const url: "/api/users/bulk/delete-by-status";
    const TSQ_url: "/api/users/bulk/delete-by-status";
    const endpointDetails: import("../../../constants").EndpointDetails;
    const RequestSchema: z.ZodObject<{
        status: z.ZodDefault<z.ZodNativeEnum<{
            readonly ACTIVE: "ACTIVE";
            readonly DISABLED: "DISABLED";
            readonly LIMITED: "LIMITED";
            readonly EXPIRED: "EXPIRED";
        }>>;
    }, "strip", z.ZodTypeAny, {
        status: "DISABLED" | "LIMITED" | "EXPIRED" | "ACTIVE";
    }, {
        status?: "DISABLED" | "LIMITED" | "EXPIRED" | "ACTIVE" | undefined;
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
//# sourceMappingURL=bulk-delete-users-by-status.command.d.ts.map