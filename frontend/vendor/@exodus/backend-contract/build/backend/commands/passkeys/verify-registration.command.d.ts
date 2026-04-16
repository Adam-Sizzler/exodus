import { z } from 'zod';
export declare namespace VerifyPasskeyRegistrationCommand {
    const url: "/api/passkeys/registration/verify";
    const TSQ_url: "/api/passkeys/registration/verify";
    const endpointDetails: import("../../constants").EndpointDetails;
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
            verified: z.ZodBoolean;
        }, "strip", z.ZodTypeAny, {
            verified: boolean;
        }, {
            verified: boolean;
        }>;
    }, "strip", z.ZodTypeAny, {
        response: {
            verified: boolean;
        };
    }, {
        response: {
            verified: boolean;
        };
    }>;
    type Response = z.infer<typeof ResponseSchema>;
}
//# sourceMappingURL=verify-registration.command.d.ts.map