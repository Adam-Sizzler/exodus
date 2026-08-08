import { z } from 'zod';
export declare namespace GetPasskeyAuthenticationOptionsCommand {
    const url: "/api/auth/passkey/authentication/options";
    const TSQ_url: "/api/auth/passkey/authentication/options";
    const endpointDetails: import("../../../constants").EndpointDetails;
    const ResponseSchema: z.ZodObject<{
        response: z.ZodUnknown;
    }, z.core.$strip>;
    type Response = z.infer<typeof ResponseSchema>;
}
//# sourceMappingURL=get-authentication-options.command.d.ts.map