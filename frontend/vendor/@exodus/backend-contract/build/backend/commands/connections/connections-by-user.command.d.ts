import { z } from 'zod';
export declare namespace ConnectionsByUserCommand {
    const url: (userId: string) => string;
    const TSQ_url: string;
    const endpointDetails: import("../../constants").EndpointDetails;
    const RequestParamSchema: z.ZodObject<{
        userId: z.ZodCoercedNumber<unknown>;
    }, z.core.$strip>;
    const ResponseSchema: z.ZodObject<{
        response: z.ZodObject<{
            jobId: z.ZodString;
        }, z.core.$strip>;
    }, z.core.$strip>;
    type RequestParam = z.infer<typeof RequestParamSchema>;
    type Response = z.infer<typeof ResponseSchema>;
}
//# sourceMappingURL=connections-by-user.command.d.ts.map