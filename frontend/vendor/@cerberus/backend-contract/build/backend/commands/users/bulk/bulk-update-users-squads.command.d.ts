import { z } from 'zod';
export declare namespace BulkUpdateUsersSquadsCommand {
    const url: "/api/users/bulk/update-squads";
    const TSQ_url: "/api/users/bulk/update-squads";
    const endpointDetails: import("../../../constants").EndpointDetails;
    const RequestSchema: z.ZodObject<{
        uuids: z.ZodArray<z.ZodString, "many">;
        activeInternalSquads: z.ZodArray<z.ZodString, "many">;
    }, "strip", z.ZodTypeAny, {
        activeInternalSquads: string[];
        uuids: string[];
    }, {
        activeInternalSquads: string[];
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
//# sourceMappingURL=bulk-update-users-squads.command.d.ts.map