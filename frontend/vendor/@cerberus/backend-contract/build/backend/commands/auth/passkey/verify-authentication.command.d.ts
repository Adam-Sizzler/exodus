import { z } from 'zod';
export declare namespace VerifyPasskeyAuthenticationCommand {
    const url: "/api/auth/passkey/authentication/verify";
    const TSQ_url: "/api/auth/passkey/authentication/verify";
    const endpointDetails: import("../../../constants").EndpointDetails;
    const RequestSchema: z.ZodObject<{
        response: z.ZodUnknown;
    }, "strip", z.ZodTypeAny, {
        response?: unknown;
    }, {
        response?: unknown;
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
//# sourceMappingURL=verify-authentication.command.d.ts.map