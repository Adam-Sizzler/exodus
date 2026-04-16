import { z } from 'zod';
export declare namespace CreateApiTokenCommand {
    const url: "/api/tokens/";
    const TSQ_url: "/api/tokens/";
    const endpointDetails: import("../../constants").EndpointDetails;
    const RequestSchema: z.ZodObject<{
        tokenName: z.ZodString;
    }, "strip", z.ZodTypeAny, {
        tokenName: string;
    }, {
        tokenName: string;
    }>;
    type Request = z.infer<typeof RequestSchema>;
    const ResponseSchema: z.ZodObject<{
        response: z.ZodObject<{
            token: z.ZodString;
            uuid: z.ZodString;
        }, "strip", z.ZodTypeAny, {
            token: string;
            uuid: string;
        }, {
            token: string;
            uuid: string;
        }>;
    }, "strip", z.ZodTypeAny, {
        response: {
            token: string;
            uuid: string;
        };
    }, {
        response: {
            token: string;
            uuid: string;
        };
    }>;
    type Response = z.infer<typeof ResponseSchema>;
}
//# sourceMappingURL=create.command.d.ts.map