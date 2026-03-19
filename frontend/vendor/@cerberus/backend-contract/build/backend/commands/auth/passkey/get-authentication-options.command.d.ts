import { z } from 'zod';
export declare namespace GetPasskeyAuthenticationOptionsCommand {
    const url: "/api/auth/passkey/authentication/options";
    const TSQ_url: "/api/auth/passkey/authentication/options";
    const endpointDetails: import("../../../constants").EndpointDetails;
    const ResponseSchema: z.ZodObject<{
        response: z.ZodUnknown;
    }, "strip", z.ZodTypeAny, {
        response?: unknown;
    }, {
        response?: unknown;
    }>;
    type Response = z.infer<typeof ResponseSchema>;
}
//# sourceMappingURL=get-authentication-options.command.d.ts.map