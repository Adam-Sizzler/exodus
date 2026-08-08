import { z } from 'zod';
export declare namespace VerifyPasskeyRegistrationCommand {
    const url: "/api/passkeys/registration/verify";
    const TSQ_url: "/api/passkeys/registration/verify";
    const endpointDetails: import("../../constants").EndpointDetails;
    const RequestBodySchema: z.ZodObject<{
        response: z.ZodUnknown;
    }, z.core.$strip>;
    const ResponseSchema: z.ZodObject<{
        response: z.ZodObject<{
            verified: z.ZodBoolean;
        }, z.core.$strip>;
    }, z.core.$strip>;
    type RequestBody = z.infer<typeof RequestBodySchema>;
    type Response = z.infer<typeof ResponseSchema>;
}
//# sourceMappingURL=verify-registration.command.d.ts.map