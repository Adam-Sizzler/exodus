import { z } from 'zod';
export declare namespace RegisterCommand {
    const url: "/api/auth/register";
    const TSQ_url: "/api/auth/register";
    const endpointDetails: import("../../constants").EndpointDetails;
    const RequestBodySchema: z.ZodObject<{
        username: z.ZodString;
        password: z.ZodString;
    }, z.core.$strip>;
    const ResponseSchema: z.ZodObject<{
        response: z.ZodObject<{
            accessToken: z.ZodString;
        }, z.core.$strip>;
    }, z.core.$strip>;
    type RequestBody = z.infer<typeof RequestBodySchema>;
    type Response = z.infer<typeof ResponseSchema>;
}
//# sourceMappingURL=register.command.d.ts.map