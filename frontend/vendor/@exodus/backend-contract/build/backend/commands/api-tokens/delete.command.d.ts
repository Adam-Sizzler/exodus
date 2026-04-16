import { z } from 'zod';
export declare namespace DeleteApiTokenCommand {
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
        response: z.ZodBoolean;
    }, "strip", z.ZodTypeAny, {
        response: boolean;
    }, {
        response: boolean;
    }>;
    type Response = z.infer<typeof ResponseSchema>;
}
//# sourceMappingURL=delete.command.d.ts.map