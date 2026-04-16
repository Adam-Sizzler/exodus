import { z } from 'zod';
export declare namespace GetConnectionKeysByUuidCommand {
    const url: (uuid: string) => string;
    const TSQ_url: string;
    const endpointDetails: import("../../../constants").EndpointDetails;
    const RequestSchema: z.ZodObject<{
        uuid: z.ZodString;
    }, "strip", z.ZodTypeAny, {
        uuid: string;
    }, {
        uuid: string;
    }>;
    type Request = z.infer<typeof RequestSchema>;
    const ResponseSchema: z.ZodObject<{
        response: z.ZodObject<{
            enabledKeys: z.ZodArray<z.ZodString, "many">;
            hiddenKeys: z.ZodArray<z.ZodString, "many">;
            disabledKeys: z.ZodArray<z.ZodString, "many">;
        }, "strip", z.ZodTypeAny, {
            enabledKeys: string[];
            hiddenKeys: string[];
            disabledKeys: string[];
        }, {
            enabledKeys: string[];
            hiddenKeys: string[];
            disabledKeys: string[];
        }>;
    }, "strip", z.ZodTypeAny, {
        response: {
            enabledKeys: string[];
            hiddenKeys: string[];
            disabledKeys: string[];
        };
    }, {
        response: {
            enabledKeys: string[];
            hiddenKeys: string[];
            disabledKeys: string[];
        };
    }>;
    type Response = z.infer<typeof ResponseSchema>;
}
//# sourceMappingURL=get-connection-keys-by-uuid.command.d.ts.map