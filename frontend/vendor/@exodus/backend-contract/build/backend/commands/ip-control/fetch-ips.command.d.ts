import { z } from 'zod';
export declare namespace FetchIpsCommand {
    const url: (uuid: string) => string;
    const TSQ_url: string;
    const endpointDetails: import("../../constants").EndpointDetails;
    const RequestSchema: z.ZodObject<{
        uuid: z.ZodString;
    }, "strip", z.ZodTypeAny, {
        uuid: string;
    }, {
        uuid: string;
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
//# sourceMappingURL=fetch-ips.command.d.ts.map