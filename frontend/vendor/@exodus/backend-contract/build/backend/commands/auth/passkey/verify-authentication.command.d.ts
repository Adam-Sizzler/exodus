import { z } from 'zod';
export declare namespace VerifyPasskeyAuthenticationCommand {
    const url: "/api/auth/passkey/authentication/verify";
    const TSQ_url: "/api/auth/passkey/authentication/verify";
    const endpointDetails: import("../../../constants").EndpointDetails;
    const RequestBodySchema: z.ZodObject<{
        response: z.ZodUnknown;
    }, z.core.$strip>;
    const ResponseSchema: z.ZodObject<{
        response: z.ZodObject<{
            accessToken: z.ZodString;
        }, z.core.$strip>;
    }, z.core.$strip>;
    type RequestBody = z.infer<typeof RequestBodySchema>;
    type Response = z.infer<typeof ResponseSchema>;
}
//# sourceMappingURL=verify-authentication.command.d.ts.map