import { z } from 'zod';
export declare namespace GetPasskeyRegistrationOptionsCommand {
    const url: "/api/passkeys/registration/options";
    const TSQ_url: "/api/passkeys/registration/options";
    const endpointDetails: import("../../constants").EndpointDetails;
    const ResponseSchema: z.ZodObject<{
        response: z.ZodUnknown;
    }, "strip", z.ZodTypeAny, {
        response?: unknown;
    }, {
        response?: unknown;
    }>;
    type Response = z.infer<typeof ResponseSchema>;
}
//# sourceMappingURL=get-registration-options.command.d.ts.map