import { z } from 'zod';
export declare namespace DeletePasskeyCommand {
    const url: "/api/passkeys/";
    const TSQ_url: "/api/passkeys/";
    const endpointDetails: import("../../constants").EndpointDetails;
    const RequestBodySchema: z.ZodObject<{
        id: z.ZodString;
    }, z.core.$strip>;
    type RequestBody = z.infer<typeof RequestBodySchema>;
}
//# sourceMappingURL=delete-passkey.command.d.ts.map