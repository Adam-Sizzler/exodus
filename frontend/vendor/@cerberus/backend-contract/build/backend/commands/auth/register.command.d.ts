import { z } from 'zod';
export declare namespace RegisterCommand {
    const url: "/api/auth/register";
    const TSQ_url: "/api/auth/register";
    const endpointDetails: import("../../constants").EndpointDetails;
    const RequestSchema: z.ZodObject<{
        username: z.ZodString;
        password: z.ZodString;
    }, "strip", z.ZodTypeAny, {
        password: string;
        username: string;
    }, {
        password: string;
        username: string;
    }>;
    type Request = z.infer<typeof RequestSchema>;
    const ResponseSchema: z.ZodObject<{
        response: z.ZodObject<{
            accessToken: z.ZodString;
        }, "strip", z.ZodTypeAny, {
            accessToken: string;
        }, {
            accessToken: string;
        }>;
    }, "strip", z.ZodTypeAny, {
        response: {
            accessToken: string;
        };
    }, {
        response: {
            accessToken: string;
        };
    }>;
    type Response = z.infer<typeof ResponseSchema>;
}
//# sourceMappingURL=register.command.d.ts.map