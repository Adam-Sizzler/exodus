import { z } from 'zod';
export declare namespace BulkNodesProfileModificationCommand {
    const url: "/api/nodes/bulk-actions/profile-modification";
    const TSQ_url: "/api/nodes/bulk-actions/profile-modification";
    const endpointDetails: import("../../../constants").EndpointDetails;
    const RequestSchema: z.ZodObject<{
        uuids: z.ZodArray<z.ZodString, "many">;
        configProfile: z.ZodObject<{
            activeConfigProfileUuid: z.ZodString;
            activeInbounds: z.ZodArray<z.ZodString, "many">;
        }, "strip", z.ZodTypeAny, {
            activeConfigProfileUuid: string;
            activeInbounds: string[];
        }, {
            activeConfigProfileUuid: string;
            activeInbounds: string[];
        }>;
    }, "strip", z.ZodTypeAny, {
        configProfile: {
            activeConfigProfileUuid: string;
            activeInbounds: string[];
        };
        uuids: string[];
    }, {
        configProfile: {
            activeConfigProfileUuid: string;
            activeInbounds: string[];
        };
        uuids: string[];
    }>;
    type Request = z.infer<typeof RequestSchema>;
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
//# sourceMappingURL=profile-modification.command.d.ts.map