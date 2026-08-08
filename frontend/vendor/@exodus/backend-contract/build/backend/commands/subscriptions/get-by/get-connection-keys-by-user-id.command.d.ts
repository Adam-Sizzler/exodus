import { z } from 'zod';
export declare namespace GetConnectionKeysByUserIdCommand {
    const url: (userId: string) => string;
    const TSQ_url: string;
    const endpointDetails: import("../../../constants").EndpointDetails;
    const RequestParamSchema: z.ZodObject<{
        userId: z.ZodCoercedNumber<unknown>;
    }, z.core.$strip>;
    const ResponseSchema: z.ZodObject<{
        response: z.ZodObject<{
            enabledKeys: z.ZodArray<z.ZodString>;
            hiddenKeys: z.ZodArray<z.ZodString>;
            disabledKeys: z.ZodArray<z.ZodString>;
        }, z.core.$strip>;
    }, z.core.$strip>;
    type RequestParam = z.infer<typeof RequestParamSchema>;
    type Response = z.infer<typeof ResponseSchema>;
}
//# sourceMappingURL=get-connection-keys-by-user-id.command.d.ts.map