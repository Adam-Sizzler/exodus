import { z } from 'zod';
export declare namespace FetchUsersIpsCommand {
    const url: (nodeUuid: string) => string;
    const TSQ_url: string;
    const endpointDetails: import("../../constants").EndpointDetails;
    const RequestSchema: z.ZodObject<{
        nodeUuid: z.ZodString;
    }, "strip", z.ZodTypeAny, {
        nodeUuid: string;
    }, {
        nodeUuid: string;
    }>;
    type Request = z.infer<typeof RequestSchema>;
    const ResponseSchema: z.ZodObject<{
        response: z.ZodObject<{
            jobId: z.ZodString;
        }, "strip", z.ZodTypeAny, {
            jobId: string;
        }, {
            jobId: string;
        }>;
    }, "strip", z.ZodTypeAny, {
        response: {
            jobId: string;
        };
    }, {
        response: {
            jobId: string;
        };
    }>;
    type Response = z.infer<typeof ResponseSchema>;
}
//# sourceMappingURL=fetch-users-ips.command.d.ts.map